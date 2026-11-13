# 028 — End-to-End Auth & Session Workflow Audit Report

**Scope:** Read-only full-stack audit of Palabrario's authentication, JWT/session lifecycle, authorization middleware, error handling, test flakiness/coverage, redundant abstractions, and client contracts across PostgreSQL/sqlc, `pkg/auth`, `core-service`, `realtime-service`, `image-service`, the Wails Go BFF (`app/`), and the Vue frontend (`app/frontend/src/`).  
**Status:** Complete (zero repository source files modified).

---

## Executive Summary

Palabrario's core cryptographic primitives (`bcrypt` default cost, `HS256` JWTs with `sha256`-derived `palabrario-internal-s2s` sub-key separation) and isolated Postgres template test harness (`CREATE DATABASE ... TEMPLATE`) provide a solid baseline, **but several cross-layer gaps break security revocation, session consistency, and adapter parity**:

1. **Suspended users (`is_active = false`) and role/dialect changes are ignored after JWT issuance.** `SignIn` checks `account.IsActive`, yet `RequireAuth`, `GetMe`, `UpdateProfile`, `ChangePassword`, `realtime-service`, and `image-service` trust stateless JWT claims for up to **72 hours** with no DB active-status check, token version, or token refresh mechanism.
2. **Hardcoded fallback `JWT_SECRET`s in sibling services.** Unlike `core-service` (which fails fast), `realtime-service` and `image-service` silently fall back to hardcoded source-code secrets (`image-service` also skips the 32-char check and ignores `JWT_EXPIRATION_HOURS` parse errors).
3. **`SignUp` lacks validation symmetry and leaks 500s on column overflow.** Password minimum length (`>= 8`), `anchor_dialect` allowlist validation, name length (`2..100`), email format, and username format are missing in `SignUp` (unlike `UpdateProfile`/`ChangePassword`), causing PostgreSQL `22001` truncation errors to surface as `500 Internal Server Error`, while `23514` check violations map to `409 Conflict`.
4. **Wails BFF retains dead JWTs on `401` and loses sessions on restart.** `App.doRequest` never clears `sessionToken` on `401 Unauthorized`, `HasActiveSession()` only checks `sessionToken != ""`, sessions live only in Go RAM, and error HTTP status codes are discarded (forcing Spanish backend errors raw into an English UI).
5. **Frontend `MockAuthAPI` vs `WailsAuthAdapter` contract divergence & missing route guards.** `MockAuthAPI.signIn` returns `null` (and ignores passwords) while `WailsAuthAdapter.signIn` throws; `getMe()` returns `null` in Mock vs throws in Wails; `IAuthAPI` forces `AuthResponse.token` to be fabricated as `""`; `router/index.ts` has no navigation guards on `/studio` or `/settings`.
6. **Zero auth tests in `app/` and frontend, zero HTTP tests for suspended live tokens / expired JWTs / role staleness, and sleep-based WebSocket timing windows in `realtime-service`.**

---

## Findings by Layer & Severity

### Layer 1: Database & sqlc (`SCHEMA_FINAL.sql`, `sql/queries/`)

