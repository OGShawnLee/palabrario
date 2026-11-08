# 001 — Translate DOMAIN.md and INFRASTRUCTURE.md to English

| Field | Value |
| :--- | :--- |
| **ID** | 001 |
| **Status** | completed |
| **Model** | Gemini 3.8 Flash, `thinking_level: high` (`gemini-3.8-flash`, Antigravity provider) |
| **Role** | Translator |
| **Spec** | [spec.md](spec.md) |

## Scope

May create:
- `.routing/DOMAIN.md` (English translation of `.context/DOMAIN.md`, without section 3.6)
- `.routing/INFRASTRUCTURE.md` (English translation of `.context/INFRASTRUCTURE.md`)

Must not touch:
- `.context/` (source files, read-only)
- Any file outside `.routing/` and this ticket folder
- Any code, `uno.config.ts`, or `packages/`

## Inputs

- [.context/DOMAIN.md](../../.context/DOMAIN.md): source. Translate everything except section 3.6.
- [.context/INFRASTRUCTURE.md](../../.context/INFRASTRUCTURE.md): source. Translate in full.
- [CONTEXT.md](../../.routing/CONTEXT.md): project context.
- Terminology glossary: not yet created. Until it exists, keep these terms unchanged: RBAC, FTS, SM-2, CEFR, XP, JWT, S2S, WebSocket, PostgreSQL, sqlc, pgxpool, Chi, Wails, Vue, UnoCSS, Docker, Vagrant, and the role names `LEARNER`, `EDITOR`, `ADMIN`. Record any other term you are unsure about in the change log.

## Task

1. Read `.context/DOMAIN.md`. Skip section 3.6 entirely: do not translate it, and do not leave a placeholder for it. Keep all other section numbers unchanged. Do not renumber the sections that follow 3.6.
2. Translate the remaining content of `.context/DOMAIN.md` into English and write it to `.routing/DOMAIN.md`.
3. Read `.context/INFRASTRUCTURE.md` and translate it in full into English. Write it to `.routing/INFRASTRUCTURE.md`.
4. Preserve the structure of both files exactly:
   - Same headings, in the same order, with the same levels.
   - Same lists and tables, with the same number of rows and items.
   - Code blocks, Mermaid diagrams, and inline code are kept unchanged. Translate only the human-readable labels inside diagrams if they are in Spanish, and keep node IDs and syntax unchanged.
   - Links keep their targets. Translate only the link text.
   - Identifiers, file names, commands, environment variables, and ports are not translated.
5. Do not summarize, shorten, or add content. A translation that removes text is a failed translation.
6. Write a change log (see below).

## Verification

All checks are listed in [spec.md](spec.md). The router runs them before accepting the ticket, not the worker's self-report.

## Change log

The worker writes its log in `.tickets/001-translate-domain-infrastructure/changelog.md` with:
- Confirmation that section 3.6 was omitted.
- Any term you were unsure about, and how you translated it.
- Any content you could not translate, with its location.
- The line counts used for the length check (translated and reference).

Proposals for new terms or glossary entries go in the change log. They are not added to the glossary by the worker.
