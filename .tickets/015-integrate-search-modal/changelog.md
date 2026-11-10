# 015 — Changelog

## `app/frontend/src/components/SearchModal.vue`
- Removed `MockWordResult`, `MockStudyDeck`, `MOCK_DICTIONARY_ENTRIES`, `MOCK_STUDY_DECKS` and the local filter computeds. The default query is now empty (was `"guagua"`).
- Imports `api` from `@/api` and the canonical `WordSearchResult`, `DeckSearchResult`, `CollectionSearchResult` types from `@bindings/palabrario/pkg/models/index.js`. No import from `@bindings/palabrario/app/app.js`.
- Added `dictionaryResults`, `deckResults`, `collectionResults`, `isLoading` refs.
- Search is debounced with `useDebounceFn` (250 ms). It calls `api.dictionary.searchWords(q, "", 20)`, `api.studio.searchDecks(q, 20)` and `api.studio.searchCollections(q, 20)` in parallel via `Promise.all`.
- Clearing the query empties all results immediately and issues no API request. A request counter discards stale responses (superseded or cleared queries); errors clear the results.
- `TabType` is now `"all" | "dictionary" | "decks" | "collections"`, with a fourth `Collections` tab button (`button-main` active, `button-side--light` inactive).
- `SelectableItem` supports `dictionary | deck | collection`, with unique ids (`word-<id_word>`, `deck-<id_deck>`, `col-<id_collection>`). Arrow/Enter/Escape navigation works across the combined list.
- Featured Entry card shows the highlighted word, falls back to the first dictionary result, and is hidden when there are no word results (the right column then takes the full width).
- New "Matching Collections" section: title, description (2-line clamp), `badge-side` "N decks", and `@owner`.
- Navigation: word -> `/dictionary?q=<lemma>`, deck -> `/studio`, collection -> `/collections`; the modal closes in each case.
- New empty states: "Start typing to search", "Searching...", "No results found". Footer shows the real result count.
- Only `region_name` (full names) is rendered; no flag emoji or country/dialect codes. Monochrome tokens and `rounded-none` geometry are unchanged.

## `app/frontend/src/views/DictionaryView.vue`
- Removed the `@/data/mockWordOfTheDay` import (`mockWordOfTheDay`, `mockHeteronymGroups`).
- Fetches `api.dictionary.getWordOfTheDay()` on mount into `word: Ref<WordDetail | null>`, with an `isLoading` flag. Renders `WordOfTheDayCard` once loaded, with loading and empty fallbacks.
- `WordOfTheDayCard` falls back to grouping `word.heteronyms` itself when no `heteronymGroups` prop is passed.

## Notes
- The router (out of this ticket's scope) has no `/collections` route yet; collection selection pushes `/collections` as specified. Today this resolves to no matched route until that route is added.
- `WordOfTheDayCard.vue` (out of scope) still imports the `HeteronymGroup` type from `@/data/mockWordOfTheDay`, and `formatFrequency` there still reads `region_code === "CL"`. Neither is rendered, but both are worth cleaning up in a later ticket.

## Verification
- `wails3 task check` in `app/` exits 0 (bindings, `vue-tsc --noEmit`, `vite build`, `go build`).
