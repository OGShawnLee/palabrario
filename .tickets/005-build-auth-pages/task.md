# 005 — Build the two-column auth layout and the interactive sign-in and sign-up pages

| Field | Value |
| :--- | :--- |
| **ID** | 005 |
| **Status** | completed |
| **Model** | Claude Sonnet 5.5, `effort: medium` (`claude-sonnet-5-5`) |
| **Role** | Implementer |
| **Spec** | [spec.md](spec.md) |

## Goal

Replace the two placeholder auth pages (`SignInView.vue` and `SignUpView.vue`, created in Ticket 004) with monochrome, mobile-first forms inside a shared two-column `AuthLayout.vue`. Include minimal Vue `<script setup>` state for form reactivity, step navigation, and region selection, without backend/API calls.

## Decisions already made

- **Existing setup from Ticket 004:** `app/frontend/src/router/index.ts`, `app/frontend/src/views/HomeView.vue`, and the placeholder `SignInView.vue` / `SignUpView.vue` already exist. `HomeView.vue` is already built with the updated hero and must **not** be modified.
- **Two-column `AuthLayout.vue` (mobile-first):**
  - Mobile (`< md`): single-column view with the `Palabrario` wordmark link at the top followed by the form `<RouterView />`.
  - Desktop (`md:` and up): two-column split (`md:grid-cols-2`, sharp 1px border, no rounded corners) where the left column is a monochrome editorial panel (wordmark link to `/`, headline, and 4 translated benefit items — text-only, no icons) and the right column renders `<RouterView />`.
- **Minimal Vue `<script setup>` interactivity:**
  - `SignInView.vue`: reactive `email` and `password` state (`v-model`), computed `isFormValid` (`email.trim().includes("@") && password.length > 0` — sign-in accepts any non-empty password, no `minlength`), and `:disabled="!isFormValid"` on the submit button.
  - `SignUpView.vue`: two-step wizard (`currentStep` `1 | 2`), optional full name toggle (`showFullName` boolean: `+ Add full name` / `Remove`), reactive form fields (`user_name`, `name`, `email`, `password`, `anchor_dialect` defaulting to `"MX"`), and step 1 validation (`user_name.trim().length >= 3 && email.trim().includes("@") && password.length >= 8`) before advancing to step 2.
- **No backend calls or error banners yet:** Neither form calls Wails bindings, composables, or `router.push`. Submitting prevents default (`@submit.prevent`). Error banners will be added when backend wiring is introduced.
- **Headroom for layout & composition:** Exact copy, field IDs, region codes, and design system constraints are specified below; you have headroom to compose the clean spacing, borders, and typography within the monochrome design system (`main`, `side`, `dull`, `white`, `black`, and shortcuts `button`, `button-main`, `button-side`, `field`).

## Scope

May create or edit:
- `app/frontend/src/layouts/AuthLayout.vue` (create)
- `app/frontend/src/components/AuthHeader.vue` (create)
- `app/frontend/src/components/FormField.vue` (create)
- `app/frontend/src/data/regions.ts` (create)
- `app/frontend/src/views/auth/SignInView.vue` (replace placeholder)
- `app/frontend/src/views/auth/SignUpView.vue` (replace placeholder)
- `app/frontend/src/router/index.ts` (edit to wrap auth routes in `AuthLayout`)
- `.tickets/005-build-auth-pages/changelog.md` (create)

Must not touch:
- `app/frontend/uno.config.ts` (propose any new shortcuts or tokens in `changelog.md`)
- `app/frontend/src/views/HomeView.vue`, `app/frontend/src/App.vue`, `app/frontend/src/main.ts`
- `app/frontend/package.json` or root `package.json`
- `app/*.go`, `app/go.mod`, `app/Taskfile.yml`, or root `Taskfile.yml`
- `local-app/`, `packages/`
- `.context/`, `.routing/`, and other tickets
- `core-service/`, `realtime-service/`, `image-service/`, `pkg/`, and `sql/`

## Inputs

