# 035 — Spec: Studio & Mastery Contracts, BFF Fixes, and Stateful Mock Parity

- [ ] `core-service/internal/service/studio_service.go` hydrates `RegionNames` on `CreateDeck`, `UpdateDeck`, `CreateCollection`, and `UpdateCollection` (verified by code inspection & `go build ./...`; `go test -race ./internal/service/...` awaits Vagrant `postgres-test` `:5433` via `task test:up`).
- [x] `app/app_studio.go` `ListMasteryItems` decodes and returns `[]models.StudyCardItem`, and the 15 unused legacy forwarding aliases (`ExplorePublicCollections`, `GetCollectionDetail`, `ExplorePublicDecks`, `GetDeckDetail`, `ClonePublicDeck`, `ListMyDecks`, `AddSenseToDeck`, `RemoveSenseFromDeck`, `ListMyCollections`, `GetStudySessionQueue`, `SubmitCardReview`, `SubmitStudyReview`, `GetMasterySummary`, `GetMyXPTransactions`, `GetMyAchievements`) are removed.
- [x] `app/app_studio_test.go` verifies `ListMasteryItems` decodes `[]models.StudyCardItem` and `go test -race ./...` (from `app/`) exits `0`.
- [x] `app/frontend/src/api/contract/studio.contract.ts` declares `listStudioDecks`, `getDeck`, `createDeck`, `updateDeck`, `deleteDeck`, `addDeckItem`, `removeDeckItem`, `listStudioCollections`, `createCollection`, `updateCollection`, and `deleteCollection` using DTOs from `@bindings/palabrario/pkg/models/index.js`.
- [x] `app/frontend/src/api/contract/mastery.contract.ts` declares `IMasteryAPI` (`getStudySession`, `submitReview`, `getMasteryDashboard`, `listMasteryItems`, `listXPLedger`), and `app/frontend/src/api/contract/index.ts` exposes `readonly mastery: IMasteryAPI`.
- [x] `app/frontend/src/api/live/wails/studio.wails.ts` (mapping `404` -> `null` on `getDeck`/`getCollection` while rethrowing `403`), `mastery.wails.ts`, and `wails.adapter.ts` implement every contract method without `$CancellablePromise` leaks.
- [x] `app/frontend/src/api/mock/mock.adapter.ts` and `fixtures/` enforce auth, `RN03`, `RN04`, and `RN05`, persist Studio + SM-2 state in `localStorage`, sync XP/level/streak to the mock session user, and contain zero `flag_emoji` or dialect codes rendered in UI strings.
- [x] `.routing/API.md` documents the expanded `IStudioAPI` and `IMasteryAPI`.
- [x] `wails3 task check` (from `app/`) exits with code `0`.
- [x] `git status --short` shows only files within the declared `Scope`.
- [x] No unauthorized `git commit`, `git push`, or `git merge` was run.
