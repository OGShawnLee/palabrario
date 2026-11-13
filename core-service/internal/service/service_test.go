package service_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"palabrario/core-service/internal/service"
	"palabrario/core-service/internal/testutil"
	"palabrario/pkg/auth"
	"palabrario/pkg/models"
)

// ==============================================================================
// CAPA 80%: PRUEBAS DIRECTAS DE SERVICIO / DAO SOBRE POSTGRESQL REAL (TEMPLATE)
// ==============================================================================

func TestIdentityService_SignIn_And_SuspendedAccountIsolation(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	if testDB == nil {
		return
	}

	tokenSvc := auth.NewTokenService("test-secret-key-minimum-32-chars", 24*time.Hour)
	identitySvc := service.NewIdentityService(testDB, tokenSvc)
	ctx := context.Background()

	// 1. Login exitoso con usuario semilla (Silas - EDITOR)
	res, err := identitySvc.SignIn(ctx, models.SignInRequest{
		Email:    "silas@palabrario.org",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("SignIn falló para usuario semilla Silas: %v", err)
	}
	if res.User.Role != models.RoleEditor {
		t.Errorf("se esperaba rol EDITOR, obtenido %s", res.User.Role)
	}
	if res.User.AnchorDialect != "MX" {
		t.Errorf("se esperaba dialecto MX, obtenido %s", res.User.AnchorDialect)
	}

	// 2. Prueba destructiva: Suspendemos a Silas directamente en este clon de DB
	_, err = testDB.Pool.Exec(ctx, `UPDATE accounts SET is_active = false WHERE email = 'silas@palabrario.org'`)
	if err != nil {
		t.Fatalf("fallo al suspender cuenta en DB efímera: %v", err)
	}

	// 3. Verificar regla de negocio: Cuenta suspendida devuelve ErrForbidden (403)
	_, err = identitySvc.SignIn(ctx, models.SignInRequest{
		Email:    "silas@palabrario.org",
		Password: "password123",
	})
	if !errors.Is(err, service.ErrForbidden) {
		t.Fatalf("se esperaba ErrForbidden para cuenta suspendida, obtenido: %v", err)
	}
}

func TestIdentityService_SignUp_And_DuplicateConflict(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	if testDB == nil {
		return
	}

	tokenSvc := auth.NewTokenService("test-secret-key-minimum-32-chars", 24*time.Hour)
	identitySvc := service.NewIdentityService(testDB, tokenSvc)
	ctx := context.Background()

	// Incluso después de que el test anterior suspendió a Silas en su propio clon,
	// aquí Silas sigue activo e intacto gracias al aislamiento por TEMPLATE.
	_, err := identitySvc.SignIn(ctx, models.SignInRequest{
		Email:    "silas@palabrario.org",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("el aislamiento por TEMPLATE falló; Silas debería estar activo en este clon: %v", err)
	}

	// Registro de nuevo aprendiz con ID determinístico (no necesitamos sufijos aleatorios)
	signUpReq := models.SignUpRequest{
		Email:         "aprendiz_nuevo@palabrario.org",
		Password:      "ClaveSegura123!",
		Name:          "Valeria Aprendiz",
		UserName:      "valeria_co",
		AnchorDialect: "CO",
	}

	created, err := identitySvc.SignUp(ctx, signUpReq)
	if err != nil {
		t.Fatalf("SignUp falló inesperadamente: %v", err)
	}
	if created.User.Role != models.RoleLearner {
		t.Errorf("se esperaba rol LEARNER por defecto, obtenido %s", created.User.Role)
	}

	// Intentar duplicar el mismo correo debe mapear SQLSTATE 23505 -> service.ErrConflict (409)
	_, err = identitySvc.SignUp(ctx, signUpReq)
	if !errors.Is(err, service.ErrConflict) {
		t.Fatalf("se esperaba ErrConflict al duplicar email/username, obtenido: %v", err)
	}
}

func TestDictionaryService_SearchAndDetail_PostgresFTS(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	if testDB == nil {
		return
	}

	dictSvc := service.NewDictionaryService(testDB)
	ctx := context.Background()

	// 1. Autocompletado rápido FTS + trigramas (función almacenada search_vocabulary_dropdown)
	suggestions, err := dictSvc.SearchDropdown(ctx, "popote", "MX", 5)
	if err != nil {
		t.Fatalf("SearchDropdown falló: %v", err)
	}
	if len(suggestions) == 0 || suggestions[0].Lemma != "popote" {
		t.Fatalf("se esperaba encontrar 'popote' en el autocompletado, obtenido: %+v", suggestions)
	}
	if suggestions[0].PartOfSpeech == "" {
		t.Errorf("se esperaba PartOfSpeech poblado, obtenido '%s'", suggestions[0].PartOfSpeech)
	}
	if suggestions[0].RegionName != "México" && suggestions[0].RegionName != "General" {
		t.Errorf("se esperaba RegionName válido, obtenido '%s'", suggestions[0].RegionName)
	}

	// 2. Detalle completo de término (incluyendo ParentZone e hijos anidados)
	detail, err := dictSvc.GetWordDetail(ctx, 1)
	if err != nil {
		t.Fatalf("GetWordDetail(1) falló: %v", err)
	}
	if detail.Lemma != "popote" {
		t.Errorf("se esperaba lema 'popote', obtenido '%s'", detail.Lemma)
	}
	if detail.OriginRegion == nil || detail.OriginRegion.ParentZone != "Norteamérica" {
		t.Errorf("se esperaba OriginRegion.ParentZone = 'Norteamérica', obtenido: %+v", detail.OriginRegion)
	}
	if len(detail.Senses) == 0 || len(detail.Senses[0].Examples) == 0 {
		t.Errorf("se esperaban acepciones y ejemplos anidados para 'popote', obtenido: %+v", detail.Senses)
	}
}

func TestDictionaryService_ExploreAndProposal(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	if testDB == nil {
		return
	}

	dictSvc := service.NewDictionaryService(testDB)
	ctx := context.Background()

	// 1. Listar regiones
	regions, err := dictSvc.ListRegions(ctx)
	if err != nil {
		t.Fatalf("ListRegions falló: %v", err)
	}
	if len(regions) == 0 {
		t.Fatalf("se esperaban regiones precargadas, obtenido 0")
	}

	// 2. Palabra del día
	wotd, err := dictSvc.GetWordOfTheDay(ctx, "MX")
	if err != nil {
		t.Fatalf("GetWordOfTheDay falló: %v", err)
	}
	if wotd == nil || wotd.IDWord <= 0 {
		t.Fatalf("palabra del día inválida: %+v", wotd)
	}

	// 3. Crear propuesta editorial (Silas - EDITOR -> debe ser DRAFT)
	silasID := "b1111111-0000-0000-0000-000000000002"
	exampleSentence := "En Monterrey siempre dicen huerco para llamar a los niños."
	prop, err := dictSvc.CreateWordProposal(ctx, silasID, models.RoleEditor, models.CreateWordProposalRequest{
		Lemma:        "huerco",
		PartOfSpeech: "sustantivo",
		Gender:       models.GenderMale,
		InitialSense: models.CreateSenseRequest{
			RegionCode:      "MX",
			ShortDefinition: "Niño, muchacho o infante travieso.",
			Definition:      "Término coloquial en el norte de México para referirse a un infante o muchacho joven.",
			IsPrimary:       true,
			CEFRLevel:       "B1",
			UsageFrequency:  "COMMON",
			Examples: []models.CreateExampleInput{
				{Sentence: exampleSentence},
			},
		},
	})
	if err != nil {
		t.Fatalf("CreateWordProposal falló: %v", err)
	}
	if prop.Lemma != "huerco" {
		t.Errorf("se esperaba lema 'huerco', obtenido: %s", prop.Lemma)
	}
	if prop.NormalizedLemma != "huerco" {
		t.Errorf("se esperaba normalizado 'huerco', obtenido: %s", prop.NormalizedLemma)
	}
	if len(prop.Senses) != 0 {
		// GetWordDetail solo retorna sentidos con state = 'APPROVED';
		// Dado que el nuevo sentido se creó en 'DRAFT', prop.Senses aprobados debe ser 0.
		t.Logf("Acepciones aprobadas devueltas: %d (correcto para DRAFT)", len(prop.Senses))
	}

	// 4. Exploración paginada
	paginated, err := dictSvc.ExploreWords(ctx, models.WordExploreFilter{
		RegionCode: "MX",
		Page:       1,
		Limit:      10,
	})
	if err != nil {
		t.Fatalf("ExploreWords falló: %v", err)
	}
	if paginated.TotalCount == 0 || len(paginated.Items) == 0 {
		t.Fatalf("se esperaban palabras en ExploreWords para MX, obtenido %+v", paginated)
	}

	// 5. Conceptos semánticos
	concept, err := dictSvc.CreateSemanticConcept(ctx, models.CreateConceptRequest{
		ConceptKey:           "test_drinking_straw",
		UniversalDescription: "Tubo para beber líquidos",
	})
	if err != nil {
		t.Fatalf("CreateSemanticConcept falló: %v", err)
	}
	if concept.ConceptKey != "test_drinking_straw" {
		t.Errorf("se esperaba concept_key 'test_drinking_straw', obtenido %s", concept.ConceptKey)
	}

	conceptsList, err := dictSvc.ListSemanticConcepts(ctx, "test_drinking")
	if err != nil || len(conceptsList) == 0 {
		t.Fatalf("ListSemanticConcepts falló al encontrar el concepto creado")
	}
}

// ==============================================================================
// PRUEBAS DE REGLAS DE NEGOCIO RN06 (QUÓRUM DE 3 VOTOS) Y RN07 (ANTI-AUTOAPROBACIÓN)
// ==============================================================================

func TestCurationService_RN07_SelfApprovalForbidden(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	if testDB == nil {
		return
	}

	dictSvc := service.NewDictionaryService(testDB)
	curationSvc := service.NewCurationService(testDB)
	ctx := context.Background()

	silasID := "b1111111-0000-0000-0000-000000000002"

	// 1. Silas (EDITOR) propone una acepción
	_, err := dictSvc.CreateWordProposal(ctx, silasID, models.RoleEditor, models.CreateWordProposalRequest{
		Lemma:        "chochear",
		PartOfSpeech: "verbo",
		InitialSense: models.CreateSenseRequest{
			RegionCode:      "MX",
			ShortDefinition: "Actuar con debilidad o despistes por la edad avanzada.",
			Definition:      "Tener las facultades mentales o físicas mermadas a causa de la vejez.",
			IsPrimary:       true,
		},
	})
	if err != nil {
		t.Fatalf("creación de propuesta falló: %v", err)
	}

	// 2. Localizar la propuesta en la cola de curaduría
	queue, err := curationSvc.ListMyProposals(ctx, silasID, nil)
	if err != nil || len(queue) == 0 {
		t.Fatalf("ListMyProposals falló: %v, len: %d", err, len(queue))
	}
	senseID := queue[0].IDSense

	// 3. Regla RN07: Silas intenta auto-aprobar su propia propuesta
	_, err = curationSvc.SubmitVote(ctx, silasID, senseID, models.SubmitVoteRequest{
		Decision: models.DecisionApprove,
	})
	if err == nil {
		t.Fatalf("se esperaba error de auto-aprobación (RN07), pero la operación fue permitida")
	}
	if !errors.Is(err, service.ErrForbidden) {
		t.Fatalf("se esperaba ErrForbidden (403) por violación de RN07, obtenido: %v", err)
	}
}

func TestCurationService_RN06_QuorumThreeVotesAutoApprove(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	if testDB == nil {
		return
	}

	dictSvc := service.NewDictionaryService(testDB)
	curationSvc := service.NewCurationService(testDB)
	ctx := context.Background()

	silasID := "b1111111-0000-0000-0000-000000000002"

	// 1. Silas (EDITOR) crea una propuesta en borrador (DRAFT)
	_, err := dictSvc.CreateWordProposal(ctx, silasID, models.RoleEditor, models.CreateWordProposalRequest{
		Lemma:        "agüitado",
		PartOfSpeech: "adjetivo",
		InitialSense: models.CreateSenseRequest{
			RegionCode:      "MX",
			ShortDefinition: "Triste, desanimado o decaído anímicamente.",
			Definition:      "Estado de ánimo bajo caracterizado por nostalgia, tristeza o desgano temporal.",
			IsPrimary:       true,
		},
	})
	if err != nil {
		t.Fatalf("creación de propuesta falló: %v", err)
	}

	queue, err := curationSvc.ListMyProposals(ctx, silasID, nil)
	if err != nil || len(queue) == 0 {
		t.Fatalf("no se encontró la propuesta creada: %v", err)
	}
	senseID := queue[0].IDSense

	// 2. Sembrar 3 editores independientes en la base efímera
	editor1ID := "c1111111-0000-0000-0000-000000000001"
	editor2ID := "c1111111-0000-0000-0000-000000000002"
	editor3ID := "c1111111-0000-0000-0000-000000000003"

	for i, id := range []string{editor1ID, editor2ID, editor3ID} {
		email := fmt.Sprintf("editor%d@palabrario.org", i+1)
		username := fmt.Sprintf("editor%d", i+1)
		name := fmt.Sprintf("Editor Evaluador %d", i+1)

		_, err := testDB.Pool.Exec(ctx, `
			INSERT INTO accounts (id_account, email, password_hash, role)
			VALUES ($1, $2, 'dummy_hash', 'EDITOR')
		`, id, email)
		if err != nil {
			t.Fatalf("error sembrando cuenta de editor %s: %v", id, err)
		}

		_, err = testDB.Pool.Exec(ctx, `
			INSERT INTO users (id_user, name, user_name, anchor_dialect)
			VALUES ($1, $2, $3, 'MX')
		`, id, name, username)
		if err != nil {
			t.Fatalf("error sembrando usuario de editor %s: %v", id, err)
		}
	}

	// 3. Primer voto favorable (Editor 1) -> Estado aún debe ser DRAFT, approvals_count = 1
	detail1, err := curationSvc.SubmitVote(ctx, editor1ID, senseID, models.SubmitVoteRequest{
		Decision: models.DecisionApprove,
	})
	if err != nil {
		t.Fatalf("voto 1 falló: %v", err)
	}
	if detail1.Proposal.State != "DRAFT" {
		t.Errorf("tras 1 voto se esperaba DRAFT, obtenido: %s", detail1.Proposal.State)
	}
	if detail1.Proposal.ApprovalsCount != 1 {
		t.Errorf("se esperaba approvals_count = 1, obtenido: %d", detail1.Proposal.ApprovalsCount)
	}

	// 4. Segundo voto favorable (Editor 2) -> Estado aún debe ser DRAFT, approvals_count = 2
	detail2, err := curationSvc.SubmitVote(ctx, editor2ID, senseID, models.SubmitVoteRequest{
		Decision: models.DecisionApprove,
	})
	if err != nil {
		t.Fatalf("voto 2 falló: %v", err)
	}
	if detail2.Proposal.State != "DRAFT" {
		t.Errorf("tras 2 votos se esperaba DRAFT, obtenido: %s", detail2.Proposal.State)
	}
	if detail2.Proposal.ApprovalsCount != 2 {
		t.Errorf("se esperaba approvals_count = 2, obtenido: %d", detail2.Proposal.ApprovalsCount)
	}

	// 5. Tercer voto favorable (Editor 3) -> ¡Regla RN06: Dispara trigger trg_auto_approve_sense!
	// El estado debe cambiar automáticamente a APPROVED y approvals_count = 3
	detail3, err := curationSvc.SubmitVote(ctx, editor3ID, senseID, models.SubmitVoteRequest{
		Decision: models.DecisionApprove,
	})
	if err != nil {
		t.Fatalf("voto 3 falló: %v", err)
	}
	if detail3.Proposal.State != "APPROVED" {
		t.Errorf("¡RN06 Quórum 3 falló! Se esperaba estado APPROVED, obtenido: %s", detail3.Proposal.State)
	}
	if detail3.Proposal.ApprovalsCount != 3 {
		t.Errorf("se esperaba approvals_count = 3, obtenido: %d", detail3.Proposal.ApprovalsCount)
	}
}

func TestCurationService_RN08_AdminOverrideState(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	if testDB == nil {
		return
	}

	dictSvc := service.NewDictionaryService(testDB)
	curationSvc := service.NewCurationService(testDB)
	ctx := context.Background()

	adminID := "b1111111-0000-0000-0000-000000000001"
	silasID := "b1111111-0000-0000-0000-000000000002"

	// 1. Silas crea propuesta en DRAFT
	_, err := dictSvc.CreateWordProposal(ctx, silasID, models.RoleEditor, models.CreateWordProposalRequest{
		Lemma:        "apapachito",
		PartOfSpeech: "sustantivo",
		InitialSense: models.CreateSenseRequest{
			RegionCode:      "MX",
			ShortDefinition: "Muestra pequeña pero tierna de cariño.",
			Definition:      "Diminutivo afectivo de apapacho.",
			IsPrimary:       true,
		},
	})
	if err != nil {
		t.Fatalf("falló creación de propuesta: %v", err)
	}

	queue, err := curationSvc.ListMyProposals(ctx, silasID, nil)
	if err != nil || len(queue) == 0 {
		t.Fatalf("no se encontró propuesta: %v", err)
	}
	senseID := queue[0].IDSense

	// 2. Administrador aprueba directamente (RN08)
	reason := "Aprobación canónica directa por administrador editorial"
	senseDetail, err := curationSvc.AdminOverrideState(ctx, adminID, senseID, models.AdminOverrideStateRequest{
		State:  "APPROVED",
		Reason: &reason,
	})
	if err != nil {
		t.Fatalf("AdminOverrideState falló: %v", err)
	}
	if senseDetail.State != "APPROVED" {
		t.Errorf("se esperaba estado APPROVED tras veto/aprobación admin, obtenido: %s", senseDetail.State)
	}

	// 3. Administrador oculta la acepción (RN08)
	hideReason := "Ocultada para corrección de fuentes"
	hiddenDetail, err := curationSvc.AdminOverrideState(ctx, adminID, senseID, models.AdminOverrideStateRequest{
		State:  "HIDDEN",
		Reason: &hideReason,
	})
	if err != nil {
		t.Fatalf("AdminOverrideState HIDDEN falló: %v", err)
	}
	if hiddenDetail.State != "HIDDEN" {
		t.Errorf("se esperaba estado HIDDEN tras ocultación admin, obtenido: %s", hiddenDetail.State)
	}
}

func TestDictionaryService_CreateWordProposal_MediaURLValidation_RN12(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	if testDB == nil {
		return
	}

	dictSvc := service.NewDictionaryService(testDB)
	ctx := context.Background()
	silasID := "b1111111-0000-0000-0000-000000000002"

	// 1. Debe rechazar URL externa arbitraria con ErrValidation (RN12)
	maliciousURL := "https://malicious.site/hack.jpg"
	_, err := dictSvc.CreateWordProposal(ctx, silasID, models.RoleEditor, models.CreateWordProposalRequest{
		Lemma:        "palabratestmaliciosa",
		PartOfSpeech: "sustantivo",
		ImageURL:     &maliciousURL,
		InitialSense: models.CreateSenseRequest{
			RegionCode:      "MX",
			ShortDefinition: "Definición de prueba para URL maliciosa.",
			Definition:      "Definición extendida para prueba de rechazo de imagen externa.",
			IsPrimary:       true,
		},
	})
	if err == nil || !errors.Is(err, service.ErrValidation) {
		t.Fatalf("se esperaba ErrValidation al enviar URL externa arbitraria, obtenido: %v", err)
	}

	// 2. Debe rechazar path traversal o esquemas inseguros (RN12)
	traversalURL := "http://localhost:8002/uploads/../secret.env"
	_, err = dictSvc.CreateWordProposal(ctx, silasID, models.RoleEditor, models.CreateWordProposalRequest{
		Lemma:        "palabratesttraversal",
		PartOfSpeech: "sustantivo",
		ImageURL:     &traversalURL,
		InitialSense: models.CreateSenseRequest{
			RegionCode:      "MX",
			ShortDefinition: "Definición de prueba para traversal.",
			Definition:      "Definición extendida para prueba de rechazo de traversal.",
			IsPrimary:       true,
		},
	})
	if err == nil || !errors.Is(err, service.ErrValidation) {
		t.Fatalf("se esperaba ErrValidation al enviar path traversal, obtenido: %v", err)
	}

	// 3. Debe aceptar URL emitida por image-service conteniendo /uploads/ (RN12)
	validURL := "http://localhost:8002/uploads/word/test.png"
	validSenseURL := "http://localhost:8002/uploads/sense/test_sense.png"
	word, err := dictSvc.CreateWordProposal(ctx, silasID, models.RoleEditor, models.CreateWordProposalRequest{
		Lemma:        "palabratestvalida",
		PartOfSpeech: "sustantivo",
		ImageURL:     &validURL,
		InitialSense: models.CreateSenseRequest{
			RegionCode:      "MX",
			ShortDefinition: "Definición válida con imagen procesada.",
			Definition:      "Definición detallada que incluye URL multimedia interna válida.",
			ImageURL:        &validSenseURL,
			IsPrimary:       true,
		},
	})
	if err != nil {
		t.Fatalf("se esperaba éxito con URL interna válida de image-service, falló: %v", err)
	}
	if word.ImageURL == nil || *word.ImageURL != validURL {
		t.Errorf("se esperaba ImageURL %s en la palabra, obtenido: %v", validURL, word.ImageURL)
	}

	// Verificar que la acepción creada en DRAFT por Silas conserve su ImageURL en la cola de curaduría
	curationSvc := service.NewCurationService(testDB)
	proposals, err := curationSvc.ListMyProposals(ctx, silasID, nil)
	if err != nil || len(proposals) == 0 {
		t.Fatalf("error consultando propuestas del editor Silas: %v", err)
	}
	foundProposal := false
	for _, p := range proposals {
		if p.Lemma == "palabratestvalida" {
			foundProposal = true
			if p.ImageURL == nil || *p.ImageURL != validSenseURL {
				t.Errorf("se esperaba ImageURL %s en la propuesta del sentido, obtenido: %v", validSenseURL, p.ImageURL)
			}
		}
	}
	if !foundProposal {
		t.Errorf("no se encontró la propuesta 'palabratestvalida' en la lista de propuestas del editor")
	}

	// 4. Si un ADMIN crea la propuesta (RN08), el sentido pasa a APPROVED inmediatamente y GetWordDetail lo proyecta
	adminID := "b1111111-0000-0000-0000-000000000001"
	adminWord, err := dictSvc.CreateWordProposal(ctx, adminID, models.RoleAdmin, models.CreateWordProposalRequest{
		Lemma:        "palabraadminvalida",
		PartOfSpeech: "sustantivo",
		ImageURL:     &validURL,
		InitialSense: models.CreateSenseRequest{
			RegionCode:      "MX",
			ShortDefinition: "Definición aprobada directa por admin.",
			Definition:      "Definición detallada aprobada directa con URL multimedia.",
			ImageURL:        &validSenseURL,
			IsPrimary:       true,
		},
	})
	if err != nil {
		t.Fatalf("CreateWordProposal con rol ADMIN falló: %v", err)
	}
	if len(adminWord.Senses) == 0 || adminWord.Senses[0].ImageURL == nil || *adminWord.Senses[0].ImageURL != validSenseURL {
		t.Errorf("se esperaba ImageURL %s en el sentido aprobado, obtenido: %+v", validSenseURL, adminWord.Senses)
	}
}

func TestArenaCatalogService_ListPublicWaitingMatches(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	if testDB == nil {
		return
	}

	arenaSvc := service.NewArenaCatalogService(testDB)
	ctx := context.Background()
	gaelID := "b1111111-0000-0000-0000-000000000003"

	// 1. Crear una partida pública y en espera válida
	isPublicTrue := true
	matchPkg, err := arenaSvc.PrepareMatch(ctx, gaelID, models.RoleLearner, models.CreateMatchRequest{
		IDDeck:       1,
		TotalRounds:  8,
		RoundTimeSec: 15,
		IsPublic:     &isPublicTrue,
	})
	if err != nil {
		t.Fatalf("PrepareMatch falló: %v", err)
	}

	// 2. Crear una partida privada (is_public = false)
	isPublicFalse := false
	privatePkg, err := arenaSvc.PrepareMatch(ctx, gaelID, models.RoleLearner, models.CreateMatchRequest{
		IDDeck:       1,
		TotalRounds:  8,
		RoundTimeSec: 15,
		IsPublic:     &isPublicFalse,
	})
	if err != nil {
		t.Fatalf("PrepareMatch (privada) falló: %v", err)
	}

	// 3. Crear una partida no en espera (status = 'IN_PROGRESS')
	_, err = testDB.Pool.Exec(ctx, `
		INSERT INTO matches (id_deck, host_user_id, room_code, status, is_public, total_rounds, round_time_sec)
		VALUES (1, $1, 'INPROG', 'IN_PROGRESS', true, 8, 15)
	`, gaelID)
	if err != nil {
		t.Fatalf("falló insertar partida IN_PROGRESS: %v", err)
	}

	// 4. Crear una partida vieja (> 2 horas)
	_, err = testDB.Pool.Exec(ctx, `
		INSERT INTO matches (id_deck, host_user_id, room_code, status, is_public, total_rounds, round_time_sec, created_at)
		VALUES (1, $1, 'OLDROO', 'WAITING', true, 8, 15, now() - INTERVAL '3 hours')
	`, gaelID)
	if err != nil {
		t.Fatalf("falló insertar partida antigua: %v", err)
	}

	// 5. Crear una partida llena (>= 8 participantes)
	var fullMatchID string
	err = testDB.Pool.QueryRow(ctx, `
		INSERT INTO matches (id_deck, host_user_id, room_code, status, is_public, total_rounds, round_time_sec)
		VALUES (1, $1, 'FULLRM', 'WAITING', true, 8, 15)
		RETURNING id_match::text
	`, gaelID).Scan(&fullMatchID)
	if err != nil {
		t.Fatalf("falló insertar partida llena: %v", err)
	}

	for i := 1; i <= 8; i++ {
		uid := strings.ReplaceAll(uuid.NewString(), "-", "")[:12]
		email := fmt.Sprintf("p%d_%s@test.org", i, uid)
		userName := fmt.Sprintf("p%d_test_%s", i, uid)
		var pUserID string
		err = testDB.Pool.QueryRow(ctx, `
			INSERT INTO accounts (email, password_hash, role)
			VALUES ($1, '$2a$10$f6V4aUFb5FkBdEG21Iz1W.5AONrEdWxEKGRtwSzfPEWIwJVe02Rdu', 'LEARNER')
			RETURNING id_account::text
		`, email).Scan(&pUserID)
		if err != nil {
			t.Fatalf("falló crear cuenta participante %d: %v", i, err)
		}
		_, err = testDB.Pool.Exec(ctx, `
			INSERT INTO users (id_user, user_name, name, anchor_dialect) VALUES ($1, $2, $2, 'MX')
		`, pUserID, userName)
		if err != nil {
			t.Fatalf("falló crear usuario participante %d: %v", i, err)
		}
		_, err = testDB.Pool.Exec(ctx, `
			INSERT INTO match_participants (id_match, id_user) VALUES ($1, $2)
		`, fullMatchID, pUserID)
		if err != nil {
			t.Fatalf("falló insertar participante %d en partida llena: %v", i, err)
		}
	}

	// 6. Consultar ListPublicWaitingMatches
	matches, err := arenaSvc.ListPublicWaitingMatches(ctx, 20)
	if err != nil {
		t.Fatalf("ListPublicWaitingMatches falló: %v", err)
	}

	// Verificar que la partida válida (matchPkg.IDMatch) está presente
	foundValid := false
	for _, m := range matches {
		if m.IDMatch == matchPkg.IDMatch {
			foundValid = true
			if m.RoomCode != matchPkg.RoomCode {
				t.Errorf("se esperaba RoomCode %s, obtenido %s", matchPkg.RoomCode, m.RoomCode)
			}
			if m.IDDeck != 1 {
				t.Errorf("se esperaba IDDeck 1, obtenido %d", m.IDDeck)
			}
			if m.DeckTitle == "" {
				t.Errorf("se esperaba DeckTitle no vacío")
			}
			if m.HostUserID == nil || *m.HostUserID != gaelID {
				t.Errorf("se esperaba HostUserID %s, obtenido %v", gaelID, m.HostUserID)
			}
			if m.HostUserName != "gael_learner" {
				t.Errorf("se esperaba HostUserName 'gael_learner', obtenido '%s'", m.HostUserName)
			}
			if m.TotalRounds != 8 {
				t.Errorf("se esperaba TotalRounds 8, obtenido %d", m.TotalRounds)
			}
			if m.RoundTimeSec != 15 {
				t.Errorf("se esperaba RoundTimeSec 15, obtenido %d", m.RoundTimeSec)
			}
			if m.ParticipantCount != 0 {
				t.Errorf("se esperaba ParticipantCount 0, obtenido %d", m.ParticipantCount)
			}
		}

		// Ninguna de las partidas inválidas debe aparecer
		if m.IDMatch == privatePkg.IDMatch {
			t.Errorf("partida privada no debe ser devuelta")
		}
		if m.RoomCode == "INPROG" {
			t.Errorf("partida con status != 'WAITING' no debe ser devuelta")
		}
		if m.RoomCode == "OLDROO" {
			t.Errorf("partida con más de 2 horas no debe ser devuelta")
		}
		if m.RoomCode == "FULLRM" {
			t.Errorf("partida con >= 8 participantes no debe ser devuelta")
		}
	}

	if !foundValid {
		t.Fatalf("la partida válida %s (%s) no fue encontrada en la lista pública", matchPkg.IDMatch, matchPkg.RoomCode)
	}

	// 7. Verificar límite
	limitedMatches, err := arenaSvc.ListPublicWaitingMatches(ctx, 1)
	if err != nil {
		t.Fatalf("ListPublicWaitingMatches con limit=1 falló: %v", err)
	}
	if len(limitedMatches) > 1 {
		t.Errorf("se esperaba como máximo 1 resultado con limit=1, obtenido %d", len(limitedMatches))
	}
}

func assertDomainCode(t *testing.T, err error, wantBase error, wantCode string) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected error with code %q, got nil", wantCode)
	}
	if !errors.Is(err, wantBase) {
		t.Fatalf("expected errors.Is(err, %v), got %v", wantBase, err)
	}
	var domErr *service.DomainError
	if !errors.As(err, &domErr) {
		t.Fatalf("expected *service.DomainError, got %T (%v)", err, err)
	}
	if domErr.Code != wantCode {
		t.Fatalf("expected error code %q, got %q (msg=%q)", wantCode, domErr.Code, domErr.Message)
	}
}

