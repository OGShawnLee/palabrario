package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/gorilla/websocket"

	"palabrario/pkg/models"
)

func testUserProfile(role models.AccountRole, dialect string) models.UserProfile {
	return models.UserProfile{
		IDUser:        "user-123",
		Email:         "ana@palabrario.org",
		Name:          "Ana Dialect",
		UserName:      "ana_dialect",
		AnchorDialect: dialect,
		Role:          role,
	}
}

func writeJSON(t *testing.T, w http.ResponseWriter, status int, payload any) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload != nil {
		if err := json.NewEncoder(w).Encode(payload); err != nil {
			t.Fatalf("failed to encode JSON response: %v", err)
		}
	}
}

func TestSignInPersistsRefreshTokenAndSignOutClears(t *testing.T) {
	sessionFile := filepath.Join(t.TempDir(), "palabrario", "session.json")

	var signOutToken atomic.Value
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/api/v1/auth/sign-in":
			writeJSON(t, w, http.StatusOK, models.AuthResponse{
				Token:        "access-1",
				RefreshToken: "refresh-1",
				User:         testUserProfile(models.RoleLearner, "MX"),
			})
		case r.Method == http.MethodPost && r.URL.Path == "/api/v1/auth/sign-out":
			var req models.SignOutRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Fatalf("decode sign-out body: %v", err)
			}
			signOutToken.Store(req.RefreshToken)
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	}))
	defer srv.Close()

	app := newAppWithConfig(srv.URL, srv.URL, srv.URL, sessionFile)
	profile, err := app.SignIn(models.SignInRequest{Email: "ana@palabrario.org", Password: "password123"})
	if err != nil {
		t.Fatalf("SignIn failed: %v", err)
	}
	if profile == nil || profile.Email != "ana@palabrario.org" {
		t.Fatalf("unexpected profile: %+v", profile)
	}

	if got := app.getAccessToken(); got != "access-1" {
		t.Fatalf("expected accessToken=access-1, got %q", got)
	}
	if got := app.getRefreshToken(); got != "refresh-1" {
		t.Fatalf("expected refreshToken=refresh-1, got %q", got)
	}

	info, err := os.Stat(sessionFile)
	if err != nil {
		t.Fatalf("expected session.json to exist: %v", err)
	}
	if runtime.GOOS != "windows" && info.Mode().Perm() != 0600 {
		t.Fatalf("expected session.json permissions 0600, got %o", info.Mode().Perm())
	}

	raw, err := os.ReadFile(sessionFile)
	if err != nil {
		t.Fatalf("read session.json: %v", err)
	}
	var disk persistedSession
	if err := json.Unmarshal(raw, &disk); err != nil {
		t.Fatalf("unmarshal session.json: %v", err)
	}
	if disk.RefreshToken != "refresh-1" {
		t.Fatalf("expected persisted refresh_token=refresh-1, got %q", disk.RefreshToken)
	}

	if err := app.SignOut(); err != nil {
		t.Fatalf("SignOut failed: %v", err)
	}
	if got := app.getAccessToken(); got != "" {
		t.Fatalf("expected accessToken cleared, got %q", got)
	}
	if got := app.getRefreshToken(); got != "" {
		t.Fatalf("expected refreshToken cleared, got %q", got)
	}
	if _, err := os.Stat(sessionFile); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected session.json deleted on SignOut, stat err=%v", err)
	}
	if got, _ := signOutToken.Load().(string); got != "refresh-1" {
		t.Fatalf("expected POST /auth/sign-out to receive refresh-1, got %q", got)
	}
}

