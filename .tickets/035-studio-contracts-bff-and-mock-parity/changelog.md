# Changelog — 035: Studio & Mastery Contracts, BFF Fixes, and Stateful Mock Parity

## Summary

Implemented Ticket 1 of the Studio suite (`035`–`042`): hydrated `RegionNames` on `StudioService` write paths, fixed the `app/app_studio.go` `ListMasteryItems` decode mismatch and removed all 15 unused legacy BFF forwarding aliases, expanded `IStudioAPI` and added `IMasteryAPI`, wired `WailsStudioAdapter` + `WailsMasteryAdapter` against regenerated bindings, built stateful `localStorage`-backed Studio + SM-2 mock parity in `MockAdapter`, and updated `.routing/API.md`.

---

## Files Modified & Created

### Backend (`core-service/`)
- **Modified** `core-service/internal/service/studio_service.go`
  - Hydrated `RegionNames` via `GetDeckByID` / `GetCollectionByID` + `nonNilStrings(...)` on write paths: `CreateDeck` (when `len(req.SenseIDs) > 0`, also populating `ItemCount`), `UpdateDeck` (also preserving `ItemCount`), `CreateCollection` (when `len(req.DeckIDs) > 0`), and `UpdateCollection`. Empty decks and collections continue returning non-nil `[]string{}`.
- **Modified** `core-service/internal/service/studio_mastery_test.go`
  - Added direct `RegionNames` assertions on the return values of `CreateDeck`, `UpdateDeck`, `CreateCollection`, and `UpdateCollection` inside `TestStudioService_DynamicRegionDetection`.

### Desktop BFF (`app/`)
- **Modified** `app/app_studio.go`
  - Fixed `ListMasteryItems(stage string)` decode target and return type from `[]models.UserMastery` to `[]models.StudyCardItem` (initializing `[]models.StudyCardItem{}` when nil).
  - Removed all 15 unused 1-line legacy forwarding aliases (`ExplorePublicCollections`, `GetCollectionDetail`, `ExplorePublicDecks`, `GetDeckDetail`, `ClonePublicDeck`, `ListMyDecks`, `AddSenseToDeck`, `RemoveSenseFromDeck`, `ListMyCollections`, `GetStudySessionQueue`, `SubmitCardReview`, `SubmitStudyReview`, `GetMasterySummary`, `GetMyXPTransactions`, `GetMyAchievements`), keeping canonical Studio/Mastery methods plus `AddDeckToCollection` / `RemoveDeckFromCollection`.
- **Created** `app/app_studio_test.go`
  - Added `httptest.NewServer` + `newAppWithConfig` unit tests (`TestListMasteryItemsDecodesStudyCardItemsWithMasteryAndExamples` and `TestGetDeckPreserves404And403APIErrors`) verifying `[]models.StudyCardItem` decoding with nested `CurrentMastery` and `Examples`, nil normalization, and `404`/`403` `*APIError` propagation.
- **Regenerated** `app/frontend/bindings/palabrario/app/app.ts` (via `wails3 task bindings` / `wails3 task check`).

### Frontend Contracts, Live Wails Adapters & Mock Parity (`app/frontend/src/api/`)
- **Modified** `app/frontend/src/api/contract/studio.contract.ts`
  - Expanded `IStudioAPI` with `listStudioDecks`, `getDeck`, `createDeck`, `updateDeck`, `deleteDeck`, `addDeckItem`, `removeDeckItem`, `listStudioCollections`, `createCollection`, `updateCollection`, and `deleteCollection` using DTOs from `@bindings/palabrario/pkg/models/index.js`.
- **Created** `app/frontend/src/api/contract/mastery.contract.ts`
  - Defined `IMasteryAPI` (`getStudySession`, `submitReview`, `getMasteryDashboard`, `listMasteryItems`, `listXPLedger`).
- **Modified** `app/frontend/src/api/contract/index.ts`
  - Re-exported `IMasteryAPI` and added `readonly mastery: IMasteryAPI` to `IPalabrarioAPI`.
