# Changelog — Ticket 014: Backend Dedicated FTS for Decks & Collections and Expanded Word Search

## Summary

Implemented full-text search (FTS) infrastructure for study decks and thematic collections, and expanded the `WordSearchResult` domain DTO to include lexical metadata (`part_of_speech` and dynamic `region_name`). Regenerated SQL queries via `sqlc`, updated domain services and HTTP routes, refreshed Wails v3 bindings, and verified 100% test pass rate on ephemeral PostgreSQL 17 (:5433).

## Key Changes

### Database Schema (`SCHEMA_FINAL.sql`)
- Added `search_vector tsvector GENERATED ALWAYS AS (setweight(to_tsvector('spanish_unaccent', coalesce(title, '')), 'A') || setweight(to_tsvector('spanish_unaccent', coalesce(description, '')), 'B')) STORED` to `decks` and `collections`.
- Created GIN index `idx_decks_search_vector ON decks USING GIN(search_vector)`.
- Created GIN index `idx_collections_search_vector ON collections USING GIN(search_vector)`.

### SQL Queries (`sql/queries/`)
- Updated `SearchVocabularyDropdown` in `dictionary.sql`:
  - Selected `w.part_of_speech` in `matched_words` CTE.
  - Selected `mw.part_of_speech` and `coalesce(r.name, 'General')::VARCHAR AS region_name` in `distinct_words` CTE.
  - Returned `dw.part_of_speech` and `dw.region_name` in outer query.
- Added `SearchDecksDropdown` in `studio.sql`:
  - Computes `region_name` dynamically from associated cards via `COALESCE((SELECT r.name FROM deck_items di JOIN senses s ON s.id_sense = di.id_sense JOIN regions r ON r.id_region = s.id_region WHERE di.id_deck = d.id_deck GROUP BY r.name ORDER BY count(*) DESC LIMIT 1), 'General')::VARCHAR AS region_name`.
  - Returns `u.user_name AS owner_user_name`.
  - Filters `d.is_public = true` with unaccent prefix and FTS vector ranking (`ts_rank`).
- Added `SearchCollectionsDropdown` in `studio.sql`:
  - Returns `u.user_name AS owner_user_name`, `deck_count`, `item_count`, `is_official`, `is_featured`.
  - Filters `c.is_public = true` with unaccent prefix and FTS vector ranking (`ts_rank`).
- Ran `sqlc generate` to regenerate `core-service/internal/db/`.

### Domain Models (`pkg/models/`)
- In `pkg/models/dictionary.go`:
  - Added `PartOfSpeech string` and `RegionName string` to `WordSearchResult`.
- In `pkg/models/studio.go`:
  - Added `DeckSearchResult` struct (`id_deck`, `title`, `description`, `region_name`, `owner_user_name`, `item_count`, `is_official`).
  - Added `CollectionSearchResult` struct (`id_collection`, `title`, `description`, `owner_user_name`, `deck_count`, `item_count`, `is_official`, `is_featured`).

### Service Layer & HTTP Handlers
- `core-service/internal/service/dictionary_service.go`:
  - Updated `SearchDropdown` to map `PartOfSpeech` and `RegionName` to `models.WordSearchResult`.
- `core-service/internal/service/studio_service.go`:
  - Added `SearchDecks` and `SearchCollections` to `StudioService` interface and implemented both methods with validation, limits, and mapping.
- `core-service/internal/api/handlers/studio_handler.go`:
  - Added `SearchDecks` (`GET /api/v1/decks/search`) and `SearchCollections` (`GET /api/v1/collections/search`).
- `core-service/internal/api/router.go`:
  - Mounted `/decks/search` and `/collections/search` before `/{id}` wildcard routes.

### Wails v3 Desktop Bridge (`app/`)
- Added `SearchDecks(q string, limit int) ([]models.DeckSearchResult, error)` to `app/app_studio.go`.
- Added `SearchCollections(q string, limit int) ([]models.CollectionSearchResult, error)` to `app/app_studio.go`.
- Executed `task client:check` (`wails3 task check`), regenerating TypeScript bindings cleanly in `app/frontend/bindings/` with 0 type-check or build errors.

### Testing & Verification
- Updated `core-service/internal/service/service_test.go`:
  - Verified `SearchDropdown` returns populated `PartOfSpeech` and `RegionName`.
- Added tests in `core-service/internal/service/studio_mastery_test.go`:
  - `TestStudioService_SearchDecks_PostgresFTS`: prefix matching, accent insensitivity (`mex` vs `méx`), description FTS, privacy filter (`is_public = false`), empty query handling.
  - `TestStudioService_SearchCollections_PostgresFTS`: prefix matching, accent insensitivity (`gastro` vs `gastronomía`), description FTS (`empanadas`), privacy filter, empty query handling.
- Executed `task test` against ephemeral PostgreSQL (:5433) in tmpfs with 100% pass rate.

### Service Documentation (`.routing/services/`)
- Updated `.routing/services/dictionary_service.md` with expanded `WordSearchResult` fields.
- Updated `.routing/services/studio_service.md` with `DeckSearchResult` and `CollectionSearchResult` DTOs, interface methods, HTTP endpoints, and `sqlc` queries.
