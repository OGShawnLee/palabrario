---
name: ticket
description: Create a Palabrario ticket (task.md and spec.md) in .tickets/. Use when a request must be handed to a worker agent. The ticket must be so clear and specific that a less capable agent can follow it without guessing.
---

# Ticket

A ticket is a folder in `.tickets/` with two files: `task.md` (what to do) and `spec.md` (how to check it is done). The worker is often a less capable model that hallucinates or diverges when anything is unclear. Your job is to remove every ambiguity before the ticket is handed out.

## Before writing: ask the user

If anything is unclear, **ask the user before writing the ticket**. Do not guess, and do not fill gaps with assumptions. Ask when:
- The goal, the expected output, or the success criteria are not exact.
- You do not know which files may be touched, or which must not be touched.
- A file path, name, section, or term is uncertain.
- The request could be read in two ways.
- The scope looks too large for one ticket.
- The request conflicts with the ground rules in [CONTEXT.md](../../../.routing/CONTEXT.md).

Ask short, specific questions, and wait for the answers. Only write the ticket when you can state the task without any "maybe" or "probably".

## Folder and naming

- Location: `.tickets/<id>-<slug>/`
- `<id>`: the next free three-digit number (`001`, `002`, ...). List `.tickets/` to find it.
- `<slug>`: lowercase words joined by hyphens that name the outcome.
- New tickets start with status `draft`.
- Examples: `.tickets/001-translate-domain-infrastructure/` and `.tickets/002-translate-interfaces/`.

## `task.md`

Use exactly this structure:

1. **Title:** `# <id> — <short outcome>`
2. **Metadata table** with these rows: ID, Status, Model, Role, Spec (a link to `spec.md`).
3. **Scope**
   - "May create or edit": exact file paths.
   - "Must not touch": exact paths or folders. Always include the Go backend and any file the worker should not modify.
4. **Inputs:** a link to each file the worker must read, using relative paths. List only what is needed. Name the exact section if only part of a file applies.
5. **Task:** numbered steps in the order to do them.
6. **Verification:** one line pointing to `spec.md`. The checks live there, not in `task.md`.
7. **Change log:** where the worker writes `changelog.md` and exactly what it must record.

## `spec.md`

- A checklist of checks, each of which can be answered yes or no.
- Each check names a file, a number, or a command and the expected result.
- Prefer checks that can be run: a file exists, a count matches, a command exits with 0, a search returns nothing.
- Include a check that changes stay inside the ticket's scope.
- Include a check for the commit when one is made: docs-only or code-only as required, English, Conventional Commits, with a detailed body.

## Writing rules for the worker's sake

- One ticket, one outcome. Split big tasks.
- Use exact paths, exact names, exact numbers. Never "the relevant files" or "a few".
- Write steps as commands: "Create", "Translate", "Run". One action per step.
- Say what **not** to do. Weaker agents tend to add or summarize on their own.
- Do not use vague words: "appropriate", "clean", "nice", "as needed", "etc.", "if possible".
- Give a concrete example when a rule is not obvious.
- Give the expected output for any command.
- State the assumptions the worker must keep: language, terms to leave unchanged, naming rules.
- Tell the worker to stop and record the problem in the change log instead of guessing.
- Keep the whole ticket short enough to be read in one go.

## Model for the ticket

Set the Model row from the roles in CONTEXT.md:
- Implementation: Claude Sonnet 5.5.
- Documentation translation: Gemini 3.8 Flash, high thinking.

## After writing

1. Re-read the ticket once as if you were a weak agent. If you would have to guess anything, fix it or ask the user.
2. Tell the user the ticket path and any open point.
3. Do not hand out the ticket. The router does that after the user confirms.