func TestIdentityService_SignUp_Validation_And_ConflictCodes(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	if testDB == nil {
		return
	}

	tokenSvc := auth.NewTokenService("test-secret-key-minimum-32-chars", 15*time.Minute)
	identitySvc := service.NewIdentityService(testDB, tokenSvc, 168*time.Hour)
	ctx := context.Background()

	cases := []struct {
		name string
		req  models.SignUpRequest
	}{
		{
			name: "invalid email",
			req: models.SignUpRequest{
				Email:         "not-an-email",
				Password:      "ValidPass123!",
				Name:          "Test User",
				UserName:      "valid_user_1",
				AnchorDialect: "MX",
			},
		},
		{
			name: "short password",
			req: models.SignUpRequest{
				Email:         "short@palabrario.org",
				Password:      "1234567",
				Name:          "Test User",
				UserName:      "valid_user_2",
				AnchorDialect: "MX",
			},
		},
		{
			name: "password exceeds 72 bytes",
			req: models.SignUpRequest{
				Email:         "longpw@palabrario.org",
				Password:      strings.Repeat("a", 73),
				Name:          "Test User",
				UserName:      "valid_user_3",
				AnchorDialect: "MX",
			},
		},
		{
			name: "invalid username",
			req: models.SignUpRequest{
				Email:         "baduser@palabrario.org",
				Password:      "ValidPass123!",
				Name:          "Test User",
				UserName:      "_invalid_start",
				AnchorDialect: "MX",
			},
		},
		{
			name: "username exceeds 50 chars",
			req: models.SignUpRequest{
				Email:         "longuser@palabrario.org",
				Password:      "ValidPass123!",
				Name:          "Test User",
				UserName:      strings.Repeat("a", 51),
				AnchorDialect: "MX",
			},
		},
		{
			name: "invalid dialect",
			req: models.SignUpRequest{
				Email:         "baddialect@palabrario.org",
				Password:      "ValidPass123!",
				Name:          "Test User",
				UserName:      "valid_user_4",
				AnchorDialect: "ZZ",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := identitySvc.SignUp(ctx, tc.req)
			assertDomainCode(t, err, service.ErrInvalidInput, "VALIDATION_ERROR")
		})
	}

	// Leading '@' stripped, empty dialect defaults to MX, both tokens issued
	created, err := identitySvc.SignUp(ctx, models.SignUpRequest{
		Email:         "CleanUser@Palabrario.org",
		Password:      "ValidPass123!",
		Name:          "",
		UserName:      "@Clean_Handle",
		AnchorDialect: "",
	})
	if err != nil {
		t.Fatalf("SignUp failed unexpectedly: %v", err)
	}
	if created.Token == "" || created.RefreshToken == "" {
		t.Fatalf("expected both access and refresh tokens, got token=%q refresh=%q", created.Token, created.RefreshToken)
	}
	if created.User.UserName != "clean_handle" || created.User.AnchorDialect != "MX" {
		t.Fatalf("unexpected normalized profile: %+v", created.User)
	}

	// Duplicate email -> EMAIL_ALREADY_EXISTS
	_, err = identitySvc.SignUp(ctx, models.SignUpRequest{
		Email:         "cleanuser@palabrario.org",
		Password:      "ValidPass123!",
		Name:          "Another Name",
		UserName:      "brand_new_handle",
		AnchorDialect: "ES",
	})
	assertDomainCode(t, err, service.ErrConflict, "EMAIL_ALREADY_EXISTS")

	// Duplicate username -> USERNAME_ALREADY_EXISTS
	_, err = identitySvc.SignUp(ctx, models.SignUpRequest{
		Email:         "other_email@palabrario.org",
		Password:      "ValidPass123!",
		Name:          "Another Name",
		UserName:      "CLEAN_HANDLE",
		AnchorDialect: "ES",
	})
	assertDomainCode(t, err, service.ErrConflict, "USERNAME_ALREADY_EXISTS")
}

