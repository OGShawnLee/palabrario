# 021 — Spec checklist

Run by the router before accepting the ticket. All items must pass.

- [ ] Route `/collections/:id` is registered under `AppLayout` in `app/frontend/src/router/index.ts`.
- [ ] Universal breadcrumb displays the path hierarchy (`Home / Collections / ...`).
- [ ] `CollectionDetailView.vue` consumes collection details and constituent decks from `api.studio.getCollection(id)` via `import { api } from "@/api"`. Zero hardcoded mock collection data in views.
- [ ] Role-based owner/admin checks dynamically read `api.auth.getMe()` rather than hardcoded boolean flags.
- [ ] `IStudioAPI` contract and adapters (Wails and Mock) support `getCollection(id: number)`, `saveCollection(id: number)`, and `unsaveCollection(id: number)`.
- [ ] `CollectionDetailView.vue` implements the approved single-column layout:
  - Header with eyebrow, large title, blurb, metadata row, and primary actions.
  - "Save to Studio" (`button-main`) toggles saved state with Lucide bookmark icon.
  - "Study First Deck →" (`button-side`) is prominently positioned.
  - Sequential list of decks with term count, title, blurb, and "Explore Deck →" action.
  - Conditional bottom management section ("Alternative view • Owner / Administrator") with "Edit Collection", "Add Deck", and "Delete".
- [ ] Mobile-first responsive reflow:
  - Header actions stack vertically on mobile (< sm).
  - Deck cards stack term info and full-width "Explore Deck →" button.
  - Management actions stack cleanly matching mobile wireframe.
- [ ] English UI copy used consistently across all labels.
- [ ] Section labels use natural sentence case (`section-label` / `text-sm font-bold text-dull`).
- [ ] Strictly zero `flag_emoji` and zero dialect codes (`CU`, `CL`, `MX`). Full country names only.
- [ ] Strict monochrome styling (sharp corners `rounded-none`, no drop shadows, 1px solid borders).
- [ ] `cd app && wails3 task check` exits with code 0.
- [ ] Changes stay strictly within the declared scope.
- [ ] `.tickets/021-build-collection-detail-route/changelog.md` exists and is populated.
- [ ] Commit (if made): Conventional Commits, English, detailed body.
