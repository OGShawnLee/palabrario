# CONTEXT

Central routing directory for all agents. Read this first.

## Files
- [DESIGN.md](DESIGN.md): monochrome design decisions (shape, color, type, buttons, fields).
- [PROGRESS.md](PROGRESS.md): pending work, next steps, and open design gaps.
- [DOMAIN.md](DOMAIN.md): domain service dispatcher and index. Directs to the modular `*_service.md` files in [services/](services/).
- [INFRASTRUCTURE.md](INFRASTRUCTURE.md): deployment, ports, Docker/Vagrant setup, and testing architecture.
- [ROUTES.md](ROUTES.md): route table, layouts, and access rules. Read only when a ticket touches a route.
- [API.md](API.md): frontend API layer, adapter contracts, static mock strategy, and runtime switching.

## Project
- **Palabrario:** a Spanish vocabulary platform with a curated dictionary, SM-2 spaced repetition study, and real-time multiplayer arenas.
- **Current branch `minibrario`:** a reduced-scope rework of the frontend. The backend is not modified.
- Domain reference: [DOMAIN.md](DOMAIN.md) (read when a route needs domain context).
- **Later, off scope:** an Astro website that will use an updated [packages/ui](../packages/ui) and a Bun monorepo setup. Not part of the current work.

## Stack
- **Backend (stable, untouched):** Go 1.24, Chi v5, sqlc, pgxpool. Services: [core-service](../core-service), [realtime-service](../realtime-service), [image-service](../image-service). PostgreSQL 17.
- **Desktop client:** Wails 3 with Go. Frontend: Vue 3, Vite, Bun, UnoCSS, configured in [uno.config.ts](../app/frontend/uno.config.ts).
- **Mobile client:** Capacitor, inside the frontend directory of the desktop app. Wails knows nothing of Capacitor, and Capacitor knows nothing of Wails.
- **Font:** Satoshi, via Fontshare.
- **Icons:** `lucide-vue-next`.

## Ground rules
- Forbidden to render `flag_emoji` or dialect codes (e.g. `CU`, `CL`, `MX`) anywhere in the app. Always render the full textual name of the region (e.g. `Cuba`, `Chile`, `Mexico`).
- Routes in scope: landing page, `auth/sign-in`, `auth/sign-up`. Others are decided as we move.
- Mobile-first design: write base styles for mobile viewports first and scale up with responsive breakpoints (`sm:`, `md:`, `lg:`).
- No route aliases: every route in `app/frontend/src/router/index.ts` has a single canonical path (do not use `alias`).
- The UI is written in English from scratch. Existing Spanish UI strings are replaced by the rework, not translated.
- The rework is route by route, one commit per route, merged to `main` when complete.
- Commits: English, Conventional Commits, with a detailed body describing what changed.
- Documentation translation commits contain documentation changes only, with no code changes.
- The academic spec ([1-PROJECT-SPECIFICATION.md](../.context/academic/1-PROJECT-SPECIFICATION.md)) stays in Spanish and is not translated.
- The Go backend is not modified during standard frontend routes, except when the user approves a dedicated full-stack or backend ticket (e.g., auth hardening, schema/query extensions).
- **No role (Router, Investigator, Implementer, Translator, or Reviewer) may commit, push, or merge without explicit user authorization.**
- [packages/ui](../packages/ui) is ignored by the rebuild. Deleting it is deferred.
- No frontend tests exist, and none are planned yet.
- Toolchain: Go, Bun, and Wails 3 are installed and the app runs.

## Agent roles
- **Router:** Coordinates the project, tracks progress, maintains `.routing/`, creates tickets, dispatches subagents, and runs final `spec.md` acceptance checks. Does not write app code and never commits without explicit user authorization. Instructions live in the [router skill](../.claude/skills/router/SKILL.md); tickets use the [ticket skill](../.claude/skills/ticket/SKILL.md).
- **Investigator (`research`):** Audits cross-layer codebases, database schemas, contracts, and edge cases before complex or cross-layer tickets are written. Produces a structured findings report; does not edit source code and never commits.
- **Implementer (`implementation`):** Default: Gemini 4 Hydrogen, `thinking_level: medium` (`gemini-4-hydrogen`). Executes one ticket at a time strictly within its declared `Scope`. Logs proposals/changes in `changelog.md` and never commits without explicit user authorization.
- **Translator (`implementation`, docs-only):** Gemini 3.8 Flash (`gemini-3.8-flash`) through the Antigravity provider, for documentation translation only. Never commits without explicit user authorization.
- **Reviewer (`review`):** Inspects completed diffs against `spec.md`, security boundaries, error paths, race conditions, and ground rules before Router acceptance (with the user holding final sign-off). Does not modify source code and never commits.
- Workers must not modify files outside their ticket's scope.
- Workers do not edit [uno.config.ts](../app/frontend/uno.config.ts) directly. They log proposed shortcuts or tokens for the router to approve.
- The router accepts a ticket only after its spec checks (and any Reviewer audit) pass, not on the worker's self-report.
