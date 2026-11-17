# 041 — Changelog: Wire Public Dictionary, Catalog, Collection & Arena Entry Points to Studio

## Summary of Changes

### Components Created
- **`app/frontend/src/components/studio/AddToDeckModal.vue`**: Built on `AppDialog.vue` (`size="md"`) to add any dictionary `SenseDetail` to a user's Studio library:
  - Displays selected word lemma, `region_name` (`badge-side`), `cefr_level` (`badge-main`), and `short_definition`.
  - Loads `api.studio.listStudioDecks()` when opened (only when authenticated) and filters to manageable owned decks (`!d.is_saved_from_community || currentUser.role === AccountRole.RoleAdmin`).
  - Lets users add `id_sense` to an existing owned deck via `api.studio.addDeckItem(deckId, new AddDeckItemRequest({ id_sense }))` or quick-create a private deck (`api.studio.createDeck` followed by `api.studio.addDeckItem`) with title validation (`3..150` chars), inline `role="alert"` error handling, and confirmation linking directly to `/studio/decks/${deck.id_deck}`.

### Components & Views Updated
- **`app/frontend/src/components/dictionary/WordOfTheDayCard.vue`**:
  - Wired each sense's `"Add to Studio"` button via `handleAddToStudio(sense)`: anonymous visitors are redirected to `/auth/sign-in?redirect=${route.fullPath}`; authenticated users open `AddToDeckModal.vue` for that `sense` (`id_sense`, `lemma`).
  - Replaced `sense.region_code` / `"CL"` branching in `formatFrequency` with `sense.region_name` so no raw dialect codes appear in the component.
- **`app/frontend/src/components/decks/FeaturedDeckCard.vue`**:
  - Added `isSaved`, `isSaving`, and `isOwned` props plus `toggle-save` emit to reflect live saved state (`"Save to Studio"` / `"✓ Saved in Studio"`, `"In Studio →"` link for owned decks).
  - Rendered up to 3 full `region_names` using `badge-side` tags and added `"Explore →"` linking to `/decks/${deck.id_deck}` while keeping `"Study Deck →"` untouched.
- **`app/frontend/src/components/decks/DeckGridItem.vue`**:
  - Defaulted public destination (`to`) to `/decks/${deck.id_deck}` when omitted so `"Explore →"` routes public catalog cards to `/decks/:id` (preserving explicit `:to` overrides such as `/studio/decks/:id` in `StudioView.vue`).
  - Rendered up to 3 full `region_names` via `badge-side` tags when `showRegions` is enabled.
- **`app/frontend/src/views/DecksView.vue`**:
  - Hydrates `savedDeckIds` and `ownedDeckIds` from `api.studio.listStudioDecks()` only when `isAuthenticated` is `true` (and re-hydrates when `currentUser` changes), preventing `401` calls for anonymous visitors.
  - Wires `FeaturedDeckCard` and `DeckGridItem` `#actions` to toggle `api.studio.saveDeck` / `api.studio.unsaveDeck` with auth-gate redirect (`/auth/sign-in?redirect=/decks`) and inline `role="alert"` error banners, enabling `show-regions` and `:to="/decks/${deck.id_deck}"`.
- **`app/frontend/src/components/collections/CollectionAdminBar.vue`**:
  - Accepts `collectionId` and `isDeleting`; routes `Edit Collection` and `Add Deck` to `/studio/collections/${collectionId}` and emits `delete` from the `Delete` button.
- **`app/frontend/src/views/CollectionDetailView.vue`**:
  - Replaced direct `api.auth.getMe()` calls with `useAuth()`; hydrates `isSaved` from `api.studio.listStudioCollections()` only when `isAuthenticated` is `true` (and updates when `currentUser` changes).
  - Auth-gates `toggleSave()` (`!isAuthenticated` -> `/auth/sign-in?redirect=${route.fullPath}`) with inline `role="alert"` error handling.
  - Routes `CollectionDeckItem` `@explore` to `/decks/${deck.id_deck}` with `show-regions` enabled, while keeping `CollectionHeader.vue` `"Study First Deck →"` untouched.
  - Wires `CollectionAdminBar` (`Edit Collection` / `Add Deck` -> `/studio/collections/${collection.id_collection}`, `Delete` -> `DeleteConfirmDialog.vue` + `api.studio.deleteCollection` -> `/collections`).
- **`app/frontend/src/components/profile/ProfileLibrarySection.vue`**:
  - Updated public deck `RouterLink` destination from `/decks` to `:to="/decks/${deck.id_deck}"`.
- **`app/frontend/src/components/arena/CreateMatchDialog.vue`**:
  - Removed the early sign-in redirect on open so both guests and signed-in users can browse available decks cleanly; when `isAuthenticated`, loads `api.studio.listStudioDecks()` and `api.studio.listDecks()` in parallel (deduplicated by `id_deck`, Studio decks listed first); when anonymous, calls only `api.studio.listDecks()` without triggering a `401`. Match creation on submit remains gated on authentication.

## Verification

### Static Pattern Check
Verified zero matches for `flag_emoji`, `region_code`, `country_code`, `/studio/study`, `@bindings/palabrario/app/app`, `due_cards_today`, `retention_percent`, and `easiness_factor` across all 9 created/modified source files.

### `wails3 task check` (executed in `app/`)
Exited with code `0`:
- `wails3 generate bindings -clean=true -ts`: `Processed: 264 Packages, 1 Service, 98 Methods, 6 Enums, 75 Models, 0 Events in 8.3848038s.`
- `bun run --cwd frontend type-check` (`vue-tsc --noEmit`): passed with 0 errors.
- `bun run --cwd frontend build` (`vite build --mode production`): `✓ 2521 modules transformed.`, `✓ built in 12.18s`.
- `go build .`: passed with 0 errors.
