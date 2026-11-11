# 023 — Add Collections Editorial Discovery Section to DictionaryView

| Field | Value |
| :--- | :--- |
| **ID** | 023 |
| **Status** | completed |
| **Model** | Claude Sonnet 5.5, `effort: medium` (`claude-sonnet-5-5`) |
| **Role** | Frontend Implementer |
| **Spec** | [spec.md](spec.md) |

> **Specification Note:**  
> This ticket builds a broadsheet cultural discovery section (`DictionaryCollectionsSection.vue`) placed directly below the Word of the Day section in `DictionaryView.vue`. The layout draws from the approved editorial designs, featuring monumental typography, a localized black stat block, hairline deck lists with `↗` indicators, and direct navigation to `/collections` and `/collections/:id`.

## Visual Design References

- **Desktop Layout:** `C:\Users\OGSha\.t3\userdata\attachments\96471626-efc2-487e-be8c-60c037a8d0b8-a6246336-be45-4abe-8290-553edbcbc216.png`
- **Mobile Layout:** `C:\Users\OGSha\.t3\userdata\attachments\96471626-efc2-487e-be8c-60c037a8d0b8-63c04bd5-7608-440f-952e-fcd767832938.png`

## Scope

May create or edit:
- `app/frontend/src/views/DictionaryView.vue` (mount `<DictionaryCollectionsSection />` below Word of the Day)
- `app/frontend/src/components/dictionary/DictionaryCollectionsSection.vue` (create the editorial showcase)
- `.tickets/023-add-collections-section-to-dictionary/changelog.md` (create)

Must not touch:
- `app/frontend/uno.config.ts`
- Go backend code, `core-service/`, `realtime-service/`, or `.context/`
- Any files outside the above declared scope

## Inputs & Strict Ground Rules

- [.routing/API.md](../../.routing/API.md): Architectural rules (import `api` from `@/api`, never import `@bindings/.../app`).
- [app/frontend/src/api/index.ts](../../app/frontend/src/api/index.ts): Universal `api` singleton (`api.studio.listCollections(featured = true)`).
- [app/frontend/src/views/DictionaryView.vue](../../app/frontend/src/views/DictionaryView.vue): Target view.
- [.routing/DESIGN.md](../../.routing/DESIGN.md):
  - **Typography:** ONLY Satoshi font family (`font-sans`). **`font-mono` is strictly forbidden** across all numbers, indexes, and tags. Use standard tabular numerals (`tabular-nums`) or regular Satoshi font weights (`font-bold`, `font-medium`, `font-light`).
  - **Language Division:** **UI labels, buttons, navigation links, and chrome are in English** (e.g., `Thematic Paths`, `Explore collection →`, `View all collections →`, `terms`, `decks`, `Official Path`). The curated cultural content (titles, descriptions, editorial quotes like *"Menos alfabeto. Más mundo."* and *"Una palabra abre una puerta. Una colección, un mundo."*) is in Spanish.
  - **Monochrome Tokens:** `main` (black `#000000`), `side` (`#404040`), `dull` (`#737373`), `light` (`#F5F5F5`), `white` (`#FFFFFF`), `ground` (`#FAFAFA`).
  - **Sharp Geometry:** `rounded-none`, no drop shadows, 1px solid hairline borders (`border-light`).
- [.routing/CONTEXT.md](../../.routing/CONTEXT.md): Strictly zero `flag_emoji` and zero regional abbreviations/codes (`CU`, `CL`, `MX`). Full country names only.

## Technical Boundaries & Data Grounding

1. **No Client-Side Deep Hydration for Region Tags:**
   `Collection` and `Deck` summary models do not contain a pre-aggregated `regions` array, and fetching all constituent senses on the client in a summary showcase is strictly forbidden (unviable N+1 network cost). Do not hallucinate or ad-hoc parse regions from senses. The thematic focus is already conveyed editorially in collection and deck titles.
2. **Real Fields & Metrics:**
   - Term counts are dynamically summed: `col.decks?.reduce((sum, d) => sum + d.item_count, 0) ?? 0`.
   - Deck count is `col.decks?.length ?? 0`.
   - Plan type is `col.is_official ? 'Official Path' : 'Community Path'`.
3. **API Integration:**
   - Call `api.studio.listCollections(true)` (with fallback to all collections if none featured).
   - Zero hardcoded mock arrays in the UI component.

## Task

