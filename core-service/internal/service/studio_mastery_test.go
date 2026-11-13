package service_test

import (
	"context"
	"errors"
	"math"
	"reflect"
	"testing"
	"time"

	"palabrario/core-service/internal/service"
	"palabrario/core-service/internal/testutil"
	"palabrario/pkg/models"
	"palabrario/pkg/utils"
)

// ==============================================================================
// PRUEBAS DE STUDIO PERSONAL: RN03, RN04, RN05
// ==============================================================================

func TestStudioService_RN03_PropertyAndPermissions(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	if testDB == nil {
		return
	}

	studioSvc := service.NewStudioService(testDB)
	ctx := context.Background()

	gaelID := "b1111111-0000-0000-0000-000000000003"  // LEARNER
	silasID := "b1111111-0000-0000-0000-000000000002" // EDITOR (no dueño)
	adminID := "b1111111-0000-0000-0000-000000000001" // ADMIN

	// 1. Gael crea una baraja
	deck, err := studioSvc.CreateDeck(ctx, gaelID, models.RoleLearner, models.CreateDeckRequest{
		Title:       "Modismos de Gael",
		Description: utils.ToPtr("Baraja personal para pruebas"),
		IsPublic:    true,
		SenseIDs:    []int{101, 102},
	})
	if err != nil {
		t.Fatalf("CreateDeck falló para Gael: %v", err)
	}

	// 2. Silas (no propietario) intenta actualizar la baraja -> Debe fallar con ErrForbidden (RN03)
	_, err = studioSvc.UpdateDeck(ctx, silasID, models.RoleEditor, deck.IDDeck, models.UpdateDeckRequest{
		Title:    "Intento de hack por Silas",
		IsPublic: true,
	})
	if err == nil || !errors.Is(err, service.ErrForbidden) {
		t.Fatalf("se esperaba ErrForbidden (RN03) al editar baraja ajena, obtenido: %v", err)
	}

	// 3. Silas intenta eliminar un elemento de la baraja -> Debe fallar con ErrForbidden (RN03)
	err = studioSvc.RemoveDeckItem(ctx, silasID, models.RoleEditor, deck.IDDeck, 101)
	if err == nil || !errors.Is(err, service.ErrForbidden) {
		t.Fatalf("se esperaba ErrForbidden (RN03) al modificar elementos de baraja ajena, obtenido: %v", err)
	}

	// 4. Admin (ADMIN) puede actualizar la baraja -> Éxito (RN03/RN08)
	updatedByAdmin, err := studioSvc.UpdateDeck(ctx, adminID, models.RoleAdmin, deck.IDDeck, models.UpdateDeckRequest{
		Title:    "Modismos de Gael (Auditado)",
		IsPublic: true,
	})
	if err != nil {
		t.Fatalf("Admin debería poder actualizar la baraja por RN08: %v", err)
	}
	if updatedByAdmin.Title != "Modismos de Gael (Auditado)" {
		t.Errorf("se esperaba título actualizado, obtenido: %s", updatedByAdmin.Title)
	}

	// 5. Gael (propietario) elimina la baraja -> Éxito
	err = studioSvc.DeleteDeck(ctx, gaelID, models.RoleLearner, deck.IDDeck)
	if err != nil {
		t.Fatalf("Gael debería poder eliminar su propia baraja: %v", err)
	}
}

