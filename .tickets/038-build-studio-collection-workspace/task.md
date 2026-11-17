# 038 — Build Studio Collection Workspace (`/studio/collections/:id`) & Path Composer

| Field | Value |
| :--- | :--- |
| **ID** | 038 |
| **Status** | completed |
| **Model** | Gemini 4 Hydrogen, `thinking_level: medium` (`gemini-4-hydrogen`) |
| **Role** | Frontend Implementer (`implementation`) |
| **Spec** | [spec.md](spec.md) |

> **Specification Note:**
> Ticket 4 of 7 in the active Studio suite (`035`–`038`, `043`, `041`–`042`; depends on `037`). Implements collection creation/metadata editing (`CU11`, `RN03`, `RN05`), ordered learning-path deck composition (`deck_ids`), and the dedicated `/studio/collections/:id` sub-route (`Home / Studio / <Collection Title>` breadcrumb), reusing existing `CollectionDeckItem.vue` aesthetics for owners/`ADMIN` and saved community subscribers (`RN04`). SM-2 Study links (`039`) are deferred.

## Scope

May create or edit:
- `app/frontend/src/router/index.ts`
- `app/frontend/src/views/StudioView.vue` (connect `New Collection` dialog trigger)
- `app/frontend/src/views/studio/StudioCollectionDetailView.vue` (create)
- `app/frontend/src/components/studio/CollectionFormDialog.vue` (create)
- `app/frontend/src/components/studio/CollectionDeckPickerDialog.vue` (create)
- `app/frontend/src/components/collections/CollectionDeckItem.vue` (optional `#actions` slot if needed for reuse without breaking public `/collections/:id`)
- `.tickets/038-build-studio-collection-workspace/changelog.md` (create)

Must not touch:
- `core-service/`, `realtime-service/`, `image-service/`, `pkg/`, `app/*.go`, `SCHEMA_FINAL.sql`, `sql/`
- `app/frontend/uno.config.ts`
- `app/frontend/src/api/`
- Any file outside the above scope
- **No role (Router, Investigator, Implementer, Translator, or Reviewer) may commit, push, or merge without explicit user authorization.** Do not run `git commit`, `git push`, or `git merge`.

## Inputs

- [.routing/CONTEXT.md](../../.routing/CONTEXT.md), [.routing/DESIGN.md](../../.routing/DESIGN.md), [.routing/ROUTES.md](../../.routing/ROUTES.md).
- [.routing/services/studio_service.md](../../.routing/services/studio_service.md) (`CU11`, `RN03`, `RN04`, `RN05`, `CreateCollectionRequest`, `UpdateCollectionRequest`).
- [app/frontend/src/views/CollectionDetailView.vue](../../app/frontend/src/views/CollectionDetailView.vue), [CollectionHeader.vue](../../app/frontend/src/components/collections/CollectionHeader.vue), and [CollectionDeckItem.vue](../../app/frontend/src/components/collections/CollectionDeckItem.vue) (UI reuse baseline).
- [app/frontend/src/components/studio/DeleteConfirmDialog.vue](../../app/frontend/src/components/studio/DeleteConfirmDialog.vue) (from `037`).
- [app/frontend/src/composables/useAuth.ts](../../app/frontend/src/composables/useAuth.ts).

## Task

1. **Register `/studio/collections/:id` in [router/index.ts](../../app/frontend/src/router/index.ts):**
   - Route name `studio-collection-detail`, `meta: { requiresAuth: true, breadcrumb: [{ label: "Home", to: "/" }, { label: "Studio", to: "/studio" }, { label: "Collection" }] }`, component `StudioCollectionDetailView.vue`.

2. **Create `CollectionFormDialog.vue` (Create & Edit modes):**
   - Uses `AppDialog.vue`.
   - Fields: `title` (`3..150` trimmed chars), `description` (optional), `is_public` toggle.
   - **`RN05` UI Gate:** Render `Official Collection` (`is_official`) and `Featured on Home` (`is_featured`) toggles **only** when `currentUser.role === "EDITOR" || currentUser.role === "ADMIN"`. Force both to `false` for `LEARNER`.
   - Submits metadata only (`deck_ids: []` on create; preserving existing `deck_ids` on update) via `api.studio.createCollection` / `api.studio.updateCollection`.
   - Add the `New Collection` action in `StudioView.vue` to open `CollectionFormDialog` and navigate to `/studio/collections/${col.id_collection}` on creation (where users add and order decks via `CollectionDeckPickerDialog`).

3. **Create `CollectionDeckPickerDialog.vue`:**
   - Fetches `api.studio.listStudioDecks()` and `api.studio.listDecks()`, merges/deduplicates by `id_deck`, excludes decks already in the collection, shows `title`, `item_count`, full `region_names`, and `Official`/`Private` tags, and emits the selected `deck_id` to append to `deck_ids`.

4. **Create `StudioCollectionDetailView.vue` (`/studio/collections/:id`):**
   - Fetches `api.studio.getCollection(id)`; sets `breadcrumbLeafLabel` to `collection.title` (clears on unmount).
   - Handles `404`/`null`, `403` (`RN04` private collection), and retryable error states.
   - Reuses [`CollectionDeckItem.vue`](../../app/frontend/src/components/collections/CollectionDeckItem.vue) (extending with an optional `#actions` slot if helpful, keeping `/collections/:id` default behavior unchanged) so deck cards share the exact public dictionary/collections visual grammar.
   - **`RN03` Capability Split:** `canManage = currentUser?.role === "ADMIN" || currentUser?.id_user === collection.id_user`:
     - When `canManage` is true: show `Edit Metadata`, `Add Deck` (`CollectionDeckPickerDialog`), accessible move-up / move-down (`↑` / `↓`) reorder controls and `Remove` per deck (persisting ordered `deck_ids` via `api.studio.updateCollection`), `Open Deck →` (`/studio/decks/:id`), and `Delete Collection` (`DeleteConfirmDialog` → `api.studio.deleteCollection` → `/studio`).
     - When `canManage` is false (community subscription): render read-only ordered decks linking to `/studio/decks/:id` and provide `Save to Studio` / `Unsave from Studio` (`api.studio.saveCollection` / `api.studio.unsaveCollection`). Zero `/studio/study` links.

5. **Verify:**
   - Run `wails3 task check` in `app/` and confirm exit code `0`.

## Verification

- Run all checks in [spec.md](spec.md).

## Change log

- Create `.tickets/038-build-studio-collection-workspace/changelog.md` recording components created/reused, deck ordering behavior, `RN03`/`RN05` UI checks, and `wails3 task check` output.
