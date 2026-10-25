package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"palabrario/pkg/auth"
	"palabrario/pkg/models"
)

func TestAuthMiddleware(t *testing.T) {
	jwtSecret := "super-secret-key-palabrario-middleware-test-2026"
	tokenService := auth.NewTokenService(jwtSecret, time.Hour)
	mw := NewAuthMiddleware(tokenService)

	nextCalled := false
	dummyHandler := mw.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
		if !ok || claims == nil {
			t.Errorf("se esperaba que los claims estuvieran en el contexto")
		}
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	}))

	t.Run("missing authorization header returns 401", func(t *testing.T) {
		nextCalled = false
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		dummyHandler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("esperaba 401, obtuvo %d", rec.Code)
		}
		if nextCalled {
			t.Errorf("el handler siguiente no debió invocarse")
		}
	})

	t.Run("invalid token format returns 401", func(t *testing.T) {
		nextCalled = false
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Basic user:pass")
		rec := httptest.NewRecorder()

		dummyHandler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("esperaba 401, obtuvo %d", rec.Code)
		}
	})

	t.Run("valid token injects claims and calls next", func(t *testing.T) {
		nextCalled = false
		token, err := tokenService.GenerateToken(&models.UserProfile{
			IDUser:        "user-uuid-xyz",
			Email:         "test@palabrario.org",
			UserName:      "tester",
			Role:          models.RoleLearner,
			AnchorDialect: "MX",
		})
		if err != nil {
			t.Fatalf("error generando token: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()

		dummyHandler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("esperaba 200 OK, obtuvo %d", rec.Code)
		}
		if !nextCalled {
			t.Errorf("el handler siguiente debió invocarse")
		}
	})
}
