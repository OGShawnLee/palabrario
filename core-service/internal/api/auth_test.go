package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"palabrario/core-service/internal/api"
	"palabrario/core-service/internal/config"
	"palabrario/core-service/internal/database"
	"palabrario/core-service/internal/testutil"
	"palabrario/pkg/models"
)

// ==============================================================================
// CAPA 20%: PRUEBAS DE HUMO HTTP END-TO-END (ROUTER + MIDDLEWARES JWT + JSON)
// ==============================================================================

func setupTestServer(t *testing.T) (*httptest.Server, *database.DB) {
	t.Helper()

	db := testutil.SetupTestDB(t)
	if db == nil {
		return nil, nil
	}

	cfg := &config.Config{
		Port:               "8000",
		Environment:        "test",
		JWTSecret:          "test-super-secret-jwt-key-min-32-chars-long",
		JWTExpirationHours: 72 * time.Hour,
	}

	router := api.NewRouter(cfg, db)
	server := httptest.NewServer(router.Setup())

	t.Cleanup(func() {
		server.Close()
	})

	return server, db
}

func TestAuth_SignIn_Success(t *testing.T) {
	server, _ := setupTestServer(t)
	if server == nil {
		return
	}

	payload := models.SignInRequest{
		Email:    "silas@palabrario.org",
		Password: "password123",
	}
	body, _ := json.Marshal(payload)

	resp, err := http.Post(server.URL+"/api/v1/auth/sign-in", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Fallo en la petición HTTP: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Se esperaba estado 200 OK, pero se obtuvo %d", resp.StatusCode)
	}

	var authResp models.AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		t.Fatalf("Fallo al decodificar respuesta JSON: %v", err)
	}

	if authResp.Token == "" {
		t.Error("El token de sesión JWT devuelto está vacío")
	}
	if authResp.User.Email != "silas@palabrario.org" {
		t.Errorf("Se esperaba email 'silas@palabrario.org', obtenido '%s'", authResp.User.Email)
	}
	if authResp.User.Role != models.RoleEditor {
		t.Errorf("Se esperaba rol 'EDITOR', obtenido '%s'", authResp.User.Role)
	}
}

func TestAuth_SignIn_InvalidPassword(t *testing.T) {
	server, _ := setupTestServer(t)
	if server == nil {
		return
	}

	payload := models.SignInRequest{
		Email:    "silas@palabrario.org",
		Password: "wrongpassword!",
	}
	body, _ := json.Marshal(payload)

	resp, err := http.Post(server.URL+"/api/v1/auth/sign-in", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Fallo en la petición HTTP: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Se esperaba 401 Unauthorized para contraseña errónea, obtenido %d", resp.StatusCode)
	}
}

func TestAuth_GetMe_ProtectedEndpoint(t *testing.T) {
	server, _ := setupTestServer(t)
	if server == nil {
		return
	}

	// 1. Obtener token válido
	payload := models.SignInRequest{
		Email:    "admin@palabrario.org",
		Password: "password123",
	}
	body, _ := json.Marshal(payload)
	loginResp, err := http.Post(server.URL+"/api/v1/auth/sign-in", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Fallo al iniciar sesión: %v", err)
	}
	defer loginResp.Body.Close()

	var authResp models.AuthResponse
	_ = json.NewDecoder(loginResp.Body).Decode(&authResp)

	// 2. Consulta autorizada con token
	req, _ := http.NewRequest(http.MethodGet, server.URL+"/api/v1/auth/me", nil)
	req.Header.Set("Authorization", "Bearer "+authResp.Token)

	meResp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Fallo en petición protegida: %v", err)
	}
	defer meResp.Body.Close()

	if meResp.StatusCode != http.StatusOK {
		t.Errorf("Se esperaba 200 OK con token válido, obtenido %d", meResp.StatusCode)
	}

	var profile models.UserProfile
	_ = json.NewDecoder(meResp.Body).Decode(&profile)
	if profile.Role != models.RoleAdmin {
		t.Errorf("Se esperaba rol ADMIN, obtenido %s", profile.Role)
	}

	// 3. Consulta sin token (debe fallar 401)
	unauthResp, err := http.Get(server.URL + "/api/v1/auth/me")
	if err != nil {
		t.Fatalf("Fallo en petición no autorizada: %v", err)
	}
	defer unauthResp.Body.Close()

	if unauthResp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Se esperaba 401 Unauthorized sin token, obtenido %d", unauthResp.StatusCode)
	}
}

func TestAuth_SignUp_Lifecycle(t *testing.T) {
	server, _ := setupTestServer(t)
	if server == nil {
		return
	}

	// Gracias al clonado efímero por TEMPLATE, no requerimos sufijos aleatorios ni DELETE manual:
	// cada test recibe una base de datos virgen y aislada que se destruye automáticamente al terminar.
	signUpReq := models.SignUpRequest{
		Email:         "nuevo_aprendiz@palabrario.org",
		Password:      "SecureTestPass123!",
		Name:          "Usuario Automatizado",
		UserName:      "autouser_co",
		AnchorDialect: "CO",
	}
	body, _ := json.Marshal(signUpReq)

	resp, err := http.Post(server.URL+"/api/v1/auth/sign-up", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Fallo en sign-up: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Se esperaba 201 Created en sign-up, obtenido %d", resp.StatusCode)
	}

	var authResp models.AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		t.Fatalf("Fallo decodificando respuesta de registro: %v", err)
	}

	if authResp.User.AnchorDialect != "CO" {
		t.Errorf("Se esperaba dialecto 'CO', obtenido '%s'", authResp.User.AnchorDialect)
	}
	if authResp.User.Role != models.RoleLearner {
		t.Errorf("Se esperaba rol LEARNER por defecto, obtenido '%s'", authResp.User.Role)
	}
}
