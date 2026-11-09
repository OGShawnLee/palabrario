# 018 — Backend: Public Active Arena Matches Listing Endpoint

| Field | Value |
| :--- | :--- |
| **ID** | 018 |
| **Status** | draft |
| **Model** | Sonnet 5.5 (`claude-sonnet-5-5`, Antigravity provider) |
| **Role** | Backend Implementer |
| **Spec** | [spec.md](spec.md) |

> **Draft Note:**  
> This backend ticket adds an endpoint to list all open, public arena rooms currently waiting for players (`status = 'WAITING' AND is_public = true`), joining deck metadata, host username, and current connected player count so the frontend lobby can display live games ready to join.

## Scope

May create or edit:
- `sql/queries/arena.sql` (add query `ListPublicWaitingMatches`)
- `core-service/internal/db/` (regenerate `sqlc`)
- `pkg/models/arena.go` (add `PublicMatchItem` or update `ArenaRoomSummary`)
- `core-service/internal/service/arena_catalog_service.go` (add `ListPublicWaitingMatches` method)
- `core-service/internal/api/handlers/arena_handler.go` (expose `GET /api/v1/arena/matches/public`)
- `core-service/internal/service/service_test.go` (add unit/service tests on RAM PostgreSQL :5433)
- `app/app_arena.go` (expose Wails Desktop bridge method and regenerate bindings)
- `.routing/services/arena_service.md` (document new `PublicMatchItem` DTO, SOA interface method, and endpoint)
- `.tickets/018-backend-arena-public-rooms/changelog.md` (create)

Must not touch:
- `app/frontend/src/` (UI components)
- `.context/`
- Any file outside the above scope

## Inputs

- [.routing/services/arena_service.md](../../.routing/services/arena_service.md): Arena domain context and existing contracts.
- [`SCHEMA_FINAL.sql`](../../SCHEMA_FINAL.sql): Table definitions for `matches`, `decks`, and `accounts`/`users`.
- [`pkg/models/arena.go`](../../pkg/models/arena.go): Existing DTOs (`MatchPackage`, `MatchSummary`, `MatchStatus`).

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
         d.id_region,
         r.name AS region_name,
         u.username AS host_username,
         m.total_rounds,
         m.round_time_sec,
         m.created_at,
         (SELECT COUNT(*) FROM match_participants mp WHERE mp.id_match = m.id_match) AS participant_count
     FROM matches m
     JOIN decks d ON d.id_deck = m.id_deck
     LEFT JOIN regions r ON r.id_region = d.id_region
     JOIN users u ON u.id_user = m.id_user_host
     WHERE m.status = 'WAITING' AND m.is_public = true
     ORDER BY m.created_at DESC
     LIMIT $1;
     ```
   - Run `sqlc generate`.

2. **DTO & Service Implementation:**
   - Define `PublicMatchItem` in `pkg/models/arena.go`.
   - Expose `ListPublicWaitingMatches(ctx context.Context, limit int) ([]models.PublicMatchItem, error)` on `ArenaCatalogService`.

3. **HTTP Handler & Wails Bridge:**
   - Expose `GET /api/v1/arena/matches/public` in `arena_handler.go` (public or `OptionalAuth`).
   - Implement `ListPublicWaitingMatches` in `app/app_arena.go` and run `wails3 task client:bindings`.

4. **Testing & Validation:**
   - Add test case in `service_test.go` verifying public waiting matches are returned and non-public or finished matches are excluded.
   - Run `task test` against the RAM PostgreSQL harness (`:5433`).

5. **Update Service Specification (`.routing/services/arena_service.md`):**
   - Document `PublicMatchItem` struct under Section 2.7.2.
   - Add `ListPublicWaitingMatches` to `ArenaCatalogService` under Section 2.7.3.
   - Add `GET /api/v1/arena/matches/public` to the endpoints table under Section 2.7.5.

## Verification

All checks are listed in [spec.md](spec.md).
