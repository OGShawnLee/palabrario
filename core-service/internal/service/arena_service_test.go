package service_test

import (
	"context"
	"errors"
	"testing"
	"strings"
	"time"

	"github.com/google/uuid"

	"palabrario/core-service/internal/db"
	"palabrario/core-service/internal/service"
	"palabrario/core-service/internal/testutil"
	"palabrario/pkg/auth"
	"palabrario/pkg/models"
	"palabrario/pkg/utils"
)

// ==============================================================================
// 1. PRUEBA DE SEPARACIÓN CRIPTOGRÁFICA S2S (SUB-FASE 4.1)
// ==============================================================================

func TestAuth_S2S_CryptographicSeparation(t *testing.T) {
	secret := "test-secret-key-palabrario-s2s-12345"
	tokenService := auth.NewTokenService(secret, time.Hour)

	user := &models.UserProfile{
		IDUser:        "b1111111-0000-0000-0000-000000000003",
		Email:         "gael@palabrario.org",
		UserName:      "gael_learner",
		AnchorDialect: "MX",
		Role:          models.RoleLearner,
	}

	// 1. Emitir token de usuario y verificarlo
	userToken, err := tokenService.GenerateToken(user)
	if err != nil {
		t.Fatalf("GenerateToken falló: %v", err)
	}

	userClaims, err := tokenService.VerifyToken(userToken)
	if err != nil {
		t.Fatalf("VerifyToken falló con token de usuario válido: %v", err)
	}
	if userClaims.UserID != user.IDUser || userClaims.Issuer != "palabrario-core-service" {
		t.Fatalf("Claims de usuario no coinciden: %+v", userClaims)
	}

	// Token de usuario DEBE fallar en VerifyInternalToken (distinta clave e issuer)
	_, err = tokenService.VerifyInternalToken(userToken)
	if err == nil {
		t.Fatalf("Se esperaba error al verificar token de usuario con VerifyInternalToken, pero pasó con éxito")
	}

	// 2. Emitir token interno S2S y verificarlo
	s2sToken, err := tokenService.GenerateInternalToken("realtime-service", "MX8921")
	if err != nil {
		t.Fatalf("GenerateInternalToken falló: %v", err)
	}

	s2sClaims, err := tokenService.VerifyInternalToken(s2sToken)
	if err != nil {
		t.Fatalf("VerifyInternalToken falló con token interno válido: %v", err)
	}
	if s2sClaims.ServiceName != "realtime-service" || s2sClaims.MatchScope != "MX8921" || s2sClaims.Issuer != "palabrario-internal-s2s" {
		t.Fatalf("Claims S2S no coinciden: %+v", s2sClaims)
	}

	// Token interno DEBE fallar en VerifyToken (distinta clave e issuer)
	_, err = tokenService.VerifyToken(s2sToken)
	if err == nil {
		t.Fatalf("Se esperaba error al verificar token S2S con VerifyToken, pero pasó con éxito")
	}
}

// ==============================================================================
// 2. PRUEBA DE LÍMITES ESTRICTOS DE ARENA: CU14 1a
// ==============================================================================

