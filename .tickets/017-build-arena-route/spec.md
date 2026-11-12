# 017 — Spec checklist

Run by the router before accepting the ticket. All items must pass.

- [x] Reusable Dialog components created in `app/frontend/src/components/ui/`:
  - `AppDialog.vue`: uses Reka UI Dialog with header class (`header`), Lucide `X` icon button (`button-icon`), and a configurable width size prop (`size` / `maxWidth`, defaulting to standard compact dialog with option for wide layout).
  - `AppDialogOverlay.vue`: wraps DialogOverlay with `dialog-overlay` shortcut.
- [x] `/arena` route renders in `ArenaView.vue` under `AppLayout` with `Home / Arena` breadcrumb.
- [x] `/arena/history` placeholder route registered with breadcrumb hierarchy.
- [x] Room code PIN entry input (`field`, 6 chars uppercase) and `"Enter Room"` button present.
- [x] `"Create Room"` button opens `CreateMatchDialog.vue` using `AppDialog`.
- [x] `CreateMatchDialog.vue` adheres to `arena_service.md` rules:
  - Decks with $< 8$ terms are disabled/locked with explanation.
  - Round selectors offer valid counts (e.g. 8, 12, 16) and times (10s, 15s, 20s).
- [x] `IArenaAPI` contract and adapters (Wails and Mock) properly implemented in `src/api/` and mounted on `IPalabrarioAPI`.
- [x] `app/frontend/src/api/mock/fixtures/arena.fixture.ts` contains authentic Spanish waiting room fixtures with zero flag emojis and zero dialect codes.
- [x] `ArenaView.vue` and `PublicRoomsTable.vue` consume open rooms from `api.arena.listPublicWaitingMatches()` via `import { api } from "@/api"`. Zero hardcoded mock room arrays exist in views.
- [x] `CreateMatchDialog.vue` loads candidate decks via `api.studio` and creates matches via `api.arena.prepareMatch()`.
- [x] Public open rooms table renders with room PIN, deck, format, participants count (`N/8`), and join action.
- [x] Refresh / Reload button is present above the public rooms table to manually refetch waiting rooms via `api.arena`.
- [x] Match history link/button pointing to `/arena/history` is present.
- [x] Strict monochrome design system: no raw colors outside theme tokens (`main`, `side`, `dull`, `light`, `ground`).
- [x] Ground rules enforced: zero flag emojis, zero dialect codes (`CU`, `CL`, `MX`); full country/region names only.
- [x] Mobile-first responsive layout (smooth stacking from mobile to desktop).
- [x] Typecheck passes: `cd app && wails3 task check` exits with code 0.
- [x] No application code outside declared scope, backend code, or `.context/` modified.
- [x] `.tickets/017-build-arena-route/changelog.md` exists and is populated.
- [x] Commit (if made): Conventional Commits, English, detailed body.

## Revision Checklist
- [x] Font hygiene: zero `font-mono` occurrences in `ArenaView.vue`, `PublicRoomsTable.vue`, `CreateMatchDialog.vue` (Satoshi font exclusively; completed by manual refactor).
- [x] Subtitles in `CreateMatchDialog.vue` replace `@TODO` with approved editorial copy (completed by manual refactor).
- [x] Backend: Typed sqlc query `GetWaitingMatchByRoomCode` added to `sql/queries/arena.sql` and `core-service/internal/repository/` regenerated.
- [x] Backend: `GetWaitingMatchByRoomCode` implemented in `ArenaCatalogService` (`core-service/internal/service/arena_catalog_service.go`).
- [x] Backend: HTTP handler `GetMatchByRoomCode` implemented in `handlers/arena_handler.go` and mounted on `GET /api/v1/arena/matches/code/{roomCode}` in `router.go`.
- [x] Backend / Bridge: Wails bridge `FindMatchRoom` implemented in `app/app_arena.go` and Wails TypeScript bindings generated.
- [x] Frontend: Inline PIN bar in `ArenaView.vue` replaced with clean `"Create Room"` and `"Join with PIN"` action buttons.
- [x] Frontend: `JoinMatchDialog.vue` implemented with `AppDialog`, PIN input, `"Find Room"` search action, room preview card (deck, host, format, capacity), and `"Join Room →"` confirmation.
- [x] Frontend: `IArenaAPI.findMatchRoom(roomCode)` defined in contract and implemented in both Mock and Wails adapters.
- [x] Full check passes: `cd app && wails3 task check` and backend compilation succeed with code 0.
