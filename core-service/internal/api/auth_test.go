package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"palabrario/core-service/internal/api"
	"palabrario/core-service/internal/config"
	"palabrario/core-service/internal/database"
	"palabrario/core-service/internal/testutil"
	"palabrario/pkg/auth"
	"palabrario/pkg/models"
	"palabrario/pkg/response"
	"palabrario/pkg/utils"
)

const testJWTSecret = "test-super-secret-jwt-key-min-32-chars-long"

func mintExpiredAccessJWT(t *testing.T, user models.UserProfile, tokenVersion int) string {
	t.Helper()
	now := time.Now().UTC()
	claims := auth.Claims{
		UserID:        user.IDUser,
		Email:         user.Email,
		UserName:      user.UserName,
		AnchorDialect: user.AnchorDialect,
		Role:          user.Role,
		TokenVersion:  tokenVersion,
		TokenType:     auth.TokenTypeAccess,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(-2 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(now.Add(-17 * time.Minute)),
			NotBefore: jwt.NewNumericDate(now.Add(-17 * time.Minute)),
			Issuer:    "palabrario-core-service",
			Subject:   user.IDUser,
		},
	}
	tok, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(testJWTSecret))
	if err != nil {
		t.Fatalf("failed to sign expired access JWT: %v", err)
	}
	return tok
}

func newTestHTTPClient() *http.Client {
	return &http.Client{Timeout: 10 * time.Second}
}

func setupTestServer(t *testing.T) (*httptest.Server, *database.DB) {
	t.Helper()

	db := testutil.SetupTestDB(t)
	if db == nil {
		return nil, nil
	}

	cfg := &config.Config{
		Port:               "8000",
		Environment:        "test",
		JWTSecret:          testJWTSecret,
		JWTExpirationHours: 168 * time.Hour,
		AccessTokenTTL:     15 * time.Minute,
		RefreshTokenTTL:    168 * time.Hour,
		AllowedOrigins:     auth.DefaultDevCORSOrigins(),
	}

	router := api.NewRouter(cfg, db)
	server := httptest.NewServer(router.Setup())

	t.Cleanup(func() {
		server.Close()
	})

	return server, db
}

func doJSONRequest(t *testing.T, client *http.Client, method, url, bearerToken string, payload any) *http.Response {
	t.Helper()
	var bodyReader *bytes.Reader
	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("failed to marshal request body: %v", err)
		}
		bodyReader = bytes.NewReader(b)
	} else {
		bodyReader = bytes.NewReader(nil)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		t.Fatalf("failed to create HTTP request: %v", err)
	}
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("HTTP %s %s failed: %v", method, url, err)
	}
	return resp
}

func decodeErrorResponse(t *testing.T, resp *http.Response) response.ErrorBody {
	t.Helper()
	var errResp response.ErrorBody
	if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
		t.Fatalf("failed to decode ErrorBody JSON: %v", err)
	}
	return errResp
}

func signInHTTP(t *testing.T, client *http.Client, serverURL, email, password string) models.AuthResponse {
	t.Helper()
	resp := doJSONRequest(t, client, http.MethodPost, serverURL+"/api/v1/auth/sign-in", "", models.SignInRequest{
		Email:    email,
		Password: password,
	})
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK on sign-in for %s, got %d", email, resp.StatusCode)
	}

	var authResp models.AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		t.Fatalf("failed to decode AuthResponse: %v", err)
	}
	return authResp
}

func TestAuth_SignIn_Success(t *testing.T) {
	server, _ := setupTestServer(t)
	if server == nil {
		return
	}

	client := newTestHTTPClient()
	authResp := signInHTTP(t, client, server.URL, "silas@palabrario.org", "password123")

	if authResp.Token == "" {
		t.Error("expected non-empty access JWT")
	}
	if authResp.RefreshToken == "" {
		t.Error("expected non-empty refresh token")
	}
	if authResp.User.Email != "silas@palabrario.org" {
		t.Errorf("expected email 'silas@palabrario.org', got '%s'", authResp.User.Email)
	}
	if authResp.User.Role != models.RoleEditor {
		t.Errorf("expected role 'EDITOR', got '%s'", authResp.User.Role)
	}
}