func TestArenaCatalogService_LimitsAndValidation(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	if testDB == nil {
		return
	}

	arenaSvc := service.NewArenaCatalogService(testDB)
	studioSvc := service.NewStudioService(testDB)
	ctx := context.Background()

	gaelID := "b1111111-0000-0000-0000-000000000003"

	// 1. Crear una baraja con solo 3 sentidos (menos del piso mínimo de 8)
	smallDeck, err := studioSvc.CreateDeck(ctx, gaelID, models.RoleLearner, models.CreateDeckRequest{
		Title:       "Baraja Pequeña",
		Description: utils.ToPtr("Solo 3 sentidos"),
		IsPublic:    true,
		SenseIDs:    []int{101, 102, 105},
	})
	if err != nil {
		t.Fatalf("CreateDeck falló: %v", err)
	}

	// Intentar preparar partida con baraja de < 8 sentidos -> Debe fallar con ErrUnprocessable (422)
	_, err = arenaSvc.PrepareMatch(ctx, gaelID, models.RoleLearner, models.CreateMatchRequest{
		IDDeck:      smallDeck.IDDeck,
		TotalRounds: 8,
	})
	if err == nil || !errors.Is(err, service.ErrUnprocessable) {
		t.Fatalf("se esperaba ErrUnprocessable al abrir partida con baraja de < 8 sentidos, obtenido: %v", err)
	}

	// 2. Usar baraja semilla 1 pero pedir 12 rondas (TotalRounds > ApprovedSenses)
	_, err = arenaSvc.PrepareMatch(ctx, gaelID, models.RoleLearner, models.CreateMatchRequest{
		IDDeck:      1,
		TotalRounds: 12,
	})
	if err == nil || !errors.Is(err, service.ErrUnprocessable) {
		t.Fatalf("se esperaba ErrUnprocessable al pedir más rondas (12) que sentidos disponibles (8), obtenido: %v", err)
	}

	// 3. Validar rangos inválidos de TotalRounds y RoundTimeSec
	_, err = arenaSvc.PrepareMatch(ctx, gaelID, models.RoleLearner, models.CreateMatchRequest{
		IDDeck:      1,
		TotalRounds: 5, // menor a 8
	})
	if err == nil || !errors.Is(err, service.ErrInvalidInput) {
		t.Fatalf("se esperaba ErrInvalidInput con TotalRounds=5, obtenido: %v", err)
	}

	_, err = arenaSvc.PrepareMatch(ctx, gaelID, models.RoleLearner, models.CreateMatchRequest{
		IDDeck:       1,
		TotalRounds:  8,
		RoundTimeSec: 40, // mayor a 30
	})
	if err == nil || !errors.Is(err, service.ErrInvalidInput) {
		t.Fatalf("se esperaba ErrInvalidInput con RoundTimeSec=40, obtenido: %v", err)
	}
}

// ==============================================================================
// 3. PRUEBA DE PREPARACIÓN DE PARTIDA 8 RONDAS Y RECUPERACIÓN 3NF
// ==============================================================================

func TestArenaCatalogService_PrepareAndGetMatchPackage(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	if testDB == nil {
		return
	}

	arenaSvc := service.NewArenaCatalogService(testDB)
	ctx := context.Background()

	gaelID := "b1111111-0000-0000-0000-000000000003"

	// 1. Preparar partida de 8 rondas con la baraja semilla 1
	pkg, err := arenaSvc.PrepareMatch(ctx, gaelID, models.RoleLearner, models.CreateMatchRequest{
		IDDeck:       1,
		TotalRounds:  8,
		RoundTimeSec: 15,
	})
	if err != nil {
		t.Fatalf("PrepareMatch falló: %v", err)
	}

	if pkg.IDMatch == "" || len(pkg.RoomCode) != 6 {
		t.Fatalf("IDMatch o RoomCode inválidos: %+v", pkg)
	}
	if pkg.TotalRounds != 8 || pkg.RoundTimeSec != 15 {
		t.Fatalf("Configuración de rondas no coincide: %+v", pkg)
	}
	if len(pkg.Questions) != 8 {
		t.Fatalf("Se esperaban 8 preguntas, obtenidas: %d", len(pkg.Questions))
	}

	// Comprobar que los 8 sentidos son únicos (sin repetición)
	seenSenses := make(map[int]bool)
	for _, q := range pkg.Questions {
		if seenSenses[q.IDSense] {
			t.Fatalf("Sentido duplicado detectado en las rondas de la partida: %d", q.IDSense)
		}
		seenSenses[q.IDSense] = true
		if len(q.Options) != 4 {
			t.Fatalf("Cada reactivo debe contener exactamente 4 opciones, obtenido: %d", len(q.Options))
		}
		if q.CorrectOption == "" {
			t.Fatalf("CorrectOption no puede estar vacía en InternalArenaQuestion")
		}
	}

	// 2. Recuperar el paquete con GetMatchPackageByRoomCode
	pkgFetched, err := arenaSvc.GetMatchPackageByRoomCode(ctx, pkg.RoomCode)
	if err != nil {
		t.Fatalf("GetMatchPackageByRoomCode falló: %v", err)
	}

	if pkgFetched.IDMatch != pkg.IDMatch || pkgFetched.RoomCode != pkg.RoomCode {
		t.Fatalf("El paquete recuperado no coincide con el preparado: %+v vs %+v", pkgFetched, pkg)
	}
	if len(pkgFetched.Questions) != 8 {
		t.Fatalf("Se esperaban 8 preguntas en el paquete recuperado, obtenidas: %d", len(pkgFetched.Questions))
	}
}

