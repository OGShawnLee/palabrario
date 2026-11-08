# 003 — Spec checklist

Run by the router before accepting the ticket. All items must pass.

## Structure
- [x] `app/go.mod` exists and its first line is `module palabrario/app`.
- [x] `app/` contains these 8 Go files: `main.go`, `app.go`, `app_arena.go`, `app_auth.go`, `app_dictionary.go`, `app_governance.go`, `app_media.go`, `app_studio.go`.
- [x] `app/greetservice.go` does not exist, and `grep -rni "greet" app --include=*.go --include=*.ts --include=*.vue` returns nothing (excluding `node_modules` and `bindings`).
- [x] Each copied Go file is identical to its source in `local-app/`. Check with `cmp`.
- [x] `go.work` contains `./app` in its `use` block, and `./local-app` is still there.
- [x] Root `package.json` `workspaces` contains `app/frontend`, and still contains `packages/*` and `local-app/frontend`.
- [x] `app/frontend/package.json` has no `@palabrario/ui` entry.
- [x] Every other dependency and script in `app/frontend/package.json` matches `local-app/frontend/package.json`. Check with `diff --strip-trailing-cr`; the only differences are `name` and the removed `@palabrario/ui` line.
- [x] `app/build/` exists and has no `ios`, `android`, or `docker` folder.
- [x] `app/build/config.yml` has `productName: "Palabrario"`.
- [x] `app/Taskfile.yml` matches `local-app/Taskfile.yml` except for the removed `packages/ui` line. Check with `diff --strip-trailing-cr`.
- [x] `app/frontend/uno.config.ts` exists and defines the theme colors `main`, `side-dark`, and `side-light`, and the shortcuts `button`, `button-main`, `button-side`, and `field`, and nothing else.
- [x] `app/frontend/src/main.ts` imports `virtual:uno.css` and `@unocss/reset/tailwind-compat.css`, and the root `package.json` lists `@unocss/reset`.
- [x] The root `Taskfile.yml` `client:*` tasks use `dir: app`.
- [x] `app/frontend/src/views/HomeView.vue` exists and contains the text `Palabrario`, `Vocabulary, dialects, and arenas.`, `Sign In`, and `Sign Up`.
- [x] `app/frontend/src/` contains only `main.ts`, `App.vue`, `vite-env.d.ts`, and `views/HomeView.vue`.

## Style rules
- [x] `grep -rn "rounded" app/frontend/src` returns nothing.
- [x] `grep -rnE "hover:|focus:|active:" app/frontend/src` returns nothing.
- [x] `grep -rn "@palabrario/ui" app/` returns nothing (excluding `node_modules`).
- [x] Visible text in `HomeView.vue` is English. Confirmed by the user.

## Compiles
- [x] `bun install` exits 0. Run it at the repo root, because `app/frontend` is a workspace.
- [x] `cd app && wails3 task check` exits 0. This single task regenerates bindings, type-checks, builds the frontend, and runs `go build .`.
- [x] `app/frontend/dist/index.html` exists after the check.
- [x] `app/frontend/bindings/palabrario/app/app.ts` exists. The bindings contain one service, `App`.

## Runs
- [x] `cd app && wails3 task dev` starts with no error in the output within 60 seconds.
- [x] `curl -s -o /dev/null -w "%{http_code}" http://localhost:9245` returns `200` while the dev task is running.
- [x] `curl -s http://localhost:9245` returns HTML that contains the text `Palabrario` (in the `<title>`).
- [x] Reviewer opens the Wails window and confirms the home page shows the heading, the tagline, and two buttons (one black, one white). Confirmed by the user ("it is working").

## Scope
- [x] `git status --short` shows changes only in: `app/`, `go.work`, `package.json`, `Taskfile.yml`, `.tickets/003-scaffold-app/`, plus the side-effect files `bun.lock` and `go.work.sum`.
- [x] `git status --short` shows no change in `local-app/`, `packages/`, `.context/`, `.routing/`, `core-service/`, `realtime-service/`, `image-service/`, `pkg/`, or `sql/`.
- [x] The main checkout at `Palabrario/` was not touched.

## Change log
- [x] `.tickets/003-scaffold-app/changelog.md` exists and covers every item listed in the ticket's Change log section.
