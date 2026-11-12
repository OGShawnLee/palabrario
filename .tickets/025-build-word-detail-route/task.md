# 025 — Build the Word Detail Route (/dictionary/word/:id)

| **ID** | 025 |
| **Status** | completed |
| **Model** | Claude Sonnet 5.5, `effort: medium` (`claude-sonnet-5-5`) |
| **Role** | Frontend Implementer |
| **Spec** | [spec.md](spec.md) |

> **Specification Note:**  
> This ticket builds the `/dictionary/word/:id` route (`WordDetailView.vue`), which renders a complete dictionary entry for a specific word by its ID. It reuses the editorial broadsheet presentation from `WordOfTheDayCard.vue`, integrates dynamic breadcrumbs (`Home / Dictionary / [Lemma]`), and wires navigation from `SearchModal.vue`.

## Scope

May create or edit:
- `app/frontend/src/views/WordDetailView.vue` (new: word detail view component)
- `app/frontend/src/components/dictionary/WordOfTheDayCard.vue` (adapt with `eyebrow`, `headingTag`, `footerLink` props, and soften `Report Word` button to `button-side`)
- `app/frontend/src/components/SearchModal.vue` (navigate dictionary search results directly to `/dictionary/word/${item.id_word}`)
- `app/frontend/src/router/index.ts` (register `/dictionary/word/:id` route with breadcrumb hierarchy, no `alias`)
- `.tickets/025-build-word-detail-route/changelog.md` (create)

Must not touch:
- `app/frontend/uno.config.ts` (shortcuts already exist; record any proposals in changelog)
- `app/*.go`, backend services, or `.context/`
- Any file outside the above scope

## Inputs

- [.routing/CONTEXT.md](../../.routing/CONTEXT.md): Architecture ground rules (strict prohibition of `flag_emoji` and regional codes; full region names required; Satoshi font only; no route aliases).
- [.routing/DESIGN.md](../../.routing/DESIGN.md): Monochrome design tokens (`main`, `side`, `dull`, `light`, `ground`), typography (`Satoshi`), surfaces, and buttons.
- [app/frontend/src/api/index.ts](../../app/frontend/src/api/index.ts): Universal `api` singleton (`api.dictionary.getWordById(id)`).
- [app/frontend/src/components/dictionary/WordOfTheDayCard.vue](../../app/frontend/src/components/dictionary/WordOfTheDayCard.vue): Existing broadsheet word card layout.
- [app/frontend/src/views/DictionaryView.vue](../../app/frontend/src/views/DictionaryView.vue): Reference for `word` + `isLoading` fallback pattern (no separate `error` ref).
- [app/frontend/src/views/CollectionDetailView.vue](../../app/frontend/src/views/CollectionDetailView.vue): Canonical pattern for `route.params.id` watching and `breadcrumbLeafLabel` lifecycle.
- [app/frontend/src/components/ui/AppBreadcrumb.vue](../../app/frontend/src/components/ui/AppBreadcrumb.vue): `breadcrumbLeafLabel` reactive ref.

## Task

1. **Parameterize `WordOfTheDayCard.vue` (`app/frontend/src/components/dictionary/WordOfTheDayCard.vue`):**
   - Add optional props using `withDefaults(defineProps<...>(), { ... })`:
     - `eyebrow?: string` (defaults to `"Word of the Day"`): render in the top-left eyebrow `<span>` instead of the hardcoded `"Word of the Day"` string.
     - `headingTag?: "h1" | "h2"` (defaults to `"h2"`): replace the hardcoded `<h2>` lemma element with `<component :is="headingTag" class="text-5xl sm:text-6xl md:text-7xl font-black text-main leading-none mt-2">`.
     - `footerLink?: "entry" | "dictionary"` (defaults to `"entry"`):
       - When `"entry"`: keep the footer `<RouterLink :to="`/dictionary/word/${word.id_word}`" class="text-sm font-semibold text-main inline-flex items-center gap-1 hover:underline sm:ml-auto">` reading `See full entry` and `→`.
       - When `"dictionary"`: render `<RouterLink to="/dictionary" class="text-sm font-semibold text-main inline-flex items-center gap-1 hover:underline sm:ml-auto">` reading `←` and `Back to Dictionary`.
   - Soften the footer `Report Word` button unconditionally (both modes) by changing its class from `button-main` to `button-side`.

2. **Create `WordDetailView.vue` (`app/frontend/src/views/WordDetailView.vue`):**
   - Import `api` from `@/api`, `RouterLink` and `useRoute` from `vue-router`, and `breadcrumbLeafLabel` from `@/components/ui/AppBreadcrumb.vue`.
   - Manage only `word = ref<WordDetail | null>(null)` and `isLoading = ref(true)` (matching `DictionaryView.vue` and `CollectionDetailView.vue`; do not add a separate `error` ref).
   - Watch `() => route.params.id` with `{ immediate: true }`:
     - Parse `const id = Number(route.params.id)`.
     - Guard invalid IDs: if `!Number.isInteger(id) || id <= 0`, set `word.value = null`, `breadcrumbLeafLabel.value = null`, `isLoading.value = false`, and return without calling `api.dictionary.getWordById`.
     - Otherwise set `isLoading.value = true`, call `word.value = await api.dictionary.getWordById(id)` inside `try`, set `word.value = null` in `catch`, and in `finally` set `breadcrumbLeafLabel.value = word.value?.lemma ?? null` and `isLoading.value = false`.
   - Reset `breadcrumbLeafLabel.value = null` inside `onUnmounted`.
   - Render inside `<main class="flex-1 w-full max-w-6xl mx-auto px-4 pt-6 pb-20">`:
     - Loading (`v-if="isLoading"`): `<div class="border border-light bg-ground p-6 text-sm text-dull">Loading dictionary entry...</div>`.
     - Not found / error fallback (`v-else-if="!word"`): `<div class="border border-light bg-ground p-8 text-center">` containing `<p class="text-sm text-dull mb-4">Word not found.</p>` and `<RouterLink to="/dictionary" class="text-sm font-semibold text-main inline-flex items-center gap-1 hover:underline">← Back to Dictionary</RouterLink>` (a `RouterLink`, not a button).
     - Loaded (`v-else`): `<WordOfTheDayCard :word="word" eyebrow="Dictionary Entry" heading-tag="h1" footer-link="dictionary" />`.

3. **Register Route in `app/frontend/src/router/index.ts`:**
   - Under the `AppLayout` children array, register (with no `alias`):
     ```ts
     {
       path: "dictionary/word/:id",
       name: "word-detail",
       component: () => import("@/views/WordDetailView.vue"),
       meta: {
         breadcrumb: [
           { label: "Home", to: "/" },
           { label: "Dictionary", to: "/dictionary" },
           { label: "Word Detail" },
         ],
       },
     },
     ```

4. **Update Search Navigation in `SearchModal.vue` (`app/frontend/src/components/SearchModal.vue`):**
   - In `selectDictionaryResult(item: WordSearchResult)`:
     - Update navigation from `router.push({ path: "/dictionary", query: { q: item.lemma } })` to:
       ```ts
       open.value = false
       router.push(`/dictionary/word/${item.id_word}`)
       ```
     - Selecting a dictionary word result now navigates straight to that word's entry page.

5. **Validation:**
   - Run `cd app && wails3 task check` (or `wails3 task check` in `app/`) to verify TypeScript, Vite build, and Go build pass with exit code 0.

6. **Change log:**
   - Document created/modified files, prop signatures, and any proposals in `.tickets/025-build-word-detail-route/changelog.md`.

## Verification

All checks are listed in [spec.md](spec.md).
