# 032 — Changelog: Harden Wails BFF Session Persistence, Auto-Refresh, Error Codes, and Tests

## Summary

Hardened the Wails v3 Go BFF (`app/`) to custody short-lived access JWTs in RAM alongside an atomic, host-persisted opaque refresh token, added single-flight `401` auto-refresh with `403 ACCOUNT_SUSPENDED` / failed-refresh eviction and `auth:expired` Wails event emission, unified HTTP execution across JSON and multipart media requests via `executeRequest` and typed `APIError`, moved Arena WebSocket authentication to the `Authorization: Bearer` header, removed `HasActiveSession` and 14 duplicate BFF aliases, regenerated TypeScript bindings, synced `.routing/DOMAIN.md`, and added comprehensive `-race` unit tests in `app/app_auth_test.go`.

## Changes Made

### 1. Typed Errors, Dual-Token Persistence & Single-Flight Refresh (`app/app.go`)
- Defined `APIError{Status int, Code string, Message string}` implementing `error` (`Error() string`).
- Replaced `sessionToken` on `App` with `accessToken` (RAM-only), `refreshToken`, `sessionFilePath`, and `refreshMu sync.Mutex`.
- Session file path & format:
  - Defaults in `NewApp()` to `filepath.Join(cfgDir, "palabrario", "session.json")` using `os.UserConfigDir()` (falling back to `os.TempDir()` only if `os.UserConfigDir()` fails).
  - JSON format: `{"refresh_token":"<opaque-token>"}`.
  - Directory created with `0700`; file written atomically via `os.CreateTemp` + `os.Chmod(..., 0600)` + `os.Rename`.
  - Gracefully ignores missing or corrupt `session.json` on startup (`loadPersistedSession`); deletes `session.json` on `clearSession()`.
  - Added unexported `newAppWithConfig(coreURL, realtimeURL, imageURL, sessionFilePath string)` for isolated unit tests.
- Replaced `doRequest` internals with `executeRequest(ctx, buildReq, target, allowRetry)` + `doRequest(ctx, method, targetURL, body, target)`:
  - Injects `Authorization: Bearer <accessToken>` and `Accept: application/json`.
  - Single-flight refresh flow: on `401 Unauthorized` when `allowRetry && tokenWasSent && getRefreshToken() != ""`, `refreshSessionIfStale` acquires `refreshMu`, skips refresh if another goroutine already rotated `accessToken`, or calls `POST /api/v1/auth/refresh` once and retries the original request once (`allowRetry = false`).
  - On `403` with `Code == "ACCOUNT_SUSPENDED"`, or `401` on an authenticated call when refresh is absent/fails, calls `clearSession()` and emits Wails event `"auth:expired"` (`{"code": ..., "reason": ...}`) with safe `application.Get() != nil` guard.
  - Wraps transport failures as `&APIError{Status: 0, Code: "NETWORK_ERROR", Message: "Unable to reach the Palabrario server. Please check your connection."}`.
- Removed `HasActiveSession()`.

### 2. Auth, Governance & Media BFF Updates (`app/app_auth.go`, `app/app_governance.go`, `app/app_media.go`)
- `SignIn` and `SignUp` store both `authResp.Token` and `authResp.RefreshToken` via `setTokens` and return `*models.UserProfile`.
- `SignOut` clears RAM + deletes `session.json`, closes active Arena WS via `LeaveArenaRoom()`, and sends best-effort `POST /api/v1/auth/sign-out` with `{RefreshToken}` (ignoring network errors).
- `GetMe` returns `(nil, nil)` when no tokens exist, exchanges a persisted `refreshToken` on startup if `accessToken == ""`, returns `(nil, nil)` on `401`/`403` expired/suspended sessions, and propagates `*APIError` on transient `5xx`/network errors.
- `UpdateProfile` (`app_auth.go`) and `RedeemInvitationTicket` (`app_governance.go`) call `refreshSession` best-effort after success so `accessToken` immediately carries updated `AnchorDialect`/`Role` claims.
- Refactored `UploadImage` (`app_media.go`) to build multipart bytes once and delegate HTTP execution, Bearer injection, `401` single-flight refresh retry, `auth:expired` eviction, and `APIError` parsing to `executeRequest`, with English error messages.

### 3. Arena WebSocket Header Auth & Duplicate Alias Removal (`app/app_arena.go`, `app_auth.go`, `app_governance.go`)
- `FindMatchRoom` inspects `*APIError` (`errors.As(err, &apiErr) && apiErr.Status == http.StatusNotFound`) instead of bilingual substring matching.
- `ConnectArenaRoom` refreshes a persisted refresh token on startup if `accessToken == ""`, passes `http.Header{"Authorization": {"Bearer " + token}}` to `websocket.DefaultDialer.DialContext` (removing `?token=` query string), retries once on `401` handshake rejection if refresh succeeds, and translates `app_arena.go` errors to English.
- Removed all 14 unused duplicate aliases:
  - `app_auth.go`: `GetMyProfile`
  - `app_arena.go`: `CreateArenaMatch`, `GetMyMatchHistory`, `GetMatchReplayDetail`, `ConnectArena`, `LeaveArena`
  - `app_governance.go`: `ListPendingEditorRequests`, `GenerateInvitationTickets`, `ListUsersForGovernance`, `UpdateUserRole`, `UpdateUserStatus`, `CreateCommunityReport`, `ListCommunityReports`, `ResolveCommunityReport`

### 4. Bindings & Documentation Sync (`app/frontend/bindings/`, `.routing/DOMAIN.md`)
- Regenerated TypeScript bindings in `app/frontend/bindings/` via `wails3 task bindings` (`1 Service, 113 Methods, 6 Enums, 75 Models`).
- Updated `.routing/DOMAIN.md` Sections 3 & 4: replaced stale `local-app/` / `./local-app` references with `app/` (`app/app.go`, `app/app_auth.go`), documented dual-token RAM + `0600` `session.json` custody + single-flight refresh + `auth:expired` eviction, and updated the Mermaid WebSocket edge to `WS /ws/arena/{roomCode} (Authorization: Bearer {jwt})`.

### 5. Unit Tests (`app/app_auth_test.go`)
- Added `app/app_auth_test.go` using `net/http/httptest.Server`, `t.TempDir()`, and `newAppWithConfig`:
  - `TestSignInPersistsRefreshTokenAndSignOutClears`
  - `TestRestartRecoveryGetMeRefreshesPersistedSession`
  - `TestGetMeUnauthenticatedAndCorruptFileAndTransientError`
  - `TestSingleFlightAutoRefreshAcrossConcurrentRequestsAndUploadImage`
  - `TestSuspendedAccountAndFailedRefreshEvictSession`
  - `TestUpdateProfileAndRedeemTicketRotateAccessToken`
  - `TestConnectArenaRoomUsesAuthorizationBearerHeader`
  - `TestFindMatchRoomUsesAPIErrorNotFoundAndNetworkErrorSanitized`

## Verification

- `wails3 task bindings` (in `app/`): exit code `0` (`Processed: 264 Packages, 1 Service, 113 Methods, 6 Enums, 75 Models, 0 Events`).
- `go build .` (in `app/`): exit code `0`.
- `go test -race -v -count=1 ./...` (in `app/`): exit code `0` (`PASS ok palabrario/app 4.817s`).