func TestAuth_SignIn_InvalidPassword(t *testing.T) {
	server, _ := setupTestServer(t)
	if server == nil {
		return
	}

	client := newTestHTTPClient()
	resp := doJSONRequest(t, client, http.MethodPost, server.URL+"/api/v1/auth/sign-in", "", models.SignInRequest{
		Email:    "silas@palabrario.org",
		Password: "wrongpassword!",
	})
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401 Unauthorized, got %d", resp.StatusCode)
	}
	errResp := decodeErrorResponse(t, resp)
	if errResp.Code != "INVALID_CREDENTIALS" {
		t.Errorf("expected code INVALID_CREDENTIALS, got %q", errResp.Code)
	}
}

func TestAuth_GetMe_ProtectedEndpoint(t *testing.T) {
	server, _ := setupTestServer(t)
	if server == nil {
		return
	}

	client := newTestHTTPClient()
	authResp := signInHTTP(t, client, server.URL, "admin@palabrario.org", "password123")

	// Authorized request with valid token
	meResp := doJSONRequest(t, client, http.MethodGet, server.URL+"/api/v1/auth/me", authResp.Token, nil)
	defer meResp.Body.Close()

	if meResp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK with valid token, got %d", meResp.StatusCode)
	}

	var profile models.UserProfile
	if err := json.NewDecoder(meResp.Body).Decode(&profile); err != nil {
		t.Fatalf("failed to decode UserProfile: %v", err)
	}
	if profile.Role != models.RoleAdmin {
		t.Errorf("expected role ADMIN, got %s", profile.Role)
	}

	// Missing token -> 401 MISSING_TOKEN
	unauthResp := doJSONRequest(t, client, http.MethodGet, server.URL+"/api/v1/auth/me", "", nil)
	defer unauthResp.Body.Close()

	if unauthResp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401 Unauthorized without token, got %d", unauthResp.StatusCode)
	}
	errResp := decodeErrorResponse(t, unauthResp)
	if errResp.Code != "MISSING_TOKEN" {
		t.Errorf("expected code MISSING_TOKEN, got %q", errResp.Code)
	}
}

