# 010 — Build the Public Collections Route (/collections)

| Field | Value |
| :--- | :--- |
| **ID** | 010 |
| **Status** | draft |
| **Model** | Claude Sonnet 5.5, `effort: medium` (`claude-sonnet-5-5`) |
| **Role** | Implementer |
| **Spec** | [spec.md](spec.md) |

> **Draft Note:**  
> This ticket is a draft placeholder. Layout, collection card design, search and featured curation filters, and deck listing interactions will be specified in detail during ticket discussion before work begins.

## Scope

May create or edit:
- `app/frontend/src/views/CollectionsView.vue` (or `app/frontend/src/views/community/CollectionsView.vue`)
- `app/frontend/src/components/collections/` (new components if required)
- `app/frontend/src/router/index.ts`
- `.tickets/010-build-collections-route/changelog.md`

Must not touch:
- `app/frontend/uno.config.ts` (propose shortcuts or tokens in changelog)
- Backend code, Go services, or `.context/`

## Inputs

- [.routing/CONTEXT.md](../../.routing/CONTEXT.md): ground rules and constraints.
- [.routing/DESIGN.md](../../.routing/DESIGN.md): design tokens and surface hierarchy.
- [.routing/DOMAIN.md](../../.routing/DOMAIN.md): Section 3.4 (Studio and collections domain rules).

## Task

1. Requirements and view layout to be defined during discussion.
2. Register `/collections` under `AppLayout` in `app/frontend/src/router/index.ts`.
3. Implement the collections showcase view following mobile-first design.
4. Validate with `cd app && wails3 task check`.
5. Document changes in `changelog.md`.

## Verification

All checks are listed in [spec.md](spec.md).

## Change log

Write `.tickets/010-build-collections-route/changelog.md` with created/modified files, check results, and shortcut proposals.
