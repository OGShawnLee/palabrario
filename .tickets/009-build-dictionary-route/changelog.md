# 009 — Build the Dictionary Route (/dictionary) — Changelog

## Summary

Implemented the `/dictionary` route featuring the hero search trigger bar and the Word of the Day (Palabra del Día) editorial presentation, adhering strictly to the monochrome design system and canonical Pan-Hispanic Lexicographical contracts (`WordDetail`, `SenseDetail`, `Example`, `HeteronymItem`).

## Scope & File Changes

### Created Files
1. `app/frontend/src/composables/useSearchModal.ts`:
   - Module-scoped singleton reactive composable managing search modal visibility.
   - Exports `{ isOpen, open, close, toggle }`.
2. `app/frontend/src/data/mockWordOfTheDay.ts`:
   - Strongly-typed mock data for *Guagua* conforming strictly to `WordDetail` (`pkg/models/dictionary.go` and TypeScript Wails bindings).
   - Contains 2 dialectal senses in Spanish (Cuba - autobús de transporte público, CEFR B1; Chile - bebé o criatura de pecho, CEFR A2) with full regional definitions, cultural notes, and contextual examples with attributions.
   - Contains categorized heteronym cross-equivalents in Spanish (transporte de pasajeros: Mexico · Camión, Argentina · Colectivo, Chile · Micro, Spain · Autobús; primera infancia: Chile · Guagua, Spain · Bebé, Mexico · Bebé, Argentina · Bebé).
3. `app/frontend/src/components/dictionary/WordOfTheDayCard.vue`:
   - Modular editorial card rendering:
     - Header eyebrows (`Word of the Day`, `Regional Lexicon`) and huge lemma headline with grammatical category.
     - Responsive 12-column grid for Etymology (col-span-4) and Friendly Blurb narrative (col-span-8), collapsing to single column on mobile.
     - 2-column dialectal senses grid on desktop (`md:grid-cols-2`), stacked on mobile.
     - Badges with full regional names (`Cuba`, `Chile`) and CEFR level tags (`B1`, `A2`).
     - Numbered senses with short definitions, frequency/register, and blockquoted contextual examples with solid left border.
     - Per-sense action buttons: "Add to Studio" (`button-small--main`) and "Report Sense" (`button-side`).
     - Categorized cross-regional heteronym pills with full textual region names.
     - Action footer with "Report Word" (`button-main`) and `RouterLink` ("See full entry →").
4. `.tickets/009-build-dictionary-route/changelog.md`:
   - This changelog document.

### Modified Files
1. `app/frontend/src/layouts/AppLayout.vue`:
   - Refactored local search state to use singleton `useSearchModal`.
   - Wired `<TopBar @open-search="open" />` and `<SearchModal v-model:open="isOpen" />`.
2. `app/frontend/src/views/DictionaryView.vue`:
   - Replaced placeholder view with full layout:
     - Hero search section: centered headline *"Every dialect. Every nuance. In one place."* and prominent combobox trigger bar with Search icon, placeholder *"Search lemmas, regional variants, or expressions..."*, and `Ctrl + K` binding badge. Wired to trigger `open()` on click, Enter, and Space keys.
     - Word of the Day section embedding `<WordOfTheDayCard>`.
3. `app/frontend/uno.config.ts`:
   - Added shortcuts `badge`, `badge-main`, `badge-side`, `button-side--light`, `button--small`, and `dialog-overlay`.
4. `.routing/DESIGN.md`:
   - Documented badges, new button variants, and dialog overlay.
5. `.routing/PROGRESS.md`:
   - Tracked new shortcuts under design system gaps.

## Design & Architecture Compliance

1. **Monochrome Palette & Sharp Shapes:**
   - Adheres strictly to `uno.config.ts` theme colors (`main`, `side`, `dull`, `light`, `ground`, `white`).
   - Zero rounded corners (`border-radius: 0`).
   - All badges (CEFR `badge-main`, region `badge-side`) adapted to monochrome system.
2. **Strict Regional Names Rule:**
   - Strictly NO flag emojis.
   - Strictly NO dialect codes (e.g., `CU`, `CL`, `MX`). All regions render as canonical full names (`Cuba`, `Chile`, `Mexico`, `Argentina`, `Spain`).
3. **Responsive Breakdown:**
   - Mobile-first approach: single-column vertical layout on mobile screens (`< md`), expanding to 12-column and 2-column grids on desktop (`md:`).

## Adopted Shortcuts in `uno.config.ts`

- `badge`: `"h-6 px-2.5 inline-flex items-center border border-solid rounded-none font-bold text-xs"`
- `badge-main`: `"badge bg-main border-black font-bold text-white"`
- `badge-side`: `"badge bg-ground border-light text-black"`
- `button-side--light`: `"button bg-white text-main border-light"`
- `button--small`: `"!h-10 px-2.5 text-xs"`
- `dialog-overlay`: `"fixed inset-0 z-50 bg-black/40 backdrop-filter backdrop-blur-xs"`

## Verification

- Ran `bun run --cwd frontend type-check`: Vue TypeScript compiler (`vue-tsc --noEmit`) exited with code 0 (no errors).
- Ran `wails3 task check`: Bindings generation + TypeScript check + Vite production build + Go build exited with code 0.
