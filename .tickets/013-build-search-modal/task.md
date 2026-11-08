# 013 — Build the Global Floating Search Modal (SearchModal.vue)

| Field | Value |
| :--- | :--- |
| **ID** | 013 |
| **Status** | completed |
| **Model** | Sonnet 5.5 (`claude-sonnet-5-5`, Antigravity provider) |
| **Role** | Implementer |
| **Spec** | [spec.md](spec.md) |

## Scope

May create or edit:
- `app/frontend/src/components/SearchModal.vue` (create)
- `app/frontend/src/layouts/AppLayout.vue` (mount `SearchModal.vue` globally)
- `app/frontend/src/components/TopBar.vue` (wire search button click to open search modal)
- `.tickets/013-build-search-modal/changelog.md` (create)

Must not touch:
- `app/frontend/uno.config.ts` (the `binding` shortcut is already added by the Router)
- `app/*.go`, backend services, or `.context/`
- Any file outside the above scope

## Inputs

- [.routing/CONTEXT.md](../../.routing/CONTEXT.md): Ground rules, including the rule that `flag_emoji` and regional codes are strictly forbidden.
- [.routing/DESIGN.md](../../.routing/DESIGN.md): Monochrome tokens (`main`, `side`, `dull`, `light`, `ground`), `binding` key shortcut, and surface hierarchy.
- [app/frontend/uno.config.ts](../../app/frontend/uno.config.ts): Theme shortcuts and colors (`binding`, `button-icon`, `field`, etc.).
- [app/frontend/src/data/regions.ts](../../app/frontend/src/data/regions.ts): Canonical region definitions.
- Visual references:
  - Desktop search modal: [Attached image: desktop layout](../../../.t3/userdata/attachments/e7510c5f-4aee-4fe2-8b04-76e954f1ea52-74c96df1-dc92-41f8-9cec-2b00a196bee8.png)
  - Mobile search modal: [Attached image: mobile layout](../../../.t3/userdata/attachments/e7510c5f-4aee-4fe2-8b04-76e954f1ea52-dc16ac37-0bf4-46d8-99e8-1f223b703843.png)
  - Keyboard hint badge: [Attached image: key binding aesthetic](../../../.t3/userdata/attachments/e7510c5f-4aee-4fe2-8b04-76e954f1ea52-95bf9cd9-39fe-4150-a298-ede5249d1a76.png)

## Deviations from Visual References

1. **No random word button:** Omit the `Random Word` button shown on the right side of the filter tabs.
2. **Tab selection:** Render only three tabs: `[All]`, `[Dictionary]`, and `[Study Decks]`. Omit the `[By Region]` tab.
3. **No flag emojis and no regional codes:** Never render flag emojis or dialect codes (e.g. `CU`, `CL`, `MX`). Render only the full textual region name (e.g. `Cuba`, `Chile`, `Mexico`).

## Task

1. **Create `SearchModal.vue`** in `app/frontend/src/components/SearchModal.vue` using Reka UI (`DialogRoot`, `DialogPortal`, `DialogOverlay`, `DialogContent`, etc.):
   - **Trigger & Visibility:**
     - Listen for `Ctrl + K` and `Cmd + K` globally using `@vueuse/core` (`useEventListener` or `onKeyStroke`).
     - Expose open/close state so `AppLayout.vue` can govern it or receive `@open-search` from `TopBar.vue`.
   - **Responsive Shell & Backdrop:**
     - **Desktop (`md:` and up):** Floats slightly above the vertical center (`top-[15%]`), centered horizontally, max width (`max-w-3xl` or `max-w-4xl`), solid sharp borders (`border border-light`), white card surface (`bg-white`), over translucent backdrop (`bg-black/40` or `bg-neutral-900/40`).
     - **Mobile (`< md`):** Full screen modal (`fixed inset-0 bg-white z-50`).
   - **Search Input Header:**
     - Search icon (`Search` from `lucide-vue-next`).
     - Input field with placeholder and autofocus on open.
     - Clear query button (`✕` / `X` from `lucide-vue-next`) visible when text is present.
     - On desktop, keyboard badge: `<span class="binding">Ctrl + K</span>`.
     - On mobile, close button / trigger.
   - **Filter Tabs:**
     - Tabs: `[All]`, `[Dictionary]`, `[Study Decks]`.
     - Active tab style: black background (`bg-main text-white`). Inactive tab style: light or transparent (`text-side border border-light`).
   - **Realistic Mock Reactive Data:**
     - Based on official models (`WordSearchResult`, `WordDetail`, `StudioDeckSummary`).
     - Provide preloaded results for query `'guagua'` (as shown in reference images) and dynamic filtering when the query changes.
   - **Results Area:**
     - **Two-column layout on Desktop (`md:grid md:grid-cols-12`):**
       - Left column: **Featured Entry Card** (`col-span-5` or `col-span-4`) with lemma, part of speech, full region name (e.g. `Cuba`), short definition, and "See full entry →" button.
       - Right column: **Live Results List** (`col-span-7` or `col-span-8`):
         - "LIVE DICTIONARY RESULTS" header with count (e.g. `05`).
         - Interactive result items showing lemma, full region name (e.g. `Cuba`, `Chile`), definition, and hover/active highlight.
         - "MATCHING STUDY DECKS" section showing matching decks with their full region name (e.g. `Caribbean Everyday Words` with `Cuba`, `Chilean Everyday Words` with `Chile`).
     - **Mobile Layout:** Stacked order matching the mobile reference screenshot.
   - **Footer:**
     - Top border (`border-t border-light`).
     - Summary count: e.g. `See all results for '{query}'`.
     - Keyboard navigation hints using the `binding` shortcut:
       - `<div class="flex items-center gap-2"><span class="binding">↑</span><span class="binding">↓</span><span class="text-xs text-black">Navigate</span></div>`
       - `<div class="flex items-center gap-2"><span class="binding">Enter</span><span class="text-xs text-black">Select</span></div>`
       - `<div class="flex items-center gap-2"><span class="binding">Esc</span><span class="text-xs text-black">Close</span></div>`
       - (On mobile, hide or simplify keyboard hints).
   - **Keyboard Navigation:**
     - `ArrowDown` and `ArrowUp` cycle through active search items.
     - `Enter` selects the highlighted item and closes the modal.
     - `Escape` closes the modal.

2. **Mount in `AppLayout.vue` and Wire `TopBar.vue`:**
   - Mount `<SearchModal />` in `app/frontend/src/layouts/AppLayout.vue`.
   - Wire the Search icon button in `app/frontend/src/components/TopBar.vue` to open the modal.

3. **Validation:**
   - Run `cd app && wails3 task check`. Must exit with code 0.

4. **Change log:**
   - Document all changes in `.tickets/013-build-search-modal/changelog.md`.

## Verification

All checks are listed in [spec.md](spec.md).

