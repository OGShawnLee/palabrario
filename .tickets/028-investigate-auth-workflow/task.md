# 028 — End-to-End Auth & Session Workflow Audit

| Field | Value |
| :--- | :--- |
| **ID** | 028 |
| **Status** | completed |
| **Model** | Gemini 4 Hydrogen, `thinking_level: high` (`gemini-4-hydrogen`) |
| **Role** | Investigator (`research`) |
| **Spec** | [spec.md](spec.md) |

> **Audit Note:**  
> This ticket runs a read-only end-to-end investigation of Palabrario's authentication, JWT/session lifecycle, authorization middleware, error handling, test flakiness, and client contracts—starting from PostgreSQL and `core-service/` through `realtime-service/`, `image-service/`, the Wails Go BFF (`app/`), and the frontend (`app/frontend/src/`). Its output (`report.md`) will drive the full-stack refactor ticket(s) that follow.

## Scope

May create or edit:
- `.tickets/028-investigate-auth-workflow/report.md` (create: structured findings, flaky test audit, and recommended refactor plan)

Must not touch:
- `core-service/`, `realtime-service/`, `image-service/`, `pkg/`, `sql/`, `SCHEMA_FINAL.sql` (read-only)
- `app/`, `app/frontend/src/`, `app/frontend/uno.config.ts` (read-only)
- `.routing/`, `.context/`, `.claude/` (read-only)
- Any other file in the repository
- **No role (Router, Investigator, Implementer, Translator, or Reviewer) may commit, push, or merge without explicit user authorization.** Do not run `git commit`, `git push`, `git merge`, or mutate git state.

## Inputs

- [.routing/CONTEXT.md](../../.routing/CONTEXT.md): Ground rules, stack, agent roles, and commit-authorization policy.
- [.routing/API.md](../../.routing/API.md): Frontend adapter architecture (`IAuthAPI`, Wails, Mock, and future Capacitor/HTTP considerations).
- [.routing/services/identity_service.md](../../.routing/services/identity_service.md): Documented Identity & Auth contract.
- [SCHEMA_FINAL.sql](../../SCHEMA_FINAL.sql): Database tables, constraints, and triggers for `accounts` and `users`.
- [sql/queries/users.sql](../../sql/queries/users.sql): sqlc queries for auth, profile lookup, and updates.
- [pkg/auth/auth.go](../../pkg/auth/auth.go), [pkg/auth/auth_test.go](../../pkg/auth/auth_test.go), [pkg/models/auth.go](../../pkg/models/auth.go), and [pkg/models/user.go](../../pkg/models/user.go): Token signing/verification, bcrypt, unit tests, and shared DTO models.
- [core-service/internal/config/config.go](../../core-service/internal/config/config.go), [core-service/internal/api/router.go](../../core-service/internal/api/router.go), [core-service/internal/api/middleware/auth_middleware.go](../../core-service/internal/api/middleware/auth_middleware.go), [core-service/internal/api/handlers/auth_handler.go](../../core-service/internal/api/handlers/auth_handler.go), [core-service/internal/api/handlers/user_handler.go](../../core-service/internal/api/handlers/user_handler.go), [core-service/internal/service/identity_service.go](../../core-service/internal/service/identity_service.go), [core-service/internal/service/errors.go](../../core-service/internal/service/errors.go), [core-service/internal/testutil/db.go](../../core-service/internal/testutil/db.go), [core-service/internal/api/auth_test.go](../../core-service/internal/api/auth_test.go), and [core-service/internal/service/service_test.go](../../core-service/internal/service/service_test.go).
- [realtime-service/internal/api/handlers/arena_handler.go](../../realtime-service/internal/api/handlers/arena_handler.go), [realtime-service/internal/api/api_test.go](../../realtime-service/internal/api/api_test.go), [image-service/internal/api/middleware/auth_middleware.go](../../image-service/internal/api/middleware/auth_middleware.go), and [image-service/internal/api/middleware/auth_middleware_test.go](../../image-service/internal/api/middleware/auth_middleware_test.go).
- [app/app.go](../../app/app.go), [app/app_auth.go](../../app/app_auth.go), and [app/app_arena.go](../../app/app_arena.go).
- [app/frontend/src/api/contract/auth.contract.ts](../../app/frontend/src/api/contract/auth.contract.ts), [app/frontend/src/api/live/wails/auth.wails.ts](../../app/frontend/src/api/live/wails/auth.wails.ts), [app/frontend/src/api/mock/mock.adapter.ts](../../app/frontend/src/api/mock/mock.adapter.ts), [app/frontend/src/composables/useAuth.ts](../../app/frontend/src/composables/useAuth.ts), [app/frontend/src/components/TopBar.vue](../../app/frontend/src/components/TopBar.vue), [app/frontend/src/views/auth/SignInView.vue](../../app/frontend/src/views/auth/SignInView.vue), [app/frontend/src/views/auth/SignUpView.vue](../../app/frontend/src/views/auth/SignUpView.vue), and [app/frontend/src/router/index.ts](../../app/frontend/src/router/index.ts).

