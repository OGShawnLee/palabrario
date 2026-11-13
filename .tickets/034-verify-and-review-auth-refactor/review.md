# 034 — End-to-End Auth Refactor Verification & Security Review

| Field | Value |
| :--- | :--- |
| **Ticket** | `034-verify-and-review-auth-refactor` |
| **Verdict** | **APPROVED** |
| **Audit Scope** | Commits `ddb2530`, `db6c4ef`, `a870b27`, `5a178c7`, `42476b1` (`git diff 77d5a8c..42476b1`, Tickets `029`–`033`) |
| **Baseline** | [028 Investigation Report](../028-investigate-auth-workflow/report.md) |

---

## Executive Summary

All critical, high, medium, and architectural blockers identified in [028-investigate-auth-workflow/report.md](../028-investigate-auth-workflow/report.md) across Tickets `029`–`033` have been verified resolved. All security, business-logic, and contract parity rules pass static inspection and runtime execution (`go test -race` across root packages, `core-service` against `postgres-test :5433`, `app/`, and `wails3 task check`), with zero source files modified by the reviewer.

Three non-blocking baseline/deferred items are documented under [Non-Blocking Observations](#non-blocking-observations) (`DB-04` local dev seeds in `SCHEMA_FINAL.sql`, `SIB-05` deterministic SHA-256 S2S subkey derivation with separated `Issuer` claims, and the deferred Capacitor `live/http/` adapter prepared by `FE-05`).

---

## Findings Matrix (`028/report.md` Traceability)

### Layer 1: Database & `sqlc` (Ticket `029` — commit `ddb2530`)

| ID | Finding Summary | Status | Verification Evidence |
| :--- | :--- | :--- | :--- |
| **DB-01** | Auth queries omit `is_active`/`token_version`; public username lookup exposes suspended users | **PASS** | `sql/queries/users.sql` projects `a.is_active` and `a.token_version` on `GetAuthAccountStateByID`, `GetUserProfileByID`, `GetAccountWithProfileByEmail`, and `UpdateUserProfile`; `GetUserProfileByUsername` filters `AND a.is_active = true` (`sql/queries/users.sql:80-98`). Regenerated `core-service/internal/db/users.sql.go` matches. |
| **DB-02** | Table definition order (`accounts` before `regions`), missing `anchor_dialect` FK | **PASS** | `regions` is defined before `accounts` (`SCHEMA_FINAL.sql:59-66`); `users.anchor_dialect` enforces `VARCHAR(20) NOT NULL DEFAULT 'PAN' REFERENCES regions(code)` (`SCHEMA_FINAL.sql:100`). |
| **DB-03** | Missing case-insensitive unique email/username indexes, `users.name`/`user_name` DB constraints, `token_version`, and `refresh_tokens` table | **PASS** | `idx_accounts_email_lower` (`SCHEMA_FINAL.sql:82`) and `idx_users_user_name_lower` (`SCHEMA_FINAL.sql:113`) enforce case-insensitive uniqueness; `GetAccountByEmail` and `GetAccountWithProfileByEmail` query `LOWER(email)`; `users.name` checks `char_length(trim(name)) BETWEEN 2 AND 100`; `users.user_name` checks `^[a-z0-9][a-z0-9_-]{2,49}$`; `accounts.token_version INTEGER NOT NULL DEFAULT 1 CHECK (token_version >= 1)` added; `refresh_tokens` table created (`SCHEMA_FINAL.sql:84-93`) with `token_hash` unique index and revocation/expiry fields. |
| **DB-04** | Seed data coupled directly into `SCHEMA_FINAL.sql` | **PASS (Non-Blocking / Deferred)** | Explicitly scoped in `028/report.md` and `SCHEMA_FINAL.sql:584-610` as deterministic local dev/test seeds; pre-production split remains a deployment hardening task rather than an auth-refactor blocker. |

### Layer 2: `core-service` Auth, Middleware & Errors (Ticket `031` — commit `a870b27`)

| ID | Finding Summary | Status | Verification Evidence |
| :--- | :--- | :--- | :--- |
| **CORE-01** | Stateless JWT ignores `is_active` and `token_version` after suspension/revocation | **PASS** | Option B `RequireAuth` (`core-service/internal/api/middleware/auth_middleware.go:56-120`) queries live DB `GetAuthAccountStateByID` on all non-read HTTP methods (`POST`/`PUT`/`DELETE`/etc.), `/api/v1/auth/*`, and `/api/v1/governance/*`, rejecting inactive accounts (`403 ACCOUNT_SUSPENDED`) and stale token versions (`401 TOKEN_REVOKED`). Verified end-to-end by `TestAuth_OptionB_SuspendedAndRoleChangeEnforcement`. |
| **CORE-02** | Long-lived (24h) non-revocable JWTs without refresh rotation or live role sync | **PASS** | `pkg/auth/auth.go` sets `DefaultAccessTokenTTL = 15 * time.Minute`, embeds `TokenVersion` and `TokenType: "access"`; `POST /api/v1/auth/refresh` atomically rotates opaque SHA-256 hashed refresh tokens under `SELECT ... FOR UPDATE`; `ResolveEditorRequest(APPROVED)`, `RedeemInvitationTicket(EDITOR)`, and `UpdateUserRoleOrStatus` (role change or deactivation) increment `token_version`; Option B syncs `claims.Role = models.AccountRole(state.Role)` before `RequireRole`. |
| **CORE-03** | `AllowedOrigins: ["*"]` paired with `AllowCredentials: true` | **PASS** | `core-service/internal/api/router.go` configures `AllowCredentials: false` and loads allowed origins from `config.CORSAllowedOrigins` (`core-service/internal/config/config.go` enforces explicit origins in `ENVIRONMENT=production`). |
| **CORE-04** | `SignUp` missing format/dialect/password length validation | **PASS** | `core-service/internal/service/identity_service.go` validates RFC email (`mail.ParseAddress` + domain dot), normalizes `@`-prefixed lowercase `user_name` (`^[a-z0-9][a-z0-9_-]{2,49}$`), enforces 2–100 char `name` (defaulting to normalized `user_name`), validates 8-dialect allowlist, and enforces 8–72 UTF-8 password bytes. |
| **CORE-05** | `MapDBError` collapses unique/check/FK violations into generic `409 Conflict` | **PASS** | `core-service/internal/service/errors.go:57-75` maps SQLSTATE `22001` (`400 INPUT_TOO_LONG`), `23514` (`400 VALIDATION_ERROR`), `23503` (`400 INVALID_REFERENCE`), and `23505` (`EMAIL_ALREADY_EXISTS` vs `USERNAME_ALREADY_EXISTS` vs `CONFLICT` by `ConstraintName`). |
| **CORE-06** | `OptionalAuth` silently swallows invalid/expired `Authorization` headers | **PASS** | `OptionalAuth` (`core-service/internal/api/middleware/auth_middleware.go:125-150`) treats requests as anonymous only when `Authorization` is absent, returning `401` (`INVALID_TOKEN_FORMAT`, `TOKEN_EXPIRED`, `INVALID_TOKEN`) when a malformed or expired header is supplied. |
| **CORE-07** | Mixed Spanish/English errors without machine-readable `code` | **PASS** | `pkg/response/json.go` provides `ErrorWithCode(w, status, code, message)`; `pkg/auth`, `core-service/internal/api/middleware`, `identity_service.go`, and handlers emit English messages with explicit codes (`INVALID_CREDENTIALS`, `ACCOUNT_SUSPENDED`, `TOKEN_EXPIRED`, `TOKEN_REVOKED`, `INVALID_TOKEN_FORMAT`, `MISSING_TOKEN`, `FORBIDDEN`, etc.). |
| **CORE-08** | Silent bcrypt `> 72` byte truncation; `ChangePassword` did not invalidate sessions | **PASS** | `pkg/auth/auth.go` enforces `MaxBcryptPasswordLen = 72` in both `HashPassword` and `CheckPassword`; `ChangePassword` (`core-service/internal/service/identity_service.go:516-541`) runs inside a transaction updating the password hash, incrementing `token_version`, and revoking all active `refresh_tokens`. |
| **CORE-09** | `GetMe` returns `404` for deleted/missing user; `GetPublicProfileByUsername` swallows secondary query errors | **PASS** | `GetMe` (`identity_service.go:394-396`) maps `pgx.ErrNoRows` to `401 UNAUTHORIZED` and rejects inactive accounts with `403 ACCOUNT_SUSPENDED`; `GetPublicProfileByUsername` (`identity_service.go:565-678`) propagates unexpected DB errors rather than ignoring them. |
| **CORE-10** | `SignIn` performs 2 sequential DB queries | **PASS** | `SignIn` (`identity_service.go:217`) uses single-query `GetAccountWithProfileByEmail`. |

### Layer 3: `pkg/auth` & Sibling Services (Ticket `030` — commit `db6c4ef`)

| ID | Finding Summary | Status | Verification Evidence |
| :--- | :--- | :--- | :--- |
| **SIB-01** | Hardcoded `"super_secret_palabrario_dev_key"` / `"super-secret-dev-key"` `JWT_SECRET` fallbacks in `realtime-service` and `image-service` | **PASS** | Both `realtime-service/internal/config/config.go` and `image-service/internal/config/config.go` require `JWT_SECRET` (`len >= 32`) with zero fallback strings; verified via repo-wide static grep (`0` production matches). |
| **SIB-02** | `image-service` short secret bypass, `CORS` credentials, silent TTL fallback, leaked `err.Error()` in Spanish | **PASS** | `image-service/internal/config/config.go` requires `len(JWT_SECRET) >= 32`, rejects malformed `ACCESS_TOKEN_TTL_MINUTES` / `JWT_EXPIRATION_HOURS`, and enforces production `CORS_ALLOWED_ORIGINS`; `image-service/internal/api/router.go` sets `AllowCredentials: false`; `image-service/internal/api/middleware/auth_middleware.go:26-36` emits English `ErrorWithCode` responses without leaking `err.Error()`. |
| **SIB-03** | `realtime-service` WebSocket `?token=` exposure and permissive `CheckOrigin: return true` | **PASS** | `realtime-service/internal/api/handlers/arena_handler.go:37-54` extracts Bearer header via `auth.ExtractBearerToken`, rejects malformed `Authorization` headers even if `?token=` is present, and `hub.go:54-64,130` validates non-empty `Origin` headers against configured `CORS_ALLOWED_ORIGINS` via `auth.IsOriginAllowed` while permitting empty `Origin` for native Wails/Go WebSocket clients. |
| **SIB-04** | Duplicated JWT Bearer parsing across services | **PASS** | Shared `auth.ExtractBearerToken` in `pkg/auth/http.go` (case-insensitive `Bearer` prefix) is used across `core-service`, `realtime-service`, and `image-service`. |
| **SIB-05** | Same `JWT_SECRET` shared between user tokens and internal S2S tokens | **PASS (Acknowledged / Non-Blocking)** | `pkg/auth/auth.go` derives a domain-separated S2S HMAC subkey via SHA-256 and enforces distinct `Issuer` (`palabrario-core-service` vs `palabrario-internal-s2s`) and `TokenType: "access"` claims so user and S2S tokens are never interchangeable. |

### Layer 4: Wails Desktop Go BFF (Ticket `032` — commit `5a178c7`)

| ID | Finding Summary | Status | Verification Evidence |
| :--- | :--- | :--- | :--- |
| **BFF-01 & BFF-02** | Single JWT persisted on disk, missing refresh/expiry lifecycle, unguarded `401`/`403` handling, `HasActiveSession` leaking state | **PASS** | `app/app.go` keeps `accessToken` strictly in RAM (`sync.RWMutex`) and persists `refreshToken` atomically to `session.json` (`0700` dir, temp file + `0600` + `os.Rename`). `executeRequest` implements single-flight `401` refresh (`refreshMu`, `refreshSessionIfStale`, single retry `allowRetry=false`), wipes RAM + `session.json` and emits `"auth:expired"` on `403 ACCOUNT_SUSPENDED` or failed refresh. `HasActiveSession` is completely removed. |
| **BFF-03** | `UpdateProfile` / `RedeemInvitationTicket` leave stale claims in JWT | **PASS** | `app/app_auth.go:100-102` (`UpdateProfile`) and `app/app_governance.go:21-23` (`RedeemInvitationTicket`) call `refreshSession` best-effort so in-memory `accessToken` reflects updated role/dialect claims immediately. |
| **BFF-04** | Lossy string-concatenated HTTP errors; fragile `strings.Contains(..., "404")` | **PASS** | Structured `APIError{Status, Code, Message}` in `app/app.go` preserves HTTP status & error code; network errors wrap with `NETWORK_ERROR`; `FindMatchRoom` in `app/app_arena.go` uses `errors.As(err, &apiErr) && apiErr.Status == http.StatusNotFound`. |
| **BFF-05** | Duplicated HTTP client logic (`doRequest` vs `UploadImage`), 14 duplicate alias methods, WebSocket `?token=` query auth, zero BFF tests | **PASS** | `doRequest` and `UploadImage` share `executeRequest`; 14 duplicate alias methods removed; `app/app_arena.go:140` sends `Authorization: Bearer <token>` header with zero `?token=`; `app/app_auth_test.go` provides 8 `-race` tests covering persistence (`0600`), single-flight refresh, `403` eviction + `auth:expired`, `SignOut`, `UpdateProfile`, and `RedeemInvitationTicket`. |

### Layer 5: Vue 3 Frontend & Router Guards (Ticket `033` — commit `42476b1`)

| ID | Finding Summary | Status | Verification Evidence |
| :--- | :--- | :--- | :--- |
| **FE-01** | `MockAuthAPI` vs `WailsAuthAdapter` divergent behavior (no password check, suspended users accepted, silent error strings) | **PASS** | `IAuthAPI` (`app/frontend/src/api/contract/auth.contract.ts`) standardizes `signIn`/`signUp` to `Promise<UserProfile>` (throwing on errors) and `getMe()` to `Promise<UserProfile \| null>`. `MockAuthAPI` (`app/frontend/src/api/mock/mock.adapter.ts`) validates credentials (`password123`/stored password), rejects suspended accounts (`mateo_vargas` / `suspended@palabrario.app`), hides suspended users from `getPublicProfileByUsername` (`null`), and mirrors backend `signUp` validation. |
| **FE-02** | `AuthResponse` token leak in contract, redundant `hasActiveSession()` double IPC call, missing `auth:expired` listener | **PASS** | `AuthResponse` and `hasActiveSession()` removed from `IAuthAPI`, `WailsAuthAdapter`, `MockAuthAPI`, and `useAuth.ts`; `loadSession()` executes a single `api.auth.getMe()` call; `useAuth` registers `onSessionExpired` once on init and exposes `sessionNotice` / `clearSessionNotice`. |
| **FE-03** | Unprotected `/studio`, `/settings`, `/arena/history`; `init()` skipped on `/auth/*` | **PASS** | `app/frontend/src/router/index.ts` awaits `init()` in `beforeEach` on every route, sets `meta: { requiresAuth: true }` on `/studio`, `/settings`, `/arena/history`, sets `meta: { guestOnly: true }` on `/auth/*`, and watches `currentUser` + `sessionNotice` to redirect off protected routes on live expiry. Redundant `init()` call removed from `TopBar.vue`. |
| **FE-04** | Incomplete frontend `SignUp` validation | **PASS** | `app/frontend/src/views/auth/SignUpView.vue` validates normalized username regex, email format, 8–72 UTF-8 byte password limit, and optional display name length (2–100 chars). |
| **FE-05** | Web/Capacitor `live/http/` adapter readiness | **PASS (Contract Prepared)** | `IAuthAPI` / `useAuth` decoupled from host token storage (`AuthResponse` removed) as specified for this refactor; Capacitor `live/http/` implementation itself remains scoped to future mobile target work. |
| **FE-06** | `ProfileView` missing `currentUser` watcher / error distinction; `ArenaView` unauthenticated room creation | **PASS** | `app/frontend/src/views/ProfileView.vue` watches both `route.params.username` and `currentUser?.id_user` with stale-request guarding and separates `404` from retryable errors; `ArenaView.vue` & `CreateMatchDialog.vue` gate room creation behind `isAuthenticated` and redirect guests to `/auth/sign-in?redirect=/arena`. |

---

## Security, Business Rule & Spec Checklist

| Requirement | Result | Notes |
| :--- | :--- | :--- |
| Zero hardcoded `JWT_SECRET` fallbacks (`core-service`, `realtime-service`, `image-service`) | **PASS** | Verified via static grep across non-test `.go` files; all three `config.go` files enforce `len >= 32`. |
| Zero `AllowedOrigins: ["*"]` + `AllowCredentials: true`; WS `Origin` validation in `hub.go` | **PASS** | Static grep confirms zero `AllowCredentials: true`; all 3 routers use `AllowCredentials: false`; `hub.go` validates non-empty `Origin` via `auth.IsOriginAllowed`. |
| `core-service` Option B (`15m` access JWT on non-governance `GET`s; live DB `is_active`, `token_version`, `role` sync on writes, `/auth/*`, `/governance/*`; `OptionalAuth` `401` on invalid headers) | **PASS** | Verified in `core-service/internal/api/middleware/auth_middleware.go` and `TestAuth_OptionB_SuspendedAndRoleChangeEnforcement`. |
| `GetUserProfileByUsername` hides suspended accounts (`404`) | **PASS** | `sql/queries/users.sql` (`AND a.is_active = true`) + `identity_service.go` `ErrNotFound`. |
| `SignIn`/`SignUp`/`RefreshToken`/`SignOut` rotate & revoke hashed `refresh_tokens`; `ChangePassword` & `is_active = false` bump `token_version` & revoke all refresh tokens | **PASS** | Verified in `identity_service.go`, `governance_service.go`, and `core-service` tests. |
| Wails BFF (`app/`): `0600` disk refresh token, RAM-only access JWT, single-flight `401` refresh, `auth:expired` + disk/RAM clear on `403 ACCOUNT_SUSPENDED` / failed refresh, WS `Authorization` header, zero duplicate alias methods | **PASS** | Verified in `app/app.go`, `app/app_auth.go`, `app/app_governance.go`, `app/app_arena.go`, and `app/app_auth_test.go`. |
| Frontend (`app/frontend/src/`): `IAuthAPI`/`MockAuthAPI`/`WailsAuthAdapter` parity, `useAuth` single `getMe()`, `onSessionExpired`, `router.beforeEach` `await init()` on all routes, guards on `/studio`, `/settings`, `/arena/history`, `/auth/*` | **PASS** | Verified across contract, adapters, composable, router, and auth/profile/arena views. |
| English auth/validation/middleware error responses with machine-readable `code`; zero Spanish auth/middleware errors in `pkg/auth`, `core-service` middleware & `identity_service.go`, `image-service` auth/config, `realtime-service` auth/config, and `app/` | **PASS** | Verified via static inspection and `Select-String` sweep. |
| Reviewer read-only compliance (only `034/review.md` created; zero source modifications; no unauthorized git ops) | **PASS** | Confirmed via `git status --short` after all checks completed. |

---

## Verification Commands & Exit Codes

All commands executed against commit `42476b1` after the Ticket `029`–`033` changes:

1. **Shared `pkg/`, `realtime-service/`, and `image-service/` race suite (`3x`)** — repo root
   - Command: `go test -race -count=3 ./pkg/... ./realtime-service/... ./image-service/...`
   - Exit code: **`0`**
   - Summary output:
     - `ok palabrario/pkg/auth 39.002s`
     - `ok palabrario/realtime-service/internal/api 6.367s`
     - `ok palabrario/realtime-service/internal/config 2.386s`
     - `ok palabrario/realtime-service/internal/service 7.782s`
     - `ok palabrario/image-service 16.745s`
     - `ok palabrario/image-service/internal/api/middleware 3.479s`
     - `ok palabrario/image-service/internal/config 2.321s`
     - `ok palabrario/image-service/internal/processor 2.470s`
     - `ok palabrario/image-service/internal/safety 3.130s`
     - `ok palabrario/image-service/internal/storage 3.097s`

2. **Wails BFF Go race suite** — `app/`
   - Command: `go test -race -count=1 ./...`
   - Exit code: **`0`**
   - Summary output: `ok palabrario/app 5.400s` (all 8 BFF auth/refresh/eviction/header tests passing)

3. **Core-service Postgres-backed race suite (`127.0.0.1:5433`)** — repo root
   - Pre-check: `Test-NetConnection -ComputerName 127.0.0.1 -Port 5433 -InformationLevel Quiet` -> `True` (`postgres-test` already running prior to verification; left running).
   - Command: `go test -race -count=1 ./core-service/...`
   - Exit code: **`0`**
   - Summary output:
     - `ok palabrario/core-service/internal/api 47.574s`
     - `ok palabrario/core-service/internal/service 55.828s`

4. **Wails bindings generation, Vue type-check, Vite production build, and Go binary build** — `app/`
   - Command: `wails3 task check`
   - Exit code: **`0`**
   - Summary output:
     - `wails3 generate bindings -clean=true -ts`: `264 Packages, 1 Service, 113 Methods, 6 Enums, 75 Models, 0 Events`
     - `bun run --cwd frontend type-check` (`vue-tsc --noEmit`): `0` errors
     - `bun run --cwd frontend build` (`vite build --mode production`): `✓ 2500 modules transformed`, `✓ built in 7.77s`
     - `go build .`: `0` errors

5. **Static Security & Localization Grep Suite** — repo root
   - Command: PowerShell `Select-String` checks for hardcoded `JWT_SECRET`s, `AllowCredentials: true`, `?token=` in `app/app_arena.go`, `hasActiveSession`, `AuthResponse` in `app/frontend/src/`, and Spanish error strings in `pkg/auth`, `core-service/internal/api/middleware`, `core-service/internal/service/identity_service.go`, `image-service/internal/{api/middleware,config}`, `realtime-service/internal/{api,config}`, and `app/`.
   - Exit code: **`0`**
   - Result: Zero hardcoded secrets, zero `AllowCredentials: true`, zero `?token=` query auth in `app/app_arena.go`, zero `HasActiveSession`/`hasActiveSession`, zero `AuthResponse` in frontend `src/`, and zero Spanish auth/middleware error strings.

6. **Working Tree Cleanliness Check** — repo root
   - Command: `git status --short`
   - Exit code: **`0`**
   - Output prior to creating `review.md`: `M .routing/PROGRESS.md` and `?? .tickets/034-verify-and-review-auth-refactor/` (zero modified source files).

---

## Non-Blocking Observations

1. **`DB-04` Local Dev Seeds in `SCHEMA_FINAL.sql` (`SCHEMA_FINAL.sql:584-610`):** Deterministic dev accounts (`admin@palabrario.org`, `silas@palabrario.org`, `gael@palabrario.org`) remain in `SCHEMA_FINAL.sql`. This is expected for local development/test DB initialization and does not block this refactor, though it should be separated before production deployment.
2. **Pre-Existing Non-Auth Spanish Domain Strings:** `image-service/internal/api/handlers/image_handler.go` / `safety/validator.go` (upload validation messages verified by `main_test.go`) and `realtime-service/internal/service/room.go:240,245` (`"La partida ya está en progreso"`, `"La sala está llena (máximo 8 jugadores)"`, matching `arena_service.md`) retain non-auth Spanish domain messages. All auth, token, session, config, and middleware error responses are English + machine-readable `code`.
3. **Future Web/Capacitor `live/http/` Adapter (`FE-05`):** The `IAuthAPI` contract and `useAuth` composable are cleanly decoupled from desktop-specific session storage so `live/http/` can be added without modifying view/guard logic when Capacitor is implemented.
