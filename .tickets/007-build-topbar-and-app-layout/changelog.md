# 007 — Change log

## Files
Created:
- `app/frontend/src/components/TopBar.vue` (brand link, Dictionary/Arena nav with `active-class="font-bold text-main"`, `Search` icon button, `anonymous` / `logged-in` state machine with mock user)
- `app/frontend/src/layouts/AppLayout.vue` (`TopBar` + `RouterView`)
- `app/frontend/src/views/DictionaryView.vue`, `app/frontend/src/views/ArenaView.vue` (placeholders)

Edited:
- `app/frontend/src/router/index.ts` (`/` parent with `AppLayout` and children `home`, `dictionary`, `arena`; `/auth` with `AuthLayout` kept as a sibling)
- `app/frontend/src/views/HomeView.vue` (`<main>`: `min-h-screen` -> `flex-1`)

## Notes
- `authState` is a local `ref` fixed to `"anonymous"`; nothing toggles it yet, so the logged-in branch is only reachable by editing the default.
- `uno.config.ts` had `light`, `ground`, and `button-icon` shortcuts added to support TopBar search button and monochrome design tokens.
- `TopBar` uses `sticky top-10` to sit below the 40px `WindowBar`.
- Refined header bottom border to `border-neutral-100` (`border-b border-solid border-neutral-100`) matching `WindowBar`.

## Validation
- `cd app && wails3 task check`: exit 0.
- Style greps (`rounded|shadow-|hover:|focus:|active:|transition|animate-` on touched files; legacy tokens and `@palabrario/ui` on `src`): no matches.
- `app/frontend/bindings/` restored with `git checkout`.

## uno.config.ts proposals
None.

---

# Revision 1

## Files
Created:
- `app/frontend/src/components/TOP-BAR.md` (grid alignment, surface hierarchy, authState x viewport matrix)
- `app/frontend/src/views/StudioView.vue`, `app/frontend/src/views/SettingsView.vue` (placeholders)

Edited:
- `app/frontend/src/components/TopBar.vue`: rewritten to `grid grid-cols-12` (left `col-span-6 md:col-span-4`, center `hidden md:flex md:col-span-4 justify-center`, right `col-span-6 md:col-span-4 justify-end`). The user's header, `max-w-6xl` container and `bg-light` miniature edits were kept. `Studio` shows in the desktop nav only when logged-in. Mobile right side is `Search` + `Menu`; the `Menu` trigger opens a full-screen Reka UI popover (`fixed inset-x-0 top-10 bottom-0 z-50 bg-white`) with an `h-20` header and `PopoverClose` (`X`). Its content depends on `authState`. Desktop: anonymous gets `Sign In`/`Sign Up`; logged-in gets a miniature that triggers a user popover (name, `@username`, `Settings`, `Sign Out`). Links close the popovers via `v-model:open`, and `Sign Out` sets `authState` to `"anonymous"`.
- `app/frontend/src/router/index.ts`: added `studio` and `settings` children under `AppLayout`.

## Notes
- `authState` is still a fixed local ref (`"anonymous"`), so the logged-in UI is only visible by changing the default.
- Reka's `PopoverContent` is positioned by a transformed wrapper, so the `fixed` full-screen mobile content may not sit flush at `top-10` and could be offset. I did not check this in a browser; verify it visually.
- Only `Search`, `Menu`, `X` are imported from `lucide-vue-next`.

## Validation
- `cd app && wails3 task check`: exit 0.
- Style greps (touched files) and legacy-token / `@palabrario/ui` greps: no matches.
- `app/frontend/bindings/` restored; `uno.config.ts` and `index.html` untouched by me.

## Revision 1 fix: mobile menu not full screen
- Cause: Reka `PopoverContent` is wrapped in a transformed floating-ui element, so a `fixed inset-x-0 top-10 bottom-0` child is positioned relative to that wrapper instead of the viewport (confirmed by the user's screenshot).
- Fix: the mobile menu now uses Reka `DialogRoot`/`DialogTrigger`/`DialogPortal`/`DialogContent`/`DialogClose` (plus an `sr-only` `DialogTitle`), which renders directly in the portal with no positioner. Classes and content are unchanged. This deviates from the ticket's "Popover" wording; the desktop user menu is still a Popover. `TOP-BAR.md` updated.

## Review: split TopBar
- `TopBarMobileMenu.vue` (new): the `< md` Reka Dialog menu, including its trigger. Owns its open state.
- `TopBarDesktopStatus.vue` (new): the `>= md` auth controls and user Popover. Owns its open state.
- `app/frontend/src/types/user.ts` (new): shared `AuthState` and `TopBarUser` types.
- Both components take `authState` and `user` props and emit `sign-out`; `TopBar.vue` keeps `authState`/`user` and handles the event. Markup and classes are unchanged.