1. **Build `DictionaryCollectionsSection.vue`:**
   - **Header (Monumental Editorial Split):**
     - Top eyebrow row: `01 / Thematic Paths` (left) and `Language in context` (right) in `text-xs font-bold text-dull tracking-normal`.
     - 2-column responsive layout:
       - Left column: Monumental headline `Menos alfabeto.` `<br/>` `Más mundo.` in `text-4xl sm:text-6xl font-black text-main tracking-tight leading-none mt-4 mb-6 sm:mb-0`.
       - Right column: Rich editorial blurb in `text-sm sm:text-base text-side leading-relaxed max-w-md`:
         > *"El español no se aprende por orden alfabético. Se recorre entre sobremesas, mercados y esquinas.*
         > *Colecciones que enlazan barajas de estudio para entender una lengua desde la cultura que la hace vivir."*
     - Structural divider: `border-b border-light my-10`.

   - **Collection Spotlight 1 (Hero Showcase):**
     - 2-column layout (stacks vertically on mobile):
       - Left column (Localized Stat Block, ~220px on desktop):
         - Solid black container `bg-black text-white p-6 sm:p-8 rounded-none flex flex-col justify-between aspect-[3/4] sm:aspect-auto sm:min-h-[280px]`.
         - Top: `text-xs font-bold text-neutral-400 tracking-wider` with `Path 01 / ${col.is_official ? 'Official Path' : 'Community'}`.
         - Middle: Giant numeral `text-6xl sm:text-7xl font-bold tabular-nums leading-none` (total terms) and `text-xs font-medium text-neutral-400 mt-1` (`terms`).
         - Bottom: `text-lg font-bold text-white` (`${col.decks?.length} decks`) and `text-xs text-neutral-400` (`Shared learning journey`).
       - Right column (Content & Constituent Deck List):
         - Title: `text-2xl sm:text-4xl font-bold text-main tracking-tight mb-3` (`col.title`).
         - Description: `text-sm sm:text-base text-side leading-relaxed mb-6` (`col.description`).
         - Decks section label: `text-xs font-bold text-dull mb-2` (`Constituent study decks`).
         - Decks list (`divide-y divide-light border-y border-light mb-6`):
           - Each constituent deck row (`py-3 flex items-center justify-between group hover:bg-white px-2 cursor-pointer transition-colors`):
             - Left: `text-xs font-bold text-dull mr-4 tabular-nums` (`01`, `02`, `03`) + `text-sm sm:text-base font-bold text-main group-hover:underline` (`deck.title`).
             - Right: Arrow `↗` in `text-sm text-dull group-hover:text-main transition-colors`. Clicking row navigates to `/decks`.
         - Primary action: `"Explore collection →"` in `text-sm font-semibold text-main hover:underline flex items-center gap-1.5 cursor-pointer` navigating to `/collections/${col.id_collection}`.

   - **Collection Spotlight 2 (Secondary Asymmetric Split, if second collection exists):**
     - Structural divider `border-b border-light my-12`.
     - 3-column / split layout:
       - Column 1: Index `02 / Official Path`, subtitle blurb, and term summary (`tabular-nums`).
       - Column 2: Title (`text-2xl font-bold text-main`), description, and `"Explore collection →"`.
       - Column 3: Hairline list of constituent decks with term counts aligned to the right in `text-sm font-bold text-main tabular-nums`.

   - **Closing Manifesto Banner:**
     - `border-t border-light pt-8 mt-14 flex flex-col sm:flex-row justify-between items-start sm:items-end gap-6`.
     - Left: `text-xl sm:text-2xl font-bold text-main tracking-tight leading-snug`:
       *"Una palabra abre una puerta. Una colección, un mundo."*
     - Right:
       - Action: `"View all collections →"` in `text-sm font-bold text-main hover:underline flex items-center gap-2 cursor-pointer` navigating to `/collections`.
       - Sub-tag: `text-xs font-medium text-dull mt-1` (`Thematic Paths · Cultural Dictionary`).

   - **Empty / Loading States:**
     - Loading: Subtle skeleton placeholders matching the spatial layout.
     - Empty: If `listCollections()` returns 0 items, render a clean empty notice without broken UI.

2. **Mount in `DictionaryView.vue`:**
   - Mount `<DictionaryCollectionsSection />` below the Word of the Day section in `max-w-5xl mx-auto px-4 pb-20`.

3. **Validation:**
   - Run `cd app && wails3 task check` to verify TypeScript, Vite build, and Go build pass with 0 errors.

4. **Change log:**
   - Document changes in `.tickets/023-add-collections-section-to-dictionary/changelog.md`.

## Verification

All checks are listed in [spec.md](spec.md).
