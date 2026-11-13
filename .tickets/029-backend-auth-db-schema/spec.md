# 029 — Spec: Harden Database Auth Schema and sqlc Queries

- [x] `SCHEMA_FINAL.sql` defines `regions` before `accounts` and sets `users.anchor_dialect VARCHAR(10) NOT NULL DEFAULT 'MX' REFERENCES regions(code)`.
- [x] `SCHEMA_FINAL.sql` adds `token_version INTEGER NOT NULL DEFAULT 1 CHECK (token_version >= 1)` to `accounts` and creates the `refresh_tokens` table with `token_hash CHAR(64) NOT NULL UNIQUE`, `expires_at`, `revoked_at`, and `idx_refresh_tokens_account`.
- [x] `SCHEMA_FINAL.sql` adds `idx_accounts_email_lower`, `idx_users_user_name_lower`, `CHECK (char_length(trim(name)) BETWEEN 2 AND 100)`, and `CHECK (user_name ~ '^[a-z0-9][a-z0-9_-]{2,49}$')`.
- [x] `sql/queries/users.sql` updates `GetUserProfileByUsername` with `AND a.is_active = true`, uses `LOWER(email) = LOWER($1)` in `GetAccountByEmail`, and adds `GetAuthAccountStateByID`, `GetAccountWithProfileByEmail`, `IncrementAccountTokenVersion`, `CreateRefreshToken`, `GetRefreshTokenForUpdate`, `RevokeRefreshTokenByID`, `RevokeRefreshTokenByHash`, and `RevokeAllRefreshTokensForAccount`.
- [x] `sqlc generate` exits with code `0` and updates `core-service/internal/db/`.
- [x] `go build ./core-service/internal/db/...` exits with code `0`.
- [x] If PostgreSQL test container (`:5433`) is available (optional check in `029`; mandatory in `031`), `SCHEMA_FINAL.sql` loads with zero SQL errors.
- [x] `git status --short` confirms no files outside `SCHEMA_FINAL.sql`, `sql/queries/users.sql`, `core-service/internal/db/`, and `.tickets/029-backend-auth-db-schema/changelog.md` were created or modified.
- [x] No unauthorized `git commit`, `git push`, or `git merge` was run (`git log -1` unchanged). (If committed after explicit user authorization: English, Conventional Commits, detailed body.)
