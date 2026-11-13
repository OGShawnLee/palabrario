# 031 — Harden Core Service Auth, Short-Lived JWTs, Validation, and English Errors

| Field | Value |
| :--- | :--- |
| **ID** | 031 |
| **Status** | completed |
| **Model** | Gemini 4 Hydrogen, `thinking_level: medium` (`gemini-4-hydrogen`) |
| **Role** | Backend Implementer (`implementation`) |
| **Spec** | [spec.md](spec.md) |

> **Specification Note:**  
> Third ticket of the end-to-end auth refactor (depends on `029` and `030`). Implements short-lived (`15m`) access JWTs + rotating opaque refresh tokens (`POST /api/v1/auth/refresh`, `POST /api/v1/auth/sign-out`), **Option B** active-account & role verification middleware (every protected route verifies the 15m JWT in memory; all `POST`/`PUT`/`DELETE` routes plus `/api/v1/auth/*` and `/api/v1/governance/*` also verify `is_active`, `token_version`, and sync live `role` against PostgreSQL), `OptionalAuth` 401 on invalid tokens, symmetric `SignUp` validation, `pgError` constraint mapping (`22001`/`23514` -> `400`, distinct email vs username `23505`), English error messages with machine-readable `code`, and full HTTP + service tests.

## Scope

May create or edit:
- `pkg/models/auth.go`
- `core-service/internal/config/config.go`
- `core-service/internal/service/errors.go`
- `core-service/internal/service/identity_service.go`
- `core-service/internal/service/governance_service.go`
- `core-service/internal/api/middleware/auth_middleware.go`
- `core-service/internal/api/router.go`
- `core-service/internal/api/handlers/auth_handler.go`
- `core-service/internal/api/handlers/user_handler.go`
- `core-service/internal/api/handlers/governance_handler.go`
- `core-service/internal/testutil/db.go`
- `core-service/internal/api/auth_test.go`
- `core-service/internal/api/governance_reports_api_test.go`
- `core-service/internal/service/service_test.go`
- `core-service/internal/service/governance_reports_test.go`
- `.routing/DOMAIN.md` (Section 3 rule 4, Section 4 Option B middleware, Section 4.2 English error contract)
- `.routing/services/identity_service.md` (Part A `REFRESH_TOKENS`/`token_version`/constraints/suspended-profile rule + Part B `2.1.1`–`2.1.4`)
- `.routing/services/governance_service.md` (Part B banner, English rules & token invalidation in `2.2.3`, Option B in `2.2.4`)
- `.routing/services/dictionary_service.md` (Part B banner, `OptionalAuth` `401` & Option B writes in `2.3.5`)
- `.routing/services/curation_service.md`, `.routing/services/studio_service.md`, `.routing/services/mastery_service.md` (Part B banner `response.ErrorWithCode` / English `{ "code", "error" }` + Option B write-strict note)
- `.tickets/031-harden-core-service-auth/changelog.md` (create)

Must not touch:
- `SCHEMA_FINAL.sql`, `sql/`, `core-service/internal/db/`
- `pkg/auth/`, `pkg/response/`, `realtime-service/`, `image-service/`, `app/`, `.context/`
- `.routing/CONTEXT.md`, `.routing/DESIGN.md`, `.routing/PROGRESS.md`, `.routing/API.md`, `.routing/ROUTES.md`, `.routing/INFRASTRUCTURE.md`, `.routing/services/arena_service.md`, `.routing/services/media_service.md`
- Any file outside the above scope
- **No role (Router, Investigator, Implementer, Translator, or Reviewer) may commit, push, or merge without explicit user authorization.** Do not run `git commit`, `git push`, or `git merge`.

## Inputs

- [.tickets/028-investigate-auth-workflow/report.md](../028-investigate-auth-workflow/report.md): Findings `CORE-01`..`CORE-10` and test audit.
- [pkg/models/auth.go](../../pkg/models/auth.go), [pkg/auth/auth.go](../../pkg/auth/auth.go), [pkg/response/json.go](../../pkg/response/json.go)
- [core-service/internal/service/errors.go](../../core-service/internal/service/errors.go), [identity_service.go](../../core-service/internal/service/identity_service.go), [governance_service.go](../../core-service/internal/service/governance_service.go)
- [core-service/internal/api/router.go](../../core-service/internal/api/router.go), [auth_middleware.go](../../core-service/internal/api/middleware/auth_middleware.go)
- [.routing/DOMAIN.md](../../.routing/DOMAIN.md) and [.routing/services/](../../.routing/services/) (`identity_service.md`, `governance_service.md`, `dictionary_service.md`, `curation_service.md`, `studio_service.md`, `mastery_service.md`)

