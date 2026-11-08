# 009 — Build the Dictionary Route (/dictionary)

| Field | Value |
| :--- | :--- |
| **ID** | 009 |
| **Status** | draft |
| **Model** | Claude Sonnet 5.5, `effort: medium` (`claude-sonnet-5-5`) |
| **Role** | Implementer |
| **Spec** | [spec.md](spec.md) |

> **Draft Note:**  
> This ticket is a draft placeholder. Specific UI layout, search interactions, filters (region/country, CEFR, grammatical category), and vocabulary term card components will be specified in detail during ticket discussion before work begins.

## Scope

May create or edit:
- `app/frontend/src/views/DictionaryView.vue`
- `app/frontend/src/components/dictionary/` (new components if required)
- `app/frontend/src/router/index.ts`
- `.tickets/009-build-dictionary-route/changelog.md`

Must not touch:
- `app/frontend/uno.config.ts` (propose shortcuts or tokens in changelog)
- `app/*.go`, backend services, or `.context/`
- Any file outside the above scope

## Inputs

- [.routing/CONTEXT.md](../../.routing/CONTEXT.md): architecture rules and mobile-first constraint.
- [.routing/DESIGN.md](../../.routing/DESIGN.md): monochrome theme tokens and component rules.
- [.routing/DOMAIN.md](../../.routing/DOMAIN.md): Section 3.1 & 3.2 (Lexicographical search and FTS contracts).
- [app/frontend/uno.config.ts](../../app/frontend/uno.config.ts): exact color tokens and shortcuts.

## Task

1. Requirements and wireframe to be defined during discussion.
2. Implement `app/frontend/src/views/DictionaryView.vue` following mobile-first monochrome rules.
3. Validate with `cd app && wails3 task check`.
4. Document changes in `changelog.md`.

## Verification

All checks are listed in [spec.md](spec.md).

## Change log

Write `.tickets/009-build-dictionary-route/changelog.md` with created/modified files, check results, and shortcut proposals.
