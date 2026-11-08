# 003 — Change log

Carried out directly by Claude Sonnet 5.5 (`effort: medium`) at the user's request, not by a delegate.

## Created
- `app/` from `wails3 init -n app -t vue -mod palabrario/app -skipgomodtidy -q` (module `palabrario/app`, Go 1.26.0).
- `app/frontend/uno.config.ts`: `presetUno`, `presetWebFonts` (Fontshare, `Supreme`), and the shortcuts `btn`, `btn-main`, `btn-side`, `field`.
- `app/frontend/src/main.ts`, `app/frontend/src/App.vue`, `app/frontend/src/views/HomeView.vue`, `app/frontend/src/vite-env.d.ts`.

## Copied unchanged from `local-app/`
- `main.go`, `app.go`, `app_arena.go`, `app_auth.go`, `app_dictionary.go`, `app_governance.go`, `app_media.go`, `app_studio.go` (verified identical with `cmp`).
- `frontend/vite.config.ts`, `frontend/tsconfig.json`.

## Copied with edits
- `app/frontend/package.json` from `local-app/frontend/package.json`: removed the `@palabrario/ui` line, and changed `name` to `palabrario-app-frontend`.
- `app/Taskfile.yml` from `local-app/Taskfile.yml`: removed the line `bun run --cwd ../packages/ui type-check` from the `check` task.
- Both edited files now have LF line endings (the sources have CRLF), because `sed -i` rewrote them. Compare with `diff --strip-trailing-cr`.

## Edited
- `go.work`: added `./app` (with `go work use ./app`).
- `package.json` (root): added `app/frontend` to `workspaces`.
- `app/build/config.yml`: `productName` set to `"Palabrario"`.
- `app/frontend/index.html`: title set to `Palabrario`.

## Deleted from the template
- `app/greetservice.go`.
- `app/build/ios`, `app/build/android`, `app/build/docker`.
- Everything inside `app/frontend/src/` (recreated above).

## Beyond the ticket (deviations)
- `app/frontend/index.html` was rewritten more than the ticket allowed. The template version loaded `/style.css`, a neon dark theme with rounded corners, and a `.bg` image div. Both are gone. The body now has `bg-white text-black font-sans antialiased min-h-screen`, and `lang="en"`.
- Removed unused template assets: `app/frontend/public/style.css`, `bg-desktop.jpg`, `bg-mobile.jpg`, `vue.svg`, `Inter-Medium.ttf`, `app/frontend/Inter Font License.txt`, and `app/frontend/README.md`. `wails.png` is kept as the icon.
- Side effects of running the tooling: `bun.lock` and `go.work.sum` were modified by `bun install` and the Go workspace. Both are tracked files outside the ticket's scope.

## Fontshare check
- `@unocss/preset-web-fonts` 66.12.0 supports the `fontshare` provider. No fallback was needed.
- The built CSS contains the `Supreme` `@font-face` with `cdn.fontshare.com` URLs. The font is therefore loaded from the CDN at runtime, so the app needs network access for it. Self-hosting the font is a possible later step.

## Template metadata for the user to decide
From `app/build/config.yml`, still at template defaults:
- `companyName: "My Company"`
- `productIdentifier: "com.mycompany.myproduct"`
- `description: "A program that does X"`
- `copyright: "(c) 2025, My Company"`
- `comments: "Some Product Comments"`

## Results
- `bun install` (run at the repo root): exit 0.
- `wails3 task check`: exit 0. Bindings: 1 service, 121 methods, 6 enums, 72 models. Type-check, Vite build, and `go build .` all passed.
- `wails3 task dev`: Vite ready on `127.0.0.1:9245`, HTTP 200, `bin/palabrario.exe` launched. The processes were stopped afterwards.

## Changes made by the user after the run
The user confirmed the app works and then changed these by hand. The ticket and spec were updated to match.
- `app/frontend/uno.config.ts`:
  - `presetUno` replaced by `presetWind3`.
  - Font changed from Supreme to Satoshi.
  - Theme colors added: `main` (`black`), `side-dark` (`#262626`), `side-light` (`#A3A3A3`).
  - Shortcuts renamed and changed: `button`, `button-main` (`bg-main text-white`), `button-side` (`bg-white text-black`, previously light gray), and `field` (`border-2 border`, `placeholder-side-light`).
- `app/frontend/src/main.ts`: imports `@unocss/reset/tailwind-compat.css` after `virtual:uno.css`. `@unocss/reset` (`^66.12.0`) was added to the root `package.json`.
- `app/frontend/src/views/HomeView.vue`: uses the new class names, and the button labels are now `Sign In` and `Sign Up`.
- Root `Taskfile.yml`: the `client:*` tasks now use `dir: app` and refer to `app/bin/palabrario.exe`.
- `.routing/CONTEXT.md`, `.routing/DESIGN.md`, `.routing/PROGRESS.md`: font references changed to Satoshi.

## Notes for the user
- `app/build/` is git-ignored by the root `.gitignore` (`build/`), like `local-app/build/`. It will not be committed. Add a negation rule if you want it tracked.
- The bindings produce one service (the `App` struct, spread over 7 files), not one module per file. The spec has been corrected.
- The `check` task still has Spanish `summary:` texts, copied from `local-app`.
- Not verified: the window's visual result. The Wails window launched, but only the user can confirm how the page looks.