## Task

1. **Models & Config (`pkg/models/auth.go`, `core-service/internal/config/config.go`):**
   - In `pkg/models/auth.go`: add `RefreshToken string `json:"refresh_token"`` to `AuthResponse`, plus `RefreshTokenRequest{RefreshToken string `json:"refresh_token"`}` and `SignOutRequest{RefreshToken string `json:"refresh_token"`}`.
   - In `core-service/internal/config/config.go`: add `AccessTokenTTL time.Duration` (`ACCESS_TOKEN_TTL_MINUTES`, default `15 * time.Minute`), `RefreshTokenTTL time.Duration` (`REFRESH_TOKEN_TTL_HOURS`, fallback `JWT_EXPIRATION_HOURS`, default `168 * time.Hour`), and `AllowedOrigins []string` (`CORS_ALLOWED_ORIGINS` via `auth.ParseCORSAllowedOrigins`). Use English config error messages. Update `router.go` to pass `rt.cfg.AccessTokenTTL` to `auth.NewTokenService` and `rt.cfg.RefreshTokenTTL` to `NewIdentityService`.

2. **English Domain Errors & PgError Mapping ([core-service/internal/service/errors.go](../../core-service/internal/service/errors.go)):**
   - Translate `Err*` sentinel descriptions to English and add `Code string` to `DomainError` (plus `NewDomainErrorWithCode(base error, code, message string)` and `HTTPStatusCodeAndMessage(err error, fallbackCode, fallbackMsg string) (int, string, string)`).
   - Update `MapDBError` to handle:
     - `pgx.ErrNoRows` -> `ErrNotFound` (`code: "NOT_FOUND"`)
     - `22001` (`string_data_right_truncation`) -> `ErrInvalidInput` (`400`, `code: "INPUT_TOO_LONG"`)
     - `23514` (`check_violation`) -> `ErrValidation` (`400`, `code: "VALIDATION_ERROR"`)
     - `23503` (`foreign_key_violation`) -> `ErrInvalidInput` (`400`, `code: "INVALID_REFERENCE"`)
     - `23505` (`unique_violation`): inspect `pgErr.ConstraintName` so `idx_accounts_email_lower` / `accounts_email_key` returns `EMAIL_ALREADY_EXISTS` (`"An account with that email already exists."`) and `idx_users_user_name_lower` / `users_user_name_key` returns `USERNAME_ALREADY_EXISTS` (`"That username is already taken."`).

