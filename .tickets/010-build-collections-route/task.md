# 010 — Build the Public Collections Route (/collections)

| Field | Value |
| :--- | :--- |
| **ID** | 010 |
| **Status** | in-progress |
| **Model** | Claude Sonnet 5.5, `effort: medium` (`claude-sonnet-5-5`) |
| **Role** | Frontend Implementer |
| **Spec** | [spec.md](spec.md) |

> **Specification Note:**  
> This ticket builds the public discovery route `/collections` following the approved desktop and mobile editorial wireframes. It introduces a universal breadcrumb component (`AppBreadcrumb.vue`) embedded in `AppLayout.vue` and controlled via `route.meta.breadcrumb`, a Featured Collection spotlight with its list of constituent learning path decks, and a clean 2-column "More Collections" grid using natural sentence-case labels (`section-label`).

## Scope

May create or edit:
- `app/frontend/src/components/ui/AppBreadcrumb.vue` (universal breadcrumb)
- `app/frontend/src/layouts/AppLayout.vue` (mount `AppBreadcrumb` under `TopBar`)
- `app/frontend/src/router/index.ts` (register `/collections` route and `meta.breadcrumb`)
- `app/frontend/src/views/CollectionsView.vue` (main public collections showcase)
- `app/frontend/src/components/collections/FeaturedCollectionCard.vue` (2-column desktop / 1-column mobile spotlight)
- `app/frontend/src/components/collections/CollectionGridItem.vue` (card for collections catalog)
- `app/frontend/src/components/collections/CollectionDetailDialog.vue` (inspection dialog showing decks and vocabulary counts)
- `app/frontend/src/api/contract/studio.contract.ts` (expand IStudioAPI with listCollections)
- `app/frontend/src/api/live/wails/studio.wails.ts` (implement listCollections)
- `app/frontend/src/api/mock/mock.adapter.ts` (implement mock listCollections)
- `app/frontend/src/api/mock/fixtures/decks.fixture.ts` (expand collection fixtures if needed)
- `app/frontend/src/vite-env.d.ts` (type augmentation for `RouteMeta.breadcrumb`)
- `.tickets/010-build-collections-route/changelog.md` (create)

Must not touch:
- `app/frontend/uno.config.ts` (propose shortcuts or tokens in changelog)
- Go backend code, `core-service/`, `realtime-service/`, or `.context/`
- Any files outside the above declared scope

## Inputs

- [.routing/API.md](../../.routing/API.md): Architectural rules (Rule 1: Never import `@bindings/.../app` in UI; always import `api` from `@/api`).
- [app/frontend/src/api/index.ts](../../app/frontend/src/api/index.ts): Universal `api` singleton (`api.studio`).
- [.routing/CONTEXT.md](../../.routing/CONTEXT.md): Architecture ground rules (strict prohibition of `flag_emoji` and regional codes; full region names required).
- [.routing/DESIGN.md](../../.routing/DESIGN.md): Monochrome design tokens (`main`, `side`, `dull`, `light`, `ground`), typography (`Satoshi`), surfaces, and buttons (`button-main`, `button-side`).
- [.routing/services/studio_service.md](../../.routing/services/studio_service.md): Canonical domain rules (`RN03` Ownership, `RN04` Resource Privacy, `RN05` Official Content), DTO contracts (`models.Collection`, `models.Deck`), and endpoints (`GET /api/v1/collections`).
- [app/frontend/uno.config.ts](../../app/frontend/uno.config.ts): Shortcuts (`header`, `button-icon`, `button-main`, `button-side`, `field`, `binding`, `dialog-overlay`).
- Visual references:
  - Desktop layout: [Desktop Reference Image](../../../.t3/userdata/attachments/627442ab-1232-401c-8e5e-25a39ad3e3e9-251d2801-12b1-485e-9a67-4c85bbc48f6b.png)
  - Mobile layout: [Mobile Reference Image](../../../.t3/userdata/attachments/627442ab-1232-401c-8e5e-25a39ad3e3e9-8e9353ec-cd92-43ab-ba83-a08c7f047a40.png)

## Deviations from Visual References

1. **No Top-Right Search Bar or Filter Buttons:**
   Universal search is globally available via `SearchModal` (`Ctrl + K` or TopBar search icon). The search input and top filter buttons from the reference header are omitted.
2. **Natural Sentence-Case Section Labels:**
   No forced all-caps (`uppercase tracking-widest`). Metadata and section headers use the `.section-label` pattern (`text-sm font-bold text-dull`).
3. **Monochrome Design System & Zero Regional Codes:**
   Strictly NO `flag_emoji` and NO regional codes (`CU`, `CL`, `MX`). Render full textual region names. Use monochrome borders and surfaces (`bg-main`, `bg-ground`, `border-light`).

## Task

1. **Type Augmentation & Router Setup:**
   - In `app/frontend/src/vite-env.d.ts` (or `types/router.d.ts`): Augment `RouteMeta` to support `breadcrumb?: Array<{ label: string; to?: string }>`.
   - In `app/frontend/src/router/index.ts`:
     - Register `/collections` pointing to `@/views/CollectionsView.vue`.
     - Set `meta: { breadcrumb: [{ label: "Home", to: "/" }, { label: "Collections" }] }`.

