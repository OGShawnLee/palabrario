# 018 — Spec checklist

Run by the router before accepting the ticket. All items must pass.

- [ ] Query `ListPublicWaitingMatches` exists in `sql/queries/arena.sql` and `sqlc generate` runs with 0 errors.
- [ ] DTO `PublicMatchItem` defined in `pkg/models/arena.go` with exact `snake_case` tags (`id_match`, `room_code`, `id_deck`, `deck_title`, `is_deck_official`, `host_user_name`, `total_rounds`, `round_time_sec`, `participant_count`).
- [ ] `ArenaCatalogService` exposes `ListPublicWaitingMatches(ctx, limit)`.
- [ ] Endpoint `GET /api/v1/arena/matches/public` mounted in `core-service/internal/api/router.go` and returns `200 OK` with JSON array.
- [ ] Rooms with `>= 8` participants or older than 2 hours are excluded from results.
- [ ] Private matches (`is_public = false`) or non-waiting matches (`status != 'WAITING'`) are excluded.
- [ ] Automated integration test in `core-service/internal/service/service_test.go` verifies the query and business rules.
- [ ] `task test` passes with 0 errors on RAM PostgreSQL (:5433).
- [ ] Wails bridge method `ListPublicWaitingMatches` exposed in `app/app_arena.go` and `task client:check` exits with code 0.
- [ ] `.routing/services/arena_service.md` updated with `PublicMatchItem`, method, and endpoint.
- [ ] Changes stay strictly within the declared scope.
- [ ] `.tickets/018-backend-arena-public-rooms/changelog.md` exists and is populated.
- [ ] Commit (if made): Conventional Commits, English, detailed body.
