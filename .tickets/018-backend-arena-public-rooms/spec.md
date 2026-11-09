# 018 — Spec checklist

Run by the router before accepting the ticket. All items must pass.

- [ ] `sql/queries/arena.sql` includes `ListPublicWaitingMatches` query with proper filters (`status = 'WAITING' AND is_public = true`).
- [ ] `sqlc generate` succeeds and typed DB models update cleanly.
- [ ] `pkg/models/arena.go` declares `PublicMatchItem`.
- [ ] `ArenaCatalogService.ListPublicWaitingMatches` implemented in `core-service/internal/service/`.
- [ ] Handler `GET /api/v1/arena/matches/public` exposed in `arena_handler.go`.
- [ ] Wails bridge method implemented in `app/app_arena.go` and TypeScript bindings regenerated under `app/frontend/bindings/`.
- [ ] Automated tests in `service_test.go` pass 100% on RAM PostgreSQL harness (`:5433`).
- [ ] `.routing/services/arena_service.md` updated with `PublicMatchItem` DTO, interface method, and endpoint documentation.
- [ ] Typecheck and check passes: `cd app && wails3 task check` exits with code 0.
- [ ] `.tickets/018-backend-arena-public-rooms/changelog.md` exists and is populated.
