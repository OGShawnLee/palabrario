# Changelog — Ticket 013: Build Global Floating Search Modal

## Created Files
- `app/frontend/src/components/SearchModal.vue`:
  - Implemented responsive floating search modal using Reka UI (`DialogRoot`, `DialogPortal`, `DialogOverlay`, `DialogContent`, `DialogTitle`, `DialogDescription`, `DialogClose`).
  - Added global keyboard listener using `@vueuse/core` (`useEventListener`) for `Ctrl + K` and `Cmd + K` to toggle the modal from anywhere.
  - Desktop view floats slightly above vertical center (`top-[15%]`), centered horizontally, max width `max-w-4xl`, solid sharp borders (`border border-light`), white card surface (`bg-white`) over translucent backdrop (`dialog-overlay`).
  - Mobile view fills the full screen (`fixed inset-0 bg-white z-50`).
  - Implemented 3 filter tabs: `All`, `Dictionary`, and `Study Decks` (strictly omitting `Random Word` button and `By Region` tab per ticket requirements) using `button-main` and `button-side--light`.
  - Search input header includes search icon, autofocus on open, query clear button (`X`), desktop `Ctrl + K` badge, and mobile close button.
  - Followed strict ground rule: ZERO flag emojis and ZERO dialect codes (`CU`, `CL`, `MX`) rendered; all regions render full text names (`Cuba`, `Chile`, `Mexico`, etc.).
  - Featured entry card with lemma, part of speech, full region name badge, definition, and "See full entry →" navigation.
  - Synchronized navigation: Featured entry dynamically updates as the user navigates through dictionary results with `ArrowDown` / `ArrowUp` or mouse hover.
  - Live Dictionary Results with uppercase counter and interactive list highlighting lemma, region name, tag (e.g. Plural), and definition.
  - Matching Study Decks section with responsive grid of deck cards showing title and full region name.
  - Footer with dynamic result counter and desktop keyboard navigation hints using `<span class="binding">` and `<span class="text-xs text-black">` wrapped in `<div class="flex items-center gap-2">`.
  - Full keyboard navigation: `ArrowDown` and `ArrowUp` cycle through active search items, `Enter` selects highlighted item and navigates/closes, `Escape` closes the modal.
  - Scoped styling block using `--uno` directives for clean semantic encapsulation (`.binding-group`, `.deck-item`, `.section-label`).

## Modified Files
- `app/frontend/uno.config.ts`:
  - Added `transformerDirectives({ applyVariable: "--uno" })` for scoped styles.
  - Added shortcuts: `binding`, `button-side--light`, and `dialog-overlay`.
- `app/frontend/src/components/TopBar.vue`:
  - Added `defineEmits<{ "open-search": [] }>()`.
  - Wired search icon button click handler to emit `open-search`.
- `app/frontend/src/layouts/AppLayout.vue`:
  - Imported `ref` and `SearchModal.vue`.
  - Declared `const isSearchOpen = ref(false)`.
  - Wired `@open-search="isSearchOpen = true"` on `<TopBar />`.
  - Mounted `<SearchModal v-model:open="isSearchOpen" />`.
- `.routing/CONTEXT.md` & `.routing/DESIGN.md`:
  - Enforced ground rule prohibiting `flag_emoji` and regional codes in favor of full region names.
  - Added `binding` key shortcut specification.

## Verification & Checks
- Ran `wails3 task check` from `app/`:
  - `wails3 generate bindings`: 0 errors
  - `vue-tsc --noEmit`: 0 errors
  - `vite build --mode production`: built in 12.23s (exit code 0)
  - `go build .`: exit code 0
- Spec validation checklist: All items verified and passing.
