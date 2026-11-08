# 013 — Spec checklist

Run by the router before accepting the ticket. All items must pass.

- [x] `app/frontend/src/components/SearchModal.vue` exists and is implemented with Reka UI Dialog.
- [x] Global keyboard shortcut (`Ctrl + K` / `Cmd + K`) opens the modal via `@vueuse/core`.
- [x] TopBar search button click opens the modal.
- [x] Desktop layout floats slightly above the center with a backdrop; mobile layout fills the full screen.
- [x] Only 3 tabs are rendered: `All`, `Dictionary`, and `Study Decks` (no `Random Word` button, no `By Region` tab).
- [x] No `flag_emoji` or regional codes (e.g. `CU`, `CL`, `MX`) are rendered; regions display only the full region name (e.g. `Cuba`, `Chile`, `Mexico`).
- [x] Keyboard hints in footer use `<span class="binding">` and `<span class="text-xs text-black">` wrapped in `<div class="flex items-center gap-2">`.
- [x] Interactive keyboard navigation: Up / Down arrow moves highlight, `Enter` selects, `Esc` closes.
- [x] Realistic mock data displays results matching official domain structures.
- [x] `cd app && wails3 task check` exits with code 0.
- [x] No files outside the scope were modified.
- [x] `.tickets/013-build-search-modal/changelog.md` exists and is populated.
- [x] Commit (if made): Conventional Commits, English, detailed body.

