# 041 — Wire Public Dictionary, Catalog, Collection & Arena Entry Points to Studio

| Field | Value |
| :--- | :--- |
| **ID** | 041 |
| **Status** | completed |
| **Model** | Gemini 4 Hydrogen, `thinking_level: medium` (`gemini-4-hydrogen`) |
| **Role** | Frontend Implementer (`implementation`) |
| **Spec** | [spec.md](spec.md) |

> **Specification Note:**
> Ticket 6 of 7 in the active Studio suite (`035`–`038`, `043`, `041`–`042`; depends on `037`, `038`, `043`). Connects all existing public and cross-route Studio stubs (`WordOfTheDayCard` `"Add to Studio"`, `FeaturedDeckCard` / `DeckGridItem` `"Save to Studio"` & `/decks/:id`, `CollectionDetailView` / `CollectionAdminBar`, `ProfileLibrarySection` deck links, and `CreateMatchDialog` in `/arena`) to live Studio APIs and `/decks/:id` / `/studio/*` routes. Existing `"Study Deck →"` / `"Study First Deck →"` labels stay unchanged (no temporary relabeling) until `039` is redesigned.

## Scope

May create or edit:
- `app/frontend/src/components/dictionary/WordOfTheDayCard.vue`
- `app/frontend/src/components/studio/AddToDeckModal.vue` (create)
- `app/frontend/src/components/decks/FeaturedDeckCard.vue`
- `app/frontend/src/components/decks/DeckGridItem.vue`
- `app/frontend/src/views/DecksView.vue`
- `app/frontend/src/views/CollectionDetailView.vue`
- `app/frontend/src/components/collections/CollectionAdminBar.vue`
- `app/frontend/src/components/profile/ProfileLibrarySection.vue` (route public deck rows to `/decks/:id`)
- `app/frontend/src/components/arena/CreateMatchDialog.vue`
- `.tickets/041-wire-public-to-studio-flows/changelog.md` (create)

Must not touch:
- `core-service/`, `realtime-service/`, `image-service/`, `pkg/`, `app/*.go`, `SCHEMA_FINAL.sql`, `sql/`
- `app/frontend/uno.config.ts`
- `app/frontend/src/api/`
- Any file outside the above scope
- **No role (Router, Investigator, Implementer, Translator, or Reviewer) may commit, push, or merge without explicit user authorization.** Do not run `git commit`, `git push`, or `git merge`.

## Inputs

- [.routing/CONTEXT.md](../../.routing/CONTEXT.md), [.routing/DESIGN.md](../../.routing/DESIGN.md), [.routing/ROUTES.md](../../.routing/ROUTES.md).
- [app/frontend/src/components/dictionary/WordOfTheDayCard.vue](../../app/frontend/src/components/dictionary/WordOfTheDayCard.vue).
- [app/frontend/src/components/decks/FeaturedDeckCard.vue](../../app/frontend/src/components/decks/FeaturedDeckCard.vue) and [DeckGridItem.vue](../../app/frontend/src/components/decks/DeckGridItem.vue).
- [app/frontend/src/views/CollectionDetailView.vue](../../app/frontend/src/views/CollectionDetailView.vue) and [CollectionAdminBar.vue](../../app/frontend/src/components/collections/CollectionAdminBar.vue).
- [app/frontend/src/components/profile/ProfileLibrarySection.vue](../../app/frontend/src/components/profile/ProfileLibrarySection.vue).
- [app/frontend/src/components/arena/CreateMatchDialog.vue](../../app/frontend/src/components/arena/CreateMatchDialog.vue).
- [app/frontend/src/composables/useAuth.ts](../../app/frontend/src/composables/useAuth.ts).

## Task

1. **Wire `"Add to Studio"` in [WordOfTheDayCard.vue](../../app/frontend/src/components/dictionary/WordOfTheDayCard.vue) (`CU10 3a`):**
   - Create `AddToDeckModal.vue` (`AppDialog`): loads `api.studio.listStudioDecks()` filtered to decks owned by `currentUser` (`!d.is_saved_from_community` or `currentUser.role === "ADMIN"`), allows picking an existing owned deck or quick-creating a new deck via `api.studio.createDeck`, and calls `api.studio.addDeckItem(deckId, { id_sense })` with inline feedback and error handling.
   - When `"Add to Studio"` is clicked by an anonymous user (`!isAuthenticated`), redirect to `/auth/sign-in` with `?redirect=${route.fullPath}` instead of opening the modal to a `401`.

2. **Wire Public `/decks` ([DecksView.vue](../../app/frontend/src/views/DecksView.vue), [FeaturedDeckCard.vue](../../app/frontend/src/components/decks/FeaturedDeckCard.vue), [DeckGridItem.vue](../../app/frontend/src/components/decks/DeckGridItem.vue)):**
   - Hydrate saved deck IDs when `isAuthenticated` (via `api.studio.listStudioDecks()`), wire `Save to Studio` / `Saved in Studio` toggle (`api.studio.saveDeck` / `unsaveDeck`, redirecting to `/auth/sign-in?redirect=/decks` if anonymous).
   - Keep `DeckGridItem.vue` `"Explore →"` routing public decks to `/decks/${deck.id_deck}` (`043`). Keep `FeaturedDeckCard.vue` `"Study Deck →"` label untouched (routing `/studio` until `039` lands—no temporary label churn).
   - Display `deck.region_names` (`badge-side`, max 3, full names) if missing on public cards.

3. **Wire [CollectionDetailView.vue](../../app/frontend/src/views/CollectionDetailView.vue), [CollectionAdminBar.vue](../../app/frontend/src/components/collections/CollectionAdminBar.vue), and [ProfileLibrarySection.vue](../../app/frontend/src/components/profile/ProfileLibrarySection.vue):**
   - Replace hardcoded `isSaved = false` in `CollectionDetailView.vue` by checking `api.studio.listStudioCollections()` when authenticated (`useAuth().currentUser`).
   - On `toggleSave` when `!isAuthenticated`, redirect to `/auth/sign-in?redirect=${route.fullPath}`.
   - Route `CollectionDeckItem` `@explore` and `ProfileLibrarySection.vue` public deck links to `/decks/${deck.id_deck}` (`043`). Keep `CollectionHeader.vue` `"Study First Deck →"` label untouched.
   - Wire `CollectionAdminBar` actions (`Edit Collection`, `Add Deck` → `/studio/collections/${collection.id_collection}`, `Delete` → confirmation + `api.studio.deleteCollection` → `/collections`).

4. **Upgrade [CreateMatchDialog.vue](../../app/frontend/src/components/arena/CreateMatchDialog.vue):**
   - When `isAuthenticated`, load both `api.studio.listStudioDecks()` and `api.studio.listDecks()` (deduplicated by `id_deck`, Studio decks listed first); when anonymous, load `api.studio.listDecks()` without triggering a `401`.

5. **Verify:**
   - Run `wails3 task check` in `app/` and confirm exit code `0`.

## Verification

- Run all checks in [spec.md](spec.md).

## Change log

- Create `.tickets/041-wire-public-to-studio-flows/changelog.md` recording each wired component, auth gate behavior, and `wails3 task check` output.
