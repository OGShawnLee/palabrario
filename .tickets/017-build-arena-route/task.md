# 017 — Build the Arena Lobby Route (/arena) and Reusable Modal Components

| Field | Value |
| :--- | :--- |
| **ID** | 017 |
| **Status** | completed |
| **Model** | Claude Sonnet 5.5, `effort: medium` (`claude-sonnet-5-5`) |
| **Role** | Frontend Implementer |
| **Spec** | [spec.md](spec.md) |

> **Specification Note:**  
> This ticket builds the `/arena` lobby route based on the desktop and mobile visual inspirations, introducing a reusable Reka UI Modal/Dialog component set (`AppDialog.vue`, `AppDialogOverlay.vue`) utilizing the new UnoCSS shortcuts (`header`, `button-icon`, `dialog-overlay`). The Arena lobby features a PIN entry bar, a primary "Create Room" button opening a spacious room configuration modal, and an active Public Rooms table. Real-time WebSocket gameplay screens are deferred to subsequent work.

## Scope

May create or edit:
- `sql/queries/arena.sql` (add GetWaitingMatchByRoomCode query)
- `core-service/internal/repository/` (regenerate via sqlc generate)
- `core-service/internal/service/arena_catalog_service.go` (implement GetWaitingMatchByRoomCode)
- `core-service/internal/api/handlers/arena_handler.go` (handler for GET /api/v1/arena/matches/code/{roomCode})
- `core-service/internal/api/router.go` (mount route)
- `app/app_arena.go` (Wails FindMatchRoom bridge method)
- `app/frontend/src/components/ui/AppDialog.vue` (new reusable dialog wrapper with `header`, title, and Lucide `X` button)
- `app/frontend/src/components/ui/AppDialogOverlay.vue` (new reusable backdrop using `dialog-overlay`)
- `app/frontend/src/views/ArenaView.vue` (implement lobby view)
- `app/frontend/src/components/arena/CreateMatchDialog.vue` (room configuration modal)
- `app/frontend/src/components/arena/JoinMatchDialog.vue` (new: dedicated private room PIN lookup & preview modal)
- `app/frontend/src/components/arena/PublicRoomsTable.vue` (or subcomponents for rooms list)
- `app/frontend/src/api/contract/arena.contract.ts` (define IArenaAPI)
- `app/frontend/src/api/contract/index.ts` (mount arena in IPalabrarioAPI)
- `app/frontend/src/api/live/wails/arena.wails.ts` (implement WailsArenaAdapter)
- `app/frontend/src/api/live/wails/wails.adapter.ts` (mount arena adapter)
- `app/frontend/src/api/mock/mock.adapter.ts` (implement MockArenaAPI with in-memory waiting rooms)
- `app/frontend/src/api/mock/fixtures/arena.fixture.ts` (mock public waiting matches fixture)
- `app/frontend/src/router/index.ts` (register `/arena` breadcrumbs and `/arena/history` placeholder)
- `.tickets/017-build-arena-route/changelog.md` (create)

Must not touch:
- `app/frontend/uno.config.ts` (shortcuts already exist; record any proposals in changelog)
- `SCHEMA_FINAL.sql` (no schema migration needed; query operates on existing tables)
- `.context/`
- Any file outside the above scope

## Inputs

- **Dependency Note:** Assumes Ticket 011 is complete and `api.studio.listDecks()` is available on `IStudioAPI`.
- [.routing/API.md](../../.routing/API.md): Architectural rules (Rule 1: Never import `@bindings/.../app` in UI; always import `api` from `@/api`).
- [app/frontend/src/api/index.ts](../../app/frontend/src/api/index.ts): Universal `api` singleton (`api.arena`, `api.studio`).
- [.routing/services/arena_service.md](../../.routing/services/arena_service.md): Canonical domain rules (`RN10`, `RN13`), DTO contracts (`CreateMatchRequest`, `MatchSummary`, `MatchStatus`, `PublicMatchItem`, `MatchPackage`), limits (min 8 approved senses, 8..20 rounds, 10..30s timers, max 8 players), and endpoints.
- [.routing/services/studio_service.md](../../.routing/services/studio_service.md): Deck contracts (`Deck`, `DeckItem`, `StudioDeckSummary`) used when selecting decks for room creation.
- [.routing/CONTEXT.md](../../.routing/CONTEXT.md): Architecture ground rules (strict prohibition of `flag_emoji` and regional codes; full region names required).
- [.routing/DESIGN.md](../../.routing/DESIGN.md): Monochrome design tokens (`main`, `side`, `dull`, `light`, `ground`), typography (`Satoshi`), surfaces, and buttons (`button-main`, `button-side`, `button-icon`).
- [app/frontend/uno.config.ts](../../app/frontend/uno.config.ts): Shortcuts (`header`, `button-icon`, `button-main`, `button-side`, `field`, `binding`, `dialog-overlay`).
- Visual references:
  - Desktop layout: [Image 1: Desktop Arena inspiration](../../../.t3/userdata/attachments/e7510c5f-4aee-4fe2-8b04-76e954f1ea52-ed6ca1a6-8ef6-4587-8b46-49610bce8526.png)
  - Mobile layout: [Image 2: Mobile Arena inspiration](../../../.t3/userdata/attachments/e7510c5f-4aee-4fe2-8b04-76e954f1ea52-dc3bfd73-faa6-419c-ae73-d3b357f70a43.png)