func TestStudioService_RN04_ResourcePrivacy(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	if testDB == nil {
		return
	}

	studioSvc := service.NewStudioService(testDB)
	ctx := context.Background()

	gaelID := "b1111111-0000-0000-0000-000000000003"  // LEARNER
	silasID := "b1111111-0000-0000-0000-000000000002" // EDITOR
	adminID := "b1111111-0000-0000-0000-000000000001" // ADMIN

	// 1. Gael crea una baraja PRIVADA
	privateDeck, err := studioSvc.CreateDeck(ctx, gaelID, models.RoleLearner, models.CreateDeckRequest{
		Title:    "Baraja Secreta de Gael",
		IsPublic: false,
		SenseIDs: []int{101},
	})
	if err != nil {
		t.Fatalf("CreateDeck privado falló: %v", err)
	}

	// 2. Silas intenta ver el detalle de la baraja privada -> Debe fallar con ErrForbidden (RN04)
	_, err = studioSvc.GetDeckDetail(ctx, &silasID, utils.ToPtr(models.RoleEditor), privateDeck.IDDeck)
	if err == nil || !errors.Is(err, service.ErrForbidden) {
		t.Fatalf("se esperaba ErrForbidden (RN04) al consultar baraja privada ajena, obtenido: %v", err)
	}

	// 3. Usuario anónimo intenta ver la baraja privada -> Debe fallar con ErrForbidden (RN04)
	_, err = studioSvc.GetDeckDetail(ctx, nil, nil, privateDeck.IDDeck)
	if err == nil || !errors.Is(err, service.ErrForbidden) {
		t.Fatalf("se esperaba ErrForbidden (RN04) para usuario no autenticado, obtenido: %v", err)
	}

	// 4. Silas intenta guardar la baraja privada -> Debe fallar con ErrForbidden (RN04)
	err = studioSvc.SaveDeck(ctx, silasID, models.RoleEditor, privateDeck.IDDeck)
	if err == nil || !errors.Is(err, service.ErrForbidden) {
		t.Fatalf("se esperaba ErrForbidden (RN04) al guardar baraja privada ajena, obtenido: %v", err)
	}

	// 5. Gael (dueño) y Admin sí pueden verla
	detailGael, err := studioSvc.GetDeckDetail(ctx, &gaelID, utils.ToPtr(models.RoleLearner), privateDeck.IDDeck)
	if err != nil {
		t.Fatalf("dueño debería ver su baraja privada: %v", err)
	}
	if detailGael.Title != "Baraja Secreta de Gael" {
		t.Errorf("título incorrecto: %s", detailGael.Title)
	}

	detailAdmin, err := studioSvc.GetDeckDetail(ctx, &adminID, utils.ToPtr(models.RoleAdmin), privateDeck.IDDeck)
	if err != nil {
		t.Fatalf("admin debería ver baraja privada por RN08: %v", err)
	}
	if detailAdmin.Title != "Baraja Secreta de Gael" {
		t.Errorf("título incorrecto: %s", detailAdmin.Title)
	}

	// 6. Gael hace pública la baraja y Silas la guarda en su Studio
	_, err = studioSvc.UpdateDeck(ctx, gaelID, models.RoleLearner, privateDeck.IDDeck, models.UpdateDeckRequest{
		Title:    "Baraja Ya No Secreta",
		IsPublic: true,
	})
	if err != nil {
		t.Fatalf("error haciendo pública la baraja: %v", err)
	}

	err = studioSvc.SaveDeck(ctx, silasID, models.RoleEditor, privateDeck.IDDeck)
	if err != nil {
		t.Fatalf("Silas debería poder guardar la baraja ahora que es pública: %v", err)
	}

	// Silas consulta sus barajas de Studio y debe ver la guardada con is_saved_from_community = true
	silasStudio, err := studioSvc.ListStudioDecks(ctx, silasID)
	if err != nil {
		t.Fatalf("error listando barajas de Silas: %v", err)
	}
	found := false
	for _, sd := range silasStudio {
		if sd.IDDeck == privateDeck.IDDeck {
			found = true
			if !sd.IsSavedFromComm {
				t.Errorf("se esperaba is_saved_from_community = true, obtenido false")
			}
			if sd.OwnerUserName != "gael_learner" {
				t.Errorf("se esperaba owner 'gael_learner', obtenido '%s'", sd.OwnerUserName)
			}
		}
	}
	if !found {
		t.Errorf("la baraja guardada no apareció en el Studio de Silas")
	}
}

func TestStudioService_RN05_OfficialFlagSecurity(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	if testDB == nil {
		return
	}

	studioSvc := service.NewStudioService(testDB)
	ctx := context.Background()

	gaelID := "b1111111-0000-0000-0000-000000000003"  // LEARNER
	silasID := "b1111111-0000-0000-0000-000000000002" // EDITOR

	// 1. Gael (LEARNER) intenta marcar baraja como oficial -> Debe forzarse a false (RN05)
	deckLearner, err := studioSvc.CreateDeck(ctx, gaelID, models.RoleLearner, models.CreateDeckRequest{
		Title:      "Intento Oficial Learner",
		IsPublic:   true,
		IsOfficial: true,
	})
	if err != nil {
		t.Fatalf("CreateDeck falló: %v", err)
	}
	if deckLearner.IsOfficial {
		t.Fatalf("Violación RN05: Un LEARNER no debe poder crear barajas oficiales (obtenido is_official = true)")
	}

	// 2. Silas (EDITOR) crea baraja oficial -> is_official = true permitido
	deckEditor, err := studioSvc.CreateDeck(ctx, silasID, models.RoleEditor, models.CreateDeckRequest{
		Title:      "Colección Oficial de Silas",
		IsPublic:   true,
		IsOfficial: true,
	})
	if err != nil {
		t.Fatalf("CreateDeck editor falló: %v", err)
	}
	if !deckEditor.IsOfficial {
		t.Fatalf("se esperaba is_official = true para EDITOR, obtenido false")
	}

	// 3. Colecciones: Gael intenta crear con is_official e is_featured -> Se fuerzan a false
	colLearner, err := studioSvc.CreateCollection(ctx, gaelID, models.RoleLearner, models.CreateCollectionRequest{
		Title:      "Ruta Dialectal de Gael",
		IsPublic:   true,
		IsOfficial: true,
		IsFeatured: true,
	})
	if err != nil {
		t.Fatalf("CreateCollection falló: %v", err)
	}
	if colLearner.IsOfficial || colLearner.IsFeatured {
		t.Fatalf("Violación RN05: Un LEARNER no puede activar is_official ni is_featured")
	}
}

