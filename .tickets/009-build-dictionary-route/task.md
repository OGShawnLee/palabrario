# 009 — Build the Dictionary Route (/dictionary)

| Field | Value |
| :--- | :--- |
| **ID** | 009 |
| **Status** | completed |
| **Model** | Gemini 3.8 Flash (`gemini-3.8-flash`, Antigravity provider) |
| **Role** | Frontend Implementer |
| **Spec** | [spec.md](spec.md) |

## Scope

May create or edit:
- `app/frontend/src/composables/useSearchModal.ts` (create: shared search modal state and controls)
- `app/frontend/src/layouts/AppLayout.vue` (wire `useSearchModal` to `SearchModal` and `TopBar`)
- `app/frontend/src/data/mockWordOfTheDay.ts` (create: typed mock data matching `WordDetail` contract)
- `app/frontend/src/views/DictionaryView.vue` (implement view)
- `app/frontend/src/components/dictionary/` (optional subcomponents: e.g. `WordOfTheDayCard.vue`)
- `.tickets/009-build-dictionary-route/changelog.md` (create)

Must not touch:
- `app/frontend/uno.config.ts` (propose shortcuts or tokens in changelog)
- `app/*.go`, backend services, or `.context/`
- Any file outside the above scope

## Inputs

- [.routing/CONTEXT.md](../../.routing/CONTEXT.md): Architecture ground rules (strict prohibition of `flag_emoji` and regional codes; full region names required).
- [.routing/DESIGN.md](../../.routing/DESIGN.md): Monochrome design tokens (`main`, `side`, `dull`, `light`, `ground`), typography (`Satoshi`), surfaces, and buttons (`button-main`, `button-side`, `binding`).
- [.routing/services/dictionary_service.md](../../.routing/services/dictionary_service.md): Pan-Hispanic Lexicographical Core, Senses, Heteronyms, and `WordDetail` / `SenseDetail` / `Example` / `HeteronymItem` contracts.
- [app/frontend/src/data/regions.ts](../../app/frontend/src/data/regions.ts): Canonical region definitions.
- [app/frontend/src/layouts/AppLayout.vue](../../app/frontend/src/layouts/AppLayout.vue): Current layout wrapping routes.
- Visual references:
  - Big search hero bar: [Image 2: Vocabulary.com style big input bar](../../../.t3/userdata/attachments/e7510c5f-4aee-4fe2-8b04-76e954f1ea52-23a0d289-1b35-4287-a0be-aa43fa19d98c.png)
  - Word of the Day desktop layout: [Image 1: Desktop editorial layout](../../../.t3/userdata/attachments/e7510c5f-4aee-4fe2-8b04-76e954f1ea52-d217cff8-84b2-40bb-811b-d6861b1dc8bd.png)
  - Word of the Day mobile layout: [Image 3: Mobile single-column layout](../../../.t3/userdata/attachments/e7510c5f-4aee-4fe2-8b04-76e954f1ea52-4ebd0a93-b7f0-49e1-bd49-c5e98cbf627b.png)

## Deviations from Visual References

1. **Search bar background & buttons:** No dark/patterned background container; clean minimal surface. No `Random Word` button.
2. **Flags and regional codes:** Strictly NO `flag_emoji` and NO dialect codes (e.g., `CU`, `CL`, `MX`). Render only full textual region names (e.g. `Cuba`, `Chile`, `Mexico`).
3. **Monochrome palette:** All colorful badges and highlights (e.g. vermilion/red badges) are adapted to the strict monochrome design system (`bg-ground`, `bg-light`, `border-light`, `text-main`, `text-side`, `text-dull`).

## Task

1. **Shared Modal State: `app/frontend/src/composables/useSearchModal.ts`:**
   - Create a module-scoped singleton composable with the following exact structure:
     ```ts
     import { ref } from "vue"

     const isOpen = ref(false)

     export function useSearchModal() {
       function open() {
         isOpen.value = true
       }
       function close() {
         isOpen.value = false
       }
       function toggle() {
         isOpen.value = !isOpen.value
       }

       return {
         isOpen,
         open,
         close,
         toggle,
       }
     }
     ```
   - In `app/frontend/src/layouts/AppLayout.vue`:
     - Replace local `const isSearchOpen = ref(false)` with `const { isOpen, open } = useSearchModal()`.
     - Bind `<SearchModal v-model:open="isOpen" />`.
     - Wire `<TopBar @open-search="open" />`.

2. **Section 1: Hero Search Bar (`DictionaryView.vue`):**
   - Centered headline above input: *"Look up a word, learn it forever."* in `text-2xl sm:text-3xl font-bold text-main tracking-tight mb-4 text-center`.
   - Big mock search trigger bar:
     - Implement as an interactive element (`<button type="button" role="combobox" aria-label="Search words or vocabulary" ...>`) styled as a prominent search input:
       `w-full max-w-3xl mx-auto h-14 sm:h-16 bg-white border border-light flex items-center px-4 sm:px-6 cursor-pointer focus:outline-none focus:border-main`
     - Left: `Search` icon from `lucide-vue-next` (`w-5 h-5 text-dull mr-3 sm:mr-4 shrink-0`).
     - Center text: *"Search words or vocabulary lists from books, exams or textbooks"* in `text-dull text-sm sm:text-base text-left flex-1 truncate`.
     - Right: Key badge `<span class="binding">Ctrl + K</span>`.
     - **Interaction:** Trigger `open()` from `useSearchModal()` on `@click` and on `@keydown.enter` / `@keydown.space`.

