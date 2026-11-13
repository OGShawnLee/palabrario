# Ticket 031 Changelog — Harden Core Service Auth, Short-Lived JWTs, Validation, and English Errors

## Summary

Implemented Ticket `031-harden-core-service-auth` across `pkg/models`, `core-service`, `testutil`, test suites, and `.routing/` contract documentation:
- Replaced 72h stateless user JWTs with short-lived (`15m` default) access JWTs carrying `token_version` and `token_type: "access"`, paired with rotating single-use opaque refresh tokens (`168h` default) persisted as SHA-256 hashes in `refresh_tokens`.
- Implemented **Option B** authorization middleware (`RequireAuth`): all protected routes verify the 15m JWT in memory, while all write routes (`Method != GET/HEAD/OPTIONS`), `/api/v1/auth/*`, and `/api/v1/governance/*` also query `GetAuthAccountStateByID` against PostgreSQL to reject suspended accounts (`403 ACCOUNT_SUSPENDED`), reject invalidated sessions (`401 TOKEN_REVOKED`), and synchronize live DB `role` onto `claims.Role` before `RequireRole`.
- Hardened `OptionalAuth` to continue anonymously only when no `Authorization` header is present, and return `401` (`INVALID_TOKEN_FORMAT` / `TOKEN_EXPIRED` / `INVALID_TOKEN`) when an invalid or expired header is sent.
- Added symmetric `SignUp` validation, `8..72`-byte bcrypt password guards, single-query `SignIn`, transactional token-version bumps + refresh-token revocation across `ChangePassword` and `GovernanceService` mutations, `pgconn.PgError` (`22001`, `23514`, `23503`, distinct `23505`) mapping, CORS allowlist (`AllowCredentials: false`), and standardized English error responses with machine-readable `code` fields via `response.ErrorWithCode`.

## Files Modified

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
- `.routing/DOMAIN.md`
- `.routing/services/identity_service.md`
- `.routing/services/governance_service.md`
- `.routing/services/dictionary_service.md`
- `.routing/services/curation_service.md`
- `.routing/services/studio_service.md`
- `.routing/services/mastery_service.md`

## Endpoints & Auth Behavior Changes

- `POST /api/v1/auth/sign-up`: Validates email (`mail.ParseAddress`), strips optional leading `@` & enforces `^[a-z0-9][a-z0-9_-]{2,49}$` on `user_name`, `2..100` rune display `name` (falling back to normalized `user_name` when omitted), `8..72`-byte `password`, and dialect allowlist `{MX, ES, AR, CO, CL, CU, CAM, PAN}` (default `MX`); returns `201` with `token`, `refresh_token`, and `user`, or `409 EMAIL_ALREADY_EXISTS` / `409 USERNAME_ALREADY_EXISTS`.
- `POST /api/v1/auth/sign-in`: Uses single-query `GetAccountWithProfileByEmail`; returns `401 INVALID_CREDENTIALS` on unknown email/wrong password, `403 ACCOUNT_SUSPENDED` when `is_active = false`, and issues `token` + `refresh_token` on `200 OK`.
- `POST /api/v1/auth/refresh`: Public endpoint performing atomic `SELECT ... FOR UPDATE` refresh token rotation (`200 OK` with fresh `token` + rotated `refresh_token`, `401 TOKEN_REVOKED` / `401 TOKEN_EXPIRED`, `403 ACCOUNT_SUSPENDED` with full account refresh-token revocation).
- `POST /api/v1/auth/sign-out`: Public endpoint revoking the supplied `refresh_token` by SHA-256 hash (`204 No Content`, tolerates empty body, `400 INVALID_JSON` on malformed payload).
- `GET /api/v1/auth/me`: Protected via Option B `RequireAuth`; returns `401 UNAUTHORIZED` when missing, `401 TOKEN_EXPIRED` / `401 TOKEN_REVOKED`, and `403 ACCOUNT_SUSPENDED` for suspended accounts.
- `PUT /api/v1/users/me`: Protected via Option B `RequireAuth`; verifies active status and syncs DB role.
- `PUT /api/v1/users/me/password`: Enforces `8..72`-byte `new_password`, verifies active status and current password (`401 INVALID_CREDENTIALS`), updates hash, increments `token_version`, and revokes all account refresh tokens transactionally (`204 No Content`).
- `GET /api/v1/users/u/{username}`: Protected via hardened `OptionalAuth` (`401` on invalid/expired Bearer header); propagates non-`pgx.ErrNoRows` secondary query errors; returns `404` for suspended accounts.
- `GovernanceService`: `ResolveEditorRequest` (`APPROVED`) and `RedeemInvitationTicket` (`EDITOR` promotion) increment `token_version` transactionally; `UpdateUserRoleOrStatus` increments `token_version` on role changes or deactivation (`is_active = false`) and revokes all refresh tokens on deactivation. All `/api/v1/governance/*` routes enforce Option B live DB checks.
- **CORS (`core-service/internal/api/router.go` & `config.go`):** Uses `CORS_ALLOWED_ORIGINS` parsed via `auth.ParseCORSAllowedOrigins` (required when `ENVIRONMENT=production`) and sets `AllowCredentials: false`.

