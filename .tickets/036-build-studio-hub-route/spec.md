# 036 — Spec: Build Studio Hub Route (`/studio`)

- [x] `app/frontend/src/views/StudioView.vue` replaces the placeholder and loads `listStudioDecks()` and `listStudioCollections()` via `@/api` with spacious broadsheet layout (`gap-16`, `gap-y-12 md:gap-x-16`), loading, error (`role="alert"` + Retry), and empty states, without temporary disabled buttons.
- [x] Both **Your Decks** and **Your Collections** sections in `StudioView.vue` support `All` / `Owned` / `Saved` filter tabs and use a quiet text `Unsave` action so `Open →` stays the single bordered action per card.
- [x] `DeckGridItem.vue` and `CollectionGridItem.vue` are reused cleanly across public catalogs and `/studio` using quiet typographic metadata lines (`Owned`/`Saved`, `Public`/`Private`, `Official`/`Featured`, full `region_names`, safe omitted-`decks` fallback) instead of cluttered badge stacks, linking to `/studio/decks/:id` / `/studio/collections/:id` without SM-2 mastery stats or `flag_emoji`/dialect codes.
- [x] `app/frontend/src/router/index.ts` keeps `meta.requiresAuth: true` on `/studio` and defines the `Home / Studio` breadcrumb.
- [x] Zero imports from `@bindings/palabrario/app/app` appear in any `.vue` file.
- [x] `wails3 task check` (from `app/`) exits with code `0`.
- [x] `git status --short` shows only files within the declared `Scope` (and `app/frontend/uno.config.ts` is untouched).
- [x] No unauthorized `git commit`, `git push`, or `git merge` was run.
