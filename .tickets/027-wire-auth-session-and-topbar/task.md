# 027 — Wire Auth Session, Mock User Registry, Auth Forms, and TopBar

| Field | Value |
| :--- | :--- |
| **ID** | 027 |
| **Status** | completed |
| **Model** | Gemini 4 Hydrogen, `thinking_level: medium` (`gemini-4-hydrogen`) |
| **Role** | Frontend Implementer |
| **Spec** | [spec.md](spec.md) |

> **Specification Note:**  
> This ticket creates the singleton `useAuth` composable (calling `hasActiveSession()` before `getMe()` once on boot and caching `currentUser` across `signIn`, `signUp`, and `signOut`), upgrades `MockAuthAPI` to persist an in-memory + `localStorage` user registry so newly registered accounts can sign out, sign back in, and open `/u/:username`, wires `SignInView.vue` and `SignUpView.vue` with monochrome inline error banners and loading states, and replaces the hardcoded `@eponiatowska` placeholder in `TopBar.vue` with live reactive session state.

## Scope

May create or edit:
- `app/frontend/src/api/mock/mock.adapter.ts` (upgrade `MockAuthAPI` to maintain a mutable user registry seeded from `mockUsers` across `signIn`, `signUp`, and `getPublicProfileByUsername`)
- `app/frontend/src/composables/useAuth.ts` (new: singleton reactive auth state composable)
- `app/frontend/src/components/TopBar.vue` (wire `authState`, `user`, and `signOut` to `useAuth()`, removing the hardcoded `eponiatowska` placeholder)
- `app/frontend/src/views/auth/SignInView.vue` (wire form submission to `useAuth().signIn`, add monochrome error banner, loading state, and redirect to `/dictionary`)
- `app/frontend/src/views/auth/SignUpView.vue` (wire form submission to `useAuth().signUp`, add monochrome error banner, loading state, and redirect to `/dictionary`)
- `.tickets/027-wire-auth-session-and-topbar/changelog.md` (create)

Must not touch:
- `app/frontend/uno.config.ts` (use existing tokens and shortcuts; record any proposals in changelog)
- `app/frontend/src/api/contract/auth.contract.ts`, `app/frontend/src/api/live/wails/auth.wails.ts`, `app/frontend/src/api/mock/fixtures/user.fixture.ts`
- `app/*.go`, `core-service/`, `realtime-service/`, `image-service/`, or `.context/`
- Any file outside the above scope

## Inputs

- [.routing/CONTEXT.md](../../.routing/CONTEXT.md): Ground rules (English UI, Satoshi font only, no route aliases, Go backend untouched).
- [.routing/DESIGN.md](../../.routing/DESIGN.md): Monochrome palette (`main`, `side`, `dull`, `light`, `ground`, `white`), sharp geometry (`rounded-none`, zero drop shadows), `button-main`, `button-side`.
- [.routing/API.md](../../.routing/API.md): `api.auth` usage rules (never import `@bindings` runtime functions directly in Vue components; import `@bindings/palabrario/pkg/models/index.js` only for types/DTO classes).
- [.routing/services/identity_service.md](../../.routing/services/identity_service.md): `SignInRequest`, `SignUpRequest`, `AuthResponse`, `UserProfile` contracts (`POST /api/v1/auth/sign-in`, `POST /api/v1/auth/sign-up`, `GET /api/v1/auth/me`).
- [app/frontend/src/composables/useSearchModal.ts](../../app/frontend/src/composables/useSearchModal.ts): Module-level singleton composable pattern reference.

## Task