## Task

1. **Audit Database Schema & Queries (`SCHEMA_FINAL.sql`, `sql/queries/users.sql`):**
   - Check column types, constraints, uniqueness indexes (`email`, `user_name`), check constraints (`anchor_dialect`, role, lengths), triggers, and `is_active` enforcement across `CreateAccount`, `CreateUserProfile`, `GetAccountByEmail`, `GetUserProfileByID`, `GetUserProfileByUsername`, and `UpdateUserProfile`.
2. **Audit Core Auth, Middleware, Handlers, and Services (`pkg/auth/`, `core-service/`):**
   - Inspect `SignUp`, `SignIn`, `GetMe`, `UpdateProfile`, `ChangePassword`, and `GetPublicProfileByUsername`.
   - Verify input validation symmetry (password minimum length, username format, optional vs. required `name`, `anchor_dialect` validation), SQL/constraint error mapping in `errors.go`, `is_active` checks on authenticated requests vs. `SignIn`, `OptionalAuth` behavior on invalid/expired tokens vs. missing tokens, CORS configuration (`AllowOrigins: ["*"]` with `AllowCredentials: true`), and error message format/language consistency for UI consumers.
   - Check token issuance, claims (`Role`, `AnchorDialect`, expiration), and whether `realtime-service` and `image-service` align with `core-service` auth middleware and WebSocket token transport.
3. **Audit Wails BFF & Frontend Auth Layer (`app/`, `app/frontend/src/`):**
   - Inspect `App.sessionToken`, `doRequest` (401 handling, error extraction, network errors), `HasActiveSession`, `SignIn`, `SignUp`, `SignOut`, `GetMe`, `GetMyProfile`, `GetPublicProfileByUsername`, and `ConnectArenaRoom`.
   - Inspect `IAuthAPI`, `WailsAuthAdapter`, `MockAuthAPI`, `useAuth`, `TopBar`, `SignInView`, `SignUpView`, and future Capacitor HTTP adapter needs: identify redundant types/methods (e.g., `AuthResponse`, `hasActiveSession()`, `GetMyProfile`), error/null mismatches between Mock and Wails, and session-persistence/hydration behavior across desktop, web mock, and mobile.
4. **Search for Flaky Auth Testing & Coverage Gaps:**
   - Audit `pkg/auth/auth_test.go`, `core-service/internal/api/auth_test.go`, `core-service/internal/api/governance_reports_test.go`, `core-service/internal/service/service_test.go`, `core-service/internal/testutil/db.go`, `realtime-service/internal/api/api_test.go`, `image-service/internal/api/middleware/auth_middleware_test.go`, and `app/` tests (if any).
   - Actively look for **flaky test patterns**: wall-clock timing / `time.Sleep` / sub-second JWT `exp`/`nbf` boundary races, parallel test DB/template collisions (`t.Parallel()` vs. shared `testutil.SetupTestDB` / `os.Setenv`), silent skips when env vars/DB are unset, hardcoded seed credentials/UUIDs, non-deterministic ordering, or uncleaned goroutines/HTTP servers.
   - List critical auth scenarios that have zero test coverage today (e.g., suspended user calling `/me`, expired/malformed JWT, concurrent `SignUp` conflicts, `UpdateProfile`/`ChangePassword`, `401` handling in `app/`).
   - Run read-only Go checks (e.g., `go test ./pkg/auth/...`, `go test ./image-service/internal/api/middleware/...`, and inspect whether `core-service` auth tests run or skip in this environment) without modifying files. Record exact command output.
5. **Write `.tickets/028-investigate-auth-workflow/report.md`:**
   - Organize findings by **Severity (`Critical` / `High` / `Medium` / `Low`)** and **Layer (`Database` / `core-service` / `realtime & image` / `Wails BFF` / `Frontend` / `Tests & Flakiness`)** with exact `file:line` citations.
   - Include a dedicated **"Flaky Auth Tests & Coverage Gaps"** section detailing every flaky pattern found (or proven absent) and untested failure modes.
   - Include a **"Remove / Simplify"** table listing redundant functions, types, workarounds, and dead code across backend, BFF, and frontend.
   - Propose an ordered **Refactor Ticket Breakdown** (e.g., Ticket `029`: `core-service` + `pkg` + `sql` + `app/` hardening & deterministic auth tests; Ticket `030`: `IAuthAPI` + `useAuth` + `MockAuthAPI` + view error handling; Reviewer pass) with exact file scopes and acceptance criteria ready to turn into tickets.

## Verification

All checks are listed in [spec.md](spec.md).

## Change log

Record all verified findings, redundant abstractions, test observations, and proposed implementation tickets directly in `.tickets/028-investigate-auth-workflow/report.md` (no separate `changelog.md` needed since no source code is modified).
