# 017 — Build the Arena Lobby Route (/arena) and Reusable Modal Components

| Field | Value |
| :--- | :--- |
| **ID** | 017 |
| **Status** | draft |
| **Model** | Sonnet 5.5 (`claude-sonnet-5-5`, Antigravity provider) |
| **Role** | Frontend Implementer |
| **Spec** | [spec.md](spec.md) |

> **Draft Note:**  
> This ticket builds the `/arena` lobby route based on the desktop and mobile visual inspirations, introducing a reusable Reka UI Modal/Dialog component set (`AppDialog.vue`, `AppDialogOverlay.vue`) utilizing the new UnoCSS shortcuts (`header`, `button-icon`, `dialog-overlay`). The Arena lobby features a PIN entry bar, a primary "Create Room" button opening a spacious room configuration modal, and an active Public Rooms table.

## Scope

May create or edit:
- `app/frontend/src/components/ui/AppDialog.vue` (new reusable dialog wrapper with `header`, title, and Lucide `X` button)
- `app/frontend/src/components/ui/AppDialogOverlay.vue` (new reusable backdrop using `dialog-overlay`)
- `app/frontend/src/views/ArenaView.vue` (implement lobby view)
- `app/frontend/src/components/arena/CreateMatchDialog.vue` (room configuration modal)
- `app/frontend/src/components/arena/PublicRoomsTable.vue` (or subcomponents for rooms list)
- `.tickets/017-build-arena-route/changelog.md` (create)

Must not touch:
- `app/frontend/uno.config.ts` (shortcuts already exist; record any proposals in changelog)
- `app/*.go`, backend services, or `.context/`
- Any file outside the above scope

## Inputs

- [.routing/services/arena_service.md](../../.routing/services/arena_service.md): Canonical domain rules (`RN10`, `RN13`), DTO contracts (`CreateMatchRequest`, `MatchSummary`, `MatchStatus`), limits (min 8 approved senses, 8..20 rounds, 10..30s timers, max 8 players), and endpoints.
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
   "Contiendas Recientes" (Match History) will live on `/arena/history` (addressed in a dedicated follow-up ticket). In `/arena`, provide a link or button pointing to `/arena/history`.
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

2. **Section 1: Hero Header & Join Room Bar:**
   - **Hero Text:**
     - Title: `Arena de Contienda` (or `Contienda Arena` in English: `Multiplayer Arena`) in `text-3xl sm:text-4xl font-black text-main tracking-tight`.
     - Subtitle blurb: *"Compete in real-time with up to 8 players in definition and image rounds."* in `text-sm sm:text-base text-side mt-1`.
   - **Join PIN Bar & Actions:**
     - On desktop, aligned to the right or below the title; on mobile, stacked cleanly:
     - 6-character uppercase PIN input field (`field font-mono text-center tracking-widest uppercase w-32 sm:w-40`).
     - `"Enter Room →"` (`button-main`).
     - Primary `"Create Room"` (`button-side`).
     - Secondary link: `"Match History →"` navigating to `/arena/history`.

3. **Section 2: Public Open Matches Table:**
   - Table header includes title *"Open Matches"* and a **Refresh / Reload button** (with Lucide `RotateCw` or `RefreshCw` icon using `button-side` or `button-icon` styling) to manually refetch the waiting public rooms list without real-time socket polling.
   - Responsive table or scannable list showing current open rooms in `WAITING` status.
   - Columns:
     - Room Code (e.g. `#MX8921` in font-mono font-bold).
     - Deck Title (and region if regional deck).
     - Format (e.g. `8 rounds · 15s`).
     - Participants slot indicator (e.g. `3 / 8`).
     - Action: `"Join"` button.
   - Empty state when no public rooms exist: clean minimal box with prompt to create a room or enter a private PIN, plus the refresh button.

4. **Section 3: Room Configuration Modal (`CreateMatchDialog.vue`):**
   - Modal title: `"Create Match Room"`.
   - **Deck Selector (Left column on desktop, stacked on mobile):**
     - Lists study decks. Decks with $< 8$ approved terms are disabled/locked with a padlock icon (`Lock` from `lucide-vue-next`) and helper text `"(minimum 8 terms required)"` per `arena_service.md`.
     - Active selected deck highlighted with `bg-main text-white` and checkmark.
   - **Game Rules (Right column on desktop):**
     - Total rounds segmented control: `8 | 12 | 16` (enforcing `TotalRounds <= ApprovedSenses`).
     - Time per round segmented control: `10s | 15s | 20s`.
     - Public / Private match toggle.
     - Action: `"Create Room and Get PIN"` (`button-main w-full`).

5. **Validation:**
   - Run `cd app && wails3 task check` to verify Vue 3 and TypeScript typecheck.

6. **Change log:**
   - Document new components and proposals in `.tickets/017-build-arena-route/changelog.md`.

## Verification

All checks are listed in [spec.md](spec.md).