1. **Upgrade `MockAuthAPI` In-Memory User Registry (`app/frontend/src/api/mock/mock.adapter.ts`):**
   - Add storage key `const MOCK_USERS_KEY = "palabrario_mock_users"`.
   - Add helpers `readUsersStorage(): UserProfile[]` and `writeUsersStorage(users: UserProfile[]): void`:
     - In `readUsersStorage()`: read `localStorage.getItem(MOCK_USERS_KEY)`. If present and valid JSON array, map items with `UserProfile.createFrom(item)`; otherwise return `clone(mockUsers)`.
     - Always clone `mockUsers` so the imported fixture array is never mutated.
   - Store `private users: UserProfile[] = readUsersStorage()` on `MockAuthAPI`.
   - In `signIn(req: SignInRequest)`:
     - Look up `candidate.email.toLowerCase() === email` inside `this.users` (not `mockUsers`).
     - If `!user || req.password.length === 0`, return `null`.
     - Persist `clone(user)` to `this.session` and return `new AuthResponse({ token: MOCK_TOKEN, user: clone(user) })`.
   - In `signUp(req: SignUpRequest)`:
     - Normalize `email = req.email.trim().toLowerCase()`, `userName = req.user_name.trim().replace(/^@/, "").toLowerCase()`, `name = req.name.trim() || userName`, and `anchorDialect = req.anchor_dialect.trim().toUpperCase() || "MX"`.
     - If `!email || !req.password || !userName`, throw `new Error("Email, username, and password are required.")`.
     - Check `this.users` for an existing user matching `email` or `user_name` (case-insensitive). If found, throw `new Error("An account with that email or username already exists.")`.
     - Create `new UserProfile({ ... })` with `role: AccountRole.RoleLearner`, push `clone(user)` into `this.users`, call `writeUsersStorage(this.users)`, call `this.persist(clone(user))`, and return `new AuthResponse({ token: MOCK_TOKEN, user: clone(user) })`.
   - In `getPublicProfileByUsername(username: string)`:
     - Look up `candidate.user_name.toLowerCase() === normalized` inside `this.users` (not `mockUsers`) so newly signed-up users immediately have a working `/u/:username` profile page.

2. **Create `useAuth` Singleton Composable (`app/frontend/src/composables/useAuth.ts`):**
   - Import `computed, ref` from `"vue"`, `api` from `"@/api"`, and types `SignInRequest, SignUpRequest, UserProfile` from `"@bindings/palabrario/pkg/models/index.js"`.
   - Declare module-level reactive refs:
     - `const currentUser = ref<UserProfile | null>(null)`
     - `const isInitialized = ref(false)`
     - `const isAuthLoading = ref(false)`
     - `let initPromise: Promise<void> | null = null`
   - Implement `async function init(): Promise<void>`:
     - If `isInitialized.value`, return immediately.
     - If `initPromise`, return `initPromise`.
     - Set `initPromise` to an async IIFE that checks `const active = await api.auth.hasActiveSession()`; if `active`, sets `currentUser.value = await api.auth.getMe()`, otherwise sets `currentUser.value = null`; wraps in `try/catch` (setting `currentUser.value = null` on failure) and in `finally` sets `isInitialized.value = true` and `initPromise = null`.
   - Implement `async function signIn(req: SignInRequest): Promise<UserProfile>`:
     - Set `isAuthLoading.value = true`.
     - In `try`, call `const res = await api.auth.signIn(req)`. If `!res`, throw `new Error("Invalid email or password.")`.
     - Set `currentUser.value = res.user`, `isInitialized.value = true`, and return `res.user`.
     - In `finally`, set `isAuthLoading.value = false`.
   - Implement `async function signUp(req: SignUpRequest): Promise<UserProfile>`:
     - Set `isAuthLoading.value = true`.
     - In `try`, call `const res = await api.auth.signUp(req)`. If `!res`, throw `new Error("Could not create account.")`.
     - Set `currentUser.value = res.user`, `isInitialized.value = true`, and return `res.user`.
     - In `finally`, set `isAuthLoading.value = false`.
   - Implement `async function signOut(): Promise<void>`:
     - Set `isAuthLoading.value = true`.
     - In `try`, call `await api.auth.signOut()`, then set `currentUser.value = null` and `isInitialized.value = true`.
     - In `finally`, set `isAuthLoading.value = false`.
   - Implement `function setUser(user: UserProfile | null): void` updating `currentUser.value = user`.
   - Export `useAuth()` returning `{ currentUser, isAuthenticated: computed(() => currentUser.value !== null), isInitialized, isAuthLoading, init, signIn, signUp, signOut, setUser }`.

