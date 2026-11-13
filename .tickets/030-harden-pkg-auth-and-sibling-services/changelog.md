# 030 — Changelog: Harden `pkg/auth`, `pkg/response`, Sibling Services, and Flaky Tests

## Modified & Created Files

### Shared Packages (`pkg/`)
- **`pkg/response/json.go`** (modified):
  - Extended `ErrorBody` with `Code string `json:"code,omitempty"``.
  - Preserved `Error(w, status, message)` and added `ErrorWithCode(w, status, code, message)`.
  - Translated doc comments to English.
- **`pkg/auth/auth.go`** (modified):
  - Added `TokenVersion int `json:"token_version"`` and `TokenType string `json:"token_type,omitempty"`` (`TokenTypeAccess = "access"`) to `Claims`.
  - Changed `NewTokenService` default duration (`duration <= 0`) from `72 * time.Hour` to `DefaultAccessTokenTTL` (`15 * time.Minute`) and added `AccessTTL() time.Duration`.
  - Added `GenerateAccessToken(user *models.UserProfile, tokenVersion int)` and kept backwards-compatible `GenerateToken(user *models.UserProfile)` delegating with `tokenVersion = 1`.
  - Updated `VerifyToken` to reject non-`"access"` `TokenType` values when present.
  - Added `GenerateRefreshToken() (raw string, hash string, err error)` (32 random bytes base64url-encoded + SHA-256 hex digest) and `HashRefreshToken(raw string) string`.
  - Enforced `MaxBcryptPasswordLen = 72` byte guard in `HashPassword` and `CheckPassword`.
  - Translated all error messages in `pkg/auth` to English.
- **`pkg/auth/http.go`** (created):
  - Added `ExtractBearerToken(r *http.Request) (token string, present bool, err error)` with case-insensitive `Bearer` scheme matching and malformed/empty token rejection.
  - Added `DefaultDevCORSOrigins() []string`, `ParseCORSAllowedOrigins(raw string) []string`, and `IsOriginAllowed(origin string, allowedOrigins []string) bool` (supporting empty `Origin` for native desktop clients and port wildcards such as `http://localhost:*` / `http://127.0.0.1:*`).
- **`pkg/auth/auth_test.go`** (modified):
  - Added unit tests for `GenerateAccessToken` (`TokenVersion`, `TokenType`, `15m` default TTL), expired JWT, future `nbf`, `alg=none`, tampered signature, wrong secret, wrong issuer, non-access `token_type`, `>72`-byte password rejection (ASCII and multi-byte UTF-8), `GenerateRefreshToken` / `HashRefreshToken` determinism, S2S internal token expiry and cross-use rejection, `ExtractBearerToken` cases, and `ParseCORSAllowedOrigins` / `IsOriginAllowed`.

### `realtime-service/`
- **`realtime-service/internal/config/config.go`** (modified):
  - Removed hardcoded fallback `JWT_SECRET`; required `JWT_SECRET` non-empty and `len >= 32`.
  - Added `AccessTokenTTL time.Duration` (`ACCESS_TOKEN_TTL_MINUTES`, default `15m`, rejects non-positive/non-integer values).
  - Added `AllowedOrigins []string` (`CORS_ALLOWED_ORIGINS`, parsed via `auth.ParseCORSAllowedOrigins`; required in `production`).
  - Translated configuration errors to English.
- **`realtime-service/internal/config/config_test.go`** (modified):
  - Replaced `os.Unsetenv` with `t.Setenv(k, "")` for subtest isolation.
  - Added tests for missing `JWT_SECRET`, short `JWT_SECRET` (`< 32`), invalid `ACCESS_TOKEN_TTL_MINUTES` / `JWT_EXPIRATION_HOURS`, missing `CORS_ALLOWED_ORIGINS` in `production`, and valid custom env vars.
- **`realtime-service/internal/api/router.go`** (modified):
  - Set `AllowCredentials: false` and `AllowedOrigins` from `cfg.AllowedOrigins` (defaulting to dev origins when empty).
  - Wired `rt.hub.SetAllowedOrigins(allowedOrigins)`.
- **`realtime-service/internal/api/handlers/arena_handler.go`** (modified):
  - Checked `Authorization` header first via `auth.ExtractBearerToken(r)`, keeping `?token=` only as fallback when `Authorization` is absent.
  - Returned English `response.ErrorWithCode` responses (`MISSING_ROOM_CODE`, `INVALID_TOKEN_FORMAT`, `MISSING_TOKEN`, `INVALID_TOKEN`).
