# 014 — Spec checklist

Run by the router before accepting the ticket. All items must pass.

- [ ] `SCHEMA_FINAL.sql` includes `search_vector` and GIN indexes for `decks` and `collections`.
- [ ] `sqlc generate` succeeds without errors.
- [ ] `WordSearchResult` includes `part_of_speech` and `region_name`.
- [ ] `SearchDecksDropdown` query and service method implemented.
- [ ] Wails v3 bindings regenerated cleanly with the updated models.
- [ ] Automated tests in `core-service/internal/service/service_test.go` cover expanded word search and deck/collection FTS.
- [ ] `task test` passes with 0 errors on RAM PostgreSQL (:5433).
- [ ] `wails3 task check` in `app/` exits with code 0.
- [ ] `.routing/services/dictionary_service.md` and `studio_service.md` updated with the new search contracts.
- [ ] Changes confined strictly to the backend, schema, tests, bindings, and service docs scope.
- [ ] `.tickets/014-backend-search-fts/changelog.md` exists and is populated.
- [ ] Commit (if made): Conventional Commits, English, detailed body.
