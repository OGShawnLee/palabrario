# 029 — Changelog: Harden Database Auth Schema and sqlc Queries

## 1. Schema Changes in `SCHEMA_FINAL.sql`

- **Section ordering (`regions` before `accounts`):**
  - Moved `CREATE TABLE regions (...)` (now Section 3) immediately above `CREATE TABLE accounts (...)` (Section 4) so `users.anchor_dialect` can reference `regions(code)` cleanly in both PostgreSQL and `sqlc`.
- **`accounts` hardening:**
  - Added `token_version INTEGER NOT NULL DEFAULT 1 CHECK (token_version >= 1)`.
  - Added case-insensitive unique index `CREATE UNIQUE INDEX idx_accounts_email_lower ON accounts (LOWER(email));` while retaining `email VARCHAR(255) NOT NULL UNIQUE` so existing `ON CONFLICT (email) DO NOTHING` seed statements continue working.
- **`refresh_tokens` table:**
  - Added immediately after `accounts`:
    - `id_token UUID PRIMARY KEY DEFAULT gen_random_uuid()`
    - `id_account UUID NOT NULL REFERENCES accounts(id_account) ON DELETE CASCADE`
    - `token_hash CHAR(64) NOT NULL UNIQUE`
    - `expires_at TIMESTAMPTZ NOT NULL`
    - `revoked_at TIMESTAMPTZ`
    - `created_at TIMESTAMPTZ NOT NULL DEFAULT now()`
  - Added index `CREATE INDEX idx_refresh_tokens_account ON refresh_tokens(id_account);`.
- **`users` hardening:**
  - Updated `anchor_dialect` to `VARCHAR(10) NOT NULL DEFAULT 'MX' REFERENCES regions(code)`.
  - Added `CHECK (char_length(trim(name)) BETWEEN 2 AND 100)` on `users.name`.
  - Added `CHECK (user_name ~ '^[a-z0-9][a-z0-9_-]{2,49}$')` on `users.user_name`.
  - Added case-insensitive unique index `CREATE UNIQUE INDEX idx_users_user_name_lower ON users (LOWER(user_name));` while retaining `user_name VARCHAR(50) NOT NULL UNIQUE` for existing `ON CONFLICT (user_name) DO NOTHING` seeds.

## 2. Updated and Added Queries in `sql/queries/users.sql`

- **Updated queries:**
  - `CreateAccount :one` — now also returns `token_version`.
  - `GetAccountByEmail :one` — now also selects `token_version` and filters case-insensitively with `WHERE LOWER(email) = LOWER($1)`.
  - `GetAccountByID :one` — now also selects `token_version`.
  - `GetUserProfileByID :one` — now also selects `a.is_active` and `a.token_version`.
  - `GetUserProfileByUsername :one` — now enforces `AND a.is_active = true` so suspended profiles are hidden (`404`).
  - `UpdateUserProfile :one` — now also selects `a.is_active` and `a.token_version`.
- **New queries:**
  - `GetAuthAccountStateByID :one` — selects `id_account, role, is_active, token_version` from `accounts` by `id_account`.
  - `GetAccountWithProfileByEmail :one` — joins `accounts` and `users` on `LOWER(a.email) = LOWER($1)` and returns `a.id_account, a.email, a.password_hash, a.role, a.is_active, a.token_version, u.name, u.user_name, u.avatar_url, u.anchor_dialect, u.total_xp, u.current_level, u.current_streak, u.longest_streak, u.created_at, u.updated_at`.
  - `IncrementAccountTokenVersion :one` — atomically increments `accounts.token_version` and updates `updated_at`, returning the new `token_version`.
  - `CreateRefreshToken :one` — inserts `(id_account, token_hash, expires_at)` into `refresh_tokens` and returns `*`.
  - `GetRefreshTokenForUpdate :one` — selects `id_token, id_account, token_hash, expires_at, revoked_at, created_at` from `refresh_tokens` by `token_hash` with `FOR UPDATE`.
  - `RevokeRefreshTokenByID :exec` — sets `revoked_at = now()` where `id_token = $1 AND revoked_at IS NULL`.
  - `RevokeRefreshTokenByHash :exec` — sets `revoked_at = now()` where `token_hash = $1 AND revoked_at IS NULL`.
  - `RevokeAllRefreshTokensForAccount :exec` — sets `revoked_at = now()` where `id_account = $1 AND revoked_at IS NULL`.

## 3. Code Generation, Build, and Verification

- **`sqlc generate`:** Exited with code `0` (run with `GOMAXPROCS=1` on Windows to avoid wazero concurrent 4 GiB virtual-memory commit failures), updating `core-service/internal/db/models.go` and `core-service/internal/db/users.sql.go`.
- **`go build ./core-service/internal/db/...`:** Exited with code `0`.
- **`go build ./core-service/...`:** Exited with code `0`. All existing service callers remain compatible without edits in this ticket.
- **Notes for Ticket `031` (`core-service` auth hardening):**
  - `db.Account` now exposes `TokenVersion int32`.
  - `db.GetUserProfileByIDRow` and `db.UpdateUserProfileRow` now expose `IsActive bool` and `TokenVersion int32`.
  - `db.RefreshToken` and the new `*db.Queries` methods (`GetAuthAccountStateByID`, `GetAccountWithProfileByEmail`, `IncrementAccountTokenVersion`, `CreateRefreshToken`, `GetRefreshTokenForUpdate`, `RevokeRefreshTokenByID`, `RevokeRefreshTokenByHash`, `RevokeAllRefreshTokensForAccount`) are ready for wiring in `identity_service.go`, `governance_service.go`, and `auth_middleware.go`.
- **PostgreSQL `:5433` schema check (optional in `029`):** Loaded `SCHEMA_FINAL.sql` into a fresh temporary database on `localhost:5433` with zero SQL errors (`accounts=3 users=3 regions=8`), then dropped the temporary database.
- **Git check:** `git status --short` confirms only `SCHEMA_FINAL.sql`, `sql/queries/users.sql`, `core-service/internal/db/models.go`, `core-service/internal/db/users.sql.go`, and `.tickets/029-backend-auth-db-schema/changelog.md` were created or modified by this ticket; no `git commit`, `git push`, or `git merge` was run (`git log -1` unchanged at `f470c79`).
