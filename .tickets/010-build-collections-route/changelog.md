# 010 — Changelog

## Added
- `app/frontend/src/components/ui/AppBreadcrumb.vue`: universal breadcrumb driven by `route.meta.breadcrumb`; ancestors are links, the leaf is bold and non-clickable; renders nothing on `/` or when the meta is empty.
- `app/frontend/src/views/CollectionsView.vue`: `/collections` page ("Study Collections"). Loads data only through `api.studio.listCollections()`; featured spotlight, "More collections" grid, loading and empty states.
- `app/frontend/src/components/collections/FeaturedCollectionCard.vue`: 2-column desktop / stacked mobile spotlight with a "Collection paths" ordered deck list.
- `app/frontend/src/components/collections/CollectionGridItem.vue`: catalog item (Official/Curated tag, title, description, "Explore →").
- `app/frontend/src/components/collections/CollectionDetailDialog.vue`: Reka UI dialog listing decks, word counts, and Close / Save Collection actions (Save is a placeholder with no behavior yet).

## Changed
- `app/frontend/src/layouts/AppLayout.vue`: mounts `<AppBreadcrumb />` directly below `<TopBar />`.
- `app/frontend/src/router/index.ts`: registers lazy `/collections` route with `meta.breadcrumb` (`Home / Collections`).
- `app/frontend/src/vite-env.d.ts`: augments `vue-router` `RouteMeta` with `breadcrumb`.
- `app/frontend/src/api/contract/studio.contract.ts`: adds `listCollections(featured?: boolean): Promise<Collection[]>`.
- `app/frontend/src/api/live/wails/studio.wails.ts`: implements `listCollections` via `Wails.ListCollections`.
- `app/frontend/src/api/mock/mock.adapter.ts`: implements `listCollections` (with simulated latency, public-only, optional featured filter).
- `app/frontend/src/api/mock/fixtures/decks.fixture.ts`: adds `mockCollectionDetails` (full `Collection` objects with constituent `Deck`s derived from existing `mockDecks`), in Spanish.

## Notes
- `AppDialog.vue` does not exist in the repo, so the dialog uses Reka UI `Dialog*` primitives directly, matching `SearchModal.vue`.
- The `Deck` model has no region field, so no region (name or code) is rendered; no flag emojis or dialect codes appear anywhere.
- No `uno.config.ts` changes were needed or made; no new shortcuts proposed.
- Verification: `cd app && wails3 task check` exits 0.
- Nothing was committed.
