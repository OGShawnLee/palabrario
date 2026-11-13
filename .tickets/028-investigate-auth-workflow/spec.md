# 028 — Spec checklist

Run by the router before accepting the ticket. All items must pass.

- [x] `.tickets/028-investigate-auth-workflow/report.md` exists and is non-empty.
- [x] `report.md` covers all 5 layers with exact file and line citations: PostgreSQL/sqlc (`SCHEMA_FINAL.sql`, `sql/queries/users.sql`), `pkg/auth` + `core-service`, sibling services (`realtime-service`, `image-service`), Wails BFF (`app/`), and frontend (`app/frontend/src/api/`, `useAuth.ts`, auth views, `TopBar.vue`).
- [x] Every finding in `report.md` has a severity rating, verified evidence from source files, concrete failure mode / UX impact, and specific fix recommendation.
- [x] `report.md` includes a **"Flaky Auth Tests & Coverage Gaps"** section auditing `pkg/auth/auth_test.go`, `core-service/internal/api/auth_test.go`, `core-service/internal/service/service_test.go`, `core-service/internal/testutil/db.go`, `realtime-service/internal/api/api_test.go`, and `image-service/internal/api/middleware/auth_middleware_test.go` for timing, concurrency, env/DB setup, seed, or silent-skip flakiness, plus exact `go test` results.
- [x] `report.md` includes a **"Remove / Simplify"** section identifying redundant methods, DTO leaks, and inconsistencies (`AuthResponse`, `hasActiveSession()`, `GetMyProfile`, `MockAuthAPI` vs. `WailsAuthAdapter` error/null divergence, etc.).
- [x] `report.md` includes a concrete, ordered **Refactor Ticket Breakdown** with exact file scopes separating `core-service` + `app/` changes from frontend contract/composable changes.
- [x] Zero source files, routing files, or `uno.config.ts` were modified (`git diff --name-only` is empty outside `.tickets/028-investigate-auth-workflow/report.md`).
- [x] Zero unauthorized `git commit`, `git push`, or `git merge` actions were performed.