3. **Section 2: Word of the Day (Palabra del Día):**
   - **Typed Mock Data (`app/frontend/src/data/mockWordOfTheDay.ts`):**
     - Define and export `mockWordOfTheDay` conforming strictly to the canonical `WordDetail` contract in `.routing/services/dictionary_service.md`:
       - `id_word: 101`
       - `lemma: "Guagua"`
       - `part_of_speech: "noun · feminine"`
       - `gender: "FEMALE"`
       - `etymology: "From Quechua wawa (infant, child) in the Andean zone; in the Caribbean, likely onomatopoeic or derived from the transport company Wa & Wa Co. Inc."`
       - `friendly_blurb: "One of the most fascinating linguistic chameleons of the Spanish language. In the Caribbean, a guagua is a noisy public bus that gets you to work; in Chile, Ecuador, and Peru, a guagua is a tender newborn baby. Say it in the wrong country and you might offer to hold someone's bus!"`
       - `senses: SenseDetail[]`:
         - Sense 01: `region_code: "CU"`, `region_name: "Cuba"`, `cefr_level: "B1"`, `usage_frequency: "EVERYDAY"`, `short_definition: "Public transit bus"`, `definition: "Omnibús de servicio urbano o interurbano destinado al transporte masivo de pasajeros."`, example sentence: *"Ayer la guagua pasó llena y tuve que caminar."* with source *"Urban speech · Cuba"*.
         - Sense 02: `region_code: "CL"`, `region_name: "Chile"`, `cefr_level: "A2"`, `usage_frequency: "EVERYDAY"`, `short_definition: "Baby or infant"`, `definition: "Niño de pecho o criatura de muy corta edad."`, example sentence: *"La guagua se quedó dormida en brazos de su mamá."* with source *"Everyday speech · Chile"*.
       - `heteronyms: HeteronymItem[]`:
         - Public transit equivalents: `Mexico · Camión`, `Argentina · Colectivo`, `Chile · Micro`, `Spain · Autobús`.
         - Infant equivalents: `Chile · Guagua`, `Spain · Bebé`, `Mexico · Bebé`, `Argentina · Bebé`.
   - **Header & Title:**
     - Top metadata row: Eyebrows *"Word of the Day"* (left) and *"Regional Lexicon"* (right) in `text-xs font-bold uppercase tracking-widest text-dull`.
     - Huge Lemma headline: `Guagua` in `text-5xl sm:text-6xl md:text-7xl font-black text-main leading-none mt-2`.
     - Part of speech: `noun · feminine` in `text-sm text-dull mt-1 mb-6`.
   - **Editorial Blurb & Etymology:**
     - Desktop (`md:grid md:grid-cols-12 md:gap-8`):
       - Left (`col-span-4`): `ETYMOLOGY` in `text-xs font-bold uppercase tracking-widest text-dull mb-2` followed by the derivation text in `text-sm text-side leading-relaxed`.
       - Right (`col-span-8`): Friendly blurb narrative paragraph in `text-base sm:text-lg text-side leading-relaxed`.
     - Mobile (`< md`): Stacked vertically matching visual reference Image 3.
   - **Dialectal Senses Grid:**
     - 2 columns on desktop (`md:grid md:grid-cols-2 md:gap-8`), stacked single-column on mobile.
     - Each sense card renders:
       - Badges row: Full region name tag (e.g. `Cuba`, `Chile` in `border border-light bg-ground px-2.5 py-0.5 text-xs font-medium text-side`), CEFR level badge (e.g. `B1`, `A2` in `bg-main text-white text-xs font-bold px-2 py-0.5`).
       - Sense numbering and short definition: e.g. `01 Public transit bus.` in `text-lg sm:text-xl font-bold text-main mt-3`.
       - Sociolinguistic register / frequency: e.g. `Daily colloquial use` in `text-xs text-dull mb-2`.
       - Detailed definition & cultural note: narrative context in `text-sm text-side leading-relaxed mb-4`.
       - Contextual Example: Blockquote in italics (e.g. *«Ayer la guagua pasó llena y tuve que caminar.»*) with left solid border (`border-l-2 border-main pl-3`), followed by source attribution (e.g. `— Urban speech · Cuba` in `text-xs text-dull`).
   - **Cross-Regional Equivalents (Heteronyms):**
     - Section label: `Equivalencies in other countries` in `text-sm font-bold text-main mb-3`.
     - Categorized or grouped pills (e.g. For transportation: `Mexico · Camión`, `Argentina · Colectivo`; For infancy: `Chile · Guagua`, `Spain · Bebé`, `Mexico · Bebé`).
     - Rendered with full region names and terms inside `border border-light bg-white px-3 py-1 text-xs text-side`.
   - **Action Footer:**
     - "Add to Studio" button: `<button class="button-main">Add to Studio</button>`.
     - "See full entry →" button/link: `<RouterLink to="/dictionary/word/101" class="text-sm font-semibold text-main inline-flex items-center gap-1">See full entry →</RouterLink>`.

4. **Validation:**
   - Run `wails3 task check` from `app/`. It must exit with code 0.

5. **Change log:**
   - Document all changes in `.tickets/009-build-dictionary-route/changelog.md`.

## Verification

All checks are listed in [spec.md](spec.md).
