# 024 — Spec checklist

Run by the router before accepting the ticket. All items must pass.

- [ ] `sql/queries/studio.sql` updated with `region_names` subqueries in:
  - `ListPublicCollections`
  - `ListStudioCollections`
  - `GetCollectionByID`
  - `ListPublicDecks`
  - `ListStudioDecks`
  - `GetDeckByID`
- [ ] Subqueries use `COALESCE(..., '{}'::TEXT[])::TEXT[] AS region_names` to guarantee non-null return.
- [ ] Subqueries limit results to top 3 regions ordered by `count(*) DESC, r.name ASC`.
- [ ] `sqlc generate` executed cleanly with 0 errors.
- [ ] `RegionNames []string `json:"region_names"`` added to `Collection` and `Deck` structs in `pkg/models/studio.go`.
- [ ] Service layer in `core-service/internal/service/studio/` maps `row.RegionNames` to `model.RegionNames`.
- [ ] Mock fixtures in `app/frontend/src/api/mock/fixtures/decks.fixture.ts` populated with realistic `region_names`.
- [ ] `cd app && wails3 task check` generates updated TypeScript bindings and exits with code 0 (TypeScript, Vite, and Go builds all pass).
- [ ] Zero flag emojis, zero country codes in region names.
- [ ] `.tickets/024-backend-dynamic-region-detection/changelog.md` exists and is populated.
- [ ] Commit (if made): Conventional Commits, English, detailed body.
