# 033 — Changelog: Align Frontend Auth Contract, Adapter Parity, Router Guards, and UX

## Summary

Aligned the frontend authentication contract (`IAuthAPI`), `WailsAuthAdapter`, `MockAuthAPI`, `useAuth`, Vue Router navigation guards, and auth-sensitive views (`SignInView`, `SignUpView`, `ProfileView`, `ArenaView`, `CreateMatchDialog`) with the hardened backend and Wails BFF session lifecycle.

## Changed Files

- `app/frontend/src/api/contract/auth.contract.ts` — Removed `AuthResponse` and `hasActiveSession()`; updated `signIn`/`signUp` to return `Promise<UserProfile>`; added `onSessionExpired(callback: (reason: string) => void): () => void`.
- `app/frontend/src/api/live/wails/auth.wails.ts` — Removed `hasActiveSession()` and `toAuthResponse()`; unwrapped `Wails.SignIn`/`Wails.SignUp` to return `UserProfile` (throwing on `null`); mapped `404` / not-found errors in `getPublicProfileByUsername` to `null` while propagating server/network errors; subscribed to `"auth:expired"` via `@wailsio/runtime` `Events.On` when `isWailsRuntime()` is true.
- `app/frontend/src/api/mock/fixtures/user.fixture.ts` — Exported `MOCK_DEMO_PASSWORD = "password123"`, `MOCK_SUSPENDED_EMAIL = "suspended@palabrario.app"`, `mockSuspendedProfile` (`mateo_vargas`), `MockAccountRecord`, and `mockAccountRecords` while keeping `mockUsers` intact.
- `app/frontend/src/api/mock/mock.adapter.ts` — Removed `MOCK_TOKEN`, `AuthResponse`, and `hasActiveSession()`. Updated `MockAuthAPI` to verify passwords, reject suspended (`is_active === false`) accounts on `signIn` and hide them (`null`) in `getPublicProfileByUsername`, apply `core-service` normalization and validation rules on `signUp` (distinguishing duplicate email vs username), and support `onSessionExpired` listeners (plus `triggerSessionExpired`).
- `app/frontend/src/composables/useAuth.ts` — Simplified `loadSession()` to call `await api.auth.getMe()` directly (single round-trip); registered `api.auth.onSessionExpired` once on module init to clear `currentUser` and set `sessionNotice`; exposed `sessionNotice` and `clearSessionNotice`.
- `app/frontend/src/router/index.ts` — Added `meta: { requiresAuth: true }` to `studio`, `settings`, and `arena-history`, and `meta: { guestOnly: true }` to `sign-in` and `sign-up`; added `router.beforeEach` calling `await init()` and enforcing auth/guest redirects; watched `currentUser` + `sessionNotice` to redirect protected routes to `/auth/sign-in` on live session expiry.
- `app/frontend/src/components/TopBar.vue` — Removed redundant `onMounted(() => void init())` (handled by `router.beforeEach`) and updated `handleSignOut()` to check `route.meta.requiresAuth` and `/arena/history`.
- `app/frontend/src/views/auth/SignInView.vue` — Rendered monochrome `sessionNotice` alert banner with dismiss button; honored safe internal `route.query.redirect` (must start with `/` and not `//`) after `signIn`, defaulting to `/dictionary`.
- `app/frontend/src/views/auth/SignUpView.vue` — Aligned `isStep1Valid` with backend rules (`^[a-zA-Z0-9][a-zA-Z0-9_-]{2,49}$` after optional `@`, email validation, `8..72` UTF-8 password bytes via `TextEncoder`, and `2..100` rune bounds on optional full name).
- `app/frontend/src/views/ProfileView.vue` — Watched both `route.params.username` and `currentUser?.id_user` (with stale-request guard) so `is_owner` updates immediately on auth changes; distinguished `404`/missing (`"Profile not found."`) from load errors (`"Unable to load profile. Please try again."` with Retry button).
- `app/frontend/src/views/ArenaView.vue` & `app/frontend/src/components/arena/CreateMatchDialog.vue` — Gated `Create Room` and dialog open/submit on `isAuthenticated`, redirecting anonymous users to `/auth/sign-in?redirect=/arena`.
- `.routing/API.md` — Added Section 6 documenting the updated `IAuthAPI` contract and `onSessionExpired` lifecycle hook.

## `IAuthAPI` Contract Changes

```ts
export interface IAuthAPI {
  signIn(req: SignInRequest): Promise<UserProfile>
  signUp(req: SignUpRequest): Promise<UserProfile>
  signOut(): Promise<void>
  getMe(): Promise<UserProfile | null>
  getPublicProfileByUsername(username: string): Promise<PublicUserProfile | null>
  onSessionExpired(callback: (reason: string) => void): () => void
}
```

Removed from `IAuthAPI`:
- `AuthResponse` return wrapper on `signIn` / `signUp`
- `hasActiveSession(): Promise<boolean>`

## Mock Fixture Accounts (Manual Testing)

All mock fixture accounts share the password `password123` (`MOCK_DEMO_PASSWORD`):

| Role / State | Email | Username | Password | `is_active` |
| :--- | :--- | :--- | :--- | :--- |
| Learner (Active) | `aprendiz@palabrario.app` | `camila_rojas` | `password123` | `true` |
| Editor (Active) | `editora@palabrario.app` | `valentina_ospina` | `password123` | `true` |
| Learner (Suspended) | `suspended@palabrario.app` | `mateo_vargas` | `password123` | `false` |

## Verification (`wails3 task check`)

Executed in `app/` — exited with code `0`:

```text
task: [bindings] wails3 generate bindings -clean=true -ts
Wails v3.0.0 › Generate Bindings
 INFO  Processed: 264 Packages, 1 Service, 113 Methods, 6 Enums, 75 Models, 0 Events in 32.6918086s.
 INFO  Output directory: C:\Users\OGSha\OneDrive\Escritorio\Developer\Web Development\Palabrario\app\frontend\bindings
task: [check] bun run --cwd frontend type-check
$ vue-tsc --noEmit
task: [check] bun run --cwd frontend build
$ vite build --mode production
vite v8.3.5 building client environment for production...
transforming...
✓ 2500 modules transformed.
rendering chunks...
[unocss] Fetching web fonts: https://api.fontshare.com/v2/css?f[]=satoshi@1&display=swap
[unocss] Web fonts fetched.
computing gzip size...
dist/index.html                                     0.86 kB │ gzip:  0.45 kB
...
✓ built in 10.37s
task: [check] go build .
```

## UnoCSS Shortcut Proposals

None (`uno.config.ts` was not modified).