func TestRestartRecoveryGetMeRefreshesPersistedSession(t *testing.T) {
	sessionFile := filepath.Join(t.TempDir(), "palabrario", "session.json")

	var refreshCalls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/api/v1/auth/sign-in":
			writeJSON(t, w, http.StatusOK, models.AuthResponse{
				Token:        "access-initial",
				RefreshToken: "refresh-initial",
				User:         testUserProfile(models.RoleLearner, "MX"),
			})
		case r.Method == http.MethodPost && r.URL.Path == "/api/v1/auth/refresh":
			refreshCalls.Add(1)
			var req models.RefreshTokenRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Fatalf("decode refresh request: %v", err)
			}
			if req.RefreshToken != "refresh-initial" {
				t.Fatalf("expected refresh-initial, got %q", req.RefreshToken)
			}
			writeJSON(t, w, http.StatusOK, models.AuthResponse{
				Token:        "access-recovered",
				RefreshToken: "refresh-rotated",
				User:         testUserProfile(models.RoleEditor, "AR"),
			})
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	}))
	defer srv.Close()

	firstApp := newAppWithConfig(srv.URL, srv.URL, srv.URL, sessionFile)
	if _, err := firstApp.SignIn(models.SignInRequest{Email: "ana@palabrario.org", Password: "password123"}); err != nil {
		t.Fatalf("SignIn failed: %v", err)
	}

	// Simulate desktop restart: new App instance has empty accessToken in RAM and reads session.json.
	restartedApp := newAppWithConfig(srv.URL, srv.URL, srv.URL, sessionFile)
	if restartedApp.getAccessToken() != "" {
		t.Fatalf("expected empty accessToken on fresh restart")
	}
	if restartedApp.getRefreshToken() != "refresh-initial" {
		t.Fatalf("expected refresh-initial loaded from disk, got %q", restartedApp.getRefreshToken())
	}

	profile, err := restartedApp.GetMe()
	if err != nil {
		t.Fatalf("GetMe failed: %v", err)
	}
	if profile == nil || profile.Role != models.RoleEditor || profile.AnchorDialect != "AR" {
		t.Fatalf("unexpected recovered profile: %+v", profile)
	}
	if refreshCalls.Load() != 1 {
		t.Fatalf("expected 1 refresh call, got %d", refreshCalls.Load())
	}
	if got := restartedApp.getAccessToken(); got != "access-recovered" {
		t.Fatalf("expected access-recovered, got %q", got)
	}
	if got := restartedApp.getRefreshToken(); got != "refresh-rotated" {
		t.Fatalf("expected refresh-rotated, got %q", got)
	}
}

func TestGetMeUnauthenticatedAndCorruptFileAndTransientError(t *testing.T) {
	sessionFile := filepath.Join(t.TempDir(), "palabrario", "session.json")

	var serverHits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serverHits.Add(1)
		writeJSON(t, w, http.StatusInternalServerError, map[string]string{
			"code":  "INTERNAL_ERROR",
			"error": "database unavailable",
		})
	}))
	defer srv.Close()

	// 1. No tokens -> returns (nil, nil) with zero HTTP calls.
	app := newAppWithConfig(srv.URL, srv.URL, srv.URL, sessionFile)
	profile, err := app.GetMe()
	if err != nil || profile != nil {
		t.Fatalf("expected (nil, nil) when logged out, got (%+v, %v)", profile, err)
	}
	if serverHits.Load() != 0 {
		t.Fatalf("expected 0 HTTP calls when unauthenticated, got %d", serverHits.Load())
	}

	// 2. Corrupt session.json -> ignored gracefully, returns (nil, nil).
	if err := os.MkdirAll(filepath.Dir(sessionFile), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(sessionFile, []byte("{corrupt"), 0600); err != nil {
		t.Fatal(err)
	}
	corruptApp := newAppWithConfig(srv.URL, srv.URL, srv.URL, sessionFile)
	profile, err = corruptApp.GetMe()
	if err != nil || profile != nil {
		t.Fatalf("expected (nil, nil) with corrupt session.json, got (%+v, %v)", profile, err)
	}

	// 3. Transient 500 during startup refresh -> returns *APIError and preserves session.json.
	app.setTokens("", "refresh-valid")
	profile, err = app.GetMe()
	if profile != nil {
		t.Fatalf("expected nil profile on 500, got %+v", profile)
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) || apiErr.Status != http.StatusInternalServerError {
		t.Fatalf("expected *APIError status 500, got %T %v", err, err)
	}
	if _, statErr := os.Stat(sessionFile); statErr != nil {
		t.Fatalf("expected session.json preserved on transient 500, err=%v", statErr)
	}
}