- [CONTEXT.md](../../.routing/CONTEXT.md): stack, mobile-first rule, and ground rules.
- [DESIGN.md](../../.routing/DESIGN.md): monochrome tokens (`main`, `side`, `dull`), Satoshi font, and shortcuts (`button`, `button-main`, `button-side`, `field`).
- [app/frontend/uno.config.ts](../../app/frontend/uno.config.ts): exact color aliases and shortcut definitions.
- [app/frontend/src/views/HomeView.vue](../../app/frontend/src/views/HomeView.vue): current landing page for typographic and spacing reference.
- [app/frontend/src/router/index.ts](../../app/frontend/src/router/index.ts): current router from Ticket 004.
- Reference views from `local-app` (for structural inspiration only — do not copy legacy classes or Spanish text):
  - [local-app/frontend/src/layouts/AuthLayout.vue](../../local-app/frontend/src/layouts/AuthLayout.vue)
  - [local-app/frontend/src/views/auth/SignInView.vue](../../local-app/frontend/src/views/auth/SignInView.vue)
  - [local-app/frontend/src/views/auth/SignUpView.vue](../../local-app/frontend/src/views/auth/SignUpView.vue)

## Task

1. **Create `app/frontend/src/data/regions.ts`.**
   - Export `interface Region { code: string; name: string }` and `export const REGIONS: Region[]` with these 8 seed regions in this exact order:
     1. `{ code: "MX", name: "Mexico" }`
     2. `{ code: "ES", name: "Spain" }`
     3. `{ code: "AR", name: "Argentina" }`
     4. `{ code: "CO", name: "Colombia" }`
     5. `{ code: "CL", name: "Chile" }`
     6. `{ code: "CU", name: "Cuba" }`
     7. `{ code: "CAM", name: "Central America" }`
     8. `{ code: "PAN", name: "Pan-Hispanic" }`

2. **Create `app/frontend/src/components/AuthHeader.vue` and `app/frontend/src/components/FormField.vue`.**
   - `AuthHeader.vue`: accepts `title: string` and `subtitle: string` props and renders a clean header (`<h1>` in `text-main` + `<p>` in `text-side`).
   - `FormField.vue`: accepts `id: string`, `label: string`, and optional `hint?: string` props. Renders a `<label :for="id">` (with an optional header slot or right-side action slot if useful for `+ Add full name` / `Remove`), the default `<slot />` for the input, and `<p v-if="hint" class="text-xs text-dull">{{ hint }}</p>`.

3. **Create the two-column `app/frontend/src/layouts/AuthLayout.vue`.**
   - Mobile-first container centered in the viewport (`min-h-screen`), scaling up to a two-column grid (`md:grid-cols-2`, e.g. `max-w-md md:max-w-4xl`) with sharp 1px solid borders (`border border-dull` or `border-black`).
   - **Left column (`<aside>`, hidden on mobile `hidden md:flex`, visible from `md:` up):**
     - Top `<RouterLink to="/">Palabrario</RouterLink>` wordmark.
     - Editorial headline: `Real Spanish as it is spoken across every region.`
     - List of 4 benefit statements (text-only monochrome editorial list, no icons):
       1. `Discover how the same idea is expressed in Mexico, Argentina, or Spain.`
       2. `Review with flashcards timed to return right before you forget them.`
       3. `Challenge your friends to live multiplayer vocabulary matches.`
       4. `Every definition is reviewed by three editors before publication.`
   - **Right column (`<main>`):**
     - Mobile-only `<RouterLink to="/" class="md:hidden ...">Palabrario</RouterLink>` at the top so mobile users always have a link back to `/`.
     - Renders `<RouterView />`.

4. **Update `app/frontend/src/router/index.ts`.**
   - Keep `createWebHashHistory()` and the `/` route (`name: "home"`, `HomeView`).
   - Wrap `/auth/sign-in` (`name: "sign-in"`) and `/auth/sign-up` (`name: "sign-up"`) as children of `{ path: "/auth", component: AuthLayout }` so both render inside `AuthLayout.vue` while URLs remain `/auth/sign-in` and `/auth/sign-up`.

5. **Build `app/frontend/src/views/auth/SignInView.vue`.**
   - Header via `AuthHeader`:
     - Title: `Welcome back`
     - Subtitle: `Pick up your streak and your decks where you left them.`
   - Reactive form state (`email`, `password`) and computed `isFormValid` (`email.trim().includes("@") && password.length > 0`).
   - Fields via `FormField`:
     - Email input: `id="email"`, `v-model`, `type="email"`, `required`, `autocomplete="email"`, `class="field w-full"`.
     - Password input: `id="password"`, `v-model`, `type="password"`, `required`, `autocomplete="current-password"`, `class="field w-full"`. Do **not** set `minlength` on sign-in password.
   - Submit button: `type="submit"`, `class="button-main w-full"`, `:disabled="!isFormValid"`, label `Sign In`.
   - Footer prompt: `New here?` with `<RouterLink to="/auth/sign-up">Create an account</RouterLink>`.