2. **Universal Breadcrumb (`AppBreadcrumb.vue`):**
   - Create `app/frontend/src/components/ui/AppBreadcrumb.vue`:
     - Reads `route.meta.breadcrumb`.
     - Displays `Home / Collections` (clickable ancestors, bold non-clickable leaf).
     - Renders within `max-w-6xl mx-auto px-4 py-2.5 flex items-center gap-1.5 text-xs text-dull`.
   - Wire `<AppBreadcrumb />` in `app/frontend/src/layouts/AppLayout.vue` directly below `<TopBar />` (hidden if `meta.breadcrumb` is empty or path is `/`).

3. **API Client Expansion (`IStudioAPI`):**
   - In `app/frontend/src/api/contract/studio.contract.ts`:
     - Add `listCollections(featured?: boolean): Promise<Collection[]>`.
     - Import `Collection` from `@bindings/palabrario/pkg/models/index.js`.
   - In `app/frontend/src/api/live/wails/studio.wails.ts`:
     - Implement `listCollections(featured = false)` delegating to `Wails.ListCollections(featured)`.
   - In `app/frontend/src/api/mock/mock.adapter.ts` and `decks.fixture.ts`:
     - Implement `listCollections(featured = false)` returning mock collection objects.
     - Ensure mock fixtures contain rich, authentic Spanish metadata (constituent decks, titles, descriptions).

4. **Collections View (`CollectionsView.vue`):**
   - **Data Fetching via API Client:**
     - Import `api` from `@/api`:
       ```ts
       import { api } from "@/api"
       ```
     - Fetch collections on mount:
       ```ts
       const collections = ref<Collection[]>([])
       const isLoading = ref(true)

       onMounted(async () => {
         try {
           collections.value = await api.studio.listCollections()
         } finally {
           isLoading.value = false
         }
       })
       ```
     - **Strict Rule:** Never hardcode mock collections or mock deck arrays inside `CollectionsView.vue` or child components. All data must flow through `api.studio`.
   - **Header:**
     - Big Title: `Study Collections` in `text-4xl sm:text-5xl font-black text-main tracking-tight`.
     - Subhead: *"Thematic vocabulary journeys and study paths exploring Spanish through everyday life."* in `text-base text-side mt-2 mb-8 sm:mb-10`.
   - **Section 1: Featured Spotlight (`FeaturedCollectionCard.vue`):**
     - Receives `collection: Collection` as prop (e.g. `collections.find(c => c.is_featured) ?? collections[0]`).
     - Two columns on desktop (`grid md:grid-cols-12 md:gap-8`), stacked on mobile.
     - Sharp rectangular container: `border border-light bg-white p-6 sm:p-8`.
     - **Left column (Hero):**
       - Eyebrow: *"Featured collection · Official"* in `text-sm font-bold text-dull mb-2`.
       - Title: `collection.title` in `text-2xl sm:text-3xl font-black text-main leading-tight`.
       - Narrative description: `collection.description` in `text-sm sm:text-base text-side leading-relaxed my-4`.
       - Action: `<button class="button-main">Explore Collection →</button>`.
     - **Right column (Collection Paths / Decks):**
       - Label: *"Collection paths"* in `text-sm font-bold text-dull border-b border-light pb-2 mb-4`.
       - Ordered list of constituent decks from `collection.decks` with short blurbs, separated by subtle dividers.
   - **Section 2: More Collections Catalog Grid:**
     - Section Title: *"More collections"* in `text-xl font-bold text-main mt-12 mb-6`.
     - 2-column grid on desktop (`md:grid md:grid-cols-2 md:gap-x-12 md:gap-y-8`), 1-column on mobile.
     - Render `collections.filter(c => c.id_collection !== featuredCollection?.id_collection)` with `CollectionGridItem.vue`:
       - Header row: Status tag (*"Official"* or *"Curated"*) on left in `text-xs font-bold text-dull`, `"Explore →"` on right in `text-xs font-semibold text-main hover:underline cursor-pointer`.
       - Title in `text-xl sm:text-2xl font-bold text-main mt-1 mb-2`.
       - Description narrative in `text-sm text-side leading-relaxed mb-4`.
       - Hairline border bottom `border-b border-light pb-6`.
   - **Section 3: Collection Detail Inspection Dialog (`CollectionDetailDialog.vue`):**
     - Opens when clicking *"Explore Collection →"* or *"Explore →"*.
     - Built using `AppDialog.vue` (Reka UI).
     - Displays full sequence of decks, item counts, and action to review or save.

5. **Validation:**
   - Run `cd app && wails3 task check` to verify TypeScript and Vue compilation.

5. **Change log:**
   - Document all changes in `.tickets/010-build-collections-route/changelog.md`.

## Verification

All checks are listed in [spec.md](spec.md).
