package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand/v2"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"palabrario/core-service/internal/database"
	"palabrario/core-service/internal/db"
	"palabrario/pkg/models"
)

const roomCodeCharset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

func generateRoomCode() string {
	b := make([]byte, 6)
	for i := range b {
		b[i] = roomCodeCharset[rand.IntN(len(roomCodeCharset))]
	}
	return string(b)
}

type ArenaCatalogService interface {
	PrepareMatch(ctx context.Context, hostUserID string, hostRole models.AccountRole, req models.CreateMatchRequest) (*models.MatchPackage, error)
	GetMatchPackageByRoomCode(ctx context.Context, roomCode string) (*models.MatchPackage, error)
	ConsolidateMatch(ctx context.Context, req models.ConsolidateMatchRequest) (*models.MatchSummary, error)
	ListMatchHistory(ctx context.Context, userID string, limit int) ([]models.MatchSummary, error)
	GetMatchReplay(ctx context.Context, viewerID string, viewerRole models.AccountRole, matchID string) (*models.MatchReplayDetail, error)
}

type arenaCatalogService struct {
	databaseInstance *database.DB
	queries          *db.Queries
}

func NewArenaCatalogService(databaseInstance *database.DB) ArenaCatalogService {
	return &arenaCatalogService{
		databaseInstance: databaseInstance,
		queries:          db.New(databaseInstance.Pool),
	}
}

