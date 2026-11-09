 014 — Backend: Dedicated FTS for Decks & Collections and Expanded Word Search

| Field | Value |
| :--- | :--- |
| **ID** | 014 |
| **Status** | draft |
| **Model** | Sonnet 5.5 (`claude-sonnet-5-5`, Antigravity provider) |
| **Role** | Backend Implementer |
| **Spec** | [spec.md](spec.md) |

> **Draft Note:**  
> This backend ticket upgrades the database schema, sqlc queries, and Go services to provide true Full-Text Search (FTS) for study decks and collections, while expanding the `WordSearchResult` contract so the search combobox receives rich lexical metadata directly from the backend. It also updates and adds automated service tests in `service_test.go` running on the RAM PostgreSQL harness (`:5433`).

## Scope

May create or edit:
- `SCHEMA_FINAL.sql` (add `search_vector tsvector` generated column and GIN indexes to `decks` and `collections`)
- `sql/queries/dictionary.sql` (update `SearchVocabularyDropdown` to return `part_of_speech` and `region_name`)
- `sql/queries/studio.sql` (add `SearchDecksDropdown` and `SearchCollectionsDropdown` with FTS ranking)
- `core-service/internal/db/` (regenerate sqlc)
- `pkg/models/dictionary.go` (expand `WordSearchResult` with `PartOfSpeech` and `RegionName`)
- `pkg/models/studio.go` (add `DeckSearchResult` / `CollectionSearchResult`)
- `core-service/internal/service/` (implement service methods in `dictionary_service.go`, `studio_service.go`)
- `core-service/internal/service/service_test.go` (update/add tests for expanded search and new FTS queries)
- `core-service/internal/service/studio_mastery_test.go` (add unit/integration tests for deck & collection search)
- `core-service/internal/api/handlers/` (expose/update endpoints in `dictionary_handler.go`, `studio_handler.go`)
- `app/` (Wails App bridge methods in `app_dictionary.go` / `app_studio.go` and regenerate bindings)
- `.routing/services/dictionary_service.md` (document expanded `WordSearchResult` fields)
- `.routing/services/studio_service.md` (document `DeckSearchResult` / `CollectionSearchResult` DTOs and endpoints)
- `.tickets/014-backend-search-fts/changelog.md` (create)

Must not touch:
- Frontend UI components (`app/frontend/src/`)
- `.context/`
- Any file outside the above scope

## Inputs

- [.routing/services/dictionary_service.md](../../.routing/services/dictionary_service.md): Lexicon, senses, anchor dialect, and dictionary service contracts (`CU05–CU06`).
- [.routing/services/studio_service.md](../../.routing/services/studio_service.md): Study decks, collections, and studio service contracts (`CU10–CU11`).
- [`SCHEMA_FINAL.sql`](../../SCHEMA_FINAL.sql): Base schema for `decks` and `collections`.
- [`core-service/internal/service/service_test.go`](../../core-service/internal/service/service_test.go): Test harness and existing dictionary tests.

## Task

1. **Database Schema Enhancements (`SCHEMA_FINAL.sql`):**
   - In `decks`: Add `search_vector tsvector GENERATED ALWAYS AS (setweight(to_tsvector('spanish_unaccent', coalesce(title, '')), 'A') || setweight(to_tsvector('spanish_unaccent', coalesce(description, '')), 'B')) STORED`.
   - Add GIN index on `decks(search_vector)`.
   - In `collections`: Add `search_vector tsvector GENERATED ALWAYS AS (setweight(to_tsvector('spanish_unaccent', coalesce(title, '')), 'A') || setweight(to_tsvector('spanish_unaccent', coalesce(description, '')), 'B')) STORED`.
   - Add GIN index on `collections(search_vector)`.

2. **SQL Queries (`sql/queries/`):**
   - In `dictionary.sql`: Update `SearchVocabularyDropdown` to include `w.part_of_speech` and `r.name AS region_name`.
   - In `studio.sql`: Add `SearchDecksDropdown` returning matching public decks ranked by prefix match and FTS vector.
   - In `studio.sql`: Add `SearchCollectionsDropdown` returning matching public collections ranked by prefix match and FTS vector.
   - Run `sqlc generate`.

3. **Domain Models & Service Layer:**
   - In `pkg/models/dictionary.go`: Add `PartOfSpeech string` and `RegionName string` to `WordSearchResult`.
   - In `pkg/models/studio.go`: Add `DeckSearchResult` (id, title, region_name, owner_username, item_count, is_official).
   - In `core-service`: Update `DictionaryService.SearchDropdown` and `StudioService.SearchDecks`.
   - Update HTTP handler `GET /api/v1/words/search` and expose `GET /api/v1/decks/search`.

4. **Wails Desktop Bridge & Bindings:**
   - Expose methods in `app/app_dictionary.go` and `app/app_studio.go`.
   - Run `wails3 task client:bindings` to generate updated TypeScript bindings under `app/frontend/bindings/`.

5. **Automated Testing Suite (`service_test.go`):**
   - Update `service_test.go` to assert that `SearchDropdown` returns `part_of_speech` and `region_name`.
   - Add dedicated test cases in `service_test.go` (or `studio_mastery_test.go`) covering FTS search on decks and collections with accent-insensitive and prefix matching.
   - Run `task test` to verify all tests pass 100% against the RAM PostgreSQL harness (`:5433`).

6. **Validation:**
   - Run `wails3 task check` from `app/` (exiting with code 0).

7. **Update Service Specifications (`.routing/services/`):**
   - In `dictionary_service.md`: Document `part_of_speech` and `region_name` on `WordSearchResult`.
   - In `studio_service.md`: Document `DeckSearchResult` and `CollectionSearchResult` DTOs and search endpoints.

8. **Change log:**
   - Document all changes and test outputs in `.tickets/014-backend-search-fts/changelog.md`.

## Verification

All checks are listed in [spec.md](spec.md).