- **`realtime-service/internal/service/hub.go`** (modified):
  - Replaced unconditional `CheckOrigin: func(r *http.Request) bool { return true }` with `ArenaHub.SetAllowedOrigins` + `ArenaHub.CheckOrigin` validating non-empty `Origin` headers via `auth.IsOriginAllowed` (while permitting empty `Origin` for native desktop clients).
  - Returned English `response.ErrorWithCode(w, http.StatusNotFound, "ROOM_UNAVAILABLE", ...)` when room hydration fails.
- **`realtime-service/internal/api/api_test.go`** (modified):
  - Replaced `super_secret_*` test secret; added tests for expired JWT (`401`), S2S internal token rejection (`401`), empty `Bearer ` header with `?token=` fallback blocked (`401`), lowercase `bearer <token>`, CORS `AllowCredentials: false` + origin restriction, and WebSocket `CheckOrigin` (`403` on disallowed origin vs upgrade success on allowed/empty origin).
- **`realtime-service/internal/service/service_test.go`** (modified):
  - Fixed all 4 flakiness sources:
    1. Replaced repeated `ReadJSON` after read-deadline timeout with a single dedicated `wsReader` goroutine per WS connection feeding a buffered `chan wsReadResult` consumed via `select` + timer.
    2. Increased `ReconnectWindow` in `TestActiveRoom_ReconnectWithinWindow_And_HostSuccession` from `80ms` to `250ms` and replaced `time.Sleep(20ms)` with deterministic `ROOM_STATE` observation (`!p.IsReady`).
    3. Removed `time.Sleep(5ms)` timing assumptions before `SUBMIT_ANSWER` and asserted deterministic score (`[500, 1000]` per round / `[1000, 2000]` total) and elapsed-time window bounds (`[0, maxMs]`).
    4. Replaced post-`PODIUM` `time.Sleep(50ms)` with an `eventually` polling helper asserting `hub.ActiveRoomsCount() == 0`.

### `image-service/`
- **`image-service/internal/config/config.go`** (modified):
  - Removed hardcoded fallback `JWT_SECRET`; required `JWT_SECRET` non-empty and `len >= 32`.
  - Added `AccessTokenTTL time.Duration` (`ACCESS_TOKEN_TTL_MINUTES`, default `15m`) and rejected invalid `ACCESS_TOKEN_TTL_MINUTES` and `JWT_EXPIRATION_HOURS` (fixing silent fallback).
  - Added `AllowedOrigins []string` (`CORS_ALLOWED_ORIGINS`, required in `production`).
  - Translated configuration errors to English.
- **`image-service/internal/config/config_test.go`** (modified):
  - Replaced `os.Unsetenv` with `t.Setenv(k, "")` and added tests for missing/short `JWT_SECRET`, invalid `ACCESS_TOKEN_TTL_MINUTES` / `JWT_EXPIRATION_HOURS`, `production` missing `CORS_ALLOWED_ORIGINS`, and valid custom env vars.
- **`image-service/internal/api/router.go`** (modified):
  - Set `AllowCredentials: false` and `AllowedOrigins` from `cfg.AllowedOrigins` (defaulting to dev origins when empty).
  - Used `cfg.AccessTokenTTL` (with `cfg.JWTExpirationHours` fallback when unset) for `auth.NewTokenService`.
- **`image-service/internal/api/middleware/auth_middleware.go`** (modified):
  - Used `auth.ExtractBearerToken(r)` and returned English `response.ErrorWithCode` (`MISSING_TOKEN`, `INVALID_TOKEN_FORMAT`, `INVALID_TOKEN`) without appending `err.Error()` in `401` bodies.
- **`image-service/internal/api/middleware/auth_middleware_test.go`** (modified):
  - Replaced `super-secret-*` test secret and added tests for empty `Bearer `, expired JWT without internal error leak, S2S internal token rejection (`401`), and lowercase `bearer <token>`.

