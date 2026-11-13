# 033 — Spec: Align Frontend Auth Contract, Adapter Parity, Router Guards, and UX

- [x] `app/frontend/src/api/contract/auth.contract.ts` exports `IAuthAPI` without `AuthResponse` or `hasActiveSession()`, with `signIn`/`signUp` returning `Promise<UserProfile>`, `getMe()` returning `Promise<UserProfile | null>`, and `onSessionExpired(...)`.
- [x] `app/frontend/src/api/live/wails/auth.wails.ts` removes `toAuthResponse` and `hasActiveSession()` and subscribes to `"auth:expired"` via `@wailsio/runtime` `Events.On`.
- [x] `app/frontend/src/api/mock/mock.adapter.ts` verifies passwords, rejects `is_active === false` accounts on `signIn`, hides `is_active === false` profiles (`null`) in `getPublicProfileByUsername`, distinguishes email vs username conflicts on `signUp`, and matches `core-service` validation rules.
- [x] `app/frontend/src/composables/useAuth.ts` removes `hasActiveSession()` round-trip in `loadSession()`, hooks `api.auth.onSessionExpired`, and exposes `sessionNotice`.
- [x] `app/frontend/src/router/index.ts` runs `await init()` inside `router.beforeEach`, protects `/studio`, `/settings`, and `/arena/history` (`meta.requiresAuth`), and guards `/auth/*` (`meta.guestOnly`).
- [x] `app/frontend/src/views/auth/SignUpView.vue` validates username regex (`^[a-zA-Z0-9][a-zA-Z0-9_-]{2,49}$` after optional `@`), email format, `8..72` UTF-8 byte passwords, and `name` bounds.
- [x] `app/frontend/src/views/ProfileView.vue` re-fetches on `currentUser?.id_user` changes and distinguishes missing profiles from load failures; `ArenaView.vue` gates `Create Room` on authentication.
- [x] `.routing/API.md` reflects the updated `IAuthAPI` contract.
- [x] `wails3 task check` (from `app/`) exits with code `0` (`bindings`, `vue-tsc --noEmit`, `vite build`, and `go build .`).
- [x] `git status --short` confirms no files outside the declared `Scope` were created or modified.
- [x] No unauthorized `git commit`, `git push`, or `git merge` was run. (If committed after explicit user authorization: English, Conventional Commits, detailed body.)
