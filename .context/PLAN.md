# 🗺️ Plan Maestro de Implementación y Estado Arquitectónico (`.context/PLAN.md`)

> **Documento de Planificación Estratégica y Hoja de Ruta a Largo Plazo (Backend-First)**  
> **Proyecto:** Palabrario — Plataforma Gamificada de Aprendizaje y Preservación del Léxico Dialectal Panhispánico  
> **Stack Backend Unificado (100% Go):** `Go 1.24 + Chi v5 + sqlc + pgxpool/v5` (`core-service :8000`, `realtime-service :8001`, `image-service :8002`) sobre `PostgreSQL 17 Alpine (:5432)`  
> **Referencias Canónicas:** [`DOMAIN.md`](./DOMAIN.md) · [`INTERFACES.md`](./INTERFACES.md) · [`SCHEMA_FINAL.sql`](../SCHEMA_FINAL.sql) · [`.context/academic/1-PROJECT-SPECIFICATION.md`](./academic/1-PROJECT-SPECIFICATION.md)

---

## 0. Principio Rector de Ingeniería: Backend-First & Contract-First

> [!IMPORTANT]
> **¿Por qué no se toca la UI (`local-app`) hasta cerrar las Fases 3.3 a 4?**
> Cuando los modelos de IA trabajan primero sobre interfaces gráficas con mockups sin un backend estable, tienden a alucinar propiedades en el frontend que la base de datos o los endpoints jamás retornan, provocando divergencia estructural.  
> Por ello, **el desarrollo sigue estrictamente el orden:**
> $$\text{1. Esquema SQL (`SCHEMA_FINAL.sql`)} \longrightarrow \text{2. Consultas Tipadas (`sqlc`)} \longrightarrow \text{3. Servicios SOA Go (`INTERFACES.md`)} \longrightarrow \text{4. Pruebas (`go test`)} \longrightarrow \text{5. Conexión UI (Fase 5)}$$

---

## 1. Diagnóstico Ejecutivo: ¿Dónde Estamos Hoy?

Al auditar el código real del repositorio frente a los **16 Casos de Uso (`CU01–CU16`)** y las **11 Reglas de Negocio (`RN01–RN11`)**, el proyecto se encuentra en la **Fase 3.3 (Mitad del Camino del Backend Nuclear)**:

```mermaid
flowchart LR
    F1["✅ Fase 1\nEsquema SQL + Triggers\n(100%)"] --> F2["✅ Fase 2\nMonorepo Go + Contratos\n(100%)"]
    F2 --> F31["✅ Fase 3.1/3.2\nCore Base + Auth JWT\n(85%)"]
    F31 --> F33["🚧 Fase 3.3 (ACTUAL)\nDiccionario & Curaduría\n(30% -> En Progreso)"]
    F33 --> F34["⏳ Fase 3.4\nStudio & Motor SM-2\n(0% en Backend)"]
    F34 --> F35["⏳ Fase 3.5\nGobernanza & Reportes\n(0% en Backend)"]
    F35 --> F4["⏳ Fase 4\nMigración Realtime a Go\n+ Arena WS & Flush\n(10% Esqueleto)"]
    F4 --> F5["⏳ Fase 5\nIntegración UI + CI/CD\n(Conectar Mocks a API)"]
    F5 --> F6["🌱 Fase 6 (Largo Plazo)\nEventos de Temporada"]
```

### 1.1. Matriz de Estado Real por Capa y Servicio

