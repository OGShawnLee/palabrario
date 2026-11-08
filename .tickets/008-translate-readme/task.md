# 008 — Translate README.md to English and update architecture references

| Field | Value |
| :--- | :--- |
| **ID** | 008 |
| **Status** | draft |
| **Model** | Gemini 3.8 Flash, `thinking_level: high` (`gemini-3.8-flash`, Antigravity provider) |
| **Role** | Translator |
| **Spec** | [spec.md](spec.md) |

> **Draft Note / Blocker:**  
> This ticket is held in `draft` status until `.routing/` has stable and complete documentation. The existing `README.md` references the older `.context/` directory and earlier monorepo layout (`local-app/`, `packages/ui`). When activated, it will translate `README.md` to English and retarget architecture links to `.routing/`.

## Scope

May create or edit:
- `README.md` (edit in place)

Must not touch:
- `.routing/` (read-only reference)
- `.context/` (legacy reference, read-only)
- `app/` or any application code
- `core-service/`, `realtime-service/`, `image-service/`, `pkg/`, `sql/`
- Any file outside `README.md` and this ticket folder

## Inputs

- [README.md](../../README.md): source in Spanish.
- [.routing/CONTEXT.md](../../.routing/CONTEXT.md): active context and architecture overview.
- [.routing/PROGRESS.md](../../.routing/PROGRESS.md): current progress and roadmap.
- Terminology glossary: retain standard terms unchanged (FTS, SM-2, RBAC, JWT, WebSocket, Wails v3, Chi, etc.).

## Task

1. Translate `README.md` into clear, professional technical English.
2. Update the "Architectural Documentation Map" table:
   - Redirect links from `.context/*` to the corresponding translated `.routing/*` files where available ([`.routing/CONTEXT.md`](../../.routing/CONTEXT.md), [`.routing/DOMAIN.md`](../../.routing/DOMAIN.md), [`.routing/INTERFACES.md`](../../.routing/INTERFACES.md), [`.routing/INFRASTRUCTURE.md`](../../.routing/INFRASTRUCTURE.md)).
   - Note any pending documents clearly.
3. Update project structure sections to reflect the current `app/` desktop client on `minibrario` while preserving the stable Go backend services description.
4. Keep all shields/badges, formatting, and tables aligned.
5. Write the change log in `changelog.md`.

## Verification

All checks are listed in [spec.md](spec.md). The router runs them before accepting the ticket, not the worker's self-report.

## Change log

Write `.tickets/008-translate-readme/changelog.md` with:
- Line counts (source vs translated).
- Documentation links updated to `.routing/`.
- Unresolved terms or questions for the user.
