# PROGRESS

Work tracker for the minibrario rework. Decisions already made are in `DESIGN.md` and `CONTEXT.md`. Current stance: absolute minimalism. Items below are deferred until the routes are rebuilt.

## Done
- [x] `.tickets/001-translate-domain-infrastructure/` and `.tickets/002-translate-interfaces/` (both `completed`): see Localization below.
- [x] `.tickets/003-scaffold-app/` (`completed`): new app under `app/` (Wails 3, Vue 3, Vite, Bun, UnoCSS) with the Go services and bindings copied from `local-app/`, a Satoshi-based `uno.config.ts` with color aliases and the shortcuts `button`, `button-main`, `button-side`, `field`, and a minimal home page. `wails3 task check` and `wails3 task dev` pass.
- [x] Set up by hand after ticket 003: `@unocss/reset` imported in `main.ts`, `presetWind3`, Satoshi as the font, and the root `Taskfile.yml` `client:*` tasks point to `app`.
- [x] Fontshare setup confirmed working with the UnoCSS web fonts preset.

## Next steps
- [ ] Rebuild the root landing page on top of the scaffold: `.tickets/004-build-landing-page/` (`draft`) — static mobile-first hero adapted from `HomeHero.vue` + Vue Router setup (`/`, `/auth/sign-in`, `/auth/sign-up`).
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

## Design system gaps
- [ ] Define the remaining component classes (cards, lists, tables, dividers, feedback messages) when their routes are built.
- [x] `field` border color defined (`border-dull`, `#737373`).
- [x] Contrast increased: theme colors are now `main` (`black`), `side` (`#404040`), and `dull` (`#737373`, ~4.6:1 contrast on white).
- [ ] The Satoshi font loads from the Fontshare CDN at runtime, so the app needs network access for it. Decide whether to self-host it.

## Scaffold follow-ups
- [ ] `app/build/config.yml` still has template values: company, product identifier, description, and copyright.
- [ ] `app/build/` is git-ignored by the root `.gitignore` (`build/`). Decide whether to add a negation rule so it is tracked.
- [ ] The `summary:` texts in `app/Taskfile.yml` are still in Spanish.
- [x] Moved `@unocss/reset` from the root `package.json` to `app/frontend/package.json`.

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