func TestAuth_SignUp_Lifecycle_And_Validation(t *testing.T) {
	server, _ := setupTestServer(t)
	if server == nil {
		return
	}

	client := newTestHTTPClient()

	// 1. Validation failures return 400 VALIDATION_ERROR
	invalidPayloads := []struct {
		name string
		req  models.SignUpRequest
	}{
		{
			name: "short password",
			req: models.SignUpRequest{
				Email:         "short@palabrario.org",
				Password:      "short",
				Name:          "Short Pass",
				UserName:      "short_pass_u",
				AnchorDialect: "MX",
			},
		},
		{
			name: "password over 72 bytes",
			req: models.SignUpRequest{
				Email:         "over72@palabrario.org",
				Password:      strings.Repeat("a", 73),
				Name:          "Long Pass",
				UserName:      "long_pass_u",
				AnchorDialect: "MX",
			},
		},
		{
			name: "invalid email",
			req: models.SignUpRequest{
				Email:         "invalid-email",
				Password:      "SecureTestPass123!",
				Name:          "Bad Email",
				UserName:      "bad_email_u",
				AnchorDialect: "MX",
			},
		},
		{
			name: "invalid username characters",
			req: models.SignUpRequest{
				Email:         "badhandle@palabrario.org",
				Password:      "SecureTestPass123!",
				Name:          "Bad Handle",
				UserName:      "bad handle!",
				AnchorDialect: "MX",
			},
		},
		{
			name: "username over 50 chars",
			req: models.SignUpRequest{
				Email:         "longhandle@palabrario.org",
				Password:      "SecureTestPass123!",
				Name:          "Long Handle",
				UserName:      strings.Repeat("a", 51),
				AnchorDialect: "MX",
			},
		},
		{
			name: "invalid dialect",
			req: models.SignUpRequest{
				Email:         "dialect@palabrario.org",
				Password:      "SecureTestPass123!",
				Name:          "Bad Dialect",
				UserName:      "bad_dialect_u",
				AnchorDialect: "ZZ",
			},
		},
	}

	for _, tc := range invalidPayloads {
		t.Run(tc.name, func(t *testing.T) {
			resp := doJSONRequest(t, client, http.MethodPost, server.URL+"/api/v1/auth/sign-up", "", tc.req)
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusBadRequest {
				t.Fatalf("expected 400 Bad Request for %s, got %d", tc.name, resp.StatusCode)
			}
			errResp := decodeErrorResponse(t, resp)
			if errResp.Code != "VALIDATION_ERROR" {
				t.Errorf("expected code VALIDATION_ERROR, got %q", errResp.Code)
			}
		})
	}

	// 2. Valid SignUp returns 201 + access & refresh tokens
	signUpReq := models.SignUpRequest{
		Email:         "nuevo_aprendiz@palabrario.org",
		Password:      "SecureTestPass123!",
		Name:          "Usuario Automatizado",
		UserName:      "autouser_co",
		AnchorDialect: "CO",
	}
	resp := doJSONRequest(t, client, http.MethodPost, server.URL+"/api/v1/auth/sign-up", "", signUpReq)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201 Created on sign-up, got %d", resp.StatusCode)
	}

	var authResp models.AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		t.Fatalf("failed to decode AuthResponse: %v", err)
	}
	if authResp.Token == "" || authResp.RefreshToken == "" {
		t.Errorf("expected both token and refresh_token on sign-up")
	}
	if authResp.User.AnchorDialect != "CO" || authResp.User.Role != models.RoleLearner {
		t.Errorf("unexpected registered profile: %+v", authResp.User)
	}

	// 3. Duplicate email -> 409 EMAIL_ALREADY_EXISTS
	dupEmailResp := doJSONRequest(t, client, http.MethodPost, server.URL+"/api/v1/auth/sign-up", "", models.SignUpRequest{
		Email:         "nuevo_aprendiz@palabrario.org",
		Password:      "SecureTestPass123!",
		Name:          "Dup Email",
		UserName:      "another_unique_handle",
		AnchorDialect: "MX",
	})
	defer dupEmailResp.Body.Close()
	if dupEmailResp.StatusCode != http.StatusConflict {
		t.Fatalf("expected 409 Conflict for duplicate email, got %d", dupEmailResp.StatusCode)
	}
	if errResp := decodeErrorResponse(t, dupEmailResp); errResp.Code != "EMAIL_ALREADY_EXISTS" {
		t.Errorf("expected EMAIL_ALREADY_EXISTS, got %q", errResp.Code)
	}

	// 4. Duplicate username -> 409 USERNAME_ALREADY_EXISTS
	dupUserResp := doJSONRequest(t, client, http.MethodPost, server.URL+"/api/v1/auth/sign-up", "", models.SignUpRequest{
		Email:         "different_email@palabrario.org",
		Password:      "SecureTestPass123!",
		Name:          "Dup User",
		UserName:      "autouser_co",
		AnchorDialect: "MX",
	})
	defer dupUserResp.Body.Close()
	if dupUserResp.StatusCode != http.StatusConflict {
		t.Fatalf("expected 409 Conflict for duplicate username, got %d", dupUserResp.StatusCode)
	}
	if errResp := decodeErrorResponse(t, dupUserResp); errResp.Code != "USERNAME_ALREADY_EXISTS" {
		t.Errorf("expected USERNAME_ALREADY_EXISTS, got %q", errResp.Code)
	}
}

