# 021 — Build the Individual Collection Detail Route (/collections/:id)

| Field | Value |
| :--- | :--- |
| **ID** | 021 |
| **Status** | draft |
| **Model** | Sonnet 5.5 (`claude-sonnet-5-5`, Antigravity provider) |
| **Role** | Frontend Implementer |
| **Spec** | [spec.md](spec.md) |

> **Specification Note:**  
> This ticket builds the dedicated route `/collections/:id` for exploring an individual learning path, its constituent study decks, and role-based management controls following the approved desktop and mobile editorial wireframes. The route integrates with universal `AppBreadcrumb.vue` (via route meta/props), features primary actions ("Save to Studio" and "Study First Deck →"), sequential deck cards, and a conditional owner/administrator management bar at the bottom.

## Scope

May create or edit:
- `app/frontend/src/views/CollectionDetailView.vue` (main collection detail view)
- `app/frontend/src/components/collections/CollectionHeader.vue` (headline, blurb, metadata, and primary CTAs)
- `app/frontend/src/components/collections/CollectionDeckItem.vue` (individual deck card in the sequential path)
- `app/frontend/src/components/collections/CollectionAdminBar.vue` (role-gated owner/administrator actions bar)
- `app/frontend/src/api/contract/studio.contract.ts` (ensure getCollection, saveCollection, unsaveCollection exist)
- `app/frontend/src/api/live/wails/studio.wails.ts` (implement getCollection, saveCollection, unsaveCollection)
- `app/frontend/src/api/mock/mock.adapter.ts` (implement mock getCollection, saveCollection, unsaveCollection)
- `app/frontend/src/router/index.ts` (register `/collections/:id` route and breadcrumb resolver)
- `.tickets/021-build-collection-detail-route/changelog.md` (create)

Must not touch:
- `app/frontend/uno.config.ts` (shortcuts already exist; record proposals in changelog)
- Go backend code, `core-service/`, `realtime-service/`, `image-service/`, or `.context/`
- Any files outside the above declared scope

## Inputs

- [.routing/API.md](../../.routing/API.md): Architectural rules (Rule 1: Never import `@bindings/.../app` in UI; always import `api` from `@/api`).
- [app/frontend/src/api/index.ts](../../app/frontend/src/api/index.ts): Universal `api` singleton (`api.studio`, `api.auth`).
- [.routing/CONTEXT.md](../../.routing/CONTEXT.md): Architecture ground rules (strict prohibition of `flag_emoji` and regional codes; full textual region names required; English UI strings).
- [.routing/DESIGN.md](../../.routing/DESIGN.md): Monochrome design tokens (`main`, `side`, `dull`, `light`, `ground`), typography (`Satoshi`), surfaces, sharp geometry (`rounded-none`), buttons (`button-main`, `button-side`), and natural sentence-case section labels (`section-label`).
- [.routing/services/studio_service.md](../../.routing/services/studio_service.md): Canonical domain rules (`RN03` Ownership, `RN04` Resource Privacy, `RN05` Official Content), DTO contracts (`models.Collection`, `models.Deck`), and endpoints (`GET /api/v1/collections/{id}`).
- Visual references:
  - Desktop layout: [Desktop Reference Image](../../../.t3/userdata/attachments/627442ab-1232-401c-8e5e-25a39ad3e3e9-db2340cb-d834-4173-8ad5-c10c42bd3bd9.png)
  - Mobile layout: [Mobile Reference Image](../../../.t3/userdata/attachments/627442ab-1232-401c-8e5e-25a39ad3e3e9-c40f9136-0a01-43fe-8772-c422e5cb670a.png)

## Deviations from Visual References

1. **English UI Strings:**
   All labels and interactive copy are in English per `CONTEXT.md` (e.g. *"Featured collection • Official curriculum"*, *"Save to Studio"*, *"Study First Deck →"*, *"Decks in this collection (3)"*, *"Explore Deck →"*, *"Alternative view • Owner / Administrator"*, *"Edit Collection"*, *"Add Deck"*, *"Delete"*).
2. **Dynamic Breadcrumb Resolution:**
   Top breadcrumb displays `Home / Collections / [Collection Title]` via `AppBreadcrumb.vue` mounted in `AppLayout.vue`.
3. **Monochrome Design System & Zero Regional Codes:**
   Strictly NO `flag_emoji` and NO regional codes (`CU`, `CL`, `MX`). Render full textual region names only. Sharp rectangular geometry (`rounded-none`, no drop shadows).

## Task

1. **Router Registration (`app/frontend/src/router/index.ts`):**
   - Register route `/collections/:id` with name `collection-detail` under `AppLayout`.
   - Set up `meta: { breadcrumb: [{ label: "Home", to: "/" }, { label: "Collections", to: "/collections" }, { label: "Collection Detail" }] }`.

