# 019 — Build the Arena Match History & Replay Route (/arena/history)

| Field | Value |
| :--- | :--- |
| **ID** | 019 |
| **Status** | draft |
| **Model** | Sonnet 5.5 (`claude-sonnet-5-5`, Antigravity provider) |
| **Role** | Frontend Implementer |
| **Spec** | [spec.md](spec.md) |

> **Draft Note:**  
> This ticket builds the dedicated match history route (`/arena/history`) and the 3NF round-by-round match replay inspection dialog (`MatchReplayDialog.vue`), satisfying Use Case `CU16` and enabling players to review past matches, scores, podium finishes, and granular telemetry.

## Scope

May create or edit:
- `app/frontend/src/views/ArenaHistoryView.vue` (new view)
- `app/frontend/src/components/arena/MatchReplayDialog.vue` (modal inspecting 3NF round-by-round replay)
- `app/frontend/src/router/index.ts` (register `/arena/history` route under `AppLayout`)
- `.routing/ROUTES.md` (record `/arena/history` route)
- `.tickets/019-build-arena-history-route/changelog.md` (create)

Must not touch:
- Backend code, Go services, or `.context/`
- Any file outside the declared scope

## Inputs

- [.routing/services/arena_service.md](../../.routing/services/arena_service.md): Contracts for `MatchSummary`, `MatchReplayDetail`, `MatchRoundReplay`, and endpoints (`GET /api/v1/arena/matches/history`, `GET /api/v1/arena/matches/{matchId}`).
- [.routing/CONTEXT.md](../../.routing/CONTEXT.md): Ground rules (zero flag emojis, full country names).
- [.routing/DESIGN.md](../../.routing/DESIGN.md): Monochrome typography and surfaces.
- [app/frontend/src/components/ui/AppDialog.vue](../../app/frontend/src/components/ui/AppDialog.vue): Reusable modal dialog component.
- Visual reference: "Contiendas Recientes" section from [Image 1: Desktop Arena inspiration](../../../.t3/userdata/attachments/e7510c5f-4aee-4fe2-8b04-76e954f1ea52-ed6ca1a6-8ef6-4587-8b46-49610bce8526.png).

## Task

1. **Register Route in Router (`router/index.ts`):**
   - Path `/arena/history`, Name `arena-history`, Component `views/ArenaHistoryView.vue`, Layout `AppLayout`.

2. **Match History View (`ArenaHistoryView.vue`):**
   - Back link to `/arena` (`"← Back to Arena"`).
   - Headline: *"Match History & Replays"* in `text-3xl font-black text-main`.
   - Table of recent matches played by the user:
     - Room Code (e.g. `#MX8921`).
     - Deck Title.
     - Format (e.g. `8 rounds · 15s`).
     - Winner & Top Score (e.g. `marianag · 1,420 pts`).
     - Date & Time.
     - Action: `"View Replay →"` button.

3. **Round-by-Round Replay Dialog (`MatchReplayDialog.vue`):**
   - Uses `AppDialog.vue` with `size="2xl"` or `maxWidth="max-w-4xl"`.
   - Displays podium breakdown (1st, 2nd, 3rd place with XP earned).
   - Stepper / list through rounds $1..N$:
     - Prompt & modality (`DEFINITION_TO_WORD`, `WORD_TO_DEFINITION`, etc.).
     - 4 options, highlighting the `correct_option`.
     - Explanation note.
     - Responses of all participants in that round (selected option, time taken in ms, correct/incorrect badge).

4. **Validation:**
   - Run `cd app && wails3 task check`.

## Verification

All checks are listed in [spec.md](spec.md).
