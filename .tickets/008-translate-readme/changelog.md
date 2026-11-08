# Changelog — Ticket 008: Translate README.md to English and Simplify Architecture

## Summary
Translated the repository `README.md` from Spanish to clear, professional technical English. Updated the architectural documentation matrix to reference only valid `.routing/` files and canonical schemas, removing dangling references to legacy `.context/` documents. Reworked Section 9 to reflect the 100% verified backend and the active frontend rewrite in `app/`. Completely removed the legacy narrative Section 10 ("Bitácora de la Imprenta").

## Changes

### 1. Header & Badges
- Translated the mission statement and academic context subtitle into English.
- Preserved all status shields, versions, and centered HTML formatting.

### 2. Architectural Documentation Map
- Redirected all active documentation links to `.routing/`:
  - [`.routing/CONTEXT.md`](../../.routing/CONTEXT.md)
  - [`.routing/PROGRESS.md`](../../.routing/PROGRESS.md)
  - [`.routing/ROUTES.md`](../../.routing/ROUTES.md)
  - [`.routing/DESIGN.md`](../../.routing/DESIGN.md)
  - [`.routing/DOMAIN.md`](../../.routing/DOMAIN.md)
  - [`.routing/INTERFACES.md`](../../.routing/INTERFACES.md)
  - [`SCHEMA_FINAL.sql`](../../SCHEMA_FINAL.sql)
  - [`.routing/INFRASTRUCTURE.md`](../../.routing/INFRASTRUCTURE.md)
  - [`.context/academic/`](../../.context/academic) (preserved for formal academic course deliverables)
- Removed legacy rows that have no active counterpart in `.routing/`:
  - `PALETTE.md`
  - `ROUTES_MAP.md`
  - `UI_FLOWS.md`
  - `UI_SPECIFICATION.md`
  - `.context/PLAN.md`

### 3. Sections 1 through 8
- **Section 1 (Project Vision & Purpose):** Fully translated dialectal nuances, pan-hispanic examples (*guagua*, *popote*, *apapachar*), and core pillar summaries.
- **Section 2 (Network Architecture & Engineering Decisions):** Translated Mermaid diagram labels, notes, and the technology rationale table (trade-offs for Go 1.24, `sqlc`, Wails v3, PostgreSQL 17, and S2S sub-key security).
- **Section 3 (Use Case & Business Rules Matrix):** Translated all module names, use case descriptions (`CU01`–`CU16`), and business rules (`RN01`–`RN13`).
- **Section 4 (Service Topology & Ports):** Translated table headers and descriptions, updating client reference to `local-app` / `app`.
- **Section 5 (Taskfile Automation Catalog):** Translated all comments and category headers inside the bash automation snippet.
- **Section 6 (Quick Start Guide):** Translated prerequisites, fail-fast environment variable rules, deployment steps, and curl probes.
- **Section 7 (Testing Strategy: 80/20 Harness on RAM):** Translated service/domain tests and HTTP smoke test descriptions; updated reference link to `.routing/INFRASTRUCTURE.md`.
- **Section 8 (Deterministic Seed Data and Accounts):** Translated table columns, role titles, and lexical example notes.

### 4. Section 9 (Project Status & Roadmap)
- Retitled section to `## 🗺️ 9. Project Status & Roadmap`.
- Documented the verified 100% completion of the backend architecture across `core-service`, `realtime-service`, `image-service`, 16 Use Cases, 13 Business Rules, and RAM test harness.
- Documented the active scratch UI rewrite in `app/` (Wails v3, Vue 3, UnoCSS, monochrome design system), linking to `.routing/PROGRESS.md` and `.routing/ROUTES.md`.

### 5. Section 10 (Removal)
- Completely deleted `## 📖 10. Bitácora de la Imprenta: Los Tejedores`.

### 6. License
- Translated course attribution and academic purpose declaration.

## Metrics
- **Source line count:** 268 lines
- **Translated line count:** 259 lines
- **Net change:** -9 lines (removed Section 10 and outdated documentation rows; added detailed Section 9 status breakdown)