| Componente / Módulo | Estado Actual Real | Lo que YA existe y compila | Lo que FALTA por construir |
| :--- | :---: | :--- | :--- |
| **1. Base de Datos ([`SCHEMA_FINAL.sql`](../SCHEMA_FINAL.sql))** | **100% Completo** | 23 tablas, extensiones `unaccent` + `pg_trgm`, vectores `tsvector` GIN, triggers `RN06` (auto-aprobación al 3er voto) y `RN07` (anti-autoaprobación), tablas de gobernanza (`editor_requests`, `invitation_tickets`, `ticket_redemptions`, `content_reports`) y guardados (`saved_decks`, `saved_collections`). | Únicamente escribir las consultas SQL en `sql/queries/*.sql` para los módulos pendientes y ejecutar `sqlc generate`. |
| **2. Modelos Compartidos ([`pkg/models/`](../pkg/models))** | **65% Completo** | Estructuras base en `auth.go`, `user.go`, `dictionary.go`, `studio.go` y `arena.go`. | Agregar los DTOs de entrada/salida especificados en la sección 2 de [`INTERFACES.md`](./INTERFACES.md) (`governance.go`, `curation.go`, y ampliaciones en `user.go`, `dictionary.go`, `studio.go`, `arena.go`). |
| **3. `core-service` — Auth & Cuentas (`CU01`)** | **70% Completo** | `POST /api/v1/auth/sign-up`, `POST /api/v1/auth/sign-in`, `GET /api/v1/auth/me`, hashing `bcrypt`, JWT HMAC-SHA256, middlewares `RequireAuth` y `RequireRole`, pruebas en `auth_test.go`. | `PUT /api/v1/users/me` (actualizar perfil y dialecto ancla `CU01`), `PUT /api/v1/users/me/password` y `GET /api/v1/users/u/{username}` (perfil público/dual). |
| **4. `core-service` — Diccionario (`CU05, CU06`)** | **30% Completo** | `GET /api/v1/words/search` (autocompletado FTS + trigramas) y `GET /api/v1/words/{id}` (detalle con sentidos `APPROVED`, ejemplos y heterónimos). | `GET /api/v1/regions`, `GET /api/v1/words/word-of-the-day`, explorador paginado `GET /api/v1/words`, y todo el bloque de autoría de `CU06` (`POST/PUT /words`, `POST/PUT/DELETE /senses`, `examples` y `concepts`). |
| **5. `core-service` — Curaduría y Moderación (`CU07–CU09`)** | **15% Completo** | Los triggers PL/pgSQL (`trg_prevent_self_approval` y `trg_auto_approve_sense`) ya operan en PostgreSQL. | Consultas en `sql/queries/curation.sql`, `CurationService` y endpoints para cola de quórum (`/curation/queue`), votación (`/vote`), veto admin (`RN08`) y reportes (`CU09`). |
| **6. `core-service` — Studio y SM-2 (`CU10–CU13`)** | **0% en Backend** | Modelos iniciales en `pkg/models/studio.go`. | Consultas en `sql/queries/studio.sql` y `mastery.sql`, `StudioService` (CRUD de barajas/colecciones y guardado comunitario) y `MasteryService` (algoritmo matemático SM-2 y ledger `xp_transactions`). |
| **7. `core-service` — Gobernanza (`CU02–CU04`)** | **0% en Backend** | Tablas `editor_requests`, `invitation_tickets` y `ticket_redemptions` en `SCHEMA_FINAL.sql`. | Consultas en `sql/queries/governance.sql`, `GovernanceService` y endpoints `/api/v1/governance/*`. |
| **8. `realtime-service` (`CU14–CU16`)** | **10% Esqueleto** | Actualmente tiene un placeholder de 9 líneas en TypeScript/Bun (`realtime-service/index.ts`). | **Migrar a Go 1.24 + Chi v5 + WebSockets** (igual que se hizo con `image-service`), implementar `ArenaHub` en RAM por PIN, tickers de ronda (`RN10`), reconexión `10s` (`RNF-CON-02`) y *Batch Flush* hacia `core-service`. |
| **9. `image-service` (Multimedia)** | **25% Esqueleto** | Ya migrado a Go 1.24 + Chi v5 en puerto `8002` con `/health`, graceful shutdown y Dockerfile multi-stage. | Endpoint `POST /api/v1/images/upload` con validación de magic bytes, limpieza EXIF y compresión WebP. |

---

## 2. Plan de Implementación Paso a Paso (Checklist Operativo para Agentes)

Cada fase está diseñada como una unidad atómica de trabajo verificable mediante `sqlc generate && go build ./... && go test -v ./core-service/...`.

---

### 📍 FASE 3.3 (Prioridad Inmediata): Completar el Módulo de Diccionario y Curaduría por Quórum (`CU05, CU06, CU07, CU08`)

> **Objetivo:** Cerrar al 100% el corazón lexicográfico y editorial de Palabrario en `core-service` siguiendo las secciones **2.3** y **2.4** de [`INTERFACES.md`](./INTERFACES.md).

#### Tareas Exactas de Implementación:
1. **Modelos Compartidos (`pkg/models/`):**
   - Ampliar [`pkg/models/dictionary.go`](../pkg/models/dictionary.go) con `WordExploreFilter`, `WordSummaryCard`, `PaginatedWords`, `CreateExampleInput`, `CreateSenseRequest`, `CreateWordProposalRequest`, `UpdateWordRequest`, `UpdateSenseRequest`, `SemanticConcept`, `CreateConceptRequest`.
   - Crear `pkg/models/curation.go` con `ReviewDecision`, `SenseApprovalVote`, `SenseQuorumItem`, `SenseQuorumDetail`, `SubmitVoteRequest`, `AdminOverrideStateRequest`, `ContentReport`, `CreateReportRequest`, `ResolveReportRequest`.
