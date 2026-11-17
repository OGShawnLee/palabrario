# 037 — Build Studio Deck Workspace (`/studio/decks/:id`) & Create/Edit Dialog

| Field | Value |
| :--- | :--- |
| **ID** | 037 |
| **Status** | completed |
| **Model** | Gemini 4 Hydrogen, `thinking_level: medium` (`gemini-4-hydrogen`) |
| **Role** | Frontend Implementer (`implementation`) |
| **Spec** | [spec.md](spec.md) |

> **Specification Note:**
> Ticket 3 of 7 in the active Studio suite (`035`–`038`, `043`, `041`–`042`; depends on `036`). Implements deck creation/metadata editing (`CU10`, `RN03`, `RN05`) and the dedicated `/studio/decks/:id` sub-route (`Home / Studio / <Deck Title>` breadcrumb) for adding/removing dictionary senses, deleting owned decks, or inspecting/unsaving community decks (`RN04`). SM-2 Study CTAs (`039`) are deferred.

## Scope

May create or edit:
- `app/frontend/src/router/index.ts`
- `app/frontend/src/views/StudioView.vue` (connect `New Deck` dialog trigger)
- `app/frontend/src/views/studio/StudioDeckDetailView.vue` (create)
- `app/frontend/src/components/studio/DeckFormDialog.vue` (create)
- `app/frontend/src/components/studio/AddSenseToDeckDialog.vue` (create)
- `app/frontend/src/components/studio/DeleteConfirmDialog.vue` (create, reusable for Decks & Collections)
- `.tickets/037-build-studio-deck-workspace/changelog.md` (create)

Must not touch:
- `core-service/`, `realtime-service/`, `image-service/`, `pkg/`, `app/*.go`, `SCHEMA_FINAL.sql`, `sql/`
- `app/frontend/uno.config.ts`
- `app/frontend/src/api/`
- Any file outside the above scope
- **No role (Router, Investigator, Implementer, Translator, or Reviewer) may commit, push, or merge without explicit user authorization.** Do not run `git commit`, `git push`, or `git merge`.

## Inputs

- [.routing/CONTEXT.md](../../.routing/CONTEXT.md), [.routing/DESIGN.md](../../.routing/DESIGN.md), [.routing/ROUTES.md](../../.routing/ROUTES.md).
- [.routing/services/studio_service.md](../../.routing/services/studio_service.md) (`RN03`, `RN04`, `RN05`, `CreateDeckRequest`, `UpdateDeckRequest`, `AddDeckItemRequest`).
- [app/frontend/src/components/ui/AppDialog.vue](../../app/frontend/src/components/ui/AppDialog.vue) and [app/frontend/src/components/ui/AppBreadcrumb.vue](../../app/frontend/src/components/ui/AppBreadcrumb.vue) (`breadcrumbLeafLabel`).
- [app/frontend/src/components/decks/FeaturedDeckCard.vue](../../app/frontend/src/components/decks/FeaturedDeckCard.vue) and [app/frontend/src/components/dictionary/WordOfTheDayCard.vue](../../app/frontend/src/components/dictionary/WordOfTheDayCard.vue) (lemma/sense display conventions).
- [app/frontend/src/composables/useAuth.ts](../../app/frontend/src/composables/useAuth.ts) (`currentUser.role`, `currentUser.id_user`).

## Task

1. **Register `/studio/decks/:id` in [router/index.ts](../../app/frontend/src/router/index.ts):**
   - Route name `studio-deck-detail`, `meta: { requiresAuth: true, breadcrumb: [{ label: "Home", to: "/" }, { label: "Studio", to: "/studio" }, { label: "Deck" }] }`, component `StudioDeckDetailView.vue`.

2. **Create `DeckFormDialog.vue` (Create & Edit modes):**
   - Uses `AppDialog.vue` (`rounded-none`, `bg-white`).
   - Fields: `title` (validated `3..150` trimmed chars), `description` (optional), `is_public` toggle/checkbox.
   - **`RN05` UI Gate:** Render the `Official Deck` (`is_official`) toggle **only** when `currentUser.role === "EDITOR" || currentUser.role === "ADMIN"`. For `LEARNER`, always submit `is_official: false`.
   - Submits via `api.studio.createDeck` or `api.studio.updateDeck`, shows inline `role="alert"` errors on failure, and emits the saved `Deck`.
   - Add the `New Deck` action in `StudioView.vue` to open `DeckFormDialog` and navigate to `/studio/decks/${deck.id_deck}` on creation.

3. **Create `AddSenseToDeckDialog.vue`:**
   - Lets the deck owner/`ADMIN` search lemmas via `api.dictionary.searchWords(query)` (debounced), fetch the selected word via `api.dictionary.getWordById(id_word)`, inspect its regional senses (full `region_name`, `cefr_level`, `short_definition`; exclude senses already in `deck.items`), and call `api.studio.addDeckItem(deckId, { id_sense })`.
   - Surface inline error banner if `addDeckItem` fails (e.g., `403` `RN03`).

4. **Create `StudioDeckDetailView.vue` (`/studio/decks/:id`):**
   - Validates positive numeric `:id`; fetches `api.studio.getDeck(id)`. Sets `breadcrumbLeafLabel` to `deck.title` (clears on unmount).
   - Handles `404`/`null`, `403` (`RN04` private deck belonging to another user), and network failures with clear monochrome alert panels + back link to `/studio`.
   - **`RN03` Capability Split:** `canManage = currentUser?.role === "ADMIN" || currentUser?.id_user === deck.id_user`:
     - When `canManage` is true: show `Edit Metadata`, `Add Card` (`AddSenseToDeckDialog`), per-card `Remove` (`api.studio.removeDeckItem`), and `Delete Deck` (`DeleteConfirmDialog` → `api.studio.deleteDeck` → `/studio`).
     - When `canManage` is false (community deck): hide mutating controls and show `Save to Studio` / `Unsave from Studio` (`api.studio.saveDeck` / `api.studio.unsaveDeck`). Do **not** add `/studio/study` links while Mastery/Study is deferred.
   - Show an Arena readiness note (`item_count >= 8` ready for Arena with a `"Play in Arena →"` link to `/arena` vs. `Add ${8 - item_count} more cards to host Arena matches`) and display each card's resolved lemma (linked to `/dictionary/word/${sense.id_word}` when available), `region_name`, `cefr_level`, `short_definition`, `definition`, and `examples` (zero `flag_emoji` or dialect codes).

5. **Verify:**
   - Run `wails3 task check` in `app/` and confirm exit code `0`.

## Verification

- Run all checks in [spec.md](spec.md).

## Change log

- Create `.tickets/037-build-studio-deck-workspace/changelog.md` recording components created, `RN03`/`RN04`/`RN05` UI enforcement, and `wails3 task check` output.
