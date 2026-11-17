# 038 — Changelog: Studio Collection Workspace (`/studio/collections/:id`) & Path Composer

## Summary of Changes

### Routes
- **`app/frontend/src/router/index.ts`**: Registered `/studio/collections/:id` (`name: "studio-collection-detail"`) under `AppLayout` with `meta: { requiresAuth: true, breadcrumb: [{ label: "Home", to: "/" }, { label: "Studio", to: "/studio" }, { label: "Collection" }] }`, lazy-loading `StudioCollectionDetailView.vue`.

### Components Created & Reused
- **`app/frontend/src/components/collections/CollectionDeckItem.vue`**: Extended with optional `actionLabel` (default `"Explore Deck →"`), `actionTo`, `stepNumber`, `showRegions`, and `#actions` slot while keeping `/collections/:id` default behavior (`@explore` button, `"Explore Deck →"`) completely unchanged.
- **`app/frontend/src/components/studio/CollectionFormDialog.vue`**: Create/Edit collection metadata dialog built on `AppDialog.vue` (`size="md"`). Validates trimmed `title` length (`3..150` chars), supports optional `description` and `is_public` toggle, sends `deck_ids: []` on creation, and preserves `(collection.decks ?? []).map(d => d.id_deck)` on metadata updates.
- **`app/frontend/src/components/studio/CollectionDeckPickerDialog.vue`**: Deck selector dialog built on `AppDialog.vue` (`size="lg"`). Fetches `api.studio.listStudioDecks()` and `api.studio.listDecks()`, merges and deduplicates by `id_deck`, excludes `existingDeckIds`, supports client-side filtering, displays `title`, `item_count`, full `region_names`, and `Official`/`Private` badges, and emits `select(deckId)` with inline error surfacing.
- **`app/frontend/src/components/studio/DeleteConfirmDialog.vue`** *(reused from Ticket 037)*: Used in `StudioCollectionDetailView.vue` to confirm permanent collection deletion before calling `api.studio.deleteCollection`.

### Views Created & Updated
- **`app/frontend/src/views/StudioView.vue`**: Preserved existing copy, layout, and styling; added the `New Collection` button alongside the `02 / Your Collections` filter tabs, wired to `CollectionFormDialog` with navigation to `/studio/collections/${collection.id_collection}` upon creation.
- **`app/frontend/src/views/studio/StudioCollectionDetailView.vue`**: Implemented `/studio/collections/:id`:
  - Validates positive integer `:id` and loads `api.studio.getCollection(id)` plus `api.studio.listStudioCollections()` (for community bookmark state) using a monotonic request counter.
  - Updates `breadcrumbLeafLabel` to `collection.title` and clears it on unmount.
  - Reuses `CollectionDeckItem.vue` for each deck in the ordered path, linking every deck to `/studio/decks/:id` (`Open Deck →`) and displaying full `region_names`.
  - Re-fetches `api.studio.getCollection(id)` after every `api.studio.updateCollection(...)` (metadata edits, deck additions, removals, and `↑`/`↓` reordering) so hydrated `decks` and `region_names` stay in sync with Go `StudioService.UpdateCollection`.
  - Contains zero `/studio/study` dead links, zero SM-2 mastery fields, zero `flag_emoji` or dialect codes, and zero `@bindings/palabrario/app/app` imports.

## Ordered Deck Path Behavior (`deck_ids`)
- **Add Deck**: Appends the selected `deck_id` (excluding duplicates) to `existingDeckIds` and persists via `api.studio.updateCollection`, then re-hydrates via `api.studio.getCollection`.
- **Reorder (`↑` / `↓`)**: Moves the selected deck one slot up or down in `existingDeckIds` (accessible `aria-label` controls, disabled at bounds and while mutating), persists the ordered `deck_ids` via `api.studio.updateCollection`, and re-hydrates via `api.studio.getCollection`.
- **Remove**: Filters out the target `id_deck` from `existingDeckIds`, persists via `api.studio.updateCollection`, and re-hydrates via `api.studio.getCollection`.

## Business Rule Enforcement (`RN03`, `RN04`, `RN05`)
- **`RN03` (Ownership & Capability Split)**: `StudioCollectionDetailView.vue` computes `canManage = currentUser?.role === AccountRole.RoleAdmin || currentUser?.id_user === collection.id_user`. Managers see `Add Deck`, `Edit Metadata`, `Delete Collection`, and per-deck `↑`/`↓`/`Remove` controls; non-owners see read-only ordered decks linking to `/studio/decks/:id` plus `Save to Studio` / `Unsave from Studio`.
- **`RN04` (Private Collection Handling)**: Distinguishes invalid `:id` / `404` (`null`), `403` private collection access errors, and network failures with dedicated monochrome alert panels and `← Back to Studio` navigation.
- **`RN05` (Official & Featured Curation Gate)**: `CollectionFormDialog.vue` renders `Official Collection` (`is_official`) and `Featured on Home` (`is_featured`) toggles only when `currentUser.role === AccountRole.RoleEditor || currentUser.role === AccountRole.RoleAdmin`, and forces both to `false` for `LEARNER` accounts in both `CollectionFormDialog.vue` and `StudioCollectionDetailView.vue`.

## Verification

### Static Pattern Check
Verified zero matches for `flag_emoji`, `region_code`, `country_code`, `/studio/study`, `@bindings/palabrario/app/app`, `due_cards_today`, `retention_percent`, and `easiness_factor` across all modified and created files.

### `wails3 task check` (executed in `app/`)
Exited with code `0`:
- `wails3 generate bindings -clean=true -ts`: `Processed: 264 Packages, 1 Service, 98 Methods, 6 Enums, 75 Models, 0 Events in 5.5127395s.`
- `bun run --cwd frontend type-check` (`vue-tsc --noEmit`): passed with 0 errors.
- `bun run --cwd frontend build` (`vite build --mode production`): `✓ 2517 modules transformed`, `✓ built in 7.25s` (including `CollectionDeckItem`, `CollectionFormDialog`, `StudioView`, and `StudioCollectionDetailView` chunks).
- `go build .`: passed with 0 errors.