func TestSingleFlightAutoRefreshAcrossConcurrentRequestsAndUploadImage(t *testing.T) {
	sessionFile := filepath.Join(t.TempDir(), "palabrario", "session.json")

	var refreshCalls atomic.Int32
	var protectedHits atomic.Int32
	var uploadHits atomic.Int32

	// Coordinate 4 concurrent calls so all hit the first 401 before refresh completes.
	const workers = 4
	var staleGate sync.WaitGroup
	staleGate.Add(workers)
	var gateOnce sync.Once

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/api/v1/auth/refresh":
			n := refreshCalls.Add(1)
			var req models.RefreshTokenRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Fatalf("decode refresh: %v", err)
			}
			if req.RefreshToken != "refresh-old" || n > 1 {
				writeJSON(t, w, http.StatusUnauthorized, map[string]string{
					"code":  "TOKEN_REVOKED",
					"error": "Refresh token has already been revoked.",
				})
				return
			}
			writeJSON(t, w, http.StatusOK, models.AuthResponse{
				Token:        "access-fresh",
				RefreshToken: "refresh-new",
				User:         testUserProfile(models.RoleLearner, "MX"),
			})

		case r.Method == http.MethodGet && r.URL.Path == "/api/v1/regions":
			protectedHits.Add(1)
			authHeader := r.Header.Get("Authorization")
			if authHeader == "Bearer access-expired" {
				staleGate.Done()
				staleGate.Wait()
				writeJSON(t, w, http.StatusUnauthorized, map[string]string{
					"code":  "TOKEN_EXPIRED",
					"error": "Access token has expired.",
				})
				return
			}
			if authHeader != "Bearer access-fresh" {
				writeJSON(t, w, http.StatusUnauthorized, map[string]string{
					"code":  "INVALID_TOKEN",
					"error": "Unexpected token: " + authHeader,
				})
				return
			}
			writeJSON(t, w, http.StatusOK, []models.Region{{Code: "MX", Name: "México"}})

		case r.Method == http.MethodPost && r.URL.Path == "/api/v1/images/upload":
			uploadHits.Add(1)
			authHeader := r.Header.Get("Authorization")
			if err := r.ParseMultipartForm(1 << 20); err != nil {
				t.Fatalf("parse multipart upload: %v", err)
			}
			file, _, err := r.FormFile("file")
			if err != nil {
				t.Fatalf("missing uploaded file: %v", err)
			}
			content, _ := io.ReadAll(file)
			_ = file.Close()
			if string(content) != "webp-payload" {
				t.Fatalf("expected multipart file payload preserved on retry, got %q", string(content))
			}

			if authHeader == "Bearer access-expired" {
				staleGate.Done()
				staleGate.Wait()
				writeJSON(t, w, http.StatusUnauthorized, map[string]string{
					"code":  "TOKEN_EXPIRED",
					"error": "Access token has expired.",
				})
				return
			}
			if authHeader != "Bearer access-fresh" {
				writeJSON(t, w, http.StatusUnauthorized, map[string]string{
					"code":  "INVALID_TOKEN",
					"error": "Unexpected token in upload: " + authHeader,
				})
				return
			}
			writeJSON(t, w, http.StatusCreated, models.ImageUploadResponse{URL: "/uploads/avatar/a.webp"})

		default:
			gateOnce.Do(func() {})
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	}))
	defer srv.Close()

	app := newAppWithConfig(srv.URL, srv.URL, srv.URL, sessionFile)
	app.setTokens("access-expired", "refresh-old")

	var wg sync.WaitGroup
	errs := make(chan error, workers)
	for i := 0; i < workers-1; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			regions, err := app.ListRegions()
			if err != nil {
				errs <- err
				return
			}
			if len(regions) != 1 || regions[0].Code != "MX" {
				errs <- errors.New("unexpected regions response")
			}
		}()
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		up, err := app.UploadImage("avatar", "me.webp", []byte("webp-payload"))
		if err != nil {
			errs <- err
			return
		}
		if up == nil || up.URL != "/uploads/avatar/a.webp" {
			errs <- errors.New("unexpected upload response")
		}
	}()

	wg.Wait()
	close(errs)
	for err := range errs {
		t.Fatalf("concurrent request failed: %v", err)
	}

	if refreshCalls.Load() != 1 {
		t.Fatalf("expected single-flight refresh to run exactly once, got %d", refreshCalls.Load())
	}
	if protectedHits.Load() != 6 {
		t.Fatalf("expected 3 ListRegions calls x 2 attempts = 6 hits, got %d", protectedHits.Load())
	}
	if uploadHits.Load() != 2 {
		t.Fatalf("expected 1 UploadImage x 2 attempts = 2 hits, got %d", uploadHits.Load())
	}
	if app.getAccessToken() != "access-fresh" || app.getRefreshToken() != "refresh-new" {
		t.Fatalf("expected rotated tokens saved, got access=%q refresh=%q", app.getAccessToken(), app.getRefreshToken())
	}
}

