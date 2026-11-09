# 016 — Spec checklist

Run by the router before accepting the ticket. All items must pass.

- [x] All 8 service files exist in `.routing/services/`:
  - `.routing/services/identity_service.md`
  - `.routing/services/governance_service.md`
  - `.routing/services/dictionary_service.md`
  - `.routing/services/curation_service.md`
  - `.routing/services/studio_service.md`
  - `.routing/services/mastery_service.md`
  - `.routing/services/arena_service.md`
  - `.routing/services/media_service.md`
- [x] No `*_service.md` files remain loose directly under `.routing/` root.
- [x] Each file in `.routing/services/` contains both its domain entities/rules (Part A) and technical contracts/DTOs/queries/endpoints (Part B).
- [x] Root `.routing/DOMAIN.md` exists and contains the dispatcher summary table linking to `services/<name>_service.md` with `TBD` for routes still under construction.
- [x] `.routing/INTERFACES.md` is removed.
- [x] `.routing/CONTEXT.md` `## Files` section is updated to reference the new domain dispatcher and `services/` directory.
- [x] `README.md` Architectural Blueprint table is updated to reference `DOMAIN.md` as the dispatcher and lists `.routing/services/*_service.md`.
- [x] `.claude/skills/router/SKILL.md` references `DOMAIN.md` and `.routing/services/*_service.md`.
- [x] Active ticket inputs (`009`, `010`, `011`, `014`) point to `.routing/services/<name>_service.md`.
- [x] `.routing/PROGRESS.md` tracks Ticket 016 and contains no reference to Ticket 017.
- [x] Zero content lost: all 16 Use Cases (`CU01–CU16`) and 13 Business Rules (`RN01–RN13`) are accounted for.
- [x] No broken links to `.routing/INTERFACES.md` in active documentation files (`README.md`, `.routing/`, `.tickets/`).
- [x] No application code, backend files, or `.context/` files modified.
- [x] `.tickets/016-modularize-domain-services/changelog.md` exists and is populated.
- [x] Commit (if made): docs-only, Conventional Commits, English, detailed body.
