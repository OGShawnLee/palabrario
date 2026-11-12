# 025 — Changelog

## Created
- `app/frontend/src/views/WordDetailView.vue`: word detail view for `/dictionary/word/:id`. Refs: `word`, `isLoading`. Guards invalid ids (`!Number.isInteger(id) || id <= 0`) without calling the API; otherwise calls `api.dictionary.getWordById(id)`. Sets `breadcrumbLeafLabel` to the lemma (or null) and resets it on unmount. Renders loading, "Word not found." fallback (`RouterLink` to `/dictionary`), or `WordOfTheDayCard`.

## Modified
- `app/frontend/src/components/dictionary/WordOfTheDayCard.vue`: new optional props via `withDefaults`:
  - `eyebrow?: string` (default `"Word of the Day"`)
  - `headingTag?: "h1" | "h2"` (default `"h2"`)
  - `footerLink?: "entry" | "dictionary"` (default `"entry"`)
  - `Report Word` button now `button-side` in both modes.
- `app/frontend/src/components/SearchModal.vue`: `selectDictionaryResult` navigates to `/dictionary/word/${item.id_word}`.
- `app/frontend/src/router/index.ts`: registered `word-detail` route (no alias) with breadcrumb Home / Dictionary / Word Detail.

## Validation
- `wails3 task check` in `app/`: TypeScript, Vite build, and `go build` completed successfully.

## Proposals
- None (`uno.config.ts` untouched).
- Note: `formatFrequency` in `WordOfTheDayCard.vue` still checks `regionCode === "CL"` (pre-existing, out of ticket scope).
