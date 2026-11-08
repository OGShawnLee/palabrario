# 003 — Scaffold the new barebones app under `app/`

| Field | Value |
| :--- | :--- |
| **ID** | 003 |
| **Status** | completed |
| **Model** | Claude Sonnet 5.5, `effort: medium` (`claude-sonnet-5-5`) |
| **Role** | Implementer (exception: the user asked the current session to carry it out directly, without a delegate) |
| **Spec** | [spec.md](spec.md) |

## Goal

Create a new desktop app under `app/` that compiles and runs. It has the same stack as `local-app/`, a new `uno.config.ts` that follows [DESIGN.md](../../.routing/DESIGN.md), Go bindings generated from the copied backend services, and one minimal home page.

## How the scaffold starts

Start from the Wails v3 template, then adapt it to mimic `local-app/`. The `build/` folder comes from the template, because it is git-ignored and missing from `local-app/` in this worktree.

The command below was tested. It creates `app/` at the repo root with `module palabrario/app`:

```
wails3 init -n app -t vue -mod palabrario/app -skipgomodtidy -q
```

Run it from the repo root. Do not add `-d`: with `-d app` the project is created one level deeper, at `app/app/`.

The template brings pieces we do not want: a greeting service, and `build/ios`, `build/android`, `build/docker`. Those folders are dead code and break `go build`. Delete them (step 2).

## Scope

May create:
- `app/` and everything inside it.
- `app/frontend/uno.config.ts`, which is the first version of the project's UnoCSS config.

May edit:
- `go.work`: add `./app` to the `use` block. Change nothing else.
- `package.json` (root): add `app/frontend` to `workspaces`, and the `@unocss/reset` dependency. Keep `packages/*` and `local-app/frontend` as they are.
- `Taskfile.yml` (root): point the `client:*` tasks at `app`.

Must not touch:
- `local-app/` (read-only; it is the reference copy).
- `packages/`.
- `.context/`, `.routing/`, and other tickets.
- `core-service/`, `realtime-service/`, `image-service/`, `pkg/`, and `sql/`.
- `go.mod` and `go.sum` at the root and in `local-app/`.
- The main checkout `Palabrario/` (a different folder, outside this worktree).

## Inputs

Read only these. Do not read other files in `local-app/` unless a check fails.
- [DESIGN.md](../../.routing/DESIGN.md): the visual rules for `uno.config.ts` and the home page.
- [CONTEXT.md](../../.routing/CONTEXT.md): stack and ground rules.
- [local-app/frontend/package.json](../../local-app/frontend/package.json): base dependency list.
- [local-app/Taskfile.yml](../../local-app/Taskfile.yml): base task list.
- [WAILS3-STACK.md](../../learning/WAILS3-STACK.md) (Spanish): lessons learned with this stack. The points that matter here:
  - The Go module must be registered in `go.work`, or bindings come out empty (`0 Services, 0 Methods`).
  - Use `go build .`, never `go build ./...`, inside a Wails project.
  - `typescript` must be 5.7 or higher.
  - Run tasks with `wails3 task <name>`.

Copy without reading: `local-app/main.go`, `local-app/app.go`, `local-app/app_arena.go`, `local-app/app_auth.go`, `local-app/app_dictionary.go`, `local-app/app_governance.go`, `local-app/app_media.go`, `local-app/app_studio.go`, `local-app/frontend/vite.config.ts`, `local-app/frontend/tsconfig.json`.

## Task

1. **Create the template.** Run the `wails3 init` command from "How the scaffold starts", from the repo root. Check that `app/go.mod` exists and its first line is `module palabrario/app`.
2. **Remove the template parts we do not want.**
   - Delete `app/greetservice.go`.
   - Delete `app/build/ios`, `app/build/android`, and `app/build/docker`.
   - Delete everything inside `app/frontend/src/`. Step 8 recreates what is needed.
3. **Register the module.** Add `./app` to `go.work` with `go work use ./app` run from the repo root. Check that `./local-app` is still in the `use` block.
4. **Copy the Go services.**
   - Copy `local-app/app.go` and `local-app/app_*.go` (7 files in total) into `app/`, unchanged.
   - Overwrite `app/main.go` with a copy of `local-app/main.go`, unchanged.
   - The copied files import `palabrario/pkg/models`. Keep those imports as they are; `go.work` resolves them.
