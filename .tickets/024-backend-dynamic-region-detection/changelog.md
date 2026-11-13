# Changelog — Ticket 024: Dynamic Region Detection in Studio Engine

## Summary
Added PostgreSQL-level dynamic regional origin detection (`region_names`) to `Collection` and `Deck` queries via `sqlc`, propagated `RegionNames []string` (`json:"region_names"`) across Go domain models, service hydration, Wails bridge normalization, TypeScript bindings, mock fixtures, and documentation.

## Changed Files

- **`sql/queries/studio.sql`**
  - Added top-3 `COALESCE(..., '{}'::TEXT[])::TEXT[] AS region_names` subquery (`GROUP BY r.id_region, r.name ORDER BY count(*) DESC, r.name ASC LIMIT 3`) to:
    - Deck queries: `ListPublicDecks`, `ListStudioDecks`, `GetDeckByID`, `ListCollectionDecks`
    - Collection queries: `ListPublicCollections`, `ListStudioCollections`, `GetCollectionByID`
- **`core-service/internal/db/studio.sql.go`**
  - Regenerated cleanly with `sqlc generate`; added `RegionNames []string` (`json:"region_names"`) to `GetCollectionByIDRow`, `GetDeckByIDRow`, `ListCollectionDecksRow`, `ListPublicCollectionsRow`, `ListPublicDecksRow`, `ListStudioCollectionsRow`, and `ListStudioDecksRow`.
- **`pkg/models/studio.go`**
  - Added `RegionNames []string \`json:"region_names"\`` to `Deck` and `Collection`.
- **`core-service/internal/service/studio_service.go`**
  - Added `nonNilStrings` helper to guarantee nil slices become `[]string{}`.
  - Hydrated `RegionNames` in `ListPublicDecks`, `ListStudioDecks`, `GetDeckDetail`, `ListPublicCollections`, `ListStudioCollections`, and `GetCollectionDetail` (both collection and nested `Decks`).
  - Initialized `RegionNames: []string{}` in `CreateDeck`, `UpdateDeck`, `CreateCollection`, and `UpdateCollection`.
- **`app/app_studio.go`**
  - Added `normalizeDeckRegions` and `normalizeCollectionRegions` helpers and invoked them across `ListCollections`, `GetCollection`, `ListDecks`, `GetDeck`, `ListStudioDecks`, `CreateDeck`, `UpdateDeck`, `ListStudioCollections`, `CreateCollection`, and `UpdateCollection` so Wails bridge responses never return `nil` for `RegionNames`.
- **`core-service/internal/service/studio_mastery_test.go`**
  - Added `TestStudioService_DynamicRegionDetection` covering:
    - Seeded deck (`Mexicanismos Indispensables`) returning `["México"]`.
    - Multi-region deck (`>3` regions) verifying `count(*) DESC, r.name ASC LIMIT 3` (`["México", "Chile", "Cuba"]`) across `GetDeckDetail`, `ListPublicDecks`, and `ListStudioDecks`.
    - Empty deck and empty collection returning non-nil `[]string{}` across `Create*`, `Update*`, `Get*Detail`, `ListPublic*`, and `ListStudio*`.
    - Multi-deck collection aggregating region counts across constituent decks (`["Chile", "España", "México"]`) and hydrating nested `Decks[i].RegionNames` in `GetCollectionDetail`, `ListPublicCollections`, and `ListStudioCollections`.
- **`app/frontend/src/api/mock/fixtures/decks.fixture.ts`**
  - Added `DECK_REGION_NAMES` with full region names (no flag emojis, no dialect/country codes) and wired `region_names` into `deckOf` and `mockDeckDetails`.
  - Added `aggregateCollectionRegions` and `collectionOf` so `mockCollectionDetails` dynamically computes the top 1–3 ordered region names from its constituent decks.
- **`app/frontend/bindings/palabrario/pkg/models/models.ts`**
  - Regenerated Wails TypeScript bindings via `wails3 task check` (`Collection`, `Deck`, and `StudioDeckSummary` include `"region_names": string[]`).
- **`.routing/services/studio_service.md`**
  - Synced `Deck` and `Collection` Go struct snippets with `RegionNames []string \`json:"region_names"\``.

## Verification Results

- `sqlc generate` (repo root): exited `0`.
- `go test -v -race -count=1 -run TestStudioService ./core-service/internal/service/...`: exited `0` (all 6 `TestStudioService_*` tests passed, including `TestStudioService_DynamicRegionDetection` in `0.74s`).
- `go test -race -count=1 -p 1 ./core-service/...`: exited `0` (`palabrario/core-service/internal/api` passed in `48.131s`, `palabrario/core-service/internal/service` passed in `54.919s`; note: running both packages in parallel with `-race` hit a transient 10s `:5433` postgres-test connection deadline in `TestArenaCatalogService_PrepareAndGetMatchPackage`, while `-p 1` passed cleanly).
- `go test -race -count=1 ./...` (inside `app/`): exited `0` (`palabrario/app` passed in `4.338s`).
- `wails3 task check` (inside `app/`): exited `0` (`wails3 generate bindings -clean=true -ts`, `vue-tsc --noEmit`, `vite build --mode production`, and `go build .` all succeeded).