3. **Wire `TopBar.vue` (`app/frontend/src/components/TopBar.vue`):**
   - Remove the hardcoded `authState` ref and hardcoded `Elena Poniatowska` (`eponiatowska`) `user` object.
   - Import `computed, onMounted` from `"vue"`, `useRoute, useRouter` from `"vue-router"`, and `useAuth` from `"@/composables/useAuth"`.
   - Call `onMounted(() => { void init() })`.
   - Derive `authState = computed<AuthState>(() => (currentUser.value ? "logged-in" : "anonymous"))`.
   - Derive `topBarUser = computed<TopBarUser>(() => ({ name: currentUser.value?.name || currentUser.value?.user_name || "", user_name: currentUser.value?.user_name || "", avatar_url: currentUser.value?.avatar_url || "" }))`.
   - Update `async function handleSignOut()`: call `await signOut()`, and if `route.path.startsWith("/studio") || route.path.startsWith("/settings")`, navigate to `"/"` via `router.push("/")`.
   - Pass `:auth-state="authState"` and `:user="topBarUser"` to `TopBarMobileMenu` and `TopBarDesktopStatus`, and `@sign-out="handleSignOut"`.

4. **Wire `SignInView.vue` (`app/frontend/src/views/auth/SignInView.vue`):**
   - Import `useRouter` from `"vue-router"` and `useAuth` from `"@/composables/useAuth"`.
   - Add `errorMessage = ref<string | null>(null)`.
   - Replace `@submit.prevent` with `@submit.prevent="onSubmit"`:
     - Guard `if (!isFormValid.value || isAuthLoading.value) return`.
     - Clear `errorMessage.value = null`.
     - Call `await signIn({ email: email.value.trim(), password: password.value })` inside `try`, then `await router.push("/dictionary")`.
     - In `catch (err)`: set `errorMessage.value = err instanceof Error ? err.message : "Unable to sign in."`.
   - Above the submit button inside `<form>`, render a sharp monochrome alert when `errorMessage` is set:
     `<div v-if="errorMessage" role="alert" class="border border-main bg-white px-4 py-3 text-sm font-medium text-main">{{ errorMessage }}</div>`.
   - Disable the submit button when `!isFormValid || isAuthLoading` and display `"Signing in..."` while `isAuthLoading` is true, `"Sign In"` otherwise.

5. **Wire `SignUpView.vue` (`app/frontend/src/views/auth/SignUpView.vue`):**
   - Import `useRouter` from `"vue-router"` and `useAuth` from `"@/composables/useAuth"`.
   - Add `errorMessage = ref<string | null>(null)`.
   - Clear `errorMessage.value = null` when moving between Step 1 and Step 2.
   - Replace `@submit.prevent` on `<form>` with `@submit.prevent="onSubmit"`:
     - If `currentStep.value === 1`, advance `if (isStep1Valid.value) currentStep.value = 2` and return.
     - Guard `if (!isStep1Valid.value || isAuthLoading.value) return`.
     - Clear `errorMessage.value = null`.
     - Compute `const displayName = form.name.trim() || form.user_name.trim()`.
     - Call `await signUp({ email: form.email.trim(), password: form.password, name: displayName, user_name: form.user_name.trim(), anchor_dialect: form.anchor_dialect })` inside `try`, then `await router.push("/dictionary")`.
     - In `catch (err)`: set `errorMessage.value = err instanceof Error ? err.message : "Unable to create account."`.
   - Render the same sharp monochrome `role="alert"` box (`border border-main bg-white px-4 py-3 text-sm font-medium text-main`) on Step 2 when `errorMessage` is non-null.
   - Disable `"Create Account"` and `"Back to Credentials"` while `isAuthLoading` is true, showing `"Creating account..."` on the submit button while loading.

6. **Validation:**
   - Run `wails3 task check` inside `app/` and confirm exit code 0.

7. **Change log:**
   - Document created/modified files and any notes in `.tickets/027-wire-auth-session-and-topbar/changelog.md`.

## Verification

All checks are listed in [spec.md](spec.md).
