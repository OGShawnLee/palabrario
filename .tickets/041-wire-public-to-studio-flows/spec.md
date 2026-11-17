# 041 — Spec: Wire Public Dictionary, Catalog, Collection & Arena Entry Points to Studio

- [x] `WordOfTheDayCard.vue` `"Add to Studio"` redirects anonymous users to `/auth/sign-in?redirect=...` and opens `AddToDeckModal.vue` for authenticated users to add `sense.id_sense` to an owned/new deck.
- [x] `FeaturedDeckCard.vue`, `DeckGridItem.vue`, and `DecksView.vue` reflect live `isSaved` state when authenticated, toggle `saveDeck`/`unsaveDeck`, and route public `"Explore →"` to `/decks/:id` while leaving `"Study Deck →"` untouched.
- [x] `CollectionDetailView.vue` hydrates `isSaved` from `listStudioCollections()`, gates `toggleSave` on authentication, routes deck exploration (and `ProfileLibrarySection.vue` public deck links) to `/decks/:id`, leaves `"Study First Deck →"` untouched, and wires `CollectionAdminBar.vue` edit/add-deck/delete actions.
- [x] `CreateMatchDialog.vue` includes Studio decks (`listStudioDecks`) alongside public decks (`listDecks`) when authenticated (deduplicated by `id_deck`), without throwing `401` for anonymous visitors.
- [x] Zero `flag_emoji`, dialect codes, or `@bindings/palabrario/app/app` imports appear in any touched file.
- [x] `wails3 task check` (from `app/`) exits with code `0`.
- [x] `git status --short` shows only files within the declared `Scope`.
- [x] No unauthorized `git commit`, `git push`, or `git merge` was run.