## Deviations from Visual References

1. **Room Configuration as Modal:**
   Rather than embedding the "Configurar Nueva Sala" configuration box permanently in the page center, room creation is triggered via a primary `"Create Room"` button which opens a large dialog (`CreateMatchDialog.vue`).
2. **Public Room Browser Table:**
   The main body of `/arena` features a table displaying open/waiting public rooms (`WAITING` status) allowing fast one-click join.
3. **Match History Moved to Dedicated Sub-Route:**
   "Contiendas Recientes" (Match History) will live on `/arena/history` (addressed in Ticket 019). In `/arena`, provide a link or button pointing to `/arena/history`.
4. **Monochrome Design System & Zero Regional Codes:**
   Strictly NO `flag_emoji` and NO regional codes (`CU`, `CL`, `MX`). Render full textual region names. Use monochrome borders and surfaces (`bg-main`, `bg-ground`, `border-light`).

## Task

1. **Reusable Dialog Components (`app/frontend/src/components/ui/`):**
   - Create `AppDialog.vue`:
     - Built on Reka UI (`DialogRoot`, `DialogPortal`, `DialogContent`, `DialogTitle`, `DialogDescription`, `DialogClose`).
     - Props:
       - `open`: boolean model (`defineModel<boolean>("open")`).
       - `title`: string.
       - `size` or `maxWidth`: customizable width size prop (e.g. `'sm'` | `'md'` | `'lg'` | `'xl'` | `'2xl'` or string class override, defaulting to standard compact dialog `max-w-lg`, with the room creation dialog configuring a wide layout `max-w-4xl` or `max-w-5xl`).
     - Header bar with UnoCSS shortcut `class="header"`:
       - Left: `DialogTitle` in `font-bold text-main text-base sm:text-lg`.
       - Right: `DialogClose` containing a Lucide `X` icon using UnoCSS shortcut `class="button-icon flex items-center justify-center cursor-pointer"`.
     - Inner content container with `p-4 sm:p-6 overflow-y-auto max-h-[85vh]`.
     - Sharp edges and clean styling: `rounded-none border border-light bg-white shadow-2xl`.
   - Create `AppDialogOverlay.vue`:
     - Wrapping Reka UI `DialogOverlay` with `class="dialog-overlay"`.

2. **API Client Support for Arena (`IArenaAPI`) & Fixtures:**
   - Create `app/frontend/src/api/contract/arena.contract.ts`:
     - Define `IArenaAPI`:
       ```ts
       export interface IArenaAPI {
         listPublicWaitingMatches(limit?: number): Promise<PublicMatchItem[]>
         prepareMatch(req: CreateMatchRequest): Promise<MatchPackage>
         joinMatch(roomCode: string): Promise<MatchPackage>
       }
       ```
     - Import DTOs from `@bindings/palabrario/pkg/models/index.js`.
   - Update `app/frontend/src/api/contract/index.ts` to include `arena: IArenaAPI` on `IPalabrarioAPI`.
   - Create `app/frontend/src/api/live/wails/arena.wails.ts`:
     - Implement `WailsArenaAdapter` delegating to `Wails.ListPublicWaitingMatches(limit)`, `Wails.PrepareMatch(req)`, and `Wails.JoinArenaMatch(roomCode)`.
   - Update `app/frontend/src/api/live/wails/wails.adapter.ts` to instantiate `WailsArenaAdapter`.
   - Create `app/frontend/src/api/mock/fixtures/arena.fixture.ts`:
     - Export initial `mockPublicWaitingMatches: PublicMatchItem[]` with authentic Spanish entries (e.g. `#MX8921` "Dichos y Refranes Mexicanos", `#CO4410` "Vocabulario Paisa", `#AR1092` "Lunfardo Rioplatense"). Ensure full region names and zero flag emojis.
   - Update `app/frontend/src/api/mock/mock.adapter.ts`:
     - Implement `MockArenaAPI` maintaining simulated open rooms in memory backed by `arena.fixture.ts`, and handling `prepareMatch` (generates new room PIN) and `joinMatch`.

