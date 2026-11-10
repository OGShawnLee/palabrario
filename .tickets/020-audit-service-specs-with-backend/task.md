# 020 — Audit and Synchronize Domain Service Markdown Specs with Backend Go Code

| Field | Value |
| :--- | :--- |
| **ID** | 020 |
| **Status** | completed |
| **Model** | Sonnet 5.5 (`claude-sonnet-5-5`, Antigravity provider) |
| **Role** | Documentation / Architect |
| **Spec** | [spec.md](spec.md) |

> **Draft Note:**  
> This documentation ticket performs a systematic audit across all 8 domain service specifications in `.routing/services/` against the authoritative Go source code (`pkg/models/`, `core-service/internal/service/`, `core-service/internal/api/handlers/`, and `sql/queries/`). It detects and fixes missing structs (like the `Deck` struct in `studio_service.md`), struct embedding gaps, mismatched field names/types, missing SOA interface methods, and unlisted HTTP/WS endpoints.

## Scope

May edit:
- `.routing/services/identity_service.md`
- `.routing/services/governance_service.md`
- `.routing/services/dictionary_service.md`
- `.routing/services/curation_service.md`
- `.routing/services/studio_service.md`
- `.routing/services/mastery_service.md`
- `.routing/services/arena_service.md`
- `.routing/services/media_service.md`
- `.tickets/020-audit-service-specs-with-backend/changelog.md` (create)

Must not touch:
- `app/` or application code
- `core-service/`, `realtime-service/`, `image-service/`, `pkg/`, `sql/` (these are the authoritative source code)
- `.context/`
- Any file outside `.routing/services/` and this ticket folder

## Inputs

- The 8 service specifications in [.routing/services/](../../.routing/services/):
  - `identity_service.md`
  - `governance_service.md`
  - `dictionary_service.md`
  - `curation_service.md`
  - `studio_service.md`
  - `mastery_service.md`
  - `arena_service.md`
  - `media_service.md`
- Authoritative Go models in [pkg/models/](../../pkg/models/):
  - `user.go`, `auth.go`, `governance.go`, `dictionary.go`, `curation.go`, `studio.go`, `arena.go`, `media.go`
- Authoritative SOA interfaces and implementations in [core-service/internal/service/](../../core-service/internal/service/):
  - `identity_service.go`, `governance_service.go`, `dictionary_service.go`, `curation_service.go`, `studio_service.go`, `mastery_service.go`, `arena_catalog_service.go`
- Authoritative HTTP handlers and routes in [core-service/internal/api/](../../core-service/internal/api/):
  - `routes.go` and `handlers/*.go`
- Authoritative satellite services:
  - `realtime-service/` and `image-service/`

## Task

1. **Systematic Struct & DTO Verification:**
   - For each service spec in `.routing/services/`:
     - Compare every Go struct documented in Part B against its canonical definition in `pkg/models/*.go`.
     - Verify that struct embedding (e.g. `Deck` inside `StudioDeckSummary`) has the embedded base struct documented in the file.
     - Verify that all JSON field tags match verbatim (`snake_case`).
     - Add any missing DTOs that are actively used in the service's API contract.

2. **SOA Interface Verification:**
   - Compare each SOA interface signature (e.g. `DictionaryService`, `StudioService`, `ArenaCatalogService`) documented in the markdown against the actual Go interface definition in `core-service/internal/service/*.go`.
   - Ensure method names, argument types, and return values match Go 1.24 code verbatim.

3. **Endpoints & Middleware Audit:**
   - Cross-check each endpoint table against `core-service/internal/api/routes.go` and handlers.
   - Verify HTTP methods, URL paths, and middleware requirements (`RequireAuth`, `RequireRole`, `RequireInternalService`, `OptionalAuth`).

4. **Detailed Audit Report & Changelog:**
   - In `.tickets/020-audit-service-specs-with-backend/changelog.md`:
     - Detail a comprehensive per-service audit table listing every verified struct, interface, and endpoint.
     - Document all discrepancies found and corrected in the markdown specs.

## Verification

All checks are listed in [spec.md](spec.md).