2. **Consultas Tipadas (`sql/queries/`):**
   - Ampliar [`sql/queries/dictionary.sql`](../sql/queries/dictionary.sql) con `ListRegions`, `GetRegionByCode`, `GetWordOfTheDayID`, `ExploreWordsPaginated`, `CountExploreWords`, `GetWordByNormalizedLemmaAndPOS`, `CreateWord`, `UpdateWordMetadata`, `CreateSense`, `GetSenseByID`, `UpdateSense`, `DeleteSense`, `CreateExample`, `UpdateExample`, `DeleteExample`, `ListSemanticConcepts`, `CreateSemanticConcept`, `LinkSenseToConcept`, `UnlinkSenseFromConcept`.
   - Crear `sql/queries/curation.sql` con `ListPendingQuorumQueue`, `ListEditorProposals`, `GetSenseQuorumItemByID`, `ListVotesBySenseID`, `UpsertSenseApproval`, `AdminUpdateSenseState`.
   - Ejecutar `sqlc generate`.
3. **Servicios SOA (`core-service/internal/service/`):**
   - Crear `core-service/internal/service/errors.go` con los errores tipados de dominio y traductor de errores `pgconn.PgError` (especialmente `SQLSTATE 23505` y excepciones `P0001` lanzadas por `trg_prevent_self_approval`).
   - Crear `core-service/internal/service/dictionary_service.go` implementando `DictionaryService` (incluyendo helper determinístico `NormalizeLemma(s string) string` que pasa a minúsculas y remueve tildes para `normalized_lemma`).
   - Crear `core-service/internal/service/curation_service.go` implementando `CurationService`.
4. **Handlers y Router (`core-service/internal/api/`):**
   - Actualizar [`dictionary_handler.go`](../core-service/internal/api/handlers/dictionary_handler.go) y crear `curation_handler.go`.
   - Registrar en [`router.go`](../core-service/internal/api/router.go) todas las rutas de las secciones **2.3.3** y **2.4.3** de `INTERFACES.md`.
5. **Verificación y Pruebas (DoD):**
   - Crear `core-service/internal/api/dictionary_curation_test.go` probando:
     1. `GET /api/v1/regions` y `GET /api/v1/words/word-of-the-day`.
     2. Editor A crea una propuesta (`POST /api/v1/words`) $\to$ nace en `DRAFT` y aún no aparece como sentido público.
     3. Editor A intenta votar su propia propuesta (`POST /api/v1/curation/senses/{id}/vote`) $\to$ recibe `403 Forbidden` (**RN07**).
     4. Editores B, C y D votan `APPROVE` $\to$ al 3er voto el trigger `RN06` promueve automáticamente el sentido a `APPROVED` (**CU08**).

---

### 📍 FASE 3.4: Studio Personal, Motor Matemático SM-2 y Ledger de XP (`CU10, CU11, CU12, CU13`)

> **Objetivo:** Implementar la gestión de barajas/colecciones y el motor de repetición espaciada SuperMemo SM-2 siguiendo las secciones **2.5** y **2.6** de [`INTERFACES.md`](./INTERFACES.md).

#### Tareas Exactas de Implementación:
1. **Modelos Compartidos:** Ampliar [`pkg/models/studio.go`](../pkg/models/studio.go) con `StudioDeckSummary`, `CreateDeckRequest`, `UpdateDeckRequest`, `AddDeckItemRequest`, `CreateCollectionRequest`, `UpdateCollectionRequest`, `StudyCardItem`, `StudySessionPayload`, `SubmitReviewRequest`, `CardReviewResult`, `XPTransactionItem`, `AchievementStatus`, `MasteryDashboard`.
2. **Consultas Tipadas:**
   - Crear `sql/queries/studio.sql` (CRUD de `decks`, `deck_items`, `saved_decks`, `collections`, `collection_decks`, `saved_collections`).
   - Crear `sql/queries/mastery.sql` (tarjetas pendientes `next_review_at <= now()`, `UpsertUserMastery`, `InsertXPTransaction`, `AddUserXPAndUpdateStreak`, `GetMasteryStageCounts`, `ListUserAchievements`).
   - Ejecutar `sqlc generate`.
3. **Servicios SOA y Tests Unitarios:**
   - Crear `core-service/internal/service/studio_service.go` validando estrictamente **RN03** (propiedad), **RN04** (privacidad) y **RN05** (`is_official`).
   - Crear `core-service/internal/service/mastery_service.go` con la función pura `CalculateSM2` y crear `mastery_service_test.go` verificando todos los casos borde de la fórmula matemática SM-2 ($q \in \{1..4\}$, piso $EF \ge 1.3$, saltos de intervalo $1 \to 6 \to \lceil I \times EF \rceil$, transición `NEW` $\to$ `LEARNING` $\to$ `REVIEWING` $\to$ `MASTERED`).
4. **Handlers y Router:** Crear `studio_handler.go` y `mastery_handler.go` y montar todas las rutas de las secciones **2.5.3** y **2.6.3** de `INTERFACES.md`.

---

### 📍 FASE 3.5: Módulo de Cuenta, Perfiles Sociales, Gobernanza y Moderación (`CU01, CU02, CU03, CU04, CU09`)