// ==============================================================================
// 4. PRUEBA DE CONSOLIDACIÓN TRANSACCIONAL, DESEMPATE, XP, SM-2 Y REPLAY 3NF
// ==============================================================================

func TestArenaCatalogService_ConsolidateAndReplay(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	if testDB == nil {
		return
	}

	arenaSvc := service.NewArenaCatalogService(testDB)
	queries := db.New(testDB.Pool)
	ctx := context.Background()

	gaelID := "b1111111-0000-0000-0000-000000000003"
	silasID := "b1111111-0000-0000-0000-000000000002"

	// 1. Preparar una partida con baraja 1
	pkg, err := arenaSvc.PrepareMatch(ctx, gaelID, models.RoleLearner, models.CreateMatchRequest{
		IDDeck:       1,
		TotalRounds:  8,
		RoundTimeSec: 15,
	})
	if err != nil {
		t.Fatalf("PrepareMatch falló: %v", err)
	}

	// 2. Simular resultados empatados en puntos (FinalScore = 800) y aciertos (4 sentidos acertados, 4 fallados),
	// pero Gael con MENOR TotalTimeSpentMs (4500 ms) que Silas (7200 ms).
	// Por criterio determinístico: Gael DEBE ser 1° Lugar (Rank 1, +150 XP) y Silas 2° Lugar (Rank 2, +100 XP).
	gaelAttempts := make([]models.ParticipantRoundAttempt, 8)
	silasAttempts := make([]models.ParticipantRoundAttempt, 8)

	for i := 0; i < 8; i++ {
		q := pkg.Questions[i]
		isGaelCorrect := (i < 4)
		isSilasCorrect := (i < 4)

		gaelOpt := q.CorrectOption
		if !isGaelCorrect {
			gaelOpt = "opcion_erronea_gael"
		}
		silasOpt := q.CorrectOption
		if !isSilasCorrect {
			silasOpt = "opcion_erronea_silas"
		}

		gaelAttempts[i] = models.ParticipantRoundAttempt{
			RoundNumber:    q.RoundNumber,
			IDSense:        q.IDSense,
			IDQuestion:     q.IDQuestion,
			SelectedOption: gaelOpt,
			IsCorrect:      isGaelCorrect,
			TimeSpentMs:    500 + i*10,
			PointsEarned:   100,
		}

		silasAttempts[i] = models.ParticipantRoundAttempt{
			RoundNumber:    q.RoundNumber,
			IDSense:        q.IDSense,
			IDQuestion:     q.IDQuestion,
			SelectedOption: silasOpt,
			IsCorrect:      isSilasCorrect,
			TimeSpentMs:    900 + i*10,
			PointsEarned:   100,
		}
	}

	correctSenses := []int{pkg.Questions[0].IDSense, pkg.Questions[1].IDSense, pkg.Questions[2].IDSense, pkg.Questions[3].IDSense}
	wrongSenses := []int{pkg.Questions[4].IDSense, pkg.Questions[5].IDSense, pkg.Questions[6].IDSense, pkg.Questions[7].IDSense}

	req := models.ConsolidateMatchRequest{
		IDMatch:  pkg.IDMatch,
		RoomCode: pkg.RoomCode,
		Participants: []models.ParticipantFinalResult{
			{
				IDUser:           silasID, // Ponemos a Silas primero deliberadamente para verificar que el sort no dependa del orden de entrada
				FinalScore:       800,
				AccuracyRate:     0.5,
				TotalTimeSpentMs: 7200,
				CorrectSenseIDs:  correctSenses,
				WrongSenseIDs:    wrongSenses,
				RoundAttempts:    silasAttempts,
			},
			{
				IDUser:           gaelID,
				FinalScore:       800,
				AccuracyRate:     0.5,
				TotalTimeSpentMs: 4500, // Menor tiempo -> Gael debe ganar el desempate
				CorrectSenseIDs:  correctSenses,
				WrongSenseIDs:    wrongSenses,
				RoundAttempts:    gaelAttempts,
			},
		},
	}

	summary, err := arenaSvc.ConsolidateMatch(ctx, req)
	if err != nil {
		t.Fatalf("ConsolidateMatch falló: %v", err)
	}

	if summary.Status != models.MatchStatusFinished {
		t.Fatalf("Se esperaba estado FINISHED en el resumen, obtenido: %s", summary.Status)
	}

	if summary.WinnerUserID == nil || *summary.WinnerUserID != gaelID {
		t.Fatalf("Se esperaba que Gael fuera el ganador por menor tiempo (4500 < 7200), obtenido: %v", summary.WinnerUserID)
	}

	if len(summary.Podium) != 2 {
		t.Fatalf("Se esperaban 2 participantes en el podio, obtenidos: %d", len(summary.Podium))
	}
	if summary.Podium[0].IDUser != gaelID || summary.Podium[0].Rank != 1 {
		t.Fatalf("El primer puesto del podio debe ser Gael con Rank 1: %+v", summary.Podium[0])
	}
	if summary.Podium[1].IDUser != silasID || summary.Podium[1].Rank != 2 {
		t.Fatalf("El segundo puesto del podio debe ser Silas con Rank 2: %+v", summary.Podium[1])
	}

	// 3. Verificar que intentar consolidar de nuevo la misma partida falle con ErrConflict
	_, err = arenaSvc.ConsolidateMatch(ctx, req)
	if err == nil || !errors.Is(err, service.ErrConflict) {
		t.Fatalf("se esperaba ErrConflict al reconsolidar una partida ya finalizada, obtenido: %v", err)
	}

	// 4. Verificar que GetMatchPackageByRoomCode falle con ErrConflict si ya finalizó
	_, err = arenaSvc.GetMatchPackageByRoomCode(ctx, pkg.RoomCode)
	if err == nil || !errors.Is(err, service.ErrConflict) {
		t.Fatalf("se esperaba ErrConflict al pedir el paquete de una partida finalizada, obtenido: %v", err)
	}

	// 5. Verificar que se grabaron los asientos de XP (+150 y +100) en xp_transactions
	matchUUID, _ := uuid.Parse(pkg.IDMatch)
	gaelUUID, _ := uuid.Parse(gaelID)
	silasUUID, _ := uuid.Parse(silasID)

	gaelTx, err := queries.ListRecentXPTransactions(ctx, db.ListRecentXPTransactionsParams{
		IDUser: gaelUUID,
		Limit:  10,
	})
	if err != nil {
		t.Fatalf("Error consultando transacciones de XP de Gael: %v", err)
	}
	foundGaelXP := false
	for _, tx := range gaelTx {
		if tx.AmountXp == 150 && tx.MetadataID.Valid && tx.MetadataID.String == matchUUID.String() {
			foundGaelXP = true
			break
		}
	}
	if !foundGaelXP {
		t.Fatalf("No se encontró la transacción de +150 XP para el ganador Gael")
	}

	silasTx, err := queries.ListRecentXPTransactions(ctx, db.ListRecentXPTransactionsParams{
		IDUser: silasUUID,
		Limit:  10,
	})
	if err != nil {
		t.Fatalf("Error consultando transacciones de XP de Silas: %v", err)
	}
	foundSilasXP := false
	for _, tx := range silasTx {
		if tx.AmountXp == 100 && tx.MetadataID.Valid && tx.MetadataID.String == matchUUID.String() {
			foundSilasXP = true
			break
		}
	}
	if !foundSilasXP {
		t.Fatalf("No se encontró la transacción de +100 XP para el 2° lugar Silas")
	}

	// 6. Verificar actualización de user_mastery con SM-2 para los sentidos correctos y erróneos
	gaelMastery, err := queries.GetUserMasteryForSense(ctx, db.GetUserMasteryForSenseParams{
		IDUser:  gaelUUID,
		IDSense: int32(correctSenses[0]),
	})
	if err != nil {
		t.Fatalf("Error consultando user_mastery para sentido acertado: %v", err)
	}
	if gaelMastery.ConsecutiveCorrect != 1 || gaelMastery.TimesReviewed != 1 {
		t.Fatalf("user_mastery no se actualizó correctamente con SM-2: %+v", gaelMastery)
	}

	// 7. Verificar GetMatchReplay completo con las 8 rondas y las respuestas exactas de ambos jugadores
	replay, err := arenaSvc.GetMatchReplay(ctx, gaelID, models.RoleLearner, pkg.IDMatch)
	if err != nil {
		t.Fatalf("GetMatchReplay falló: %v", err)
	}

	if replay.IDMatch != pkg.IDMatch || len(replay.Rounds) != 8 {
		t.Fatalf("El replay no contiene las 8 rondas esperadas: %+v", replay)
	}

	for i, r := range replay.Rounds {
		if len(r.Answers) != 2 {
			t.Fatalf("Ronda %d debe tener exactamente 2 respuestas de participantes, obtenidas: %d", i+1, len(r.Answers))
		}
		// Verificar que están las opciones elegidas por Gael y Silas
		var foundGael, foundSilas bool
		for _, ans := range r.Answers {
			if ans.IDUser == gaelID {
				foundGael = true
				if i < 4 && !ans.IsCorrect {
					t.Fatalf("Ronda %d: Gael debía tener respuesta correcta", i+1)
				}
				if i >= 4 && ans.IsCorrect {
					t.Fatalf("Ronda %d: Gael debía tener respuesta errónea", i+1)
				}
			}
			if ans.IDUser == silasID {
				foundSilas = true
			}
		}
		if !foundGael || !foundSilas {
			t.Fatalf("Ronda %d: no se encontraron las respuestas de ambos jugadores", i+1)
		}
	}

	// 8. Verificar ListMatchHistory
	history, err := arenaSvc.ListMatchHistory(ctx, gaelID, 10)
	if err != nil {
		t.Fatalf("ListMatchHistory falló: %v", err)
	}
	if len(history) == 0 {
		t.Fatalf("ListMatchHistory no devolvió ninguna partida para Gael")
	}
	if history[0].IDMatch != pkg.IDMatch {
		t.Fatalf("La partida más reciente del historial debería ser %s, obtenida: %s", pkg.IDMatch, history[0].IDMatch)
	}
}

