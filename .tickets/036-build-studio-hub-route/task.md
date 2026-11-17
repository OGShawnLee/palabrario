# 036 — Build Studio Hub Route (`/studio`)

| Field | Value |
| :--- | :--- |
| **ID** | 036 |
| **Status** | completed |
| **Model** | Gemini 4 Hydrogen, `thinking_level: medium` (`gemini-4-hydrogen`) |
| **Role** | Frontend Implementer (`implementation`) |
| **Spec** | [spec.md](spec.md) |

> **Specification Note:**
> Ticket 2 of 7 in the active Studio suite (`035`–`038`, `043`, `041`–`042`; depends on `035`). Replaces the `StudioView.vue` placeholder on `/studio` with a calm, spacious monochrome Swiss-broadsheet Studio Hub for browsing and filtering owned & saved **Decks** and **Collections** (`All` / `Owned` / `Saved` on both sections), reusing [`DeckGridItem.vue`](../../app/frontend/src/components/decks/DeckGridItem.vue) and [`CollectionGridItem.vue`](../../app/frontend/src/components/collections/CollectionGridItem.vue) directly with inline UnoCSS utilities (`uno.config.ts` untouched) and `Home / Studio` breadcrumbs. Temporary disabled `New Deck` / `New Collection` buttons are omitted here and added directly in `037` / `038` alongside their dialogs.

## Scope

May create or edit:
- `app/frontend/src/views/StudioView.vue`
- `app/frontend/src/components/decks/DeckGridItem.vue` (parameterize cleanly for reuse in both `/decks` and `/studio`)
- `app/frontend/src/components/collections/CollectionGridItem.vue` (parameterize cleanly for reuse in both `/collections` and `/studio`)
- `app/frontend/src/router/index.ts` (add `breadcrumb` meta for `studio`; keep `requiresAuth: true`)
- `.tickets/036-build-studio-hub-route/changelog.md` (update)

Must not touch:
- `core-service/`, `realtime-service/`, `image-service/`, `pkg/`, `app/*.go`, `SCHEMA_FINAL.sql`, `sql/`
- `app/frontend/uno.config.ts` (stick to inline UnoCSS utilities + existing shortcuts)
- `app/frontend/src/api/` (already prepared in `035`)
- Any file outside the above scope
- **No role (Router, Investigator, Implementer, Translator, or Reviewer) may commit, push, or merge without explicit user authorization.** Do not run `git commit`, `git push`, or `git merge`.

## Inputs

- [.routing/CONTEXT.md](../../.routing/CONTEXT.md), [.routing/DESIGN.md](../../.routing/DESIGN.md), [.routing/ROUTES.md](../../.routing/ROUTES.md), [.routing/API.md](../../.routing/API.md).
- [.routing/services/studio_service.md](../../.routing/services/studio_service.md) (`RN03`, `RN04`, `RN05`, `StudioDeckSummary`, `Collection`).
- [app/frontend/uno.config.ts](../../app/frontend/uno.config.ts).
- [app/frontend/src/views/StudioView.vue](../../app/frontend/src/views/StudioView.vue).
- [app/frontend/src/components/decks/DeckGridItem.vue](../../app/frontend/src/components/decks/DeckGridItem.vue) and [app/frontend/src/components/collections/CollectionGridItem.vue](../../app/frontend/src/components/collections/CollectionGridItem.vue).
- [app/frontend/src/components/profile/ProfileLibrarySection.vue](../../app/frontend/src/components/profile/ProfileLibrarySection.vue) (inline UnoCSS broadsheet styling reference).
- [app/frontend/src/composables/useAuth.ts](../../app/frontend/src/composables/useAuth.ts) (`currentUser.id_user`).

## Task

1. **Ensure Studio breadcrumb & auth gate in [router/index.ts](../../app/frontend/src/router/index.ts):**
   - Keep `meta: { requiresAuth: true }` on `/studio` and add `breadcrumb: [{ label: "Home", to: "/" }, { label: "Studio" }]`.

