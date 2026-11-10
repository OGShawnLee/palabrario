# 015 — Spec checklist

Run by the router before accepting the ticket. All items must pass.

- [x] `SearchModal.vue` imports `api` from `@/api` (`api.dictionary.searchWords`, `api.studio.searchDecks`, `api.studio.searchCollections`) and canonical models. Zero direct imports from `@bindings/.../app.js`.
- [x] `DictionaryView.vue` imports `api` from `@/api` (`api.dictionary.getWordOfTheDay()`); static `mockWordOfTheDay.ts` import is removed from the view.
- [x] Local mock data interfaces (`MockWordResult`, `MockStudyDeck`) and static mock arrays are completely removed from `SearchModal.vue`.
- [x] 4 filter tabs exist and toggle correctly: `All`, `Dictionary`, `Study Decks`, `Collections`.
- [x] Search queries are triggered through `api` with debouncing (200–250ms) using `@vueuse/core`.
- [x] Clearing the query clears all results cleanly without unneeded network/IPC requests.
- [x] Matching Collections section renders when collections match, showing collection title, deck count badge, and owner.
- [x] Featured entry card synchronizes with active keyboard navigation when highlighting words, and falls back cleanly when highlighting non-words.
- [x] Full keyboard navigation (`↑`, `↓`, `Enter`, `Esc`) operates seamlessly across all three result types.
- [x] Selection of words, decks, and collections navigates to the proper route (`/dictionary`, `/decks` or `/studio`, `/collections`) and closes the modal.
- [x] Strictly NO `flag_emoji` and NO regional codes (`CU`, `CL`, `MX`); all regions render as full names (`Cuba`, `Chile`, etc.).
- [x] Strict monochrome styling upheld: no unauthorized colors, sharp solid borders (`rounded-none`).
- [x] `wails3 task check` in `app/` exits with code 0.
- [x] No files modified outside specified scope.
- [x] `.tickets/015-integrate-search-modal/changelog.md` exists and is populated.
- [x] Commit (if made): Conventional Commits, English, detailed body.
