package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"palabrario/pkg/auth"
	"palabrario/pkg/models"
	"palabrario/pkg/response"
)

func TestAuthMiddleware(t *testing.T) {
	jwtSecret := "test-jwt-secret-key-palabrario-minimum-32-chars"
	tokenService := auth.NewTokenService(jwtSecret, time.Hour)
	mw := NewAuthMiddleware(tokenService)

	nextCalled := false
	dummyHandler := mw.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
		if !ok || claims == nil {
			t.Errorf("expected claims in request context")
		}
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	}))

	decodeError := func(t *testing.T, rec *httptest.ResponseRecorder) response.ErrorBody {
		t.Helper()
		var body response.ErrorBody
		if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode ErrorBody: %v", err)
		}
		return body
	}

	t.Run("missing authorization header returns 401", func(t *testing.T) {
		nextCalled = false
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		dummyHandler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", rec.Code)
		}
		if nextCalled {
			t.Fatal("next handler must not be called")
		}
		errBody := decodeError(t, rec)
		if errBody.Code != "MISSING_TOKEN" {
			t.Errorf("expected code MISSING_TOKEN, got %q", errBody.Code)
		}
	})

	t.Run("invalid token format returns 401", func(t *testing.T) {
		nextCalled = false
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Basic user:pass")
		rec := httptest.NewRecorder()

		dummyHandler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", rec.Code)
		}
		if nextCalled {
			t.Fatal("next handler must not be called")
		}
		errBody := decodeError(t, rec)
		if errBody.Code != "INVALID_TOKEN_FORMAT" {
			t.Errorf("expected code INVALID_TOKEN_FORMAT, got %q", errBody.Code)
		}
	})

	t.Run("empty Bearer token returns 401", func(t *testing.T) {
		nextCalled = false
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer   ")
		rec := httptest.NewRecorder()

		dummyHandler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", rec.Code)
		}
		if nextCalled {
			t.Fatal("next handler must not be called")
		}
		errBody := decodeError(t, rec)
		if errBody.Code != "INVALID_TOKEN_FORMAT" {
			t.Errorf("expected code INVALID_TOKEN_FORMAT, got %q", errBody.Code)
		}
	})

	t.Run("expired JWT returns 401 without leaking internal error", func(t *testing.T) {
		nextCalled = false
		past := time.Now().UTC().Add(-5 * time.Minute)
		expiredClaims := auth.Claims{
			UserID:        "user-uuid-xyz",
			Email:         "test@palabrario.org",
			UserName:      "tester",
			Role:          models.RoleLearner,
			AnchorDialect: "MX",
			TokenVersion:  1,
			TokenType:     auth.TokenTypeAccess,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(past.Add(time.Minute)),
				IssuedAt:  jwt.NewNumericDate(past),
				NotBefore: jwt.NewNumericDate(past),
				Issuer:    "palabrario-core-service",
			},
		}
		expiredToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims).SignedString([]byte(jwtSecret))
		if err != nil {
			t.Fatalf("failed to sign expired token: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer "+expiredToken)
		rec := httptest.NewRecorder()

		dummyHandler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", rec.Code)
		}
		if nextCalled {
			t.Fatal("next handler must not be called")
		}
		rawBody := rec.Body.String()
		if strings.Contains(rawBody, "token has invalid claims") || strings.Contains(rawBody, "token is expired") {
			t.Fatalf("401 body leaked internal jwt error: %s", rawBody)
		}
		errBody := decodeError(t, rec)
		if errBody.Code != "INVALID_TOKEN" || errBody.Error != "Invalid or expired token" {
			t.Errorf("unexpected error body: %+v", errBody)
		}
	})

	t.Run("S2S internal token rejected on user endpoint with 401", func(t *testing.T) {
		nextCalled = false
		s2sTok, err := tokenService.GenerateInternalToken("realtime-service", "MX1234")
		if err != nil {
			t.Fatalf("failed to generate S2S token: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer "+s2sTok)
		rec := httptest.NewRecorder()

		dummyHandler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", rec.Code)
		}
		if nextCalled {
			t.Fatal("next handler must not be called")
		}
		errBody := decodeError(t, rec)
		if errBody.Code != "INVALID_TOKEN" {
			t.Errorf("expected code INVALID_TOKEN, got %q", errBody.Code)
		}
	})

	t.Run("lowercase bearer scheme accepted and calls next", func(t *testing.T) {
		nextCalled = false
		token, err := tokenService.GenerateToken(&models.UserProfile{
			IDUser:        "user-uuid-xyz",
			Email:         "test@palabrario.org",
			UserName:      "tester",
			Role:          models.RoleLearner,
			AnchorDialect: "MX",
		})
		if err != nil {
			t.Fatalf("error generating token: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "bearer "+token)
		rec := httptest.NewRecorder()

		dummyHandler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200 OK, got %d", rec.Code)
		}
		if !nextCalled {
			t.Fatal("next handler should have been called")
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
			t.Fatalf("error generating token: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()

		dummyHandler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200 OK, got %d", rec.Code)
		}
		if !nextCalled {
			t.Fatal("next handler should have been called")
		}
	})
}
