# 021 — Changelog

## Added
- `views/CollectionDetailView.vue`: `/collections/:id` page. Loads the collection via `api.studio.getCollection(id)` and the session via `api.auth.getMe()` (reloads when the id changes); loading and not-found states; save toggle via `saveCollection` / `unsaveCollection`.
- `components/collections/CollectionHeader.vue`: eyebrow, title, blurb, metadata line (deck count, term count, updated date), "Save to Studio" / "✓ Saved in Studio" and "Study First Deck →" (stacked on mobile).
- `components/collections/CollectionDeckItem.vue`: deck card with term count, title, blurb, "Explore Deck →".
- `components/collections/CollectionAdminBar.vue`: owner/admin bar (Edit Collection, Add Deck, Delete). Buttons are placeholders with no behavior yet.

## Changed
- `router/index.ts`: registers `collections/:id` (`collection-detail`) with breadcrumb meta.
- `api/contract/studio.contract.ts`, `api/live/wails/studio.wails.ts`, `api/mock/mock.adapter.ts`: added `getCollection`, `saveCollection`, `unsaveCollection` (Wails delegates to `Wails.GetCollection/SaveCollection/UnsaveCollection`; mock tracks saved ids in memory).
- `api/mock/fixtures/decks.fixture.ts`: the `palabrario` owner id changed to `...0099` because it collided with the mock learner's id, which would have made the learner owner of official collections. `caribe_viva`'s collection stays owned by the mock editor (`...0002`), so the owner bar can be seen when signed in as the editor.
- `SearchModal.vue`: selecting a collection navigates to `/collections/:id`.
- `CollectionsView.vue`, `FeaturedCollectionCard.vue`, `CollectionGridItem.vue`: dialog state removed; explore navigates with `router.push` to `/collections/:id`.

## Removed
- `components/collections/CollectionDetailDialog.vue` (retired interim dialog).

## Deviations / notes
- **Out-of-scope edit:** `components/ui/AppBreadcrumb.vue` now exports a `breadcrumbLeafLabel` ref that replaces the last crumb's label. This was needed to show `Home / Collections / [Collection Title]`, because static route meta can't hold a loaded title. The detail view sets it after load and clears it on unmount. Please approve or suggest another approach.
- **Owner check:** the `Collection` DTO has no `owner_user_name`, so ownership is `currentUser.id_user === collection.id_user` (or `RoleAdmin`).
- **Saved state:** the DTO has no saved flag, so the Save button starts as "Save to Studio" on every visit and only tracks saves made in the current view.
- "Study First Deck →" and "Explore Deck →" go to `/studio`, since no deck route exists yet.
- Metadata omits a curator name (not in the DTO); it shows "Official curriculum" or "Community curated" instead.
- No `uno.config.ts` changes or proposals.
- Verification: `cd app && wails3 task check` exits 0. Nothing committed.