func TestAuth_RefreshToken_SignOut_And_ExpiredAccessJWT(t *testing.T) {
	server, _ := setupTestServer(t)
	if server == nil {
		return
	}

	client := newTestHTTPClient()
	login := signInHTTP(t, client, server.URL, "gael@palabrario.org", "password123")

	// 1. Refresh rotates token pair
	refreshResp := doJSONRequest(t, client, http.MethodPost, server.URL+"/api/v1/auth/refresh", "", models.RefreshTokenRequest{
		RefreshToken: login.RefreshToken,
	})
	defer refreshResp.Body.Close()
	if refreshResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK on refresh, got %d", refreshResp.StatusCode)
	}
	var rotated models.AuthResponse
	if err := json.NewDecoder(refreshResp.Body).Decode(&rotated); err != nil {
		t.Fatalf("failed decoding refresh response: %v", err)
	}
	if rotated.Token == "" || rotated.RefreshToken == "" || rotated.RefreshToken == login.RefreshToken {
		t.Fatalf("expected rotated refresh token, got %+v", rotated)
	}

	// 2. Reusing previous refresh token returns 401 TOKEN_REVOKED
	reuseResp := doJSONRequest(t, client, http.MethodPost, server.URL+"/api/v1/auth/refresh", "", models.RefreshTokenRequest{
		RefreshToken: login.RefreshToken,
	})
	defer reuseResp.Body.Close()
	if reuseResp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized on refresh reuse, got %d", reuseResp.StatusCode)
	}
	if errResp := decodeErrorResponse(t, reuseResp); errResp.Code != "TOKEN_REVOKED" {
		t.Errorf("expected TOKEN_REVOKED, got %q", errResp.Code)
	}

	// 3. SignOut revokes active refresh token
	signOutResp := doJSONRequest(t, client, http.MethodPost, server.URL+"/api/v1/auth/sign-out", "", models.SignOutRequest{
		RefreshToken: rotated.RefreshToken,
	})
	defer signOutResp.Body.Close()
	if signOutResp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected 204 No Content on sign-out, got %d", signOutResp.StatusCode)
	}

	// Empty body sign-out also returns 204
	emptySignOut := doJSONRequest(t, client, http.MethodPost, server.URL+"/api/v1/auth/sign-out", "", nil)
	defer emptySignOut.Body.Close()
	if emptySignOut.StatusCode != http.StatusNoContent {
		t.Fatalf("expected 204 No Content on empty sign-out, got %d", emptySignOut.StatusCode)
	}

	// Signed-out token cannot be refreshed
	afterSignOut := doJSONRequest(t, client, http.MethodPost, server.URL+"/api/v1/auth/refresh", "", models.RefreshTokenRequest{
		RefreshToken: rotated.RefreshToken,
	})
	defer afterSignOut.Body.Close()
	if afterSignOut.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized after sign-out, got %d", afterSignOut.StatusCode)
	}

	// 4. Expired access JWT returns 401 TOKEN_EXPIRED
	expiredJWT := mintExpiredAccessJWT(t, login.User, 1)
	expResp := doJSONRequest(t, client, http.MethodGet, server.URL+"/api/v1/auth/me", expiredJWT, nil)
	defer expResp.Body.Close()
	if expResp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized for expired access JWT, got %d", expResp.StatusCode)
	}
	if errResp := decodeErrorResponse(t, expResp); errResp.Code != "TOKEN_EXPIRED" {
		t.Errorf("expected TOKEN_EXPIRED, got %q", errResp.Code)
	}
}

