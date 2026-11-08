# 006 — Build the frameless macOS-style WindowBar and Wails v3 window service

| Field | Value |
| :--- | :--- |
| **ID** | 006 |
| **Status** | completed |
| **Model** | Gemini 3.8 Flash, `thinking_level: high` (`gemini-3.8-flash`, Antigravity provider) |
| **Role** | Implementer |
| **Spec** | [spec.md](spec.md) |

## Goal

Hide the stock OS window titlebar (`Frameless: true` in `app/main.go`), add Wails v3 window control methods in `app/app_window.go` (`OnMinimize`, `OnMaximize`, `OnClose`), create a clean macOS-style `WindowBar.vue` component (`h-10 bg-white border-b border-solid border-dull`) with the 3 traffic-light buttons on the left and a draggable bar area, render it globally in `app/frontend/src/App.vue`, and delete the temporary reference file `app/frontend/src/components/WindowBar.md`.

## Exception to general rules

Unlike earlier frontend-only tickets, this ticket **is explicitly authorized** to edit `app/main.go`, create `app/app_window.go`, and commit the regenerated Wails bindings in `app/frontend/bindings/` (since 3 new Go methods are added to `App`). Also, the 3 traffic-light buttons in `WindowBar.vue` use `rounded-full` (12x12 circles, `w-3 h-3 rounded-full`) with the macOS colors from the reference (`#FF5F56` / `#E0443E`, `#FFBD2E` / `#DEA123`, `#27C93F` / `#1AAB29`).

## Scope

May create, edit, or delete:
- `app/main.go` (edit: add `Frameless: true` to `WebviewWindowOptions` and update `BackgroundColour` to white `application.NewRGB(255, 255, 255)`)
- `app/app_window.go` (create: `OnMinimize`, `OnMaximize`, `OnClose` on `*App`)
- `app/frontend/bindings/` (updated by `wails3 task check` with the 3 new `App` methods)
- `app/frontend/src/components/WindowBar.vue` (create)
- `app/frontend/src/components/WindowBar.md` (delete reference file once implemented)
- `app/frontend/src/App.vue` (edit: render `<WindowBar />` above `<RouterView />`)
- `.tickets/006-build-window-bar/changelog.md` (create)

Must not touch:
- `app/frontend/uno.config.ts`
- `app/app.go` and the other `app/app_*.go` files
- `app/frontend/src/views/`, `app/frontend/src/layouts/`, `app/frontend/src/router/`
- `local-app/`, `packages/`
- `.context/`, `.routing/`, and other tickets
- `core-service/`, `realtime-service/`, `image-service/`, `pkg/`, and `sql/`

## Inputs

- [CONTEXT.md](../../.routing/CONTEXT.md): stack and ground rules.
- [DESIGN.md](../../.routing/DESIGN.md): monochrome tokens (`main`, `side`, `dull`, `light`).
- [app/frontend/src/components/WindowBar.md](../../app/frontend/src/components/WindowBar.md): reference file (read before deleting in Step 4).
- [app/main.go](../../app/main.go): window initialization in Wails v3.
- [app/app.go](../../app/app.go): `App` struct definition.

## Task

1. **Create `app/app_window.go` with Wails v3 runtime methods.**
   - In Wails v3 (`github.com/wailsapp/wails/v3/pkg/application`), window operations go through `application.Get()`:
     ```go
     package main

     import "github.com/wailsapp/wails/v3/pkg/application"

     func (a *App) OnMinimize() {
     	if app := application.Get(); app != nil {
     		if w := app.Window.Current(); w != nil {
     			w.Minimise()
     		}
     	}
     }

     func (a *App) OnMaximize() {
     	if app := application.Get(); app != nil {
     		if w := app.Window.Current(); w != nil {
     			w.ToggleMaximise()
     		}
     	}
     }

     func (a *App) OnClose() {
     	if app := application.Get(); app != nil {
     		app.Quit()
     	}
     }
     ```

2. **Make the window frameless in `app/main.go`.**
   - In `app.Window.NewWithOptions(application.WebviewWindowOptions{ ... })`:
     - Add `Frameless: true,`
     - Change `BackgroundColour` from `application.NewRGB(255, 253, 245)` to `application.NewRGB(255, 255, 255)` (white `#FFFFFF`).
   - Leave all other lines in `app/main.go` unchanged.

