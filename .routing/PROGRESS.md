# PROGRESS

Pending work for the minibrario rework. Decisions already made are in `DESIGN.md` and `CONTEXT.md`. Current stance: absolute minimalism. No styling is applied yet; these items are deferred until the routes are rebuilt.

## Next steps
- [ ] Configure `uno.config.ts`: font (Supreme via Fontshare preset), colors, border and sharp radius, and the core component shortcuts (buttons, fields).
- [ ] Rebuild the root landing page.
- [ ] Rebuild `auth/sign-in`.
- [ ] Rebuild `auth/sign-up`.
- [ ] Decide what happens to the Capacitor adapter (`capacitor.adapter.ts`) during the frontend rebuild.

## Agent workflow
- [ ] Define the router-worker skill. Roles are in `CONTEXT.md`.
  - Tickets: one folder per feature or task under `.tickets/`, committed to Git. Each has `task.md` (what to do and the context needed) and `spec.md` (how to check it is done).
  - Ticket status: `draft`, `in-progress`, or `completed`.
  - Execution: sequential, route by route, while `uno.config.ts` is still being built. Parallel work is deferred until interactivity and backend integration begin.
  - Ticket template (draft; to be researched before adoption):
    - Header: id, title, status, and the model that runs it.
    - Scope: files it may touch, and files it must not touch.
    - Inputs: context files to read.
    - Task: ordered steps.
    - Spec: checks that decide whether it's done.
    - Change log: proposals made, such as new `uno.config.ts` shortcuts, for the user to approve.
  - Open: whether a Gemini 4 Argon reviewer is worth adding once the first tickets run.

## Localization and documentation
- [ ] Translate the codebase documentation to English (`.context/`, `README.md`, `AGENT_WORKFLOW.md`, `proposals/`). Translators: Gemini 3.8 Flash. Flash tends to delete large chunks, so each translation ticket must check that the translated length is within 100 lines of the Spanish original.
  - [x] `.tickets/001-translate-domain-infrastructure/` (`completed`) — translated `.context/DOMAIN.md` (excluding section 3.6) to `.routing/DOMAIN.md` and `.context/INFRASTRUCTURE.md` to `.routing/INFRASTRUCTURE.md`.
  - [x] `.tickets/002-translate-interfaces/` (`completed`) — translated `.context/INTERFACES.md` into `.routing/INTERFACES.md`.
  - Pending translation targets: `README.md`, `AGENT_WORKFLOW.md`, `proposals/`.
- [ ] Create a terminology glossary before translating, so technical terms (RBAC, FTS, SM-2, quorum) stay consistent.
  - Proposed entries from `001-translate-domain-infrastructure` (pending user review): `Canonical Lemma` (*Lema Canónico*), `Sense` (*Acepción / Sentido Dialectal*), `Anchor Dialect` (*Dialecto Ancla*), `Engaging Cultural Blurb` (*Reseña Cultural Amena*), `Transdialectal Heteronymy` (*Heteronimia Transdialectal*), `Prompt` (*Reactivo*), `Deck` (*Baraja*), `Ledger` (*Libro mayor*).
  - Proposed entries from `002-translate-interfaces` (pending user review): `Invitation Ticket` (*Boleto de invitación*), `Internal Prompt` (*Reactivo interno*), `Batch Flush` (*Descarga por lotes*), `Session Vault` (*Bóveda de sesión*).
- [ ] Confirm the Fontshare setup works with the UnoCSS web fonts preset, using Supreme.

## Design system gaps
- [ ] Define the remaining component classes (cards, lists, tables, dividers, feedback messages) when their routes are built.

## Hover and focus states
- [ ] Define one focus indicator for all interactive elements (monochrome, visible, keyboard-operable). Required by RNF-USA-01 and WCAG 2.4.7. Do not remove focus styling.
- [ ] Decide hover behavior (currently: none planned, optional).
- [ ] Decide active/pressed behavior (currently: none planned, optional).

## Loading states
- [ ] Every operation over 200 ms shows a loading indicator (RNF-USA-02).
- [ ] Define the loading indicator style (spinner or skeleton) in monochrome.
- [ ] Arena round timer needs a visible countdown state.

## Medal coloring
- [ ] Medal tiers (Bronze, Silver, Gold, Platinum, Diamond) currently depend on color. In monochrome, differentiate them by label, weight, or shape.
- [ ] Decide the medal representation before the achievements route is rebuilt.
