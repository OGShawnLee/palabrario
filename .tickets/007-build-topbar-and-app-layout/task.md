# 007 — Build the TopBar component, universal AppLayout, and Dictionary/Arena route placeholders

| Field | Value |
| :--- | :--- |
| **ID** | 007 |
| **Status** | completed |
| **Model** | Claude Sonnet 5.5, `effort: medium` (`claude-sonnet-5-5`) |
| **Role** | Implementer |
| **Spec** | [spec.md](spec.md) |

## Goal

Create a `TopBar.vue` navigation component and a universal `AppLayout.vue` layout that wraps all non-auth routes (`/`, `/dictionary`, `/arena`, `/studio`, `/settings`), while `/auth/sign-in` and `/auth/sign-up` remain wrapped only in `AuthLayout.vue`. Add placeholder views for `/dictionary`, `/arena`, `/studio`, and `/settings`.

## Scope

May create or edit:
- `app/frontend/src/components/TopBar.vue` (create / revise)
- `app/frontend/src/components/TopBarMobileMenu.vue` (create: mobile full-screen Reka UI Dialog menu)
- `app/frontend/src/components/TopBarDesktopStatus.vue` (create: desktop auth buttons & user status Reka UI Popover)
- `app/frontend/src/types/user.ts` (create: shared `AuthState` and `TopBarUser` types)
- `app/frontend/src/components/TOP-BAR.md` (create in Revision 1: architecture & state machine documentation)
- `app/frontend/src/layouts/AppLayout.vue` (create)
- `app/frontend/src/views/DictionaryView.vue` (create placeholder)
- `app/frontend/src/views/ArenaView.vue` (create placeholder)
- `app/frontend/src/views/StudioView.vue` (create placeholder in Revision 1)
- `app/frontend/src/views/SettingsView.vue` (create placeholder in Revision 1)
- `app/frontend/src/router/index.ts` (edit)
- `app/frontend/src/views/HomeView.vue` (edit `<main>` height class from `min-h-screen` to `flex-1` so it fits cleanly below `WindowBar` and `TopBar`)
- `.tickets/007-build-topbar-and-app-layout/changelog.md` (create / update)

Must not touch:
- `app/frontend/uno.config.ts` (already updated by the user with `ground: "#FAFAFA"`, `light: "#F5F5F5"`, `button-side: "button bg-ground text-main"`, and `button-icon: "h-12 min-w-12 w-12 bg-light text-black disabled:opacity-50"`)
- `app/frontend/index.html`, `app/frontend/src/components/WindowBar.vue`, `app/frontend/src/App.vue`, `app/frontend/src/main.ts`
- `app/frontend/src/layouts/AuthLayout.vue`, `app/frontend/src/views/auth/SignInView.vue`, `app/frontend/src/views/auth/SignUpView.vue`
- `app/frontend/package.json` or root `package.json`
- `app/*.go`, `app/go.mod`, `app/Taskfile.yml`, or root `Taskfile.yml`
- `local-app/`, `packages/`
- `.context/`, `.routing/`, and other tickets
- `core-service/`, `realtime-service/`, `image-service/`, `pkg/`, and `sql/`

## Inputs

- [CONTEXT.md](../../.routing/CONTEXT.md): stack, mobile-first rule, and ground rules.
- [DESIGN.md](../../.routing/DESIGN.md): monochrome tokens (`main`, `side`, `dull`, `ground`, `light`, `white`), global `index.html` root styling, and button shortcuts (`button`, `button-main`, `button-side`, `button-icon`).
- [app/frontend/uno.config.ts](../../app/frontend/uno.config.ts): exact color aliases and shortcut definitions.
- [app/frontend/index.html](../../app/frontend/index.html): global root `<body>` classes (`bg-ground text-side font-sans antialiased min-h-screen`).
- [app/frontend/src/components/TopBar.vue](../../app/frontend/src/components/TopBar.vue): includes the user's manual post-007 edits (`bg-ground/50`, `border-light`, `max-w-6xl` inner container, and simplified 48x48 user miniature).
- [app/frontend/src/components/WindowBar.vue](../../app/frontend/src/components/WindowBar.vue): reference for `grid-cols-12` center alignment (`col-span-4` left, `col-span-4` center, `col-span-4` right).
- [app/frontend/src/router/index.ts](../../app/frontend/src/router/index.ts): current router.