3. **Router Registration & Breadcrumbs (`app/frontend/src/router/index.ts`):**
   - Update `/arena` route with `meta: { breadcrumb: [{ label: "Home", to: "/" }, { label: "Arena" }] }`.
   - Register placeholder route `/arena/history` (`name: "arena-history"`) with `meta: { breadcrumb: [{ label: "Home", to: "/" }, { label: "Arena", to: "/arena" }, { label: "Match History" }] }` rendering a minimal view or pointing to `ArenaView.vue` until Ticket 019 implements the full replay view.

4. **Section 1: Hero Header & Join Room Bar:**
   - **Hero Text:**
     - Title: `Multiplayer Arena` in `text-3xl sm:text-4xl font-black text-main tracking-tight`.
     - Subtitle blurb: *"Compete in real-time with up to 8 players in definition and image rounds."* in `text-sm sm:text-base text-side mt-1`.
   - **Join PIN Bar & Actions:**
     - 6-character uppercase PIN input field (`field font-mono text-center tracking-widest uppercase w-32 sm:w-40`).
     - `"Enter Room →"` (`button-main`), calls `api.arena.joinMatch(pin)`.
     - Primary `"Create Room"` (`button-side`), opens `CreateMatchDialog.vue`.
     - Secondary link: `"Match History →"` navigating to `/arena/history`.
   - **Note on Real-time Gameplay:** Real-time WebSocket connection (`ConnectArenaRoom`) and in-game match screens are deferred to a subsequent ticket. Upon calling `joinMatch` or `prepareMatch`, handle the returned `MatchPackage` gracefully without throwing (e.g. closing dialog or keeping local state).

5. **Section 2: Public Open Matches Table (`PublicRoomsTable.vue`):**
   - **Data Fetching via API Client:**
     - Import `api` from `@/api`:
       ```ts
       import { api } from "@/api"
       ```
     - Fetch matches on mount and on Refresh click:
       ```ts
       const rooms = ref<PublicMatchItem[]>([])
       const isLoading = ref(true)

       async function fetchRooms() {
         isLoading.value = true
         try {
           rooms.value = await api.arena.listPublicWaitingMatches(20)
         } finally {
           isLoading.value = false
         }
       }
       ```
     - **Strict Rule:** Never hardcode mock open room arrays in the Vue component. All room data must flow through `api.arena`.
   - Table header includes title *"Open Matches"* and a **Refresh / Reload button** (with Lucide `RotateCw` or `RefreshCw` icon using `button-side` or `button-icon` styling) to manually call `fetchRooms()`.
   - Responsive table or scannable list showing current open rooms in `WAITING` status.
   - Columns:
     - Room Code (e.g. `#MX8921` in font-mono font-bold).
     - Deck Title (and region if regional deck).
     - Format (e.g. `8 rounds · 15s`).
     - Participants slot indicator (e.g. `3 / 8`).
     - Action: `"Join"` button (calls `api.arena.joinMatch(room.room_code)`).
   - Empty state when no public rooms exist: clean minimal box with prompt to create a room or enter a private PIN, plus the refresh button.

6. **Section 3: Room Configuration Modal (`CreateMatchDialog.vue`):**
   - Modal title: `"Create Match Room"`.
   - **Deck Selector (Left column on desktop, stacked on mobile):**
     - Fetch candidate decks using `api.studio.listDecks(...)`. Zero hardcoded decks in the modal.
     - Decks with $< 8$ approved terms are disabled/locked with a padlock icon (`Lock` from `lucide-vue-next`) and helper text `"(minimum 8 terms required)"` per `arena_service.md`.
     - Active selected deck highlighted with `bg-main text-white` and checkmark.
   - **Game Rules (Right column on desktop):**
     - Total rounds segmented control: `8 | 12 | 16` (enforcing `TotalRounds <= ApprovedSenses`).
     - Time per round segmented control: `10s | 15s | 20s`.
     - Public / Private match toggle.
     - Action: `"Create Room and Get PIN"` (`button-main w-full`), calls `api.arena.prepareMatch(payload)`.