// ==============================================================================
// PRUEBAS DEL MOTOR SM-2 Y TRANSACCIÓN DE MAESTRÍA (CU12, CU13, RN09, RN11)
// ==============================================================================

func almostEqual(a, b float32) bool {
	return math.Abs(float64(a-b)) < 0.001
}

func TestMasteryService_CalculateSM2_MathematicalPrecision(t *testing.T) {
	now := time.Now().UTC()

	// 1. Primera revisión con dificultad 4 (Fácil):
	// q = 5, EF inicial = 2.5
	// Consecutive = 1 -> Interval = 1
	// EF' = 2.5 + (0.1 - 0) = 2.6
	// Stage = LEARNING, Mastery% = 40.0, XP = 20
	m1, xp1 := service.CalculateSM2(nil, 4, now)
	if m1.ConsecutiveCorrect != 1 {
		t.Errorf("se esperaba consecutive 1, obtenido %d", m1.ConsecutiveCorrect)
	}
	if m1.CurrentIntervalDays != 1 {
		t.Errorf("se esperaba intervalo 1, obtenido %d", m1.CurrentIntervalDays)
	}
	if !almostEqual(m1.EasinessFactor, 2.6) {
		t.Errorf("se esperaba EF 2.6, obtenido %f", m1.EasinessFactor)
	}
	if m1.ProgressStage != models.StageLearning {
		t.Errorf("se esperaba stage LEARNING, obtenido %s", m1.ProgressStage)
	}
	if xp1 != 20 {
		t.Errorf("se esperaba XP 20 para calificación Fácil (4), obtenido %d", xp1)
	}

	// 2. Segunda revisión consecutiva con dificultad 3 (Bueno):
	// q = 4, Consecutive = 2 -> Interval = 6
	// EF' = 2.6 + (0.1 - 1*(0.08 + 0.02)) = 2.6 + 0.0 = 2.6
	// Stage = LEARNING (consecutive < 3), XP = 15
	m2, xp2 := service.CalculateSM2(&m1, 3, now)
	if m2.ConsecutiveCorrect != 2 {
		t.Errorf("se esperaba consecutive 2, obtenido %d", m2.ConsecutiveCorrect)
	}
	if m2.CurrentIntervalDays != 6 {
		t.Errorf("se esperaba intervalo 6, obtenido %d", m2.CurrentIntervalDays)
	}
	if !almostEqual(m2.EasinessFactor, 2.6) {
		t.Errorf("se esperaba EF 2.6, obtenido %f", m2.EasinessFactor)
	}
	if m2.ProgressStage != models.StageLearning {
		t.Errorf("se esperaba stage LEARNING, obtenido %s", m2.ProgressStage)
	}
	if xp2 != 15 {
		t.Errorf("se esperaba XP 15 para calificación Bueno (3), obtenido %d", xp2)
	}

	// 3. Tercera revisión consecutiva con dificultad 4 (Fácil):
	// q = 5, Consecutive = 3 -> Interval = ceil(6 * 2.6) = ceil(15.6) = 16
	// EF' = 2.6 + 0.1 = 2.7
	// Stage = REVIEWING (consecutive >= 3 && interval >= 6), Mastery% = 75.0, XP = 20
	m3, xp3 := service.CalculateSM2(&m2, 4, now)
	if m3.ConsecutiveCorrect != 3 {
		t.Errorf("se esperaba consecutive 3, obtenido %d", m3.ConsecutiveCorrect)
	}
	if m3.CurrentIntervalDays != 16 {
		t.Errorf("se esperaba intervalo 16, obtenido %d", m3.CurrentIntervalDays)
	}
	if m3.ProgressStage != models.StageReviewing {
		t.Errorf("se esperaba stage REVIEWING, obtenido %s", m3.ProgressStage)
	}
	if m3.MasteryPercentage != 75.0 {
		t.Errorf("se esperaba 75%% mastery, obtenido %f", m3.MasteryPercentage)
	}
	if xp3 != 20 {
		t.Errorf("se esperaba XP 20, obtenido %d", xp3)
	}

	// 4. Cuarta revisión: Fallo (calificación 1):
	// q = 1 -> Consecutive reinicia a 0, Interval reinicia a 1
	// EF' = max(1.3, 2.7 - 0.20) = 2.5
	// Stage = LEARNING, Mastery% = 20.0, XP consolación = 5
	m4, xp4 := service.CalculateSM2(&m3, 1, now)
	if m4.ConsecutiveCorrect != 0 {
		t.Errorf("se esperaba consecutive 0 tras fallo, obtenido %d", m4.ConsecutiveCorrect)
	}
	if m4.CurrentIntervalDays != 1 {
		t.Errorf("se esperaba intervalo 1 tras fallo, obtenido %d", m4.CurrentIntervalDays)
	}
	if m4.ProgressStage != models.StageLearning {
		t.Errorf("se esperaba stage LEARNING tras fallo, obtenido %s", m4.ProgressStage)
	}
	if !almostEqual(m4.EasinessFactor, 2.5) {
		t.Errorf("se esperaba EF 2.5 tras fallo, obtenido %f", m4.EasinessFactor)
	}
	if xp4 != 5 {
		t.Errorf("se esperaba XP consolación 5 tras fallo, obtenido %d", xp4)
	}
}

