# 019 — Build the Arena Match History & Replay Route (/arena/history)

| Field | Value |
| :--- | :--- |
| **ID** | 019 |
| **Status** | draft |
| **Model** | Claude Sonnet 5.5, `effort: medium` (`claude-sonnet-5-5`) |
| **Role** | Frontend Implementer |
| **Spec** | [spec.md](spec.md) |

> **Draft Note:**  
> This ticket builds the dedicated match history route (`/arena/history`) and the 3NF round-by-round match replay inspection dialog (`MatchReplayDialog.vue`), satisfying Use Case `CU16` and enabling players to review past matches, scores, podium finishes, and granular telemetry.

## Scope

May create or edit:
- `app/frontend/src/views/ArenaHistoryView.vue` (new view)
- `app/frontend/src/components/arena/MatchReplayDialog.vue` (modal inspecting 3NF round-by-round replay)
- `app/frontend/src/api/contract/arena.contract.ts` (support listMatchHistory, getMatchReplay)
- `app/frontend/src/api/live/wails/arena.wails.ts` (implement listMatchHistory, getMatchReplay)
- `app/frontend/src/api/mock/mock.adapter.ts` (implement mock listMatchHistory, getMatchReplay)
- `app/frontend/src/router/index.ts` (register `/arena/history` route under `AppLayout`)
- `.routing/ROUTES.md` (record `/arena/history` route)
- `.tickets/019-build-arena-history-route/changelog.md` (create)

Must not touch:
- Backend code, Go services, or `.context/`
- Any file outside the declared scope

## Inputs

- [.routing/API.md](../../.routing/API.md): Architectural rules (Rule 1: Never import `@bindings/.../app` in UI; always import `api` from `@/api`).
- [app/frontend/src/api/index.ts](../../app/frontend/src/api/index.ts): Universal `api` singleton (`api.arena`).
- [.routing/services/arena_service.md](../../.routing/services/arena_service.md): Contracts for `MatchSummary`, `MatchReplayDetail`, `MatchRoundReplay`, and endpoints (`GET /api/v1/arena/matches/history`, `GET /api/v1/arena/matches/{matchId}`).
- [.routing/CONTEXT.md](../../.routing/CONTEXT.md): Ground rules (zero flag emojis, full country names).
- [.routing/DESIGN.md](../../.routing/DESIGN.md): Monochrome typography and surfaces.
- [app/frontend/src/components/ui/AppDialog.vue](../../app/frontend/src/components/ui/AppDialog.vue): Reusable modal dialog component.
- Visual reference: "Contiendas Recientes" section from [Image 1: Desktop Arena inspiration](../../../.t3/userdata/attachments/e7510c5f-4aee-4fe2-8b04-76e954f1ea52-ed6ca1a6-8ef6-4587-8b46-49610bce8526.png).

## Task

1. **API Client Support for Match History & Replays:**
   - In `app/frontend/src/api/contract/arena.contract.ts`:
     - Add `listMatchHistory(limit?: number): Promise<MatchSummary[]>`.
     - Add `getMatchReplay(matchId: string): Promise<MatchReplayDetail | null>`.
   - In `app/frontend/src/api/live/wails/arena.wails.ts`:
     - Delegate to `Wails.ListMatchHistory(limit ?? 20)` and `Wails.GetMatchReplay(matchId)`.
   - In `app/frontend/src/api/mock/mock.adapter.ts`:
     - Implement mock `listMatchHistory` and `getMatchReplay` returning authentic Spanish round data, prompts, and score breakdowns.

2. **Register Route in Router (`router/index.ts`):**
   - Path `/arena/history`, Name `arena-history`, Component `views/ArenaHistoryView.vue`, Layout `AppLayout`.

3. **Match History View (`ArenaHistoryView.vue`):**
   - Import `api` from `@/api`:
     ```ts
     import { api } from "@/api"
     ```
   - Fetch matches on mount: `matches.value = await api.arena.listMatchHistory(20)`.
   - **Strict Rule:** Never hardcode mock match history in `ArenaHistoryView.vue`. All history items must flow through `api.arena`.
   - Back link to `/arena` (`"← Back to Arena"`).
   - Headline: *"Match History & Replays"* in `text-3xl font-black text-main`.
   - Table of recent matches played by the user:
     - Room Code (e.g. `#MX8921`).
     - Deck Title.
     - Format (e.g. `8 rounds · 15s`).
     - Winner & Top Score (e.g. `marianag · 1,420 pts`).
     - Date & Time.
     - Action: `"View Replay →"` button (fetches replay detail and opens dialog).

4. **Round-by-Round Replay Dialog (`MatchReplayDialog.vue`):**
   - Receives `matchId: string` or `replay: MatchReplayDetail`.
   - If `matchId` passed, calls `await api.arena.getMatchReplay(props.matchId)`. Zero hardcoded static round replays in the component.
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