7. **Validation:**
   - Run `cd app && wails3 task check` to verify Vue 3 and TypeScript typecheck.

8. **Change log:**
   - Document new components and proposals in `.tickets/017-build-arena-route/changelog.md`.

## Revision Guide: Backend Lookup Endpoint & Join Private Room Dialog

> **Note:** Font hygiene (removal of `font-mono`) and form input subtitles in `CreateMatchDialog.vue` have already been completed by the reviewer's manual refactor.

### Part A: Backend Implementation (`sqlc`, `core-service`, and Wails bridge)

#### Task 1: Typed SQL Query (`sql/queries/arena.sql`)
Add `GetWaitingMatchByRoomCode` to `sql/queries/arena.sql`:
```sql
-- name: GetWaitingMatchByRoomCode :one
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
WHERE UPPER(m.room_code) = UPPER($1)
  AND m.status = 'WAITING'
  AND m.created_at >= now() - INTERVAL '2 hours'
LIMIT 1;
```
Run `sqlc generate` in the root repository to regenerate `core-service/internal/repository/`.

#### Task 2: Service Layer (`core-service/internal/service/arena_catalog_service.go`)
1. In `ArenaCatalogService` interface, ensure `GetWaitingMatchByRoomCode` is declared:
   ```go
   GetWaitingMatchByRoomCode(ctx context.Context, roomCode string) (*models.PublicMatchItem, error)
   ```
2. Implement `GetWaitingMatchByRoomCode` in `arenaCatalogService`:
   - Normalize `roomCode := strings.TrimSpace(strings.ToUpper(strings.TrimPrefix(rawRoomCode, "#")))`.
   - Call `s.queries.GetWaitingMatchByRoomCode(ctx, roomCode)`.
   - If `err != nil`, return `MapDBError(err, "Partida no encontrada")`.
   - Map row fields to `models.PublicMatchItem` and return `&item, nil`.

#### Task 3: HTTP Handler & Router (`core-service/internal/api/`)
1. In `core-service/internal/api/handlers/arena_handler.go`, add `GetMatchByRoomCode`:
   ```go
   // GetMatchByRoomCode maneja GET /api/v1/arena/matches/code/{roomCode}
   func (h *ArenaHandler) GetMatchByRoomCode(w http.ResponseWriter, r *http.Request) {
       roomCode := chi.URLParam(r, "roomCode")
       match, err := h.arenaService.GetWaitingMatchByRoomCode(r.Context(), roomCode)
       if err != nil {
           status, msg := service.HTTPStatusAndMessage(err, "Partida no encontrada")
           response.Error(w, status, msg)
           return
       }
       response.JSON(w, http.StatusOK, match)
   }
   ```
2. In `core-service/internal/api/router.go`, mount the public route next to `/matches/public`:
   ```go
   r.Get("/matches/public", arenaHandler.ListPublicMatches)
   r.Get("/matches/code/{roomCode}", arenaHandler.GetMatchByRoomCode)
   ```

#### Task 4: Wails Bridge (`app/app_arena.go`)
In `app/app_arena.go`, implement `FindMatchRoom`:
```go
// FindMatchRoom consulta una partida en espera por su código PIN.
// Retorna (nil, nil) si no existe para que TypeScript lo reciba como `null`.
func (a *App) FindMatchRoom(roomCode string) (*models.PublicMatchItem, error) {
    roomCode = strings.TrimSpace(strings.ToUpper(strings.TrimPrefix(roomCode, "#")))
    if roomCode == "" {
        return nil, nil
    }

    targetURL := a.coreURL + "/api/v1/arena/matches/code/" + url.PathEscape(roomCode)
    var item models.PublicMatchItem
    err := a.doRequest(a.context(), http.MethodGet, targetURL, nil, &item)
    if err != nil {
        if strings.Contains(strings.ToLower(err.Error()), "not found") || strings.Contains(strings.ToLower(err.Error()), "no encontrada") {
            return nil, nil
        }
        return nil, err
    }
    return &item, nil
}
```
Run `cd app && wails3 generate bindings -clean=true -ts` (or `wails3 task check`) to generate updated TypeScript bindings for `FindMatchRoom`.