// ==============================================================================
// 5. PRUEBA DE PRIVACIDAD EN REPLAY: CU16 (ESPECTADOR EN PARTIDA PRIVADA)
// ==============================================================================

func TestArenaCatalogService_GetMatchReplay_PrivateMatchAuthorization(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	if testDB == nil {
		return
	}

	arenaSvc := service.NewArenaCatalogService(testDB)
	ctx := context.Background()

	gaelID := "b1111111-0000-0000-0000-000000000003"  // Participante
	silasID := "b1111111-0000-0000-0000-000000000002" // No participante
	adminID := "b1111111-0000-0000-0000-000000000001" // Admin (no participante)

	isPublicFalse := false
	pkg, err := arenaSvc.PrepareMatch(ctx, gaelID, models.RoleLearner, models.CreateMatchRequest{
		IDDeck:       1,
		TotalRounds:  8,
		RoundTimeSec: 15,
		IsPublic:     &isPublicFalse,
	})
	if err != nil {
		t.Fatalf("PrepareMatch falló: %v", err)
	}

	// Consolidar partida con solo Gael como participante
	_, err = arenaSvc.ConsolidateMatch(ctx, models.ConsolidateMatchRequest{
		IDMatch:  pkg.IDMatch,
		RoomCode: pkg.RoomCode,
		Participants: []models.ParticipantFinalResult{
			{
				IDUser:           gaelID,
				FinalScore:       500,
				AccuracyRate:     1.0,
				TotalTimeSpentMs: 2000,
				CorrectSenseIDs:  []int{pkg.Questions[0].IDSense},
				WrongSenseIDs:    []int{},
				RoundAttempts:    []models.ParticipantRoundAttempt{},
			},
		},
	})
	if err != nil {
		t.Fatalf("ConsolidateMatch falló: %v", err)
	}

	// 1. Silas (Learner no participante) intenta ver el replay privado -> ErrForbidden
	_, err = arenaSvc.GetMatchReplay(ctx, silasID, models.RoleLearner, pkg.IDMatch)
	if err == nil || !errors.Is(err, service.ErrForbidden) {
		t.Fatalf("Se esperaba ErrForbidden para usuario no participante en replay privado, obtenido: %v", err)
	}

	// 2. Gael (Learner participante) ve el replay privado -> Éxito
	replayGael, err := arenaSvc.GetMatchReplay(ctx, gaelID, models.RoleLearner, pkg.IDMatch)
	if err != nil {
		t.Fatalf("Gael (participante) debería poder ver el replay privado: %v", err)
	}
	if replayGael.IDMatch != pkg.IDMatch {
		t.Fatalf("IDMatch no coincide: %s vs %s", replayGael.IDMatch, pkg.IDMatch)
	}

	// 3. Admin (no participante pero RoleAdmin) ve el replay privado -> Éxito
	replayAdmin, err := arenaSvc.GetMatchReplay(ctx, adminID, models.RoleAdmin, pkg.IDMatch)
	if err != nil {
		t.Fatalf("Admin debería poder ver el replay privado: %v", err)
	}
	if replayAdmin.IDMatch != pkg.IDMatch {
		t.Fatalf("IDMatch no coincide: %s vs %s", replayAdmin.IDMatch, pkg.IDMatch)
	}
}

