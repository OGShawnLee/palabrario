# 015 — Frontend: Integrate SearchModal & DictionaryView with API Client

| Field | Value |
| :--- | :--- |
| **ID** | 015 |
| **Status** | in-progress |
| **Model** | Claude Sonnet 5.5 (`claude-sonnet-5-5`, Claude provider) |
| **Role** | Frontend Implementer |
| **Spec** | [spec.md](spec.md) |

> **Draft Note:**  
> This ticket connects `SearchModal.vue` and `DictionaryView.vue` to the universal API client (`import { api } from "@/api"`), eliminating all hardcoded mock arrays and local mock interfaces. It utilizes `api.dictionary.searchWords`, `api.studio.searchDecks`, `api.studio.searchCollections`, and `api.dictionary.getWordOfTheDay()`, seamlessly supporting both live Wails execution and offline exploration without hardcoded UI mock data.

## Scope

May create or edit:
- `app/frontend/src/components/SearchModal.vue`
- `app/frontend/src/views/DictionaryView.vue`
- `.tickets/015-integrate-search-modal/changelog.md` (create)

Must not touch:
- `app/frontend/uno.config.ts` (shortcuts already exist)
- Backend code, Go services, or `.context/`
- Any file outside the above scope

## Inputs

- [.routing/API.md](../../.routing/API.md): Architectural rules (Rule 1: Never import `@bindings/.../app` in UI; always import `api` from `@/api`).
- [app/frontend/src/api/index.ts](../../app/frontend/src/api/index.ts): Universal `api` singleton (`api.dictionary`, `api.studio`, `api.auth`).
- [app/frontend/bindings/palabrario/pkg/models/index.ts](../../app/frontend/bindings/palabrario/pkg/models/index.ts): Canonical models (`WordSearchResult`, `DeckSearchResult`, `CollectionSearchResult`, `WordDetail`).
- [app/frontend/src/components/SearchModal.vue](../../app/frontend/src/components/SearchModal.vue): Existing search modal to be wired.
- [app/frontend/src/views/DictionaryView.vue](../../app/frontend/src/views/DictionaryView.vue): Dictionary route view to be wired to `api.dictionary.getWordOfTheDay()`.
- [.routing/CONTEXT.md](../../.routing/CONTEXT.md): Ground rules (strictly NO `flag_emoji`, NO regional codes like `CU`, `CL`, `MX`; full region names only).
- [.routing/DESIGN.md](../../.routing/DESIGN.md): Monochrome tokens (`main`, `side`, `dull`, `light`, `ground`), `section-label`, `badge`, `badge-side`, `binding`.

## Task

1. **Remove Local Mock Data & Import Universal API Client:**
   - In `app/frontend/src/components/SearchModal.vue`:
     - Delete local interfaces `MockWordResult` and `MockStudyDeck`.
     - Delete static arrays `MOCK_DICTIONARY_ENTRIES` and `MOCK_STUDY_DECKS`.
     - Import `api` from `@/api`:
       ```ts
       import { api } from "@/api"
       ```
     - Import canonical DTO models from `@bindings/palabrario/pkg/models/index.js`:
       `type WordSearchResult`, `type DeckSearchResult`, `type CollectionSearchResult`.
     - **Strict Rule:** Never import from `@bindings/palabrario/app/app.js` in Vue components. All requests flow through `api`.

2. **Update Filter Tabs to 4 Categories:**
   - Update `TabType` definition:
     ```ts
     type TabType = "all" | "dictionary" | "decks" | "collections"
     ```
   - Render 4 tab filter buttons in the modal sub-header:
     - `All`
     - `Dictionary`
     - `Study Decks`
     - `Collections`
   - Active tab uses `button-main`, inactive tabs use `button-side--light`.

