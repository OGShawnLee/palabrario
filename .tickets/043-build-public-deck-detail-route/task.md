# 043 — Build Public Deck Detail Route (`/decks/:id`)

| Field | Value |
| :--- | :--- |
| **ID** | 043 |
| **Status** | completed |
| **Model** | Gemini 4 Hydrogen, `thinking_level: medium` (`gemini-4-hydrogen`) |
| **Role** | Frontend Implementer (`implementation`) |
| **Spec** | [spec.md](spec.md) |

> **Specification Note:**
> Ticket 5 of 7 in the active Studio/Decks suite (`035`–`038`, `043`, `041`–`042`). Creates the missing public `/decks/:id` detail route (`name: "deck-detail"`, `Public` access, `Home / Decks / <Deck Title>` breadcrumb) so public catalog links (`DeckGridItem.vue` on `/decks`, `CollectionDeckItem.vue` on `/collections/:id`, and `ProfileLibrarySection.vue` on `/u/:username`) inspect public decks without forcing anonymous visitors through `/auth/sign-in` into `/studio/decks/:id`.

## Scope

May create or edit:
- `app/frontend/src/router/index.ts` (register `/decks/:id`)
- `app/frontend/src/views/DeckDetailView.vue` (create)
- `.tickets/043-build-public-deck-detail-route/changelog.md` (create)

Must not touch:
- `core-service/`, `realtime-service/`, `image-service/`, `pkg/`, `app/*.go`, `SCHEMA_FINAL.sql`, `sql/`
- `app/frontend/uno.config.ts`
- `app/frontend/src/api/`
- Any file outside the above scope
- **No role (Router, Investigator, Implementer, Translator, or Reviewer) may commit, push, or merge without explicit user authorization.** Do not run `git commit`, `git push`, or `git merge`.

## Inputs

- [.routing/CONTEXT.md](../../.routing/CONTEXT.md), [.routing/DESIGN.md](../../.routing/DESIGN.md), [.routing/ROUTES.md](../../.routing/ROUTES.md).
- [.routing/services/studio_service.md](../../.routing/services/studio_service.md) (`RN03`, `RN04`, `Deck`, `DeckItem`).
- [app/frontend/src/views/CollectionDetailView.vue](../../app/frontend/src/views/CollectionDetailView.vue) and [app/frontend/src/views/studio/StudioDeckDetailView.vue](../../app/frontend/src/views/studio/StudioDeckDetailView.vue) (from `037`, UI alignment reference).
- [app/frontend/src/components/ui/AppBreadcrumb.vue](../../app/frontend/src/components/ui/AppBreadcrumb.vue) (`breadcrumbLeafLabel`).
- [app/frontend/src/composables/useAuth.ts](../../app/frontend/src/composables/useAuth.ts).

## Task

1. **Register `/decks/:id` in [router/index.ts](../../app/frontend/src/router/index.ts):**
   - Route name `deck-detail`, `meta: { breadcrumb: [{ label: "Home", to: "/" }, { label: "Decks", to: "/decks" }, { label: "Deck Detail" }] }`, component `DeckDetailView.vue` (public route, no `requiresAuth`).

2. **Create [DeckDetailView.vue](../../app/frontend/src/views/DeckDetailView.vue) (`/decks/:id`):**
   - Validate positive numeric `:id` and fetch `api.studio.getDeck(id)`. Set `breadcrumbLeafLabel` to `deck.title` (clear on unmount).
   - Handle invalid ID, `404`/`null`, `403` (`RN04` private deck), and retryable error states cleanly with monochrome alert banner and back link to `/decks`.
   - When `isAuthenticated`, hydrate `isSaved` via `api.studio.listStudioDecks()` and compute `canManage = currentUser?.role === "ADMIN" || currentUser?.id_user === deck.id_user`:
     - Show `Save to Studio` / `✓ Saved in Studio` button (`api.studio.saveDeck` / `unsaveDeck`; if anonymous, redirect to `/auth/sign-in?redirect=${route.fullPath}`).
     - When `canManage` is true, show an `"Open in Studio →"` link to `/studio/decks/${deck.id_deck}`.
     - When `deck.item_count >= 8`, show a `"Play in Arena →"` link to `/arena`.
   - Render deck metadata (`is_official`, `item_count` terms, full `region_names` via `badge-side`, zero `flag_emoji` or dialect codes) and each card (`lemma` resolved via `api.dictionary.getWordById` linking to `/dictionary/word/${sense.id_word}`, `region_name`, `cefr_level`, `short_definition`, `definition`, and `examples`).

3. **Verify:**
   - Run `wails3 task check` in `app/` and confirm exit code `0`.

## Verification

- Run all checks in [spec.md](spec.md).

## Change log

- Create `.tickets/043-build-public-deck-detail-route/changelog.md` recording `/decks/:id` implementation, auth/save gates, and `wails3 task check` output.
