# 016 — Modularize DOMAIN.md and INTERFACES.md into Domain Services in .routing/services/

| Field | Value |
| :--- | :--- |
| **ID** | 016 |
| **Status** | completed |
| **Model** | Claude Sonnet 5.5, `effort: medium` (`claude-sonnet-5-5`, Antigravity provider) |
| **Role** | Documentation / Architect |
| **Spec** | [spec.md](spec.md) |

> **Architectural Note:**  
> This ticket modularizes the domain business rules and technical interface contracts into 8 self-contained `*_service.md` files located in a dedicated [`.routing/services/`](../../.routing/services/) directory. Root [`.routing/DOMAIN.md`](../../.routing/DOMAIN.md) serves as the lightweight master service dispatcher table (with `TBD` placeholders for routes still under construction), and all repository documentation/tickets are updated to point to the new structure.

## Scope

May create, edit, or delete:
- `.routing/services/` (create directory)
- `.routing/services/identity_service.md` (create / relocate)
- `.routing/services/governance_service.md` (create / relocate)
- `.routing/services/dictionary_service.md` (create / relocate)
- `.routing/services/curation_service.md` (create / relocate)
- `.routing/services/studio_service.md` (create / relocate)
- `.routing/services/mastery_service.md` (create / relocate)
- `.routing/services/arena_service.md` (create / relocate)
- `.routing/services/media_service.md` (create / relocate)
- Remove any root `.routing/*_service.md` files so `.routing/` is clean
- `.routing/DOMAIN.md` (master service dispatcher pointing to `services/*_service.md`)
- `.routing/INTERFACES.md` (delete after distributing into `services/`)
- `.routing/CONTEXT.md` (update `## Files` and reference `services/`)
- `.routing/PROGRESS.md` (update progress tracker)
- `README.md` (update Architectural Blueprint table)
- `.claude/skills/router/SKILL.md` (update domain files list)
- Active tickets in `.tickets/`:
  - `.tickets/009-build-dictionary-route/task.md` (point to `.routing/services/dictionary_service.md`)
  - `.tickets/010-build-collections-route/task.md` (point to `.routing/services/studio_service.md`)
  - `.tickets/011-build-decks-route/task.md` (point to `.routing/services/studio_service.md`)
  - `.tickets/014-backend-search-fts/task.md` (point to `.routing/services/dictionary_service.md` and `studio_service.md`)
- `.tickets/016-modularize-domain-services/changelog.md` (document changes and line counts)

Must not touch:
- `app/` or any application code
- `core-service/`, `realtime-service/`, `image-service/`, `pkg/`, `sql/`
- `.context/` (read-only reference archive)
- Any files outside the declared scope

## Inputs

- [.routing/DOMAIN.md](../../.routing/DOMAIN.md): Source domain entities, bounded contexts, and business rules (`RN01–RN13`).
- Source SOA contracts, DTOs, `sqlc` queries, and HTTP/WS endpoints (from former `INTERFACES.md`).
- [.routing/ROUTES.md](../../.routing/ROUTES.md): Current route statuses (built vs planned/TBD).
- [.routing/CONTEXT.md](../../.routing/CONTEXT.md): Project context.
- [README.md](../../README.md): Master architectural blueprint.

## Task

1. **Place all 8 Modular Service Files in `.routing/services/`:**
   Ensure each file lives under `.routing/services/` and contains both Part A (Domain Model & business rules) and Part B (SOA Technical Contracts, DTOs, `sqlc` queries, and HTTP/WS endpoints):
   - `.routing/services/identity_service.md`: Domain Context 1 (Accounts & Identity, `RN01` Anchor Dialect) + Interface Module 1 (`IdentityService`, `CU01`).
   - `.routing/services/governance_service.md`: Domain Context 1 (Governance, `RN02` Editor Promotion, VIP Tickets, Last Admin Protection) + Interface Module 2 (`GovernanceService`, `CU02–CU04`).
   - `.routing/services/dictionary_service.md`: Domain Context 2 (Lexicon, Senses, Heteronyms, FTS Search) + Interface Module 3 (`DictionaryService`, `CU05–CU06`).
   - `.routing/services/curation_service.md`: Domain Context 3 (Quorum of 3, `RN06–RN08`, Authorship, Community Reports) + Interface Module 4 (`CurationService`, `CU07–CU09`).
   - `.routing/services/studio_service.md`: Domain Context 4 (Study Decks & Collections, `RN03–RN05` Privacy & Cloning) + Interface Module 5 (`StudioService`, `CU10–CU11`).
   - `.routing/services/mastery_service.md`: Domain Context 4 (SM-2 Spaced Repetition Engine `RN09`, Immutable XP Ledger `RN11`, Streaks, Achievements) + Interface Module 6 (`MasteryService`, `CU12–CU13`).
   - `.routing/services/arena_service.md`: Domain Context 5 (Realtime Multiplayer Arena, 3NF Match Rounds & Replay, S2S Token, `RN10` Monotonic Timer, `RN13` Image Quorum) + Interface Module 7 (`ArenaCatalogService` + `realtime-service`, `CU14–CU16`).
   - `.routing/services/media_service.md`: Media Domain Rules (`RN12` Sanitization, NSFW Filter) + Interface Module 8 (`image-service` on `:8002`).

   *Ensure any relative links inside these files (e.g. back to `../DOMAIN.md` or `../../SCHEMA_FINAL.sql`) properly reflect the new subdirectory depth.*

