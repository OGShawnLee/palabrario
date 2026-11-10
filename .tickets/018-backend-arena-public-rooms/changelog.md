# 018 — Changelog: Backend Public Arena Rooms Listing Endpoint

## Summary
Implemented the backend endpoint and desktop bridge method to list public, open Arena matches currently waiting for players (`status = 'WAITING' AND is_public = true`), with participant counts < 8, filtering out stale rooms older than 2 hours.

## Changes Implemented

1. **SQL Query (`sql/queries/arena.sql`):**
   - Added typed query `ListPublicWaitingMatches`:
     - Selects match fields (`id_match`, `room_code`, `id_deck`, `total_rounds`, `round_time_sec`, `created_at`).
     - Joins `decks` for `deck_title` and `is_deck_official`.
     - Left joins `users` for `host_user_name` (coalesced to `'Unknown'`).
     - Calculates `participant_count` using subquery count on `match_participants`.
     - Filters by `status = 'WAITING'`, `is_public = true`, `created_at >= now() - INTERVAL '2 hours'`, and `participant_count < 8`.
     - Orders by `m.created_at DESC` with configurable `LIMIT $1`.
   - Executed `sqlc generate` producing `db.ListPublicWaitingMatchesRow` and `Queries.ListPublicWaitingMatches` in `core-service/internal/db/arena.sql.go`.

2. **DTO Definition (`pkg/models/arena.go`):**
   - Defined `PublicMatchItem` with typed JSON attributes in `snake_case`:
     - `id_match`, `room_code`, `id_deck`, `deck_title`, `is_deck_official`, `host_user_id` (pointer, omitempty), `host_user_name`, `total_rounds`, `round_time_sec`, `created_at` (`time.Time`), `participant_count` (`int`).

3. **Service Layer (`core-service/internal/service/arena_catalog_service.go`):**
   - Added `ListPublicWaitingMatches(ctx context.Context, limit int) ([]models.PublicMatchItem, error)` to interface and struct implementation.
   - Defaulted limit to 20 if `limit <= 0 || limit > 50`.
   - Mapped `db.ListPublicWaitingMatchesRow` to `models.PublicMatchItem`.

4. **HTTP Handler & Routing (`core-service/internal/api/`):**
   - In `handlers/arena_handler.go`: Added `ListPublicMatches(w http.ResponseWriter, r *http.Request)` parsing `?limit=` (default 20, min 1, max 50).
   - In `router.go`: Mounted public endpoint `GET /api/v1/arena/matches/public` mapped to `arenaHandler.ListPublicMatches`.

5. **Wails Desktop Bridge (`app/app_arena.go`):**
   - Added bridge method `ListPublicWaitingMatches(limit int) ([]models.PublicMatchItem, error)`.
   - Executed `task client:check` to regenerate Wails TS bindings and verify frontend build (`vue-tsc` and `vite build`).

6. **Automated Integration Tests (`core-service/internal/service/service_test.go`):**
   - Added `TestArenaCatalogService_ListPublicWaitingMatches`:
     - Verified public waiting match is returned with full metadata and correct participant count.
     - Verified exclusion of private matches (`is_public = false`).
     - Verified exclusion of active matches (`status = 'IN_PROGRESS'`).
     - Verified exclusion of expired rooms (> 2 hours).
     - Verified exclusion of full rooms (>= 8 participants).
     - Verified limit parameter behavior (`limit = 1`).

7. **Documentation Updates:**
   - Updated `.routing/services/arena_service.md` Sections 2.7.2 (DTO), 2.7.3 (interface), 2.7.5 (endpoints table), and 2.7.6 (`sqlc` query list).

## Verification Results
- `sqlc generate`: 0 errors.
- `task test`: All unit & integration tests passed against PostgreSQL in RAM (:5433).
- `task client:check`: Wails TS bindings generated cleanly; TypeScript type-check and Vite production build succeeded with exit code 0.
