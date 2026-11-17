# 035 — Studio & Mastery Contracts, BFF Fixes, and Stateful Mock Parity

| Field | Value |
| :--- | :--- |
| **ID** | 035 |
| **Status** | completed |
| **Model** | Gemini 4 Hydrogen, `thinking_level: medium` (`gemini-4-hydrogen`) |
| **Role** | Full-Stack Implementer (`implementation`) |
| **Spec** | [spec.md](spec.md) |

> **Specification Note:**
> Ticket 1 of 8 for the Studio suite (`035`–`042`). Fixes backend/BFF contract gaps (`app/app_studio.go` `ListMasteryItems` decode mismatch, removes 15 unused legacy BFF forwarding aliases in `app/app_studio.go`, and hydrates `RegionNames` on `StudioService` write paths), expands `app/frontend/src/api/contract/` (`IStudioAPI` + new `IMasteryAPI`), wires `WailsAdapter`, and builds stateful Spanish mock fixtures + SM-2 simulation in `MockAdapter`.

## Scope

May create or edit:
- `core-service/internal/service/studio_service.go`
- `core-service/internal/service/studio_mastery_test.go`
- `app/app_studio.go`
- `app/app_studio_test.go` (create)
- `app/frontend/bindings/` (regenerated only via `wails3 task check` / `wails3 task bindings`)
- `app/frontend/src/api/contract/index.ts`
- `app/frontend/src/api/contract/studio.contract.ts`
- `app/frontend/src/api/contract/mastery.contract.ts` (create)
- `app/frontend/src/api/live/wails/wails.adapter.ts`
- `app/frontend/src/api/live/wails/studio.wails.ts`
- `app/frontend/src/api/live/wails/mastery.wails.ts` (create)
- `app/frontend/src/api/mock/mock.adapter.ts`
- `app/frontend/src/api/mock/fixtures/decks.fixture.ts`
- `app/frontend/src/api/mock/fixtures/mastery.fixture.ts` (create)
- `.routing/API.md`
- `.tickets/035-studio-contracts-bff-and-mock-parity/changelog.md` (create)

Must not touch:
- `SCHEMA_FINAL.sql`, `sql/`, `pkg/`, `realtime-service/`, `image-service/`, `core-service/internal/api/`, `core-service/internal/db/`
- `app/frontend/uno.config.ts`
- `app/frontend/src/views/`, `app/frontend/src/components/`, `app/frontend/src/router/`
- Any file outside the above scope
- **No role (Router, Investigator, Implementer, Translator, or Reviewer) may commit, push, or merge without explicit user authorization.** Do not run `git commit`, `git push`, or `git merge`.

## Inputs

- [.routing/CONTEXT.md](../../.routing/CONTEXT.md) and [.routing/API.md](../../.routing/API.md).
- [.routing/services/studio_service.md](../../.routing/services/studio_service.md) and [.routing/services/mastery_service.md](../../.routing/services/mastery_service.md).
- [core-service/internal/service/studio_service.go](../../core-service/internal/service/studio_service.go) (`CreateDeck`, `UpdateDeck`, `CreateCollection`, `UpdateCollection`, `GetDeckDetail`, `GetCollectionDetail`).
- [core-service/internal/service/mastery_service.go](../../core-service/internal/service/mastery_service.go) (`CalculateSM2`, `ListMasteryItems`).
- [core-service/internal/service/studio_mastery_test.go](../../core-service/internal/service/studio_mastery_test.go).
- [app/app_studio.go](../../app/app_studio.go) and [app/app_auth_test.go](../../app/app_auth_test.go).
- [app/frontend/src/api/contract/studio.contract.ts](../../app/frontend/src/api/contract/studio.contract.ts), [app/frontend/src/api/live/wails/auth.wails.ts](../../app/frontend/src/api/live/wails/auth.wails.ts), [app/frontend/src/api/live/wails/studio.wails.ts](../../app/frontend/src/api/live/wails/studio.wails.ts), and [app/frontend/src/api/mock/mock.adapter.ts](../../app/frontend/src/api/mock/mock.adapter.ts).

