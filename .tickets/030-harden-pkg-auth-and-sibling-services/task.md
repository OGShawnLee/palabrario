# 030 — Harden pkg/auth, pkg/response, Sibling Services, and Flaky Tests

| Field | Value |
| :--- | :--- |
| **ID** | 030 |
| **Status** | completed |
| **Model** | Gemini 4 Hydrogen, `thinking_level: medium` (`gemini-4-hydrogen`) |
| **Role** | Backend Implementer (`implementation`) |
| **Spec** | [spec.md](spec.md) |

> **Specification Note:**  
> Second ticket of the end-to-end auth refactor ([028 report](../028-investigate-auth-workflow/report.md)). Eliminates hardcoded fallback `JWT_SECRET`s in `realtime-service` and `image-service`, standardizes `pkg/auth` (English errors, `TokenVersion` claim, `TokenType: "access"`, 15-minute access TTL + opaque refresh token helpers, bcrypt 72-byte check, shared `ExtractBearerToken` & `AllowedCORSOrigins`), adds machine-readable `code` to `pkg/response`, fixes CORS + `image-service` error leakage + `realtime-service` WebSocket origin/header hygiene, fixes the flaky/under-tested suites in `pkg/auth`, `realtime-service`, and `image-service`, and syncs `.routing/services/arena_service.md` & `.routing/services/media_service.md`.

## Scope

May create or edit:
- `pkg/response/json.go`
- `pkg/auth/auth.go`, `pkg/auth/auth_test.go` (and a helper file like `pkg/auth/http.go` if splitting HTTP/config helpers cleanly)
- `realtime-service/internal/config/config.go`, `realtime-service/internal/config/config_test.go`
- `realtime-service/internal/api/router.go`, `realtime-service/internal/api/handlers/arena_handler.go`, `realtime-service/internal/api/api_test.go`
- `realtime-service/internal/service/hub.go`, `realtime-service/internal/service/service_test.go`
- `image-service/internal/config/config.go`, `image-service/internal/config/config_test.go`
- `image-service/internal/api/router.go`, `image-service/internal/api/middleware/auth_middleware.go`, `image-service/internal/api/middleware/auth_middleware_test.go`
- `.routing/services/arena_service.md` (Section 2.7.1, 2.7.5, Part B banner)
- `.routing/services/media_service.md` (Part B banner, Section 2.8.2)
- `.tickets/030-harden-pkg-auth-and-sibling-services/changelog.md` (create)

Must not touch:
- `SCHEMA_FINAL.sql`, `sql/`, `core-service/`, `app/`, `.context/`
- `.routing/CONTEXT.md`, `.routing/DESIGN.md`, `.routing/PROGRESS.md`, `.routing/DOMAIN.md`, `.routing/API.md`, `.routing/ROUTES.md`, `.routing/INFRASTRUCTURE.md`, and any `.routing/services/*.md` other than `arena_service.md` and `media_service.md`
- Any file outside the above scope
- **No role (Router, Investigator, Implementer, Translator, or Reviewer) may commit, push, or merge without explicit user authorization.** Do not run `git commit`, `git push`, or `git merge`.

## Inputs

- [.tickets/028-investigate-auth-workflow/report.md](../028-investigate-auth-workflow/report.md): Findings `CORE-03`, `CORE-08`, `SIB-01`..`SIB-05`, and Flaky Test items 1, 5, 6, 7.
- [pkg/response/json.go](../../pkg/response/json.go)
- [pkg/auth/auth.go](../../pkg/auth/auth.go) and [pkg/auth/auth_test.go](../../pkg/auth/auth_test.go)
- [realtime-service/internal/config/config.go](../../realtime-service/internal/config/config.go)
- [realtime-service/internal/api/handlers/arena_handler.go](../../realtime-service/internal/api/handlers/arena_handler.go) and [realtime-service/internal/service/hub.go](../../realtime-service/internal/service/hub.go)
- [image-service/internal/config/config.go](../../image-service/internal/config/config.go) and [image-service/internal/api/middleware/auth_middleware.go](../../image-service/internal/api/middleware/auth_middleware.go)
- [.routing/services/arena_service.md](../../.routing/services/arena_service.md) and [.routing/services/media_service.md](../../.routing/services/media_service.md)

## Task

1. **Structured English Error Responses ([pkg/response/json.go](../../pkg/response/json.go)):**
   - Extend `ErrorBody` with `Code string `json:"code,omitempty"`` alongside `Error string `json:"error"``.
   - Keep `Error(w, status, message)` working for existing callers and add `ErrorWithCode(w, status int, code string, message string)`.

2. **Harden Shared Auth Primitives ([pkg/auth/auth.go](../../pkg/auth/auth.go)):**
   - Add `TokenVersion int `json:"token_version"`` to `Claims` (and keep backwards-compatible `GenerateToken` setting `TokenVersion: 1` or accept `tokenVersion` via `UserProfile` / new `GenerateAccessToken(user *models.UserProfile, tokenVersion int)` so existing callers keep compiling until Ticket `031`).
   - Change `NewTokenService(secretKey string, duration time.Duration)` default duration when `duration <= 0` to `15 * time.Minute` (short-lived access token).
   - Add opaque refresh token helpers: `GenerateRefreshToken() (raw string, hash string, err error)` (32 random bytes base64url-encoded + SHA-256 hex 64-char hash) and `HashRefreshToken(raw string) string`.
   - Enforce bcrypt max 72-byte guard in `HashPassword(password string)` returning a clear English error when `len([]byte(password)) > 72`.
   - Translate all `pkg/auth` error strings to English.
   - Export `ExtractBearerToken(r *http.Request) (token string, present bool, err error)` supporting case-insensitive `Bearer <token>` (`strings.EqualFold`) and rejecting malformed schemes or empty tokens.
   - Export `ParseCORSAllowedOrigins(raw string) []string` (if empty, returns dev localhost defaults `"http://localhost:*"`, `"http://127.0.0.1:*"`, `"wails://wails.localhost"`, `"http://wails.localhost"`, `"capacitor://localhost"` in non-production, or requires explicit origins).