func TestAuth_OptionalAuth_RejectsInvalidAndExpiredHeaders(t *testing.T) {
	server, _ := setupTestServer(t)
	if server == nil {
		return
	}

	client := newTestHTTPClient()

	// 1. Anonymous request without header succeeds
	anonResp := doJSONRequest(t, client, http.MethodGet, server.URL+"/api/v1/users/u/gael_learner", "", nil)
	defer anonResp.Body.Close()
	if anonResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK without Authorization header, got %d", anonResp.StatusCode)
	}

	// 2. Malformed header returns 401 INVALID_TOKEN_FORMAT
	reqBadFmt, _ := http.NewRequest(http.MethodGet, server.URL+"/api/v1/users/u/gael_learner", nil)
	reqBadFmt.Header.Set("Authorization", "Basic abcdef")
	badFmtResp, err := client.Do(reqBadFmt)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer badFmtResp.Body.Close()
	if badFmtResp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401 on malformed Authorization header, got %d", badFmtResp.StatusCode)
	}
	if errResp := decodeErrorResponse(t, badFmtResp); errResp.Code != "INVALID_TOKEN_FORMAT" {
		t.Errorf("expected INVALID_TOKEN_FORMAT, got %q", errResp.Code)
	}

	// 3. Invalid Bearer token returns 401 INVALID_TOKEN on OptionalAuth routes
	invalidResp := doJSONRequest(t, client, http.MethodGet, server.URL+"/api/v1/words/search?q=popote", "invalid.jwt.token", nil)
	defer invalidResp.Body.Close()
	if invalidResp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401 on invalid Bearer token in OptionalAuth, got %d", invalidResp.StatusCode)
	}
	if errResp := decodeErrorResponse(t, invalidResp); errResp.Code != "INVALID_TOKEN" {
		t.Errorf("expected INVALID_TOKEN, got %q", errResp.Code)
	}

	// 4. Expired Bearer token returns 401 TOKEN_EXPIRED on OptionalAuth
	login := signInHTTP(t, client, server.URL, "gael@palabrario.org", "password123")
	expiredJWT := mintExpiredAccessJWT(t, login.User, 1)
	expResp := doJSONRequest(t, client, http.MethodGet, server.URL+"/api/v1/users/u/gael_learner", expiredJWT, nil)
	defer expResp.Body.Close()
	if expResp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401 on expired Bearer token in OptionalAuth, got %d", expResp.StatusCode)
	}
	if errResp := decodeErrorResponse(t, expResp); errResp.Code != "TOKEN_EXPIRED" {
		t.Errorf("expected TOKEN_EXPIRED, got %q", errResp.Code)
	}
}

