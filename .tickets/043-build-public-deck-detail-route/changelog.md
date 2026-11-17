# 043 — Changelog: Public Deck Detail Route (`/decks/:id`)

## Summary of Changes

### Routes
- **`app/frontend/src/router/index.ts`**: Registered public route `/decks/:id` (`name: "deck-detail"`) under `AppLayout` without `requiresAuth`, with `meta: { breadcrumb: [{ label: "Home", to: "/" }, { label: "Decks", to: "/decks" }, { label: "Deck Detail" }] }`, lazy-loading `DeckDetailView.vue`.

### Views Created
- **`app/frontend/src/views/DeckDetailView.vue`**: Implemented the public deck inspection view (`/decks/:id`):
  - Validates positive integer `:id` and fetches `api.studio.getDeck(id)` using a monotonic request counter.
  - Updates `breadcrumbLeafLabel` to `deck.title` on load and clears it on invalid ID, error, or unmount.
  - Handles invalid `:id`, `404` (`null`), `403` (`RN04` private deck belonging to another user), and network failures with monochrome alert panels, `Retry` (on network failure), and `← Back to Decks` (`/decks`).
  - Displays deck header metadata (`Public`/`Private`, `Official`/`Community`, `item_count` terms, updated date, full `region_names` via `badge-side`).
  - Displays Arena readiness banner and `"Play in Arena →"` link to `/arena` when `deck.item_count >= 8`.
  - Resolves card lemmas via `api.dictionary.getWordById` (linking to `/dictionary/word/${sense.id_word}`), alongside `region_name` (`badge-side`), `cefr_level` (`badge-main`), `short_definition`, `definition`, and `examples`.
  - Contains zero `/studio/study` dead links, zero SM-2 mastery fields, zero `flag_emoji` or dialect codes, and zero `@bindings/palabrario/app/app` imports.

## Auth & Save Gates (`RN03`, `RN04`)
- **Guest-safe `listStudioDecks` gate**: Calls `api.studio.listStudioDecks()` only when `isAuthenticated.value` is `true` so anonymous visitors never trigger authenticated Studio requests.
- **Save / Unsave flow**: When `canManage` is false, renders `Save to Studio` / `✓ Saved in Studio` (`api.studio.saveDeck` / `api.studio.unsaveDeck`); clicking `Save to Studio` while unauthenticated redirects to `/auth/sign-in?redirect=${encodeURIComponent(route.fullPath)}`.
- **Owner / Admin link (`RN03`)**: Computes `canManage = currentUser?.role === AccountRole.RoleAdmin || currentUser?.id_user === deck.id_user` and renders `"Open in Studio →"` linking to `/studio/decks/${deck.id_deck}` instead of bookmarking controls.

## Verification

### Static Pattern Check
Verified zero matches for `flag_emoji`, `region_code`, `country_code`, `/studio/study`, `@bindings/palabrario/app/app`, `due_cards_today`, `retention_percent`, and `easiness_factor` across `app/frontend/src/router/index.ts` and `app/frontend/src/views/DeckDetailView.vue`.

### `wails3 task check` (executed in `app/`)
Exited with code `0`:
- `wails3 generate bindings -clean=true -ts`: `Processed: 264 Packages, 1 Service, 98 Methods, 6 Enums, 75 Models, 0 Events in 27.069718s.`
- `bun run --cwd frontend type-check` (`vue-tsc --noEmit`): passed with 0 errors.
- `bun run --cwd frontend build` (`vite build --mode production`): `✓ 2519 modules transformed`, `✓ built in 7.45s` (including `DeckDetailView-9VG37B6R.js`).
- `go build .`: passed with 0 errors.