func TestMasteryService_StudySessionAndReview_Transaction(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	if testDB == nil {
		return
	}

	masterySvc := service.NewMasteryService(testDB)
	ctx := context.Background()

	gaelID := "b1111111-0000-0000-0000-000000000003" // LEARNER

	// 1. Obtener sesión de estudio de la baraja oficial sembrada (Deck 1: Mexicanismos Indispensables)
	session, err := masterySvc.GetStudySession(ctx, gaelID, models.RoleLearner, utils.ToPtr(1), 10)
	if err != nil {
		t.Fatalf("GetStudySession falló: %v", err)
	}
	if session.TotalDue == 0 || len(session.SessionCards) == 0 {
		t.Fatalf("se esperaban tarjetas en la sesión para la baraja 1, obtenido TotalDue: %d", session.TotalDue)
	}

	firstCard := session.SessionCards[0]
	t.Logf("Estudiando primera tarjeta: sense_id=%d, lemma='%s'", firstCard.IDSense, firstCard.Lemma)

	// 2. Ejecutar revisión SM-2 (Fácil: 4) en transacción ACID
	reviewRes, err := masterySvc.SubmitReview(ctx, gaelID, models.SubmitReviewRequest{
		IDSense:       firstCard.IDSense,
		QualityRating: 4,
		TimeSpentMs:   2500,
	})
	if err != nil {
		t.Fatalf("SubmitReview falló: %v", err)
	}

	if reviewRes.XPAwarded != 20 {
		t.Errorf("se esperaba XPAwarded 20, obtenido %d", reviewRes.XPAwarded)
	}
	if reviewRes.UpdatedMastery.ConsecutiveCorrect != 1 {
		t.Errorf("se esperaba 1 acierto consecutivo, obtenido %d", reviewRes.UpdatedMastery.ConsecutiveCorrect)
	}
	if reviewRes.NewTotalXP <= 0 {
		t.Errorf("se esperaba NewTotalXP > 0, obtenido %d", reviewRes.NewTotalXP)
	}
	if reviewRes.CurrentStreak <= 0 {
		t.Errorf("se esperaba CurrentStreak > 0, obtenido %d", reviewRes.CurrentStreak)
	}

	// 3. Consultar tablero de maestría (MasteryDashboard)
	dashboard, err := masterySvc.GetMasteryDashboard(ctx, gaelID)
	if err != nil {
		t.Fatalf("GetMasteryDashboard falló: %v", err)
	}
	if dashboard.TotalWordsStudied < 1 {
		t.Errorf("se esperaba TotalWordsStudied >= 1, obtenido %d", dashboard.TotalWordsStudied)
	}
	if len(dashboard.RecentXP) == 0 {
		t.Fatalf("se esperaba al menos una transacción reciente de XP en el tablero")
	}
	if dashboard.RecentXP[0].AmountXP != 20 {
		t.Errorf("se esperaba RecentXP[0].AmountXP = 20, obtenido %d", dashboard.RecentXP[0].AmountXP)
	}
	if len(dashboard.Achievements) == 0 {
		t.Fatalf("se esperaba catálogo de logros en el tablero")
	}

	// 4. Consultar libro mayor de XP (RN11)
	ledger, err := masterySvc.ListXPLedger(ctx, gaelID, 10)
	if err != nil {
		t.Fatalf("ListXPLedger falló: %v", err)
	}
	if len(ledger) == 0 {
		t.Fatalf("el ledger inmutable de XP no contiene la transacción registrada")
	}
	if ledger[0].Reason != models.XPReasonQuizCorrect {
		t.Errorf("se esperaba razón QUIZ_CORRECT en el ledger, obtenido: %s", ledger[0].Reason)
	}

	// 5. Consultar ítems de maestría por etapa
	items, err := masterySvc.ListMasteryItems(ctx, gaelID, utils.ToPtr("LEARNING"))
	if err != nil {
		t.Fatalf("ListMasteryItems falló: %v", err)
	}
	if len(items) == 0 {
		t.Fatalf("se esperaba al menos una tarjeta en etapa LEARNING")
	}
	if items[0].IDSense != firstCard.IDSense {
		t.Errorf("se esperaba sense %d en LEARNING, obtenido %d", firstCard.IDSense, items[0].IDSense)
	}
}

