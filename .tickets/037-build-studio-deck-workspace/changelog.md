# 037 — Changelog: Studio Deck Workspace (`/studio/decks/:id`) & Create/Edit Dialog

## Summary of Changes

### Routes
- **`app/frontend/src/router/index.ts`**: Registered `/studio/decks/:id` (`name: "studio-deck-detail"`) under `AppLayout` with `meta: { requiresAuth: true, breadcrumb: [{ label: "Home", to: "/" }, { label: "Studio", to: "/studio" }, { label: "Deck" }] }`, lazy-loading `StudioDeckDetailView.vue`.

### Components Created
- **`app/frontend/src/components/studio/DeleteConfirmDialog.vue`**: Reusable monochrome confirmation dialog built on `AppDialog.vue` (`size="sm"`), supporting configurable `title`, `message`, `confirmLabel`, `confirmingLabel`, `cancelLabel`, `isDeleting`, and inline `role="alert"` `error` rendering (ready for reuse by Decks and Ticket 038 Collections).
- **`app/frontend/src/components/studio/DeckFormDialog.vue`**: Create/Edit deck metadata dialog built on `AppDialog.vue` (`size="md"`), validating trimmed `title` length (`3..150` characters), optional `description`, and `is_public` toggle, submitting `CreateDeckRequest` / `UpdateDeckRequest` through `api.studio.createDeck` and `api.studio.updateDeck` with inline `role="alert"` errors.
- **`app/frontend/src/components/studio/AddSenseToDeckDialog.vue`**: Dictionary-to-deck sense picker dialog built on `AppDialog.vue` (`size="xl"`). Performs debounced lemma search via `api.dictionary.searchWords`, loads full regional senses via `api.dictionary.getWordById`, filters out senses already present in the deck (`existingSenseIds`) to prevent duplicates, displays `region_name`, `cefr_level`, `short_definition`, and `definition`, and calls `api.studio.addDeckItem` with inline `role="alert"` error surfacing.

### Views Created & Updated
- **`app/frontend/src/views/StudioView.vue`**: Preserved existing copy, layout, and styling; wired the `New Deck` button in the `Study Decks` section header to open `DeckFormDialog` and navigate to `/studio/decks/${deck.id_deck}` on creation.
- **`app/frontend/src/views/studio/StudioDeckDetailView.vue`**: Implemented the `/studio/decks/:id` workspace:
  - Validates positive integer `:id` and fetches `api.studio.getDeck(id)` (plus `api.studio.listStudioDecks()` to hydrate community `isSaved` status) with a monotonic request counter to prevent stale updates.
  - Updates `breadcrumbLeafLabel` to `deck.title` and clears it on unmount.
  - Renders resolved lemmas (via `api.dictionary.getWordById`, linked to `/dictionary/word/${sense.id_word}`), `region_name`, `cefr_level`, `short_definition`, `definition`, and `examples`.
  - Displays the Arena readiness banner (`"Play in Arena →"` link to `/arena` when `item_count >= 8`, or remaining card count prompt when `< 8`).
  - Contains zero `/studio/study` dead links, zero SM-2 mastery fields, zero `flag_emoji` or dialect codes, and zero `@bindings/palabrario/app/app` imports.

## Business Rule Enforcement (`RN03`, `RN04`, `RN05`)
- **`RN03` (Ownership & Capability Split)**: `StudioDeckDetailView.vue` computes `canManage = currentUser?.role === AccountRole.RoleAdmin || currentUser?.id_user === deck.id_user`. Mutating controls (`Add Card`, `Edit Metadata`, `Delete Deck`, and per-card `Remove`) render only when `canManage` is true; non-owners see `Save to Studio` / `Unsave from Studio` (`api.studio.saveDeck` / `api.studio.unsaveDeck`) instead.
- **`RN04` (Private Resource Handling)**: `StudioDeckDetailView.vue` distinguishes invalid `:id` / `404` (`null`) from `403` private-deck errors and network failures, rendering dedicated monochrome panels with `← Back to Studio` navigation (and `Retry` on network errors).
- **`RN05` (Official Curation Gate)**: `DeckFormDialog.vue` renders the `Official Deck` (`is_official`) checkbox only when `currentUser.role === AccountRole.RoleEditor || currentUser.role === AccountRole.RoleAdmin`, and forces `is_official: false` on submit for `LEARNER` accounts.

## Verification

### Static Pattern Check
Verified zero matches for `flag_emoji`, `region_code`, `country_code`, `/studio/study`, `@bindings/palabrario/app/app`, `due_cards_today`, `retention_percent`, and `easiness_factor` across all created and modified files.

### `wails3 task check` (executed in `app/`)
Exited with code `0`:
- `wails3 generate bindings -clean=true -ts`: `Processed: 264 Packages, 1 Service, 98 Methods, 6 Enums, 75 Models, 0 Events in 6.2950773s.`
- `bun run --cwd frontend type-check` (`vue-tsc --noEmit`): passed with 0 errors.
- `bun run --cwd frontend build` (`vite build --mode production`): `✓ 2511 modules transformed`, `✓ built in 5.09s` (including `DeckFormDialog`, `StudioView`, and `StudioDeckDetailView` chunks).
- `go build .`: passed with 0 errors.
