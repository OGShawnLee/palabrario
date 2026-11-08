# 008 — Translate README.md to English and Simplify Architecture

| Field | Value |
| :--- | :--- |
| **ID** | 008 |
| **Status** | completed |
| **Model** | Gemini 3.8 Flash, `thinking_level: high` (`gemini-3.8-flash`, Antigravity provider) |
| **Role** | Translator |
| **Spec** | [spec.md](spec.md) |

## Scope

May edit:
- `README.md`
- `.tickets/008-translate-readme/changelog.md` (create)

Must not touch:
- `app/` or any application code
- `core-service/`, `realtime-service/`, `image-service/`, `pkg/`, `sql/`
- `.context/`
- `.routing/` (read-only reference)
- Any file outside `README.md` and this ticket folder

## Inputs

- [README.md](../../README.md): Spanish source.
- [.routing/CONTEXT.md](../../.routing/CONTEXT.md): architecture context.
- [.routing/PROGRESS.md](../../.routing/PROGRESS.md): active progress.
- [.routing/ROUTES.md](../../.routing/ROUTES.md): active routes.

## Task

1. Translate `README.md` in full into clear, professional technical English.
2. **Update the Architectural Documentation Map table**:
   - Link only verified files: `.routing/CONTEXT.md`, `.routing/PROGRESS.md`, `.routing/ROUTES.md`, `.routing/DESIGN.md`, `.routing/DOMAIN.md`, `.routing/INTERFACES.md`, `SCHEMA_FINAL.sql`, `.routing/INFRASTRUCTURE.md`, and `.context/academic/`.
   - Remove references to legacy files that have no `.routing/` counterpart (`PALETTE.md`, `UI_FLOWS.md`, `UI_SPECIFICATION.md`, and `.context/PLAN.md`).
3. **Rework Section 9**:
   - Retitle to `## 🗺️ 9. Project Status & Roadmap`.
   - State that the backend architecture is 100% complete and verified (Go 1.24, PostgreSQL 17, `core-service`, `realtime-service`, `image-service`, 16 Use Cases, 13 Business Rules, RAM test suite).
   - State that the frontend is currently undergoing a complete from-scratch UI rewrite under `app/` (Wails v3, Vue 3, UnoCSS, monochrome design system) tracked in `.routing/PROGRESS.md` and `.routing/ROUTES.md`.
4. **Remove Section 10**:
   - Completely delete `## 📖 10. Bitácora de la Imprenta: Los Tejedores`.
5. Keep sections 1 through 8, the badges, Docker instructions, test harness, seed data tables, and License intact in English.
6. Write `.tickets/008-translate-readme/changelog.md`.

## Verification

All checks are listed in [spec.md](spec.md).

## Change log

Write `.tickets/008-translate-readme/changelog.md` detailing lines translated, removed sections, and link updates.