---

### Part B: Frontend Implementation

#### Task 5: Replace Inline PIN Bar in `ArenaView.vue` with Clean Action Buttons
In `app/frontend/src/views/ArenaView.vue`:
- Remove the inline `<form>` containing the naked PIN input and `"Enter Room →"` button.
- In its place, render a clean, balanced action cluster:
  - Primary button: `"Create Room"` (`button-main cursor-pointer`) $\to$ opens `CreateMatchDialog.vue`.
  - Secondary button: `"Join with PIN"` (`button-side cursor-pointer`) $\to$ opens `JoinMatchDialog.vue`.
  - Tertiary link: `"Match History →"` (`text-sm font-medium text-side underline sm:ml-auto`) $\to$ navigates to `/arena/history`.

#### Task 6: Create Dedicated Private Room Join Modal (`JoinMatchDialog.vue`)
Create `app/frontend/src/components/arena/JoinMatchDialog.vue`:
- Built on `AppDialog.vue` with `size="md"`.
- Modal title: `"Join Private Room"`.
- Modal description: `"Enter the 6-character room PIN provided by the host."`
- **Upper Section — PIN Input & Search:**
  - Standard Satoshi text input: `field tracking-widest uppercase text-center font-bold text-lg w-full h-12` (maxlength 6, placeholder `"PIN"`, auto-uppercase, Satoshi font, no `font-mono`).
  - Search action button: `"Find Room"` (`button-main h-12 px-5 cursor-pointer`). Disabled until exactly 6 characters are entered or while search is in flight.
- **Lower Section — Room Preview Card & Confirmation:**
  - While searching: Minimal text `"Searching for room..."` (`text-sm text-dull`).
  - When not found: Error message (`border border-main p-3 text-sm text-main`) stating `"No active room found with PIN #{code}."`.
  - When room is found:
    - Display a clean monochrome preview card (`border border-light bg-ground p-4 flex flex-col gap-3`):
      - Header: Room code badge (`#MX8921` in `font-bold tracking-widest text-main`) and deck title (truncated).
      - Host blurb: `"Hosted by @" + room.host_user_name` (`text-xs text-dull`).
      - Format & Players: `"{total_rounds} rounds · {round_time_sec}s"` and `"{participant_count} / 8 players"` (`tabular-nums`).
    - Primary action at bottom: `"Join Room →"` (`button-main w-full h-12 cursor-pointer`).
      - Calls `api.arena.joinMatch(room.room_code)`, emits `'joined', matchPackage`, closes dialog, and clears local preview.
      - If room is full (`participant_count >= 8`), display `"(Room is full)"` and disable the button.

#### Task 7: API Layer Contract & Adapters for `findMatchRoom`
1. In `app/frontend/src/api/contract/arena.contract.ts`, expand `IArenaAPI`:
   ```ts
   export interface IArenaAPI {
     listPublicWaitingMatches(limit?: number): Promise<PublicMatchItem[]>
     findMatchRoom(roomCode: string): Promise<PublicMatchItem | null>
     prepareMatch(req: CreateMatchRequest): Promise<MatchPackage>
     joinMatch(roomCode: string): Promise<MatchPackage>
   }
   ```
2. In `app/frontend/src/api/live/wails/arena.wails.ts`:
   - Implement `findMatchRoom(roomCode: string): Promise<PublicMatchItem | null>`:
     - Normalizes room code.
     - Calls `(Wails as any).FindMatchRoom(roomCode)`.
3. In `app/frontend/src/api/mock/mock.adapter.ts`:
   - Implement `findMatchRoom(roomCode: string): Promise<PublicMatchItem | null>`:
     - Normalizes `roomCode` (trims whitespace, strips `#`, uppercase).
     - Searches `this.rooms` and `mockPublicWaitingMatches`.
     - Returns `clone(room)` if found, or `null` if not found.

#### Task 8: Verification & Changelog
- Verify that `cd app && wails3 task check` exits with code 0.
- Update `.tickets/017-build-arena-route/changelog.md` with backend queries, handlers, `JoinMatchDialog.vue`, and `findMatchRoom`.

## Verification

All checks are listed in [spec.md](spec.md).