func (s *arenaCatalogService) PrepareMatch(ctx context.Context, hostUserID string, hostRole models.AccountRole, req models.CreateMatchRequest) (*models.MatchPackage, error) {
	if req.IDDeck <= 0 {
		return nil, NewDomainError(ErrInvalidInput, "El identificador de baraja debe ser positivo")
	}

	totalRounds := req.TotalRounds
	if totalRounds == 0 {
		totalRounds = 8
	}
	if totalRounds < 8 || totalRounds > 20 {
		return nil, NewDomainError(ErrInvalidInput, "El número de rondas debe estar entre 8 y 20")
	}

	roundTimeSec := req.RoundTimeSec
	if roundTimeSec == 0 {
		roundTimeSec = 15
	}
	if roundTimeSec < 10 || roundTimeSec > 30 {
		return nil, NewDomainError(ErrInvalidInput, "El tiempo por ronda debe estar entre 10 y 30 segundos")
	}

	isPublic := true
	if req.IsPublic != nil {
		isPublic = *req.IsPublic
	}

	hostUUID, err := ParseDomainUUID(hostUserID, "ID de anfitrión inválido")
	if err != nil {
		return nil, err
	}

	// 1. Obtener la baraja y verificar privacidad (RN04)
	deck, err := s.queries.GetDeckByID(ctx, int32(req.IDDeck))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, NewDomainError(ErrNotFound, "Baraja no encontrada")
		}
		return nil, MapDBError(err, "Error al consultar la baraja")
	}

	if !deck.IsPublic && deck.IDUser.String() != hostUserID && hostRole != models.RoleAdmin {
		return nil, NewDomainError(ErrForbidden, "No tienes permiso para abrir una sala con una baraja privada ajena")
	}

	// 2. Obtener sentidos aprobados de la baraja
	senses, err := s.queries.GetApprovedDeckSensesForArena(ctx, int32(req.IDDeck))
	if err != nil {
		return nil, MapDBError(err, "Error al consultar los sentidos de la baraja")
	}

	// 3. Regla estricta de límites de Arena (CU14 1a)
	if len(senses) < 8 || len(senses) < totalRounds {
		return nil, NewDomainError(ErrUnprocessable, "La baraja debe contener al menos 8 sentidos aprobados y al menos tantos sentidos como el número de rondas configuradas")
	}

	// 4. Barajar aleatoriamente y seleccionar totalRounds sentidos únicos
	shuffledSenses := make([]db.GetApprovedDeckSensesForArenaRow, len(senses))
	copy(shuffledSenses, senses)
	rand.Shuffle(len(shuffledSenses), func(i, j int) {
		shuffledSenses[i], shuffledSenses[j] = shuffledSenses[j], shuffledSenses[i]
	})
	selectedSenses := shuffledSenses[:totalRounds]

	selectedIDs := make([]int32, len(selectedSenses))
	for i, sense := range selectedSenses {
		selectedIDs[i] = sense.IDSense
	}

	// 5. Consultar preguntas curadas para los sentidos seleccionados
	curatedQuestions, err := s.queries.GetCuratedQuestionsBySenseIDs(ctx, selectedIDs)
	if err != nil {
		return nil, MapDBError(err, "Error al consultar preguntas curadas")
	}

	curatedMap := make(map[int32]db.Question)
	for _, q := range curatedQuestions {
		// Mantener la primera pregunta curada por sentido si hubiera múltiples
		if _, exists := curatedMap[q.IDSense]; !exists {
			curatedMap[q.IDSense] = q
		}
	}

	questions := make([]models.InternalArenaQuestion, totalRounds)
	for i, sense := range selectedSenses {
		roundNum := i + 1
		q, hasCurated := curatedMap[sense.IDSense]

		var (
			modality      string
			prompt        string
			correctOption string
			explanation   string
			idQuestion    *int
			options       []string
		)

		var customDistractors []string
		if hasCurated && len(q.CustomDistractors) > 0 {
			_ = json.Unmarshal(q.CustomDistractors, &customDistractors)
		}

		if hasCurated && len(customDistractors) == 3 {
			modality = string(q.Type)
			prompt = q.Prompt
			if q.Type == db.QuestionTypeDEFINITIONTOWORD {
				correctOption = sense.Lemma
			} else {
				correctOption = sense.ShortDefinition
			}
			explanation = q.Explanation
			qID := int(q.IDQuestion)
			idQuestion = &qID

			opts := append([]string{correctOption}, customDistractors...)
			rand.Shuffle(len(opts), func(a, b int) { opts[a], opts[b] = opts[b], opts[a] })
			options = opts
		} else {
			// Fallback dinámico obligatorio (WORD_TO_DEFINITION)
			modality = "WORD_TO_DEFINITION"
			prompt = fmt.Sprintf("¿Cuál es la definición correcta de «%s» (%s %s)?", sense.Lemma, sense.FlagEmoji, sense.RegionName)
			correctOption = sense.ShortDefinition
			explanation = fmt.Sprintf("En %s (%s), «%s» (%s) significa: %s", sense.RegionName, sense.RegionCode, sense.Lemma, sense.PartOfSpeech, sense.Definition)
			idQuestion = nil

			// Seleccionar 3 distractores distintos de la misma baraja
			candidates := make([]string, 0, len(senses))
			seenDefs := make(map[string]struct{})
			seenDefs[correctOption] = struct{}{}

			for _, s := range senses {
				if s.IDSense == sense.IDSense {
					continue
				}
				if _, seen := seenDefs[s.ShortDefinition]; !seen {
					seenDefs[s.ShortDefinition] = struct{}{}
					candidates = append(candidates, s.ShortDefinition)
				}
			}

			rand.Shuffle(len(candidates), func(a, b int) { candidates[a], candidates[b] = candidates[b], candidates[a] })
			distractors := candidates
			if len(distractors) > 3 {
				distractors = distractors[:3]
			}

			opts := append([]string{correctOption}, distractors...)
			rand.Shuffle(len(opts), func(a, b int) { opts[a], opts[b] = opts[b], opts[a] })
			options = opts
		}

		questions[i] = models.InternalArenaQuestion{
			RoundNumber:   roundNum,
			IDSense:       int(sense.IDSense),
			IDQuestion:    idQuestion,
			Lemma:         sense.Lemma,
			Modality:      modality,
			Prompt:        prompt,
			Options:       options,
			CorrectOption: correctOption,
			Explanation:   explanation,
			TimeLimitSecs: roundTimeSec,
		}
	}

	// 6. Transacción ACID para crear la partida y persistir las rondas en 3NF
	tx, err := s.databaseInstance.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("error al iniciar transacción de preparación de sala: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := s.queries.WithTx(tx)

	var match db.Match
	var roomCode string
	for attempts := 0; attempts < 5; attempts++ {
		roomCode = generateRoomCode()
		match, err = qtx.CreateMatch(ctx, db.CreateMatchParams{
			IDDeck:       int32(req.IDDeck),
			HostUserID:   UUIDToPgtype(hostUUID),
			RoomCode:     roomCode,
			IsPublic:     isPublic,
			TotalRounds:  int32(totalRounds),
			RoundTimeSec: int32(roundTimeSec),
		})
		if err == nil {
			break
		}
	}
	if err != nil {
		return nil, MapDBError(err, "Error al crear la partida en base de datos")
	}

	for _, q := range questions {
		var qID pgtype.Int4
		if q.IDQuestion != nil {
			qID = pgtype.Int4{Int32: int32(*q.IDQuestion), Valid: true}
		}

		err = qtx.InsertMatchRound(ctx, db.InsertMatchRoundParams{
			IDMatch:       match.IDMatch,
			RoundNumber:   int16(q.RoundNumber),
			IDSense:       int32(q.IDSense),
			IDQuestion:    qID,
			Modality:      db.QuestionType(q.Modality),
			Prompt:        q.Prompt,
			OptionA:       q.Options[0],
			OptionB:       q.Options[1],
			OptionC:       q.Options[2],
			OptionD:       q.Options[3],
			CorrectOption: q.CorrectOption,
			Explanation:   q.Explanation,
		})
		if err != nil {
			return nil, MapDBError(err, "Error al insertar la ronda de la partida")
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("error al confirmar transacción de preparación de sala: %w", err)
	}

	return &models.MatchPackage{
		IDMatch:      match.IDMatch.String(),
		RoomCode:     match.RoomCode,
		HostUserID:   hostUserID,
		IDDeck:       req.IDDeck,
		DeckTitle:    deck.Title,
		IsPublic:     match.IsPublic,
		TotalRounds:  int(match.TotalRounds),
		RoundTimeSec: int(match.RoundTimeSec),
		Questions:    questions,
	}, nil
}

func (s *arenaCatalogService) GetMatchPackageByRoomCode(ctx context.Context, roomCode string) (*models.MatchPackage, error) {
	code := strings.ToUpper(strings.TrimSpace(roomCode))
	if code == "" {
		return nil, NewDomainError(ErrInvalidInput, "El código de sala es requerido")
	}

	match, err := s.queries.GetMatchByRoomCode(ctx, code)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, NewDomainError(ErrNotFound, "Partida no encontrada para el código especificado")
		}
		return nil, MapDBError(err, "Error al consultar la partida por código de sala")
	}

	if match.Status == db.MatchStatusFINISHED {
		return nil, NewDomainError(ErrConflict, "La partida ya ha finalizado")
	}

	rounds, err := s.queries.ListMatchRoundsByMatchID(ctx, match.IDMatch)
	if err != nil {
		return nil, MapDBError(err, "Error al consultar las rondas de la partida")
	}

	questions := make([]models.InternalArenaQuestion, len(rounds))
	for i, r := range rounds {
		var qID *int
		if r.IDQuestion.Valid {
			id := int(r.IDQuestion.Int32)
			qID = &id
		}

		questions[i] = models.InternalArenaQuestion{
			RoundNumber:   int(r.RoundNumber),
			IDSense:       int(r.IDSense),
			IDQuestion:    qID,
			Lemma:         r.Lemma,
			Modality:      string(r.Modality),
			Prompt:        r.Prompt,
			Options:       []string{r.OptionA, r.OptionB, r.OptionC, r.OptionD},
			CorrectOption: r.CorrectOption,
			Explanation:   r.Explanation,
			TimeLimitSecs: int(match.RoundTimeSec),
		}
	}

	hostUserIDStr := ""
	if match.HostUserID.Valid {
		hostUserIDStr = uuid.UUID(match.HostUserID.Bytes).String()
	}

	return &models.MatchPackage{
		IDMatch:      match.IDMatch.String(),
		RoomCode:     match.RoomCode,
		HostUserID:   hostUserIDStr,
		IDDeck:       int(match.IDDeck),
		DeckTitle:    match.DeckTitle,
		IsPublic:     match.IsPublic,
		TotalRounds:  int(match.TotalRounds),
		RoundTimeSec: int(match.RoundTimeSec),
		Questions:    questions,
	}, nil
}

