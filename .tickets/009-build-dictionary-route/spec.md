# 009 — Spec checklist

Run by the router before accepting the ticket. All items must pass.

- [ ] `/dictionary` renders under `AppLayout` without console errors.
- [ ] Hero search bar renders with headline "Look up a word, learn it forever." and big mock input.
- [ ] Clicking or focusing the mock search input opens the global `SearchModal`.
- [ ] Word of the Day section renders lemma, etymology, friendly blurb, senses, examples, and heteronym equivalents.
- [ ] Mobile-first responsive layout: single-column stacked on mobile screens, 2-column layout on desktop (`md:`).
- [ ] Strictly NO `flag_emoji` and NO regional codes (`CU`, `CL`, `MX`); regions display as full names (`Cuba`, `Chile`, etc.).
- [ ] Strict monochrome styling upheld: no unauthorized colors, no rounded corners (sharp solid borders).
- [ ] Data structures conform 1:1 to `models.WordDetail` contracts.
- [ ] `wails3 task check` in `app/` exits with code 0.
- [ ] No files modified outside specified scope.
- [ ] `.tickets/009-build-dictionary-route/changelog.md` exists and is populated.
- [ ] Commit (if made): Conventional Commits, English, detailed body.
