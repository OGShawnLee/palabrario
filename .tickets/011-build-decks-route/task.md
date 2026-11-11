# 011 — Build the Public Decks Route (/decks)

| Field | Value |
| :--- | :--- |
| **ID** | 011 |
| **Status** | completed |
| **Model** | Claude Sonnet 5.5, `effort: medium` (`claude-sonnet-5-5`) |
| **Role** | Frontend Implementer |
| **Spec** | [spec.md](spec.md) |

> **Specification Note:**  
> This ticket builds the public discovery route `/decks` following the approved desktop and mobile editorial wireframes. The design is directly aligned with the collections route (`/collections`), featuring a universal breadcrumb (`Home / Decks`), a 2-column Featured Deck spotlight card with constituent term previews, and a 2-column "More Decks" catalog grid with clean hairline dividers.

## Scope

May create or edit:
- `app/frontend/src/views/DecksView.vue` (main public decks catalog)
- `app/frontend/src/components/decks/FeaturedDeckCard.vue` (spotlight card with term preview column)
- `app/frontend/src/components/decks/DeckGridItem.vue` (catalog card for secondary decks)
- `app/frontend/src/api/contract/studio.contract.ts` (expand IStudioAPI with listDecks)
- `app/frontend/src/api/live/wails/studio.wails.ts` (implement listDecks)
- `app/frontend/src/api/mock/mock.adapter.ts` (implement mock listDecks)
- `app/frontend/src/api/mock/fixtures/decks.fixture.ts` (ensure mock decks contain items for term previews)
- `app/frontend/src/router/index.ts` (register `/decks` route and breadcrumb metadata)
- `.tickets/011-build-decks-route/changelog.md` (create)

Must not touch:
- `app/frontend/uno.config.ts` (shortcuts already exist; record proposals in changelog)
- Go backend code, `core-service/`, `realtime-service/`, or `.context/`
- Any files outside the above declared scope

## Inputs

- [.routing/API.md](../../.routing/API.md): Architectural rules (Rule 1: Never import `@bindings/.../app` in UI; always import `api` from `@/api`).
- [app/frontend/src/api/index.ts](../../app/frontend/src/api/index.ts): Universal `api` singleton (`api.studio`).
- [app/frontend/src/views/CollectionsView.vue](../../app/frontend/src/views/CollectionsView.vue): **Authoritative layout reference.** `DecksView.vue` must mirror its outer container, spacing, typography, and loading/empty patterns.
- [app/frontend/src/components/collections/FeaturedCollectionCard.vue](../../app/frontend/src/components/collections/FeaturedCollectionCard.vue): **Authoritative spotlight reference.** `FeaturedDeckCard.vue` must mirror its 7/5 12-column grid, padding, and hairline borders.
- [app/frontend/src/components/collections/CollectionGridItem.vue](../../app/frontend/src/components/collections/CollectionGridItem.vue): **Authoritative grid item reference.** `DeckGridItem.vue` must mirror its top metadata flex row, typography, and bottom border.
- [.routing/CONTEXT.md](../../.routing/CONTEXT.md): Architecture ground rules (strict prohibition of `flag_emoji` and regional codes; full textual region names required; English UI strings).
- [.routing/DESIGN.md](../../.routing/DESIGN.md): Monochrome design tokens (`main`, `side`, `dull`, `light`, `ground`), typography (`Satoshi`), surfaces, sharp geometry (`rounded-none`), buttons (`button-main`, `button-side`), and natural sentence-case section labels (`section-label`).
- [.routing/services/studio_service.md](../../.routing/services/studio_service.md): Canonical domain rules (`RN03` Ownership, `RN05` Official Content), DTO contracts (`models.Deck`, `models.DeckItem`), and endpoints (`GET /api/v1/decks`).
- Visual references:
  - Desktop layout: [Desktop Reference Image](../../../.t3/userdata/attachments/96471626-efc2-487e-be8c-60c037a8d0b8-ebe01332-5c9c-4d77-ac9e-02904498cf3a.png)
  - Mobile layout: [Mobile Reference Image](../../../.t3/userdata/attachments/96471626-efc2-487e-be8c-60c037a8d0b8-ef9ef017-8dad-4215-8246-0e423b97012c.png)

## Deviations from Visual References

1. **No Top-Right Search Bar or Filter Segmented Control:**
   Universal search is globally available via `SearchModal` (`Ctrl + K` or TopBar search icon). Mirroring `/collections`, the top-right search box and filter buttons from the reference header are omitted.
2. **English UI Strings:**
   All labels and interactive copy are in English per `CONTEXT.md` (e.g. *"Study Decks"*, *"Featured Deck • Official • 5 terms"*, *"Terms in this deck"*, *"Preview • 3 of 5 terms"*, *"Study Deck →"*, *"Save to Studio"*, *"More Decks"*, *"Explore →"*).
