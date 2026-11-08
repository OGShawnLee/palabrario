# 006 — Change Log: Build Frameless WindowBar & Wails v3 Window Service

## Overview

Implemented the frameless macOS-style `WindowBar.vue` component and exposed Wails v3 window management methods (`OnMinimize`, `OnMaximize`, `OnClose`) from Go through Wails TypeScript bindings.

## Files Created, Edited, or Deleted

- **Created**:
  - `app/app_window.go`: Implements `OnMinimize()`, `OnMaximize()`, and `OnClose()` on `*App` using `github.com/wailsapp/wails/v3/pkg/application` (`application.Get()`, `app.Window.Current().Minimise()`, `app.Window.Current().ToggleMaximise()`, and `app.Quit()`).
  - `app/frontend/src/components/WindowBar.vue`: Minimal macOS-style window header bar (`h-10 bg-white border-b border-solid border-dull px-4 flex items-center select-none`) with drag regions (`--wails-draggable: drag; -webkit-app-region: drag;`), double-click to maximize, and 3 traffic-light buttons on the left (`#FF5F56`, `#FFBD2E`, `#27C93F` with `rounded-full`) marked `no-drag`.
  - `.tickets/006-build-window-bar/changelog.md`: This change log.

- **Edited**:
  - `app/main.go`: Configured `Frameless: true` and set `BackgroundColour: application.NewRGB(255, 255, 255)` in `WebviewWindowOptions`.
  - `app/frontend/src/App.vue`: Wired `<WindowBar />` above `<RouterView />` in a flex column layout.
  - `app/frontend/bindings/palabrario/app/app.ts`: Regenerated via `wails3 generate bindings -clean=true -ts`, exporting `OnMinimize`, `OnMaximize`, and `OnClose`.

- **Deleted**:
  - `app/frontend/src/components/WindowBar.md`: Removed temporary reference file after implementation.

- **Unchanged**:
  - `app/frontend/uno.config.ts`: Untouched as mandated by scope.

## Wails v3 Bindings Confirmation

- `app/frontend/bindings/palabrario/app/app.ts` now exports:
  - `export function OnClose(): $CancellablePromise<void>`
  - `export function OnMaximize(): $CancellablePromise<void>`
  - `export function OnMinimize(): $CancellablePromise<void>`
- Unrelated formatting diffs in `app/frontend/bindings/palabrario/pkg/models/` were restored via `git checkout`.

## Validation & Verification

- `cd app && wails3 task check` passed with exit code 0:
  - `wails3 generate bindings -clean=true -ts`: 0 errors.
  - `bun run --cwd frontend type-check` (`vue-tsc --noEmit`): 0 errors.
  - `bun run --cwd frontend build` (`vite build --mode production`): 0 errors, generated `dist/` bundle.
  - `go build .`: compiled cleanly.
- `WindowBar.vue` audit:
  - Header classes: `sticky top-0 z-50 h-10 w-full bg-white border-b border-solid border-neutral-100 px-4 grid grid-cols-12 items-center select-none`
  - Drag region: `--wails-draggable: drag; -webkit-app-region: drag;` on `<header>` and `--wails-draggable: no-drag; -webkit-app-region: no-drag;` on controls container.
  - Double click calls `OnMaximize`.
  - Verified no occurrences of `Illima`, `DEV`, `Inbox`, `dark:`, `shadow`, `transition`, `hover:`, `focus:`, or `active:`.

## User Manual Refinements

- Updated header bottom border from `border-dull` to `border-neutral-100` (`border-b border-solid border-neutral-100`).
- Added centered `Palabrario` branding text across columns 5-8 inside a 12-column grid layout:
  - Col 1-4: Traffic-light button controls (`no-drag`).
  - Col 5-8: `<span class="col-span-4 text-center font-black tracking-tight text-main leading-tight">Palabrario</span>`.
  - Col 9-12: Empty column spacer `<span class="col-span-4"></span>` to preserve horizontal balance.

## Token & Shortcut Proposals for `uno.config.ts`

- Currently, macOS traffic-light colors are inlined as hex arbitrary values:
  - Close: `bg-[#FF5F56] border-[#E0443E]`
  - Minimize: `bg-[#FFBD2E] border-[#DEA123]`
  - Maximize: `bg-[#27C93F] border-[#1AAB29]`
- Proposal if tokenized in future design system updates:
  ```ts
  theme: {
    colors: {
      "traffic-close": "#FF5F56",
      "traffic-close-border": "#E0443E",
      "traffic-minimize": "#FFBD2E",
      "traffic-minimize-border": "#DEA123",
      "traffic-maximize": "#27C93F",
      "traffic-maximize-border": "#1AAB29",
    }
  }
  ```
