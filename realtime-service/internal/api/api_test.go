package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"palabrario/pkg/auth"
	"palabrario/realtime-service/internal/api/handlers"
	"palabrario/realtime-service/internal/config"
	"palabrario/realtime-service/internal/service"
)

func TestAPIRoutes(t *testing.T) {
	cfg := &config.Config{
		Port:               "8001",
		Environment:        "test",
		JWTSecret:          "super_secret_jwt_key_palabrario_change_in_production_min_32_chars",
		JWTExpirationHours: 1 * time.Hour,
		CoreServiceURL:     "http://localhost:8000",
	}

	tokenService := auth.NewTokenService(cfg.JWTSecret, cfg.JWTExpirationHours)
	coreClient := service.NewCoreClient(cfg.CoreServiceURL, nil, tokenService)
	hub := service.NewArenaHub(coreClient, tokenService)

	router := NewRouter(cfg, hub, tokenService)
	handler := router.Setup()

	t.Run("GET /health returns 200 and valid JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("esperado código 200, obtenido %d", rec.Code)
		}

		var resp handlers.HealthResponse
		if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
			t.Fatalf("error decodificando HealthResponse: %v", err)
		}

		if resp.Status != "ok" || resp.Service != "realtime-service" || resp.ActiveRooms != 0 {
			t.Errorf("respuesta de salud inesperada: %+v", resp)
		}
	})

	t.Run("GET / returns 200 and welcome message", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("esperado código 200, obtenido %d", rec.Code)
		}

		var resp map[string]string
		if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
			t.Fatalf("error decodificando respuesta raíz: %v", err)
		}

		if resp["message"] != "Bienvenido a Palabrario Realtime Service (Golang)" {
			t.Errorf("mensaje de bienvenida inesperado: %s", resp["message"])
		}
	})

	t.Run("GET /ws/arena/{roomCode} without token returns 401 Unauthorized", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/ws/arena/MX1234", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("esperado código 401 sin token, obtenido %d", rec.Code)
		}
	})

	t.Run("GET /ws/arena/{roomCode} with invalid token returns 401 Unauthorized", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/ws/arena/MX1234?token=invalid.token.here", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("esperado código 401 con token inválido, obtenido %d", rec.Code)
		}
	})

	t.Run("GET /ws/arena/{roomCode} with invalid Bearer header returns 401 Unauthorized", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/ws/arena/MX1234", nil)
		req.Header.Set("Authorization", "Bearer invalid-bearer-token")
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("esperado código 401 con Bearer inválido, obtenido %d", rec.Code)
		}
	})
}
