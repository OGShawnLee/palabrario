# PROGRESS

Work tracker for the minibrario rework. Decisions already made are in `DESIGN.md` and `CONTEXT.md`. Current stance: absolute minimalism. Items below are deferred until the routes are rebuilt.

## Done
- [x] `.tickets/001-translate-domain-infrastructure/` and `.tickets/002-translate-interfaces/` (both `completed`): see Localization below.
- [x] `.tickets/003-scaffold-app/` (`completed`): new app under `app/` (Wails 3, Vue 3, Vite, Bun, UnoCSS) with the Go services and bindings copied from `local-app/`, a Satoshi-based `uno.config.ts` with color aliases and the shortcuts `button`, `button-main`, `button-side`, `field`, and a minimal home page. `wails3 task check` and `wails3 task dev` pass.
- [x] Set up by hand after ticket 003: `@unocss/reset` imported in `main.ts`, `presetWind3`, Satoshi as the font, and the root `Taskfile.yml` `client:*` tasks point to `app`.
- [x] Fontshare setup confirmed working with the UnoCSS web fonts preset.
- [x] `.tickets/004-build-landing-page/` (`completed`): static mobile-first monochrome hero on `/` (`HomeView.vue`), Vue Router hash history setup (`router/index.ts`), and route placeholders for `/auth/sign-in` and `/auth/sign-up`.
- [x] `.tickets/005-build-auth-pages/` (`completed`): two-column mobile-first `AuthLayout.vue` (inverted left editorial panel on desktop), `AuthHeader.vue`, `FormField.vue`, `regions.ts`, reactive `SignInView.vue`, and 2-step `SignUpView.vue` with optional full-name toggle and region grid.
- [x] `.tickets/006-build-window-bar/` (`completed`): frameless window setup (`app/main.go`), Wails v3 window methods (`OnMinimize`, `OnMaximize`, `OnClose` in `app/app_window.go`), and macOS-style `WindowBar.vue` wired into `App.vue`.
- [x] `.tickets/007-build-topbar-and-app-layout/` (`completed`): responsive `TopBar.vue` (`max-w-6xl` + `grid-cols-12` center alignment), `TopBarMobileMenu.vue` (full-screen Reka UI `Dialog`), `TopBarDesktopStatus.vue` (Reka UI user status `Popover`), `TOP-BAR.md` state machine docs, universal `AppLayout.vue` wrapping non-auth routes, and `/dictionary`, `/arena`, `/studio`, `/settings` route placeholders.
- [x] `.tickets/008-translate-readme/` (`completed`): translated `README.md` to English, updated architectural map links to `.routing/`, reworked Section 9 (Project Status & Roadmap), and removed Section 10.
- [x] `.tickets/012-translate-routes-map/` (`completed`): created standalone `.routing/ROUTES.md` canonical route inventory table (9 active/planned routes) and linked from `.routing/CONTEXT.md`.
- [x] `.tickets/013-build-search-modal/` (`completed`): built responsive floating search combobox modal (`SearchModal.vue`) with Reka UI Dialog, `@vueuse/core` shortcuts, dynamic featured entry sync on keyboard navigation, and `--uno` scoped styles.
- [x] `.tickets/016-modularize-domain-services/` (`completed`): modularized monolithic `DOMAIN.md` and `INTERFACES.md` into 8 self-contained `*_service.md` files in `.routing/services/`, rewrote root `DOMAIN.md` as the route dispatcher with `TBD` placeholders, deleted `INTERFACES.md`, and updated all project documentation and active ticket references.
- [x] `.tickets/018-backend-arena-public-rooms/` (`completed`): backend schema, sqlc query, and endpoint `GET /api/v1/arena/matches/public` for listing waiting public matches.
- [x] `.tickets/014-backend-search-fts/` (`completed`): backend schema & sqlc FTS for decks & collections, plus expanded `WordSearchResult` (`part_of_speech`, `region_name`).
- [x] `.tickets/020-audit-service-specs-with-backend/` (`completed`): audited and synchronized all 8 `.routing/services/*_service.md` specifications against authoritative Go source code (`pkg/models/`, `core-service/internal/service/`, `core-service/internal/api/`, and `sql/queries/`).
- [x] `.tickets/009-build-dictionary-route/` (`completed`): built the `/dictionary` route (`DictionaryView.vue`), singleton `useSearchModal.ts` composable, typed `mockWordOfTheDay.ts`, and responsive `WordOfTheDayCard.vue` with dialectal senses, full region names, and cross-regional heteronym equivalents.
- [x] `.tickets/022-scaffold-api-adapter-architecture/` (`completed`): scaffolded the `src/api/` adapter architecture (pure TypeScript contracts, authentic Spanish mock fixtures, live `WailsAdapter`, and runtime mode switcher) for Auth, Dictionary, and Studio.
- [x] `.tickets/015-integrate-search-modal/` (`completed`): frontend integration wiring `SearchModal.vue` and `DictionaryView.vue` to the universal API client (`api.dictionary`, `api.studio`) with debouncing, eliminating all hardcoded mock arrays from views.
- [x] `.tickets/010-build-collections-route/` (`completed`): built the `/collections` public catalog backed by `api.studio.listCollections()`, universal `AppBreadcrumb.vue` mounted in `AppLayout.vue`, `FeaturedCollectionCard.vue`, `CollectionGridItem.vue`, and inspection dialog.
- [x] `.tickets/021-build-collection-detail-route/` (`completed`): built the `/collections/:id` individual learning path detail view backed by `api.studio.getCollection(id)` and `api.auth.getMe()`, dynamic leaf breadcrumb, and retired the temporary dialog.
- [x] `.tickets/011-build-decks-route/` (`completed`): built the `/decks` public catalog backed by `api.studio.listDecks()`, universal breadcrumb `Home / Decks`, `FeaturedDeckCard.vue` with term previews, and `DeckGridItem.vue`.
- [x] `.tickets/023-add-collections-section-to-dictionary/` (`completed`): built broadsheet cultural editorial discovery section (`DictionaryCollectionsSection.vue`) in `DictionaryView.vue`, wired `saveDeck` / `unsaveDeck` across API adapters, and updated `$figma-prompt-craft` skill.
- [x] `.tickets/017-build-arena-route/` (`completed`): built `/arena` lobby route, reusable Dialog components (`AppDialog.vue`, `AppDialogOverlay.vue`), `CreateMatchDialog.vue` (with subtitles), `JoinMatchDialog.vue` (dedicated private room PIN lookup & preview), `PublicRoomsTable.vue`, and end-to-end backend lookup query/endpoint (`GetWaitingMatchByRoomCode` / `GET /api/v1/arena/matches/code/{roomCode}`) with Wails bridge `FindMatchRoom`.
- [x] `.tickets/025-build-word-detail-route/` (`completed`): built `/dictionary/word/:id` (`WordDetailView.vue`) with invalid-id guard, dynamic leaf breadcrumb, parameterized `WordOfTheDayCard.vue` (`eyebrow`, `headingTag`, `footerLink`, softened `Report Word` button), and updated `SearchModal.vue` dictionary navigation.
- [x] `.tickets/026-build-profile-route/` (`completed`): built `/u/:username` (`ProfileView.vue`) in editorial Swiss-broadsheet style, wired `IAuthAPI.getPublicProfileByUsername` across Wails and Mock adapters, rendered unlocked-only achievements with temporary monochrome tier badges, added bottom owner-only strip, and added `"Profile"` links to TopBar menus.
- [x] `.tickets/027-wire-auth-session-and-topbar/` (`completed`): created `useAuth()` singleton composable (`hasActiveSession()` before `getMe()`, deduplicated `init()`), upgraded `MockAuthAPI` with `palabrario_mock_users` persistence, wired `SignInView.vue` and `SignUpView.vue` with `role="alert"` error banners and loading states, and replaced the hardcoded `eponiatowska` placeholder in `TopBar.vue`.

