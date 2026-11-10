# 018 — Backend: Public Active Arena Matches Listing Endpoint

| Field | Value |
| :--- | :--- |
| **ID** | 018 |
| **Status** | completed |
| **Model** | Gemini 3.8 Flash (High) (`gemini-3.8-flash`, Antigravity provider) |
| **Role** | Backend Implementer |
| **Spec** | [spec.md](spec.md) |

> **Specification Note:**  
> This backend ticket adds an endpoint to list open, public arena rooms currently waiting for players (`status = 'WAITING' AND is_public = true`). It joins deck metadata, host username (`users.user_name`), and current participant count (`COUNT(*) < 8`), filtering out stale/abandoned rooms older than 2 hours so the frontend lobby displays active, joinable games.

## Scope

May create or edit:
- `sql/queries/arena.sql` (add query `ListPublicWaitingMatches`)
- `core-service/internal/db/` (regenerate via `sqlc generate`)
- `pkg/models/arena.go` (add `PublicMatchItem` DTO)
- `core-service/internal/service/arena_catalog_service.go` (expose `ListPublicWaitingMatches` on interface and service)
- `core-service/internal/api/handlers/arena_handler.go` (expose handler `ListPublicMatches`)
- `core-service/internal/api/router.go` (mount `GET /api/v1/arena/matches/public`)
- `core-service/internal/service/service_test.go` (add automated integration test)
- `app/app_arena.go` (expose desktop bridge method and regenerate bindings)
- `.routing/services/arena_service.md` (document `PublicMatchItem`, interface method, and endpoint)
- `.tickets/018-backend-arena-public-rooms/changelog.md` (create)

Must not touch:
- `app/frontend/src/` (frontend UI components)
- `realtime-service/`, `image-service/`, `.context/`
- Any file outside the above declared scope

## Inputs

- [.routing/services/arena_service.md](../../.routing/services/arena_service.md): Arena domain rules (`RN10`, `RN13`), DTO contracts, and endpoints.
- [`SCHEMA_FINAL.sql`](../../SCHEMA_FINAL.sql): Base schema for `matches` (`id_match`, `host_user_id`, `room_code`, `status`, `is_public`), `decks`, and `users` (`user_name`).
- [`core-service/internal/service/service_test.go`](../../core-service/internal/service/service_test.go): Test harness with RAM PostgreSQL (:5433).

## Task

1. **SQL Query (`sql/queries/arena.sql`):**
   - Add query `ListPublicWaitingMatches`:
     ```sql
     -- name: ListPublicWaitingMatches :many
     SELECT 
         m.id_match,
         m.room_code,
         m.id_deck,
         d.title AS deck_title,
         d.is_official AS is_deck_official,
         m.host_user_id,
         COALESCE(u.user_name, 'Unknown')::VARCHAR AS host_user_name,
         m.total_rounds,
         m.round_time_sec,
         m.created_at,
         (SELECT count(*)::INT FROM match_participants mp WHERE mp.id_match = m.id_match) AS participant_count
     FROM matches m
     JOIN decks d ON d.id_deck = m.id_deck
     LEFT JOIN users u ON u.id_user = m.host_user_id
     WHERE m.status = 'WAITING' 
       AND m.is_public = true
       AND m.created_at >= now() - INTERVAL '2 hours'
       AND (SELECT count(*)::INT FROM match_participants mp WHERE mp.id_match = m.id_match) < 8
     ORDER BY m.created_at DESC
     LIMIT $1;
     ```
   - Run `sqlc generate`.

2. **DTO in `pkg/models/arena.go`:**
   - Define `PublicMatchItem`:
     ```go
     type PublicMatchItem struct {
         IDMatch          string    `json:"id_match"`
         RoomCode         string    `json:"room_code"`
         IDDeck           int       `json:"id_deck"`
         DeckTitle        string    `json:"deck_title"`
         IsDeckOfficial   bool      `json:"is_deck_official"`
         HostUserID       *string   `json:"host_user_id,omitempty"`
         HostUserName     string    `json:"host_user_name"`
         TotalRounds      int       `json:"total_rounds"`
         RoundTimeSec     int       `json:"round_time_sec"`
         CreatedAt        time.Time `json:"created_at"`
         ParticipantCount int       `json:"participant_count"`
     }
     ```

3. **Service Layer (`core-service/internal/service/arena_catalog_service.go`):**
   - Add method to `ArenaCatalogService` interface:
     ```go
     ListPublicWaitingMatches(ctx context.Context, limit int) ([]models.PublicMatchItem, error)
     ```
   - Implement `ListPublicWaitingMatches` defaulting limit to 20 if `limit <= 0 || limit > 50`.

4. **HTTP Route & Handler (`core-service/internal/api/`):**
   - In `handlers/arena_handler.go`: Add `ListPublicMatches(w http.ResponseWriter, r *http.Request)`.
   - In `router.go`: Mount `r.Get("/matches/public", arenaHandler.ListPublicMatches)` under `/api/v1/arena` (public endpoint).

5. **Wails Desktop Bridge (`app/app_arena.go`):**
   - Expose `ListPublicWaitingMatches(limit int) ([]models.PublicMatchItem, error)` on the App bridge.
   - Run `task client:check` to regenerate Wails TS bindings and verify build.

6. **Automated Testing:**
   - Add test case in `core-service/internal/service/service_test.go` verifying that:
     - Public waiting matches with `< 8` participants are returned.
     - Private matches (`is_public = false`) or active/finished matches (`status != 'WAITING'`) are excluded.
   - Run `task test` to verify all tests pass against PostgreSQL (`:5433`).

7. **Update Service Specification (`.routing/services/arena_service.md`):**
   - Document `PublicMatchItem` struct under Section 2.7.2.
   - Add `ListPublicWaitingMatches` to `ArenaCatalogService` under Section 2.7.3.
   - Add `GET /api/v1/arena/matches/public` to the endpoints table under Section 2.7.5.

8. **Change log:**
   - Document all changes and outputs in `.tickets/018-backend-arena-public-rooms/changelog.md`.

## Verification

All checks are listed in [spec.md](spec.md).