3. **Identity & Governance Services (`identity_service.go`, `governance_service.go`):**
   - In `IdentityService` interface and `identityService`:
     - Add `RefreshToken(ctx context.Context, refreshToken string) (*models.AuthResponse, error)` and `SignOut(ctx context.Context, refreshToken string) error`.
     - **Symmetric `SignUp` validation:** normalize `email`, strip optional leading `@` from `user_name` and lowercase, trim `name` (defaulting to normalized `user_name` only if empty is omitted — wait: require `len(name) in [2,100]` or match frontend: if `name == ""` fall back to `userName`, then enforce `len([]rune(name)) in [2,100]`); validate `email` format with `mail.ParseAddress`, `user_name` regex `^[a-z0-9][a-z0-9_-]{2,49}$`, password byte length `8..72`, and `anchor_dialect` in `{MX,ES,AR,CO,CL,CU,CAM,PAN}` (default `MX` when empty). Inside tx, issue access JWT (`TokenVersion: 1`) + persisted hashed refresh token.
     - **`SignIn`:** replace two queries with single `GetAccountWithProfileByEmail` call. Reject `!is_active` with `403 Forbidden` (`code: "ACCOUNT_SUSPENDED"`, `"This account has been suspended by an administrator."`). Reject missing user / bad password with `401` (`code: "INVALID_CREDENTIALS"`, `"Invalid email or password."`). Issue access JWT (`TokenVersion: account.TokenVersion`) + persisted hashed refresh token.
     - **`RefreshToken`:** hash incoming raw token via `auth.HashRefreshToken`, inside `pgx.Tx` lock row via `GetRefreshTokenForUpdate`, reject missing/revoked/expired tokens with `401` (`code: "TOKEN_EXPIRED"` or `"TOKEN_REVOKED"`), load `GetUserProfileByID`, reject `!profile.IsActive` with `403 ACCOUNT_SUSPENDED` (and revoke all refresh tokens for account), revoke used refresh token (`RevokeRefreshTokenByID`), create new refresh token, and return new `AuthResponse{Token, RefreshToken, User}` signed with current DB `role`, `anchor_dialect`, and `token_version`.
     - **`SignOut`:** if non-empty, revoke matching hashed token via `RevokeRefreshTokenByHash`.
     - **`GetMe`:** map `pgx.ErrNoRows` to `401` (`code: "UNAUTHORIZED"`) and `!row.IsActive` to `403` (`code: "ACCOUNT_SUSPENDED"`).
     - **`ChangePassword`:** enforce password byte length `8..72` on `NewPassword`, reject `!account.IsActive` with `403 ACCOUNT_SUSPENDED`, check current password (`401 INVALID_CREDENTIALS`), then inside tx update password hash, call `IncrementAccountTokenVersion`, and call `RevokeAllRefreshTokensForAccount`.
     - **`GetPublicProfileByUsername`:** propagate non-`pgx.ErrNoRows` errors from secondary queries instead of swallowing `err != nil` silently.
     - **Translate all messages in `identity_service.go` and auth-related `governance_service.go` errors to English.**
   - In `governance_service.go`:
     - In `UpdateUserRoleOrStatus`: when `req.IsActive != nil && !*req.IsActive`, call `IncrementAccountTokenVersion` + `RevokeAllRefreshTokensForAccount` inside a transaction; when `req.Role != nil` (promotion or demotion), call `IncrementAccountTokenVersion` inside the transaction.
     - In `ResolveEditorRequest` (`APPROVED`) and `RedeemInvitationTicket` (promoted to `EDITOR`): also call `IncrementAccountTokenVersion` inside the same transaction.

4. **Option B Auth Middleware, CORS, Handlers & Router:**
   - Update `auth_middleware.go` to take `*db.Queries` (or `AuthStateStore` interface with `GetAuthAccountStateByID`) alongside `*auth.TokenService`.
   - Implement `verifyAccessToken(r)` using `auth.ExtractBearerToken(r)` (`401 MISSING_TOKEN` / `INVALID_TOKEN` / `TOKEN_EXPIRED`).
   - Implement **Option B** active-account & role check helper (`enforceActiveAccount` inside `RequireActiveAccount` middleware or inside `RequireAuth` checking method/path):
     - Every protected route runs `RequireAuth` (in-memory `15m` JWT verification).
     - Routes matching `(r.Method != http.MethodGet && r.Method != http.MethodHead && r.Method != http.MethodOptions)` **OR** `/api/v1/auth/` **OR** `/api/v1/governance/` additionally query `GetAuthAccountStateByID`:
       - If missing -> `401 UNAUTHORIZED`.
       - If `!state.IsActive` -> `403 ACCOUNT_SUSPENDED` (`"This account has been suspended by an administrator."`).
       - If `int(state.TokenVersion) != claims.TokenVersion` -> `401 TOKEN_REVOKED` (`"Session is no longer valid. Please sign in again."`).
       - Sync `claims.Role = models.AccountRole(state.Role)` before downstream `RequireRole` runs (`RF04` immediate role enforcement on writes & governance).
   - Update `OptionalAuth`: if `Authorization` header is absent, proceed anonymously; if present and invalid/expired, return `401` with `ErrorWithCode`.
   - Register `POST /api/v1/auth/refresh` and `POST /api/v1/auth/sign-out` in `router.go`, switch CORS to `AllowedOrigins: rt.cfg.AllowedOrigins`, `AllowCredentials: false`, and update `auth_handler.go`, `user_handler.go`, and `governance_handler.go` to use `HTTPStatusCodeAndMessage` + `response.ErrorWithCode`.

