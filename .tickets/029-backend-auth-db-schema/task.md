# 029 — Harden Database Auth Schema and sqlc Queries

| Field | Value |
| :--- | :--- |
| **ID** | 029 |
| **Status** | completed |
| **Model** | Gemini 4 Hydrogen, `thinking_level: medium` (`gemini-4-hydrogen`) |
| **Role** | Backend Implementer (`implementation`) |
| **Spec** | [spec.md](spec.md) |

> **Specification Note:**  
> First ticket of the end-to-end auth refactor ([028 report](../028-investigate-auth-workflow/report.md)). Adds `accounts.token_version`, the `refresh_tokens` table, `users.anchor_dialect` foreign key to `regions(code)`, `CHECK` constraints on `users.name` and `users.user_name`, case-insensitive unique indexes on `accounts(LOWER(email))` and `users(LOWER(user_name))`, hides suspended users (`is_active = false`) in `GetUserProfileByUsername`, and adds the typed `sqlc` queries needed for short-lived JWT + opaque refresh token auth and write-strict `is_active`/role verification.

## Scope

May create or edit:
- `SCHEMA_FINAL.sql`
- `sql/queries/users.sql`
- `core-service/internal/db/` (exclusively via `sqlc generate` — never hand-edit generated files)
- `.tickets/029-backend-auth-db-schema/changelog.md` (create)

Must not touch:
- `core-service/internal/service/`, `core-service/internal/api/`, `core-service/internal/config/`, `core-service/internal/testutil/`
- `pkg/`, `realtime-service/`, `image-service/`, `app/`, `.routing/`, `.context/`
- Any file outside the above scope
- **No role (Router, Investigator, Implementer, Translator, or Reviewer) may commit, push, or merge without explicit user authorization.** Do not run `git commit`, `git push`, or `git merge`.

## Inputs

- [.tickets/028-investigate-auth-workflow/report.md](../028-investigate-auth-workflow/report.md): Findings `DB-01`, `DB-02`, `DB-03`, `CORE-10`.
- [SCHEMA_FINAL.sql](../../SCHEMA_FINAL.sql): `accounts`, `users`, `regions` definitions and seed data.
- [sql/queries/users.sql](../../sql/queries/users.sql): Existing account and user profile `sqlc` queries.
- [sqlc.yaml](../../sqlc.yaml): Code generation configuration (`sql/queries/` -> `core-service/internal/db/`).

## Task

1. **Reorder `regions` Before `accounts` & Harden `accounts`/`users` in [SCHEMA_FINAL.sql](../../SCHEMA_FINAL.sql):**
   - Move `CREATE TABLE regions (...)` (lines 130–137) immediately above `CREATE TABLE accounts (...)` (line 59) so `users.anchor_dialect` can reference `regions(code)` cleanly in both PostgreSQL and `sqlc`.
   - On `accounts`:
     - Add `token_version INTEGER NOT NULL DEFAULT 1` (`CHECK (token_version >= 1)`).
     - Add `CREATE UNIQUE INDEX idx_accounts_email_lower ON accounts (LOWER(email));`.
   - On `users`:
     - Change `anchor_dialect` to `VARCHAR(10) NOT NULL DEFAULT 'MX' REFERENCES regions(code)`.
     - Add `CHECK (char_length(trim(name)) BETWEEN 2 AND 100)`.
     - Add `CHECK (user_name ~ '^[a-z0-9][a-z0-9_-]{2,49}$')`.
     - Add `CREATE UNIQUE INDEX idx_users_user_name_lower ON users (LOWER(user_name));`.
   - Add `refresh_tokens` table right after `accounts`:
     ```sql
     CREATE TABLE refresh_tokens (
       id_token UUID PRIMARY KEY DEFAULT gen_random_uuid(),
       id_account UUID NOT NULL REFERENCES accounts(id_account) ON DELETE CASCADE,
       token_hash CHAR(64) NOT NULL UNIQUE,
       expires_at TIMESTAMPTZ NOT NULL,
       revoked_at TIMESTAMPTZ,
       created_at TIMESTAMPTZ NOT NULL DEFAULT now()
     );
     CREATE INDEX idx_refresh_tokens_account ON refresh_tokens(id_account);
     ```
   - Verify the existing `INSERT INTO` seed rows in `SCHEMA_FINAL.sql` still satisfy all new constraints (`admin`, `silas_editor`, `gael_learner`, `MX`, names length `2..100`). Note: existing `admin` username is 5 chars (`^[a-z0-9][a-z0-9_-]{2,49}$` matches 3..50 chars).