- **Modified** `app/frontend/src/api/live/wails/studio.wails.ts`
  - Implemented all expanded `IStudioAPI` methods without `$CancellablePromise` leaks; `getDeck` and `getCollection` catch `404` / not-found errors (`"not found"`, `"no encontrada"`, `"no encontrado"`, `"not_found"`, `"404"`) and return `null` while rethrowing `403` (`RN04`) and server errors.
- **Created** `app/frontend/src/api/live/wails/mastery.wails.ts`
  - Implemented `WailsMasteryAdapter` delegating to `@bindings/palabrario/app/app.js`.
- **Modified** `app/frontend/src/api/live/wails/wails.adapter.ts`
  - Registered `readonly mastery = new WailsMasteryAdapter()`.
- **Modified** `app/frontend/src/api/mock/fixtures/decks.fixture.ts`
  - Added `camila_rojas` (`00000000-0000-4000-8000-000000000001`) public deck (#7), private deck (#8), and owned collection (#4); assigned official fixtures (#1, #6, collection #1) to `valentina_ospina` (`00000000-0000-4000-8000-000000000002`); normalized Deck 6 region `"Mexico"` to `"México"`; moved `editora_andina` owner UUID to `...0004`; exported `mockInitialSavedDeckIdsByUser` and `mockInitialSavedCollectionIdsByUser`.
- **Created** `app/frontend/src/api/mock/fixtures/mastery.fixture.ts`
  - Added `buildStudyCardItem`, `mockStudyCardCatalog` (sourced from `mockWords`, `flag_emoji: ""`, full region names), `mockMasteryByUser` (`LEARNING`, `REVIEWING`, `MASTERED`), and `mockXPLedgerByUser`.
- **Modified** `app/frontend/src/api/mock/mock.adapter.ts`
  - Refactored `MockAuthAPI` to read live from `localStorage` (`SESSION_KEY` / `MOCK_USERS_KEY`) so XP/level/streak updates performed by `MockMasteryAPI.submitReview()` are immediately reflected in `getMe()` and `getPublicProfileByUsername()`.
  - Implemented `localStorage`-backed `MockStudioAPI` enforcing authentication (`"Authentication required."`), `RN03`, `RN04`, and `RN05`, plus top-3 `region_names` recomputation.
  - Implemented `MockMasteryAPI` (`readonly mastery = new MockMasteryAPI()`) with deterministic SM-2 scheduling, XP/level/streak/achievement persistence, and `404`/`403` handling on deck study sessions.

### Documentation
- **Modified** `.routing/API.md`
  - Updated sections 1–2 for `api.mastery`, `mastery.contract.ts`, `mastery.wails.ts`, and `mastery.fixture.ts`, and added sections 7–8 documenting `IStudioAPI` and `IMasteryAPI` contracts and guarantees.

---

## Verification & Command Outputs

### 1. `app/` Unit & Race Tests (`go test -race -v ./...`)
Command (from `app/`): `go test -race -v ./...`
Result: **PASS** (exit code `0`)
```text
=== RUN   TestSignInPersistsRefreshTokenAndSignOutClears
--- PASS: TestSignInPersistsRefreshTokenAndSignOutClears (0.04s)
=== RUN   TestRestartRecoveryGetMeRefreshesPersistedSession
--- PASS: TestRestartRecoveryGetMeRefreshesPersistedSession (0.04s)
=== RUN   TestGetMeUnauthenticatedAndCorruptFileAndTransientError
--- PASS: TestGetMeUnauthenticatedAndCorruptFileAndTransientError (0.04s)
=== RUN   TestSingleFlightAutoRefreshAcrossConcurrentRequestsAndUploadImage
--- PASS: TestSingleFlightAutoRefreshAcrossConcurrentRequestsAndUploadImage (0.03s)
=== RUN   TestSuspendedAccountAndFailedRefreshEvictSession
--- PASS: TestSuspendedAccountAndFailedRefreshEvictSession (0.04s)
=== RUN   TestUpdateProfileAndRedeemTicketRotateAccessToken
--- PASS: TestUpdateProfileAndRedeemTicketRotateAccessToken (0.03s)
=== RUN   TestConnectArenaRoomUsesAuthorizationBearerHeader
--- PASS: TestConnectArenaRoomUsesAuthorizationBearerHeader (0.02s)
=== RUN   TestFindMatchRoomUsesAPIErrorNotFoundAndNetworkErrorSanitized
--- PASS: TestFindMatchRoomUsesAPIErrorNotFoundAndNetworkErrorSanitized (0.01s)
=== RUN   TestListMasteryItemsDecodesStudyCardItemsWithMasteryAndExamples
--- PASS: TestListMasteryItemsDecodesStudyCardItemsWithMasteryAndExamples (0.01s)
=== RUN   TestGetDeckPreserves404And403APIErrors
--- PASS: TestGetDeckPreserves404And403APIErrors (0.01s)
PASS
ok  	palabrario/app	4.213s
```

### 2. Wails Pipeline (`wails3 task check`)
Command (from `app/`): `wails3 task check`
Result: **PASS** (exit code `0`)
```text
task: [bindings] wails3 generate bindings -clean=true -ts
Wails v3.0.0 › Generate Bindings
 INFO  Processed: 264 Packages, 1 Service, 98 Methods, 6 Enums, 75 Models, 0 Events in 5.3987194s.
 INFO  Output directory: C:\Users\OGSha\OneDrive\Escritorio\Developer\Web Development\Palabrario\app\frontend\bindings
task: [check] bun run --cwd frontend type-check
$ vue-tsc --noEmit
task: [check] bun run --cwd frontend build
$ vite build --mode production
✓ 2502 modules transformed.
✓ built in 6.16s
task: [check] go build .
```

### 3. `core-service/` Compile & Service Tests (`go test -race ...`)
- `go build ./...` (from `core-service/`): **PASS** (exit code `0`).
- `go test -race -v -run "TestMasteryService_CalculateSM2_MathematicalPrecision|TestErrors_MapDBError" ./internal/service/...` (from `core-service/`): **PASS** (exit code `0`, `ok palabrario/core-service/internal/service 3.765s`).
- `go test -race -v ./internal/service/...` (from `core-service/`): exited `1` because the ephemeral PostgreSQL test container (`localhost:5433`) is offline in this environment (`docker` CLI is absent on Windows and `vagrant status` shows `default poweroff (virtualbox)`), causing `testutil.SetupTestDB` to fail with `connectex: No connection could be made because the target machine actively refused it`. Run `task test:up` once the Vagrant VM is running to execute the DB-backed assertions in `TestStudioService_DynamicRegionDetection`.

---

## Edge Cases & Findings Discovered
1. **Spanish 404 Messages in `StudioService`:** `StudioService.GetDeckDetail` and `GetCollectionDetail` wrap `pgx.ErrNoRows` with `MapDBError(err, "Baraja no encontrada")` / `"Colección no encontrada"`. `isNotFoundError` in `app/frontend/src/api/live/wails/studio.wails.ts` checks both English (`not found`, `not_found`, `404`) and Spanish (`no encontrada`, `no encontrado`) substrings so `404` maps cleanly to `null` while `403` (`RN04`) and `5xx` errors continue to throw.
2. **`MockAuthAPI` Instance Caching:** `MockAuthAPI` previously cached `this.session` and `this.accounts` in private fields, which would hide XP/level/streak updates written to `SESSION_KEY` and `MOCK_USERS_KEY` by `MockMasteryAPI.submitReview()`. Switching `MockAuthAPI` reads to `readStorage()` and `readAccountsStorage()` ensures cross-domain consistency.
3. **Mock Owner UUID Collision & Region Spelling:** In `decks.fixture.ts`, `editora_andina` previously shared UUID `...0003` with suspended account `mateo_vargas` (moved to `...0004`), and Deck 6 used `"Mexico"` instead of `"México"`.
4. **Unused Legacy Aliases Verified Safe to Remove:** Full-text scan of `app/` confirmed zero references to any of the 15 removed forwarding aliases in `app/app_studio.go`.