## In progress
- None.

## Draft tickets
- [ ] `.tickets/019-build-arena-history-route/` (`draft`): expand `IArenaAPI` and build dedicated `/arena/history` route and 3NF round-by-round replay modal (`MatchReplayDialog.vue`) backed by `api.arena`.
- [ ] `.tickets/024-backend-dynamic-region-detection/` (`ready`): implement dynamic top 1-3 dialect region detection in PostgreSQL / sqlc for Collections and Decks (`region_names text[]`).

## Next steps
- [ ] Implement live WebSocket match room (`ConnectArenaRoom`) and interactive in-game multiplayer match screen after `/arena` lobby route is built.
- [ ] Add router navigation guards (`guestOnly`, `requireAuth`) or build `/settings` / `/studio` next.
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
  - Pending translation targets: `AGENT_WORKFLOW.md`, `proposals/`.
- [ ] Create a terminology glossary before translating, so technical terms (RBAC, FTS, SM-2, quorum) stay consistent.
  - Proposed entries from `001-translate-domain-infrastructure` (pending user review): `Canonical Lemma` (*Lema Canónico*), `Sense` (*Acepción / Sentido Dialectal*), `Anchor Dialect` (*Dialecto Ancla*), `Engaging Cultural Blurb` (*Reseña Cultural Amena*), `Transdialectal Heteronymy` (*Heteronimia Transdialectal*), `Prompt` (*Reactivo*), `Deck` (*Baraja*), `Ledger` (*Libro mayor*).
  - Proposed entries from `002-translate-interfaces` (pending user review): `Invitation Ticket` (*Boleto de invitación*), `Internal Prompt` (*Reactivo interno*), `Batch Flush` (*Descarga por lotes*), `Session Vault` (*Bóveda de sesión*).

## Design system gaps
- [ ] Define the remaining component classes (cards, lists, tables, dividers, feedback messages) when their routes are built.
- [x] `field` border color defined (`border-dull`, `#737373`).
- [x] Contrast increased: theme colors are now `main` (`black`), `side` (`#404040`), and `dull` (`#737373`, ~4.6:1 contrast on white).
- [x] `header` shortcut added to `uno.config.ts` (`h-12 px-4 flex items-center justify-between border-b border-light`).
- [x] `section-label` shortcut added to `uno.config.ts` (`text-sm font-bold text-dull mb-3`) and all-caps styling discouraged in `DESIGN.md`.
- [x] `badge`, `badge-main`, and `badge-side` shortcuts added to `uno.config.ts` for monochrome metadata and CEFR tags.
- [x] `button-side--light` and `button--small` shortcuts added to `uno.config.ts`.
- [x] `dialog-overlay` shortcut added to `uno.config.ts` for modal and dialog backdrops.
- [ ] The Satoshi font loads from the Fontshare CDN at runtime, so the app needs network access for it. Decide whether to self-host it.

## Scaffold follow-ups
- [ ] `app/build/config.yml` still has template values: company, product identifier, description, and copyright.
- [x] Added negation rule `!app/build/` to root `.gitignore` so Wails build configuration is tracked.
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
- [ ] Medal tiers (Bronze, Silver, Gold, Platinum, Diamond) currently depend on color. In monochrome, differentiate them by label, weight, or shape (currently rendered as temporary monochrome `badge-side` labels on `/u/:username`).
- [ ] Decide the permanent monochrome medal tier visual treatment (weight, border, or iconography) before the achievements route is rebuilt.