2. **Remove any stray service files in `.routing/` root:**
   - Ensure `.routing/` root contains ONLY `DOMAIN.md`, `CONTEXT.md`, `DESIGN.md`, `INFRASTRUCTURE.md`, `PROGRESS.md`, `ROUTES.md`, and the `services/` directory.

3. **Configure `.routing/DOMAIN.md` as the Root Dispatcher:**
   - Retain the polycentric Spanish linguistic premise and shared cross-cutting architectural sections (anti-hallucination contract rules, SOA/Wails IPC diagram, directory structure, and domain error status table).
   - Summary table must link directly into `services/`:

| Route / Context | Service File | Domain & Service | Key Responsibilities & Use Cases |
| :--- | :--- | :--- | :--- |
| `/auth/sign-in`, `/auth/sign-up` | [`services/identity_service.md`](services/identity_service.md) | `IdentityService` | Accounts, profiles, passwords, `RN01` Anchor Dialect (`CU01`) |
| `TBD` | [`services/governance_service.md`](services/governance_service.md) | `GovernanceService` | Roles, promotion, invitation tickets, admin directory (`CU02–CU04`) |
| `/dictionary` | [`services/dictionary_service.md`](services/dictionary_service.md) | `DictionaryService` | Lemmas, senses, FTS search, Word of the Day, heteronyms (`CU05, CU06`) |
| `TBD` | [`services/curation_service.md`](services/curation_service.md) | `CurationService` | Quorum review (`RN06–RN08`), proposals, community reports (`CU07–CU09`) |
| `TBD` | [`services/studio_service.md`](services/studio_service.md) | `StudioService` | Decks & collections catalog, studio management, privacy (`CU10, CU11`) |
| `TBD` | [`services/mastery_service.md`](services/mastery_service.md) | `MasteryService` | Spaced repetition SM-2 engine (`RN09`), immutable XP ledger (`RN11`) (`CU12, CU13`) |
| `TBD` | [`services/arena_service.md`](services/arena_service.md) | `ArenaCatalogService` + `realtime-service` | Real-time WebSocket hub, rounds, speed scoring, 3NF replay (`CU14–CU16`) |
| `TBD` *(Media uploads)* | [`services/media_service.md`](services/media_service.md) | `image-service` | EXIF sanitization, magic bytes, WebP, NSFW filter (`RN12`) |

4. **Delete `.routing/INTERFACES.md`:**
   - Confirm it is deleted and all contents are preserved in `services/`.

5. **Update Repository Documentation & Ticket Inputs:**
   - In `.routing/CONTEXT.md`: Update `## Files` to describe `DOMAIN.md` as the dispatcher and reference `.routing/services/`.
   - In `README.md`: Update the blueprint table to reference `[`.routing/DOMAIN.md`](.routing/DOMAIN.md)` as the dispatcher and list `.routing/services/*_service.md`.
   - In `.claude/skills/router/SKILL.md`: Reference `DOMAIN.md` and `.routing/services/*_service.md`.
   - In active tickets:
     - `.tickets/009-build-dictionary-route/task.md` $\to$ `[.routing/services/dictionary_service.md](../../.routing/services/dictionary_service.md)`
     - `.tickets/010-build-collections-route/task.md` $\to$ `[.routing/services/studio_service.md](../../.routing/services/studio_service.md)`
     - `.tickets/011-build-decks-route/task.md` $\to$ `[.routing/services/studio_service.md](../../.routing/services/studio_service.md)`
     - `.tickets/014-backend-search-fts/task.md` $\to$ `[.routing/services/dictionary_service.md](../../.routing/services/dictionary_service.md)` and `[.routing/services/studio_service.md](../../.routing/services/studio_service.md)`
   - In `.routing/PROGRESS.md`: Track Ticket 016 as active/completed; ensure Ticket 017 is removed.

6. **Validation & Changelog:**
   - Verify zero content loss: all 16 Use Cases (`CU01–CU16`) and 13 Business Rules (`RN01–RN13`) accounted for.
   - Verify zero broken links across active docs.
   - Record line counts and file relocations in `.tickets/016-modularize-domain-services/changelog.md`.

## Verification

All checks are listed in [spec.md](spec.md).