func TestStudioService_SearchDecks_PostgresFTS(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	if testDB == nil {
		return
	}

	studioSvc := service.NewStudioService(testDB)
	ctx := context.Background()

	adminID := "b1111111-0000-0000-0000-000000000001" // ADMIN

	// Deck 1 ya existe en la semilla: 'Mexicanismos Indispensables' con descripción 'Vocabulario esencial y modismos de México que todo viajero o estudiante debe dominar.'
	// 1. Búsqueda por prefijo exacto sin acento
	results, err := studioSvc.SearchDecks(ctx, "mex", 10)
	if err != nil {
		t.Fatalf("SearchDecks('mex') falló: %v", err)
	}
	if len(results) == 0 {
		t.Fatalf("se esperaba al menos 1 resultado para 'mex', obtenido: %+v", results)
	}
	if results[0].Title != "Mexicanismos Indispensables" {
		t.Errorf("se esperaba 'Mexicanismos Indispensables', obtenido '%s'", results[0].Title)
	}
	if results[0].RegionName != "México" {
		t.Errorf("se esperaba RegionName 'México', obtenido '%s'", results[0].RegionName)
	}
	if results[0].OwnerUserName == "" {
		t.Errorf("se esperaba owner_user_name poblado, obtenido vacío")
	}

	// 2. Búsqueda con tildes / accent-insensitive ('méx' vs 'mex')
	resultsAccented, err := studioSvc.SearchDecks(ctx, "méx", 10)
	if err != nil {
		t.Fatalf("SearchDecks('méx') falló: %v", err)
	}
	if len(resultsAccented) == 0 || resultsAccented[0].IDDeck != results[0].IDDeck {
		t.Errorf("búsqueda insensible a acentos falló: %+v", resultsAccented)
	}

	// 3. Búsqueda FTS dentro de la descripción ('viajero' o 'estudiante')
	ftsResults, err := studioSvc.SearchDecks(ctx, "viajero", 10)
	if err != nil {
		t.Fatalf("SearchDecks('viajero') falló: %v", err)
	}
	if len(ftsResults) == 0 || ftsResults[0].IDDeck != 1 {
		t.Errorf("FTS en descripción de baraja falló: %+v", ftsResults)
	}

	// 4. Barajas privadas no deben aparecer en la búsqueda pública
	_, err = studioSvc.CreateDeck(ctx, adminID, models.RoleAdmin, models.CreateDeckRequest{
		Title:       "Secret Mexican Vault",
		Description: utils.ToPtr("Baraja secreta para pruebas privadas"),
		IsPublic:    false,
	})
	if err != nil {
		t.Fatalf("CreateDeck privado falló: %v", err)
	}

	secretResults, err := studioSvc.SearchDecks(ctx, "Secret", 10)
	if err != nil {
		t.Fatalf("SearchDecks('Secret') falló: %v", err)
	}
	if len(secretResults) > 0 {
		t.Errorf("las barajas privadas no deben aparecer en búsqueda pública, obtenido: %+v", secretResults)
	}

	// 5. Query vacía devuelve slice vacío sin error
	emptyResults, err := studioSvc.SearchDecks(ctx, "   ", 10)
	if err != nil {
		t.Fatalf("SearchDecks('') falló: %v", err)
	}
	if len(emptyResults) != 0 {
		t.Errorf("se esperaba resultado vacío para query en blanco, obtenido: %+v", emptyResults)
	}
}

