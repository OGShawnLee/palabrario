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
- `app/frontend/src/views/DecksView.vue` (or `app/frontend/src/views/community/DecksView.vue`)
- `app/frontend/src/components/decks/` (new components if required)
- `app/frontend/src/router/index.ts`
- `.tickets/011-build-decks-route/changelog.md`

Must not touch:
- `app/frontend/uno.config.ts` (propose shortcuts or tokens in changelog)
- Backend code, Go services, or `.context/`

## Inputs

- [.routing/CONTEXT.md](../../.routing/CONTEXT.md): ground rules and constraints.
- [.routing/DESIGN.md](../../.routing/DESIGN.md): design tokens and surface hierarchy.
- [.routing/services/studio_service.md](../../.routing/services/studio_service.md): Studio decks, cards, and SM-2 readiness.

## Task

1. Requirements and view layout to be defined during discussion.
2. Register `/decks` under `AppLayout` in `app/frontend/src/router/index.ts`.
3. Implement the public decks catalog following mobile-first design.
4. Validate with `cd app && wails3 task check`.
5. Document changes in `changelog.md`.

## Verification

All checks are listed in [spec.md](spec.md).

## Change log

Write `.tickets/011-build-decks-route/changelog.md` with created/modified files, check results, and shortcut proposals.
