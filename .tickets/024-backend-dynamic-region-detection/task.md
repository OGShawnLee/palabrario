# 024 — Dynamic Region Detection in Studio Engine (Collections & Decks)

| Field | Value |
| :--- | :--- |
| **ID** | 024 |
| **Status** | completed |
| **Model** | Gemini 4 Hydrogen, `thinking_level: medium` (`gemini-4-hydrogen`) |
| **Role** | Backend / Database Engineer (`implementation`) |
| **Spec** | [spec.md](spec.md) |

> **Specification Note:**  
> This ticket adds dynamic regional origin detection to Collections and Decks directly in PostgreSQL via `sqlc`. By aggregating sense regions through leading primary key joins, the database engine returns the top 1–3 most prominent dialect regions (`region_names text[]`) without client-side N+1 hydration.

## Scope

May create or edit:
- `sql/queries/studio.sql` (update `ListPublicCollections`, `ListStudioCollections`, `GetCollectionByID`, `ListCollectionDecks`, `ListPublicDecks`, `ListStudioDecks`, `GetDeckByID`)
- `pkg/models/studio.go` (add `RegionNames []string` with JSON tag `json:"region_names"` to `Collection` and `Deck`)
- `core-service/internal/db/studio.sql.go` (regenerate via `sqlc generate`)
- `core-service/internal/db/models.go` (regenerate via `sqlc generate` if changed)
- `core-service/internal/service/studio_service.go` (hydrate `RegionNames` in `Deck` and `Collection` mappings, guaranteeing non-nil `[]string{}`)
- `core-service/internal/service/studio_mastery_test.go` (add/update tests verifying `RegionNames` aggregation, top-3 ordering, and empty-slice behavior)
- `app/app_studio.go` (ensure Wails bridge methods normalize/return non-nil `RegionNames` if needed)
- `app/frontend/src/api/mock/fixtures/decks.fixture.ts` (add `region_names` to mock collections and decks)
- `app/frontend/bindings/` (regenerated cleanly by `wails3 task check`)
- `.routing/services/studio_service.md` (add `RegionNames` field to `Deck` and `Collection` model docs)
- `.tickets/024-backend-dynamic-region-detection/changelog.md` (create)

Must not touch:
- `SCHEMA_FINAL.sql` (this is a non-destructive query-level aggregation; no DDL migration needed)
- `app/frontend/uno.config.ts`
- Any files outside the above declared scope
- **No role (Router, Investigator, Implementer, Translator, or Reviewer) may commit, push, or merge without explicit user authorization.** Do not run `git commit`, `git push`, or `git merge`.

## Inputs

- [.routing/services/studio_service.md](../../.routing/services/studio_service.md): Studio domain specification.
- [sql/queries/studio.sql](../../sql/queries/studio.sql): Existing query file (see `SearchDecksDropdown` for the existing single-region precursor).
- [pkg/models/studio.go](../../pkg/models/studio.go): Canonical Go domain models.
- [core-service/internal/service/studio_service.go](../../core-service/internal/service/studio_service.go): Studio SOA implementation and SQL row mappings.
- [core-service/internal/service/studio_mastery_test.go](../../core-service/internal/service/studio_mastery_test.go): Studio service integration tests against PostgreSQL.
- [.routing/CONTEXT.md](../../.routing/CONTEXT.md): Full regional names only, zero flag emojis, zero dialect codes.

## Architectural Context

In `sql/queries/studio.sql`, `SearchDecksDropdown` already computes `r.name` via an inline scalar subquery. Extend this pattern to compute top 1–3 distinct regions ordered by sense frequency (`count(*) DESC, r.name ASC`):

### Collection Subquery Snippet
```sql
COALESCE((
    SELECT array_agg(sub.name)
    FROM (
        SELECT r.name
        FROM collection_decks cd
        JOIN deck_items di ON di.id_deck = cd.id_deck
        JOIN senses s ON s.id_sense = di.id_sense
        JOIN regions r ON r.id_region = s.id_region
        WHERE cd.id_collection = c.id_collection
        GROUP BY r.id_region, r.name
        ORDER BY count(*) DESC, r.name ASC
        LIMIT 3
    ) sub
), '{}'::TEXT[])::TEXT[] AS region_names
```

### Deck Subquery Snippet
```sql
COALESCE((
    SELECT array_agg(sub.name)
    FROM (
        SELECT r.name
        FROM deck_items di
        JOIN senses s ON s.id_sense = di.id_sense
        JOIN regions r ON r.id_region = s.id_region
        WHERE di.id_deck = d.id_deck
        GROUP BY r.id_region, r.name
        ORDER BY count(*) DESC, r.name ASC
        LIMIT 3
    ) sub
), '{}'::TEXT[])::TEXT[] AS region_names
```

## Task

1. **SQL Queries Update (`sql/queries/studio.sql`):**
   - In `ListPublicCollections`, `ListStudioCollections`, and `GetCollectionByID`: add the Collection `region_names` subquery.
   - In `ListPublicDecks`, `ListStudioDecks`, `GetDeckByID`, and `ListCollectionDecks`: add the Deck `region_names` subquery.

2. **Generate `sqlc` Code:**
   - Run `sqlc generate` from the repo root.
   - Verify `RegionNames []string` is generated on the respective row structs in `core-service/internal/db/studio.sql.go`.

3. **Domain Models & Service Hydration:**
   - In `pkg/models/studio.go`: Add `RegionNames []string` with tag `` `json:"region_names"` `` to `Collection` and `Deck`.
   - In `core-service/internal/service/studio_service.go`: map `row.RegionNames` into `model.RegionNames` across `ListPublicDecks`, `ListStudioDecks`, `GetDeckDetail`, `ListPublicCollections`, `ListStudioCollections`, `GetCollectionDetail` (both collection and nested decks), and initialize empty `RegionNames: []string{}` in `CreateDeck`, `UpdateDeck`, `CreateCollection`, and `UpdateCollection` so slices never serialize as `null`.
   - In `core-service/internal/service/studio_mastery_test.go`: add integration assertions verifying `RegionNames` returns top 1–3 ordered full region names for populated decks/collections and `[]string{}` when empty.
   - Run `go test -race -count=1 ./core-service/...` against the running `postgres-test` container (`:5433`, start with `task test:up` only if down).

4. **Mock Fixtures, Docs & Frontend Bindings:**
   - In `app/frontend/src/api/mock/fixtures/decks.fixture.ts`: Populate realistic full region names (no flag emojis, no dialect codes) in `region_names` on `mockDeckDetails`, `deckOf`, and `mockCollectionDetails`.
   - In `.routing/services/studio_service.md`: Update the `Deck` and `Collection` Go snippets to include `RegionNames []string`.
   - Run `wails3 task check` inside `app/` to regenerate Wails TS bindings (`app/frontend/bindings/`) and verify clean `vue-tsc`, Vite, and Go builds.

5. **Changelog:**
   - Document all changed files, tests, and verification command outputs in `.tickets/024-backend-dynamic-region-detection/changelog.md`.

## Verification

All checks are listed in [spec.md](spec.md).

