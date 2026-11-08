# 007 — Spec checklist

Run by the router before accepting the ticket. All items must pass.

## Initial Implementation (Completed)
- [x] These files exist: `app/frontend/src/components/TopBar.vue`, `app/frontend/src/layouts/AppLayout.vue`, `app/frontend/src/views/DictionaryView.vue`, `app/frontend/src/views/ArenaView.vue`.
- [x] `app/frontend/src/router/index.ts` nests `home` (`/`), `dictionary` (`/dictionary`), and `arena` (`/arena`) under `AppLayout`, while `/auth/sign-in` and `/auth/sign-up` stay under `AuthLayout`.
- [x] `AppLayout.vue` renders `<TopBar />` and `<RouterView />`.

## Revision 1 — Responsive TopBar, 12-Col Centering, Reka UI Popovers & Documentation
- [x] `app/frontend/src/components/TOP-BAR.md` exists and documents the 12-column grid alignment, surface hierarchy (`bg-ground` vs `bg-white`), and the `authState` × viewport state machine.
- [x] `app/frontend/src/views/StudioView.vue` and `app/frontend/src/views/SettingsView.vue` exist and are registered as `/studio` and `/settings` children under `AppLayout` in `app/frontend/src/router/index.ts`.
- [x] `TopBar.vue` preserves the outer `<header class="h-20 border-b border-solid border-light bg-ground/50 backdrop-filter backdrop-blur-sm sticky top-10 z-40">` and inner `max-w-6xl h-full w-full mx-auto px-4` container, using `grid grid-cols-12 items-center` (`col-span-6 md:col-span-4` left, `hidden md:flex md:col-span-4 justify-center` center, `col-span-6 md:col-span-4 justify-end` right) so the desktop nav links center-align with `WindowBar`.
- [x] Desktop center nav renders `Dictionary` and `Arena` when `anonymous`, and `Dictionary`, `Arena`, and `Studio` when `logged-in`.
- [x] Mobile right side (`< md`) renders only the `Search` (`button-icon`) and `Menu` (`button-icon`) buttons.
- [x] Clicking the mobile `Menu` button opens a full-screen Reka UI `Dialog` (`DialogRoot` / `DialogContent`, `bg-white`, positioned at `top-10` below `WindowBar` with no Floating UI transform wrapper) with:
  - A top header (`h-20 px-4 border-b border-solid border-light flex items-center justify-between`) containing the `Palabrario` logo and a `DialogClose` button (`button-icon` with `X` icon).
  - `logged-in` content: vertical links (`Dictionary`, `Arena`, `Studio`, `Settings`), user block (`h-12 w-12` user miniature + `gap-3` + vertical `user.name` and `@user.user_name` in `text-dull`), and `Sign Out`.
  - `anonymous` content: vertical links (`Dictionary`, `Arena`), `Sign In` (`button-main`), and `Sign Up` (`button-side`).
- [x] Desktop right side (`>= md`) renders `Search` (`button-icon`) plus:
  - `anonymous`: `Sign In` (`button-main`) and `Sign Up` (`button-side`).
  - `logged-in`: clickable `h-12 w-12 bg-light` user miniature triggering a Reka UI `Popover` (`p-4 bg-white border border-solid border-light flex flex-col gap-3`) with vertical `user.name` + `@user.user_name`, `Settings` link, and `Sign Out`.

## Style rules
- [x] `grep -rnE "rounded|shadow-|hover:|focus:|active:|transition|animate-" app/frontend/src/components/TopBar.vue app/frontend/src/layouts/AppLayout.vue app/frontend/src/views/DictionaryView.vue app/frontend/src/views/ArenaView.vue app/frontend/src/views/StudioView.vue app/frontend/src/views/SettingsView.vue` returns nothing.
- [x] `grep -rnE "side-dark|side-light|text-brand|text-ink|bg-surface|bg-paper" app/frontend/src` returns nothing.
- [x] `grep -rn "@palabrario/ui" app/frontend/src` returns nothing.
- [x] `app/frontend/uno.config.ts` and `app/frontend/index.html` are not modified by the worker.

## Build and Scope
- [x] `cd app && wails3 task check` exits 0.
- [x] `.tickets/007-build-topbar-and-app-layout/changelog.md` includes a Revision 1 section.
