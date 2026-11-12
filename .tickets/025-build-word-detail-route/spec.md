# 025 — Spec checklist

Run by the router before accepting the ticket. All items must pass.

- [x] `app/frontend/src/views/WordDetailView.vue` exists and mounts under `AppLayout`.
- [x] Route `/dictionary/word/:id` registered in `app/frontend/src/router/index.ts` with no `alias` and breadcrumb hierarchy (`Home` -> `Dictionary` -> `Word Detail`).
- [x] Navigating to `/dictionary/word/:id` dynamically updates the breadcrumb leaf with the word's lemma via `breadcrumbLeafLabel`.
- [x] Unmounting `WordDetailView.vue` resets `breadcrumbLeafLabel` to null.
- [x] `WordDetailView.vue` guards invalid `:id` values (`!Number.isInteger(id) || id <= 0`) without calling the API, and fetches valid entries via `api.dictionary.getWordById(id)` using `import { api } from "@/api"`. Zero hardcoded word fixtures exist in views.
- [x] Loading state renders while data is fetching (`"Loading dictionary entry..."`).
- [x] Not found / error fallback (`!word`) renders `"Word not found."` and a `<RouterLink to="/dictionary">` (not a button) to return to `/dictionary`.
- [x] `WordOfTheDayCard.vue` supports `eyebrow` (defaults to `"Word of the Day"`), `headingTag` (`"h1" | "h2"`, defaults to `"h2"`), and `footerLink` (`"entry" | "dictionary"`, defaults to `"entry"`).
- [x] `WordOfTheDayCard.vue` renders the `Report Word` footer button with `button-side` (not `button-main`).
- [x] When rendered in `WordDetailView.vue`, `WordOfTheDayCard` shows eyebrow `"Dictionary Entry"`, renders the lemma as `<h1>`, and renders `← Back to Dictionary` linking to `/dictionary` in the footer.
- [x] Selecting a dictionary result in `SearchModal.vue` navigates to `/dictionary/word/${item.id_word}`.
- [x] "See full entry →" link in `WordOfTheDayCard.vue` on `/dictionary` links to `/dictionary/word/${word.id_word}`.
- [x] Strict font hygiene: zero `font-mono` occurrences; single font family (Satoshi) honored.
- [x] Ground rules enforced: zero `flag_emoji`, zero dialect codes (`CU`, `CL`, `MX`); full textual region names only.
- [x] Typecheck and build pass: `wails3 task check` in `app/` exits with code 0.
- [x] No application code outside declared scope, backend code, or `.context/` modified.
- [x] `.tickets/025-build-word-detail-route/changelog.md` exists and is populated.
- [x] Commit (if made): Conventional Commits, English, detailed body.