### Service Documentation (`.routing/services/`)
- **`.routing/services/arena_service.md`** (modified):
  - Updated Part B banner (`response.ErrorWithCode` `{ "code", "error" }` in English), Section `2.7.1` (`Authorization: Bearer` primary WS auth + `?token=` fallback, `CheckOrigin` against `CORS_ALLOWED_ORIGINS`, `Claims.TokenVersion`, `15m` default `AccessTokenTTL`, aligned room-full message `"La sala está llena (máximo 8 jugadores)"`), and Section `2.7.5` (`/ws/arena/{roomCode}` pre-upgrade `ErrorWithCode` statuses & origin validation).
- **`.routing/services/media_service.md`** (modified):
  - Updated Part B banner (`response.ErrorWithCode` `{ "code", "error" }` in English) and Section `2.8.2` (`POST /api/v1/images/upload` `401` English `ErrorWithCode` responses without internal error leakage + `CORS_ALLOWED_ORIGINS` with `AllowCredentials: false`).

## New Error Codes & Helpers
- **`pkg/response`**: `ErrorBody.Code` (`json:"code,omitempty"`), `ErrorWithCode(w, status, code, message)`
- **`pkg/auth`**: `Claims.TokenVersion`, `Claims.TokenType`, `DefaultAccessTokenTTL` (`15m`), `MaxBcryptPasswordLen` (`72`), `TokenTypeAccess` (`"access"`), `TokenService.AccessTTL()`, `TokenService.GenerateAccessToken(user, tokenVersion)`, `GenerateRefreshToken()`, `HashRefreshToken(raw)`, `ExtractBearerToken(r)`, `DefaultDevCORSOrigins()`, `ParseCORSAllowedOrigins(raw)`, `IsOriginAllowed(origin, allowedOrigins)`
- **HTTP Error Codes**: `MISSING_TOKEN`, `INVALID_TOKEN_FORMAT`, `INVALID_TOKEN`, `MISSING_ROOM_CODE`, `ROOM_UNAVAILABLE`

## Verification

### `go test -race -count=1 ./pkg/... ./realtime-service/... ./image-service/...`
```text
ok  	palabrario/pkg/auth	23.989s
?   	palabrario/pkg/models	[no test files]
?   	palabrario/pkg/response	[no test files]
?   	palabrario/pkg/utils	[no test files]
?   	palabrario/realtime-service	[no test files]
ok  	palabrario/realtime-service/internal/api	7.230s
?   	palabrario/realtime-service/internal/api/handlers	[no test files]
ok  	palabrario/realtime-service/internal/config	2.255s
ok  	palabrario/realtime-service/internal/service	9.611s
ok  	palabrario/image-service	11.490s
?   	palabrario/image-service/internal/api	[no test files]
?   	palabrario/image-service/internal/api/handlers	[no test files]
ok  	palabrario/image-service/internal/api/middleware	8.053s
ok  	palabrario/image-service/internal/config	5.948s
ok  	palabrario/image-service/internal/processor	4.813s
ok  	palabrario/image-service/internal/safety	2.303s
ok  	palabrario/image-service/internal/storage	2.342s
```

### `go test -race -count=3 ./pkg/... ./realtime-service/... ./image-service/...`
```text
ok  	palabrario/pkg/auth	40.858s
?   	palabrario/pkg/models	[no test files]
?   	palabrario/pkg/response	[no test files]
?   	palabrario/pkg/utils	[no test files]
?   	palabrario/realtime-service	[no test files]
ok  	palabrario/realtime-service/internal/api	6.739s
?   	palabrario/realtime-service/internal/api/handlers	[no test files]
ok  	palabrario/realtime-service/internal/config	2.691s
ok  	palabrario/realtime-service/internal/service	7.630s
ok  	palabrario/image-service	13.492s
?   	palabrario/image-service/internal/api	[no test files]
?   	palabrario/image-service/internal/api/handlers	[no test files]
ok  	palabrario/image-service/internal/api/middleware	3.039s
ok  	palabrario/image-service/internal/config	3.299s
ok  	palabrario/image-service/internal/processor	3.508s
ok  	palabrario/image-service/internal/safety	2.341s
ok  	palabrario/image-service/internal/storage	3.645s
```

### Secret & Scope Checks
- `git grep -n "super_secret\|super-secret" -- "realtime-service/*.go" "realtime-service/**/*.go" "image-service/*.go" "image-service/**/*.go"` returned no matches (exit code 1).
- `go build ./core-service/... ./realtime-service/... ./image-service/... ./pkg/...` exited with code `0`.