## Task

1. **Fix `RegionNames` hydration on write paths in [studio_service.go](../../core-service/internal/service/studio_service.go):**
   - In `CreateDeck` (when `len(req.SenseIDs) > 0`), `UpdateDeck`, `CreateCollection` (when `len(req.DeckIDs) > 0`), and `UpdateCollection`, query `GetDeckByID` / `GetCollectionByID` (with `nonNilStrings(...)`) after writes so returned `Deck` / `Collection` values carry accurate `RegionNames` instead of hardcoded `[]string{}`. Keep `nonNilStrings` guaranteed for empty decks/collections.
   - Add assertions in [studio_mastery_test.go](../../core-service/internal/service/studio_mastery_test.go) verifying non-empty `RegionNames` directly on the return values of `CreateDeck`, `UpdateDeck`, `CreateCollection`, and `UpdateCollection`.

2. **Fix BFF decode bug & remove legacy aliases in [app_studio.go](../../app/app_studio.go) + add tests ([app_studio_test.go](../../app/app_studio_test.go)):**
   - Change `ListMasteryItems(stage string)` (line 428) return type and decode target from `[]models.UserMastery` to `[]models.StudyCardItem` to match `GET /api/v1/studio/mastery/items` (`MasteryService.ListMasteryItems`), initializing `[]models.StudyCardItem{}` when nil.
   - Remove the 15 unused 1-line forwarding aliases (`ExplorePublicCollections`, `GetCollectionDetail`, `ExplorePublicDecks`, `GetDeckDetail`, `ClonePublicDeck`, `ListMyDecks`, `AddSenseToDeck`, `RemoveSenseFromDeck`, `ListMyCollections`, `GetStudySessionQueue`, `SubmitCardReview`, `SubmitStudyReview`, `GetMasterySummary`, `GetMyXPTransactions`, `GetMyAchievements`), keeping only the canonical methods (plus `AddDeckToCollection` and `RemoveDeckFromCollection`).
   - Create `app/app_studio_test.go` (using `httptest.NewServer` + `newAppWithConfig` as in `app_auth_test.go`) testing `ListMasteryItems` decoding `[]models.StudyCardItem` with nested `CurrentMastery` and `Examples`.

3. **Expand Frontend Contracts (`contract/`):**
   - Update [studio.contract.ts](../../app/frontend/src/api/contract/studio.contract.ts) (`IStudioAPI`) adding:
     - `listStudioDecks(): Promise<StudioDeckSummary[]>`
     - `getDeck(id: number): Promise<Deck | null>`
     - `createDeck(req: CreateDeckRequest): Promise<Deck>`
     - `updateDeck(id: number, req: UpdateDeckRequest): Promise<Deck>`
     - `deleteDeck(id: number): Promise<void>`
     - `addDeckItem(deckId: number, req: AddDeckItemRequest): Promise<DeckItem>`
     - `removeDeckItem(deckId: number, senseId: number): Promise<void>`
     - `listStudioCollections(): Promise<Collection[]>`
     - `createCollection(req: CreateCollectionRequest): Promise<Collection>`
     - `updateCollection(id: number, req: UpdateCollectionRequest): Promise<Collection>`
     - `deleteCollection(id: number): Promise<void>`
   - Create `app/frontend/src/api/contract/mastery.contract.ts` exporting `IMasteryAPI`:
     - `getStudySession(deckId?: number, limit?: number): Promise<StudySessionPayload>`
     - `submitReview(req: SubmitReviewRequest): Promise<CardReviewResult>`
     - `getMasteryDashboard(): Promise<MasteryDashboard>`
     - `listMasteryItems(stage?: MasteryStage): Promise<StudyCardItem[]>`
     - `listXPLedger(limit?: number): Promise<XPTransactionItem[]>`
   - Update [contract/index.ts](../../app/frontend/src/api/contract/index.ts) to export `IMasteryAPI` and add `readonly mastery: IMasteryAPI` to `IPalabrarioAPI`. Import all DTOs directly from `@bindings/palabrario/pkg/models/index.js`.