func TestSuspendedAccountAndFailedRefreshEvictSession(t *testing.T) {
	t.Run("403 ACCOUNT_SUSPENDED clears RAM and deletes session.json", func(t *testing.T) {
		sessionFile := filepath.Join(t.TempDir(), "palabrario", "session.json")
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			writeJSON(t, w, http.StatusForbidden, map[string]string{
				"code":  "ACCOUNT_SUSPENDED",
				"error": "This account has been suspended by an administrator.",
			})
		}))
		defer srv.Close()

		app := newAppWithConfig(srv.URL, srv.URL, srv.URL, sessionFile)
		app.setTokens("access-1", "refresh-1")

		err := app.ChangePassword(models.ChangePasswordRequest{
			CurrentPassword: "password123",
			NewPassword:     "newpassword123",
		})
		var apiErr *APIError
		if !errors.As(err, &apiErr) || apiErr.Status != http.StatusForbidden || apiErr.Code != "ACCOUNT_SUSPENDED" {
			t.Fatalf("expected 403 ACCOUNT_SUSPENDED APIError, got %T %v", err, err)
		}
		if app.getAccessToken() != "" || app.getRefreshToken() != "" {
			t.Fatalf("expected tokens cleared on suspension")
		}
		if _, statErr := os.Stat(sessionFile); !errors.Is(statErr, os.ErrNotExist) {
			t.Fatalf("expected session.json deleted on suspension, err=%v", statErr)
		}
	})

	t.Run("Failed 401 refresh clears RAM and deletes session.json", func(t *testing.T) {
		sessionFile := filepath.Join(t.TempDir(), "palabrario", "session.json")
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/api/v1/regions":
				writeJSON(t, w, http.StatusUnauthorized, map[string]string{
					"code":  "TOKEN_EXPIRED",
					"error": "Access token has expired.",
				})
			case "/api/v1/auth/refresh":
				writeJSON(t, w, http.StatusUnauthorized, map[string]string{
					"code":  "TOKEN_REVOKED",
					"error": "Refresh token has already been revoked.",
				})
			default:
				t.Fatalf("unexpected path: %s", r.URL.Path)
			}
		}))
		defer srv.Close()

		app := newAppWithConfig(srv.URL, srv.URL, srv.URL, sessionFile)
		app.setTokens("access-1", "refresh-revoked")

		_, err := app.ListRegions()
		var apiErr *APIError
		if !errors.As(err, &apiErr) || apiErr.Status != http.StatusUnauthorized || apiErr.Code != "TOKEN_REVOKED" {
			t.Fatalf("expected 401 TOKEN_REVOKED APIError, got %T %v", err, err)
		}
		if app.getAccessToken() != "" || app.getRefreshToken() != "" {
			t.Fatalf("expected tokens cleared after failed refresh")
		}
		if _, statErr := os.Stat(sessionFile); !errors.Is(statErr, os.ErrNotExist) {
			t.Fatalf("expected session.json deleted after failed refresh, err=%v", statErr)
		}
	})
}

func TestUpdateProfileAndRedeemTicketRotateAccessToken(t *testing.T) {
	sessionFile := filepath.Join(t.TempDir(), "palabrario", "session.json")

	var refreshCount atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPut && r.URL.Path == "/api/v1/users/me":
			if r.Header.Get("Authorization") != "Bearer access-1" {
				t.Fatalf("expected access-1 on UpdateProfile, got %q", r.Header.Get("Authorization"))
			}
			writeJSON(t, w, http.StatusOK, testUserProfile(models.RoleLearner, "CO"))

		case r.Method == http.MethodPost && r.URL.Path == "/api/v1/governance/tickets/redeem":
			if r.Header.Get("Authorization") != "Bearer access-2" {
				t.Fatalf("expected rotated access-2 on RedeemInvitationTicket, got %q", r.Header.Get("Authorization"))
			}
			writeJSON(t, w, http.StatusOK, testUserProfile(models.RoleEditor, "CO"))

		case r.Method == http.MethodPost && r.URL.Path == "/api/v1/auth/refresh":
			n := refreshCount.Add(1)
			if n == 1 {
				writeJSON(t, w, http.StatusOK, models.AuthResponse{
					Token:        "access-2",
					RefreshToken: "refresh-2",
					User:         testUserProfile(models.RoleLearner, "CO"),
				})
				return
			}
			writeJSON(t, w, http.StatusOK, models.AuthResponse{
				Token:        "access-3",
				RefreshToken: "refresh-3",
				User:         testUserProfile(models.RoleEditor, "CO"),
			})

		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	}))
	defer srv.Close()

	app := newAppWithConfig(srv.URL, srv.URL, srv.URL, sessionFile)
	app.setTokens("access-1", "refresh-1")

	updated, err := app.UpdateProfile(models.UpdateProfileRequest{Name: "Ana Dialect", AnchorDialect: "CO"})
	if err != nil {
		t.Fatalf("UpdateProfile failed: %v", err)
	}
	if updated.AnchorDialect != "CO" {
		t.Fatalf("expected CO dialect, got %q", updated.AnchorDialect)
	}
	if app.getAccessToken() != "access-2" {
		t.Fatalf("expected accessToken rotated to access-2, got %q", app.getAccessToken())
	}

	promoted, err := app.RedeemInvitationTicket(models.RedeemTicketRequest{Code: "EDIT-2026"})
	if err != nil {
		t.Fatalf("RedeemInvitationTicket failed: %v", err)
	}
	if promoted.Role != models.RoleEditor {
		t.Fatalf("expected EDITOR role, got %q", promoted.Role)
	}
	if app.getAccessToken() != "access-3" || app.getRefreshToken() != "refresh-3" {
		t.Fatalf("expected access-3/refresh-3 after ticket redeem, got %q/%q", app.getAccessToken(), app.getRefreshToken())
	}
}

