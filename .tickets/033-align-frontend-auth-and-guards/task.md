# 033 — Align Frontend Auth Contract, Adapter Parity, Router Guards, and UX

| Field | Value |
| :--- | :--- |
| **ID** | 033 |
| **Status** | completed |
| **Model** | Gemini 4 Hydrogen, `thinking_level: medium` (`gemini-4-hydrogen`) |
| **Role** | Frontend Implementer (`implementation`) |
| **Spec** | [spec.md](spec.md) |

> **Specification Note:**  
> Fifth ticket of the end-to-end auth refactor (depends on `032`). Removes the leaky `AuthResponse` and `hasActiveSession()` abstractions from `IAuthAPI`, aligns `MockAuthAPI` and `WailsAuthAdapter` to identical throw/return semantics (including password verification and `is_active` suspension checks in Mock), wires reactive `auth:expired` handling through `IAuthAPI.onSessionExpired`, adds `router.beforeEach` session hydration + route guards on `/studio`, `/settings`, `/arena/history`, and `/auth/*`, aligns `SignUpView` validation rules with backend, and updates `ProfileView` and `ArenaView` auth UX.

## Scope

May create or edit:
- `app/frontend/src/api/contract/auth.contract.ts`
- `app/frontend/src/api/live/wails/auth.wails.ts`
- `app/frontend/src/api/mock/mock.adapter.ts`
- `app/frontend/src/api/mock/fixtures/user.fixture.ts`
- `app/frontend/src/composables/useAuth.ts`
- `app/frontend/src/router/index.ts`
- `app/frontend/src/components/TopBar.vue`
- `app/frontend/src/views/auth/SignInView.vue`
- `app/frontend/src/views/auth/SignUpView.vue`
- `app/frontend/src/views/ProfileView.vue`
- `app/frontend/src/views/ArenaView.vue`
- `app/frontend/src/components/arena/CreateMatchDialog.vue`
- `.routing/API.md` (document updated `IAuthAPI` returning `UserProfile`, `onSessionExpired`, and removal of `hasActiveSession`/`AuthResponse` from frontend contract)
- `.tickets/033-align-frontend-auth-and-guards/changelog.md` (create)

Must not touch:
- `SCHEMA_FINAL.sql`, `sql/`, `pkg/`, `core-service/`, `realtime-service/`, `image-service/`, `app/*.go`
- `app/frontend/uno.config.ts` (record any shortcut proposals in `changelog.md`)
- Any file outside the above scope
- **No role (Router, Investigator, Implementer, Translator, or Reviewer) may commit, push, or merge without explicit user authorization.** Do not run `git commit`, `git push`, or `git merge`.

## Inputs

- [.tickets/028-investigate-auth-workflow/report.md](../028-investigate-auth-workflow/report.md): Findings `FE-01`..`FE-06`.
- [.routing/CONTEXT.md](../../.routing/CONTEXT.md) and [.routing/DESIGN.md](../../.routing/DESIGN.md).
- [.routing/API.md](../../.routing/API.md): Adapter pattern & three synchronization pillars.
- [app/frontend/src/api/contract/auth.contract.ts](../../app/frontend/src/api/contract/auth.contract.ts) and [app/frontend/src/api/live/wails/auth.wails.ts](../../app/frontend/src/api/live/wails/auth.wails.ts).
- [app/frontend/src/api/mock/mock.adapter.ts](../../app/frontend/src/api/mock/mock.adapter.ts) and [app/frontend/src/api/mock/fixtures/user.fixture.ts](../../app/frontend/src/api/mock/fixtures/user.fixture.ts).
- [app/frontend/src/composables/useAuth.ts](../../app/frontend/src/composables/useAuth.ts) and [app/frontend/src/router/index.ts](../../app/frontend/src/router/index.ts).

## Task

1. **Simplify `IAuthAPI` ([auth.contract.ts](../../app/frontend/src/api/contract/auth.contract.ts)) & [auth.wails.ts](../../app/frontend/src/api/live/wails/auth.wails.ts):**
   - Remove `AuthResponse` import/return type and `hasActiveSession()` from `IAuthAPI`.
   - Define:
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
   - In `WailsAuthAdapter`:
     - Remove `hasActiveSession()` and `toAuthResponse()`.
     - `signIn` and `signUp`: await `Wails.SignIn(req)` / `Wails.SignUp(req)`, throw `new Error(...)` if null/error.
     - `getMe`: return `await Wails.GetMe()` (`UserProfile | null`).
     - `onSessionExpired(callback)`: subscribe via `Events.On("auth:expired", ...)` from `@wailsio/runtime` when `isWailsRuntime()` is true (returning the unsubscribe function); otherwise return `() => {}`.