4. **Implement Wails Adapters (`live/wails/`):**
   - Implement all new `IStudioAPI` methods in [studio.wails.ts](../../app/frontend/src/api/live/wails/studio.wails.ts). In `getDeck` and `getCollection`, catch `404` / `"not found"` errors (matching the helper pattern in `auth.wails.ts`) and return `null`, while rethrowing `403` (`RN04`) and server errors.
   - Create `app/frontend/src/api/live/wails/mastery.wails.ts` (`WailsMasteryAdapter`) delegating `getStudySession`, `submitReview`, `getMasteryDashboard`, `listMasteryItems`, and `listXPLedger` to `@bindings/palabrario/app/app.js`.
   - Register `readonly mastery = new WailsMasteryAdapter()` in [wails.adapter.ts](../../app/frontend/src/api/live/wails/wails.adapter.ts).

5. **Implement Stateful Mock Studio & Mastery Parity (`mock/`):**
   - Create `app/frontend/src/api/mock/fixtures/mastery.fixture.ts` with realistic Spanish `StudyCardItem` fixtures (drawn from `mockWords`, full region names, zero rendered `flag_emoji`), initial `UserMastery` progress states (`LEARNING`, `REVIEWING`, `MASTERED`), and `XPTransactionItem[]` ledger entries.
   - Update [decks.fixture.ts](../../app/frontend/src/api/mock/fixtures/decks.fixture.ts) so `camila_rojas` (`00000000-0000-4000-8000-000000000001`) owns at least 1 public deck and 1 private deck (and 1 owned collection), alongside `valentina_ospina` (`00000000-0000-4000-8000-000000000002`) owning official content.
   - Update [mock.adapter.ts](../../app/frontend/src/api/mock/mock.adapter.ts):
     - Share active session user lookup (reading `readStorage()`) across `MockStudioAPI` and `MockMasteryAPI`, rejecting unauthenticated calls to Studio/Mastery methods with `new Error("Authentication required.")`.
     - Implement all new `IStudioAPI` methods with `localStorage`-backed persistence, enforcing `RN03` (`id_user === user.id_user || user.role === ADMIN`), `RN04` (private `getDeck`/`getCollection`/`saveDeck`/`saveCollection` throw forbidden error for non-owner non-admin; return `null` when ID does not exist), and `RN05` (`is_official`/`is_featured` forced `false` for `LEARNER`), plus recomputing `region_names` (top 3 by count, alphabetical tiebreak).
     - Implement `MockMasteryAPI` (`readonly mastery: IMasteryAPI` on `MockAdapter`) with deterministic SM-2 math (`1..4` ratings -> `q in {1,3,4,5}`, intervals `1 -> 6 -> ceil(I*EF)`, `LEARNING`/`REVIEWING`/`MASTERED` stage transitions, XP `+5`/`+15`/`+20`, level recalculation, achievement checks, and syncing updated `total_xp`/`current_level`/`current_streak` back into `SESSION_KEY` and `MOCK_USERS_KEY`).

6. **Update Docs & Verify:**
   - Update [.routing/API.md](../../.routing/API.md) to document `IMasteryAPI` and the expanded `IStudioAPI` methods.
   - Run `go test -race ./internal/service/...` in `core-service/`, `go test -race ./...` in `app/`, and `wails3 task check` in `app/` and confirm exit code `0`.

## Verification

- Run all checks in [spec.md](spec.md).

## Change log

- Create `.tickets/035-studio-contracts-bff-and-mock-parity/changelog.md` recording:
  - Files modified/created.
  - Exact command outputs for `go test -race` and `wails3 task check`.
  - Any edge cases discovered in `StudioService`, `app_studio.go`, or mock fixtures.
