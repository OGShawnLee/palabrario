# 026 — Build the Public Profile Route (/u/:username)

| Field | Value |
| :--- | :--- |
| **ID** | 026 |
| **Status** | completed |
| **Model** | Claude Sonnet 5.5, `effort: medium` (`claude-sonnet-5-5`) |
| **Role** | Frontend Implementer |
| **Spec** | [spec.md](spec.md) |

> **Specification Note:**  
> This ticket builds the `/u/:username` public profile route (`ProfileView.vue`) in Palabrario’s editorial Swiss-broadsheet monochrome style, wires `getPublicProfileByUsername(username)` across `IAuthAPI` (Wails and Mock adapters), displays the masthead (with square avatar beside the `@username` handle), `01 / The practice` progression ledger (inverted Level/XP monolith with exact `total_xp % 1000` level bar, Consistency, and The Arena), `02 / Milestones` unlocked achievements (`is_unlocked === true` only with temporary monochrome tier badges), `03 / The public library` (featured collection monolith + authored collections/decks), the owner-only footer strip (`is_owner === true` linking to `/settings` and `/studio`), and `"Profile"` links in `TopBarDesktopStatus.vue` and `TopBarMobileMenu.vue`.

## Scope

May create or edit:
- `app/frontend/src/api/contract/auth.contract.ts` (add `getPublicProfileByUsername(username: string): Promise<PublicUserProfile | null>`)
- `app/frontend/src/api/live/wails/auth.wails.ts` (implement `getPublicProfileByUsername` via `Wails.GetPublicProfileByUsername(username)`)
- `app/frontend/src/api/mock/fixtures/user.fixture.ts` (add mock `AchievementStatus[]` and `mockArenaStatsByUser` fixture data for `mockUsers`)
- `app/frontend/src/api/mock/mock.adapter.ts` (implement `MockAuthAPI.getPublicProfileByUsername(username)` with case-insensitive lookup and `is_owner` enrichment)
- `app/frontend/src/views/ProfileView.vue` (new: public profile view + bottom owner-only strip)
- `app/frontend/src/components/profile/ProfileHeader.vue` (new: broadsheet masthead with square avatar beside `@username`, role badge, member-since date, and full `anchor_dialect` region column)
- `app/frontend/src/components/profile/ProfileStatsStrip.vue` (new: `01 / The practice` 3-column ledger — inverted Level/XP monolith + Consistency + The Arena)
- `app/frontend/src/components/profile/ProfileBadgesSection.vue` (new: `02 / Milestones` unlocked-only achievements grid with hairline dividers and temporary monochrome tier labels)
- `app/frontend/src/components/profile/ProfileLibrarySection.vue` (new: `03 / The public library` — inverted black monolith for the first collection, stacked remaining collections, and hairline-divided public study decks strip)
- `app/frontend/src/components/TopBarDesktopStatus.vue` (add `"Profile"` link to `/u/${user.user_name}`)
- `app/frontend/src/components/TopBarMobileMenu.vue` (add `"Profile"` link to `/u/${user.user_name}`)
- `app/frontend/src/router/index.ts` (register `u/:username` route with no `alias` and `Home / Profile / User` breadcrumb)
- `.tickets/026-build-profile-route/changelog.md` (create)

Must not touch:
- `app/frontend/uno.config.ts` (use existing tokens/shortcuts; record any proposals in changelog)
- `app/*.go`, `core-service/`, `realtime-service/`, `image-service/`, or `.context/`
- Any file outside the above scope

## Inputs