func TestAuth_OptionB_SuspendedAndRoleChangeEnforcement(t *testing.T) {
	server, db := setupTestServer(t)
	if server == nil {
		return
	}

	client := newTestHTTPClient()
	ctx := context.Background()

	adminLogin := signInHTTP(t, client, server.URL, "admin@palabrario.org", "password123")
	gaelLogin := signInHTTP(t, client, server.URL, "gael@palabrario.org", "password123")
	silasLogin := signInHTTP(t, client, server.URL, "silas@palabrario.org", "password123")

	// 1. Direct DB role sync check (without token_version bump) to prove claims.Role is synced from DB on Option B routes
	if _, err := db.Pool.Exec(ctx, `UPDATE accounts SET role = 'ADMIN' WHERE email = 'silas@palabrario.org'`); err != nil {
		t.Fatalf("failed promoting silas in DB: %v", err)
	}
	silasAdminList := doJSONRequest(t, client, http.MethodGet, server.URL+"/api/v1/governance/users", silasLogin.Token, nil)
	defer silasAdminList.Body.Close()
	if silasAdminList.StatusCode != http.StatusOK {
		t.Fatalf("expected promoted user JWT to pass Option B RequireRole(ADMIN), got %d", silasAdminList.StatusCode)
	}

	if _, err := db.Pool.Exec(ctx, `UPDATE accounts SET role = 'LEARNER' WHERE email = 'silas@palabrario.org'`); err != nil {
		t.Fatalf("failed demoting silas in DB: %v", err)
	}
	silasDemotedWrite := doJSONRequest(t, client, http.MethodPost, server.URL+"/api/v1/words", silasLogin.Token, map[string]any{})
	defer silasDemotedWrite.Body.Close()
	if silasDemotedWrite.StatusCode != http.StatusForbidden {
		t.Fatalf("expected demoted user JWT to fail RequireRole(EDITOR,ADMIN) with 403, got %d", silasDemotedWrite.StatusCode)
	}

	// 2. Governance promotion bumps token_version -> old JWT gets 401 TOKEN_REVOKED
	promResp := doJSONRequest(t, client, http.MethodPut, server.URL+"/api/v1/governance/users/"+gaelLogin.User.IDUser+"/access", adminLogin.Token, models.UpdateUserRoleOrStatusRequest{
		Role: utils.ToPtr(models.RoleEditor),
	})
	defer promResp.Body.Close()
	if promResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK promoting gael via governance, got %d", promResp.StatusCode)
	}

	staleWrite := doJSONRequest(t, client, http.MethodPut, server.URL+"/api/v1/users/me", gaelLogin.Token, models.UpdateProfileRequest{
		Name:          "Gael Stale",
		AnchorDialect: "MX",
	})
	defer staleWrite.Body.Close()
	if staleWrite.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401 TOKEN_REVOKED on write after governance role change, got %d", staleWrite.StatusCode)
	}
	if errResp := decodeErrorResponse(t, staleWrite); errResp.Code != "TOKEN_REVOKED" {
		t.Errorf("expected TOKEN_REVOKED, got %q", errResp.Code)
	}

	// Re-login Gael to obtain fresh JWT + refresh token
	gaelFresh := signInHTTP(t, client, server.URL, "gael@palabrario.org", "password123")

	// 3. Suspend Gael via governance -> 403 ACCOUNT_SUSPENDED on /auth/me, writes, /governance/*, 401 TOKEN_REVOKED on already-revoked refresh token, and 404 on /users/u/{username}
	suspResp := doJSONRequest(t, client, http.MethodPut, server.URL+"/api/v1/governance/users/"+gaelLogin.User.IDUser+"/access", adminLogin.Token, models.UpdateUserRoleOrStatusRequest{
		IsActive: utils.ToPtr(false),
	})
	defer suspResp.Body.Close()
	if suspResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK suspending gael, got %d", suspResp.StatusCode)
	}

	suspendedChecks := []struct {
		name   string
		method string
		path   string
		body   any
	}{
		{name: "GET /auth/me", method: http.MethodGet, path: "/api/v1/auth/me"},
		{name: "PUT /users/me", method: http.MethodPut, path: "/api/v1/users/me", body: models.UpdateProfileRequest{Name: "Gael", AnchorDialect: "MX"}},
		{name: "POST /decks", method: http.MethodPost, path: "/api/v1/decks", body: map[string]any{"title": "Deck"}},
		{name: "GET /governance/editor-requests/me", method: http.MethodGet, path: "/api/v1/governance/editor-requests/me"},
	}

	for _, sc := range suspendedChecks {
		t.Run(sc.name, func(t *testing.T) {
			r := doJSONRequest(t, client, sc.method, server.URL+sc.path, gaelFresh.Token, sc.body)
			defer r.Body.Close()
			if r.StatusCode != http.StatusForbidden {
				t.Fatalf("expected 403 ACCOUNT_SUSPENDED on %s, got %d", sc.name, r.StatusCode)
			}
			if errResp := decodeErrorResponse(t, r); errResp.Code != "ACCOUNT_SUSPENDED" {
				t.Errorf("expected ACCOUNT_SUSPENDED on %s, got %q", sc.name, errResp.Code)
			}
		})
	}

	refreshRevokedByGovernance := doJSONRequest(t, client, http.MethodPost, server.URL+"/api/v1/auth/refresh", "", models.RefreshTokenRequest{
		RefreshToken: gaelFresh.RefreshToken,
	})
	defer refreshRevokedByGovernance.Body.Close()
	if refreshRevokedByGovernance.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401 TOKEN_REVOKED on /auth/refresh after governance suspension revoked tokens, got %d", refreshRevokedByGovernance.StatusCode)
	}

	// Also verify RefreshToken returns 403 ACCOUNT_SUSPENDED if an unrevoked refresh token exists for a suspended account
	if _, err := db.Pool.Exec(ctx, `UPDATE accounts SET is_active = false WHERE email = 'silas@palabrario.org'`); err != nil {
		t.Fatalf("failed suspending silas in DB: %v", err)
	}
	refreshSusp := doJSONRequest(t, client, http.MethodPost, server.URL+"/api/v1/auth/refresh", "", models.RefreshTokenRequest{
		RefreshToken: silasLogin.RefreshToken,
	})
	defer refreshSusp.Body.Close()
	if refreshSusp.StatusCode != http.StatusForbidden {
		t.Fatalf("expected 403 ACCOUNT_SUSPENDED on /auth/refresh for suspended user, got %d", refreshSusp.StatusCode)
	}
	if errResp := decodeErrorResponse(t, refreshSusp); errResp.Code != "ACCOUNT_SUSPENDED" {
		t.Errorf("expected ACCOUNT_SUSPENDED, got %q", errResp.Code)
	}

	pubSusp := doJSONRequest(t, client, http.MethodGet, server.URL+"/api/v1/users/u/gael_learner", "", nil)
	defer pubSusp.Body.Close()
	if pubSusp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404 NOT_FOUND for suspended user public profile, got %d", pubSusp.StatusCode)
	}
}