func TestConnectArenaRoomUsesAuthorizationBearerHeader(t *testing.T) {
	sessionFile := filepath.Join(t.TempDir(), "palabrario", "session.json")

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	var gotAuth atomic.Value
	var gotQueryToken atomic.Value
	var gotPath atomic.Value

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && r.URL.Path == "/api/v1/auth/refresh" {
			writeJSON(t, w, http.StatusOK, models.AuthResponse{
				Token:        "ws-access-token",
				RefreshToken: "ws-refresh-rotated",
				User:         testUserProfile(models.RoleLearner, "MX"),
			})
			return
		}

		if strings.HasPrefix(r.URL.Path, "/ws/arena/") {
			gotAuth.Store(r.Header.Get("Authorization"))
			gotQueryToken.Store(r.URL.Query().Get("token"))
			gotPath.Store(r.URL.Path)

			conn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				t.Fatalf("ws upgrade failed: %v", err)
			}
			defer conn.Close()
			for {
				if _, _, err := conn.ReadMessage(); err != nil {
					return
				}
			}
		}

		t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
	}))
	defer srv.Close()

	app := newAppWithConfig(srv.URL, srv.URL, srv.URL, sessionFile)
	// Start with only persisted refresh token to also verify startup refresh before WS dial.
	app.setTokens("", "ws-refresh-initial")
	app.SetContext(context.Background())

	if err := app.ConnectArenaRoom("ROOM42"); err != nil {
		t.Fatalf("ConnectArenaRoom failed: %v", err)
	}
	defer func() { _ = app.LeaveArenaRoom() }()

	if auth, _ := gotAuth.Load().(string); auth != "Bearer ws-access-token" {
		t.Fatalf("expected Authorization header 'Bearer ws-access-token', got %q", auth)
	}
	if qt, _ := gotQueryToken.Load().(string); qt != "" {
		t.Fatalf("expected empty ?token= query param, got %q", qt)
	}
	if p, _ := gotPath.Load().(string); p != "/ws/arena/ROOM42" {
		t.Fatalf("expected /ws/arena/ROOM42, got %q", p)
	}
}

func TestFindMatchRoomUsesAPIErrorNotFoundAndNetworkErrorSanitized(t *testing.T) {
	sessionFile := filepath.Join(t.TempDir(), "palabrario", "session.json")

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/arena/matches/code/MISSING":
			writeJSON(t, w, http.StatusNotFound, map[string]string{
				"code":  "NOT_FOUND",
				"error": "Match room could not be located.",
			})
		case "/api/v1/arena/matches/code/BOOM":
			writeJSON(t, w, http.StatusInternalServerError, map[string]string{
				"code":  "INTERNAL_ERROR",
				"error": "Internal failure",
			})
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))

	app := newAppWithConfig(srv.URL, srv.URL, srv.URL, sessionFile)

	item, err := app.FindMatchRoom("#missing")
	if err != nil || item != nil {
		t.Fatalf("expected (nil, nil) on 404, got (%+v, %v)", item, err)
	}

	_, err = app.FindMatchRoom("BOOM")
	var apiErr *APIError
	if !errors.As(err, &apiErr) || apiErr.Status != http.StatusInternalServerError {
		t.Fatalf("expected *APIError 500, got %T %v", err, err)
	}

	// Close server to test English sanitized network error.
	srv.Close()
	_, err = app.ListRegions()
	if !errors.As(err, &apiErr) || apiErr.Code != "NETWORK_ERROR" {
		t.Fatalf("expected NETWORK_ERROR APIError, got %T %v", err, err)
	}
	if err.Error() != "Unable to reach the Palabrario server. Please check your connection." {
		t.Fatalf("unexpected network error message: %q", err.Error())
	}
}