## Initial Task (Completed)

1. Created initial `TopBar.vue`, `AppLayout.vue`, `DictionaryView.vue`, `ArenaView.vue`, and updated `router/index.ts` and `HomeView.vue`.

---

## Revision 1 — Responsive Mobile/Desktop TopBar, Center Alignment, Reka UI Popovers & `TOP-BAR.md`

### Context: User Manual Edits Already Applied

Before starting this revision, note the changes the user already made directly in the codebase (do not revert these):
1. **`app/frontend/index.html`:** Global root styling now lives on `<body class="bg-ground text-side font-sans antialiased min-h-screen">`. `App.vue` and `AppLayout.vue` inherit `bg-ground` and `font-sans` from `index.html`.
2. **`app/frontend/uno.config.ts` & `.routing/DESIGN.md`:**
   - `ground: "#FAFAFA"` (page canvas, `WindowBar`, `TopBar` translucent base `bg-ground/50`, `button-side` background).
   - `light: "#F5F5F5"` (`WindowBar`, `TopBar`, and Popover borders `border-light`, plus `button-icon` background `bg-light`).
   - `white: "#FFFFFF"` (`bg-white` is reserved for elevated foreground surfaces like Popovers so they stand out cleanly against `bg-ground`).
3. **`app/frontend/src/components/TopBar.vue`:**
   - Outer `<header>` uses `h-20 border-b border-solid border-light bg-ground/50 backdrop-filter backdrop-blur-sm sticky top-10 z-40`.
   - Inner container uses `max-w-6xl h-full w-full mx-auto px-4`.
   - Desktop `logged-in` state in the bar was simplified to just the clickable 48x48 user miniature (`h-12 w-12 bg-light flex items-center justify-center text-sm font-bold text-main shrink-0 overflow-hidden`).

### Revision Tasks

1. **Fix Desktop Navigation Centering in `TopBar.vue` (`max-w-6xl` + `grid-cols-12`).**
   - In the current `flex justify-between` layout, the wider right-side controls push `Dictionary` and `Arena` to the left of the centered `WindowBar` title.
   - Change the inner `max-w-6xl h-full w-full mx-auto px-4` wrapper to use `grid grid-cols-12 items-center gap-4`:
     - **Left column (`col-span-6 md:col-span-4 flex items-center`):** `<RouterLink to="/" class="text-2xl font-black tracking-tight text-main">Palabrario</RouterLink>`
     - **Center column (`hidden md:flex md:col-span-4 items-center justify-center gap-6`):**
       - `<RouterLink to="/dictionary" class="text-sm font-medium text-side" active-class="font-bold text-main">Dictionary</RouterLink>`
       - `<RouterLink to="/arena" class="text-sm font-medium text-side" active-class="font-bold text-main">Arena</RouterLink>`
       - `<RouterLink v-if="authState === 'logged-in'" to="/studio" class="text-sm font-medium text-side" active-class="font-bold text-main">Studio</RouterLink>` (rendered in desktop nav **only** when `authState === "logged-in"`).
     - **Right column (`col-span-6 md:col-span-4 flex items-center justify-end gap-3`):**
       - Always starts with the Search icon button:
         ```vue
         <button type="button" class="button-icon inline-flex items-center justify-center" aria-label="Search">
           <Search class="w-5 h-5" />
         </button>
         ```
       - **On Mobile (`< md`):** Hide the desktop auth controls (`hidden md:flex`) and show the Mobile Menu trigger (`md:hidden`) right next to the Search button.
       - **On Desktop (`>= md`):** Hide the Mobile Menu trigger and show the desktop auth state machine (`anonymous` buttons or `logged-in` User Status Popover).

