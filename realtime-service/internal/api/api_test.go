package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"

	"palabrario/pkg/auth"
	"palabrario/pkg/models"
	"palabrario/pkg/response"
	"palabrario/realtime-service/internal/api/handlers"
	"palabrario/realtime-service/internal/config"
	"palabrario/realtime-service/internal/service"
)

func TestAPIRoutes(t *testing.T) {
	jwtSecret := "test-jwt-secret-key-palabrario-minimum-32-chars"
	cfg := &config.Config{
		Port:               "8001",
		Environment:        "test",
		JWTSecret:          jwtSecret,
		AccessTokenTTL:     15 * time.Minute,
		JWTExpirationHours: 1 * time.Hour,
		CoreServiceURL:     "http://localhost:8000",
		AllowedOrigins:     auth.ParseCORSAllowedOrigins(""),
	}

	tokenService := auth.NewTokenService(cfg.JWTSecret, cfg.AccessTokenTTL)
	coreClient := service.NewCoreClient(cfg.CoreServiceURL, nil, tokenService)
	hub := service.NewArenaHub(coreClient, tokenService)

	router := NewRouter(cfg, hub, tokenService)
	handler := router.Setup()

	decodeError := func(t *testing.T, rec *httptest.ResponseRecorder) response.ErrorBody {
		t.Helper()
		var body response.ErrorBody
		if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode ErrorBody: %v", err)
		}
		return body
	}

	t.Run("GET /health returns 200 and valid JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", rec.Code)
		}

		var resp handlers.HealthResponse
		if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
			t.Fatalf("error decoding HealthResponse: %v", err)
		}

		if resp.Status != "ok" || resp.Service != "realtime-service" || resp.ActiveRooms != 0 {
			t.Errorf("unexpected health response: %+v", resp)
		}
	})

	t.Run("GET / returns 200 and welcome message", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", rec.Code)
		}

		var resp map[string]string
		if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
			t.Fatalf("error decoding root response: %v", err)
		}

		if resp["message"] != "Bienvenido a Palabrario Realtime Service (Golang)" {
			t.Errorf("unexpected welcome message: %s", resp["message"])
		}
	})

	t.Run("GET /ws/arena/{roomCode} without token returns 401 Unauthorized", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/ws/arena/MX1234", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401 without token, got %d", rec.Code)
		}
		if errBody := decodeError(t, rec); errBody.Code != "MISSING_TOKEN" {
			t.Errorf("expected MISSING_TOKEN, got %q", errBody.Code)
		}
	})

	t.Run("GET /ws/arena/{roomCode} with invalid query token returns 401 Unauthorized", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/ws/arena/MX1234?token=invalid.token.here", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401 with invalid query token, got %d", rec.Code)
		}
		if errBody := decodeError(t, rec); errBody.Code != "INVALID_TOKEN" {
			t.Errorf("expected INVALID_TOKEN, got %q", errBody.Code)
		}
	})

	t.Run("GET /ws/arena/{roomCode} with invalid Bearer header returns 401 Unauthorized", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/ws/arena/MX1234", nil)
		req.Header.Set("Authorization", "Bearer invalid-bearer-token")
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401 with invalid Bearer, got %d", rec.Code)
		}
		if errBody := decodeError(t, rec); errBody.Code != "INVALID_TOKEN" {
			t.Errorf("expected INVALID_TOKEN, got %q", errBody.Code)
		}
	})

	t.Run("GET /ws/arena/{roomCode} with empty Bearer header returns 401 even if query token is set", func(t *testing.T) {
		validTok, err := tokenService.GenerateToken(&models.UserProfile{
			IDUser:        "u-1",
			Email:         "u1@palabrario.org",
			UserName:      "u1",
			Role:          models.RoleLearner,
			AnchorDialect: "MX",
		})
		if err != nil {
			t.Fatalf("GenerateToken failed: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, "/ws/arena/MX1234?token="+validTok, nil)
		req.Header.Set("Authorization", "Bearer   ")
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401 with empty Bearer header, got %d", rec.Code)
		}
		if errBody := decodeError(t, rec); errBody.Code != "INVALID_TOKEN_FORMAT" {
			t.Errorf("expected INVALID_TOKEN_FORMAT, got %q", errBody.Code)
		}
	})

	t.Run("GET /ws/arena/{roomCode} with expired JWT returns 401 Unauthorized", func(t *testing.T) {
		past := time.Now().UTC().Add(-5 * time.Minute)
		expiredClaims := auth.Claims{
			UserID:        "u-1",
			Email:         "u1@palabrario.org",
			UserName:      "u1",
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
		expiredTok, err := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims).SignedString([]byte(jwtSecret))
		if err != nil {
			t.Fatalf("failed to sign expired token: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, "/ws/arena/MX1234", nil)
		req.Header.Set("Authorization", "Bearer "+expiredTok)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401 with expired JWT, got %d", rec.Code)
		}
		if errBody := decodeError(t, rec); errBody.Code != "INVALID_TOKEN" {
			t.Errorf("expected INVALID_TOKEN, got %q", errBody.Code)
		}
	})

	t.Run("GET /ws/arena/{roomCode} rejects S2S internal token with 401 Unauthorized", func(t *testing.T) {
		s2sTok, err := tokenService.GenerateInternalToken("core-service", "MX1234")
		if err != nil {
			t.Fatalf("GenerateInternalToken failed: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, "/ws/arena/MX1234", nil)
		req.Header.Set("Authorization", "Bearer "+s2sTok)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401 with S2S token on user endpoint, got %d", rec.Code)
		}
	})

	t.Run("CORS sets AllowCredentials=false and restricts origins", func(t *testing.T) {
		reqAllowed := httptest.NewRequest(http.MethodOptions, "/health", nil)
		reqAllowed.Header.Set("Origin", "http://localhost:5173")
		reqAllowed.Header.Set("Access-Control-Request-Method", "GET")
		recAllowed := httptest.NewRecorder()
		handler.ServeHTTP(recAllowed, reqAllowed)

		if got := recAllowed.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:5173" {
			t.Errorf("expected allowed origin echoed, got %q", got)
		}
		if creds := recAllowed.Header().Get("Access-Control-Allow-Credentials"); creds == "true" {
			t.Errorf("AllowCredentials must be false, got %q", creds)
		}

		reqDenied := httptest.NewRequest(http.MethodOptions, "/health", nil)
		reqDenied.Header.Set("Origin", "https://evil.example.com")
		reqDenied.Header.Set("Access-Control-Request-Method", "GET")
		recDenied := httptest.NewRecorder()
		handler.ServeHTTP(recDenied, reqDenied)

		if got := recDenied.Header().Get("Access-Control-Allow-Origin"); got != "" {
			t.Errorf("expected disallowed origin not to be echoed, got %q", got)
		}
	})
}

func TestWebSocketOriginAndLowercaseBearer(t *testing.T) {
	jwtSecret := "test-jwt-secret-key-palabrario-minimum-32-chars"
	tokenService := auth.NewTokenService(jwtSecret, 15*time.Minute)

	mockMatchPkg := &models.MatchPackage{
		IDMatch:      "match-api-1",
		RoomCode:     "MX4321",
		HostUserID:   "user-1",
		IDDeck:       1,
		DeckTitle:    "Origin Check Deck",
		TotalRounds:  1,
		RoundTimeSec: 10,
		Questions: []models.InternalArenaQuestion{
			{
				RoundNumber:   1,
				IDSense:       101,
				Prompt:        "Q1",
				Options:       []string{"A", "B", "C", "D"},
				CorrectOption: "A",
				TimeLimitSecs: 10,
			},
		},
	}

	mockCore := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/package") {
			response.JSON(w, http.StatusOK, mockMatchPkg)
			return
		}
		http.NotFound(w, r)
	}))
	defer mockCore.Close()

	cfg := &config.Config{
		Port:               "8001",
		Environment:        "production",
		JWTSecret:          jwtSecret,
		AccessTokenTTL:     15 * time.Minute,
		JWTExpirationHours: 1 * time.Hour,
		CoreServiceURL:     mockCore.URL,
		AllowedOrigins:     []string{"https://app.palabrario.org", "wails://wails.localhost"},
	}

	coreClient := service.NewCoreClient(cfg.CoreServiceURL, mockCore.Client(), tokenService)
	hub := service.NewArenaHub(coreClient, tokenService)
	server := httptest.NewServer(NewRouter(cfg, hub, tokenService).Setup())
	defer server.Close()

	userTok, err := tokenService.GenerateToken(&models.UserProfile{
		IDUser:        "user-1",
		Email:         "u1@palabrario.org",
		UserName:      "u1",
		Role:          models.RoleLearner,
		AnchorDialect: "MX",
	})
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/arena/MX4321"

	t.Run("disallowed Origin rejected with 403", func(t *testing.T) {
		headers := http.Header{
			"Authorization": {"Bearer " + userTok},
			"Origin":        {"https://evil.example.com"},
		}
		ws, resp, err := websocket.DefaultDialer.Dial(wsURL, headers)
		if err == nil {
			_ = ws.Close()
			t.Fatal("expected WebSocket dial from disallowed origin to fail")
		}
		if resp == nil || resp.StatusCode != http.StatusForbidden {
			t.Fatalf("expected 403 Forbidden on disallowed origin, got resp=%v err=%v", resp, err)
		}
	})

	t.Run("lowercase bearer with allowed Origin connects cleanly", func(t *testing.T) {
		headers := http.Header{
			"Authorization": {"bearer " + userTok},
			"Origin":        {"wails://wails.localhost"},
		}
		ws, resp, err := websocket.DefaultDialer.Dial(wsURL, headers)
		if err != nil {
			t.Fatalf("expected WebSocket dial with lowercase bearer and allowed origin to succeed, got resp=%v err=%v", resp, err)
		}
		_ = ws.Close()
	})

	t.Run("empty Origin allowed for native Go desktop clients", func(t *testing.T) {
		headers := http.Header{
			"Authorization": {"Bearer " + userTok},
		}
		ws, resp, err := websocket.DefaultDialer.Dial(wsURL, headers)
		if err != nil {
			t.Fatalf("expected WebSocket dial with empty origin to succeed, got resp=%v err=%v", resp, err)
		}
		_ = ws.Close()
	})
}