5. **Adapt the frontend config.**
   - Overwrite `app/frontend/package.json` with a copy of `local-app/frontend/package.json`, with these changes only:
     - Remove the dependency `"@palabrario/ui": "workspace:*"`.
     - Change `"name"` to `"palabrario-app-frontend"`.
   - Keep every other dependency and script exactly as in the source.
   - Overwrite `app/frontend/vite.config.ts` and `app/frontend/tsconfig.json` with copies of the `local-app/frontend` ones. If `vite.config.ts` imports anything from `@palabrario/ui`, remove that import and record it in the change log.
   - In `app/frontend/index.html`, set the `<title>` to `Palabrario`. Change nothing else in that file.
   - Add `"app/frontend"` to the root `package.json` `workspaces` list.
6. **Create `app/frontend/uno.config.ts`.**
   - Use `defineConfig` from `unocss`, with `presetWind3` and `presetWebFonts`.
   - Set the sans font to `Satoshi` from Fontshare. Check that the installed `@unocss/preset-web-fonts` supports the `fontshare` provider. If not, load Satoshi with an `@font-face` rule in `app/frontend/index.html` instead, and record the choice in the change log.
   - Define these theme colors: `main` = `black`, `side-dark` = `#262626`, `side-light` = `#A3A3A3`.
   - Define these shortcuts, with values taken from [DESIGN.md](../../.routing/DESIGN.md):
     - `button`: `h-12 px-4 border-2 border-black border-solid rounded-none inline-flex items-center justify-center disabled:opacity-50`
     - `button-main`: `button bg-main text-white`
     - `button-side`: `button bg-white text-black`
     - `field`: `h-12 px-4 border-2 border border-solid rounded-none placeholder-side-light disabled:opacity-50`
   - Do not add any other shortcut, and do not add hover, focus, or active variants.
7. **Create the home page.**
   - Create `app/frontend/src/views/HomeView.vue`.
   - Content, in this order:
     1. A heading with the text `Palabrario`.
     2. A paragraph with the text `Vocabulary, dialects, and arenas.`
     3. A button with class `button-main` and the text `Sign In`.
     4. A button with class `button-side` and the text `Sign Up`.
   - Neither button calls the backend or any binding. The home page is a shell only.
   - All text is English. No icons.
8. **Wire the page.** Create `app/frontend/src/main.ts` and `app/frontend/src/App.vue` so `HomeView` is the only screen. In `main.ts`, import `virtual:uno.css` and then `@unocss/reset/tailwind-compat.css` (the `@unocss/reset` package is a dependency in the root `package.json`). Keep it as small as possible: no router, no store.
9. **Taskfile and build config.**
   - Overwrite `app/Taskfile.yml` with a copy of `local-app/Taskfile.yml`, then remove the line `bun run --cwd ../packages/ui type-check` from the `check` task. Leave every other line unchanged.
   - In the root `Taskfile.yml`, point the `client:*` tasks at `app` (`dir: app`, and the binary path `app/bin/palabrario.exe`).
   - In `app/build/config.yml`, set `productName` to `"Palabrario"`. Do not change other values, and list the template's company and copyright values in the change log for the user to decide.
10. **Validate.**
    - Run `cd app/frontend && bun install`.
    - Run `cd app && wails3 task check`. It regenerates bindings, type-checks the frontend, builds the frontend, and runs `go build .`. It must exit 0.
    - If it fails, record the exact error in the change log and stop. Do not guess a fix outside this ticket's scope.
11. **Write the change log** (see below).

## Rules

- Do not add dependencies beyond those in the source `package.json`.
- Never run `go build ./...`. Use `go build .` inside `app/`.
- Run tasks with `wails3 task <name>`, not with a separate `task` binary.
- Do not write comments in code except where a line is not obvious.
- Do not use `rounded-*`, `hover:`, `focus:`, or `active:` classes anywhere in `app/frontend/src`.
- Do not create extra pages, components, stores, or helpers.
- If a command fails, do not guess a fix. Record the exact error in the change log and stop.

## Verification

All checks are in [spec.md](spec.md). The router runs them. The worker does not accept its own work.

## Change log

Write `.tickets/003-scaffold-app/changelog.md` with:
- Each file deleted, each file copied, and each file created.
- Every edit made to a copied file, with the line changed.
- The result of the Fontshare provider check (supported, or fallback used).
- The template's company and copyright values from `app/build/config.yml`.
- Any command that failed, with its exact error text.
- Any question or proposal for the user.