2. **Implement the Mobile Full-Screen Reka UI Popover (`< md`).**
   - Use Reka UI's `PopoverRoot`, `PopoverTrigger`, `PopoverPortal`, `PopoverContent`, and `PopoverClose` from `"reka-ui"`, and import `Search`, `Menu`, and `X` from `"lucide-vue-next"`.
   - **Mobile Trigger (`md:hidden`):**
     - `<PopoverTrigger as-child>` wrapping `<button type="button" class="button-icon inline-flex md:hidden items-center justify-center" aria-label="Open menu"><Menu class="w-5 h-5" /></button>`.
   - **Mobile Popover Content (`PopoverPortal` -> `PopoverContent`):**
     - Position as a full-screen overlay covering `TopBar` and the viewport below the `h-10` `WindowBar`:
       - Classes on `PopoverContent`: `fixed inset-x-0 top-10 bottom-0 z-50 w-screen bg-white flex flex-col md:hidden` (use `:side-offset="0"` or fixed positioning so it sits flush at `top-10` right below `WindowBar`).
     - **Popover Top Header (`h-20` matching `TopBar`):**
       - `<div class="h-20 w-full border-b border-solid border-light px-4 flex items-center justify-between shrink-0">`
       - Left: `<RouterLink to="/" class="text-2xl font-black tracking-tight text-main">Palabrario</RouterLink>` (clicking it or any nav link should also close the mobile popover, e.g. via `v-model:open="mobileMenuOpen"`).
       - Right: `<PopoverClose as-child><button type="button" class="button-icon inline-flex items-center justify-center" aria-label="Close menu"><X class="w-5 h-5" /></button></PopoverClose>`
     - **Popover Body (`flex-1 p-4 flex flex-col gap-6 overflow-y-auto`):**
       - Driven by `authState`:
       - **When `authState === "logged-in"`:**
         1. Vertical navigation links (`flex flex-col gap-4`):
            - `<RouterLink to="/dictionary" ...>Dictionary</RouterLink>`
            - `<RouterLink to="/arena" ...>Arena</RouterLink>`
            - `<RouterLink to="/studio" ...>Studio</RouterLink>`
            - `<RouterLink to="/settings" ...>Settings</RouterLink>`
         2. User identity wrapper (`flex items-center gap-3 pt-4 border-t border-solid border-light`):
            - User miniature (`h-12 w-12 bg-light flex items-center justify-center text-sm font-bold text-main shrink-0 overflow-hidden`): `<img v-if="user.avatar_url" ... />` or `<span>{{ user.name.charAt(0) }}</span>`.
            - Vertical user text stack (`flex flex-col leading-tight`):
              - `<span class="text-sm font-bold text-main">{{ user.name }}</span>`
              - `<span class="text-xs text-dull">@{{ user.user_name }}</span>`
         3. Sign Out link/button (`text-left text-sm font-medium text-side`): `Sign Out` (clicking sets `authState = "anonymous"` and closes the popover, or acts as a clean action button/link).
       - **When `authState === "anonymous"`:**
         1. Vertical navigation links (`flex flex-col gap-4`):
            - `<RouterLink to="/dictionary" ...>Dictionary</RouterLink>`
            - `<RouterLink to="/arena" ...>Arena</RouterLink>`
         2. Auth action buttons (`flex flex-col gap-3 pt-4 border-t border-solid border-light`):
            - `<RouterLink to="/auth/sign-in" class="button-main w-full">Sign In</RouterLink>`
            - `<RouterLink to="/auth/sign-up" class="button-side w-full">Sign Up</RouterLink>`

