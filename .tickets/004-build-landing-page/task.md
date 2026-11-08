# 004 — Build the landing page and set up Vue Router

| Field | Value |
| :--- | :--- |
| **ID** | 004 |
| **Status** | completed |
| **Model** | Claude Sonnet 5.5, `effort: medium` (`claude-sonnet-5-5`) |
| **Role** | Implementer |
| **Spec** | [spec.md](spec.md) |

## Scope

May create or edit:
- `app/frontend/src/router/index.ts` (create)
- `app/frontend/src/views/HomeView.vue` (edit)
- `app/frontend/src/views/auth/SignInView.vue` (create placeholder)
- `app/frontend/src/views/auth/SignUpView.vue` (create placeholder)
- `app/frontend/src/App.vue` (edit)
- `app/frontend/src/main.ts` (edit)
- `.tickets/004-build-landing-page/changelog.md` (create)

Must not touch:
- `app/frontend/uno.config.ts`
- `app/frontend/package.json` or root `package.json`
- `app/*.go`, `app/go.mod`, `app/Taskfile.yml`, or root `Taskfile.yml`
- `local-app/`
- `packages/`
- `.context/`, `.routing/`, and other tickets
- `core-service/`, `realtime-service/`, `image-service/`, `pkg/`, and `sql/`

## Inputs

- [CONTEXT.md](../../.routing/CONTEXT.md): project stack, mobile-first ground rule, and scope rules.
- [DESIGN.md](../../.routing/DESIGN.md): monochrome color tokens (`main`, `side`, `dull`), typography (`Satoshi`), and button shortcuts (`button`, `button-main`, `button-side`).
- [app/frontend/uno.config.ts](../../app/frontend/uno.config.ts): exact theme colors and shortcuts.
- [packages/ui/src/components/ui/HomeHero.vue](../../packages/ui/src/components/ui/HomeHero.vue): reference hero structure (simplified and translated in Step 3 below).

## Task

1. **Create the auth route placeholders.**
   - Create `app/frontend/src/views/auth/SignInView.vue` with a minimal template:
     ```vue
     <template>
       <main class="min-h-screen flex items-center justify-center p-4">
         <h1 class="text-2xl font-bold text-main">Sign In</h1>
       </main>
     </template>
     ```
   - Create `app/frontend/src/views/auth/SignUpView.vue` with a minimal template:
     ```vue
     <template>
       <main class="min-h-screen flex items-center justify-center p-4">
         <h1 class="text-2xl font-bold text-main">Sign Up</h1>
       </main>
     </template>
     ```
   - Do not add `<script>` blocks, forms, or backend calls to either placeholder. They will be built in their own tickets.

2. **Create the router.**
   - Create `app/frontend/src/router/index.ts` using `createRouter` and `createWebHashHistory` from `vue-router`.
   - Define exactly three routes:
     1. `{ path: "/", name: "home", component: HomeView }` (import `HomeView` from `@/views/HomeView.vue`).
     2. `{ path: "/auth/sign-in", name: "sign-in", component: () => import("@/views/auth/SignInView.vue") }`.
     3. `{ path: "/auth/sign-up", name: "sign-up", component: () => import("@/views/auth/SignUpView.vue") }`.
   - Export the router instance as the default export. Do not add navigation guards, composables, or extra routes.

3. **Wire the router into the app.**
   - In `app/frontend/src/main.ts`, import `router` from `./router` and register it with `.use(router)` before `.mount("#app")`. Keep the `virtual:uno.css` and `@unocss/reset/tailwind-compat.css` imports in their exact order at the top of `main.ts`.
   - In `app/frontend/src/App.vue`, replace the direct `HomeView` import and usage with `<RouterView />` imported from `vue-router`.

4. **Rebuild `app/frontend/src/views/HomeView.vue`.**
   - Build a static, mobile-first monochrome hero page adapted from `HomeHero.vue`.
   - Do not add props, emits, search input, region chips, or icons.
   - Structure and classes (mobile-first, scaling up with `sm:` and `md:`):
     - Outer `<main>` container: `min-h-screen w-full flex flex-col items-center justify-center text-center px-4 py-12 sm:px-6 md:py-16`
     - Inner content wrapper `<div>`: `w-full max-w-4xl mx-auto flex flex-col items-center gap-6 md:gap-8`
     - **Eyebrow `<p>`:**
       - Classes: `text-xs font-bold tracking-widest uppercase text-dull select-none`
       - Exact text: `The dictionary of our differences`
     - **Headline `<h1>`:**
       - Classes: `text-4xl sm:text-6xl md:text-7xl font-black tracking-tight text-main leading-none`
       - Three `<span class="block">` lines inside `<h1>` with exact text:
         1. `One shared language.`
         2. `Twenty ways to say`
         3. `the same thing.`
     - **Subtitle `<p>`:**
       - Classes: `text-base sm:text-lg md:text-xl text-side leading-relaxed max-w-2xl mx-auto font-normal`
       - Exact text: `In Mexico it's chamba, in Argentina it's laburo, in Spain it's curro. Explore the living vocabulary that unites and divides five hundred million Spanish speakers.`
     - **Action buttons container `<div>`:**
       - Classes: `flex flex-col sm:flex-row items-center justify-center gap-3 w-full max-w-xs sm:max-w-md mx-auto pt-2`
       - Two `<RouterLink>` elements (imported from `vue-router`):
         1. `<RouterLink to="/auth/sign-in" class="button-main w-full sm:w-auto">Sign In</RouterLink>`
         2. `<RouterLink to="/auth/sign-up" class="button-side w-full sm:w-auto">Sign Up</RouterLink>`

5. **Validate.**
   - Run `cd app && wails3 task check`. It must exit with code 0.
   - If it fails, record the exact error in `changelog.md` and stop.

6. **Write the change log** (see below).

## Rules

- Follow mobile-first styling: base classes target mobile screens; `sm:` and `md:` refine layout for larger screens.
- Use only the color tokens defined in `app/frontend/uno.config.ts` (`main`, `side`, `dull`, `white`, `black`) and the shortcuts `button-main` and `button-side`.
- Do not use `rounded-*` (except `rounded-none` via shortcuts), `shadow-*`, `hover:`, `focus:`, or `active:` classes anywhere in `app/frontend/src`.
- Do not import anything from `@palabrario/ui` or `lucide-vue-next`.
- Do not modify `app/frontend/uno.config.ts`. If you think a new shortcut or token is needed, propose it in `changelog.md`.

## Verification

All checks are listed in [spec.md](spec.md). The router runs them before accepting the ticket, not the worker's self-report.

## Change log

Write `.tickets/004-build-landing-page/changelog.md` with:
- Each file created or edited.
- Output status of `wails3 task check`.
- Any `uno.config.ts` shortcut or token proposals for the user to review.
- Any problems encountered.
