# 006 — Spec checklist

Run by the router before accepting the ticket. All items must pass.

## Go Service & Window Options
- [x] `app/app_window.go` exists and defines `OnMinimize()`, `OnMaximize()`, and `OnClose()` on `*App` using `github.com/wailsapp/wails/v3/pkg/application`.
- [x] `app/main.go` sets `Frameless: true` and `BackgroundColour: application.NewRGB(255, 255, 255)` in `WebviewWindowOptions`.
- [x] `app/frontend/bindings/palabrario/app/app.ts` exports `OnMinimize`, `OnMaximize`, and `OnClose`.

## WindowBar Component & App Wiring
- [x] `app/frontend/src/components/WindowBar.vue` exists and `app/frontend/src/components/WindowBar.md` has been deleted.
- [x] `WindowBar.vue` uses `h-10`, `bg-white`, `border-b border-solid border-neutral-100` (updated by user from `border-dull`), and contains centered `Palabrario` branding text.
- [x] `WindowBar.vue` sets drag attributes (`--wails-draggable: drag` / `-webkit-app-region: drag`) on the header and `no-drag` on the button group.
- [x] `WindowBar.vue` renders the 3 traffic-light buttons (`#FF5F56`, `#FFBD2E`, `#27C93F`) wired to `OnClose`, `OnMinimize`, and `OnMaximize`, and double-clicking the bar calls `OnMaximize`.
- [x] `WindowBar.vue` contains no `Illima`, `DEV`, `Inbox`, `dark:`, `shadow`, `transition`, `hover:`, `focus:`, or `active:` strings.
- [x] `app/frontend/src/App.vue` renders `<WindowBar />` above `<RouterView />`.

## Build and Scope
- [x] `cd app && wails3 task check` exits 0.
- [x] `app/frontend/uno.config.ts` is unchanged.
- [x] `.tickets/006-build-window-bar/changelog.md` exists and covers every required item.