3. **Implement the Desktop User Status Reka UI Popover (`>= md`, when `authState === "logged-in"`).**
   - On desktop (`hidden md:flex items-center gap-3`):
     - **When `authState === "anonymous"`:**
       - `<RouterLink to="/auth/sign-in" class="button-main">Sign In</RouterLink>`
       - `<RouterLink to="/auth/sign-up" class="button-side">Sign Up</RouterLink>`
     - **When `authState === "logged-in"`:**
       - Render a Reka UI `PopoverRoot`:
         - `<PopoverTrigger as-child>` wrapping a `<button type="button" aria-label="User menu" class="h-12 w-12 bg-light flex items-center justify-center text-sm font-bold text-main shrink-0 overflow-hidden">` containing the `<img>` or `{{ user.name.charAt(0) }}`.
         - `<PopoverPortal>` -> `<PopoverContent align="end" :side-offset="8" class="z-50 min-w-48 p-4 bg-white border border-solid border-light flex flex-col gap-3">`:
           1. Vertical user info (`flex flex-col leading-tight`):
              - `<span class="text-sm font-bold text-main">{{ user.name }}</span>`
              - `<span class="text-xs text-dull">@{{ user.user_name }}</span>`
           2. Divider: `<div class="border-t border-solid border-light" />`
           3. `<RouterLink to="/settings" class="text-sm font-medium text-side" active-class="font-bold text-main">Settings</RouterLink>`
           4. `<button type="button" class="text-left text-sm font-medium text-side">Sign Out</button>`

4. **Create `StudioView.vue` and `SettingsView.vue` Placeholders & Register in `router/index.ts`.**
   - Create `app/frontend/src/views/StudioView.vue`:
     ```vue
     <template>
       <main class="flex-1 flex items-center justify-center p-4">
         <h1 class="text-2xl font-bold text-main">Studio</h1>
       </main>
     </template>
     ```
   - Create `app/frontend/src/views/SettingsView.vue`:
     ```vue
     <template>
       <main class="flex-1 flex items-center justify-center p-4">
         <h1 class="text-2xl font-bold text-main">Settings</h1>
       </main>
     </template>
     ```
   - Register both under `AppLayout` children in `app/frontend/src/router/index.ts`:
     - `{ path: "studio", name: "studio", component: () => import("@/views/StudioView.vue") }`
     - `{ path: "settings", name: "settings", component: () => import("@/views/SettingsView.vue") }`

5. **Create `app/frontend/src/components/TOP-BAR.md`.**
   - Document:
     - The `max-w-6xl` + `grid-cols-12` layout and how `col-span-4` on desktop keeps the center navigation aligned with `WindowBar.vue`.
     - The surface hierarchy (`bg-ground/50` bar over `bg-ground` page canvas, with `bg-white` Popovers bordered by `border-light`).
     - The complete `authState` (`"anonymous"` vs `"logged-in"`) × Viewport (`< md` Mobile vs `>= md` Desktop) state machine matrix, including the Mobile Full-Screen Popover and Desktop User Status Popover.

6. **Validate and Update `.tickets/007-build-topbar-and-app-layout/changelog.md`.**
   - Run `cd app && wails3 task check` (must exit 0).
   - Restore `app/frontend/bindings/` if modified by `wails3 task check`.
   - Update `changelog.md` with a Revision 1 section.

## Rules

- Mobile-first styling: base classes target mobile; `sm:` and `md:` refine larger screens.
- Use only the color tokens `main`, `side`, `dull`, `ground`, `light`, `white`, and `black`, and the shortcuts `button`, `button-main`, `button-side`, `button-icon`, and `field`.
- Do not use `rounded-*` (except `rounded-none` via shortcuts), `shadow-*`, `hover:`, `focus:`, `active:`, `transition-*`, or `animate-*` classes in any file touched by this ticket.
- Only `Search`, `Menu`, and `X` from `lucide-vue-next` may be imported in `TopBar.vue`. Use `reka-ui` for Popovers. Do not import anything from `@palabrario/ui`.
- Do not modify `app/frontend/uno.config.ts` or `app/frontend/index.html`.
