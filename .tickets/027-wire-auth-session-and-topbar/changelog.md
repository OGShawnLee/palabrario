# 027 — Change log

## Created
- `app/frontend/src/composables/useAuth.ts`: singleton reactive auth composable exposing `{ currentUser, isAuthenticated, isInitialized, isAuthLoading, init, signIn, signUp, signOut, setUser }`, calling `api.auth.hasActiveSession()` before `api.auth.getMe()` in `init()`, deduplicating concurrent `init()` calls, and persisting session state across `signIn`, `signUp`, and `signOut`.

## Modified
- `app/frontend/src/api/mock/mock.adapter.ts`: upgraded `MockAuthAPI` with a mutable `users` registry seeded from `mockUsers` and persisted under `localStorage` key `palabrario_mock_users`; `signIn`, `signUp`, and `getPublicProfileByUsername` now read from `this.users`; `signUp` normalizes fields, validates required fields, checks duplicate email/username (case-insensitive), and persists the newly created user.
- `app/frontend/src/components/TopBar.vue`: removed hardcoded `eponiatowska` placeholder and static `logged-in` state, initialized `useAuth()` on mount, derived `authState` and `topBarUser` reactively, and redirected from `/studio` or `/settings` to `/` on sign-out.
- `app/frontend/src/views/auth/SignInView.vue`: wired form submission to `useAuth().signIn(...)`, added a monochrome `role="alert"` error banner, disabled submit while `isAuthLoading` (`"Signing in..."`), and routed to `/dictionary` on success.
- `app/frontend/src/views/auth/SignUpView.vue`: wired form submission to `useAuth().signUp(...)` (defaulting empty optional `name` to `user_name`), cleared `errorMessage` across step transitions, added a monochrome `role="alert"` error banner on Step 2, disabled Step 2 actions while `isAuthLoading` (`"Creating account..."`), and routed to `/dictionary` on success.
- `app/frontend/src/composables/useAuth.ts` (refactor): extracted `loadSession()` from `init()` to replace the inline async IIFE, guarded `currentUser` updates inside `loadSession()` with `!isInitialized.value` so an in-flight `init()` cannot overwrite state set by `signIn`, `signUp`, `signOut`, or `setUser`, and updated `setUser()` to mark `isInitialized = true`.

## Notes / proposals
- No `uno.config.ts` changes; existing tokens (`border-main`, `bg-white`, `text-main`, `button-main`, `button-side`) covered the inline error banner and loading states.

## Validation
- `wails3 task check` in `app/`: exited with code 0 (`wails3 generate bindings -clean=true -ts`, `bun run --cwd frontend type-check`, `bun run --cwd frontend build`, and `go build .` all passed).
