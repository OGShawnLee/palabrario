# CONTEXT

Central routing directory for all agents. Read this first.

## Files
- [DESIGN.md](DESIGN.md): monochrome design decisions (shape, color, type, buttons, fields).
- [PROGRESS.md](PROGRESS.md): pending work, next steps, and open design gaps.
- [DOMAIN.md](DOMAIN.md): domain entities, bounded contexts, and business rules.
- [INTERFACES.md](INTERFACES.md): SOA contracts, DTOs, `sqlc` queries, and HTTP/WS endpoints.
- [INFRASTRUCTURE.md](INFRASTRUCTURE.md): deployment, ports, Docker/Vagrant setup, and testing architecture.

## Project
- **Palabrario:** a Spanish vocabulary platform with a curated dictionary, SM-2 spaced repetition study, and real-time multiplayer arenas.
- **Current branch `minibrario`:** a reduced-scope rework of the frontend. The backend is not modified.
- Domain reference: [DOMAIN.md](DOMAIN.md) (read when a route needs domain context).
- **Later, off scope:** an Astro website that will use an updated [packages/ui](../packages/ui) and a Bun monorepo setup. Not part of the current work.

## Stack
- **Backend (stable, untouched):** Go 1.24, Chi v5, sqlc, pgxpool. Services: [core-service](../core-service), [realtime-service](../realtime-service), [image-service](../image-service). PostgreSQL 17.
- **Desktop client:** Wails 3 with Go. Frontend: Vue 3, Vite, Bun, UnoCSS, configured in [uno.config.ts](../uno.config.ts).
- **Mobile client:** Capacitor, inside the frontend directory of the desktop app. Wails knows nothing of Capacitor, and Capacitor knows nothing of Wails.
- **Font:** Supreme, via Fontshare.
- **Icons:** `lucide-vue-next`.

## Ground rules
- Routes in scope: landing page, `auth/sign-in`, `auth/sign-up`. Others are decided as we move.
- The UI is written in English from scratch. Existing Spanish UI strings are replaced by the rework, not translated.
- The rework is route by route, one commit per route, merged to `main` when complete.
- Commits: English, Conventional Commits, with a detailed body describing what changed.
- Documentation translation commits contain documentation changes only, with no code changes.
- The academic spec ([1-PROJECT-SPECIFICATION.md](../.context/academic/1-PROJECT-SPECIFICATION.md)) stays in Spanish and is not translated.
- The Go backend is not modified.
- [packages/ui](../packages/ui) is ignored by the rebuild. Deleting it is deferred.
- No frontend tests exist, and none are planned yet.
- Toolchain: Go, Bun, and Wails 3 are installed and the app runs.

## Agent roles
- **Router:** Claude Opus 5.5 (`claude-opus-5-5`) or Gemini 4 Argon (`gemini-4-argon`), to be chosen. Tracks progress, acts as the domain expert, creates tickets, and hands them out. Its instructions are in the [router skill](../.claude/skills/router/SKILL.md). Tickets are created with the [ticket skill](../.claude/skills/ticket/SKILL.md).
- **Workers:** Claude Sonnet 5.5 (`claude-sonnet-5-5`), one ticket at a time, no reviewer for now.
- **Translators:** Gemini 3.8 Flash (`gemini-3.8-flash`) through the Antigravity provider, for documentation translation.
- **Reviewer:** the user. The user reads the code and may spawn an agent.
- Workers must not modify files outside their ticket's scope.
- Workers do not edit [uno.config.ts](../uno.config.ts) directly. They log proposed shortcuts or tokens for the router to approve.
- The router accepts a ticket only after its spec checks pass, not on the worker's self-report.
