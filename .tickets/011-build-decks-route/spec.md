# 011 — Spec checklist

Run by the router before accepting the ticket. All items must pass.

- [x] Route `/decks` is registered under `AppLayout` in `app/frontend/src/router/index.ts` with `meta.breadcrumb`.
- [x] Universal breadcrumb displays the path hierarchy (`Home / Decks`).
- [x] `DecksView.vue` consumes decks from `api.studio.listDecks()` via `import { api } from "@/api"`. Zero hardcoded mock decks or static arrays exist in views or components.
- [x] `IStudioAPI` contract and adapters (Wails and Mock) properly support `listDecks(q?: string, officialOnly?: boolean)`.
- [x] `DecksView.vue` implements the approved single-column layout:
  - Layout, spacing, typography, and loading/empty state boxes match `CollectionsView.vue` (`max-w-6xl mx-auto px-4 pt-10 pb-20`).
  - Big headline `Study Decks` with subhead narrative.
  - Featured deck spotlight card (`FeaturedDeckCard.vue`) with 2-column desktop (7/5 grid split) / 1-column mobile layout matching `FeaturedCollectionCard.vue`.
  - Left column displays eyebrow tag, deck title, blurb, "Study Deck →" (`button-main`), and "Save to Studio" (`button-side`).
  - Right column displays "Terms in this deck", preview count (`Preview · 3 of N terms`), and up to 3 term preview cards with hairline dividers.
- [x] Secondary decks catalog grid (`DeckGridItem.vue`):
  - 2-column grid on desktop, 1-column on mobile.
  - Each item displays status/term count tag, "Explore →" action button, large title, blurb, and hairline bottom border.
- [x] Mobile-first responsive reflow:
  - Featured card stacks CTAs vertically on mobile (< sm).
  - Featured card stacks left hero column and right terms preview column on mobile (< md).
  - "More Decks" grid reflows to 1 column on mobile.
- [x] English UI copy used consistently across all labels.
- [x] Section labels use natural sentence case (`section-label` / `text-sm font-bold text-dull`).
- [x] Strictly zero `flag_emoji` and zero dialect codes (`CU`, `CL`, `MX`). Full country names only.
- [x] Strict monochrome styling (sharp corners `rounded-none`, no drop shadows, 1px solid borders).
- [x] `cd app && wails3 task check` exits with code 0.
- [x] Changes stay strictly within the declared scope.
- [x] `.tickets/011-build-decks-route/changelog.md` exists and is populated.
- [x] Commit (if made): Conventional Commits, English, detailed body.
