# 010 — Spec checklist

Run by the router before accepting the ticket. All items must pass.

- [ ] Route `/collections` is registered under `AppLayout` in `app/frontend/src/router/index.ts` with `meta.breadcrumb`.
- [ ] Universal breadcrumb component `AppBreadcrumb.vue` created and mounted in `AppLayout.vue`.
- [ ] Breadcrumb correctly hides on `/` (Home) and displays hierarchy (e.g. `Home / Collections`) on `/collections`.
- [ ] `CollectionsView.vue` implements the approved layout:
  - Big headline `Study Collections` with proper vertical whitespace.
  - Featured collection spotlight card (`FeaturedCollectionCard.vue`) with 2-column desktop / 1-column mobile layout.
  - Right column lists constituent decks with descriptions and hairline dividers.
- [ ] `CollectionsView.vue` consumes collections from `api.studio.listCollections()` via `import { api } from "@/api"`. Zero hardcoded mock arrays or static collections exist inside the Vue view or its child components.
- [ ] `IStudioAPI` contract and adapters (Wails and Mock) properly support `listCollections(featured?: boolean)`.
- [ ] Inspection dialog (`CollectionDetailDialog.vue`) opens on "Explore Collection" click and lists path details.
- [ ] Section labels and badges use natural sentence case (`section-label` / `text-sm font-bold text-dull`) instead of all-caps.
- [ ] Strictly zero `flag_emoji` and zero dialect codes (`CU`, `CL`, `MX`). Full country names only.
- [ ] Strict monochrome styling (no rounded corners, no colored badges).
- [ ] `cd app && wails3 task check` exits with code 0.
- [ ] Changes stay strictly within the declared scope.
- [ ] `.tickets/010-build-collections-route/changelog.md` exists and is populated.
- [ ] Commit (if made): Conventional Commits, English, detailed body.
