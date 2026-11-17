# 043 — Spec: Build Public Deck Detail Route (`/decks/:id`)

- [x] `app/frontend/src/router/index.ts` registers public route `/decks/:id` (`name: "deck-detail"`) with `Home / Decks / Deck Detail` breadcrumb.
- [x] `DeckDetailView.vue` loads `api.studio.getDeck(id)`, sets `breadcrumbLeafLabel` to `deck.title`, and handles invalid ID / `404` / `403` (`RN04`) / retryable errors.
- [x] `DeckDetailView.vue` hydrates `isSaved` when authenticated, gates `Save to Studio` on authentication (redirecting guests to `/auth/sign-in?redirect=...`), links owners/`ADMIN` to `/studio/decks/:id`, and links `8+` card decks to `/arena`.
- [x] Renders deck cards with resolved lemmas (`/dictionary/word/:id_word`), full `region_name`, `cefr_level`, `short_definition`, and `examples` with zero `flag_emoji`, dialect codes, or `@bindings/palabrario/app/app` imports.
- [x] `wails3 task check` (from `app/`) exits with code `0`.
- [x] `git status --short` shows only files within the declared `Scope`.
- [x] No unauthorized `git commit`, `git push`, or `git merge` was run.
