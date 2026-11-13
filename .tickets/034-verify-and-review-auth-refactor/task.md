# 034 — End-to-End Auth Refactor Verification and Security Review

| Field | Value |
| :--- | :--- |
| **ID** | 034 |
| **Status** | completed |
| **Model** | Gemini 4 Hydrogen, `thinking_level: high` (`gemini-4-hydrogen`) |
| **Role** | Auth Refactor Reviewer (`review`) |
| **Spec** | [spec.md](spec.md) |

> **Specification Note:**  
> Final read-only Reviewer gate across Tickets `029`–`033`. Audits the complete git diff across database, `pkg/auth`, `pkg/response`, `core-service`, `realtime-service`, `image-service`, Wails BFF (`app/`), Vue frontend (`app/frontend/src/`), and `.routing/` docs against every finding in [028 report](../028-investigate-auth-workflow/report.md) and the user's confirmed architectural rules (Option B active-account enforcement on writes + `/auth/*` + `/governance/*`, hidden suspended profiles, short-lived `15m` access tokens + persistent host refresh tokens, English error responses with `code`).

## Scope

May create or edit:
- `.tickets/034-verify-and-review-auth-refactor/review.md` (create: structured pass/fail review findings only)

Must not touch:
- Any repository source file (`SCHEMA_FINAL.sql`, `sql/`, `pkg/`, `core-service/`, `realtime-service/`, `image-service/`, `app/`, `.routing/`, `.context/`)
- **Reviewer strictly audits and runs read-only verification commands.** If any blocker is found, record it in `review.md` and return to the Router—do not patch source code directly.
- **No role (Router, Investigator, Implementer, Translator, or Reviewer) may commit, push, or merge without explicit user authorization.** Do not run `git commit`, `git push`, or `git merge`.

## Inputs

- [.tickets/028-investigate-auth-workflow/report.md](../028-investigate-auth-workflow/report.md)
- [.tickets/029-backend-auth-db-schema/changelog.md](../029-backend-auth-db-schema/changelog.md)
- [.tickets/030-harden-pkg-auth-and-sibling-services/changelog.md](../030-harden-pkg-auth-and-sibling-services/changelog.md)
- [.tickets/031-harden-core-service-auth/changelog.md](../031-harden-core-service-auth/changelog.md)
- [.tickets/032-harden-wails-bff-auth/changelog.md](../032-harden-wails-bff-auth/changelog.md)
- [.tickets/033-align-frontend-auth-and-guards/changelog.md](../033-align-frontend-auth-and-guards/changelog.md)
- [.routing/CONTEXT.md](../../.routing/CONTEXT.md), [.routing/DOMAIN.md](../../.routing/DOMAIN.md), [.routing/services/](../../.routing/services/), [.routing/API.md](../../.routing/API.md)

## Task

1. **Audit Diff Against Security & Business Rules (Read-Only):**
   - Inspect `git diff` across all layers and confirm:
     - Zero hardcoded `JWT_SECRET` fallbacks in `core-service`, `realtime-service`, and `image-service`.
     - Zero `AllowedOrigins: ["*"]` paired with `AllowCredentials: true`, and `hub.go` validates WebSocket origins.
     - `core-service` enforces DB `is_active`, `token_version`, and live `role` on all `POST`/`PUT`/`DELETE` routes, `/api/v1/auth/*`, and `/api/v1/governance/*` (Option B), while trusting `15m` access JWTs on non-governance `GET`s.
     - `GetUserProfileByUsername` hides suspended accounts (`404`).
     - `SignIn`/`SignUp`/`RefreshToken`/`SignOut` cleanly rotate and revoke hashed `refresh_tokens`; `ChangePassword` and `is_active = false` bump `token_version` and revoke all refresh tokens.
     - `app/` persists refresh token with `0600` permissions, performs single-flight auto-refresh on `401`, emits `auth:expired` + clears disk/RAM on `403 ACCOUNT_SUSPENDED` or failed refresh, passes WS JWT in `Authorization` header, and exposes zero duplicate alias methods.
     - `IAuthAPI`, `MockAuthAPI`, and `WailsAuthAdapter` share identical error/return semantics; router guards hydrate session on boot (`init()`) for all routes including `/auth/*`.
     - All auth/validation/middleware error responses use English + machine-readable `code`.

2. **Run Full Verification Suites:**
   - Run `go test -race -count=3 ./pkg/... ./realtime-service/... ./image-service/...` from repo root.
   - Run `go test -race -count=1 ./...` inside `app/`.
   - With `postgres-test (:5433)` active (`task test:up`), run `task test:fast` (or `go test -race -count=1 ./core-service/...`) and clean up (`task test:down`) if started only for verification.
   - Run `wails3 task check` inside `app/`.

3. **Write Review Verdict:**
   - Save findings, command outputs, and either **`APPROVED`** or **`BLOCKED` (with exact file:line findings)** to `.tickets/034-verify-and-review-auth-refactor/review.md`.

## Verification

Run every check in [spec.md](spec.md).

## Change log

Record the full security & quality review verdict in `.tickets/034-verify-and-review-auth-refactor/review.md` (do not create `changelog.md` or edit any source file).