6. **Build `app/frontend/src/views/auth/SignUpView.vue`.**
   - Two-step wizard controlled by `const currentStep = ref<1 | 2>(1)`:
     - Compact step indicator at the top showing `1. Credentials` and `2. Regional variant`, visually distinguishing the active step using monochrome typography/borders (`text-main font-bold` vs `text-dull`).
   - **Step 1 (`currentStep === 1` — Credentials):**
     - Header (`AuthHeader`): title `Create your account`, subtitle `Join in a few minutes.`
     - Username (`FormField` `id="user_name"`, label `Username`, hint `Letters, numbers, or hyphens`): `<input id="user_name" v-model="form.user_name" type="text" required minlength="3" autocomplete="username" class="field w-full" />`.
     - Optional full name toggle:
       - When `showFullName` is `false`, show a button `+ Add full name` next to the Username label (or below it) that sets `showFullName = true`.
       - When `showFullName` is `true`, render `FormField` with `id="name"`, label `Full name (optional)`, a `Remove` button that resets `showFullName = false; form.name = ""`, and `<input id="name" v-model="form.name" type="text" autocomplete="name" class="field w-full" />`.
     - Email (`FormField` `id="email"`, label `Email`): `<input id="email" v-model="form.email" type="email" required autocomplete="email" class="field w-full" />`.
     - Password (`FormField` `id="password"`, label `Password`, hint `At least 8 characters`): `<input id="password" v-model="form.password" type="password" required minlength="8" autocomplete="new-password" class="field w-full" />`.
     - Primary button to advance: `type="button"`, `class="button-main w-full"`, `:disabled="!isStep1Valid"`, label `Continue to region`, clicking sets `currentStep = 2`.
     - Footer prompt: `Already have an account?` with `<RouterLink to="/auth/sign-in">Sign In</RouterLink>`.
   - **Step 2 (`currentStep === 2` — Regional variant):**
     - Header (`AuthHeader`): title `Where is your Spanish from?`, subtitle `We tailor your decks and challenges to your region.`
     - Region selector grid (`grid grid-cols-2 gap-2`): renders a `<button type="button">` for each `region in REGIONS`, with `:class="form.anchor_dialect === region.code ? 'button-main' : 'button-side'"`, `:aria-pressed="form.anchor_dialect === region.code"`, and `@click="form.anchor_dialect = region.code"`.
     - Action buttons:
       - Submit button: `type="submit"`, `class="button-main w-full"`, label `Create Account`.
       - Back button: `type="button"`, `class="button-side w-full"`, label `Back to credentials`, clicking sets `currentStep = 1`.

7. **Validate.**
   - Run `cd app && wails3 task check`. It must exit with code 0.
   - Note: `wails3 task check` may touch `app/frontend/bindings/`. Restore `app/frontend/bindings/` with `git checkout -- app/frontend/bindings/` afterwards so the working tree only contains changes in this ticket's scope.

8. **Write the change log** in `.tickets/005-build-auth-pages/changelog.md`.

## Rules

- Mobile-first styling: base classes target mobile; `sm:` and `md:` refine larger viewports.
- Use only the color tokens `main`, `side`, `dull`, `white`, and `black`, and the shortcuts `button`, `button-main`, `button-side`, and `field`.
- Do not use `rounded-*` (except `rounded-none` via shortcuts), `shadow-*`, `hover:`, `focus:`, `active:`, `transition-*`, or `animate-*` classes anywhere in `app/frontend/src`.
- Do not import anything from `@palabrario/ui` or `lucide-vue-next`.
- Do not call any backend/Wails API, composable, or `router.push` / `router.replace`.
- Do not modify `app/frontend/uno.config.ts`. Record any proposed shortcuts or tokens in `changelog.md`.

## Verification

All checks are in [spec.md](spec.md). The router runs them before accepting the ticket, not the worker's self-report.

## Change log

Write `.tickets/005-build-auth-pages/changelog.md` with:
- Each file created or edited.
- Layout and interactivity choices made.
- Output status of `wails3 task check`.
- Any shortcut or token proposals for `uno.config.ts` for the user to review.
