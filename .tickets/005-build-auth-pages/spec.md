# 005 — Spec checklist

Run by the router before accepting the ticket. All items must pass.

## Files
- [x] These files exist: `app/frontend/src/layouts/AuthLayout.vue`, `app/frontend/src/components/AuthHeader.vue`, `app/frontend/src/components/FormField.vue`, `app/frontend/src/data/regions.ts`.
- [x] `app/frontend/src/views/auth/SignInView.vue` and `SignUpView.vue` contain interactive forms, not the placeholder `<h1>` only.
- [x] `.tickets/005-build-auth-pages/changelog.md` exists and covers every item in the ticket's Change log section.

## Router & Layout
- [x] `app/frontend/src/router/index.ts` uses `createWebHashHistory`.
- [x] The router has a parent route `/auth` with component `AuthLayout` and two children, `sign-in` and `sign-up`.
- [x] The route names `home`, `sign-in`, and `sign-up` exist, and `HomeView.vue` is unchanged and not wrapped in `AuthLayout`.
- [x] `AuthLayout.vue` uses a mobile-first two-column grid (`md:grid-cols-2`) with the left editorial panel (`hidden md:flex`) containing the `Palabrario` link to `/`, the headline `Real Spanish as it is spoken across every region.`, and the 4 benefit statements, and the right column rendering `<RouterView />` (plus a mobile `md:hidden` `Palabrario` link to `/`).

## Content & Interactivity
- [x] `regions.ts` lists 8 regions with these codes in this order: `MX`, `ES`, `AR`, `CO`, `CL`, `CU`, `CAM`, `PAN`.
- [x] `SignInView.vue` has reactive `email` and `password` bindings (`v-model`), computed `isFormValid`, and two inputs (`id="email"` and `id="password"`). The password input has **no** `minlength` attribute, and the `Sign In` submit button (`button-main w-full`) is `:disabled="!isFormValid"`.
- [x] `SignUpView.vue` implements the 2-step wizard (`currentStep` `1 | 2`) and optional full name toggle (`showFullName` with `+ Add full name` / `Remove`).
- [x] Step 1 of `SignUpView.vue` has inputs with ids `user_name` (`minlength="3"`), `name` (when `showFullName` is true), `email`, and `password` (`minlength="8"`), plus a `Continue to region` button (`button-main w-full`, `:disabled="!isStep1Valid"`).
- [x] Step 2 of `SignUpView.vue` renders the regions with `v-for` over `REGIONS`, each button has `:aria-pressed`, the selected one uses `button-main` while others use `button-side`, and includes `Create Account` (`button-main w-full`) and `Back to credentials` (`button-side w-full`).
- [x] Neither form calls an API, a composable, or `router.push` / `router.replace`. `grep -rnE "api|useAuth|router\.(push|replace)" app/frontend/src/views/auth app/frontend/src/components app/frontend/src/layouts` returns nothing.

## Style rules
- [x] `grep -rnE "rounded|shadow-|hover:|focus:|active:|transition|animate-" app/frontend/src/views/auth app/frontend/src/components app/frontend/src/layouts` returns nothing.
- [x] `grep -rnE "side-dark|side-light|text-brand|text-ink|bg-surface|bg-paper|input-editorial|btn-primary|btn-secondary" app/frontend/src` returns nothing.
- [x] `grep -rnE "@palabrario/ui|lucide-vue-next" app/frontend/src` returns nothing.
- [x] `grep -rnP "[áéíóúñ¿¡]" app/frontend/src/views/auth app/frontend/src/components app/frontend/src/layouts app/frontend/src/data` returns nothing (all text is English).
- [x] `app/frontend/uno.config.ts` is unchanged.

## Build and Scope
- [x] `cd app && wails3 task check` exits 0.
- [x] `git status --short` shows changes only in `app/frontend/src/` (the files listed in Scope) and `.tickets/005-build-auth-pages/`.
- [x] No change in `app/frontend/views/HomeView.vue`, `local-app/`, `packages/`, `.context/`, `app/*.go`, `core-service/`, `realtime-service/`, `image-service/`, `pkg/`, or `sql/`.
