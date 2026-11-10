# Changelog — Ticket 020: Audit Service Specs with Backend

Audited and synchronized all 8 modular service specification documents (`.routing/services/*_service.md`) against authoritative Go backend source code (`pkg/models/`, `core-service/internal/service/`, `core-service/internal/api/router.go`, and `sql/queries/*.sql`). Documentation only. Zero application code modified.

## Summary of Findings & Corrections

Across all 8 files, the documentation has been updated to reflect exact Go types, field tags, SOA method signatures, `sqlc` query names, and HTTP middleware behavior.

### 1. Identity Service (`identity_service.md`)
- **DTOs:** Added and synchronized `PublicUserProfile` (embedding `[]models.Deck` and `[]models.Collection`), `AccountSummary`, and `UserProfileDetail` matching `pkg/models/user.go`.
- **SOA Interface:** Documented `IdentityService` interface matching `core-service/internal/service/identity_service.go` verbatim.
- **Middleware & Endpoints:** Explicitly documented `RequireAuth` for protected user routes (`/users/me`, `/users/avatar`) and `OptionalAuth` for public routes with authenticated enhancements.
- **Persistence Queries:** Listed `users.sql` queries including user search and public profile aggregate queries (`CountUserPrivateDecks`, `ListUserPublicDecks`, `ListUserPublicCollections`, `GetUserArenaStats`).

### 2. Governance Service (`governance_service.md`)
- **DTOs:** Synchronized `InvitationTicketDetail`, `EditorCandidate`, and `VoteCandidateRequest` matching `pkg/models/governance.go`.
- **SOA Interface:** Documented `GovernanceService` interface matching `core-service/internal/service/governance_service.go` verbatim.
- **Middleware & Endpoints:** Clarified `RequireAuth + RequireRole(ADMIN, TRUSTED_EDITOR)` on ticket generation, candidate review, and voting endpoints.
- **Persistence Queries:** Listed `governance.sql` queries verbatim.

### 3. Dictionary Service (`dictionary_service.md`)
- **DTOs:** Added `image_url` field to `CreateWord` and `UpdateWordMetadata`; verified all `snake_case` tags against `pkg/models/dictionary.go`.
- **Endpoints & Middleware:**
  - Clarified `OptionalAuth` dialect fallback logic on `GET /api/v1/words/word-of-the-day` and `GET /api/v1/words/search` (prefers query param `?dialect=`, then `claims.AnchorDialect`, then MX).
  - Clarified pagination parameter validation (`1..50`, default 20) on `GET /api/v1/words`.
  - Updated dictionary authoring endpoints with `RequireAuth + RequireRole(EDITOR, ADMIN)`.
- **Persistence Queries:** Documented full list of `dictionary.sql` queries, including `GetExampleByID` and `GetExamplesBySenseID`.

### 4. Curation Service (`curation_service.md`)
- **DTOs:**
  - Added optional `ImageURL` to `SenseQuorumItem` and documented `PendingSenseItem = SenseQuorumItem` canonical alias.
  - Corrected moderation report reason enum value (`INACCURATE` vs draft `INACCURATE_DEFINITION`).
  - Added explicit hygiene note: `flag_emoji` exists on wire contract for backwards compatibility but MUST NEVER be rendered in UI.
- **SOA Interface:** Documented `CurationService` interface matching `core-service/internal/service/curation_service.go` verbatim.
- **Endpoints & Middleware:** Documented `RequireAuth + RequireRole` middleware on all curation routes, plus `409 Conflict` behavior when resolving a non-pending report.
- **Persistence Queries:** Listed `curation.sql` queries.

### 5. Studio Service (`studio_service.md`)
- **DTOs:**
  - Added canonical `Collection` struct definition (`id_collection`, `id_user`, `title`, `description`, `is_public`, `is_featured`, `is_official`, `decks`, timestamps).
  - Clarified that `StudioDeckSummary` embeds `Deck` (fields flattened in JSON).
- **SOA Interface:** Documented `StudioService` interface verbatim (`ListPublicDecks`, `ListStudioDecks`, `GetDeckDetail`, `CreateDeck`, `UpdateDeck`, `DeleteDeck`, `AddDeckItem`, `RemoveDeckItem`, `SaveDeck`, `UnsaveDeck`, `ListPublicCollections`, `ListStudioCollections`, `GetCollectionDetail`, etc.).
- **Persistence Queries:** Documented `studio.sql` queries and cross-service queries (`GetExamplesForSenses` in `mastery.sql`).

### 6. Mastery Service (`mastery_service.md`)
- **DTOs:** Synchronized `UserMastery`, `ReviewAttempt`, `XPTransaction`, `Achievement`, and `LeaderboardEntry` matching `pkg/models/mastery.go`.
- **SOA Interface:** Documented `MasteryService` interface matching `core-service/internal/service/mastery_service.go` verbatim.
- **Persistence Queries:** Listed `mastery.sql` queries.

### 7. Arena Service (`arena_service.md`)
- **DTOs:** Synchronized `PublicMatchItem`, `ArenaRoomSummary`, `MatchParticipant`, and `RoundDetail` matching `pkg/models/arena.go`.
- **SOA Interface:** Documented `ArenaCatalogService` interface verbatim.
- **S2S Security:** Documented `realtime-service` to `core-service` HMAC-SHA256 JWT handshake for batch match consolidation.
- **Persistence Queries:** Listed `arena.sql` queries.

### 8. Media Service (`media_service.md`)
- **Endpoints:** Documented `image-service` endpoints (`POST /api/v1/media/upload`, `GET /media/{key}`).
- **Constraints (`RN12`):** Documented client/server validation: max 5MB, format conversion to WebP, dimension normalization (max 1920px).

## Scope Compliance

- Audited and updated all 8 `.routing/services/*.md` files.
- No backend code, models, SQL migrations, or frontend Vue code modified.