func TestAuth_ChangePassword_HTTP_And_CORS(t *testing.T) {
	server, _ := setupTestServer(t)
	if server == nil {
		return
	}

	client := newTestHTTPClient()
	gaelLogin := signInHTTP(t, client, server.URL, "gael@palabrario.org", "password123")

	// 1. Invalid JSON body -> 400 INVALID_JSON
	reqBadJSON, _ := http.NewRequest(http.MethodPut, server.URL+"/api/v1/users/me/password", strings.NewReader("{invalid"))
	reqBadJSON.Header.Set("Authorization", "Bearer "+gaelLogin.Token)
	reqBadJSON.Header.Set("Content-Type", "application/json")
	badJSONResp, err := client.Do(reqBadJSON)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer badJSONResp.Body.Close()
	if badJSONResp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 on invalid JSON, got %d", badJSONResp.StatusCode)
	}

	// 2. Wrong current password -> 401 INVALID_CREDENTIALS
	wrongResp := doJSONRequest(t, client, http.MethodPut, server.URL+"/api/v1/users/me/password", gaelLogin.Token, models.ChangePasswordRequest{
		CurrentPassword: "wrongpassword",
		NewPassword:     "NewPassword123!",
	})
	defer wrongResp.Body.Close()
	if wrongResp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401 INVALID_CREDENTIALS, got %d", wrongResp.StatusCode)
	}
	if errResp := decodeErrorResponse(t, wrongResp); errResp.Code != "INVALID_CREDENTIALS" {
		t.Errorf("expected INVALID_CREDENTIALS, got %q", errResp.Code)
	}

	// 3. Valid change -> 204 No Content & invalidates old JWT + refresh token
	okResp := doJSONRequest(t, client, http.MethodPut, server.URL+"/api/v1/users/me/password", gaelLogin.Token, models.ChangePasswordRequest{
		CurrentPassword: "password123",
		NewPassword:     "NewPassword123!",
	})
	defer okResp.Body.Close()
	if okResp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected 204 No Content on password change, got %d", okResp.StatusCode)
	}

	revokedMe := doJSONRequest(t, client, http.MethodGet, server.URL+"/api/v1/auth/me", gaelLogin.Token, nil)
	defer revokedMe.Body.Close()
	if revokedMe.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401 TOKEN_REVOKED for prior JWT after password change, got %d", revokedMe.StatusCode)
	}

	// 4. Sign in with new password works
	_ = signInHTTP(t, client, server.URL, "gael@palabrario.org", "NewPassword123!")

	// 5. CORS check: AllowCredentials must be false and allowed origin echoed
	corsReq, _ := http.NewRequest(http.MethodOptions, server.URL+"/api/v1/auth/sign-in", nil)
	corsReq.Header.Set("Origin", "http://localhost:5173")
	corsReq.Header.Set("Access-Control-Request-Method", "POST")
	corsResp, err := client.Do(corsReq)
	if err != nil {
		t.Fatalf("CORS preflight failed: %v", err)
	}
	defer corsResp.Body.Close()
	if gotOrigin := corsResp.Header.Get("Access-Control-Allow-Origin"); gotOrigin != "http://localhost:5173" {
		t.Errorf("expected Access-Control-Allow-Origin http://localhost:5173, got %q", gotOrigin)
	}
	if gotCreds := corsResp.Header.Get("Access-Control-Allow-Credentials"); gotCreds == "true" {
		t.Errorf("expected Access-Control-Allow-Credentials not to be 'true', got %q", gotCreds)
	}
}