3. **Wire Debounced Reactive Backend Queries via API Client:**
   - Create reactive result refs:
     - `dictionaryResults = ref<WordSearchResult[]>([])`
     - `deckResults = ref<DeckSearchResult[]>([])`
     - `collectionResults = ref<CollectionSearchResult[]>([])`
     - `isLoading = ref(false)`
   - Use `@vueuse/core` (`useDebounceFn` with 200–250ms) to trigger search on input changes:
     - When `query.trim()` is empty:
       - Clear results: `dictionaryResults.value = []`, `deckResults.value = []`, `collectionResults.value = []`.
       - Do not issue API requests.
     - When `query.trim()` is not empty:
       - Set `isLoading.value = true`.
       - Call in parallel via `api`:
         ```ts
         const [words, decks, collections] = await Promise.all([
           api.dictionary.searchWords(trimmedQuery, "", 20),
           api.studio.searchDecks(trimmedQuery, 20),
           api.studio.searchCollections(trimmedQuery, 20),
         ])
         ```
       - Populate the reactive refs with the returned results and reset `isLoading.value = false`.

4. **Wire DictionaryView to API Client:**
   - In `app/frontend/src/views/DictionaryView.vue`:
     - Remove `import { mockWordOfTheDay, mockHeteronymGroups } from "@/data/mockWordOfTheDay"`.
     - Import `api` from `@/api` and `WordDetail` from `@bindings/palabrario/pkg/models/index.js`.
     - Fetch word on mount:
       ```ts
       const word = ref<WordDetail | null>(null)
       const isLoading = ref(true)

       onMounted(async () => {
         try {
           word.value = await api.dictionary.getWordOfTheDay()
         } finally {
           isLoading.value = false
         }
       })
       ```
     - Pass `word` to `WordOfTheDayCard.vue` when loaded.

5. **Unified Selectable Items & Keyboard Navigation:**
   - Update `SelectableItem`:
     ```ts
     interface SelectableItem {
       id: string
       type: "dictionary" | "deck" | "collection"
       item: WordSearchResult | DeckSearchResult | CollectionSearchResult
     }
     ```
   - Populate `selectableItems` based on active `selectedTab`:
     - If `all` or `dictionary`: add words with unique prefix (e.g. `word-${item.id_word}`).
     - If `all` or `decks`: add decks with unique prefix (e.g. `deck-${deck.id_deck}`).
     - If `all` or `collections`: add collections with unique prefix (e.g. `col-${col.id_collection}`).
   - Maintain keyboard navigation (`ArrowUp`, `ArrowDown`, `Enter`, `Escape`) across the combined list.
   - Synchronize Featured Entry (left card):
     - If the currently highlighted item is a word, display it in the Featured Entry card.
     - Otherwise, display the first word from `dictionaryResults` (or hide the featured card if `dictionaryResults` is empty).

5. **Render Collections Results UI:**
   - Add a `"Matching Collections"` section visible when `selectedTab === 'all' || selectedTab === 'collections'` and `collectionResults.length > 0`.
   - Card layout matching the monochrome design system:
     - Collection title in `text-sm font-bold text-main`.
     - Description or subtitle in `text-xs text-side leading-relaxed line-clamp-2 mt-1`.
     - Metadata badges: total decks (`${col.deck_count} decks` using `badge-side`), owner username (`@${col.owner_user_name}` in `text-xs text-dull`).
   - Clicking or pressing `Enter` on a collection navigates to `/collections` and closes the modal.

6. **Selection Navigation Handling:**
   - Dictionary item selected: close modal and navigate to `/dictionary` with query parameter:
     `router.push({ path: '/dictionary', query: { q: item.lemma } })`.
   - Deck item selected: close modal and navigate to `/decks` (or `/studio`).
   - Collection item selected: close modal and navigate to `/collections`.

7. **Zero Regional Codes & Flag Emojis Rule:**
   - Strictly NO flag emojis and NO dialect codes (`CU`, `CL`, `MX`).
   - Always display full region names (`item.region_name`) provided by the backend models.

8. **Validation:**
   - Run `wails3 task check` from `app/`. It must exit with code 0.

9. **Change log:**
   - Document all changes in `.tickets/015-integrate-search-modal/changelog.md`.

## Verification

All checks are listed in [spec.md](spec.md).
