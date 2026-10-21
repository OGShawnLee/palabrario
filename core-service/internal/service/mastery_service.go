package service

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"palabrario/core-service/internal/database"
	"palabrario/core-service/internal/db"
	"palabrario/pkg/models"
)

type MasteryService interface {
	GetStudySession(ctx context.Context, userID string, userRole models.AccountRole, deckID *int, limit int) (*models.StudySessionPayload, error)
	SubmitReview(ctx context.Context, userID string, req models.SubmitReviewRequest) (*models.CardReviewResult, error)
	GetMasteryDashboard(ctx context.Context, userID string) (*models.MasteryDashboard, error)
	ListMasteryItems(ctx context.Context, userID string, stage *string) ([]models.StudyCardItem, error)
	ListXPLedger(ctx context.Context, userID string, limit int) ([]models.XPTransactionItem, error)
}

type masteryService struct {
	databaseInstance *database.DB
	queries          *db.Queries
}

func NewMasteryService(databaseInstance *database.DB) MasteryService {
	return &masteryService{
		databaseInstance: databaseInstance,
		queries:          db.New(databaseInstance.Pool),
	}
}

// CalculateSM2 implementa la función determinística pura del algoritmo SuperMemo SM-2
// siguiendo exactamente las especificaciones de la sección 2.6.2 de .context/INTERFACES.md.
func CalculateSM2(prev *models.UserMastery, qualityRating int, now time.Time) (models.UserMastery, int) {
	// 1. Mapeo de Escala UI (1..4) a Escala SM-2 (1..5):
	// 1 (Otra vez / Fallo) -> q = 1 (Incorrecto)
	// 2 (Difícil)          -> q = 3 (Correcto con esfuerzo alto)
	// 3 (Bueno)            -> q = 4 (Correcto con ritmo normal)
	// 4 (Fácil)            -> q = 5 (Recuerdo instantáneo perfecto)
	var q int
	switch qualityRating {
	case 1:
		q = 1
	case 2:
		q = 3
	case 3:
		q = 4
	case 4:
		q = 5
	default:
		q = 1
	}

	// Estado inicial si prev == nil
	ef := float32(2.5)
	consecutiveCorrect := 0
	timesReviewed := 0
	prevInterval := 1
	var senseID int
	var userID string

	if prev != nil {
		ef = prev.EasinessFactor
		consecutiveCorrect = prev.ConsecutiveCorrect
		timesReviewed = prev.TimesReviewed
		prevInterval = prev.CurrentIntervalDays
		senseID = prev.IDSense
		userID = prev.IDUser
	}

	timesReviewed++

	var currentIntervalDays int
	var progressStage models.MasteryStage
	var masteryPercentage float32
	var xpAwarded int

	if qualityRating == 1 { // q < 3 (Fallo)
		consecutiveCorrect = 0
		currentIntervalDays = 1
		ef = maxFloat32(1.3, ef-0.20)
		progressStage = models.StageLearning
		masteryPercentage = 20.0
		xpAwarded = 5
	} else { // q >= 3 (Acierto)
		consecutiveCorrect++
		if consecutiveCorrect == 1 {
			currentIntervalDays = 1
		} else if consecutiveCorrect == 2 {
			currentIntervalDays = 6
		} else {
			currentIntervalDays = int(math.Ceil(float64(prevInterval) * float64(ef)))
		}

		// EF' = max(1.3, EF + (0.1 - (5-q)*(0.08 + (5-q)*0.02)))
		diff := float32(5 - q)
		efAdjustment := 0.1 - diff*(0.08+diff*0.02)
		ef = maxFloat32(1.3, ef+efAdjustment)

		// Clasificación de etapa (ProgressStage)
		if consecutiveCorrect >= 5 && currentIntervalDays >= 21 {
			progressStage = models.StageMastered
			masteryPercentage = 100.0
		} else if consecutiveCorrect >= 3 && currentIntervalDays >= 6 {
			progressStage = models.StageReviewing
			masteryPercentage = 75.0
		} else {
			progressStage = models.StageLearning
			masteryPercentage = 40.0
		}

		if qualityRating == 4 {
			xpAwarded = 20
		} else {
			xpAwarded = 15
		}
	}

	nextReviewAt := now.AddDate(0, 0, currentIntervalDays)

	return models.UserMastery{
		IDUser:              userID,
		IDSense:             senseID,
		ProgressStage:       progressStage,
		MasteryPercentage:   masteryPercentage,
		ConsecutiveCorrect:  consecutiveCorrect,
		TimesReviewed:       timesReviewed,
		EasinessFactor:      ef,
		CurrentIntervalDays: currentIntervalDays,
		NextReviewAt:        nextReviewAt,
		LastReviewedAt:      now,
	}, xpAwarded
}

