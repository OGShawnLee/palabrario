# 026 — Change log

## Created
- `app/frontend/src/views/ProfileView.vue`
- `app/frontend/src/components/profile/ProfileHeader.vue`
- `app/frontend/src/components/profile/ProfileStatsStrip.vue`
- `app/frontend/src/components/profile/ProfileBadgesSection.vue`
- `app/frontend/src/components/profile/ProfileLibrarySection.vue`

## Modified
- `api/contract/auth.contract.ts`: added `getPublicProfileByUsername`.
- `api/live/wails/auth.wails.ts`: delegates to `Wails.GetPublicProfileByUsername`.
- `api/mock/fixtures/user.fixture.ts`: added `mockAchievements`, `mockArenaStatsByUser`, and `mockUnlockedCodesByUser` (per-user unlocked codes, an addition beyond the ticket so each mock profile unlocks a different set).
- `api/mock/mock.adapter.ts`: `MockAuthAPI.getPublicProfileByUsername` (case-insensitive, `@` stripped, owner-only fields `email`, `private_decks_count`, `due_reviews_today` set to null for non-owners; mock owner values 2 and 12 are placeholders).
- `router/index.ts`: `u/:username` route, no alias.
- `TopBarDesktopStatus.vue`, `TopBarMobileMenu.vue`: `Profile` link.

## Notes / proposals
- No `uno.config.ts` changes; existing `section-label`, `badge-main`, `badge-side` shortcuts were enough.
- Mock deck/collection owner ids do not match the mock user ids, so the mock library sections show their empty states. Aligning `OWNER_IDS` in `decks.fixture.ts` (out of scope) would populate them.
- Tier badges are temporary monochrome labels; permanent treatment is tracked in `PROGRESS.md` under "Medal coloring".

## Validation
- `wails3 task check` in `app/`: frontend build and `go build` completed with no errors.