3. **Monochrome Design System & Zero Regional Codes:**
   Strictly NO `flag_emoji` and NO regional codes (`CU`, `CL`, `MX`). Render full textual region names (*"Mexico"*, *"Chile"*). Sharp rectangular geometry (`rounded-none`, no drop shadows).

## Task

1. **API Client Expansion (`IStudioAPI`):**
   - In `app/frontend/src/api/contract/studio.contract.ts`:
     - Add `listDecks(q?: string, officialOnly?: boolean): Promise<Deck[]>`.
     - Import `Deck` from `@bindings/palabrario/pkg/models/index.js`.
   - In `app/frontend/src/api/live/wails/studio.wails.ts`:
     - Delegate to `Wails.ListDecks(q ?? "", officialOnly ?? false)`.
   - In `app/frontend/src/api/mock/mock.adapter.ts` and `decks.fixture.ts`:
     - Implement mock `listDecks(...)` returning filtered decks.
     - Ensure mock decks have preview `items` populated with authentic lemmas, full `region_name` strings, and definitions.

2. **Router Setup (`app/frontend/src/router/index.ts`):**
   - Register `/decks` under `AppLayout` with `name: "decks"`.
   - Set up `meta: { breadcrumb: [{ label: "Home", to: "/" }, { label: "Decks" }] }`.

3. **Public Decks View Container (`DecksView.vue`):**
   - Centered container layout (`max-w-6xl mx-auto px-4 pt-10 pb-20 flex-1 w-full`).
   - Title: `Study Decks` in `text-4xl sm:text-5xl font-black text-main tracking-tight`.
   - Subhead: *"Vocabulary decks created by the editorial team and community to practice or take into the Arena."* in `text-base text-side mt-2 mb-8 sm:mb-10`.
   - **Data Fetching via API Client:**
     - Import `api` from `@/api`.
     - Fetch decks on mount via `api.studio.listDecks()`.
     - Reactive states: `isLoading`, empty state (*"No decks are available right now."*).
     - Featured deck computation: `decks.find(d => d.is_official) ?? decks[0] ?? null`.
     - Secondary decks computation: `decks.filter(d => d.id_deck !== featured?.id_deck)`.
     - **Strict Rule:** Never hardcode mock deck arrays in views or components. All data must flow through `api.studio`.

4. **Featured Deck Spotlight Card (`FeaturedDeckCard.vue`):**
   - Two columns on desktop (`grid md:grid-cols-12 md:gap-8`), single column on mobile.
   - Bordered container: `border border-light bg-white p-6 sm:p-8`.
   - **Left Column (Deck Details & Primary CTAs):**
     - Eyebrow: `Featured Deck · [Official / Community] · [N] terms` in `section-label !mb-3`.
     - Deck Title: `text-2xl sm:text-3xl font-black text-main leading-tight mb-3`.
     - Narrative Description: `text-sm sm:text-base text-side leading-relaxed mb-6`.
     - Primary CTAs (horizontal row on desktop, stacked on mobile):
       - `<button class="button-main cursor-pointer" @click="router.push('/studio')">Study Deck →</button>`
       - `<button class="button-side cursor-pointer">Save to Studio</button>`
   - **Right Column (Term Preview Column):**
     - Section label: `"Terms in this deck"` (`section-label`).
     - Subhead: `"Preview • 3 of [N] terms"` in `text-xs text-dull mb-4`.
     - List of up to 3 preview items from `deck.items`:
       - Lemma in `text-sm font-bold text-main`.
       - Region + Definition: `[Region Name] · [Definition]` in `text-xs text-side leading-relaxed mt-0.5`.
       - Hairline bottom divider `border-b border-light pb-3 mb-3 last:border-b-0 last:pb-0 last:mb-0`.

5. **Secondary Decks Catalog Grid (`DeckGridItem.vue`):**
   - Section heading in view: `"More Decks"` (`text-xl font-bold text-main mt-12 mb-6`).
   - 2-column grid on desktop (`grid grid-cols-1 gap-y-8 md:grid-cols-2 md:gap-x-12`), 1-column on mobile.
   - Component layout (`border-b border-light pb-6`):
     - Top metadata row: Status + term count tag (`[Official / Community] · [N] terms`) in `text-xs font-bold text-dull`, and `"Explore →"` button on right in `text-xs font-semibold text-main hover:underline cursor-pointer` (navigates to `/studio`).
     - Title in `text-xl sm:text-2xl font-bold text-main mt-1 mb-2`.
     - Description blurb in `text-sm text-side leading-relaxed`.

6. **Validation:**
   - Run `cd app && wails3 task check` to verify TypeScript and Vue compilation with 0 errors.

7. **Change log:**
   - Document all changes in `.tickets/011-build-decks-route/changelog.md`.

## Verification

All checks are listed in [spec.md](spec.md).
