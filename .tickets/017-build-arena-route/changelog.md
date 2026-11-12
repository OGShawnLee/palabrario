# 017 — Changelog

## Added
- `components/ui/AppDialog.vue`: reusable Reka UI dialog (`header` bar, title, Lucide `X` close with `button-icon`, `size` prop: `sm | md | lg | xl | 2xl` or a custom `max-w-*` class; default `md` = `max-w-lg`).
- `components/ui/AppDialogOverlay.vue`: `DialogOverlay` with the `dialog-overlay` shortcut.
- `components/arena/PublicRoomsTable.vue`: "Open Matches" list with Refresh button, empty state, and Join action; data only via `api.arena.listPublicWaitingMatches(20)`.
- `components/arena/CreateMatchDialog.vue`: room configuration modal (decks from `api.studio.listDecks()`, decks with < 8 terms locked, rounds 8/12/16 limited to deck size, 10/15/20s, public/private, creates via `api.arena.prepareMatch`).
- `api/contract/arena.contract.ts` (`IArenaAPI`), mounted on `IPalabrarioAPI`.
- `api/live/wails/arena.wails.ts` (`WailsArenaAdapter`) wired in `wails.adapter.ts`. Note: the Wails binding for room creation is `PrepareMatch` (as specified).
- `api/mock/fixtures/arena.fixture.ts` and `MockArenaAPI` (in-memory rooms, PIN generation, join).
- Routes: `/arena` breadcrumb meta; `/arena/history` placeholder (renders `ArenaView.vue` until Ticket 019).

## Changed
- `views/ArenaView.vue`: lobby with hero, PIN bar, Create Room, Match History link, public rooms table.

## Notes / Proposals
- Deck terms are approximated by `Deck.item_count` (the DTO exposes no approved-sense count); the backend remains authoritative (422).
- Real-time WebSocket gameplay is deferred; returned `MatchPackage` is kept in local state and shown as a "Ready" notice.
- No `uno.config.ts` changes needed.

## Revision: Backend Lookup Endpoint & Join Private Room Dialog

### Backend
- `sql/queries/arena.sql`: added `GetWaitingMatchByRoomCode`; regenerated sqlc (`core-service/internal/db/arena.sql.go`).
- `ArenaCatalogService.GetWaitingMatchByRoomCode`: normalizes the PIN (trim, strip `#`, uppercase), maps rows to `models.PublicMatchItem`, maps errors via `MapDBError(err, "Partida no encontrada")`.
- `ArenaHandler.GetMatchByRoomCode` and public route `GET /api/v1/arena/matches/code/{roomCode}`.
- Wails bridge `App.FindMatchRoom` (returns `nil, nil` on not found, so TypeScript receives `null`).

### Frontend
- `IArenaAPI.findMatchRoom` with Wails and mock implementations (mock searches live rooms, fixtures and private rooms).
- New `components/arena/JoinMatchDialog.vue` (PIN input, Find Room, loading/error/preview states, Join Room with full-room guard).
- `ArenaView.vue`: inline PIN form replaced by Create Room / Join with PIN / Match History actions; removed leftover `font-mono`.