| ID | Severity | Finding | Location |
| :--- | :--- | :--- | :--- |
| DB-01 | **High** | `GetUserProfileByID`, `GetUserProfileByUsername`, and `UpdateUserProfile` do not check or project `a.is_active` | [users.sql](../../sql/queries/users.sql#L26-L90) |
| DB-02 | **High** | `users.anchor_dialect` has no FK to `regions(code)` or `CHECK` constraint | [SCHEMA_FINAL.sql](../../SCHEMA_FINAL.sql#L74), [SCHEMA_FINAL.sql](../../SCHEMA_FINAL.sql#L130-L137) |
| DB-03 | **Medium** | Asymmetric case sensitivity & missing SQL `CHECK` constraints on `email`, `user_name`, and `name` | [SCHEMA_FINAL.sql](../../SCHEMA_FINAL.sql#L61-L72), [users.sql](../../sql/queries/users.sql#L11-L14), [users.sql](../../sql/queries/users.sql#L45-L62) |
| DB-04 | **Low** | Demo seed accounts use active status with hardcoded `password123` bcrypt hash | [SCHEMA_FINAL.sql](../../SCHEMA_FINAL.sql#L583-L594) |

- **DB-01 (`High`) — Suspended accounts invisible to profile & read queries:**
  - **Evidence:** [GetUserProfileByID](../../sql/queries/users.sql#L26-L43), [GetUserProfileByUsername](../../sql/queries/users.sql#L45-L62), and [UpdateUserProfile](../../sql/queries/users.sql#L64-L90) join `accounts a ON a.id_account = u.id_user` without `AND a.is_active = true` and without selecting `a.is_active`. `GetAccountByID` ([users.sql](../../sql/queries/users.sql#L16-L19)) selects `is_active`, but `ChangePassword` never inspects it.
  - **Failure Mode / UX Impact:** Once suspended by an Admin, a user can still call `/api/v1/auth/me`, `/api/v1/users/me`, `/api/v1/users/me/password`, and their public profile continues resolving as normal with `is_owner: true`.
  - **Fix:** Select `a.is_active` (or add `GetAuthAccountStateByID` returning `is_active`, `role`, `token_version`/`updated_at`) and enforce `is_active = true` in authenticated identity queries and/or auth middleware.

- **DB-02 (`High`) — Unconstrained `users.anchor_dialect` allows invalid dialect codes:**
  - **Evidence:** [SCHEMA_FINAL.sql](../../SCHEMA_FINAL.sql#L74) declares `anchor_dialect VARCHAR(10) NOT NULL DEFAULT 'MX'` with no foreign key to [regions(code)](../../SCHEMA_FINAL.sql#L132) and no `CHECK` constraint.
  - **Failure Mode / UX Impact:** `SignUp` does not validate `anchor_dialect` in Go, so arbitrary values like `"ZZ"` or `"INVALID"` persist into Postgres and break dialect-scoped queries/UI region labels.
  - **Fix:** Add `REFERENCES regions(code)` (or `CHECK (anchor_dialect IN ('MX','ES','AR','CO','CL','CU','CAM','PAN'))` if seed ordering prevents FK, though `regions` is seeded before `users` in `SCHEMA_FINAL.sql`).

- **DB-03 (`Medium`) — Missing DB constraints & inconsistent case-insensitive uniqueness:**
  - **Evidence:** [accounts.email](../../SCHEMA_FINAL.sql#L61) (`VARCHAR(255) NOT NULL UNIQUE`) and [users.user_name](../../SCHEMA_FINAL.sql#L72) (`VARCHAR(50) NOT NULL UNIQUE`) rely on Go `strings.ToLower` rather than `CITEXT` or functional `UNIQUE (LOWER(...))` indexes; [GetUserProfileByUsername](../../sql/queries/users.sql#L62) filters with `LOWER(u.user_name) = LOWER($1)` whereas [GetAccountByEmail](../../sql/queries/users.sql#L14) uses exact `email = $1`. No `CHECK (char_length(name) >= 2)` or username character constraint exists at the DB layer.
  - **Failure Mode / UX Impact:** Any future seed/script or SQL insertion with mixed-case email/username creates unreachable logins or duplicate usernames under `LOWER(user_name)` lookups.
  - **Fix:** Add `CHECK (char_length(trim(name)) BETWEEN 2 AND 100)`, `CHECK (user_name ~ '^[a-z0-9][a-z0-9_-]{2,49}$')`, and `CREATE UNIQUE INDEX` on `LOWER(email)` / `LOWER(user_name)`.

- **DB-04 (`Low`) — Shared demo seed credentials in schema DDL:**
  - **Evidence:** [SCHEMA_FINAL.sql](../../SCHEMA_FINAL.sql#L583-L594) inserts `admin@palabrario.org`, `silas@palabrario.org`, and `gael@palabrario.org` with the same `password123` bcrypt hash.
  - **Failure Mode / UX Impact:** Safe for local dev/template test DB, dangerous if `SCHEMA_FINAL.sql` is applied verbatim to a staging/production database.
  - **Fix:** Split `SCHEMA_FINAL.sql` DDL from dev seed data (`sql/seeds/dev.sql`) before production deployment.

---

### Layer 2: `pkg/auth` & `core-service`

| ID | Severity | Finding | Location |
| :--- | :--- | :--- | :--- |
| CORE-01 | **Critical** | `RequireAuth` never verifies account `is_active`, `role`, or revocation (`token_version`) against DB | [auth_middleware.go](../../core-service/internal/api/middleware/auth_middleware.go#L22-L45), [identity_service.go](../../core-service/internal/service/identity_service.go#L150-L265) |
| CORE-02 | **Critical** | Role promotion/demotion & dialect updates remain stale in JWT claims for 72h; no token refresh/re-issue | [auth.go](../../pkg/auth/auth.go#L45-L74), [auth_middleware.go](../../core-service/internal/api/middleware/auth_middleware.go#L65-L82), [governance_service.go](../../core-service/internal/service/governance_service.go#L180-L188), [governance_service.go](../../core-service/internal/service/governance_service.go#L280-L314), [governance_service.go](../../core-service/internal/service/governance_service.go#L462-L469) |
| CORE-03 | **High** | `AllowedOrigins: ["*"]` paired with `AllowCredentials: true` across all services | [router.go](../../core-service/internal/api/router.go#L39-L46) |
| CORE-04 | **High** | `SignUp` omits password length (`>= 8`), `anchor_dialect` allowlist, name bounds, email & username format checks | [identity_service.go](../../core-service/internal/service/identity_service.go#L40-L82) vs [identity_service.go](../../core-service/internal/service/identity_service.go#L178-L241) |
| CORE-05 | **High** | PostgreSQL `22001` truncation returns `500`, `23514` maps to `409`, and `SignUp` conflates email vs username conflicts | [errors.go](../../core-service/internal/service/errors.go#L51-L63), [identity_service.go](../../core-service/internal/service/identity_service.go#L66-L82) |
| CORE-06 | **Medium** | `OptionalAuth` silently swallows malformed/expired Bearer tokens instead of rejecting invalid credentials | [auth_middleware.go](../../core-service/internal/api/middleware/auth_middleware.go#L49-L63) |
| CORE-07 | **Medium** | All backend auth/middleware errors are Spanish despite an English-only UI mandate | [auth_middleware.go](../../core-service/internal/api/middleware/auth_middleware.go#L26-L38), [identity_service.go](../../core-service/internal/service/identity_service.go#L47-L270), [errors.go](../../core-service/internal/service/errors.go#L14-L21) |
| CORE-08 | **Medium** | `ChangePassword` does not revoke existing JWTs; `bcrypt` silently truncates passwords > 72 bytes | [identity_service.go](../../core-service/internal/service/identity_service.go#L230-L265), [auth.go](../../pkg/auth/auth.go#L158-L169) |
| CORE-09 | **Medium** | `GetPublicProfileByUsername` swallows secondary query errors silently; `GetMe` returns `404` on missing user | [identity_service.go](../../core-service/internal/service/identity_service.go#L156-L159), [identity_service.go](../../core-service/internal/service/identity_service.go#L286-L393) |
| CORE-10 | **Low** | `SignIn` executes 2 round-trip queries (`GetAccountByEmail` + `GetMe`) on every login | [identity_service.go](../../core-service/internal/service/identity_service.go#L121-L137) |

- **CORE-01 (`Critical`) — Suspended accounts retain full API access for up to 72 hours:**
  - **Evidence:** `SignIn` checks `!account.IsActive` ([identity_service.go](../../core-service/internal/service/identity_service.go#L130-L132)), and Admins can suspend users via `UpdateUserRoleOrStatus` ([governance_service.go](../../core-service/internal/service/governance_service.go#L285-L294)). However, `RequireAuth` ([auth_middleware.go](../../core-service/internal/api/middleware/auth_middleware.go#L22-L45)) is purely stateless, and `GetMe`, `UpdateProfile`, and `ChangePassword` ([identity_service.go](../../core-service/internal/service/identity_service.go#L150-L265)) never check `is_active`.
  - **Failure Mode / UX Impact:** A suspended user continues calling authenticated endpoints, editing their profile, creating Arena rooms, and uploading images until their 72h token naturally expires.
  - **Fix:** Add an `accounts.token_version` (incremented on suspension, role change, or password change) and/or check `is_active` in `RequireAuth` (via lightweight indexed PK query or short-TTL cache) and return `401`/`403`.

- **CORE-02 (`Critical`) — Role & `anchor_dialect` staleness in 72h JWT:**
  - **Evidence:** `RequireRole` checks `claims.Role` from the JWT ([auth_middleware.go](../../core-service/internal/api/middleware/auth_middleware.go#L74)). When `ResolveEditorRequest` ([governance_service.go](../../core-service/internal/service/governance_service.go#L180-L188)), `RedeemInvitationTicket` ([governance_service.go](../../core-service/internal/service/governance_service.go#L462-L469)), or `UpdateUserRoleOrStatus` ([governance_service.go](../../core-service/internal/service/governance_service.go#L280-L294)) updates `accounts.role`, no new token is issued and the old JWT stays valid. Similarly, `UpdateProfile` updates `users.anchor_dialect` without re-issuing a JWT containing the new `AnchorDialect` claim (which `/words/search`, `/words/word-of-the-day`, and `realtime-service` read from JWT claims).
  - **Failure Mode / UX Impact:** Demoted editors/admins keep moderation/editorial endpoints unlocked for 72h; promoted learners get `403 Forbidden` until signing out and back in; users who change region in Settings keep seeing old-region search ranking.
  - **Fix:** Add `POST /api/v1/auth/refresh` (or re-issue `AuthResponse` from `RedeemInvitationTicket` & `UpdateProfile`), and hydrate authoritative `Role`/`IsActive` from DB (or invalidate via `token_version` claim) in `RequireAuth`.

- **CORE-03 (`High`) — Insecure/invalid CORS wildcard + credentials:**
  - **Evidence:** [core-service/internal/api/router.go](../../core-service/internal/api/router.go#L39-L46), [realtime-service/internal/api/router.go](../../realtime-service/internal/api/router.go#L43-L49), and [image-service/internal/api/router.go](../../image-service/internal/api/router.go#L82-L88) all set `AllowedOrigins: []string{"*"}` with `AllowCredentials: true`.
  - **Failure Mode / UX Impact:** In browser/Capacitor WebView deployments, `Access-Control-Allow-Origin: *` combined with `Access-Control-Allow-Credentials: true` violates the Fetch CORS spec (browsers reject credentialed responses or middlewares dynamically reflect arbitrary `Origin` headers). Palabrario uses `Authorization: Bearer` headers, not cookies.
  - **Fix:** Set `AllowCredentials: false` and load `CORS_ALLOWED_ORIGINS` from config (strict allowlist in production; localhost/wails origins in dev).

- **CORE-04 (`High`) — Incomplete `SignUp` validation vs `UpdateProfile`/`ChangePassword`/Frontend:**
  - **Evidence:** [SignUp](../../core-service/internal/service/identity_service.go#L40-L51) only checks non-empty strings. In contrast, `ChangePassword` requires `len(req.NewPassword) >= 8` ([identity_service.go](../../core-service/internal/service/identity_service.go#L239)), `UpdateProfile` enforces `len(name) in [2,100]` and the 8-dialect allowlist ([identity_service.go](../../core-service/internal/service/identity_service.go#L185-L196)), and `SignUpView.vue` hints `"Letters, numbers, or hyphens"` and `minlength="3"`.
  - **Failure Mode / UX Impact:** Direct API callers (or buggy clients) can register with a 1-char password (`"a"`), invalid email (`"foo"`), `@handle` prefix, 1-char name, or invalid dialect (`"ZZ"`).
  - **Fix:** Extract shared validators in `identity_service.go` (or `pkg/models`): email format (`net/mail` or regex), password `8..72` bytes, trimmed `user_name` (strip optional `@`, match `^[a-z0-9][a-z0-9_-]{2,49}$`), `name` length `2..100`, and `anchor_dialect` in the 8-region allowlist.

- **CORE-05 (`High`) — Missing `pgconn.PgError` code/constraint differentiation in `MapDBError`:**
  - **Evidence:** [errors.go](../../core-service/internal/service/errors.go#L51-L63) maps `23505` using a single `defaultMsg` without checking `pgErr.ConstraintName`, maps `23514` (`check_violation`) to `ErrConflict` (`409`), and ignores `22001` (`string_data_right_truncation`) and `23503` (`foreign_key_violation`), which fall through to `500 Internal Server Error` ([errors.go](../../core-service/internal/service/errors.go#L88-L90)). Meanwhile `SignUp` passes the exact same generic message `"El correo o nombre de usuario ya está registrado"` for both `CreateAccount` and `CreateUserProfile` ([identity_service.go](../../core-service/internal/service/identity_service.go#L71), [identity_service.go](../../core-service/internal/service/identity_service.go#L81)).
  - **Failure Mode / UX Impact:** Entering a 51-char username or 101-char name returns `500` `"No se pudo registrar la cuenta"`; duplicate email vs duplicate username cannot be distinguished in the UI.
  - **Fix:** Handle `22001` -> `ErrInvalidInput` (`400`), `23514` -> `ErrValidation` (`400`), inspect `pgErr.ConstraintName` on `23505` (`accounts_email_key` vs `users_user_name_key`), and return distinct messages from `CreateAccount` vs `CreateUserProfile`.

- **CORE-06 (`Medium`) — `OptionalAuth` silently masks expired/invalid Bearer tokens:**
  - **Evidence:** [auth_middleware.go](../../core-service/internal/api/middleware/auth_middleware.go#L49-L63) calls `VerifyToken` and ignores errors, letting requests to `/users/u/{username}`, `/words/search`, `/words/word-of-the-day`, `/decks/{id}`, `/collections/{id}` proceed anonymously even when an `Authorization` header was explicitly sent and failed verification.
  - **Failure Mode / UX Impact:** An expired desktop/mobile session viewing its own profile (`/u/my_name`) silently gets `is_owner: false` (private counts/email hidden) instead of receiving `401` to trigger session cleanup/re-login.
  - **Fix:** If no `Authorization` header is present, continue anonymously; if an `Authorization` header IS present and malformed/expired, return `401 Unauthorized` (so clients know their session expired).

- **CORE-07 (`Medium`) — Error language mismatch (Spanish backend vs English UI):**
  - **Evidence:** All error strings in [auth_middleware.go](../../core-service/internal/api/middleware/auth_middleware.go#L26-L102), [identity_service.go](../../core-service/internal/service/identity_service.go#L47-L270), and [errors.go](../../core-service/internal/service/errors.go#L14-L21) are hardcoded in Spanish, while [CONTEXT.md](../../.routing/CONTEXT.md#L31) specifies an English UI and [SignInView.vue](../../app/frontend/src/views/auth/SignInView.vue#L23) / [SignUpView.vue](../../app/frontend/src/views/auth/SignUpView.vue#L57) render `err.message` verbatim.
  - **Failure Mode / UX Impact:** English UI displays alerts like `"Credenciales incorrectas"` or `"El correo o nombre de usuario ya está registrado"`.
  - **Fix:** Return stable machine-readable error codes (`{ "code": "INVALID_CREDENTIALS", "error": "..." }`) and/or standardize user-facing API messages in English (mapped cleanly in frontend error normalizer).

- **CORE-08 (`Medium`) — Password change doesn't invalidate live tokens; bcrypt 72-byte silent truncation:**
  - **Evidence:** [ChangePassword](../../core-service/internal/service/identity_service.go#L230-L265) updates `password_hash` without bumping a token version/revocation timestamp; [HashPassword](../../pkg/auth/auth.go#L158-L164) passes arbitrary-length passwords to `bcrypt.GenerateFromPassword`, which silently truncates beyond 72 bytes (and `golang.org/x/crypto/bcrypt` returns `ErrPasswordTooLong` in newer versions or truncates depending on version — either way unchecked in validation).
  - **Failure Mode / UX Impact:** Compromised tokens stay valid after a user changes their password; passphrases >72 bytes produce unexpected `500` or weaken verification.
  - **Fix:** Validate `len([]byte(password)) <= 72` in `SignUp` and `ChangePassword`; bump `accounts.token_version` on `ChangePassword`.

- **CORE-09 (`Medium`) — Silent partial degradation & misleading `404` on `/auth/me`:**
  - **Evidence:** `GetMe` maps `pgx.ErrNoRows` via `MapDBError` to `ErrNotFound` (`404`) ([identity_service.go](../../core-service/internal/service/identity_service.go#L156-L159)). `GetPublicProfileByUsername` wraps 6 secondary queries with `if err == nil` ([identity_service.go](../../core-service/internal/service/identity_service.go#L286-L393)).
  - **Failure Mode / UX Impact:** A deleted account with an unexpired JWT causes `404` on `/auth/me` instead of `401` (preventing standard `401` auto-sign-out). Transient DB failures on badges/decks/due counts silently return empty arrays/`null` with HTTP `200`.
  - **Fix:** Return `ErrUnauthorized` from `GetMe` when the token's user row no longer exists or is inactive; propagate (or explicitly log/surface) DB errors in `GetPublicProfileByUsername`.

- **CORE-10 (`Low`) — Extra DB query on `SignIn`:**
  - **Evidence:** `SignIn` calls `GetAccountByEmail` ([identity_service.go](../../core-service/internal/service/identity_service.go#L121)) and then calls `s.GetMe` (`GetUserProfileByID`, [identity_service.go](../../core-service/internal/service/identity_service.go#L134)).
  - **Fix:** Single `GetAuthProfileByEmail` query returning `password_hash`, `is_active`, and full `UserProfile` columns.

---

### Layer 3: Sibling Services (`realtime-service`, `image-service`)

| ID | Severity | Finding | Location |
| :--- | :--- | :--- | :--- |
| SIB-01 | **Critical** | `realtime-service` and `image-service` fall back to hardcoded `JWT_SECRET`s instead of failing fast | [realtime config.go](../../realtime-service/internal/config/config.go#L39-L45), [image config.go](../../image-service/internal/config/config.go#L41-L44) vs [core config.go](../../core-service/internal/config/config.go#L41-L60) |
| SIB-02 | **High** | `image-service` skips `len(JWT_SECRET) >= 32` check, ignores invalid `JWT_EXPIRATION_HOURS`, and leaks `err.Error()` on `401` | [image config.go](../../image-service/internal/config/config.go#L41-L70), [image auth_middleware.go](../../image-service/internal/api/middleware/auth_middleware.go#L26-L40) |
| SIB-03 | **High** | WebSocket JWT passed in URL query string (`?token=...`) and `CheckOrigin` returns `true` unconditionally | [arena_handler.go](../../realtime-service/internal/api/handlers/arena_handler.go#L36-L42), [hub.go](../../realtime-service/internal/service/hub.go#L18-L24), [app_arena.go](../../app/app_arena.go#L149-L151) |
| SIB-04 | **Medium** | `Bearer` header parsing is case-sensitive in `realtime-service` and `image-service`, unlike `core-service` | [arena_handler.go](../../realtime-service/internal/api/handlers/arena_handler.go#L39), [image auth_middleware.go](../../image-service/internal/api/middleware/auth_middleware.go#L26) vs [core auth_middleware.go](../../core-service/internal/api/middleware/auth_middleware.go#L31) |
| SIB-05 | **Low** | S2S secret derived deterministically from user `JWT_SECRET` (`sha256(secret + ":palabrario-s2s-internal")`) | [auth.go](../../pkg/auth/auth.go#L47-L51) |

- **SIB-01 (`Critical`) — Hardcoded default `JWT_SECRET`s permit token forgery & environment drift:**
  - **Evidence:** [realtime-service/internal/config/config.go](../../realtime-service/internal/config/config.go#L40-L42) defaults to `"super_secret_jwt_key_palabrario_change_in_production_min_32_chars"`; [image-service/internal/config/config.go](../../image-service/internal/config/config.go#L42-L44) defaults to `"super-secret-jwt-key-palabrario-2026"`. Only `core-service` enforces `requireEnv("JWT_SECRET")`.
  - **Failure Mode / UX Impact:** Any environment starting `realtime-service` or `image-service` without `JWT_SECRET` silently trusts publicly committed keys (allowing anyone to forge admin/user JWTs against WebSocket Arena and image uploads) while rejecting real tokens signed by `core-service`.
  - **Fix:** Require `JWT_SECRET` (`>= 32` chars) in `Load()` across all three services (or fail fast whenever `ENVIRONMENT != "development"` and require non-empty in all).

- **SIB-02 (`High`) — `image-service` weak config validation & raw JWT error leakage:**
  - **Evidence:** [image-service/internal/config/config.go](../../image-service/internal/config/config.go#L41-L70) does not validate `len(jwtSecret) >= 32` and silently falls back to `72h` if `JWT_EXPIRATION_HOURS` is invalid (unlike `core-service` and `realtime-service`). [image-service/internal/api/middleware/auth_middleware.go](../../image-service/internal/api/middleware/auth_middleware.go#L39) returns `"token inválido o expirado: " + err.Error()`.
  - **Failure Mode / UX Impact:** Weak signing keys allowed in `image-service`; internal JWT parsing details returned in HTTP 401 bodies.
  - **Fix:** Enforce `len(jwtSecret) >= 32`, reject non-positive/non-integer `JWT_EXPIRATION_HOURS`, and share a canonical `pkg/auth` HTTP middleware helper without appending `err.Error()`.

- **SIB-03 (`High`) — WebSocket JWT transmitted in query string (`?token=...`) & open `CheckOrigin`:**
  - **Evidence:** [app/app_arena.go](../../app/app_arena.go#L149-L151) appends `?token=<JWT>` to `wsURL` even though Go's `websocket.DefaultDialer.DialContext` supports custom `Authorization: Bearer <token>` headers and [arena_handler.go](../../realtime-service/internal/api/handlers/arena_handler.go#L38-L41) already checks `Authorization`. Meanwhile [hub.go](../../realtime-service/internal/service/hub.go#L21-L23) returns `true` for all origins.
  - **Failure Mode / UX Impact:** JWTs in URL query strings leak into `chimiddleware.Logger` logs, reverse proxies, and browser history; unconstrained `CheckOrigin` enables cross-site WebSocket hijacking in web contexts.
  - **Fix:** For Wails desktop (`app_arena.go`), pass `http.Header{"Authorization": {"Bearer " + token}}` instead of `?token=`. For future browser/Capacitor JS WebSocket clients (which cannot set `Authorization` headers on `new WebSocket()`), exchange the Bearer token for a single-use short-lived (15–30s) opaque WS ticket or authenticate via first WebSocket message, and restrict `CheckOrigin` in non-dev environments.

- **SIB-04 (`Medium`) — Inconsistent `Bearer` prefix parsing:**
  - **Evidence:** `core-service` uses `strings.EqualFold(parts[0], "Bearer")` ([core auth_middleware.go](../../core-service/internal/api/middleware/auth_middleware.go#L31)), whereas `realtime-service` ([arena_handler.go](../../realtime-service/internal/api/handlers/arena_handler.go#L39)) and `image-service` ([image auth_middleware.go](../../image-service/internal/api/middleware/auth_middleware.go#L26)) use case-sensitive `strings.HasPrefix(authHeader, "Bearer ")`.
  - **Fix:** Consolidate HTTP Bearer extraction into `pkg/auth`.

---

### Layer 4: Wails Go BFF (`app/`)

| ID | Severity | Finding | Location |
| :--- | :--- | :--- | :--- |
| BFF-01 | **Critical** | `doRequest` & `UploadImage` never clear `sessionToken` on `401 Unauthorized`; `HasActiveSession` reports stale `true` | [app.go](../../app/app.go#L74-L149), [app_media.go](../../app/app_media.go#L61-L83) |
| BFF-02 | **High** | `sessionToken` is RAM-only (no persistence across desktop restarts) | [app.go](../../app/app.go#L28), [app_auth.go](../../app/app_auth.go#L15-L39) |
| BFF-03 | **High** | `RedeemInvitationTicket` and `UpdateProfile` leave stale `Role` / `AnchorDialect` in `sessionToken` | [app_governance.go](../../app/app_governance.go#L16-L22), [app_auth.go](../../app/app_auth.go#L70-L76) |
| BFF-04 | **Medium** | `doRequest` discards HTTP status codes and wraps raw Go network errors in Spanish | [app.go](../../app/app.go#L112-L129), [app_arena.go](../../app/app_arena.go#L56) |
| BFF-05 | **Medium** | Redundant BFF wrapper aliases & duplicated HTTP/auth code in `app_media.go` | [app_auth.go](../../app/app_auth.go#L50-L53), [app_arena.go](../../app/app_arena.go#L22-L248), [app_governance.go](../../app/app_governance.go#L45-L114), [app_media.go](../../app/app_media.go#L53-L84) |

- **BFF-01 (`Critical`) — Stale `sessionToken` on `401` desynchronizes `App` and frontend:**
  - **Evidence:** [doRequest](../../app/app.go#L118-L129) checks `if resp.StatusCode >= 400` and returns `errors.New(errResp.Error)` without clearing `a.sessionToken` when `resp.StatusCode == http.StatusUnauthorized`. `UploadImage` ([app_media.go](../../app/app_media.go#L72-L83)) duplicates the same omission. `HasActiveSession()` ([app.go](../../app/app.go#L147-L149)) only checks `a.getSessionToken() != ""`.
  - **Failure Mode / UX Impact:** Once a token expires, is invalidated, or belongs to a suspended user whose request returns `401`, `HasActiveSession()` keeps returning `true` while every protected call fails; `useAuth` swallows `getMe()` errors without informing Go to clear the dead token, and Wails never emits an `auth:unauthorized` event.
  - **Fix:** In `doRequest` / shared HTTP helper, when `resp.StatusCode == http.StatusUnauthorized` and a `sessionToken` was sent, clear `a.setSessionToken("")` and emit a Wails event (`auth:expired`) so `useAuth` immediately resets `currentUser` and redirects if on a protected route.

- **BFF-02 (`High`) — Ephemeral in-memory desktop session:**
  - **Evidence:** [App.sessionToken](../../app/app.go#L28) is only stored in Go memory (`sync.RWMutex`), never persisted to disk/OS keychain, contrasting directly with `MockAuthAPI` persisting `localStorage` (`palabrario_mock_session`).
  - **Failure Mode / UX Impact:** Desktop users are silently logged out every time they close and reopen Palabrario, whereas 72h tokens imply persistent sessions.
  - **Fix:** Document whether RAM-only is intentional or persist the token securely via OS keyring / encrypted app-data file (and validate via `GetMe()` on startup).

- **BFF-03 (`High`) — Token claims stay stale after `RedeemInvitationTicket` / `UpdateProfile`:**
  - **Evidence:** [RedeemInvitationTicket](../../app/app_governance.go#L16-L22) returns an updated `UserProfile` (`Role: "EDITOR"`) while `a.sessionToken` retains the prior `LEARNER` JWT; [UpdateProfile](../../app/app_auth.go#L70-L76) similarly leaves old `AnchorDialect` in `a.sessionToken`.
  - **Fix:** Return refreshed `AuthResponse` from backend on role/dialect change (or call `RefreshToken` in BFF) and update `a.setSessionToken(newToken)`.

- **BFF-04 (`Medium`) — HTTP status codes stripped & raw network errors leaked:**
  - **Evidence:** [app.go](../../app/app.go#L114,L125-L128) converts all HTTP errors into untyped `errors.New(errResp.Error)` and wraps transport errors as `fmt.Errorf("error de red al conectar con el servidor: %w", err)`. Because status codes are lost, [app_arena.go](../../app/app_arena.go#L56) resorts to fragile bilingual substring matching (`strings.Contains(msg, "not found") || strings.Contains(msg, "no encontrada")`).
  - **Fix:** Introduce a structured `APIError{Status int, Code string, Message string}` in `app/` and sanitize connection errors before sending to the UI.

- **BFF-05 (`Medium`) — Duplicate Wails exported bindings & duplicated HTTP client code:**
  - **Evidence:** `GetMyProfile` ([app_auth.go](../../app/app_auth.go#L50-L53)), `CreateArenaMatch`, `GetMyMatchHistory`, `GetMatchReplayDetail`, `ConnectArena`, `LeaveArena` ([app_arena.go](../../app/app_arena.go#L22-L247)), and `ListPendingEditorRequests`, `GenerateInvitationTickets`, `ListUsersForGovernance` ([app_governance.go](../../app/app_governance.go#L45-L114)) + `UpdateUserRole`, `UpdateUserStatus`, `CreateCommunityReport`, `ListCommunityReports`, `ResolveCommunityReport` duplicate canonical methods and bloat `app/frontend/bindings/palabrario/app/app.ts`. `UploadImage` ([app_media.go](../../app/app_media.go#L53-L84)) duplicates request creation, Bearer injection, and JSON error parsing from `doRequest`.
  - **Fix:** Remove unused aliases from `App` (regenerate Wails bindings) and extract `executeRequest(req, target)` shared by `doRequest` and `UploadImage`.

---

### Layer 5: Frontend (`app/frontend/src/`)

| ID | Severity | Finding | Location |
| :--- | :--- | :--- | :--- |
| FE-01 | **Critical** | `MockAuthAPI` vs `WailsAuthAdapter` diverge on `signIn` (`null` vs throw, no password check), `getMe` (`null` vs throw), and `signUp` validation | [mock.adapter.ts](../../app/frontend/src/api/mock/mock.adapter.ts#L110-L170) vs [auth.wails.ts](../../app/frontend/src/api/live/wails/auth.wails.ts#L22-L36) |
| FE-02 | **High** | `IAuthAPI` leaks `AuthResponse.token` (fabricated `""` in Wails, ignored in `useAuth`) & forces redundant `hasActiveSession()` | [auth.contract.ts](../../app/frontend/src/api/contract/auth.contract.ts#L9-L16), [auth.wails.ts](../../app/frontend/src/api/live/wails/auth.wails.ts#L13-L15), [useAuth.ts](../../app/frontend/src/composables/useAuth.ts#L10-L65) |
| FE-03 | **High** | No Vue Router navigation guards (`beforeEach`) on `/studio`, `/settings`, or `/auth/*`; `init()` skipped on `/auth/*` entry | [router/index.ts](../../app/frontend/src/router/index.ts#L81-L92), [TopBar.vue](../../app/frontend/src/components/TopBar.vue#L16-L33) |
| FE-04 | **High** | `SignUpView` allows invalid usernames (`@`, spaces, >50 chars) despite hint, and unvalidated Step 2 submit | [SignUpView.vue](../../app/frontend/src/views/auth/SignUpView.vue#L23-L76) |
| FE-05 | **High** | Capacitor/HTTP adapter documented in `.routing/API.md` & deps installed, but missing from `api/config.ts` and `api/live/` | [api/config.ts](../../app/frontend/src/api/config.ts#L1-L43), [package.json](../../app/frontend/package.json#L14-L16), [API.md](../../.routing/API.md#L46-L48) |
| FE-06 | **Medium** | `ProfileView` swallows all errors into `"Profile not found."` and doesn't react to auth state hydration; `ArenaView`/`CreateMatchDialog` lack auth gating | [ProfileView.vue](../../app/frontend/src/views/ProfileView.vue#L18-L37), [ArenaView.vue](../../app/frontend/src/views/ArenaView.vue#L56-L65), [CreateMatchDialog.vue](../../app/frontend/src/components/arena/CreateMatchDialog.vue#L48-L80) |

- **FE-01 (`Critical`) — Behavioral contract mismatch between Mock and Live adapters:**
  - **Evidence:**
    - `MockAuthAPI.signIn` ([mock.adapter.ts](../../app/frontend/src/api/mock/mock.adapter.ts#L110-L117)) **accepts any non-empty password** for known mock users, has no suspended-account handling, and returns `null` on failure (`useAuth.ts:41-42` translates `null` to `"Invalid email or password."`). Meanwhile `WailsAuthAdapter.signIn` ([auth.wails.ts](../../app/frontend/src/api/live/wails/auth.wails.ts#L22-L24)) **rejects with an Error** containing Spanish backend text (`"Credenciales incorrectas"` or `"Esta cuenta ha sido suspendida..."`).
    - `MockAuthAPI.getMe()` ([mock.adapter.ts](../../app/frontend/src/api/mock/mock.adapter.ts#L166-L169)) returns `null` when unauthenticated, whereas `WailsAuthAdapter.getMe()` ([auth.wails.ts](../../app/frontend/src/api/live/wails/auth.wails.ts#L34-L36)) throws `401` if called without an active token.
    - `MockAuthAPI.signUp` ([mock.adapter.ts](../../app/frontend/src/api/mock/mock.adapter.ts#L121-L136)) strips `@`, defaults blank `name` to `userName`, checks uniqueness, and throws English errors; `core-service` does not strip `@` and throws Spanish errors.
  - **Failure Mode / UX Impact:** UI code tested against Mock behaves differently in Wails (and vice versa); suspended/wrong-password/error paths cannot be tested accurately in Mock mode.
  - **Fix:** Align `IAuthAPI` on a single deterministic contract: `signIn`/`signUp` return `Promise<UserProfile>` and throw normalized errors; `getMe()` returns `Promise<UserProfile | null>` (returning `null` on `401` / missing session in both Mock and Wails). Enforce mock passwords and `is_active` checks in `MockAuthAPI`.

- **FE-02 (`High`) — Leaky `AuthResponse` & redundant `hasActiveSession()` in `IAuthAPI`:**
  - **Evidence:** [auth.contract.ts](../../app/frontend/src/api/contract/auth.contract.ts#L9-L16) forces `signIn`/`signUp` to return `Promise<AuthResponse | null>` and exposes `hasActiveSession()`. Because Wails keeps the JWT in Go RAM, [auth.wails.ts](../../app/frontend/src/api/live/wails/auth.wails.ts#L13-L15) fabricates `new AuthResponse({ token: "", user })`. [useAuth.ts](../../app/frontend/src/composables/useAuth.ts#L44,L59) discards `res.token` and makes two sequential IPC/network round-trips on startup (`hasActiveSession()` then `getMe()`, [useAuth.ts](../../app/frontend/src/composables/useAuth.ts#L12-L13)).
  - **Fix:** Remove `AuthResponse` and `hasActiveSession()` from `IAuthAPI`. Have `loadSession()` call `api.auth.getMe()` directly (one IPC/HTTP call; token storage stays encapsulated inside each transport adapter — Go RAM/keyring for Wails, `@capacitor/preferences` for future HTTP adapter).

- **FE-03 (`High`) — Missing Vue Router guards & TopBar-coupled `init()`:**
  - **Evidence:** [router/index.ts](../../app/frontend/src/router/index.ts#L81-L92) defines `/studio`, `/settings`, `/auth/sign-in`, `/auth/sign-up` with zero `meta.requiresAuth` or `beforeEach` guards. Protected-route redirect only happens inside `TopBar.handleSignOut()` ([TopBar.vue](../../app/frontend/src/components/TopBar.vue#L30-L32)). Moreover, `init()` is only called in [TopBar.vue](../../app/frontend/src/components/TopBar.vue#L16-L18), which is NOT mounted inside `AuthLayout`.
  - **Failure Mode / UX Impact:** Unauthenticated users can directly navigate to `#/studio` or `#/settings`; logged-in users landing directly on `#/auth/sign-in` aren't redirected away and `currentUser` stays `null` until hitting an `AppLayout` route.
  - **Fix:** Call `await init()` inside a router `beforeEach` guard (`router/index.ts`), add `meta: { requiresAuth: true }` on `/studio`, `/settings`, `/arena/history` (if private), and `meta: { guestOnly: true }` on `/auth/*`.

- **FE-04 (`High`) — `SignUpView` form validation gaps:**
  - **Evidence:** [SignUpView.vue](../../app/frontend/src/views/auth/SignUpView.vue#L23-L25) only validates `form.user_name.trim().length >= 3 && form.email.trim().includes("@") && form.password.length >= 8`, ignoring the UI hint `"Letters, numbers, or hyphens"` ([SignUpView.vue](../../app/frontend/src/views/auth/SignUpView.vue#L72)), allowing spaces, `@`, special chars, >50-char usernames, and >72-byte passwords.
  - **Fix:** Match backend validation rules cleanly in `isStep1Valid` (e.g., `/^[a-zA-Z0-9_-]{3,50}$/` after stripping optional leading `@`, RFC-like email regex, `8..72` password bytes, and `2..100` name if provided).

- **FE-05 (`High`) — Missing `live/http/` Capacitor adapter implementation despite `.routing/API.md` & installed deps:**
  - **Evidence:** [package.json](../../app/frontend/package.json#L14-L16) includes `@capacitor/core` and `@capacitor/preferences`, and [API.md](../../.routing/API.md#L46-L48) specifies `APIMode = "auto" | "wails" | "http" | "mock"`, yet [api/config.ts](../../app/frontend/src/api/config.ts#L1,L38-L42) only supports `"wails" | "mock"` and falls back to `"mock"` in all non-Wails browsers/Capacitor builds.
  - **Fix:** Prepare `IAuthAPI` without Wails-specific token leaks (`AuthResponse`) now so the upcoming `live/http/` adapter can store the JWT in `@capacitor/preferences` (`SecureStorage` / `Preferences`) transparently.

- **FE-06 (`Medium`) — `ProfileView` & `ArenaView` auth UX gaps:**
  - **Evidence:** [ProfileView.vue](../../app/frontend/src/views/ProfileView.vue#L27-L37) watches only `route.params.username`, not `currentUser`, and catches all errors as `"Profile not found."` (masking 500/network errors or expired `OptionalAuth` failures). [ArenaView.vue](../../app/frontend/src/views/ArenaView.vue#L56-L58) and [CreateMatchDialog.vue](../../app/frontend/src/components/arena/CreateMatchDialog.vue#L61-L80) allow anonymous users to attempt `prepareMatch` (hitting a raw backend `401` in live mode while succeeding in mock mode).
  - **Fix:** Re-fetch profile if `currentUser?.id_user` changes while on `/u/:username`, distinguish network errors from `404`, and gate `Create Room` / `prepareMatch` behind `isAuthenticated` (prompting Sign In when anonymous).

---

## Flaky Auth Tests & Coverage Gaps

### Live `go test` Execution Results (Read-Only Environment Snapshot)

| Command | Result | Summary / Observed Output |
| :--- | :--- | :--- |
| `go test -v -count=1 ./pkg/auth/...` | **PASS** | `TestPasswordHashingAndJWT` passed (`0.24s`, total `5.032s`) |
| `go test -v -count=1 ./image-service/...` | **PASS** | All 6 packages passed (`image-service`, `internal/api/middleware`, `internal/config`, `internal/processor`, `internal/safety`, `internal/storage`) |
| `go test -v -count=1 ./realtime-service/...` | **PASS** | `internal/api`, `internal/config`, `internal/service` passed (`api` `4.324s`, `config` `3.222s`, `service` `4.941s`) |
| `go test -v -count=1 -run TestAuth ./core-service/internal/api/...` | **FAIL** | All 4 `TestAuth_*` tests failed at `setupTestServer` / `SetupTestDB`: `postgres-test (:5433)` not running (`connectex: No connection could be made because the target machine actively refused it`) |
| `go test -v -count=1 -run TestIdentityService ./core-service/internal/service/...` | **FAIL** | All 3 `TestIdentityService_*` tests failed at `SetupTestDB`: `postgres-test (:5433)` not running |
| `go test -v -count=1 ./...` (in `app/`) | **PASS (no tests)** | `? palabrario/app [no test files]` |

### Audit of Test Files & Flaky Patterns

*(Note: `task.md` references `core-service/internal/api/governance_reports_test.go`; in the repository this is split into [core-service/internal/api/governance_reports_api_test.go](../../core-service/internal/api/governance_reports_api_test.go) and [core-service/internal/service/governance_reports_test.go](../../core-service/internal/service/governance_reports_test.go).)*

1. **`realtime-service/internal/service/service_test.go` — Tight millisecond sleeps, `ReconnectWindow` race, & `gorilla/websocket` read-deadline corruption (High Flakiness Risk):**
   - **Wall-clock sleep & tight window:** [service_test.go](../../realtime-service/internal/service/service_test.go#L255-L256,L282,L294) sets `hub.ReconnectWindow = 80 * time.Millisecond` and calls `time.Sleep(20 * time.Millisecond)` before dialing a fresh TCP/WebSocket connection ([line 296](../../realtime-service/internal/service/service_test.go#L296)). Under Windows/CI CPU load, dialing + JWT verification + `ROOM_STATE` broadcast can exceed the remaining `60ms`, triggering premature disconnect cleanup and failing the reconnection assertion.
   - **Post-PODIUM cleanup sleep:** [service_test.go](../../realtime-service/internal/service/service_test.go#L623-L625) uses `time.Sleep(50 * time.Millisecond)` before asserting `hub.ActiveRoomsCount() == 0`, instead of polling with `Eventually`/deadline.
   - **Corrupted WebSocket reader on poll timeout:** [waitForEventType](../../realtime-service/internal/service/service_test.go#L43-L57) and [waitForRoomStateCondition](../../realtime-service/internal/service/service_test.go#L59-L79) set a `300ms` read deadline in a `for` loop and ignore errors (`if err := ws.ReadJSON(&ev); err == nil`). In `gorilla/websocket`, once `ReadJSON` hits a `net.Error` timeout, the underlying connection reader enters a permanently failed state ("repeated read on failed websocket connection"), so any temporary stall guarantees a timeout failure on all subsequent loop iterations.
   - **Timing-dependent `TotalTimeSpentMs > 0`:** [service_test.go](../../realtime-service/internal/service/service_test.go#L502,L532,L604) sleeps `5ms` hoping `TotalTimeSpentMs > 0`, which is fragile on coarse OS timers if room state timing/clock jitter collapses below 1ms.

2. **`core-service/internal/service/service_test.go` — `time.Now().UnixNano()` collision risk in tight loop (Medium Flakiness Risk):**
   - **Evidence:** [service_test.go](../../core-service/internal/service/service_test.go#L622-L624) generates 8 synthetic users in a tight loop using `fmt.Sprintf("p%d_%d@test.org", i, time.Now().UnixNano())` and `fmt.Sprintf("p%d_test_%d", i, time.Now().UnixNano())` twice per iteration. (Mitigated slightly by prefix `p%d_`, but the pattern is non-deterministic and error-prone if copied without `i`).
   - **Fix:** Use `uuid.NewString()` or deterministic index strings inside isolated template DBs.

3. **`core-service/internal/api/auth_test.go` & `governance_reports_api_test.go` — Shared `http.DefaultClient` & swallowed JSON decode errors (Low/Medium Risk):**
   - **Evidence:** [auth_test.go](../../core-service/internal/api/auth_test.go#L127,L133,L144) and [governance_reports_api_test.go](../../core-service/internal/api/governance_reports_api_test.go#L30,L64,L84) use `http.DefaultClient` (no explicit timeout) via `http.Post`/`http.Get` and discard `json.NewDecoder(...).Decode(...)` errors (`_ = ...`). On unexpected status codes, `authResp.Token` silently decodes as `""`, causing confusing cascading failures rather than failing immediately with response body diagnostics.
   - **Positive note:** Unlike fragile suites that use `t.Skipf`, `SetupTestDB` intentionally calls `t.Fatalf` ([db.go](../../core-service/internal/testutil/db.go#L66)) to prevent Go's test cache from poisoning `(cached)` passes when Postgres `:5433` is down. Also no tests call `t.Parallel()` unsafely alongside `SetupTestDB`.

4. **`core-service/internal/testutil/db.go` — Good isolation, minor lock-error & pool-limit risks (Low Risk):**
   - **Evidence:** Uses `pg_advisory_lock(84920174)` to serialize `CREATE DATABASE ... TEMPLATE palabrario_test` across packages ([db.go](../../core-service/internal/testutil/db.go#L80-L95)) and drops clones via `t.Cleanup` (`WITH (FORCE)`). However, `pg_advisory_lock` / `pg_advisory_unlock` errors are silently ignored (`_, _ = conn.Exec(...)`), and if `t.Parallel()` is ever added to tests, unbounded concurrent `adminPool` clones can hit Postgres `max_connections`.

5. **`realtime-service/internal/config/config_test.go` & `image-service/internal/config/config_test.go` — Direct `os.Unsetenv` instead of `t.Setenv` (Low Risk):**
   - **Evidence:** Both [realtime config_test.go](../../realtime-service/internal/config/config_test.go#L11-L16) and [image config_test.go](../../image-service/internal/config/config_test.go#L10-L15) mutate process-wide env vars with `os.Unsetenv` in subtest 1 (preventing safe `t.Parallel()`). Also, `image-service`'s `TestConfigLoad` does **not** test invalid `JWT_EXPIRATION_HOURS` or short `JWT_SECRET`, leaving those bugs undetected.

6. **`pkg/auth/auth_test.go` — Only covers happy-path round-trip (Critical Coverage Gap):**
   - **Evidence:** Only 1 test (`TestPasswordHashingAndJWT`, [auth_test.go](../../pkg/auth/auth_test.go#L11-L41)). Zero tests for expired JWTs, future `nbf`, `alg=none` or RSA-confusion attacks, wrong secret, wrong `iss`, S2S token expiration, or invalid password hashes. Also, watch for **JWT second-level truncation** if future expiry tests use sub-second durations (`jwt.NewNumericDate` serializes Unix seconds as integers in [auth.go](../../pkg/auth/auth.go#L69-L71)).

7. **`image-service/internal/api/middleware/auth_middleware_test.go` & `realtime-service/internal/api/api_test.go` — Missing expired/S2S-cross-use & case-insensitive Bearer tests:**
   - Neither middleware test verifies: expired user JWT (`401`), S2S internal token rejected on user endpoints (`401`), `bearer <token>` lowercase scheme, or missing/empty token after `Bearer `.

### Critical Untested Auth Scenarios (Zero Coverage Today)

- **Suspended account with valid existing JWT** calling `/api/v1/auth/me`, `/api/v1/users/me`, `/api/v1/users/me/password`, `realtime-service` WS, or `image-service` `/upload`.
- **Role escalation/demotion token staleness** (`LEARNER` promoted via `RedeemInvitationTicket` / `ResolveEditorRequest`, or `ADMIN`/`EDITOR` demoted via `UpdateUserRoleOrStatus` while holding an old JWT).
- **`OptionalAuth` behavior** with an expired/malformed token vs missing token on `/api/v1/users/u/{username}`.
- **`RequireRole` & `RequireInternalService` HTTP tests** (`401` vs `403` status codes in `core-service/internal/api/`).
- **`SignUp` validation failures** (short password `< 8`, >72-byte password, invalid `anchor_dialect`, invalid email, invalid username characters, >50-char username / >100-char name currently causing `500`, concurrent `SignUp` conflict race).
- **`UpdateProfile` & `ChangePassword` HTTP handler tests** in `core-service/internal/api/` (`ChangePassword` has zero HTTP handler tests).
- **`app/` BFF auth suite (0 tests):** `App.doRequest` `401` handling, `SignIn`/`SignUp`/`SignOut`/`GetMe` token state, `UploadImage` auth header & `401`, `ConnectArenaRoom` header vs query token.
- **Frontend auth suite (0 tests):** `useAuth` `init()` race/concurrent callers, `WailsAuthAdapter` vs `MockAuthAPI` parity, router guards, `SignInView`/`SignUpView` error rendering.

---

## Remove / Simplify

| Item | Location(s) | Why Remove / Simplify | Recommendation |
| :--- | :--- | :--- | :--- |
| **`AuthResponse` in `IAuthAPI` + `toAuthResponse`** | [auth.contract.ts](../../app/frontend/src/api/contract/auth.contract.ts#L11-L12), [auth.wails.ts](../../app/frontend/src/api/live/wails/auth.wails.ts#L11-L15), [mock.adapter.ts](../../app/frontend/src/api/mock/mock.adapter.ts#L38,L116,L158) | Wails fabricates `token: ""`, Mock invents `"mock-session-token"`, and `useAuth.ts` completely ignores `res.token` | Have `IAuthAPI.signIn` / `signUp` return `Promise<UserProfile>`. Keep JWTs private inside transport adapters (Go RAM for Wails, `@capacitor/preferences` for future HTTP) |
| **`hasActiveSession()` in `IAuthAPI` & `HasActiveSession()` in Wails** | [auth.contract.ts](../../app/frontend/src/api/contract/auth.contract.ts#L10), [useAuth.ts](../../app/frontend/src/composables/useAuth.ts#L12-L13), [app.go](../../app/app.go#L147-L149) | Forces 2 sequential IPC calls on startup and cannot detect expired/revoked tokens (`sessionToken != ""` stays `true`) | Remove `hasActiveSession()` from `IAuthAPI` (and `HasActiveSession()` from `App` unless needed synchronously); call `getMe()` directly on `init()` and return `null` on `401` |
| **Duplicate BFF alias methods (`GetMyProfile`, etc.)** | [app_auth.go](../../app/app_auth.go#L50-L53), [app_arena.go](../../app/app_arena.go#L22,L65,L103,L190,L245), [app_governance.go](../../app/app_governance.go#L46,L77,L112,L140,L155,L175,L189,L207) | 14 thin 1-line aliases clutter `App` and bloat generated Wails TS bindings (`app.ts`) | Keep one canonical method per operation and regenerate `app/frontend/bindings` |
| **Duplicated HTTP & auth logic in `UploadImage`** | [app_media.go](../../app/app_media.go#L53-L84) vs [app.go](../../app/app.go#L87-L139) | Copies Bearer token injection, `Accept: application/json`, and `errResp` JSON parsing | Extract shared `App.executeRequest(req, target)` used by both `doRequest` and `UploadImage` (with unified `401` token cleanup) |
| **Mock vs Live error/null divergence** | [mock.adapter.ts](../../app/frontend/src/api/mock/mock.adapter.ts#L110-L170) vs [auth.wails.ts](../../app/frontend/src/api/live/wails/auth.wails.ts#L22-L36) | Mock returns `null` on bad login/unauthenticated `getMe()` and skips password verification; Wails throws on both | Standardize both adapters to throw `Error` on failed `signIn`/`signUp`, return `null` from `getMe()` when unauthenticated, and validate mock passwords + `is_active` |
| **Duplicated `SignIn` DB lookup** | [identity_service.go](../../core-service/internal/service/identity_service.go#L121,L134) | `SignIn` runs `GetAccountByEmail` then calls `GetMe` (`GetUserProfileByID`) | Merge into a single `GetAccountWithProfileByEmail` sqlc query |
| **Duplicate Bearer parsing & `JWT_SECRET` loading** | [core auth_middleware.go](../../core-service/internal/api/middleware/auth_middleware.go#L24-L40), [image auth_middleware.go](../../image-service/internal/api/middleware/auth_middleware.go#L25-L41), [arena_handler.go](../../realtime-service/internal/api/handlers/arena_handler.go#L38-L41) | 3 inconsistent copies (case-sensitive vs `EqualFold`, raw `err.Error()` leak in image-service, 3 divergent `config.Load` validators) | Extract `auth.ExtractBearerToken(r)` and `auth.LoadJWTConfig()` into `pkg/auth` |

---

## Ordered Refactor Ticket Breakdown

### Ticket 029 — Backend Auth Hardening, DB Constraints, BFF Session Lifecycle & Deterministic Go Tests
**Priority:** 1 (Do first — stabilizes backend/BFF contracts and security invariants)  
**Exact File Scope:**
- Database & sqlc: `SCHEMA_FINAL.sql`, `sql/queries/users.sql`, `core-service/internal/db/` (regenerated via `sqlc generate`)
- Shared Auth & Models: `pkg/auth/auth.go`, `pkg/auth/auth_test.go`
- Core Service: `core-service/internal/config/config.go`, `core-service/internal/api/router.go`, `core-service/internal/api/middleware/auth_middleware.go`, `core-service/internal/api/handlers/auth_handler.go`, `core-service/internal/api/handlers/user_handler.go`, `core-service/internal/api/handlers/governance_handler.go`, `core-service/internal/service/errors.go`, `core-service/internal/service/identity_service.go`, `core-service/internal/service/governance_service.go`, `core-service/internal/testutil/db.go`, `core-service/internal/api/auth_test.go`, `core-service/internal/api/governance_reports_api_test.go`, `core-service/internal/service/service_test.go`
- Sibling Services: `realtime-service/internal/config/config.go`, `realtime-service/internal/config/config_test.go`, `realtime-service/internal/api/router.go`, `realtime-service/internal/api/handlers/arena_handler.go`, `realtime-service/internal/api/api_test.go`, `realtime-service/internal/service/hub.go`, `realtime-service/internal/service/service_test.go`, `image-service/internal/config/config.go`, `image-service/internal/config/config_test.go`, `image-service/internal/api/router.go`, `image-service/internal/api/middleware/auth_middleware.go`, `image-service/internal/api/middleware/auth_middleware_test.go`
- Wails BFF: `app/app.go`, `app/app_auth.go`, `app/app_arena.go`, `app/app_governance.go`, `app/app_media.go`, plus new `app/app_auth_test.go`

**Key Changes & Acceptance Criteria:**
1. **Revocation & DB state enforcement:** Add `accounts.token_version` (or active/role DB check in `RequireAuth`) + `is_active = true` validation across `RequireAuth`, `GetMe`, `UpdateProfile`, and `ChangePassword`; bump `token_version` on suspension, role change, and `ChangePassword`; add `POST /api/v1/auth/refresh` (and refresh `sessionToken` in Wails after `UpdateProfile` / `RedeemInvitationTicket`).
2. **Input validation & SQL error mapping:** Enforce symmetric validation on `SignUp` (email format, password `8..72` bytes, `user_name` regex/length, `name` `2..100`, 8-dialect allowlist) + DB `CHECK`/FK constraints and `LOWER()` unique indexes; map `22001`/`23514` -> `400` and constraint-specific `23505` -> `409`; make `OptionalAuth` return `401` if a malformed/expired Bearer header is provided; return structured/English-consistent error messages.
3. **Cross-service alignment:** Remove hardcoded fallback `JWT_SECRET`s in `realtime-service` and `image-service`, share `ExtractBearerToken` in `pkg/auth`, remove `err.Error()` leak in `image-service`, fix `AllowedOrigins: ["*"]` + `AllowCredentials: true` across all 3 routers, and pass Wails WebSocket JWT via `Authorization: Bearer` header in `app_arena.go`.
4. **Wails BFF cleanup & `401` sync:** Clear `sessionToken` and emit `auth:expired` on `401` in shared `executeRequest` (covering `doRequest` + `UploadImage`), return typed `APIError` with HTTP status, decide/implement desktop session persistence if desired, remove duplicate BFF alias methods, and add `app/app_auth_test.go` using `httptest.Server`.
5. **Test hardening:** Fix `realtime-service/internal/service/service_test.go` sleep/timeout WebSocket flakiness, replace `UnixNano()` in `core-service/internal/service/service_test.go`, convert `os.Unsetenv` to `t.Setenv`, and add unit/integration tests for expired/forged/S2S-cross-use JWTs, suspended live JWT rejection, role-change invalidation, `OptionalAuth` invalid token, `RequireRole`, `SignUp` bounds, and `ChangePassword` HTTP.

---

### Ticket 030 — Frontend Auth Contract Simplification, Adapter Parity, Router Guards & Auth UX
**Priority:** 2 (Depends on Ticket `029` Wails binding regeneration & error contract)  
**Exact File Scope:**
- Generated bindings (via `task bindings` / Wails CLI): `app/frontend/bindings/palabrario/...`
- API Contract & Adapters: `app/frontend/src/api/config.ts`, `app/frontend/src/api/contract/auth.contract.ts`, `app/frontend/src/api/live/wails/auth.wails.ts`, `app/frontend/src/api/mock/mock.adapter.ts`, `app/frontend/src/api/mock/fixtures/user.fixture.ts` (if adding mock passwords/active flag)
- State, Navigation & Views: `app/frontend/src/composables/useAuth.ts`, `app/frontend/src/router/index.ts`, `app/frontend/src/components/TopBar.vue`, `app/frontend/src/views/auth/SignInView.vue`, `app/frontend/src/views/auth/SignUpView.vue`, `app/frontend/src/views/ProfileView.vue`, `app/frontend/src/views/ArenaView.vue`, `app/frontend/src/components/arena/CreateMatchDialog.vue`

**Key Changes & Acceptance Criteria:**
1. **Contract & composable cleanup:** Remove `AuthResponse` and `hasActiveSession()` from `IAuthAPI`; `signIn`/`signUp` return `Promise<UserProfile>` and throw on failure; `getMe()` returns `Promise<UserProfile | null>` cleanly (`null` when unauthenticated/401). Subscribe `useAuth` to Wails `auth:expired` events.
2. **Mock vs Wails parity:** Update `MockAuthAPI` to verify passwords, respect `is_active`, reject with errors matching live `IAuthAPI` error semantics, and strip `@` consistently with `core-service` (`SignUpView` + backend).
3. **Navigation guards & hydration:** Call `init()` in `router.beforeEach` (not only `TopBar.onMounted`), protect `/studio` and `/settings` with `meta: { requiresAuth: true }`, and guard `/auth/*` with `meta: { guestOnly: true }`.
4. **View validation & error hygiene:** Enforce username regex (`^[a-zA-Z0-9_-]{3,50}$`) and `8..72` password checks in `SignUpView.vue`, handle `ProfileView.vue` non-404 errors and reactive owner re-hydration, gate Arena room creation behind auth, and prepare `api/config.ts` cleanly for future `live/http/` Capacitor adapter (`@capacitor/preferences`).
5. **Verification:** `npm run type-check` and `npm run build` pass cleanly in `app/frontend/`, plus manual/automated verification in both `mock` and `wails` modes.

---

### Ticket 031 — Post-Refactor Auth Verification & Security Review
**Priority:** 3 (Reviewer pass after Tickets `029` and `030`)  
**Exact File Scope:** Read-only audit across all Ticket `029` and `030` touched files + test runner verification.  
**Acceptance Criteria:**
- Re-run `go test -race -count=1 ./pkg/auth/... ./image-service/... ./realtime-service/... ./app/...`, `task test` (with `postgres-test` `:5433` up), and frontend `type-check`/`build`.
- Verify zero remaining hardcoded JWT fallbacks, `401` auto-cleanup end-to-end, immediate effect of account suspension and role changes, and full `MockAuthAPI` / `WailsAuthAdapter` parity.
