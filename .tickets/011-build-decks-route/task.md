 011 — Build the Public Decks Route (/decks)

| Field | Value |
| :--- | :--- |
| **ID** | 011 |
| **Status** | draft |
| **Model** | Claude Sonnet 5.5, `effort: medium` (`claude-sonnet-5-5`) |
| **Role** | Implementer |
| **Spec** | [spec.md](spec.md) |

> **Draft Note:**  
> This ticket is a draft placeholder. Layout, deck cards, search bar, official/community filters, and card count indicators will be specified in detail during ticket discussion before work begins.

## Scope

May create or edit:
- `app/frontend/src/views/DecksView.vue`
- `app/frontend/src/components/decks/` (new components if required)
- `app/frontend/src/api/contract/studio.contract.ts` (support listDecks)
- `app/frontend/src/api/live/wails/studio.wails.ts` (implement listDecks)
- `app/frontend/src/api/mock/mock.adapter.ts` (implement mock listDecks)
- `app/frontend/src/router/index.ts`
- `.tickets/011-build-decks-route/changelog.md`

Must not touch:
- `app/frontend/uno.config.ts` (propose shortcuts or tokens in changelog)
- Backend code, Go services, or `.context/`

## Inputs

- [.routing/API.md](../../.routing/API.md): Architectural rules (Rule 1: Never import `@bindings/.../app` in UI; always import `api` from `@/api`).
- [app/frontend/src/api/index.ts](../../app/frontend/src/api/index.ts): Universal `api` singleton (`api.studio`).
- [.routing/CONTEXT.md](../../.routing/CONTEXT.md): ground rules and constraints.
- [.routing/DESIGN.md](../../.routing/DESIGN.md): design tokens and surface hierarchy.
- [.routing/services/studio_service.md](../../.routing/services/studio_service.md): Studio decks, cards, and SM-2 readiness.

## Task

1. **API Client Expansion (`IStudioAPI`):**
   - In `app/frontend/src/api/contract/studio.contract.ts`:
     - Add `listDecks(q?: string, officialOnly?: boolean): Promise<Deck[]>`.
     - Import `Deck` from `@bindings/palabrario/pkg/models/index.js`.
   - In `app/frontend/src/api/live/wails/studio.wails.ts`:
     - Delegate to `Wails.ListDecks(q ?? "", officialOnly ?? false)`.
   - In `app/frontend/src/api/mock/mock.adapter.ts` and `decks.fixture.ts`:
     - Implement mock `listDecks(...)` returning filtered decks.
2. **Router Setup:**
   - Register `/decks` under `AppLayout` in `app/frontend/src/router/index.ts`.
3. **Decks Catalog View (`DecksView.vue`):**
   - Import `api` from `@/api`:
     ```ts
     import { api } from "@/api"
     ```
   - Fetch decks via `api.studio.listDecks(...)` with reactive loading and filter state (`officialOnly`, keyword search).
   - **Strict Rule:** Never hardcode mock decks or static arrays inside `DecksView.vue` or subcomponents. All data must flow through `api.studio`.
4. Validate with `cd app && wails3 task check`.
5. Document changes in `changelog.md`.

## Verification

All checks are listed in [spec.md](spec.md).

## Change log

Write `.tickets/011-build-decks-route/changelog.md` with created/modified files, check results, and shortcut proposals.