> **Objetivo:** Completar la gestión de perfil propio/público, el flujo de ascenso a Editor, los boletos de invitación y la resolución de reportes siguiendo las secciones **2.1**, **2.2** y **2.4** de [`INTERFACES.md`](./INTERFACES.md).

#### Tareas Exactas de Implementación:
1. **Modelos y Consultas:**
   - Crear `pkg/models/governance.go` y ampliar `pkg/models/user.go`.
   - Ampliar `sql/queries/users.sql`, `sql/queries/curation.sql` (reportes) y crear `sql/queries/governance.sql`.
   - Ejecutar `sqlc generate`.
2. **Servicios SOA y Handlers:**
   - Crear `core-service/internal/service/identity_service.go` y `governance_service.go`.
   - Crear `user_handler.go` y `governance_handler.go` y montar los endpoints `/api/v1/users/*`, `/api/v1/governance/*` y `/api/v1/reports/*`.
   - Verificar con tests de integración **RN01** (perfil propio), **RN02** (unicidad de solicitud `PENDING`), protección del último `ADMIN` (`CU03 4a`) y redención atómica de boletos (`CU04`).

---

### 📍 FASE 4: Migración de `realtime-service` a Go 1.24 y Arena Multijugador (`CU14, CU15, CU16`)

> **Objetivo:** Unificar el 100% del backend en Go 1.24 migrando `realtime-service` (`:8001`) de Bun a Go e implementando el motor de salas por PIN en RAM y el *Batch Flush* transaccional hacia `core-service` según la sección **2.7** de [`INTERFACES.md`](./INTERFACES.md).

#### Tareas Exactas de Implementación:
1. **Migración de Infraestructura de `realtime-service` a Go 1.24:**
   - Eliminar los artefactos de Bun/Node en `realtime-service/` (`index.ts`, `package.json`, `bun.lock`, `tsconfig.json`, `node_modules`).
   - Crear `realtime-service/main.go` y `realtime-service/Dockerfile` (multi-stage build con `context: .` en [`docker-compose.yml`](../docker-compose.yml), idéntico en patrón a `core-service` e `image-service`).
   - Añadir `github.com/gorilla/websocket` (o `github.com/coder/websocket`) al `go.mod` raíz.
2. **Coordinación en `core-service` (`ArenaCatalogService`):**
   - Crear `sql/queries/arena.sql`, `arena_catalog_service.go` y `arena_handler.go`.
   - Implementar `POST /api/v1/arena/matches/prepare` (valida $\ge 4$ sentidos aprobados, genera preguntas híbridas y registra la sala `WAITING` con su `room_code`), `GET /api/v1/arena/matches/{roomCode}/package` y `POST /api/v1/arena/matches/consolidate` (*Batch Flush* ACID en `matches`, `match_participants`, `xp_transactions` y `user_mastery`).
3. **Motor WebSocket en `realtime-service` (`realtime-service/internal/arena/`):**
   - Implementar `hub.go` (`ArenaHub` con `sync.RWMutex`) y `room.go` (`ActiveRoom` con bucle de eventos por canales, temporizador por ronda, cálculo de puntaje por velocidad **RN10**, ventana de reconexión de 10 segundos **RNF-CON-02** y disparo automático del *Batch Flush* al emitir `PODIUM`).
4. **Completar `image-service` (`:8002`):**
   - Implementar `POST /api/v1/images/upload` validando magic bytes y limpiando metadatos EXIF (**RNF-SEG-04**).

---

### 📍 FASE 5: Integración End-to-End con el Cliente Nativo Wails v2 (`local-app`) y Pipeline CI/CD

> **Objetivo:** Una vez que el 100% de los endpoints de `INTERFACES.md` estén probados y estables en Go, reemplazar los datos mock de `local-app/frontend` por clientes HTTP/WS tipados contra los contratos reales.

1. Crear una capa de cliente API tipada en `local-app/frontend/src/services/api.ts` cuyos tipos TypeScript sean espejo exacto de [`pkg/models/*.go`](../pkg/models).
2. Conectar secuencialmente: Auth & Perfil $\to$ Diccionario & Quórum $\to$ Studio & SM-2 $\to$ Arena WebSockets.
3. Empaquetar el binario nativo con `wails build` y validar el despliegue completo en `docker compose up --build`.

---

### 🌱 FASE 6 (Incremento Post-MVP a Largo Plazo): Sistema de Eventos Culturales de Temporada

> **Objetivo:** Incorporar festivales estacionales (*Día de Muertos*, *Feria de las Flores*, *San Juan*) y galería de arte hispanoamericano según [`proposals/SEASONAL_EVENTS_SPECIFICATION.md`](../proposals/SEASONAL_EVENTS_SPECIFICATION.md), únicamente tras haber cerrado y verificado las Fases 3.3 a 5.