func TestArenaCatalogService_PrepareMatch_GracefulDegradation_RN13(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	if testDB == nil {
		return
	}

	arenaSvc := service.NewArenaCatalogService(testDB)
	ctx := context.Background()
	gaelID := "b1111111-0000-0000-0000-000000000003"
	silasID := "b1111111-0000-0000-0000-000000000002"

	// 1. Configurar solo 2 sentidos con imagen (< 4 imágenes en la baraja)
	_, err := testDB.Pool.Exec(ctx, `
		UPDATE senses SET image_url = 'http://localhost:8002/uploads/sense/img101.png' WHERE id_sense = 101;
		UPDATE senses SET image_url = 'http://localhost:8002/uploads/sense/img102.png' WHERE id_sense = 102;
	`)
	if err != nil {
		t.Fatalf("falló actualización de imágenes en sentidos: %v", err)
	}

	// 2. Insertar pregunta curada WORD_TO_IMAGE para el sentido 101
	_, err = testDB.Pool.Exec(ctx, `
		INSERT INTO questions (id_sense, id_editor, type, prompt, custom_distractors, explanation)
		VALUES (101, $1, 'WORD_TO_IMAGE', '¿Cuál de estas imágenes representa popote?', '[]', 'Popote explicación')
	`, silasID)
	if err != nil {
		t.Fatalf("falló inserción de pregunta curada WORD_TO_IMAGE: %v", err)
	}

	// 3. Preparar partida de 8 rondas con la baraja 1
	pkg, err := arenaSvc.PrepareMatch(ctx, gaelID, models.RoleLearner, models.CreateMatchRequest{
		IDDeck:       1,
		TotalRounds:  8,
		RoundTimeSec: 15,
	})
	if err != nil {
		t.Fatalf("PrepareMatch falló: %v", err)
	}

	// 4. RN13: La pregunta del sentido 101 DEBE haberse degradado a WORD_TO_DEFINITION
	foundSense101 := false
	for _, q := range pkg.Questions {
		if q.IDSense == 101 {
			foundSense101 = true
			if q.Modality != "WORD_TO_DEFINITION" {
				t.Errorf("RN13 violada: se esperaba degradación a WORD_TO_DEFINITION con < 4 imágenes, obtenido: %s", q.Modality)
			}
			if len(q.Options) != 4 {
				t.Errorf("se esperaban 4 opciones en la ronda degradada, obtenido: %d", len(q.Options))
			}
			if strings.HasPrefix(q.CorrectOption, "http") {
				t.Errorf("la opción correcta no debería ser una URL de imagen tras degradarse, obtenido: %s", q.CorrectOption)
			}
		}
		// Ninguna ronda debe ser WORD_TO_IMAGE en una baraja con < 4 imágenes
		if q.Modality == "WORD_TO_IMAGE" {
			t.Errorf("RN13 violada: ninguna ronda debe ser WORD_TO_IMAGE cuando la baraja tiene < 4 imágenes")
		}
	}
	if !foundSense101 {
		t.Fatalf("no se encontró el sentido 101 en las preguntas de la partida preparada")
	}
}

