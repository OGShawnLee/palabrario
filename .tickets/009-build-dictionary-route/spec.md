# 009 — Spec checklist

Run by the router before accepting the ticket. All items must pass.

- [x] `app/frontend/src/composables/useSearchModal.ts` exists, has module-scoped reactive state, and exports `{ isOpen, open, close, toggle }`.
- [x] `app/frontend/src/layouts/AppLayout.vue` uses `useSearchModal` and binds `<SearchModal v-model:open="isOpen" />` and `<TopBar @open-search="open" />`.
- [x] `/dictionary` renders under `AppLayout` without console errors.
- [x] Hero search bar renders with headline "Every dialect. Every nuance. In one place." and keyboard binding badge `Ctrl + K`.
- [x] Clicking or keyboard-triggering (Enter/Space) the hero mock search bar opens the global `SearchModal`.
- [x] `app/frontend/src/data/mockWordOfTheDay.ts` exists and strictly conforms to `WordDetail` contracts (`senses: SenseDetail[]`, `heteronyms: HeteronymItem[]`).
- [x] Word of the Day section renders lemma, etymology, friendly blurb, senses, examples, and heteronym equivalents.
- [x] Mobile-first responsive layout: single-column stacked on mobile screens, 2-column layout on desktop (`md:`).
- [x] Strictly NO `flag_emoji` and NO regional codes (`CU`, `CL`, `MX`); regions display as full names (`Cuba`, `Chile`, etc.).
- [x] Strict monochrome styling upheld: no unauthorized colors, no rounded corners (sharp solid borders).
- [x] `wails3 task check` in `app/` exits with code 0.
- [x] No files modified outside specified scope.
- [x] `.tickets/009-build-dictionary-route/changelog.md` exists and is populated.
- [x] Commit (if made): Conventional Commits, English, detailed body.
