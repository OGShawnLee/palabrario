# 027 — Spec checklist

Run by the router before accepting the ticket. All items must pass.

- [x] `MockAuthAPI` in `app/frontend/src/api/mock/mock.adapter.ts` maintains a mutable `users` list seeded from `mockUsers` (persisted via `palabrario_mock_users`), checks duplicate email/username on `signUp`, and reads from that list in `signIn`, `signUp`, and `getPublicProfileByUsername`.
- [x] `app/frontend/src/composables/useAuth.ts` exists, exposes `{ currentUser, isAuthenticated, isInitialized, isAuthLoading, init, signIn, signUp, signOut, setUser }`, calls `hasActiveSession()` before `getMe()` inside `init()`, and deduplicates concurrent `init()` calls.
- [x] `app/frontend/src/components/TopBar.vue` has zero hardcoded `eponiatowska` or static `logged-in` state, initializes `useAuth()` on mount, and derives `authState`, `user`, and `signOut` from `useAuth()`.
- [x] `app/frontend/src/views/auth/SignInView.vue` submits to `useAuth().signIn(...)`, displays a monochrome `role="alert"` error message on failure, disables submit while `isAuthLoading`, and routes to `/dictionary` on success.
- [x] `app/frontend/src/views/auth/SignUpView.vue` submits to `useAuth().signUp(...)` (defaulting empty optional `name` to `user_name`), displays a monochrome `role="alert"` error message on failure, disables Step 2 buttons while `isAuthLoading`, and routes to `/dictionary` on success.
- [x] Design hygiene: zero `font-mono`, zero `flag_emoji`, zero rounded corners (`rounded-none`), and zero `@bindings` runtime calls inside `.vue` or composable files.
- [x] Typecheck and build pass: `wails3 task check` in `app/` exits with code 0.
- [x] No application code outside declared scope, backend code, or `.context/` modified.
- [x] `.tickets/027-wire-auth-session-and-topbar/changelog.md` exists and is populated.
- [x] Commit (if made): Conventional Commits, English, detailed body.
