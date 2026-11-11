# 023 — Spec checklist

Run by the router before accepting the ticket. All items must pass.

- [x] `DictionaryCollectionsSection.vue` created and mounted inside `DictionaryView.vue` below the Word of the Day section in `max-w-5xl mx-auto px-4 pb-20`.
- [x] Collections data loaded dynamically via `api.studio.listCollections()` through `import { api } from "@/api"`. Zero hardcoded mock arrays exist in view or component.
- [x] Broadsheet cultural editorial layout implemented based on attached references:
  - Header: Eyebrow `01 / Thematic Paths`, headline `Menos alfabeto. Más mundo.`, and editorial statement.
  - Featured Spotlight 1: Localized black stat box (`bg-black text-white rounded-none`) with total terms and deck count, title, description, constituent deck list with hairlines (`↗`), and `"Explore collection →"`.
  - Featured Spotlight 2 (if present): Asymmetric split with deck list and right-aligned term counts.
  - Closing Manifesto Banner: Quote *"Una palabra abre una puerta. Una colección, un mundo."* and `"View all collections →"`.
- [x] Strictly ONLY Satoshi font (`font-sans` with `tabular-nums` where numbers are formatted). **Zero `font-mono`**.
- [x] Language Division strictly enforced: UI chrome, buttons, tags, and links in **English** (`Thematic Paths`, `Explore collection →`, `View all collections →`, `decks`, `terms`); editorial content and quotes in **Spanish**.
- [x] Navigation wired correctly:
  - `"Explore collection →"` navigates to `/collections/:id`.
  - Constituent deck items navigate to `/decks`.
  - `"View all collections →"` navigates to `/collections`.
- [x] Responsive layout: stacks cleanly on mobile (< 640px) and displays expansive multi-column on desktop.
- [x] Zero unviable client-side hydration for country tags (no N+1 fetching of senses from client).
- [x] Strictly zero `flag_emoji` and zero dialect codes (`CU`, `CL`, `MX`). Full country names only.
- [x] Strict monochrome styling (sharp geometry `rounded-none`, no drop shadows, 1px solid borders).
- [x] `cd app && wails3 task check` exits with code 0.
- [x] Changes stay strictly within the declared scope.
- [x] `.tickets/023-add-collections-section-to-dictionary/changelog.md` exists and is populated.
- [x] Commit (if made): Conventional Commits, English, detailed body.