func TestStudioService_SearchCollections_PostgresFTS(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	if testDB == nil {
		return
	}

	studioSvc := service.NewStudioService(testDB)
	ctx := context.Background()

	adminID := "b1111111-0000-0000-0000-000000000001" // ADMIN

	// Creamos una colección pública
	col, err := studioSvc.CreateCollection(ctx, adminID, models.RoleAdmin, models.CreateCollectionRequest{
		Title:       "Gastronomía Callejera Panhispánica",
		Description: utils.ToPtr("Ruta gastronómica de tacos, arepas, empanadas y comida callejera."),
		IsPublic:    true,
		IsFeatured:  true,
		IsOfficial:  true,
		DeckIDs:     []int{1},
	})
	if err != nil {
		t.Fatalf("CreateCollection falló: %v", err)
	}

	// 1. Prefijo sin tilde ('gastro' sobre 'Gastronomía')
	results, err := studioSvc.SearchCollections(ctx, "gastro", 10)
	if err != nil {
		t.Fatalf("SearchCollections('gastro') falló: %v", err)
	}
	if len(results) == 0 {
		t.Fatalf("se esperaba al menos 1 resultado para 'gastro', obtenido: %+v", results)
	}
	if results[0].IDCollection != col.IDCollection {
		t.Errorf("se esperaba colección id %d, obtenido %d", col.IDCollection, results[0].IDCollection)
	}
	if results[0].OwnerUserName == "" {
		t.Errorf("se esperaba owner_user_name poblado")
	}
	if results[0].DeckCount != 1 {
		t.Errorf("se esperaba deck_count = 1, obtenido %d", results[0].DeckCount)
	}

	// 2. Prefijo con tilde ('gastronomía')
	resultsAccented, err := studioSvc.SearchCollections(ctx, "gastronomía", 10)
	if err != nil {
		t.Fatalf("SearchCollections('gastronomía') falló: %v", err)
	}
	if len(resultsAccented) == 0 || resultsAccented[0].IDCollection != col.IDCollection {
		t.Errorf("búsqueda insensible a acentos falló: %+v", resultsAccented)
	}

	// 3. FTS en descripción ('empanadas' o 'callejera')
	ftsResults, err := studioSvc.SearchCollections(ctx, "empanadas", 10)
	if err != nil {
		t.Fatalf("SearchCollections('empanadas') falló: %v", err)
	}
	if len(ftsResults) == 0 || ftsResults[0].IDCollection != col.IDCollection {
		t.Errorf("FTS en descripción de colección falló: %+v", ftsResults)
	}

	// 4. Colección privada no debe aparecer
	_, err = studioSvc.CreateCollection(ctx, adminID, models.RoleAdmin, models.CreateCollectionRequest{
		Title:       "Colección Oculta Secreta",
		Description: utils.ToPtr("Descripción privada clasificada"),
		IsPublic:    false,
		DeckIDs:     []int{1},
	})
	if err != nil {
		t.Fatalf("CreateCollection privado falló: %v", err)
	}

	privateResults, err := studioSvc.SearchCollections(ctx, "Oculta", 10)
	if err != nil {
		t.Fatalf("SearchCollections('Oculta') falló: %v", err)
	}
	if len(privateResults) > 0 {
		t.Errorf("colección privada no debe aparecer en búsqueda pública: %+v", privateResults)
	}

	// 5. Query vacía devuelve slice vacío
	emptyResults, err := studioSvc.SearchCollections(ctx, "   ", 10)
	if err != nil {
		t.Fatalf("SearchCollections('') falló: %v", err)
	}
	if len(emptyResults) != 0 {
		t.Errorf("se esperaba slice vacío para query vacía, obtenido: %+v", emptyResults)
	}
}