func (s *arenaCatalogService) ConsolidateMatch(ctx context.Context, req models.ConsolidateMatchRequest) (*models.MatchSummary, error) {
	matchUUID, err := ParseDomainUUID(req.IDMatch, "ID de partida inválido")
	if err != nil {
		return nil, err
	}

	match, err := s.queries.GetMatchByID(ctx, matchUUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, NewDomainError(ErrNotFound, "Partida no encontrada")
		}
		return nil, MapDBError(err, "Error al consultar la partida para consolidación")
	}

	if match.Status == db.MatchStatusFINISHED {
		return nil, NewDomainError(ErrConflict, "Esta partida ya fue consolidada previamente")
	}

	// Criterio determinístico de desempate:
	// 1. Mayor FinalScore DESC
	// 2. Mayor len(CorrectSenseIDs) DESC
	// 3. Menor TotalTimeSpentMs ASC
	participants := make([]models.ParticipantFinalResult, len(req.Participants))
	copy(participants, req.Participants)

	sort.SliceStable(participants, func(i, j int) bool {
		if participants[i].FinalScore != participants[j].FinalScore {
			return participants[i].FinalScore > participants[j].FinalScore
		}
		if len(participants[i].CorrectSenseIDs) != len(participants[j].CorrectSenseIDs) {
			return len(participants[i].CorrectSenseIDs) > len(participants[j].CorrectSenseIDs)
		}
		return participants[i].TotalTimeSpentMs < participants[j].TotalTimeSpentMs
	})

	for i := range participants {
		participants[i].FinalRank = i + 1
	}

	var winnerUUID pgtype.UUID
	if len(participants) > 0 {
		winnerParsed, pErr := ParseDomainUUID(participants[0].IDUser, "ID de ganador inválido")
		if pErr != nil {
			return nil, pErr
		}
		winnerUUID = UUIDToPgtype(winnerParsed)
	}

	tx, err := s.databaseInstance.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("error al iniciar transacción de consolidación: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := s.queries.WithTx(tx)

	// 1. Marcar partida como finalizada
	updatedMatch, err := qtx.MarkMatchFinished(ctx, db.MarkMatchFinishedParams{
		IDMatch:      matchUUID,
		WinnerUserID: winnerUUID,
	})
	if err != nil {
		return nil, MapDBError(err, "Error al finalizar la partida")
	}

	now := time.Now().UTC()

	// 2. Procesar a cada participante
	for _, p := range participants {
		pUUID, pErr := ParseDomainUUID(p.IDUser, "ID de usuario participante inválido")
		if pErr != nil {
			return nil, pErr
		}

		var xpEarned int
		var xpReason db.XpReason
		switch p.FinalRank {
		case 1:
			xpEarned = 150
			xpReason = db.XpReasonMATCHWIN
		case 2:
			xpEarned = 100
			xpReason = db.XpReasonQUIZCORRECT
		case 3:
			xpEarned = 75
			xpReason = db.XpReasonQUIZCORRECT
		default:
			xpEarned = 40
			xpReason = db.XpReasonQUIZCORRECT
		}

		// Insertar/actualizar participante
		err = qtx.UpsertMatchParticipant(ctx, db.UpsertMatchParticipantParams{
			IDMatch:          matchUUID,
			IDUser:           pUUID,
			Score:            int32(p.FinalScore),
			FinalRank:        int32(p.FinalRank),
			CorrectAnswers:   int32(len(p.CorrectSenseIDs)),
			TotalTimeSpentMs: int32(p.TotalTimeSpentMs),
			XpEarned:         int32(xpEarned),
		})
		if err != nil {
			return nil, MapDBError(err, "Error al persistir participante de la partida")
		}

		// Persistir desglose ronda por ronda en quiz_attempts para Replay
		for _, att := range p.RoundAttempts {
			qRating := int16(1)
			if att.IsCorrect {
				qRating = 3
			}

			var qID pgtype.Int4
			if att.IDQuestion != nil {
				qID = pgtype.Int4{Int32: int32(*att.IDQuestion), Valid: true}
			}

			err = qtx.UpsertMatchRoundQuizAttempt(ctx, db.UpsertMatchRoundQuizAttemptParams{
				IDUser:         pUUID,
				IDSense:        int32(att.IDSense),
				IDQuestion:     qID,
				IDMatch:        UUIDToPgtype(matchUUID),
				RoundNumber:    pgtype.Int2{Int16: int16(att.RoundNumber), Valid: true},
				SelectedOption: att.SelectedOption,
				IsCorrect:      att.IsCorrect,
				QualityRating:  pgtype.Int2{Int16: qRating, Valid: true},
				TimeSpentMs:    int32(att.TimeSpentMs),
				PointsEarned:   int32(att.PointsEarned),
			})
			if err != nil {
				return nil, MapDBError(err, "Error al guardar el intento de ronda en quiz_attempts")
			}
		}

		// Asiento en xp_transactions y actualización de racha/saldo
		_, err = qtx.InsertXPTransaction(ctx, db.InsertXPTransactionParams{
			IDUser:     pUUID,
			AmountXp:   int32(xpEarned),
			Reason:     xpReason,
			MetadataID: StringToText(matchUUID.String()),
		})
		if err != nil {
			return nil, MapDBError(err, "Error al registrar transacción de XP")
		}

		_, err = qtx.AddUserXPAndUpdateStreak(ctx, db.AddUserXPAndUpdateStreakParams{
			AmountXp: int32(xpEarned),
			IDUser:   pUUID,
		})
		if err != nil {
			return nil, MapDBError(err, "Error al actualizar la experiencia y racha del usuario")
		}

		// Actualizar maestría SM-2 (RN09)
		for _, senseID := range p.CorrectSenseIDs {
			var prevMastery *models.UserMastery
			prevRow, err := qtx.GetUserMasteryForSense(ctx, db.GetUserMasteryForSenseParams{
				IDUser:  pUUID,
				IDSense: int32(senseID),
			})
			if err == nil {
				prevMastery = &models.UserMastery{
					IDUser:              prevRow.IDUser.String(),
					IDSense:             int(prevRow.IDSense),
					ProgressStage:       models.MasteryStage(prevRow.ProgressStage),
					MasteryPercentage:   prevRow.MasteryPercentage,
					ConsecutiveCorrect:  int(prevRow.ConsecutiveCorrect),
					TimesReviewed:       int(prevRow.TimesReviewed),
					EasinessFactor:      prevRow.EasinessFactor,
					CurrentIntervalDays: int(prevRow.CurrentIntervalDays),
					NextReviewAt:        prevRow.NextReviewAt.Time,
					LastReviewedAt:      prevRow.LastReviewedAt.Time,
				}
			}

			nextMastery, _ := CalculateSM2(prevMastery, 3, now)
			_, err = qtx.UpsertUserMastery(ctx, db.UpsertUserMasteryParams{
				IDUser:              pUUID,
				IDSense:             int32(senseID),
				ProgressStage:       db.MasteryStage(nextMastery.ProgressStage),
				MasteryPercentage:   nextMastery.MasteryPercentage,
				ConsecutiveCorrect:  int32(nextMastery.ConsecutiveCorrect),
				TimesReviewed:       int32(nextMastery.TimesReviewed),
				EasinessFactor:      nextMastery.EasinessFactor,
				CurrentIntervalDays: int32(nextMastery.CurrentIntervalDays),
				NextReviewAt:        pgtype.Timestamptz{Time: nextMastery.NextReviewAt, Valid: true},
				LastReviewedAt:      pgtype.Timestamptz{Time: nextMastery.LastReviewedAt, Valid: true},
			})
			if err != nil {
				return nil, MapDBError(err, "Error al actualizar la maestría del sentido acertado")
			}
		}

		for _, senseID := range p.WrongSenseIDs {
			var prevMastery *models.UserMastery
			prevRow, err := qtx.GetUserMasteryForSense(ctx, db.GetUserMasteryForSenseParams{
				IDUser:  pUUID,
				IDSense: int32(senseID),
			})
			if err == nil {
				prevMastery = &models.UserMastery{
					IDUser:              prevRow.IDUser.String(),
					IDSense:             int(prevRow.IDSense),
					ProgressStage:       models.MasteryStage(prevRow.ProgressStage),
					MasteryPercentage:   prevRow.MasteryPercentage,
					ConsecutiveCorrect:  int(prevRow.ConsecutiveCorrect),
					TimesReviewed:       int(prevRow.TimesReviewed),
					EasinessFactor:      prevRow.EasinessFactor,
					CurrentIntervalDays: int(prevRow.CurrentIntervalDays),
					NextReviewAt:        prevRow.NextReviewAt.Time,
					LastReviewedAt:      prevRow.LastReviewedAt.Time,
				}
			}

			nextMastery, _ := CalculateSM2(prevMastery, 1, now)
			_, err = qtx.UpsertUserMastery(ctx, db.UpsertUserMasteryParams{
				IDUser:              pUUID,
				IDSense:             int32(senseID),
				ProgressStage:       db.MasteryStage(nextMastery.ProgressStage),
				MasteryPercentage:   nextMastery.MasteryPercentage,
				ConsecutiveCorrect:  int32(nextMastery.ConsecutiveCorrect),
				TimesReviewed:       int32(nextMastery.TimesReviewed),
				EasinessFactor:      nextMastery.EasinessFactor,
				CurrentIntervalDays: int32(nextMastery.CurrentIntervalDays),
				NextReviewAt:        pgtype.Timestamptz{Time: nextMastery.NextReviewAt, Valid: true},
				LastReviewedAt:      pgtype.Timestamptz{Time: nextMastery.LastReviewedAt, Valid: true},
			})
			if err != nil {
				return nil, MapDBError(err, "Error al actualizar la maestría del sentido errado")
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("error al confirmar transacción de consolidación: %w", err)
	}

	// 3. Consultar podio resultante
	participantRows, err := s.queries.ListMatchParticipantsByMatchIDs(ctx, []uuid.UUID{matchUUID})
	if err != nil {
		return nil, MapDBError(err, "Error al consultar los participantes consolidados")
	}

	podium := make([]models.ArenaParticipant, len(participantRows))
	for i, r := range participantRows {
		podium[i] = models.ArenaParticipant{
			IDUser:      r.IDUser.String(),
			UserName:    r.UserName,
			Name:        r.Name,
			AvatarURL:   r.AvatarUrl,
			Score:       int(r.Score),
			Rank:        int(r.FinalRank),
			IsHost:      r.IsHost,
			IsReady:     true,
			DialectCode: r.AnchorDialect,
			FlagEmoji:   r.FlagEmoji,
		}
	}

	var winnerIDStr *string
	if updatedMatch.WinnerUserID.Valid {
		wStr := uuid.UUID(updatedMatch.WinnerUserID.Bytes).String()
		winnerIDStr = &wStr
	}

	var finishedAtStr string
	if updatedMatch.FinishedAt.Valid {
		finishedAtStr = updatedMatch.FinishedAt.Time.Format(time.RFC3339)
	}

	return &models.MatchSummary{
		IDMatch:      updatedMatch.IDMatch.String(),
		RoomCode:     updatedMatch.RoomCode,
		IDDeck:       int(updatedMatch.IDDeck),
		DeckTitle:    match.DeckTitle,
		WinnerUserID: winnerIDStr,
		Status:       models.MatchStatus(updatedMatch.Status),
		TotalRounds:  int(updatedMatch.TotalRounds),
		RoundTimeSec: int(updatedMatch.RoundTimeSec),
		Podium:       podium,
		FinishedAt:   finishedAtStr,
	}, nil
}

func (s *arenaCatalogService) ListMatchHistory(ctx context.Context, userID string, limit int) ([]models.MatchSummary, error) {
	uUUID, err := ParseDomainUUID(userID, "ID de usuario inválido")
	if err != nil {
		return nil, err
	}

	if limit <= 0 || limit > 50 {
		limit = 20
	}

	matches, err := s.queries.ListUserMatchHistory(ctx, db.ListUserMatchHistoryParams{
		IDUser: uUUID,
		Limit:  int32(limit),
	})
	if err != nil {
		return nil, MapDBError(err, "Error al listar el historial de partidas")
	}

	if len(matches) == 0 {
		return []models.MatchSummary{}, nil
	}

	matchIDs := make([]uuid.UUID, len(matches))
	for i, m := range matches {
		matchIDs[i] = m.IDMatch
	}

	participants, err := s.queries.ListMatchParticipantsByMatchIDs(ctx, matchIDs)
	if err != nil {
		return nil, MapDBError(err, "Error al consultar los participantes del historial")
	}

	participantsByMatch := make(map[uuid.UUID][]models.ArenaParticipant)
	for _, p := range participants {
		participantsByMatch[p.IDMatch] = append(participantsByMatch[p.IDMatch], models.ArenaParticipant{
			IDUser:      p.IDUser.String(),
			UserName:    p.UserName,
			Name:        p.Name,
			AvatarURL:   p.AvatarUrl,
			Score:       int(p.Score),
			Rank:        int(p.FinalRank),
			IsHost:      p.IsHost,
			IsReady:     true,
			DialectCode: p.AnchorDialect,
			FlagEmoji:   p.FlagEmoji,
		})
	}

	summaries := make([]models.MatchSummary, len(matches))
	for i, m := range matches {
		var winnerIDStr *string
		if m.WinnerUserID.Valid {
			w := uuid.UUID(m.WinnerUserID.Bytes).String()
			winnerIDStr = &w
		}

		var finishedAtStr string
		if m.FinishedAt.Valid {
			finishedAtStr = m.FinishedAt.Time.Format(time.RFC3339)
		}

		podium := participantsByMatch[m.IDMatch]
		if podium == nil {
			podium = []models.ArenaParticipant{}
		}

		summaries[i] = models.MatchSummary{
			IDMatch:      m.IDMatch.String(),
			RoomCode:     m.RoomCode,
			IDDeck:       int(m.IDDeck),
			DeckTitle:    m.DeckTitle,
			WinnerUserID: winnerIDStr,
			Status:       models.MatchStatus(m.Status),
			TotalRounds:  int(m.TotalRounds),
			RoundTimeSec: int(m.RoundTimeSec),
			Podium:       podium,
			FinishedAt:   finishedAtStr,
		}
	}

	return summaries, nil
}

func (s *arenaCatalogService) GetMatchReplay(ctx context.Context, viewerID string, viewerRole models.AccountRole, matchID string) (*models.MatchReplayDetail, error) {
	matchUUID, err := ParseDomainUUID(matchID, "ID de partida inválido")
	if err != nil {
		return nil, err
	}

	match, err := s.queries.GetMatchByID(ctx, matchUUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, NewDomainError(ErrNotFound, "Partida no encontrada")
		}
		return nil, MapDBError(err, "Error al consultar la partida para replay")
	}

	// 1. Participantes
	participants, err := s.queries.ListMatchParticipantsByMatchIDs(ctx, []uuid.UUID{matchUUID})
	if err != nil {
		return nil, MapDBError(err, "Error al consultar participantes para replay")
	}

	// Validación de autorización del espectador (CU16):
	// - Si viewerRole == RoleAdmin o match.IsPublic, el acceso está permitido.
	// - Si es privada (!match.IsPublic), viewerID debe estar presente en la lista de participants.
	if viewerRole != models.RoleAdmin && !match.IsPublic {
		isParticipant := false
		for _, p := range participants {
			if p.IDUser.String() == viewerID {
				isParticipant = true
				break
			}
		}
		if !isParticipant {
			return nil, NewDomainError(ErrForbidden, "No tienes permiso para consultar el replay de esta partida privada")
		}
	}

	podium := make([]models.ArenaParticipant, len(participants))
	for i, p := range participants {
		podium[i] = models.ArenaParticipant{
			IDUser:      p.IDUser.String(),
			UserName:    p.UserName,
			Name:        p.Name,
			AvatarURL:   p.AvatarUrl,
			Score:       int(p.Score),
			Rank:        int(p.FinalRank),
			IsHost:      p.IsHost,
			IsReady:     true,
			DialectCode: p.AnchorDialect,
			FlagEmoji:   p.FlagEmoji,
		}
	}

	// 2. Rondas en 3NF
	rounds, err := s.queries.ListMatchRoundsByMatchID(ctx, matchUUID)
	if err != nil {
		return nil, MapDBError(err, "Error al consultar rondas para replay")
	}

	// 3. Intentos de ronda
	attempts, err := s.queries.ListMatchQuizAttemptsByMatchID(ctx, UUIDToPgtype(matchUUID))
	if err != nil {
		return nil, MapDBError(err, "Error al consultar intentos de quiz para replay")
	}

	attemptsByRound := make(map[int16][]models.ParticipantRoundAnswer)
	for _, att := range attempts {
		rNum := att.RoundNumber.Int16
		attemptsByRound[rNum] = append(attemptsByRound[rNum], models.ParticipantRoundAnswer{
			IDUser:         att.IDUser.String(),
			UserName:       att.UserName,
			Name:           att.Name,
			AvatarURL:      att.AvatarUrl,
			DialectCode:    att.AnchorDialect,
			FlagEmoji:      att.FlagEmoji,
			SelectedOption: att.SelectedOption,
			IsCorrect:      att.IsCorrect,
			TimeSpentMs:    int(att.TimeSpentMs),
			PointsEarned:   int(att.PointsEarned),
		})
	}

	replayRounds := make([]models.MatchRoundReplay, len(rounds))
	for i, r := range rounds {
		answers := attemptsByRound[r.RoundNumber]
		if answers == nil {
			answers = []models.ParticipantRoundAnswer{}
		}

		replayRounds[i] = models.MatchRoundReplay{
			RoundNumber:   int(r.RoundNumber),
			IDSense:       int(r.IDSense),
			Lemma:         r.Lemma,
			Modality:      string(r.Modality),
			Prompt:        r.Prompt,
			Options:       []string{r.OptionA, r.OptionB, r.OptionC, r.OptionD},
			CorrectOption: r.CorrectOption,
			Explanation:   r.Explanation,
			Answers:       answers,
		}
	}

	var winnerIDStr *string
	if match.WinnerUserID.Valid {
		w := uuid.UUID(match.WinnerUserID.Bytes).String()
		winnerIDStr = &w
	}

	var finishedAtStr string
	if match.FinishedAt.Valid {
		finishedAtStr = match.FinishedAt.Time.Format(time.RFC3339)
	}

	return &models.MatchReplayDetail{
		MatchSummary: models.MatchSummary{
			IDMatch:      match.IDMatch.String(),
			RoomCode:     match.RoomCode,
			IDDeck:       int(match.IDDeck),
			DeckTitle:    match.DeckTitle,
			WinnerUserID: winnerIDStr,
			Status:       models.MatchStatus(match.Status),
			TotalRounds:  int(match.TotalRounds),
			RoundTimeSec: int(match.RoundTimeSec),
			Podium:       podium,
			FinishedAt:   finishedAtStr,
		},
		Rounds: replayRounds,
	}, nil
}
