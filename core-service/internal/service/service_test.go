package service_test

import (
	"context"
	"errors"
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
