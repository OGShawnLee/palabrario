# 024 — Dynamic Region Detection in Studio Engine (Collections & Decks)

| Field | Value |
| :--- | :--- |
| **ID** | 024 |
| **Status** | ready |
| **Model** | Claude Sonnet 5.5, `effort: medium` (`claude-sonnet-5-5`) |
| **Role** | Backend / Database Engineer |
| **Spec** | [spec.md](spec.md) |

> **Specification Note:**  
> This ticket adds dynamic regional origin detection to Collections and Decks directly in PostgreSQL via `sqlc`. By aggregating sense regions through leading primary key joins, the database engine returns the top 1–3 most prominent dialect regions (`region_names text[]`) without client-side N+1 hydration.

## Scope

May create or edit:
- `sql/queries/studio.sql` (update `ListPublicCollections`, `ListStudioCollections`, `GetCollectionByID`, `ListPublicDecks`, `ListStudioDecks`, `GetDeckByID`)
- `pkg/models/studio.go` (add `RegionNames []string `json:"region_names"`` to `Collection` and `Deck`)
- `core-service/internal/repository/studio.sql.go` (regenerate via `sqlc generate`)
- `core-service/internal/repository/models.go` (regenerate via `sqlc generate`)
- `core-service/internal/service/studio/` (hydrate `RegionNames` in service mapping)
- `app/app_studio.go` (ensure Wails bridge returns `RegionNames`)
- `app/frontend/src/api/mock/fixtures/decks.fixture.ts` (add `region_names` to mock collections and decks)
- `.tickets/024-backend-dynamic-region-detection/changelog.md` (create)

Must not touch:
- `SCHEMA_FINAL.sql` (this is a non-destructive query-level aggregation; no DDL migration needed)
- `app/frontend/uno.config.ts`
- Any files outside the above declared scope

## Inputs

- [.routing/services/studio_service.md](../../.routing/services/studio_service.md): Studio domain specification.
- [sql/queries/studio.sql](../../sql/queries/studio.sql): Existing query file (see lines 310–319 for the existing single-region precursor in `SearchDecksDropdown`).
- [pkg/models/studio.go](../../pkg/models/studio.go): Canonical Go domain models.
- [.routing/CONTEXT.md](../../.routing/CONTEXT.md): Full regional names only, zero flag emojis, zero dialect codes.

## Architectural Context

In `sql/queries/studio.sql`, `SearchDecksDropdown` already computes `r.name` via an inline scalar subquery. We extend this pattern to compute top 1–3 distinct regions ordered by sense frequency:

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
   - In `ListPublicDecks`, `ListStudioDecks`, and `GetDeckByID`: add the Deck `region_names` subquery.

2. **Generate Repository Code:**
   - Run `sqlc generate` (or run through `task generate` / `wails3 task check`).
   - Verify `region_names []string` is generated on the respective row structs in `core-service/internal/repository/`.

3. **Domain Models & Service Hydration:**
   - In `pkg/models/studio.go`: Add `RegionNames []string `json:"region_names"`` to `Collection` and `Deck`.
   - In `core-service/internal/service/studio/` mapping functions: map `row.RegionNames` to `model.RegionNames` (ensuring non-nil empty slice `[]string{}` if empty).
   - In `app/app_studio.go`: Verify Wails bridge returns the updated models.

4. **Mock Fixtures & Frontend Bindings:**
   - In `app/frontend/src/api/mock/fixtures/decks.fixture.ts`: Populate realistic `region_names` on `mockCollectionDetails` and `mockDeckDetails` (e.g. `["Argentina", "Chile"]`).
   - Run `cd app && wails3 task check` to regenerate Wails TS bindings and verify clean build.

5. **Changelog:**
   - Document changes in `.tickets/024-backend-dynamic-region-detection/changelog.md`.

## Verification

All checks are listed in [spec.md](spec.md).
