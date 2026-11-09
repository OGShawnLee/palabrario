 009 — Build the Dictionary Route (/dictionary)

| Field | Value |
| :--- | :--- |
| **ID** | 009 |
| **Status** | draft |
| **Model** | Sonnet 5.5 (`claude-sonnet-5-5`, Antigravity provider) |
| **Role** | Frontend Implementer |
| **Spec** | [spec.md](spec.md) |

## Scope

May create or edit:
- `app/frontend/src/views/DictionaryView.vue` (implement view)
- `app/frontend/src/components/dictionary/` (optional subcomponents: e.g. `WordOfTheDayCard.vue`)
- `app/frontend/src/composables/useSearchModal.ts` (shared search modal trigger, optional)
- `app/frontend/src/layouts/AppLayout.vue` (wire search modal trigger if composable/event used)
- `.tickets/009-build-dictionary-route/changelog.md` (create)

Must not touch:
- `app/frontend/uno.config.ts` (propose shortcuts or tokens in changelog)
- `app/*.go`, backend services, or `.context/`
- Any file outside the above scope

## Inputs

- [.routing/CONTEXT.md](../../.routing/CONTEXT.md): Architecture ground rules (strict prohibition of `flag_emoji` and regional codes; full region names required).
- [.routing/DESIGN.md](../../.routing/DESIGN.md): Monochrome design tokens (`main`, `side`, `dull`, `light`, `ground`), typography (`Satoshi`), surfaces, and buttons (`button-main`, `button-side`).
- [.routing/services/dictionary_service.md](../../.routing/services/dictionary_service.md): Pan-Hispanic Lexicographical Core, Senses, Heteronyms, and WordDetail / SenseDetail / Example / HeteronymItem contracts.
- [app/frontend/src/data/regions.ts](../../app/frontend/src/data/regions.ts): Canonical region definitions.
- Visual references:
  - Big search hero bar: [Image 2: Vocabulary.com style big input bar](../../../.t3/userdata/attachments/e7510c5f-4aee-4fe2-8b04-76e954f1ea52-23a0d289-1b35-4287-a0be-aa43fa19d98c.png)
  - Word of the Day desktop layout: [Image 1: Desktop editorial layout](../../../.t3/userdata/attachments/e7510c5f-4aee-4fe2-8b04-76e954f1ea52-d217cff8-84b2-40bb-811b-d6861b1dc8bd.png)
  - Word of the Day mobile layout: [Image 3: Mobile single-column layout](../../../.t3/userdata/attachments/e7510c5f-4aee-4fe2-8b04-76e954f1ea52-4ebd0a93-b7f0-49e1-bd49-c5e98cbf627b.png)

## Deviations from Visual References

1. **Search bar background & buttons:** No dark/patterned background container; clean minimal surface. No `Random Word` button.
2. **Flags and regional codes:** Strictly NO `flag_emoji` and NO dialect codes (e.g., `CU`, `CL`, `MX`). Render only full textual region names (e.g. `Cuba`, `Chile`, `Mexico`).
3. **Monochrome palette:** All colorful badges and highlights (e.g. vermilion/red badges) are adapted to the strict monochrome design system (`bg-ground`, `bg-light`, `border-light`, `text-main`, `text-side`, `text-dull`).

## Task

1. **Section 1: Hero Search Bar:**
   - Centered headline above input: *"Look up a word, learn it forever."* in `text-2xl sm:text-3xl font-bold text-main tracking-tight`.
   - Big mock search input bar (`w-full max-w-3xl mx-auto h-14 sm:h-16 bg-white border border-light flex items-center px-4 sm:px-6 cursor-pointer`):
     - Left: `Search` icon from `lucide-vue-next`.
     - Placeholder: *"Search words or vocabulary lists from books, exams or textbooks"* in `text-dull text-sm sm:text-base`.
     - Right: Key badge `<span class="binding">Ctrl + K</span>` or *"Word Finder"* label.
   - **Interaction:** Clicking, focusing, or pressing any key in the mock bar triggers the global `SearchModal` (via a shared composable or custom event).

2. **Section 2: Word of the Day (Palabra del Día):**
   - **Data Model:** Use realistic mock data structured 1:1 against canonical backend contracts (`models.WordDetail` containing `senses: SenseDetail[]`, `heteronyms: HeteronymItem[]`, etc.) based on the *Guagua* exemplar.
   - **Header & Title:**
     - Top metadata row: Eyebrows *"Word of the Day"* (left) and *"Regional Lexicon"* (right) in `text-xs font-bold uppercase tracking-widest text-dull`.
     - Huge Lemma headline: `Guagua` in `text-5xl sm:text-6xl md:text-7xl font-black text-main leading-none mt-2`.
     - Part of speech: `noun · feminine` in `text-sm text-dull mt-1 mb-6`.
   - **Editorial Blurb & Etymology:**
     - Desktop (`md:grid md:grid-cols-12 md:gap-8`):
       - Left (`col-span-4`): `ETIMOLOGÍA` / `ETYMOLOGY` in uppercase bold label with concise derivation.
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

3. **Validation:**
   - Run `wails3 task check` from `app/`. It must exit with code 0.

4. **Change log:**
   - Document all changes in `.tickets/009-build-dictionary-route/changelog.md`.

## Verification

All checks are listed in [spec.md](spec.md).