2. **Parameterize & refine [`DeckGridItem.vue`](../../app/frontend/src/components/decks/DeckGridItem.vue) and [`CollectionGridItem.vue`](../../app/frontend/src/components/collections/CollectionGridItem.vue) for calm broadsheet density:**
   - **Backward-compatible public defaults:** Preserve existing `/decks` and `/collections` rendering when optional Studio props/slots are omitted.
   - **Avoid badge clutter (no multi-box stacks):** Do **not** stack multiple bordered `.badge` pills at the top and bottom of each card. In Studio mode:
     - Render ownership, visibility, official/featured status, term/deck counts, and `@owner_user_name` (when saved) as a **single quiet typographic metadata line** (`text-xs font-bold text-dull`, separated by ` · `, e.g. `Owned · Private · 3 terms` or `Saved · Official · 48 terms · @valentina_ospina`), keeping at most one subtle accent only if genuinely helpful (or pure clean typography matching public cards).
     - In `CollectionGridItem.vue`, check `Array.isArray(collection.decks)` before rendering `"X decks"` so live `ListStudioCollections` items (where `decks` is omitted) do **not** falsely show `"0 decks"`—fall back cleanly to `Official` / `Curated` instead.
     - Render up to 3 full `region_names` as quiet inline muted text (`text-xs text-dull`, e.g. `Argentina · Uruguay`) with generous top spacing rather than bordered `badge-side` boxes.
   - **Generous card breathing room:** Increase internal vertical rhythm (`pb-8`, generous spacing around the title `mt-2.5 mb-2.5` and metadata/regions `mt-4`) so cards feel airy rather than cramped.
   - **Single bordered CTA per card:** Keep `Open →` as the primary `button-side button--small` link (`to="/studio/decks/:id"` / `"/studio/collections/:id"`), and let `#actions` host quiet secondary actions (`Unsave`).
   - **Ground Rules:** Strictly follow [DESIGN.md](../../.routing/DESIGN.md). Never render `flag_emoji` or dialect codes (`MX`, `CU`)—always full region names. Never render SM-2 mastery fields (`due_cards_today`, `retention_percent`) or `/studio/study` links. Never import `@bindings/palabrario/app/app`.

3. **Refine [StudioView.vue](../../app/frontend/src/views/StudioView.vue) whitespace & action hierarchy:**
   - Load `api.studio.listStudioDecks()` and `api.studio.listStudioCollections()` in parallel with generous Swiss-broadsheet whitespace (`max-w-6xl mx-auto px-4 pt-10 pb-24 flex flex-col gap-16`), wider grid gaps (`gap-y-12 md:gap-x-16`), and spacious section header dividers (`pb-5 mb-8`).
   - Editorial header with concise Workspace title/intro and subtle fallback links to `/decks` and `/collections` (no temporary disabled `New Deck` / `New Collection` placeholder buttons).
   - **Section `01 / Your Decks`:** `All` / `Owned` / `Saved` filter tabs + grid using `DeckGridItem.vue` linking to `/studio/decks/:id`, with a quiet text-style `Unsave` control (`text-xs font-medium text-dull hover:text-main px-2 py-1 min-h-10 inline-flex items-center cursor-pointer disabled:opacity-50`) inside `#actions` so it doesn't compete with `Open →`.
   - **Section `02 / Your Collections`:** `All` / `Owned` / `Saved` filter tabs (ownership derived from `collection.id_user === currentUser?.id_user`) + grid using `CollectionGridItem.vue` linking to `/studio/collections/:id`, with the same quiet `Unsave` text action inside `#actions`.
   - Provide clear loading state, retryable `role="alert"` error banner, and distinct empty states per tab/section linking to `/decks` and `/collections`.

4. **Verify:**
   - Run `wails3 task check` in `app/` and confirm exit code `0`.

## Verification

- Run all checks in [spec.md](spec.md).

## Change log

- Update `.tickets/036-build-studio-hub-route/changelog.md` recording the whitespace/de-crowding pass, `CollectionGridItem.vue` omitted-`decks` fix, and `wails3 task check` output.