func maxFloat32(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

func (s *masteryService) SubmitReview(ctx context.Context, userID string, req models.SubmitReviewRequest) (*models.CardReviewResult, error) {
	if req.IDSense <= 0 {
		return nil, NewDomainError(ErrInvalidInput, "El identificador de sentido debe ser positivo")
	}
	if req.QualityRating < 1 || req.QualityRating > 4 {
		return nil, NewDomainError(ErrInvalidInput, "La calificación de calidad debe estar entre 1 y 4")
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, NewDomainError(ErrInvalidInput, "Identificador de usuario inválido")
	}

	// 1. Obtener estado previo de maestría si existe
	var prevMastery *models.UserMastery
	prevRow, err := s.queries.GetUserMasteryForSense(ctx, db.GetUserMasteryForSenseParams{
		IDUser:  userUUID,
		IDSense: int32(req.IDSense),
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

	now := time.Now().UTC()
	updatedMastery, xpAwarded := CalculateSM2(prevMastery, req.QualityRating, now)
	updatedMastery.IDUser = userID
	updatedMastery.IDSense = req.IDSense

	// 2. Iniciar transacción ACID pgx.Tx
	tx, err := s.databaseInstance.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("error al iniciar transacción de revisión: %w", err)
	}
	defer tx.Rollback(ctx)

	txQueries := s.queries.WithTx(tx)

	// 2.1. Upsert en user_mastery
	upsertedRow, err := txQueries.UpsertUserMastery(ctx, db.UpsertUserMasteryParams{
		IDUser:              userUUID,
		IDSense:             int32(req.IDSense),
		ProgressStage:       db.MasteryStage(updatedMastery.ProgressStage),
		MasteryPercentage:   updatedMastery.MasteryPercentage,
		ConsecutiveCorrect:  int32(updatedMastery.ConsecutiveCorrect),
		TimesReviewed:       int32(updatedMastery.TimesReviewed),
		EasinessFactor:      updatedMastery.EasinessFactor,
		CurrentIntervalDays: int32(updatedMastery.CurrentIntervalDays),
		NextReviewAt:        pgtype.Timestamptz{Time: updatedMastery.NextReviewAt, Valid: true},
		LastReviewedAt:      pgtype.Timestamptz{Time: updatedMastery.LastReviewedAt, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("error al actualizar maestría en DB: %w", err)
	}

	// 2.2. Insertar en xp_transactions (RN11)
	metadataID := fmt.Sprintf("sense:%d", req.IDSense)
	_, err = txQueries.InsertXPTransaction(ctx, db.InsertXPTransactionParams{
		IDUser:     userUUID,
		AmountXp:   int32(xpAwarded),
		Reason:     db.XpReason(models.XPReasonQuizCorrect),
		MetadataID: pgtype.Text{String: metadataID, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("error al insertar transacción en ledger de XP: %w", err)
	}

	// 2.3. Actualizar saldo, racha y nivel en users
	gamifRow, err := txQueries.AddUserXPAndUpdateStreak(ctx, db.AddUserXPAndUpdateStreakParams{
		AmountXp: int32(xpAwarded),
		IDUser:   userUUID,
	})
	if err != nil {
		return nil, fmt.Errorf("error al actualizar XP y racha del usuario: %w", err)
	}

	// 2.4. Evaluación de logros
	var newlyUnlocked []models.AchievementStatus
	if updatedMastery.ProgressStage == models.StageMastered {
		_ = txQueries.UnlockUserAchievement(ctx, db.UnlockUserAchievementParams{
			IDUser:        userUUID,
			IDAchievement: 1,
		})
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("error al confirmar transacción de revisión: %w", err)
	}

	return &models.CardReviewResult{
		UpdatedMastery: models.UserMastery{
			IDUser:              upsertedRow.IDUser.String(),
			IDSense:             int(upsertedRow.IDSense),
			ProgressStage:       models.MasteryStage(upsertedRow.ProgressStage),
			MasteryPercentage:   upsertedRow.MasteryPercentage,
			ConsecutiveCorrect:  int(upsertedRow.ConsecutiveCorrect),
			TimesReviewed:       int(upsertedRow.TimesReviewed),
			EasinessFactor:      upsertedRow.EasinessFactor,
			CurrentIntervalDays: int(upsertedRow.CurrentIntervalDays),
			NextReviewAt:        upsertedRow.NextReviewAt.Time,
			LastReviewedAt:      upsertedRow.LastReviewedAt.Time,
		},
		XPAwarded:          xpAwarded,
		NewTotalXP:         int(gamifRow.TotalXp),
		NewLevel:           int(gamifRow.CurrentLevel),
		CurrentStreak:      int(gamifRow.CurrentStreak),
		NewlyUnlockedBadge: newlyUnlocked,
	}, nil
}

func (s *masteryService) GetStudySession(ctx context.Context, userID string, userRole models.AccountRole, deckID *int, limit int) (*models.StudySessionPayload, error) {
	if limit <= 0 {
		limit = 15
	}
	if limit > 50 {
		limit = 50
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, NewDomainError(ErrInvalidInput, "Identificador de usuario inválido")
	}

	var sessionCards []models.StudyCardItem
	var totalDue int

	if deckID != nil && *deckID > 0 {
		// Validar RN04 si la baraja es privada
		deckRow, err := s.queries.GetDeckByID(ctx, int32(*deckID))
		if err != nil {
			return nil, MapDBError(err, "Baraja no encontrada")
		}
		if !deckRow.IsPublic && deckRow.IDUser.String() != userID && userRole != models.RoleAdmin {
			return nil, NewDomainError(ErrForbidden, "Violación de regla RN04: No tienes acceso a esta baraja privada")
		}

		dueCount, err := s.queries.CountDueStudySessionCardsByDeck(ctx, db.CountDueStudySessionCardsByDeckParams{
			DeckID: int32(*deckID),
			UserID: userUUID,
		})
		if err != nil {
			return nil, fmt.Errorf("error contando tarjetas pendientes: %w", err)
		}
		totalDue = int(dueCount)

		cardsRows, err := s.queries.GetStudySessionCardsByDeck(ctx, db.GetStudySessionCardsByDeckParams{
			DeckID:   int32(*deckID),
			UserID:   userUUID,
			LimitVal: int32(limit),
		})
		if err != nil {
			return nil, fmt.Errorf("error obteniendo tarjetas de estudio: %w", err)
		}

		sessionCards = s.buildStudyCardItemsFromDeck(ctx, userID, cardsRows)
	} else {
		dueCount, err := s.queries.CountDueStudySessionCardsAll(ctx, userUUID)
		if err != nil {
			return nil, fmt.Errorf("error contando tarjetas pendientes en todo el studio: %w", err)
		}
		totalDue = int(dueCount)

		cardsRows, err := s.queries.GetStudySessionCardsAll(ctx, db.GetStudySessionCardsAllParams{
			UserID:   userUUID,
			LimitVal: int32(limit),
		})
		if err != nil {
			return nil, fmt.Errorf("error obteniendo tarjetas de estudio: %w", err)
		}

		sessionCards = s.buildStudyCardItemsFromAll(ctx, userID, cardsRows)
	}

	return &models.StudySessionPayload{
		DeckID:       deckID,
		TotalDue:     totalDue,
		SessionCards: sessionCards,
	}, nil
}

func (s *masteryService) buildStudyCardItemsFromDeck(ctx context.Context, userID string, rows []db.GetStudySessionCardsByDeckRow) []models.StudyCardItem {
	if len(rows) == 0 {
		return []models.StudyCardItem{}
	}

	senseIDs := make([]int32, len(rows))
	for i, r := range rows {
		senseIDs[i] = r.IDSense
	}
	examplesRows, _ := s.queries.GetExamplesForSenses(ctx, senseIDs)
	examplesMap := make(map[int32][]models.Example)
	for _, ex := range examplesRows {
		examplesMap[ex.IDSense] = append(examplesMap[ex.IDSense], models.Example{
			IDExample: int(ex.IDExample),
			IDSense:   int(ex.IDSense),
			Sentence:  ex.Sentence,
			Source:    TextToPtr(ex.Source),
		})
	}

	items := make([]models.StudyCardItem, len(rows))
	for i, r := range rows {
		var currentMastery *models.UserMastery
		if r.TimesReviewed > 0 {
			stageStr := "NEW"
			if r.ProgressStage != nil {
				stageStr = fmt.Sprint(r.ProgressStage)
			}
			currentMastery = &models.UserMastery{
				IDUser:              userID,
				IDSense:             int(r.IDSense),
				ProgressStage:       models.MasteryStage(stageStr),
				MasteryPercentage:   r.MasteryPercentage,
				ConsecutiveCorrect:  int(r.ConsecutiveCorrect),
				TimesReviewed:       int(r.TimesReviewed),
				EasinessFactor:      r.EasinessFactor,
				CurrentIntervalDays: int(r.CurrentIntervalDays),
				NextReviewAt:        r.NextReviewAt.Time,
				LastReviewedAt:      r.LastReviewedAt.Time,
			}
		}

		items[i] = models.StudyCardItem{
			IDSense:         int(r.IDSense),
			IDWord:          int(r.IDWord),
			Lemma:           r.Lemma,
			PartOfSpeech:    r.PartOfSpeech,
			RegionCode:      r.RegionCode,
			RegionName:      r.RegionName,
			FlagEmoji:       r.FlagEmoji,
			ShortDefinition: r.ShortDefinition,
			Definition:      r.Definition,
			CulturalNote:    TextToPtr(r.CulturalNote),
			Examples:        examplesMap[r.IDSense],
			CurrentMastery:  currentMastery,
		}
	}
	return items
}

func (s *masteryService) buildStudyCardItemsFromAll(ctx context.Context, userID string, rows []db.GetStudySessionCardsAllRow) []models.StudyCardItem {
	if len(rows) == 0 {
		return []models.StudyCardItem{}
	}

	senseIDs := make([]int32, len(rows))
	for i, r := range rows {
		senseIDs[i] = r.IDSense
	}
	examplesRows, _ := s.queries.GetExamplesForSenses(ctx, senseIDs)
	examplesMap := make(map[int32][]models.Example)
	for _, ex := range examplesRows {
		examplesMap[ex.IDSense] = append(examplesMap[ex.IDSense], models.Example{
			IDExample: int(ex.IDExample),
			IDSense:   int(ex.IDSense),
			Sentence:  ex.Sentence,
			Source:    TextToPtr(ex.Source),
		})
	}

	items := make([]models.StudyCardItem, len(rows))
	for i, r := range rows {
		var currentMastery *models.UserMastery
		if r.TimesReviewed > 0 {
			stageStr := "NEW"
			if r.ProgressStage != nil {
				stageStr = fmt.Sprint(r.ProgressStage)
			}
			currentMastery = &models.UserMastery{
				IDUser:              userID,
				IDSense:             int(r.IDSense),
				ProgressStage:       models.MasteryStage(stageStr),
				MasteryPercentage:   r.MasteryPercentage,
				ConsecutiveCorrect:  int(r.ConsecutiveCorrect),
				TimesReviewed:       int(r.TimesReviewed),
				EasinessFactor:      r.EasinessFactor,
				CurrentIntervalDays: int(r.CurrentIntervalDays),
				NextReviewAt:        r.NextReviewAt.Time,
				LastReviewedAt:      r.LastReviewedAt.Time,
			}
		}

		items[i] = models.StudyCardItem{
			IDSense:         int(r.IDSense),
			IDWord:          int(r.IDWord),
			Lemma:           r.Lemma,
			PartOfSpeech:    r.PartOfSpeech,
			RegionCode:      r.RegionCode,
			RegionName:      r.RegionName,
			FlagEmoji:       r.FlagEmoji,
			ShortDefinition: r.ShortDefinition,
			Definition:      r.Definition,
			CulturalNote:    TextToPtr(r.CulturalNote),
			Examples:        examplesMap[r.IDSense],
			CurrentMastery:  currentMastery,
		}
	}
	return items
}

func (s *masteryService) GetMasteryDashboard(ctx context.Context, userID string) (*models.MasteryDashboard, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, NewDomainError(ErrInvalidInput, "Identificador de usuario inválido")
	}

	counts, err := s.queries.GetMasteryStageCounts(ctx, userUUID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo estadísticas de maestría: %w", err)
	}

	gamif, err := s.queries.GetUserGamificationProfile(ctx, userUUID)
	if err != nil {
		return nil, MapDBError(err, "Perfil de usuario no encontrado")
	}

	recentXP, err := s.queries.ListRecentXPTransactions(ctx, db.ListRecentXPTransactionsParams{
		IDUser: userUUID,
		Limit:  10,
	})
	if err != nil {
		return nil, fmt.Errorf("error listando transacciones de XP: %w", err)
	}

	xpItems := make([]models.XPTransactionItem, len(recentXP))
	for i, tx := range recentXP {
		xpItems[i] = models.XPTransactionItem{
			IDTransaction: tx.IDTransaction.String(),
			AmountXP:      int(tx.AmountXp),
			Reason:        tx.Reason,
			MetadataID:    TextToPtr(tx.MetadataID),
			CreatedAt:     tx.CreatedAt.Time,
		}
	}

	achRows, err := s.queries.ListUserAchievements(ctx, userUUID)
	if err != nil {
		return nil, fmt.Errorf("error listando logros: %w", err)
	}

	achievements := make([]models.AchievementStatus, len(achRows))
	for i, a := range achRows {
		var unlockedAt *time.Time
		if a.UnlockedAt.Valid {
			t := a.UnlockedAt.Time
			unlockedAt = &t
		}
		achievements[i] = models.AchievementStatus{
			IDAchievement: fmt.Sprintf("%d", a.IDAchievement),
			Code:          a.Slug,
			Name:          a.Title,
			Description:   a.Description,
			Category:      a.Category,
			Tier:          a.Tier,
			IconName:      a.BadgeIconUrl,
			XPReward:      int(a.XpReward),
			IsUnlocked:    a.IsUnlocked,
			UnlockedAt:    unlockedAt,
		}
	}

	return &models.MasteryDashboard{
		TotalWordsStudied: int(counts.TotalStudied),
		DueTodayCount:     int(counts.DueToday),
		CountNew:          int(counts.CountNew),
		CountLearning:     int(counts.CountLearning),
		CountReviewing:    int(counts.CountReviewing),
		CountMastered:     int(counts.CountMastered),
		GlobalRetention:   counts.GlobalRetention,
		TotalXP:           int(gamif.TotalXp),
		CurrentLevel:      int(gamif.CurrentLevel),
		CurrentStreak:     int(gamif.CurrentStreak),
		LongestStreak:     int(gamif.LongestStreak),
		RecentXP:          xpItems,
		Achievements:      achievements,
	}, nil
}

func (s *masteryService) ListMasteryItems(ctx context.Context, userID string, stage *string) ([]models.StudyCardItem, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, NewDomainError(ErrInvalidInput, "Identificador de usuario inválido")
	}

	var stageText pgtype.Text
	if stage != nil && strings.TrimSpace(*stage) != "" {
		st := strings.ToUpper(strings.TrimSpace(*stage))
		stageText = pgtype.Text{String: st, Valid: true}
	}

	rows, err := s.queries.ListMasteryCardsByStage(ctx, db.ListMasteryCardsByStageParams{
		UserID: userUUID,
		Stage:  stageText,
	})
	if err != nil {
		return nil, fmt.Errorf("error listando tarjetas de maestría: %w", err)
	}

	senseIDs := make([]int32, len(rows))
	for i, r := range rows {
		senseIDs[i] = r.IDSense
	}
	examplesRows, _ := s.queries.GetExamplesForSenses(ctx, senseIDs)
	examplesMap := make(map[int32][]models.Example)
	for _, ex := range examplesRows {
		examplesMap[ex.IDSense] = append(examplesMap[ex.IDSense], models.Example{
			IDExample: int(ex.IDExample),
			IDSense:   int(ex.IDSense),
			Sentence:  ex.Sentence,
			Source:    TextToPtr(ex.Source),
		})
	}

	items := make([]models.StudyCardItem, len(rows))
	for i, r := range rows {
		items[i] = models.StudyCardItem{
			IDSense:         int(r.IDSense),
			IDWord:          int(r.IDWord),
			Lemma:           r.Lemma,
			PartOfSpeech:    r.PartOfSpeech,
			RegionCode:      r.RegionCode,
			RegionName:      r.RegionName,
			FlagEmoji:       r.FlagEmoji,
			ShortDefinition: r.ShortDefinition,
			Definition:      r.Definition,
			CulturalNote:    TextToPtr(r.CulturalNote),
			Examples:        examplesMap[r.IDSense],
			CurrentMastery: &models.UserMastery{
				IDUser:              userID,
				IDSense:             int(r.IDSense),
				ProgressStage:       models.MasteryStage(r.ProgressStage),
				MasteryPercentage:   r.MasteryPercentage,
				ConsecutiveCorrect:  int(r.ConsecutiveCorrect),
				TimesReviewed:       int(r.TimesReviewed),
				EasinessFactor:      r.EasinessFactor,
				CurrentIntervalDays: int(r.CurrentIntervalDays),
				NextReviewAt:        r.NextReviewAt.Time,
				LastReviewedAt:      r.LastReviewedAt.Time,
			},
		}
	}

	return items, nil
}

func (s *masteryService) ListXPLedger(ctx context.Context, userID string, limit int) ([]models.XPTransactionItem, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, NewDomainError(ErrInvalidInput, "Identificador de usuario inválido")
	}

	txRows, err := s.queries.ListRecentXPTransactions(ctx, db.ListRecentXPTransactionsParams{
		IDUser: userUUID,
		Limit:  int32(limit),
	})
	if err != nil {
		return nil, fmt.Errorf("error listando transacciones de XP: %w", err)
	}

	items := make([]models.XPTransactionItem, len(txRows))
	for i, r := range txRows {
		items[i] = models.XPTransactionItem{
			IDTransaction: r.IDTransaction.String(),
			AmountXP:      int(r.AmountXp),
			Reason:        r.Reason,
			MetadataID:    TextToPtr(r.MetadataID),
			CreatedAt:     r.CreatedAt.Time,
		}
	}
	return items, nil
}
