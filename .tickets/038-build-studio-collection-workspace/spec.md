# 038 — Spec: Build Studio Collection Workspace (`/studio/collections/:id`) & Path Composer

- [x] `app/frontend/src/router/index.ts` registers `/studio/collections/:id` (`name: "studio-collection-detail"`) with `meta: { requiresAuth: true }` and `Home / Studio / Collection` breadcrumb.
- [x] `CollectionFormDialog.vue` validates `title` (`3..150` chars), restricts `is_official` and `is_featured` visibility to `EDITOR`/`ADMIN` (`RN05`), and preserves `deck_ids` on metadata updates.
- [x] `StudioView.vue` opens `CollectionFormDialog` from its `New Collection` action.
- [x] `CollectionDeckPickerDialog.vue` lets owners/`ADMIN` pick from Studio/public decks without duplicates.
- [x] `StudioCollectionDetailView.vue` reuses `CollectionDeckItem.vue`, enforces `RN03` (`canManage = ADMIN || id_user === collection.id_user`), persists ordered `deck_ids` on add/remove/reorder via `updateCollection`, links decks to `/studio/decks/:id`, supports `deleteCollection` (with confirmation), and supports `Save`/`Unsave` for community subscribers (with zero `/studio/study` dead links).
- [x] Zero `flag_emoji`, dialect codes, or `@bindings/palabrario/app/app` imports appear in the new components.
- [x] `wails3 task check` (from `app/`) exits with code `0`.
- [x] `git status --short` shows only files within the declared `Scope`.
- [x] No unauthorized `git commit`, `git push`, or `git merge` was run.
