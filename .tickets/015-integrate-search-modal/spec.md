# 015 — Spec checklist

Run by the router before accepting the ticket. All items must pass.

- [ ] `SearchModal.vue` imports models from `@bindings` and has no local mock data interfaces.
- [ ] Backend search methods (`api.SearchWords`, etc.) are wired with debouncing.
- [ ] Featured entry card accurately synchronizes with active keyboard navigation.
- [ ] No `flag_emoji` or dialect codes are rendered.
- [ ] Full keyboard navigation (`↑`, `↓`, `Enter`, `Esc`) remains functional.
- [ ] `wails3 task check` in `app/` exits with code 0.
- [ ] Changes stay within declared scope.
- [ ] `.tickets/015-integrate-search-modal/changelog.md` exists and is populated.
- [ ] Commit (if made): Conventional Commits, English, detailed body.