2. **API Client Support for Collections:**
   - In `app/frontend/src/api/contract/studio.contract.ts`:
     - Add `getCollection(id: number): Promise<Collection | null>`.
     - Add `saveCollection(id: number): Promise<void>`.
     - Add `unsaveCollection(id: number): Promise<void>`.
   - In `app/frontend/src/api/live/wails/studio.wails.ts`:
     - Delegate to `Wails.GetCollection(id)`, `Wails.SaveCollection(id)`, `Wails.UnsaveCollection(id)`.
   - In `app/frontend/src/api/mock/mock.adapter.ts` and `decks.fixture.ts`:
     - Implement mock `getCollection(id)` returning full collection with constituent decks.
     - Implement mock `saveCollection(id)` and `unsaveCollection(id)` updating session state.

3. **Collection Detail View Container (`CollectionDetailView.vue`):**
   - Centered container layout (`max-w-4xl mx-auto px-4 py-8 sm:py-12 flex flex-col gap-8`).
   - **Data Fetching via API Client:**
     - Import `api` from `@/api`:
       ```ts
       import { api } from "@/api"
       ```
     - On mount or when `route.params.id` changes, fetch collection and user session:
       ```ts
       const collection = ref<Collection | null>(null)
       const currentUser = ref<UserProfile | null>(null)
       const isLoading = ref(true)

       onMounted(async () => {
         try {
           const id = Number(route.params.id)
           const [col, user] = await Promise.all([
             api.studio.getCollection(id),
             api.auth.getMe(),
           ])
           collection.value = col
           currentUser.value = user
         } finally {
           isLoading.value = false
         }
       })
       ```
     - Determine management permissions reactively:
       `isOwnerOrAdmin = computed(() => currentUser.value?.role === AccountRole.RoleAdmin || currentUser.value?.user_name === collection.value?.owner_user_name)`
     - **Strict Rule:** Never hardcode mock collection data in `CollectionDetailView.vue`. All data must be fetched through `api.studio`. Show loading state when `isLoading` is true.

3. **Collection Header Component (`CollectionHeader.vue`):**
   - Eyebrow tag: *"Featured collection • Official curriculum"* in `section-label` (`text-sm font-bold text-dull mb-2`).
   - Large Title: `text-4xl sm:text-5xl font-black text-main leading-tight tracking-tight`.
   - Description blurb: `text-base sm:text-lg text-side leading-relaxed max-w-3xl my-4`.
   - Metadata line: *"Curated by Academia Mexicana • 3 decks • 36 terms • Updated 2 days ago"* in `text-xs sm:text-sm text-dull font-medium mb-6`.
   - Primary Call-to-Actions:
     - On Desktop: Horizontal flex row (`flex items-center gap-3`).
     - On Mobile: Full-width stacked buttons (`flex flex-col gap-3 w-full`).
     - Button 1: `<button class="button-main flex items-center justify-center gap-2">` with Lucide `Bookmark` icon and text `"Save to Studio"` (toggles to `"✓ Saved in Studio"`).
     - Button 2: `<button class="button-side flex items-center justify-center gap-2">` with text `"Study First Deck →"`.

4. **Sequential Deck List & Cards (`CollectionDeckItem.vue`):**
   - Section heading: `"Decks in this collection (3)"` in `section-label` (`text-sm font-bold text-dull mb-4`).
   - Vertical stack of rectangular cards (`border border-light bg-white p-6 rounded-none flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4`):
     - Left content:
       - Term count tag: `"12 terms"` in `text-xs text-dull font-medium mb-1`.
       - Deck Title: `text-xl font-bold text-main mb-1`.
       - Description: `text-sm text-side leading-relaxed`.
     - Right action:
       - `<button class="button-side w-full sm:w-auto inline-flex items-center justify-center gap-1.5 cursor-pointer">Explore Deck →</button>`.

5. **Role-Gated Management Bar (`CollectionAdminBar.vue`):**
   - Separated by top hairline divider (`border-t border-light pt-8 mt-8 flex flex-col gap-4`).
   - Eyebrow: `"Alternative view • Owner / Administrator"` in `text-xs font-semibold text-dull`.
   - Action buttons row (`flex flex-col sm:flex-row items-stretch sm:items-center gap-3`):
     - `<button class="button-side inline-flex items-center justify-center gap-2"><Edit3 class="w-4 h-4" /> Edit Collection</button>`
     - `<button class="button-side inline-flex items-center justify-center gap-2"><Plus class="w-4 h-4" /> Add Deck</button>`
     - `<button class="text-sm font-semibold text-dull hover:text-black py-2 px-3 text-center cursor-pointer">Delete</button>`
   - On mobile: Buttons stack vertically with min-44px touch height, matching the mobile reference.

6. **Validation:**
   - Run `cd app && wails3 task check` to verify TypeScript and Vue compilation with 0 errors.

7. **Change log:**
   - Document all changes in `.tickets/021-build-collection-detail-route/changelog.md`.

## Verification

All checks are listed in [spec.md](spec.md).
