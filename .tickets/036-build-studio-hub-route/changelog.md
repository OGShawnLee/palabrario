# Changelog — 036: Build Studio Hub Route (`/studio`)

## Summary

Replaced the `StudioView.vue` placeholder on `/studio` with the monochrome Swiss-broadsheet Studio Hub for browsing, filtering (`All` / `Owned` / `Saved`), and unsaving owned & saved **Decks** and **Collections**, reusing parameterized `DeckGridItem.vue` and `CollectionGridItem.vue` components and adding the `Home / Studio` breadcrumb in `router/index.ts`.

---

## Files Modified & Created

### Routing (`app/frontend/src/router/index.ts`)
- **Modified** `app/frontend/src/router/index.ts`
  - Added `breadcrumb: [{ label: "Home", to: "/" }, { label: "Studio" }]` to the `/studio` route while preserving `requiresAuth: true`.

### Reusable Grid Cards (`app/frontend/src/components/`)
- **Modified** `app/frontend/src/components/decks/DeckGridItem.vue`
  - Parameterized `deck` prop to accept `Deck | StudioDeckSummary`.
  - Added optional `to` (defaults to `"/studio"`), `actionLabel` (defaults to `"Explore →"`), `ownership` (`"owned" | "saved" | null`), `showVisibility` (`boolean`), and `showRegions` (`boolean`) props, plus an optional `#actions` slot.
  - Preserved exact backward-compatible rendering on `/decks` when optional Studio props/slots are omitted.
  - In Studio mode, replaced multi-badge stacks with a single quiet typographic metadata line (`text-xs font-bold text-dull tabular-nums`, separated by ` · `, e.g. `Owned · Private · 3 terms` or `Saved · Public · Official · 48 terms · @valentina_ospina`) and renders up to 3 full `region_names` as quiet inline muted text (`text-xs text-dull`, joined by ` · `) with generous internal vertical rhythm (`pb-8`, `mt-2.5 mb-2.5` around the title, `mt-4` before regions). Never renders `flag_emoji`, dialect codes, or SM-2 mastery fields (`due_cards_today`, `retention_percent`).
- **Modified** `app/frontend/src/components/collections/CollectionGridItem.vue`
  - Added optional `to` (defaults to `/collections/${collection.id_collection}`), `actionLabel` (defaults to `"Explore →"`), `ownership` (`"owned" | "saved" | null`), `showVisibility` (`boolean`), and `showRegions` (`boolean`) props, plus an optional `#actions` slot.
  - Preserved exact backward-compatible rendering on `/collections` when optional Studio props/slots are omitted.
  - In Studio mode, replaced multi-badge stacks with a single quiet typographic metadata line (`Owned`/`Saved`, `Public`/`Private`, `Official`, `Featured`) and up to 3 full `region_names` as quiet inline muted text (`text-xs text-dull`) with generous card breathing room (`pb-8`, `mt-2.5 mb-2.5`, `mt-4`).
  - **Omitted-`decks` fallback fix:** Checks `Array.isArray(collection.decks)` before rendering `"X decks"` so `ListStudioCollections` items (where `decks` is omitted) do not falsely display `"0 decks"` and instead fall back cleanly to `Official` / `Curated`.

### Studio Hub View (`app/frontend/src/views/StudioView.vue`)
- **Modified** `app/frontend/src/views/StudioView.vue`
  - Loads `api.studio.listStudioDecks()` and `api.studio.listStudioCollections()` in parallel via `@/api` with race-safe request tracking, loading state, and a retryable `role="alert"` error banner.
  - Applies generous Swiss-broadsheet spacing (`max-w-6xl mx-auto px-4 pt-10 pb-24 flex flex-col gap-16`), spacious section header dividers (`gap-6 border-b border-light pb-5 mb-8`), and wider grid gaps (`gap-y-12 md:grid-cols-2 md:gap-x-16`).
  - Renders the editorial Workspace header (`pb-10`) with subtle fallback links to `/decks` and `/collections` (without temporary disabled `New Deck` / `New Collection` placeholder buttons).
  - **Section `01 / Your Decks`:** `All` / `Owned` / `Saved` accessible `role="tablist"` filters with counts; renders `DeckGridItem.vue` linking to `/studio/decks/${deck.id_deck}`, with a quiet text-style `Unsave` control (`text-xs font-medium text-dull hover:text-main px-2 py-1 min-h-10 inline-flex items-center cursor-pointer disabled:opacity-50`) inside `#actions` so `Open →` remains the single bordered CTA per card.
  - **Section `02 / Your Collections`:** `All` / `Owned` / `Saved` accessible `role="tablist"` filters (ownership derived from `collection.id_user === currentUser?.id_user`); renders `CollectionGridItem.vue` linking to `/studio/collections/${collection.id_collection}`, with the same quiet text-style `Unsave` action inside `#actions`.
  - Provides distinct empty states (`p-8`) per tab and section linking to `/decks` and `/collections`, plus a dismissible `role="alert"` banner if an unsave call fails.

---

## Verification & Command Outputs

### 1. Wails Pipeline (`wails3 task check`)
Command (from `app/`): `wails3 task check`
Result: **PASS** (exit code `0`)
```text
task: [bindings] wails3 generate bindings -clean=true -ts
Wails v3.0.0 › Generate Bindings
 INFO  Processed: 264 Packages, 1 Service, 98 Methods, 6 Enums, 75 Models, 0 Events in 5.6590885s.
 INFO  Output directory: C:\Users\OGSha\OneDrive\Escritorio\Developer\Web Development\Palabrario\app\frontend\bindings
task: [check] bun run --cwd frontend type-check
$ vue-tsc --noEmit
task: [check] bun run --cwd frontend build
$ vite build --mode production
vite v8.3.5 building client environment for production...
transforming...
✓ 2503 modules transformed.
rendering chunks...
computing gzip size...
...
dist/assets/DeckGridItem-BjoJ7JDo.js                1.95 kB │ gzip:  0.98 kB
dist/assets/CollectionGridItem-rb9eJewu.js          2.05 kB │ gzip:  0.99 kB
dist/assets/StudioView-pprlyhrw.js                  8.74 kB │ gzip:  2.79 kB
...
✓ built in 8.06s
task: [check] go build .
```

### 2. Ground Rules & Scope Checks
- `Select-String -Path "src\**\*.vue" -Pattern "@bindings/palabrario/app/app" -SimpleMatch` (from `app/frontend/`): only pre-existing `src\components\WindowBar.vue` matched; zero matches in `StudioView.vue`, `DeckGridItem.vue`, `CollectionGridItem.vue`, or any other view.
- `Select-String -Path "src\views\StudioView.vue","src\components\decks\DeckGridItem.vue","src\components\collections\CollectionGridItem.vue","src\router\index.ts" -Pattern "flag_emoji|due_cards_today|retention_percent|/studio/study|@bindings/palabrario/app/app"`: zero matches.
- `app/frontend/uno.config.ts`: untouched.
