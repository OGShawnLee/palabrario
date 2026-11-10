# 017 — Spec checklist

Run by the router before accepting the ticket. All items must pass.

- [ ] Reusable Dialog components created in `app/frontend/src/components/ui/`:
  - `AppDialog.vue`: uses Reka UI Dialog with header class (`header`), Lucide `X` icon button (`button-icon`), and a configurable width size prop (`size` / `maxWidth`, defaulting to standard compact dialog with option for wide layout).
  - `AppDialogOverlay.vue`: wraps DialogOverlay with `dialog-overlay` shortcut.
- [ ] `/arena` route renders in `ArenaView.vue` under `AppLayout`.
- [ ] Room code PIN entry input (`field`, 6 chars uppercase) and `"Enter Room"` button present.
- [ ] `"Create Room"` button opens `CreateMatchDialog.vue` using `AppDialog`.
- [ ] `CreateMatchDialog.vue` adheres to `arena_service.md` rules:
  - Decks with $< 8$ terms are disabled/locked with explanation.
  - Round selectors offer valid counts (e.g. 8, 12, 16) and times (10s, 15s, 20s).
- [ ] `IArenaAPI` contract and adapters (Wails and Mock) properly implemented in `src/api/` and mounted on `IPalabrarioAPI`.
- [ ] `ArenaView.vue` and `PublicRoomsTable.vue` consume open rooms from `api.arena.listPublicWaitingMatches()` via `import { api } from "@/api"`. Zero hardcoded mock room arrays exist in views.
- [ ] `CreateMatchDialog.vue` loads candidate decks via `api.studio` and creates matches via `api.arena.prepareMatch()`.
- [ ] Public open rooms table renders with room PIN, deck, format, participants count (`N/8`), and join action.
- [ ] Refresh / Reload button is present above the public rooms table to manually refetch waiting rooms via `api.arena`.
- [ ] Match history link/button pointing to `/arena/history` is present.
- [ ] Strict monochrome design system: no raw colors outside theme tokens (`main`, `side`, `dull`, `light`, `ground`).
- [ ] Ground rules enforced: zero flag emojis, zero dialect codes (`CU`, `CL`, `MX`); full country/region names only.
- [ ] Mobile-first responsive layout (smooth stacking from mobile to desktop).
- [ ] Typecheck passes: `cd app && wails3 task check` exits with code 0.
- [ ] No application code outside declared scope, backend code, or `.context/` modified.
- [ ] `.tickets/017-build-arena-route/changelog.md` exists and is populated.