5. **Tests & Routing Docs Sync:**
   - Fix `testutil/db.go` (`pg_advisory_lock`/`unlock` error checks) and `service_test.go` (replace `time.Now().UnixNano()` synthetic emails/usernames with `uuid.NewString()`).
   - Fix `auth_test.go` and `governance_reports_api_test.go` (`&http.Client{Timeout: 10 * time.Second}`, check `Decode` errors, update any hardcoded Spanish error assertions to English/codes).
   - Add tests in `core-service/internal/api/auth_test.go` and `core-service/internal/service/` covering:
     - Suspended user denied `403 ACCOUNT_SUSPENDED` on `/auth/me`, `POST`/`PUT`/`DELETE` writes, `/governance/*`, `/auth/refresh`, and `/users/u/{username}` (`404`).
     - Role change (`LEARNER` <-> `EDITOR` <-> `ADMIN`) takes immediate effect via DB check on `RequireRole` + `/governance/*` + writes (old role rejected, promoted role accepted/refreshed).
     - Access token expiration (`401 TOKEN_EXPIRED`), refresh token rotation & reuse rejection (`401`), `SignOut` revoking refresh token, and `ChangePassword` revoking prior refresh tokens & bumping `token_version`.
     - `OptionalAuth` returning `401` on invalid/expired Bearer header vs `200` with no header.
     - `SignUp` validation errors (`400` short password, `>72`-byte password, bad email, bad username, invalid dialect, `>50` char username -> `400`, duplicate email `EMAIL_ALREADY_EXISTS` vs duplicate username `USERNAME_ALREADY_EXISTS`).
     - `ChangePassword` HTTP handler tests (`204`, wrong password `401`, invalid body `400`).
   - Update routing docs in lockstep:
     - [.routing/DOMAIN.md](../../.routing/DOMAIN.md): Section 3 rule 4 (English `{ "code", "error" }` format via `response.Error` / `response.ErrorWithCode`), Section 4 Mermaid Option B middleware note, and Section 4.2 error table (`400 INPUT_TOO_LONG`/`VALIDATION_ERROR`/`INVALID_REFERENCE`, `401 TOKEN_EXPIRED`/`TOKEN_REVOKED`, `403 ACCOUNT_SUSPENDED`, `409 EMAIL_ALREADY_EXISTS`/`USERNAME_ALREADY_EXISTS`).
     - [.routing/services/identity_service.md](../../.routing/services/identity_service.md): **Part A** (add `REFRESH_TOKENS` to ER diagram + document `RefreshToken`, `Account.token_version`, `User.anchor_dialect` FK + `CHECK` constraints + `LOWER()` unique indexes, and suspended `/u/:username` `404` rule) and **Part B** (Part B banner, `2.1.1`–`2.1.4` new DTOs, `IdentityService` methods, `sqlc` queries, `/auth/refresh` & `/auth/sign-out` endpoints, `403`/`401` codes).
     - [.routing/services/governance_service.md](../../.routing/services/governance_service.md): update Part B banner, translate quoted Spanish conflict strings in `2.2.3` to English, document `IncrementAccountTokenVersion` / `RevokeAllRefreshTokensForAccount` transactional rules, and note Option B live DB enforcement on `/api/v1/governance/*` in `2.2.4`.
     - [.routing/services/dictionary_service.md](../../.routing/services/dictionary_service.md): update Part B banner and clarify in `2.3.5` that `OptionalAuth` allows anonymous requests without an `Authorization` header **but returns `401` on invalid/expired Bearer tokens**, and writes (`POST`/`PUT`/`DELETE`) run Option B live DB checks.
     - [.routing/services/curation_service.md](../../.routing/services/curation_service.md), [.routing/services/studio_service.md](../../.routing/services/studio_service.md), and [.routing/services/mastery_service.md](../../.routing/services/mastery_service.md): update the Part B `> [!IMPORTANT]` banner (`response.ErrorWithCode` English `{ "code", "error" }` + Option B write-strict auth note).

## Verification

Run every check in [spec.md](spec.md).

## Change log

Write `.tickets/031-harden-core-service-auth/changelog.md` recording all changed endpoints, error codes, middleware behavior, docs updates, and exact `go test` results against `:5433`.
