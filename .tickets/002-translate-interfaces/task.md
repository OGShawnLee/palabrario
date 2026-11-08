# 002 — Translate INTERFACES.md to English

| Field | Value |
| :--- | :--- |
| **ID** | 002 |
| **Status** | completed |
| **Model** | Gemini 3.8 Flash, `thinking_level: high` (`gemini-3.8-flash`, Antigravity provider) |
| **Role** | Translator |
| **Spec** | [spec.md](spec.md) |

## Scope

May create:
- `.routing/INTERFACES.md` (English translation of `.context/INTERFACES.md`)

Must not touch:
- `.context/` (source files, read-only)
- Any file outside `.routing/` and this ticket folder
- Any code, `uno.config.ts`, or `packages/`

## Inputs

- [.context/INTERFACES.md](../../.context/INTERFACES.md): source. Translate in full.
- [CONTEXT.md](../../.routing/CONTEXT.md): project context.
- Terminology glossary: not yet created. Until it exists, keep these terms unchanged: RBAC, FTS, SM-2, CEFR, XP, JWT, S2S, WebSocket, PostgreSQL, sqlc, pgxpool, Chi, Wails, Vue, UnoCSS, Docker, Vagrant, and the role names `LEARNER`, `EDITOR`, `ADMIN`. Record any other term you are unsure about in the change log.

## Task

1. Read `.context/INTERFACES.md` and translate it in full into English. Write it to `.routing/INTERFACES.md`.
2. Preserve the structure exactly:
   - Same headings, in the same order, with the same levels and numbering.
   - Same lists and tables, with the same number of rows and items.
   - Code blocks, Mermaid diagrams, and inline code are kept unchanged. Translate only the human-readable labels inside diagrams if they are in Spanish, and keep node IDs and syntax unchanged.
   - Links keep their targets. Translate only the link text.
   - Identifiers, file names, endpoints, commands, environment variables, field names, and ports are not translated.
3. Do not summarize, shorten, or add content. A translation that removes text is a failed translation.
4. Write a change log (see below).

## Verification

All checks are listed in [spec.md](spec.md). The router runs them before accepting the ticket, not the worker's self-report.

## Change log

The worker writes its log in `.tickets/002-translate-interfaces/changelog.md` with:
- Any term you were unsure about, and how you translated it.
- Any content you could not translate, with its location.
- The line counts used for the length check (translated and reference).

Proposals for new terms or glossary entries go in the change log. They are not added to the glossary by the worker.