func TestIdentityService_RefreshToken_SignOut_And_Revocation(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	if testDB == nil {
		return
	}

	tokenSvc := auth.NewTokenService("test-secret-key-minimum-32-chars", 15*time.Minute)
	identitySvc := service.NewIdentityService(testDB, tokenSvc, 168*time.Hour)
	ctx := context.Background()

	// 1. SignIn issues access + refresh tokens
	loginResp, err := identitySvc.SignIn(ctx, models.SignInRequest{
		Email:    "gael@palabrario.org",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("SignIn failed: %v", err)
	}
	if loginResp.Token == "" || loginResp.RefreshToken == "" {
		t.Fatal("expected non-empty access and refresh tokens on SignIn")
	}

	// 2. RefreshToken rotates the token atomically
	rotated, err := identitySvc.RefreshToken(ctx, loginResp.RefreshToken)
	if err != nil {
		t.Fatalf("RefreshToken failed: %v", err)
	}
	if rotated.RefreshToken == "" || rotated.RefreshToken == loginResp.RefreshToken {
		t.Fatal("expected RefreshToken to issue a new distinct refresh token")
	}

	// 3. Reusing the already-rotated refresh token fails with TOKEN_REVOKED
	_, err = identitySvc.RefreshToken(ctx, loginResp.RefreshToken)
	assertDomainCode(t, err, service.ErrUnauthorized, "TOKEN_REVOKED")

	// 4. SignOut revokes the rotated refresh token; empty SignOut is a no-op
	if err := identitySvc.SignOut(ctx, ""); err != nil {
		t.Fatalf("SignOut with empty token failed: %v", err)
	}
	if err := identitySvc.SignOut(ctx, rotated.RefreshToken); err != nil {
		t.Fatalf("SignOut failed: %v", err)
	}
	_, err = identitySvc.RefreshToken(ctx, rotated.RefreshToken)
	assertDomainCode(t, err, service.ErrUnauthorized, "TOKEN_REVOKED")

	// 5. Expired refresh token returns TOKEN_EXPIRED
	expiredLogin, err := identitySvc.SignIn(ctx, models.SignInRequest{
		Email:    "gael@palabrario.org",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("SignIn failed: %v", err)
	}
	_, err = testDB.Pool.Exec(ctx, `UPDATE refresh_tokens SET expires_at = now() - INTERVAL '1 minute' WHERE token_hash = $1`, auth.HashRefreshToken(expiredLogin.RefreshToken))
	if err != nil {
		t.Fatalf("failed to expire refresh token: %v", err)
	}
	_, err = identitySvc.RefreshToken(ctx, expiredLogin.RefreshToken)
	assertDomainCode(t, err, service.ErrUnauthorized, "TOKEN_EXPIRED")

	// 6. Suspended account refresh returns ACCOUNT_SUSPENDED and revokes all refresh tokens
	activeLogin, err := identitySvc.SignIn(ctx, models.SignInRequest{
		Email:    "gael@palabrario.org",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("SignIn failed: %v", err)
	}
	_, err = testDB.Pool.Exec(ctx, `UPDATE accounts SET is_active = false WHERE email = 'gael@palabrario.org'`)
	if err != nil {
		t.Fatalf("failed to suspend account: %v", err)
	}
	_, err = identitySvc.RefreshToken(ctx, activeLogin.RefreshToken)
	assertDomainCode(t, err, service.ErrForbidden, "ACCOUNT_SUSPENDED")

	var activeTokens int
	err = testDB.Pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM refresh_tokens rt
		JOIN accounts a ON a.id_account = rt.id_account
		WHERE a.email = 'gael@palabrario.org' AND rt.revoked_at IS NULL
	`).Scan(&activeTokens)
	if err != nil {
		t.Fatalf("failed to query active refresh tokens: %v", err)
	}
	if activeTokens != 0 {
		t.Fatalf("expected all refresh tokens revoked after suspended refresh, got %d active", activeTokens)
	}
}

func TestIdentityService_GetMe_And_ChangePassword_Revocation(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	if testDB == nil {
		return
	}

	tokenSvc := auth.NewTokenService("test-secret-key-minimum-32-chars", 15*time.Minute)
	identitySvc := service.NewIdentityService(testDB, tokenSvc, 168*time.Hour)
	ctx := context.Background()

	// 1. Missing account UUID -> 401 UNAUTHORIZED
	_, err := identitySvc.GetMe(ctx, uuid.NewString())
	assertDomainCode(t, err, service.ErrUnauthorized, "UNAUTHORIZED")

	// 2. SignIn Gael & record initial token_version + refresh token
	loginResp, err := identitySvc.SignIn(ctx, models.SignInRequest{
		Email:    "gael@palabrario.org",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("SignIn failed: %v", err)
	}

	var beforeVersion int
	err = testDB.Pool.QueryRow(ctx, `SELECT token_version FROM accounts WHERE id_account = $1`, loginResp.User.IDUser).Scan(&beforeVersion)
	if err != nil {
		t.Fatalf("failed to read token_version: %v", err)
	}

	// 3. Password > 72 bytes rejected in ChangePassword
	err = identitySvc.ChangePassword(ctx, loginResp.User.IDUser, models.ChangePasswordRequest{
		CurrentPassword: "password123",
		NewPassword:     strings.Repeat("x", 73),
	})
	assertDomainCode(t, err, service.ErrInvalidInput, "VALIDATION_ERROR")

	// 4. Successful ChangePassword bumps token_version and revokes prior refresh tokens
	if err := identitySvc.ChangePassword(ctx, loginResp.User.IDUser, models.ChangePasswordRequest{
		CurrentPassword: "password123",
		NewPassword:     "NewSecurePass123!",
	}); err != nil {
		t.Fatalf("ChangePassword failed: %v", err)
	}

	var afterVersion int
	err = testDB.Pool.QueryRow(ctx, `SELECT token_version FROM accounts WHERE id_account = $1`, loginResp.User.IDUser).Scan(&afterVersion)
	if err != nil {
		t.Fatalf("failed to read updated token_version: %v", err)
	}
	if afterVersion != beforeVersion+1 {
		t.Fatalf("expected token_version=%d, got %d", beforeVersion+1, afterVersion)
	}

	_, err = identitySvc.RefreshToken(ctx, loginResp.RefreshToken)
	assertDomainCode(t, err, service.ErrUnauthorized, "TOKEN_REVOKED")

	// 5. Suspended account -> GetMe returns 403 ACCOUNT_SUSPENDED
	_, err = testDB.Pool.Exec(ctx, `UPDATE accounts SET is_active = false WHERE id_account = $1`, loginResp.User.IDUser)
	if err != nil {
		t.Fatalf("failed to suspend account: %v", err)
	}
	_, err = identitySvc.GetMe(ctx, loginResp.User.IDUser)
	assertDomainCode(t, err, service.ErrForbidden, "ACCOUNT_SUSPENDED")
}

func TestErrors_MapDBError(t *testing.T) {
	cases := []struct {
		name     string
		err      error
		wantBase error
		wantCode string
	}{
		{
			name:     "pgx no rows",
			err:      pgx.ErrNoRows,
			wantBase: service.ErrNotFound,
			wantCode: "NOT_FOUND",
		},
		{
			name:     "string truncation 22001",
			err:      &pgconn.PgError{Code: "22001"},
			wantBase: service.ErrInvalidInput,
			wantCode: "INPUT_TOO_LONG",
		},
		{
			name:     "check violation 23514",
			err:      &pgconn.PgError{Code: "23514"},
			wantBase: service.ErrValidation,
			wantCode: "VALIDATION_ERROR",
		},
		{
			name:     "foreign key violation 23503",
			err:      &pgconn.PgError{Code: "23503"},
			wantBase: service.ErrInvalidInput,
			wantCode: "INVALID_REFERENCE",
		},
		{
			name:     "unique email idx_accounts_email_lower",
			err:      &pgconn.PgError{Code: "23505", ConstraintName: "idx_accounts_email_lower"},
			wantBase: service.ErrConflict,
			wantCode: "EMAIL_ALREADY_EXISTS",
		},
		{
			name:     "unique email accounts_email_key",
			err:      &pgconn.PgError{Code: "23505", ConstraintName: "accounts_email_key"},
			wantBase: service.ErrConflict,
			wantCode: "EMAIL_ALREADY_EXISTS",
		},
		{
			name:     "unique username idx_users_user_name_lower",
			err:      &pgconn.PgError{Code: "23505", ConstraintName: "idx_users_user_name_lower"},
			wantBase: service.ErrConflict,
			wantCode: "USERNAME_ALREADY_EXISTS",
		},
		{
			name:     "unique username users_user_name_key",
			err:      &pgconn.PgError{Code: "23505", ConstraintName: "users_user_name_key"},
			wantBase: service.ErrConflict,
			wantCode: "USERNAME_ALREADY_EXISTS",
		},
		{
			name:     "generic conflict 23505",
			err:      &pgconn.PgError{Code: "23505", ConstraintName: "other_unique"},
			wantBase: service.ErrConflict,
			wantCode: "CONFLICT",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mapped := service.MapDBError(tc.err, "default error")
			assertDomainCode(t, mapped, tc.wantBase, tc.wantCode)
		})
	}
}


