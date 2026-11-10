# 019 — Spec checklist

Run by the router before accepting the ticket. All items must pass.

- [ ] Route `/arena/history` registered in `app/frontend/src/router/index.ts` under `AppLayout`.
- [ ] `ArenaHistoryView.vue` consumes match history from `api.arena.listMatchHistory()` via `import { api } from "@/api"`. Zero hardcoded mock match history arrays in views.
- [ ] `"View Replay"` triggers `MatchReplayDialog.vue` which fetches detail via `api.arena.getMatchReplay(matchId)`. Zero hardcoded round replays in the component.
- [ ] `IArenaAPI` contract and adapters (Wails and Mock) properly support `listMatchHistory` and `getMatchReplay`.
- [ ] Replay modal shows podium summary and round-by-round replay with participant responses and timings.
- [ ] Zero flag emojis, full country names used.
- [ ] Monochrome design tokens strictly observed.
- [ ] Typecheck passes: `cd app && wails3 task check` exits with code 0.
- [ ] `.tickets/019-build-arena-history-route/changelog.md` exists and is populated.