func TestArenaCatalogService_PrepareMatch_WordToImage_RN13(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	if testDB == nil {
		return
	}

	arenaSvc := service.NewArenaCatalogService(testDB)
	ctx := context.Background()
	gaelID := "b1111111-0000-0000-0000-000000000003"
	silasID := "b1111111-0000-0000-0000-000000000002"

	// 1. Configurar 4 sentidos con imagen única (>= 4 imágenes en la baraja)
	_, err := testDB.Pool.Exec(ctx, `
		UPDATE senses SET image_url = 'http://localhost:8002/uploads/sense/img101.png' WHERE id_sense = 101;
		UPDATE senses SET image_url = 'http://localhost:8002/uploads/sense/img102.png' WHERE id_sense = 102;
		UPDATE senses SET image_url = 'http://localhost:8002/uploads/sense/img105.png' WHERE id_sense = 105;
		UPDATE senses SET image_url = 'http://localhost:8002/uploads/sense/img108.png' WHERE id_sense = 108;
	`)
	if err != nil {
		t.Fatalf("falló actualización de imágenes en sentidos: %v", err)
	}

	// 2. Insertar pregunta curada WORD_TO_IMAGE para el sentido 101
	_, err = testDB.Pool.Exec(ctx, `
		INSERT INTO questions (id_sense, id_editor, type, prompt, custom_distractors, explanation)
		VALUES (101, $1, 'WORD_TO_IMAGE', '¿Cuál de estas imágenes representa popote?', '[]', 'Popote explicación')
	`, silasID)
	if err != nil {
		t.Fatalf("falló inserción de pregunta curada WORD_TO_IMAGE: %v", err)
	}

	// 3. Preparar partida de 8 rondas con la baraja 1
	pkg, err := arenaSvc.PrepareMatch(ctx, gaelID, models.RoleLearner, models.CreateMatchRequest{
		IDDeck:       1,
		TotalRounds:  8,
		RoundTimeSec: 15,
	})
	if err != nil {
		t.Fatalf("PrepareMatch falló: %v", err)
	}

	visualRoundsCount := 0
	for _, q := range pkg.Questions {
		if q.Modality == "WORD_TO_IMAGE" {
			visualRoundsCount++
			if !strings.HasPrefix(q.CorrectOption, "http://localhost:8002/uploads/") {
				t.Errorf("CorrectOption de WORD_TO_IMAGE debe ser una URL de imagen válida, obtenido: %s", q.CorrectOption)
			}
			if len(q.Options) != 4 {
				t.Fatalf("se esperaban 4 opciones en ronda WORD_TO_IMAGE, obtenido: %d", len(q.Options))
			}
			seenOpts := make(map[string]bool)
			foundCorrect := false
			for _, opt := range q.Options {
				if seenOpts[opt] {
					t.Errorf("opción duplicada en ronda visual: %s", opt)
				}
				seenOpts[opt] = true
				if !strings.HasPrefix(opt, "http://localhost:8002/uploads/") {
					t.Errorf("la opción visual debe ser una URL de imagen, obtenido: %s", opt)
				}
				if opt == q.CorrectOption {
					foundCorrect = true
				}
			}
			if !foundCorrect {
				t.Errorf("la opción correcta %s no se encontró en las opciones: %+v", q.CorrectOption, q.Options)
			}
		}
	}

	if visualRoundsCount == 0 {
		t.Fatalf("se esperaban rondas con modalidad WORD_TO_IMAGE al tener >= 4 imágenes en la baraja")
	}
}