func TestStudioService_DynamicRegionDetection(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	if testDB == nil {
		return
	}

	studioSvc := service.NewStudioService(testDB)
	ctx := context.Background()

	adminID := "b1111111-0000-0000-0000-000000000001" // ADMIN
	gaelID := "b1111111-0000-0000-0000-000000000003"  // LEARNER

	// 1. Seeded Deck 1 ('Mexicanismos Indispensables') contains only México senses.
	seededDeck, err := studioSvc.GetDeckDetail(ctx, &gaelID, utils.ToPtr(models.RoleLearner), 1)
	if err != nil {
		t.Fatalf("GetDeckDetail(1) failed: %v", err)
	}
	if !reflect.DeepEqual(seededDeck.RegionNames, []string{"México"}) {
		t.Errorf("expected seeded deck RegionNames [México], got %#v", seededDeck.RegionNames)
	}

	// 2. Empty deck guarantees non-nil []string{} across Create, Update, Get, ListPublic, ListStudio.
	emptyDeck, err := studioSvc.CreateDeck(ctx, gaelID, models.RoleLearner, models.CreateDeckRequest{
		Title:    "Baraja Vacía Sin Regiones",
		IsPublic: true,
	})
	if err != nil {
		t.Fatalf("CreateDeck empty failed: %v", err)
	}
	if emptyDeck.RegionNames == nil || len(emptyDeck.RegionNames) != 0 {
		t.Errorf("expected non-nil empty RegionNames on CreateDeck, got %#v", emptyDeck.RegionNames)
	}

	updatedEmptyDeck, err := studioSvc.UpdateDeck(ctx, gaelID, models.RoleLearner, emptyDeck.IDDeck, models.UpdateDeckRequest{
		Title:    "Baraja Vacía Actualizada",
		IsPublic: true,
	})
	if err != nil {
		t.Fatalf("UpdateDeck empty failed: %v", err)
	}
	if updatedEmptyDeck.RegionNames == nil || len(updatedEmptyDeck.RegionNames) != 0 {
		t.Errorf("expected non-nil empty RegionNames on UpdateDeck, got %#v", updatedEmptyDeck.RegionNames)
	}

	emptyDeckDetail, err := studioSvc.GetDeckDetail(ctx, &gaelID, utils.ToPtr(models.RoleLearner), emptyDeck.IDDeck)
	if err != nil {
		t.Fatalf("GetDeckDetail empty failed: %v", err)
	}
	if emptyDeckDetail.RegionNames == nil || len(emptyDeckDetail.RegionNames) != 0 {
		t.Errorf("expected non-nil empty RegionNames on GetDeckDetail, got %#v", emptyDeckDetail.RegionNames)
	}

	// 3. Populated multi-region deck (>3 regions) verifies count DESC, name ASC, LIMIT 3:
	// Senses: 101(MX), 102(MX), 103(CU), 104(CL), 106(ES) => MX:2, CL:1, CU:1, ES:1 -> top 3: [México, Chile, Cuba]
	multiDeck, err := studioSvc.CreateDeck(ctx, adminID, models.RoleAdmin, models.CreateDeckRequest{
		Title:    "Baraja Multirregional",
		IsPublic: true,
		SenseIDs: []int{101, 102, 103, 104, 106},
	})
	if err != nil {
		t.Fatalf("CreateDeck multi failed: %v", err)
	}

	// Secondary deck: 104(CL), 106(ES) => [Chile, España]
	secondaryDeck, err := studioSvc.CreateDeck(ctx, adminID, models.RoleAdmin, models.CreateDeckRequest{
		Title:    "Baraja Chile y España",
		IsPublic: true,
		SenseIDs: []int{104, 106},
	})
	if err != nil {
		t.Fatalf("CreateDeck secondary failed: %v", err)
	}

	expectedMulti := []string{"México", "Chile", "Cuba"}
	multiDetail, err := studioSvc.GetDeckDetail(ctx, &adminID, utils.ToPtr(models.RoleAdmin), multiDeck.IDDeck)
	if err != nil {
		t.Fatalf("GetDeckDetail multi failed: %v", err)
	}
	if !reflect.DeepEqual(multiDetail.RegionNames, expectedMulti) {
		t.Errorf("expected GetDeckDetail RegionNames %v, got %#v", expectedMulti, multiDetail.RegionNames)
	}

	publicDecks, err := studioSvc.ListPublicDecks(ctx, nil, nil)
	if err != nil {
		t.Fatalf("ListPublicDecks failed: %v", err)
	}
	foundMultiPublic := false
	foundEmptyPublic := false
	for _, d := range publicDecks {
		if d.IDDeck == multiDeck.IDDeck {
			foundMultiPublic = true
			if !reflect.DeepEqual(d.RegionNames, expectedMulti) {
				t.Errorf("ListPublicDecks multi RegionNames = %#v, want %v", d.RegionNames, expectedMulti)
			}
		}
		if d.IDDeck == emptyDeck.IDDeck {
			foundEmptyPublic = true
			if d.RegionNames == nil || len(d.RegionNames) != 0 {
				t.Errorf("ListPublicDecks empty RegionNames = %#v, want []string{}", d.RegionNames)
			}
		}
	}
	if !foundMultiPublic || !foundEmptyPublic {
		t.Fatalf("expected both multi (%v) and empty (%v) decks in ListPublicDecks", foundMultiPublic, foundEmptyPublic)
	}

	studioDecks, err := studioSvc.ListStudioDecks(ctx, adminID)
	if err != nil {
		t.Fatalf("ListStudioDecks failed: %v", err)
	}
	foundMultiStudio := false
	for _, sd := range studioDecks {
		if sd.IDDeck == multiDeck.IDDeck {
			foundMultiStudio = true
			if !reflect.DeepEqual(sd.RegionNames, expectedMulti) {
				t.Errorf("ListStudioDecks multi RegionNames = %#v, want %v", sd.RegionNames, expectedMulti)
			}
		}
	}
	if !foundMultiStudio {
		t.Fatalf("expected multiDeck in ListStudioDecks for admin")
	}

	// 4. Empty collection guarantees non-nil []string{} across Create, Update, Get, ListPublic, ListStudio.
	emptyCol, err := studioSvc.CreateCollection(ctx, gaelID, models.RoleLearner, models.CreateCollectionRequest{
		Title:    "Colección Vacía",
		IsPublic: true,
	})
	if err != nil {
		t.Fatalf("CreateCollection empty failed: %v", err)
	}
	if emptyCol.RegionNames == nil || len(emptyCol.RegionNames) != 0 {
		t.Errorf("expected non-nil empty RegionNames on CreateCollection, got %#v", emptyCol.RegionNames)
	}

	updatedEmptyCol, err := studioSvc.UpdateCollection(ctx, gaelID, models.RoleLearner, emptyCol.IDCollection, models.UpdateCollectionRequest{
		Title:    "Colección Vacía Editada",
		IsPublic: true,
	})
	if err != nil {
		t.Fatalf("UpdateCollection empty failed: %v", err)
	}
	if updatedEmptyCol.RegionNames == nil || len(updatedEmptyCol.RegionNames) != 0 {
		t.Errorf("expected non-nil empty RegionNames on UpdateCollection, got %#v", updatedEmptyCol.RegionNames)
	}

	// 5. Populated collection aggregating multiDeck (MX:2, CL:1, CU:1, ES:1) + secondaryDeck (CL:1, ES:1)
	// Aggregated counts: Chile:2, España:2, México:2, Cuba:1 -> alphabetical tiebreak LIMIT 3: [Chile, España, México]
	popCol, err := studioSvc.CreateCollection(ctx, adminID, models.RoleAdmin, models.CreateCollectionRequest{
		Title:      "Ruta Dialectal Variada",
		IsPublic:   true,
		IsOfficial: true,
		DeckIDs:    []int{multiDeck.IDDeck, secondaryDeck.IDDeck},
	})
	if err != nil {
		t.Fatalf("CreateCollection populated failed: %v", err)
	}

	expectedColRegions := []string{"Chile", "España", "México"}
	colDetail, err := studioSvc.GetCollectionDetail(ctx, &adminID, utils.ToPtr(models.RoleAdmin), popCol.IDCollection)
	if err != nil {
		t.Fatalf("GetCollectionDetail populated failed: %v", err)
	}
	if !reflect.DeepEqual(colDetail.RegionNames, expectedColRegions) {
		t.Errorf("GetCollectionDetail RegionNames = %#v, want %v", colDetail.RegionNames, expectedColRegions)
	}
	if len(colDetail.Decks) != 2 {
		t.Fatalf("expected 2 nested decks, got %d", len(colDetail.Decks))
	}
	if !reflect.DeepEqual(colDetail.Decks[0].RegionNames, expectedMulti) {
		t.Errorf("nested deck[0] RegionNames = %#v, want %v", colDetail.Decks[0].RegionNames, expectedMulti)
	}
	if !reflect.DeepEqual(colDetail.Decks[1].RegionNames, []string{"Chile", "España"}) {
		t.Errorf("nested deck[1] RegionNames = %#v, want [Chile España]", colDetail.Decks[1].RegionNames)
	}

	publicCols, err := studioSvc.ListPublicCollections(ctx, nil)
	if err != nil {
		t.Fatalf("ListPublicCollections failed: %v", err)
	}
	foundPopPublic := false
	foundEmptyColPublic := false
	for _, c := range publicCols {
		if c.IDCollection == popCol.IDCollection {
			foundPopPublic = true
			if !reflect.DeepEqual(c.RegionNames, expectedColRegions) {
				t.Errorf("ListPublicCollections pop RegionNames = %#v, want %v", c.RegionNames, expectedColRegions)
			}
		}
		if c.IDCollection == emptyCol.IDCollection {
			foundEmptyColPublic = true
			if c.RegionNames == nil || len(c.RegionNames) != 0 {
				t.Errorf("ListPublicCollections empty RegionNames = %#v, want []string{}", c.RegionNames)
			}
		}
	}
	if !foundPopPublic || !foundEmptyColPublic {
		t.Fatalf("expected both populated (%v) and empty (%v) collections in ListPublicCollections", foundPopPublic, foundEmptyColPublic)
	}

	studioCols, err := studioSvc.ListStudioCollections(ctx, adminID)
	if err != nil {
		t.Fatalf("ListStudioCollections failed: %v", err)
	}
	foundPopStudio := false
	for _, c := range studioCols {
		if c.IDCollection == popCol.IDCollection {
			foundPopStudio = true
			if !reflect.DeepEqual(c.RegionNames, expectedColRegions) {
				t.Errorf("ListStudioCollections pop RegionNames = %#v, want %v", c.RegionNames, expectedColRegions)
			}
		}
	}
	if !foundPopStudio {
		t.Fatalf("expected popCol in ListStudioCollections for admin")
	}
}