- [.routing/CONTEXT.md](../../.routing/CONTEXT.md): Ground rules (never render `flag_emoji` or dialect codes; full region names only; Satoshi font only; no route aliases; mobile-first responsive breakpoints).
- [.routing/DESIGN.md](../../.routing/DESIGN.md): Monochrome palette (`main`, `side`, `dull`, `light`, `ground`, `white`), sharp geometry (`rounded-none`, zero drop shadows), natural sentence-case section labels (`text-sm font-bold text-dull`, avoid forced all-caps except acronyms like `XP`), `badge-main`, `badge-side`, `button-side`.
- [.routing/services/identity_service.md](../../.routing/services/identity_service.md): `PublicUserProfile` DTO contract and `CU01` / `RN01` rules (`GET /api/v1/users/u/{username}`). Note: `PublicUserProfile` has **no** bio, quote, member number, collection term count, or collection/deck region fields — render only real DTO properties.
- [.routing/services/mastery_service.md](../../.routing/services/mastery_service.md) and [sql/queries/mastery.sql](../../sql/queries/mastery.sql#L180-L185): `AchievementStatus` DTO contract and backend level formula (`current_level = (total_xp / 1000) + 1`, so progress within the current level is `total_xp % 1000` out of `1,000 XP`).
- [app/frontend/src/data/regions.ts](../../app/frontend/src/data/regions.ts): `REGIONS` map for translating `anchor_dialect` codes into full English region names.
- [app/frontend/src/views/CollectionDetailView.vue](../../app/frontend/src/views/CollectionDetailView.vue): Reference for `route.params` watching, `isLoading` + `null` fallback, and `breadcrumbLeafLabel` lifecycle.

## Task

1. **Extend `IAuthAPI` and Both Adapters:**
   - In `app/frontend/src/api/contract/auth.contract.ts`:
     - Import `PublicUserProfile` from `@bindings/palabrario/pkg/models/index.js`.
     - Add `getPublicProfileByUsername(username: string): Promise<PublicUserProfile | null>` to `IAuthAPI`.
   - In `app/frontend/src/api/live/wails/auth.wails.ts`:
     - Implement `getPublicProfileByUsername(username: string)` delegating to `Wails.GetPublicProfileByUsername(username)`.
   - In `app/frontend/src/api/mock/fixtures/user.fixture.ts`:
     - Export `mockAchievements: AchievementStatus[]` containing at least 4 entries spanning `is_unlocked: true` and `is_unlocked: false` across tiers (`BRONZE`, `SILVER`, `GOLD`, `DIAMOND`) with authentic Spanish achievement names (`"Voz de la comunidad"`, `"Constancia pura"`, `"Primera victoria"`, etc.).
     - Export `mockArenaStatsByUser: Record<string, { matches_won: number; podium_count: number }>`.
   - In `app/frontend/src/api/mock/mock.adapter.ts`:
     - Implement `MockAuthAPI.getPublicProfileByUsername(username: string)`:
       - Trim and strip a leading `@`, return `null` if empty.
       - Match against `mockUsers` case-insensitively (`candidate.user_name.toLowerCase() === normalized`). Return `null` if not found.
       - Compute `isOwner = this.session?.id_user === user.id_user`.
       - Populate `public_decks` from `mockDeckDetails.filter((d) => d.is_public && d.id_user === user.id_user)` and `public_collections` from `mockCollectionDetails.filter((c) => c.is_public && c.id_user === user.id_user)`.
       - Populate `email`, `private_decks_count`, and `due_reviews_today` only when `isOwner === true`; omit or set `null` when `isOwner === false`.

2. **Build Profile Components (`app/frontend/src/components/profile/`):**
   > **Design Hygiene:** Do **not** invent a static bio/quote or member catalog number (`Nº 0248`), do **not** render forced all-caps eyebrows (use natural sentence case such as `01 / The practice`, `Anchor dialect`, `Only you can see this`, keeping `XP` uppercase), and keep all corners sharp (`rounded-none`, zero drop shadows). Write mobile-first base styles and scale up at `md:`.
   - **`ProfileHeader.vue`:**
     - Props: `{ profile: PublicUserProfile }`.
     - Broadsheet masthead (`border-b border-light pb-8 grid grid-cols-1 gap-6 md:grid-cols-12 md:items-end`):
       - **Left (`md:col-span-8 flex flex-col gap-5`):**
         - Oversized display `<h1>` (`text-4xl sm:text-6xl md:text-7xl font-black tracking-tight text-main leading-none`) rendering `profile.name`.
         - Identity row beneath `<h1>` (`flex flex-wrap items-center gap-3`):
           - Sharp square avatar beside the handle (`h-10 w-10 bg-light border border-light flex items-center justify-center text-sm font-bold text-main shrink-0 overflow-hidden`) rendering `<img>` when `profile.avatar_url` is set or `profile.name.charAt(0).toUpperCase()`.
           - `@${profile.user_name}` (`text-sm font-medium text-main`), role badge (`badge-main` for `ADMIN`/`EDITOR`, `badge-side` for `LEARNER`, formatted as `"Learner"`, `"Editor"`, `"Admin"`), and `"Member since ${formattedDate}"` (`text-xs text-dull`).
       - **Right (`md:col-span-4 border-t border-light pt-4 md:border-t-0 md:border-l md:pt-0 md:pl-8 flex items-end justify-between gap-4`):**
         - Sentence-case label `Anchor dialect` (`text-sm font-bold text-dull`) above the full English region name resolved via `REGIONS` in `@/data/regions` (fallback `"Unknown region"`, `text-xl font-bold text-main` — never render raw codes or `flag_emoji`), paired with a `Globe` Lucide icon (`w-5 h-5 text-main shrink-0`).
   - **`ProfileStatsStrip.vue`:**
     - Props: `{ profile: PublicUserProfile }`.
     - Section header (`mb-6`): sentence-case label `01 / The practice` (`section-label`) and heading `A little, every day.` (`text-2xl sm:text-3xl font-bold text-main tracking-tight`).
     - 3-column broadsheet ledger (`grid grid-cols-1 gap-6 md:grid-cols-12 md:gap-0 border-b border-light pb-10`):
       - **Column 1 — Level & XP Monolith (`md:col-span-4 bg-main text-white p-6 flex flex-col justify-between gap-6`):**
         - Sentence-case micro-label `Language in progress` (`text-xs font-medium text-white/70`).
         - `Level` label, oversized `profile.current_level` (`text-5xl font-black tabular-nums leading-none`) paired on the right with `${profile.total_xp.toLocaleString()} XP` (`text-lg font-bold tabular-nums`).
         - Real progress bar to next level (`const xpIntoLevel = profile.total_xp % 1000`, `const xpPercent = Math.round((xpIntoLevel / 1000) * 100)`): sharp 2px bar (`h-0.5 w-full bg-white/25` with inner `bg-white` fill at `${xpPercent}%`) + caption `${xpIntoLevel.toLocaleString()} / 1,000 XP to Level ${profile.current_level + 1}` and `${xpPercent}%` (`tabular-nums`).
       - **Column 2 — Consistency (`md:col-span-4 border-b border-light pb-6 md:border-b-0 md:border-r md:px-8 md:py-2 flex flex-col justify-between gap-4`):**
         - Top row: `Consistency` (`text-sm font-bold text-dull`) + `Flame` Lucide icon (`w-4 h-4 text-main`).
         - Main stat: oversized `profile.current_streak` (`text-5xl font-bold text-main tabular-nums leading-none`) with inline `day streak` and subtext `Current streak`.
         - Hairline footer (`border-t border-light pt-3 flex items-baseline justify-between`): `Personal best` (`text-xs text-dull`) and `${profile.longest_streak} days` (`text-lg font-bold text-main tabular-nums`).
       - **Column 3 — The Arena (`md:col-span-4 md:pl-8 md:py-2 flex flex-col justify-between gap-4`):**
         - Top row: `The Arena` (`text-sm font-bold text-dull`) + `Trophy` Lucide icon (`w-4 h-4 text-main`).
         - Main stat: oversized `profile.matches_won` (`text-5xl font-bold text-main tabular-nums leading-none`) with inline `wins` and subtext `Words under pressure.`.
         - Hairline footer (`border-t border-light pt-3 flex items-baseline justify-between`): `Podium finishes` (`text-xs text-dull`) and `${profile.podium_count}` (`text-lg font-bold text-main tabular-nums`).
   - **`ProfileBadgesSection.vue`:**
     - Props: `{ badges: AchievementStatus[] }`.
     - Filter strictly to unlocked achievements: `const achieved = computed(() => props.badges.filter((b) => b.is_unlocked))`.
     - Header row (`flex items-baseline justify-between gap-4 mb-6`): `02 / Milestones` (`section-label`) + `Earned, not given.` (`text-2xl sm:text-3xl font-bold text-main tracking-tight`), and right-aligned `${achieved.length} unlocked` (`text-xs text-dull tabular-nums`).
     - Empty state when `achieved.length === 0`: `"No achievements unlocked yet."` in a quiet `border border-light bg-white p-6 text-sm text-dull`.
     - Populated grid (`grid grid-cols-1 divide-y divide-light md:grid-cols-3 md:divide-y-0 md:divide-x border-t border-b border-light`):
       - Each item (`py-6 md:px-6 md:first:pl-0 md:last:pr-0 flex flex-col justify-between gap-3`):
         - Top row: category Lucide icon (`Award`, `Flame`, or `Trophy`, `w-5 h-5 text-main`) + temporary monochrome tier pill (`Gold`/`Platinum`/`Diamond` → `badge-main`; `Bronze`/`Silver` → `badge-side`, sentence-cased, zero color classes).
         - Title `badge.name` (`text-lg font-bold text-main`), `badge.description` (`text-sm text-side`).
         - Footer (`text-xs text-dull tabular-nums`): `Unlocked ${formattedDate} · +${badge.xp_reward} XP`.
   - **`ProfileLibrarySection.vue`:**
     - Props: `{ collections: Collection[]; decks: Deck[]; authorName: string }`.
     - Header row (`flex flex-col sm:flex-row sm:items-baseline justify-between gap-2 mb-6`): `03 / The public library` (`section-label`) + `Words worth keeping.` (`text-2xl sm:text-3xl font-bold text-main tracking-tight`), and right-aligned `Selected and authored by ${firstName}` (`text-xs text-dull`).
     - **Public Collections:**
       - Empty state when `collections.length === 0`: `"No public collections yet."` in `border border-light bg-white p-6 text-sm text-dull mb-8`.
       - Populated (`grid grid-cols-1 gap-8 md:grid-cols-12 mb-10`):
         - First item (`collections[0]`, `collections.length === 1 ? 'md:col-span-12' : 'md:col-span-6'`): Inverted black monolith `<RouterLink :to="`/collections/${col.id_collection}`" class="bg-main text-white p-6 sm:p-8 flex flex-col justify-between gap-8 min-h-60">` showing `Collection / 01` + `ArrowUpRight` icon, `col.title` (`text-2xl sm:text-3xl font-bold leading-tight`), optional `col.description` (`text-sm text-white/75`), and bottom hairline (`border-t border-white/20 pt-4 flex items-center justify-between text-xs`) with `col.is_official ? "Official collection" : "Community collection"` and `Explore collection`.
         - Remaining items (`v-if="collections.length > 1"`, `collections.slice(1)`, `md:col-span-6 divide-y divide-light border-t border-b border-light`): each `<RouterLink :to="`/collections/${col.id_collection}`">` renders index `Collection / 0${index + 2}`, `col.title` (`text-xl font-bold text-main`), optional `col.description` (`text-sm text-side line-clamp-2`), and `ArrowUpRight` icon.
     - **Public Study Decks:**
       - Empty state when `decks.length === 0`: `"No public decks yet."` in `border border-light bg-white p-6 text-sm text-dull`.
       - Populated (`border-t border-light pt-8 grid grid-cols-1 gap-6 md:grid-cols-12`):
         - Left intro (`md:col-span-4`): `Public study decks` (`text-sm font-bold text-dull mb-1`) and `Study decks curated by ${firstName}.` (`text-sm text-side`).
         - Right deck strip (`md:col-span-8 grid grid-cols-1 divide-y divide-light md:grid-cols-2 md:divide-y-0 md:divide-x`): each `<RouterLink to="/decks" class="py-4 md:py-0 md:px-6 md:first:pl-0 flex items-start justify-between gap-4">` shows `Layers` Lucide icon, `deck.title` (`text-base font-bold text-main`), `${deck.item_count} cards` (`tabular-nums text-xs text-dull`) + optional `Official` `badge-side`, and `ArrowUpRight` icon.

3. **Create `ProfileView.vue` (`app/frontend/src/views/ProfileView.vue`):**
   - Watch `() => route.params.username` with `{ immediate: true }`.
   - Normalize `const raw = String(route.params.username ?? "").trim().replace(/^@/, "")`.
   - Guard empty `raw`: set `profile.value = null`, `breadcrumbLeafLabel.value = null`, `isLoading.value = false`, and return without calling the API.
   - Otherwise call `await api.auth.getPublicProfileByUsername(raw)` in `try`, set `profile.value = null` in `catch`, and in `finally` set `breadcrumbLeafLabel.value = profile.value ? `@${profile.value.user_name}` : null` and `isLoading.value = false`.
   - Reset `breadcrumbLeafLabel.value = null` in `onUnmounted`.
   - Render inside `<main class="flex-1 w-full max-w-6xl mx-auto px-4 pt-6 pb-20 flex flex-col gap-12">`:
     - Loading (`v-if="isLoading"`): `"Loading profile..."`.
     - Not found (`v-else-if="!profile"`): `"Profile not found."` with `<RouterLink to="/dictionary">← Back to Dictionary</RouterLink>`.
     - Loaded (`v-else`): `ProfileHeader`, `ProfileStatsStrip`, `ProfileBadgesSection`, `ProfileLibrarySection`, and—when `profile.is_owner === true`—the bottom **Owner Strip** (`border-t border-light pt-6 flex flex-col gap-4 md:flex-row md:items-center md:justify-between`):
       - Label with `Lock` Lucide icon (`w-3.5 h-3.5`): `Only you can see this` (`text-xs font-bold text-dull`).
       - Owner metadata (`flex flex-wrap items-center gap-x-6 gap-y-2 text-sm text-main tabular-nums`): `profile.email`, `${profile.private_decks_count ?? 0} private decks`, `${profile.due_reviews_today ?? 0} reviews due today`.
       - Owner actions (`flex flex-wrap items-center gap-4 text-sm font-medium text-main`): `<RouterLink to="/settings">` (`Account Settings`) and `<RouterLink to="/studio">` (`Open Studio`), each with `ArrowUpRight` icon.

4. **Register Route & Update Navigation Menus:**
   - In `app/frontend/src/router/index.ts` under `AppLayout.children`, add (with no `alias`):
     ```ts
     {
       path: "u/:username",
       name: "profile",
       component: () => import("@/views/ProfileView.vue"),
       meta: {
         breadcrumb: [
           { label: "Home", to: "/" },
           { label: "Profile" },
           { label: "User" },
         ],
       },
     },
     ```
   - In `app/frontend/src/components/TopBarDesktopStatus.vue` (inside `PopoverContent` above `Settings`):
     - Add `<RouterLink :to="`/u/${user.user_name}`" class="text-sm font-medium text-side" active-class="font-bold text-main" @click="open = false">Profile</RouterLink>`.
   - In `app/frontend/src/components/TopBarMobileMenu.vue` (inside the logged-in `<nav>`):
     - Add `<RouterLink :to="`/u/${user.user_name}`" class="text-base font-medium text-side" active-class="font-bold text-main" @click="open = false">Profile</RouterLink>`.

5. **Validation:**
   - Run `wails3 task check` inside `app/` and confirm exit code 0.

6. **Change log:**
   - Document created/modified files and any proposals in `.tickets/026-build-profile-route/changelog.md`.

## Verification

All checks are listed in [spec.md](spec.md).
