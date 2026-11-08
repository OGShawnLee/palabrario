# 004 — Spec checklist

Run by the router before accepting the ticket. All items must pass.

## Files and structure
- [x] `app/frontend/src/router/index.ts` exists and configures `createWebHashHistory` with routes `/`, `/auth/sign-in`, and `/auth/sign-up`.
- [x] `app/frontend/src/views/auth/SignInView.vue` exists and contains `Sign In`.
- [x] `app/frontend/src/views/auth/SignUpView.vue` exists and contains `Sign Up`.
- [x] `app/frontend/src/main.ts` imports `virtual:uno.css`, `@unocss/reset/tailwind-compat.css`, and registers `router`.
- [x] `app/frontend/src/App.vue` renders `<RouterView />`.

## Landing page content and design rules
- [x] `app/frontend/src/views/HomeView.vue` contains the exact eyebrow text `The dictionary of our differences`.
- [x] `app/frontend/src/views/HomeView.vue` contains the three headline lines: `One shared language.`, `Twenty ways to say`, and `the same thing.`.
- [x] `app/frontend/src/views/HomeView.vue` contains the subtitle `In Mexico it's chamba, in Argentina it's laburo, in Spain it's curro. Explore the living vocabulary that unites and divides five hundred million Spanish speakers.`.
- [x] `app/frontend/src/views/HomeView.vue` renders `<RouterLink to="/auth/sign-in"` with `button-main` (`Sign In`) and `<RouterLink to="/auth/sign-up"` with `button-side` (`Sign Up`).
- [x] `app/frontend/src/views/HomeView.vue` uses mobile-first responsive classes (`flex-col sm:flex-row`, `w-full sm:w-auto`, `text-4xl sm:text-6xl md:text-7xl`).
- [x] `grep -rn "rounded" app/frontend/src` returns nothing.
- [x] `grep -rn "shadow" app/frontend/src` returns nothing.
- [x] `grep -rnE "hover:|focus:|active:" app/frontend/src` returns nothing.
- [x] `grep -rnE "side-dark|side-light|text-ink|bg-surface|bg-brand" app/frontend/src` returns nothing.
- [x] `grep -rn "@palabrario/ui" app/` returns nothing (excluding `node_modules`).

## Build and check
- [x] `cd app && wails3 task check` exits with code 0.

## Scope and change log
- [x] `app/frontend/uno.config.ts` was not modified.
- [x] Changes are limited to `app/frontend/src/` and `.tickets/004-build-landing-page/`.
- [x] `.tickets/004-build-landing-page/changelog.md` exists and covers every required item.
