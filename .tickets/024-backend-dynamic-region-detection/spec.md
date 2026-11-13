# 024 — Spec checklist

Run by the router before accepting the ticket. All items must pass.

- [x] `sql/queries/studio.sql` updated with `region_names` subqueries in:
  - `ListPublicCollections`
  - `ListStudioCollections`
  - `GetCollectionByID`
  - `ListCollectionDecks`
  - `ListPublicDecks`
  - `ListStudioDecks`
  - `GetDeckByID`
- [x] Subqueries use `COALESCE(..., '{}'::TEXT[])::TEXT[] AS region_names` to guarantee non-null return.
- [x] Subqueries limit results to top 3 regions ordered by `count(*) DESC, r.name ASC`.
- [x] `sqlc generate` executed cleanly with 0 errors (`core-service/internal/db/studio.sql.go` regenerated).
- [x] `RegionNames []string` (`json:"region_names"`) added to `Collection` and `Deck` structs in `pkg/models/studio.go` and documented in `.routing/services/studio_service.md`.
- [x] Service layer in `core-service/internal/service/studio_service.go` maps `row.RegionNames` to `model.RegionNames` (non-nil `[]string{}`).
- [x] `go test -race -count=1 ./core-service/...` passes (including new/updated `RegionNames` tests in `studio_mastery_test.go`).
- [x] Mock fixtures in `app/frontend/src/api/mock/fixtures/decks.fixture.ts` populated with realistic `region_names`.
- [x] `wails3 task check` (from `app/`) generates updated TypeScript bindings (`app/frontend/bindings/`) and exits with code `0` (`bindings`, `vue-tsc --noEmit`, `vite build`, and `go build .`).
- [x] Zero flag emojis, zero country/dialect codes in mock `region_names`.
- [x] `.tickets/024-backend-dynamic-region-detection/changelog.md` exists and is populated.
- [x] No unauthorized `git commit`, `git push`, or `git merge` was run. (If committed after explicit user authorization: Conventional Commits, English, detailed body.)