## Machine-Readable Error Codes & PgError Mapping

- `HTTPStatusCodeAndMessage` + `response.ErrorWithCode` emit `{ "code": "...", "error": "..." }` in English.
- `MapDBError` mappings:
  - `pgx.ErrNoRows` -> `404 NOT_FOUND`
  - `22001` (`string_data_right_truncation`) -> `400 INPUT_TOO_LONG`
  - `23514` (`check_violation`) -> `400 VALIDATION_ERROR`
  - `23503` (`foreign_key_violation`) -> `400 INVALID_REFERENCE`
  - `23505` (`unique_violation`): `idx_accounts_email_lower` / `accounts_email_key` -> `409 EMAIL_ALREADY_EXISTS`; `idx_users_user_name_lower` / `users_user_name_key` -> `409 USERNAME_ALREADY_EXISTS`; other constraints -> `409 CONFLICT`
  - `P0001` -> `403 FORBIDDEN`
- Middleware/Auth codes: `MISSING_TOKEN`, `INVALID_TOKEN_FORMAT`, `INVALID_TOKEN`, `TOKEN_EXPIRED`, `TOKEN_REVOKED`, `INVALID_CREDENTIALS`, `UNAUTHORIZED`, `FORBIDDEN`, `ACCOUNT_SUSPENDED`, `INVALID_JSON`, `VALIDATION_ERROR`, `INVALID_INPUT`.

## Documentation Updates

- `.routing/DOMAIN.md`: Updated Section 3 rule 4 to English `{ "code", "error" }`, Section 4 Mermaid diagram for Option B middleware, and Section 4.2 table with error codes.
- `.routing/services/identity_service.md`: Updated Part A (`REFRESH_TOKENS`, `token_version`, `anchor_dialect` FK/constraints, `LOWER()` unique indexes, suspended `/u/:username` `404` behavior) and Part B (`2.1.1`–`2.1.4` DTOs, `RefreshToken`/`SignOut`, `sqlc` query list, HTTP endpoint status/code table).
- `.routing/services/governance_service.md`: Updated Part B banner, translated `2.2.3` conflict messages to English, documented `IncrementAccountTokenVersion` & `RevokeAllRefreshTokensForAccount`, and noted Option B live DB checks in `2.2.4`.
- `.routing/services/dictionary_service.md`: Updated Part B banner and clarified `OptionalAuth` `401` and Option B write checks in `2.3.3`.
- `.routing/services/curation_service.md`, `.routing/services/studio_service.md`, `.routing/services/mastery_service.md`: Updated Part B banners to `response.ErrorWithCode` English `{ "code", "error" }` + Option B write-strict auth note.

## Verification (`postgres-test :5433`)

- `go build ./core-service/... ./pkg/...` — exited with code `0`
- `go vet ./core-service/... ./pkg/...` — exited with code `0`
- `go test -race -count=1 ./core-service/...` — exited with code `0`:

```text
?   	palabrario/core-service	[no test files]
ok  	palabrario/core-service/internal/api	49.915s
?   	palabrario/core-service/internal/api/handlers	[no test files]
?   	palabrario/core-service/internal/api/middleware	[no test files]
?   	palabrario/core-service/internal/config	[no test files]
?   	palabrario/core-service/internal/database	[no test files]
?   	palabrario/core-service/internal/db	[no test files]
ok  	palabrario/core-service/internal/service	56.317s
?   	palabrario/core-service/internal/testutil	[no test files]
```

