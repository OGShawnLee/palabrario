# 019 — Spec checklist

Run by the router before accepting the ticket. All items must pass.

- [ ] Route `/arena/history` registered in `app/frontend/src/router/index.ts` under `AppLayout`.
- [ ] `ArenaHistoryView.vue` renders match history table with room PIN, deck, format, winner, and replay button.
- [ ] `"View Replay"` triggers `MatchReplayDialog.vue` using `AppDialog.vue`.
- [ ] Replay modal shows podium summary and round-by-round replay with participant responses and timings.
- [ ] Zero flag emojis, full country names used.
- [ ] Monochrome design tokens strictly observed.
- [ ] Typecheck passes: `cd app && wails3 task check` exits with code 0.
- [ ] `.tickets/019-build-arena-history-route/changelog.md` exists and is populated.
