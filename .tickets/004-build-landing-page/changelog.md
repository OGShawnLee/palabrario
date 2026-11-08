# 004 — Change log

## Files
- Created `app/frontend/src/router/index.ts` (hash history; routes `/`, `/auth/sign-in`, `/auth/sign-up`).
- Created `app/frontend/src/views/auth/SignInView.vue` and `SignUpView.vue` (placeholders).
- Edited `app/frontend/src/views/HomeView.vue` (static monochrome landing hero with RouterLinks).
- Edited `app/frontend/src/App.vue` (renders `<RouterView />`).
- Edited `app/frontend/src/main.ts` (registers router; imports now `virtual:uno.css` then `@unocss/reset/tailwind-compat.css`, as the spec requires. The original file imported `@unocss/reset/tailwind.css` first.)

## Validation
- `cd app && wails3 task check`: exit code 0 (bindings, vue-tsc, vite build, go build all passed).
- Grep for `rounded|shadow|hover:|focus:|active:|side-dark|side-light|text-ink|bg-surface|bg-brand` in `app/frontend/src`: no matches.

## uno.config.ts proposals
None.

## Problems
None. Note: `wails3 task check` regenerates `app/frontend/bindings`; it may show up in git status.