2. **Update and Add Queries in [sql/queries/users.sql](../../sql/queries/users.sql):**
   - Update `CreateAccount :one` and `GetAccountByID :one` and `GetAccountByEmail :one` to also project `token_version`, and change `GetAccountByEmail` filter to `WHERE LOWER(email) = LOWER($1)`.
   - Add `GetAuthAccountStateByID :one`:
     ```sql
     -- name: GetAuthAccountStateByID :one
     SELECT id_account, role, is_active, token_version
     FROM accounts
     WHERE id_account = $1;
     ```
   - Add `GetAccountWithProfileByEmail :one` returning `a.id_account, a.email, a.password_hash, a.role, a.is_active, a.token_version, u.name, u.user_name, u.avatar_url, u.anchor_dialect, u.total_xp, u.current_level, u.current_streak, u.longest_streak, u.created_at, u.updated_at` filtered by `WHERE LOWER(a.email) = LOWER($1)`.
   - Add `IncrementAccountTokenVersion :one`:
     ```sql
     -- name: IncrementAccountTokenVersion :one
     UPDATE accounts
     SET token_version = token_version + 1, updated_at = now()
     WHERE id_account = $1
     RETURNING token_version;
     ```
   - Update `GetUserProfileByID :one` and `UpdateUserProfile :one` to include `a.is_active` and `a.token_version`.
   - Update `GetUserProfileByUsername :one` to enforce `AND a.is_active = true` so suspended profiles are hidden (`404`).
   - Add refresh token queries:
     - `CreateRefreshToken :one` (`INSERT INTO refresh_tokens (id_account, token_hash, expires_at) VALUES ($1, $2, $3) RETURNING *`)
     - `GetRefreshTokenForUpdate :one` (`SELECT id_token, id_account, token_hash, expires_at, revoked_at, created_at FROM refresh_tokens WHERE token_hash = $1 FOR UPDATE`)
     - `RevokeRefreshTokenByID :exec` (`UPDATE refresh_tokens SET revoked_at = now() WHERE id_token = $1 AND revoked_at IS NULL`)
     - `RevokeRefreshTokenByHash :exec` (`UPDATE refresh_tokens SET revoked_at = now() WHERE token_hash = $1 AND revoked_at IS NULL`)
     - `RevokeAllRefreshTokensForAccount :exec` (`UPDATE refresh_tokens SET revoked_at = now() WHERE id_account = $1 AND revoked_at IS NULL`)

3. **Regenerate `sqlc` Bindings:**
   - Run `sqlc generate` from the repository root (`Palabrario`).
   - Do not modify `core-service/internal/service/` yet (Note: `UpdateAccountPassword` or other untouched queries must stay compatible; only update callers in Ticket `031` if needed, or confirm `sqlc generate` succeeds cleanly and document any expected struct field additions in `changelog.md`). Verify `go build ./core-service/internal/db/...` exits `0`.

## Verification

Run every check in [spec.md](spec.md).

## Change log

Write `.tickets/029-backend-auth-db-schema/changelog.md` recording:
1. Exact schema changes and constraints added to `SCHEMA_FINAL.sql`.
2. New and modified queries in `sql/queries/users.sql`.
3. Output of `sqlc generate` and `go build ./core-service/internal/db/...` (plus note whether `core-service` services need any field updates in Ticket `031`).
