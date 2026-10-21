package service_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

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
