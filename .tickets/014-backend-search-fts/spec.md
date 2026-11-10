# 014 — Spec checklist

Run by the router before accepting the ticket. All items must pass.

- [x] `SCHEMA_FINAL.sql` includes `search_vector` and GIN indexes for `decks` and `collections`.
- [x] `sqlc generate` succeeds without errors.
- [x] `WordSearchResult` includes `part_of_speech` and `region_name`.
- [x] `SearchDecksDropdown` query and service method implemented.
- [x] Wails v3 bindings regenerated cleanly with the updated models.
- [x] Automated tests in `core-service/internal/service/service_test.go` cover expanded word search and deck/collection FTS.
- [x] `task test` passes with 0 errors on RAM PostgreSQL (:5433).
- [x] `wails3 task check` in `app/` exits with code 0.
- [x] `.routing/services/dictionary_service.md` and `studio_service.md` updated with the new search contracts.
- [x] Changes confined strictly to the backend, schema, tests, bindings, and service docs scope.
- [x] `.tickets/014-backend-search-fts/changelog.md` exists and is populated.
- [x] Commit (if made): Conventional Commits, English, detailed body.