3. **Create `app/frontend/src/components/WindowBar.vue`.**
   - Import `OnClose`, `OnMaximize`, and `OnMinimize` directly from `@bindings/palabrario/app/app.js` (generated in Step 5).
   - Purge all reference jargon from `WindowBar.md` (no `useRoute`, no `Illima Local` breadcrumb, no `DEV` status pill, no dark mode classes, no shadows, no transitions/hover/active effects).
   - Render a sticky/fixed top `<header>` bar:
     - Classes: `sticky top-0 z-50 h-10 w-full bg-white border-b border-solid border-dull px-4 flex items-center select-none`
     - Drag region attribute/style: `style="--wails-draggable: drag; -webkit-app-region: drag;"` and `@dblclick="OnMaximize"`.
     - Inside `<header>`, render the traffic-light button group on the left (`class="flex items-center gap-2" style="--wails-draggable: no-drag; -webkit-app-region: no-drag;"`):
       1. Close button: `<button type="button" aria-label="Close window" title="Close" class="w-3 h-3 rounded-full bg-[#FF5F56] border border-solid border-[#E0443E]" @click="OnClose" />`
       2. Minimize button: `<button type="button" aria-label="Minimize window" title="Minimize" class="w-3 h-3 rounded-full bg-[#FFBD2E] border border-solid border-[#DEA123]" @click="OnMinimize" />`
       3. Maximize button: `<button type="button" aria-label="Maximize window" title="Maximize" class="w-3 h-3 rounded-full bg-[#27C93F] border border-solid border-[#1AAB29]" @click="OnMaximize" />`

4. **Delete `app/frontend/src/components/WindowBar.md`.**
   - Remove the temporary reference markdown file `app/frontend/src/components/WindowBar.md`.

5. **Wire `WindowBar.vue` into `app/frontend/src/App.vue`.**
   - Edit `app/frontend/src/App.vue` to render `<WindowBar />` above `<RouterView />`:
     ```vue
     <script setup lang="ts">
     import { RouterView } from "vue-router"
     import WindowBar from "@/components/WindowBar.vue"
     </script>

     <template>
       <div class="min-h-screen w-full flex flex-col bg-white text-main">
         <WindowBar />
         <div class="flex-1 flex flex-col">
           <RouterView />
         </div>
       </div>
     </template>
     ```

6. **Validate and regenerate bindings.**
   - Run `cd app && wails3 task check`.
   - This step runs `wails3 generate bindings -clean=true -ts`, which adds `OnMinimize`, `OnMaximize`, and `OnClose` to `app/frontend/bindings/palabrario/app/app.ts`, then runs `vue-tsc --noEmit`, `vite build`, and `go build .`.
   - Keep the updated `app/frontend/bindings/palabrario/app/app.ts` (since it exposes the 3 new methods); if `bindings/palabrario/pkg/models/` has unrelated formatting diffs, you may restore `app/frontend/bindings/palabrario/pkg/models/` with `git checkout`.

7. **Write the change log** in `.tickets/006-build-window-bar/changelog.md`.

## Rules

- Do not include any `Illima`, `Inbox`, `DEV`, or Spanish strings from `WindowBar.md`.
- Do not use `shadow-*`, `hover:`, `focus:`, `active:`, `transition-*`, or ` dark:` classes.
- Only the 3 traffic-light dots (`w-3 h-3 rounded-full`) may use `rounded-full`.
- Do not modify `app/frontend/uno.config.ts`. Propose any color tokens (such as macOS traffic-light aliases if desired) in `changelog.md`.

## Verification

All checks are in [spec.md](spec.md). The router runs them before accepting the ticket, not the worker's self-report.

## Change log

Write `.tickets/006-build-window-bar/changelog.md` with:
- Each file created, edited, or deleted.
- Confirmation that Wails v3 bindings were generated for `OnMinimize`, `OnMaximize`, and `OnClose`.
- Output status of `wails3 task check`.
- Any shortcut or token proposals for `uno.config.ts`.
