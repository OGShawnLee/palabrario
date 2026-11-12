# 026 — Spec checklist

Run by the router before accepting the ticket. All items must pass.

- [x] `IAuthAPI.getPublicProfileByUsername(username)` declared in `app/frontend/src/api/contract/auth.contract.ts` and implemented in both `WailsAuthAdapter` (`app/frontend/src/api/live/wails/auth.wails.ts`) and `MockAuthAPI` (`app/frontend/src/api/mock/mock.adapter.ts`).
- [x] `app/frontend/src/views/ProfileView.vue` exists, guards empty usernames, loads data via `api.auth.getPublicProfileByUsername(...)`, updates `breadcrumbLeafLabel` to `@username`, and resets it on unmount.
- [x] Loading (`"Loading profile..."`) and not-found (`"Profile not found."` + `<RouterLink to="/dictionary">`) fallback states render cleanly.
- [x] Route `/u/:username` registered in `app/frontend/src/router/index.ts` with no `alias` and `Home / Profile / User` breadcrumb metadata.
- [x] `ProfileHeader.vue` renders `profile.name`, square avatar beside `@${profile.user_name}`, role badge, member-since date, and resolves `anchor_dialect` to its full English region name via `REGIONS` in `app/frontend/src/data/regions.ts` (zero raw dialect codes, zero `flag_emoji`, zero invented bio/quotes/member numbers).
- [x] `ProfileStatsStrip.vue` renders the inverted Level/XP monolith (`total_xp % 1000` progress to next level), Consistency (`current_streak`, `longest_streak`), and The Arena (`matches_won`, `podium_count`) using `tabular-nums` and mobile-first breakpoints.
- [x] `ProfileBadgesSection.vue` filters `profile.unlocked_badges` strictly to `is_unlocked === true` and renders monochrome tier pills (`badge-main` / `badge-side` for `Bronze`, `Silver`, `Gold`, `Platinum`, `Diamond`) with zero color classes.
- [x] `ProfileLibrarySection.vue` renders public collections (first collection in inverted black monolith, remaining stacked, linking to `/collections/:id`) and public decks (`item_count` cards, linking to `/decks`) using real DTO fields only.
- [x] Owner-only strip in `ProfileView.vue` renders `email`, `private_decks_count`, `due_reviews_today`, `"Account Settings"` (`→ /settings`), and `"Open Studio"` (`→ /studio`) strictly when `profile.is_owner === true`.
- [x] `TopBarDesktopStatus.vue` and `TopBarMobileMenu.vue` include a `"Profile"` link to `/u/${user.user_name}` when signed in.
- [x] Design hygiene: zero `font-mono` occurrences (Satoshi font only), sharp `rounded-none` geometry (zero shadows), and natural sentence-case section/eyebrow labels (no forced `uppercase` except `XP`).
- [x] Typecheck and build pass: `wails3 task check` in `app/` exits with code 0.
- [x] No application code outside declared scope, backend code, or `.context/` modified.
- [x] `.tickets/026-build-profile-route/changelog.md` exists and is populated.
- [x] Commit (if made): Conventional Commits, English, detailed body.
