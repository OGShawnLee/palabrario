# 037 — Spec: Build Studio Deck Workspace (`/studio/decks/:id`) & Create/Edit Dialog

- [x] `app/frontend/src/router/index.ts` registers `/studio/decks/:id` (`name: "studio-deck-detail"`) with `meta: { requiresAuth: true }` and `Home / Studio / Deck` breadcrumb.
- [x] `DeckFormDialog.vue` validates `title` length (`3..150` chars), restricts `is_official` visibility to `EDITOR`/`ADMIN` (`RN05`), and creates/updates decks via `api.studio`.
- [x] `StudioView.vue` opens `DeckFormDialog` from its `New Deck` action.
- [x] `AddSenseToDeckDialog.vue` searches dictionary words via `api.dictionary.searchWords`/`getWordById`, prevents duplicate sense additions, displays full `region_name` (`cefr_level`), and calls `api.studio.addDeckItem`.
- [x] `StudioDeckDetailView.vue` enforces `RN03` (`canManage = ADMIN || id_user === deck.id_user`), handles private `403`/`404` states (`RN04`), supports `removeDeckItem` + `deleteDeck` (with confirmation) for owners/`ADMIN`, offers `Save`/`Unsave` for community decks, and renders a `"Play in Arena →"` link to `/arena` when `item_count >= 8` (with zero `/studio/study` dead links).
- [x] Zero `flag_emoji`, dialect codes, or `@bindings/palabrario/app/app` imports exist in the new components.
- [x] `wails3 task check` (from `app/`) exits with code `0`.
- [x] `git status --short` shows only files within the declared `Scope`.
- [x] No unauthorized `git commit`, `git push`, or `git merge` was run.