3. **Eliminate Hardcoded Secrets & Fix Config Validation (`realtime-service` & `image-service`):**
   - In `realtime-service/internal/config/config.go` and `image-service/internal/config/config.go`:
     - Remove both hardcoded `JWT_SECRET` defaults completely. Require `JWT_SECRET` non-empty and `len(jwtSecret) >= 32` (matching `core-service`). Translate config errors to English.
     - Add `AccessTokenTTL time.Duration` parsed from `ACCESS_TOKEN_TTL_MINUTES` (optional positive int, default `15 * time.Minute`); reject non-positive/invalid `ACCESS_TOKEN_TTL_MINUTES` and invalid `JWT_EXPIRATION_HOURS` (fixing `image-service`'s silent ignore).
     - Add `AllowedOrigins []string` parsed from `CORS_ALLOWED_ORIGINS`.
   - Update `realtime-service/internal/config/config_test.go` and `image-service/internal/config/config_test.go`:
     - Replace `os.Unsetenv` with `t.Setenv(k, "")` so subtests stay isolated.
     - Test missing `JWT_SECRET`, short `JWT_SECRET` (`< 32`), invalid `ACCESS_TOKEN_TTL_MINUTES` / `JWT_EXPIRATION_HOURS`, and valid custom env vars.

4. **Harden Sibling Middleware, WebSockets & CORS (`realtime-service` & `image-service`):**
   - Set `AllowCredentials: false` and `AllowedOrigins: cfg.AllowedOrigins` (never wildcard `*` + credentials) in `realtime-service/internal/api/router.go` and `image-service/internal/api/router.go`.
   - In `image-service/internal/api/middleware/auth_middleware.go`: use `auth.ExtractBearerToken(r)` and return English `response.ErrorWithCode(w, http.StatusUnauthorized, ...)` **without** appending `err.Error()` to the HTTP body.
   - In `realtime-service/internal/api/handlers/arena_handler.go`: check `Authorization` header first via `auth.ExtractBearerToken(r)`, keeping `?token=` only as fallback when header is absent, and return English `ErrorWithCode` responses.
   - In `realtime-service/internal/service/hub.go`: replace `CheckOrigin: func(r *http.Request) bool { return true }` with an origin check injected/configured via `ArenaHub` / `upgrader` against `cfg.AllowedOrigins` (allowing empty `Origin` header for native Go desktop clients like Wails). Keep any room-full/error strings aligned with `arena_service.md`.

5. **Fix Flaky Tests & Expand Auth Coverage:**
   - In `pkg/auth/auth_test.go`: add tests for expired JWT, future `nbf`, `alg=none` / tampered signature, wrong secret, wrong issuer, `>72`-byte password rejection, `GenerateRefreshToken` + `HashRefreshToken` determinism, S2S internal token expiry/cross-use rejection, and `ExtractBearerToken` case-insensitive parsing.
   - In `image-service/internal/api/middleware/auth_middleware_test.go` and `realtime-service/internal/api/api_test.go`: add tests for expired JWT (`401`), S2S internal token rejected on user endpoints (`401`), lowercase `bearer <token>`, empty `Bearer ` token, and origin checks.
   - In `realtime-service/internal/service/service_test.go`: fix all 4 flakiness sources without weakening assertions:
     - Avoid `ReadJSON` loop corruption on `net.Error` timeouts (drive reads via a single dedicated reader goroutine per WS connection pushing `(WSEvent, error)` into a buffered channel, then `select` on the channel with deadline!).
     - Increase `ReconnectWindow` from `80ms` to `250ms` in `TestReconnectWithinWindow` (still well under real `10s`) and replace `time.Sleep(20ms)` + `time.Sleep(50ms)` cleanup sleep with channel/polling deadlines (`eventually` helper).
     - Replace `time.Sleep(5ms)` timing assumptions with deterministic assertions or clock/window bounds.

6. **Sync Service Documentation (`.routing/services/arena_service.md`, `.routing/services/media_service.md`):**
   - In [.routing/services/arena_service.md](../../.routing/services/arena_service.md): update the Part B banner to mention `response.ErrorWithCode` (`{ "code", "error" }` in English), update Sections `2.7.1` and `2.7.5` to state `Authorization: Bearer {jwt}` header is primary (with `?token=` fallback), document WebSocket `CheckOrigin` validation against `CORS_ALLOWED_ORIGINS` (allowing empty `Origin` for native desktop clients), and note `Claims.TokenVersion` + `15m` default `AccessTokenTTL`. Keep the room-full message note consistent with `hub.go`.
   - In [.routing/services/media_service.md](../../.routing/services/media_service.md): update the Part B banner (`response.ErrorWithCode` in English) and Section `2.8.2` to note `401` English `ErrorWithCode` responses on `POST /api/v1/images/upload` without leaking internal JWT errors in the body.

## Verification

Run every check in [spec.md](spec.md).

## Change log

Write `.tickets/030-harden-pkg-auth-and-sibling-services/changelog.md` recording every modified file, new helpers/codes in `pkg/auth` and `pkg/response`, flaky test fixes, routing doc updates, and exact `go test -race -count=1` outputs.