2. **Align `MockAuthAPI` ([mock.adapter.ts](../../app/frontend/src/api/mock/mock.adapter.ts) & [user.fixture.ts](../../app/frontend/src/api/mock/fixtures/user.fixture.ts)):**
   - Remove `MOCK_TOKEN`, `AuthResponse`, and `hasActiveSession()`.
   - Track mock account credentials & operational flags (`password` and `is_active`) alongside `UserProfile` (e.g., export `MOCK_DEMO_PASSWORD = "password123"` and a suspended test user fixture like `suspended@palabrario.app` with `is_active: false` in `user.fixture.ts`, keeping `mockUsers: UserProfile[]` clean for existing fixtures).
   - In `MockAuthAPI.signIn`: validate non-empty email & password (`throw new Error("Email and password are required.")`), verify matching stored password (`throw new Error("Invalid email or password.")`), and check `is_active` (`throw new Error("This account has been suspended by an administrator.")`). Never return `null`.
   - In `MockAuthAPI.signUp`: apply identical normalization & validation as `core-service` (`user_name` after stripping optional `@` must match `/^[a-z0-9][a-z0-9_-]{2,49}$/`, email regex, password 8..72 UTF-8 bytes, `name` trimmed `2..100` chars after falling back to `userName`, allowed dialects `MX|ES|AR|CO|CL|CU|CAM|PAN`). Distinguish duplicate email (`"An account with that email already exists."`) from duplicate username (`"That username is already taken."`). Store the new account's password + `is_active: true`.
   - In `MockAuthAPI.getPublicProfileByUsername`: return `null` (`404`) when the target account has `is_active === false`.
   - Implement `onSessionExpired(callback)` in `MockAuthAPI`.

3. **Upgrade `useAuth` ([useAuth.ts](../../app/frontend/src/composables/useAuth.ts)):**
   - Update `loadSession()` to call `await api.auth.getMe()` directly (1 round-trip instead of `hasActiveSession()` + `getMe()`).
   - Register `api.auth.onSessionExpired((reason) => { ... })` once on module initialization to set `currentUser.value = null`, store `sessionNotice` (reactive ref cleared on next sign-in or dismiss), and trigger navigation if currently on a protected route.
   - Expose `sessionNotice` from `useAuth()` so `SignInView.vue` can display a banner when redirected after session expiry or suspension.

4. **Router Guards & TopBar ([router/index.ts](../../app/frontend/src/router/index.ts), [TopBar.vue](../../app/frontend/src/components/TopBar.vue)):**
   - Remove the `onMounted(() => void init())` call from `TopBar.vue` (or keep it harmlessly idempotent, since `beforeEach` will run `await init()` for *every* route, including `/auth/*`).
   - Add `meta: { requiresAuth: true }` to `studio`, `settings`, and `arena-history`.
   - Add `meta: { guestOnly: true }` to `sign-in` and `sign-up`.
   - Add `router.beforeEach(async (to) => { ... })`:
     - Call `await init()`.
     - If `to.meta.requiresAuth && !isAuthenticated.value`, return `{ name: "sign-in", query: { redirect: to.fullPath } }`.
     - If `to.meta.guestOnly && isAuthenticated.value`, return `{ name: "dictionary" }`.
   - Subscribe `onSessionExpired` handler (or watch `currentUser` in `AppLayout` / router helper) so when a live `auth:expired` event clears `currentUser` while on a `meta.requiresAuth` route, it redirects to `/auth/sign-in` immediately.

5. **Auth & Profile/Arena Views (`SignInView.vue`, `SignUpView.vue`, `ProfileView.vue`, `ArenaView.vue`, `CreateMatchDialog.vue`):**
   - In `SignInView.vue`: display `sessionNotice` alert if present (using existing `border border-main bg-white px-4 py-3 text-sm font-medium text-main` style), and after `signIn` honor safe internal `route.query.redirect` (must start with `/` and not `//`) or default to `/dictionary`.
   - In `SignUpView.vue`: tighten `isStep1Valid` to match backend (`^[a-zA-Z0-9][a-zA-Z0-9_-]{2,49}$` after stripping optional leading `@`, RFC email check `/^[^\s@]+@[^\s@]+\.[^\s@]+$/`, UTF-8 password byte length `8..72` via `TextEncoder`, and if `showFullName` is enabled, `form.name.trim().length` in `2..100` or empty).
   - In `ProfileView.vue`:
     - Watch both `() => route.params.username` and `() => currentUser.value?.id_user` so signing in/out on `/u/:username` refreshes `is_owner`.
     - Distinguish `404` / not found (`"Profile not found."`) from network/server errors (`"Unable to load profile. Please try again."` with retry).
   - In `ArenaView.vue` & `CreateMatchDialog.vue`:
     - When `Create Room` is clicked by an anonymous user (`!isAuthenticated`), navigate to `/auth/sign-in?redirect=/arena` (or show an inline monochrome sign-in prompt) instead of opening `CreateMatchDialog` to hit a backend `401`.
   - Update [.routing/API.md](../../.routing/API.md) to reflect the updated `IAuthAPI` contract and `onSessionExpired` lifecycle hook.
   - Run `wails3 task check` (in `app/`) and verify exit code `0`.

## Verification

Run every check in [spec.md](spec.md).

## Change log

Write `.tickets/033-align-frontend-auth-and-guards/changelog.md` recording all changed files, `IAuthAPI` changes, mock fixture accounts (email/password for manual testing), and `wails3 task check` output.
