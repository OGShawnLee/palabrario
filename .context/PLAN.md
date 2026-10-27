# 🗺️ Plan Maestro de Implementación y Estado Arquitectónico (`.context/PLAN.md`)

> **Propósito de este documento:**  
> Este archivo es el **tablero de control operativo y mapa de ruta paso a paso** de Palabrario.  
> Cada agente que inicie una sesión de trabajo debe consultar este documento para saber con precisión milimétrica **qué fase está completada**, **cuál es la prioridad inmediata** y **qué alcance exacto le corresponde construir**, evitando duplicidad de esfuerzo o dispersión.

---

## 0. Principio Rector de Ingeniería: Backend-First & Cadena Tipada Isomórfica (`Go -> Wails Bindings (@bindings) -> @palabrario/ui -> Desktop / Mobile / Web SSR`)

Palabrario sigue una metodología estricta de construcción en cascada tipada sobre un monorepo híbrido (**Go Workspaces + Bun Workspaces**):
1. **Base de Datos Relacional y Contratos de Dominio:** Diseñados e indexados en PostgreSQL 17 ([`SCHEMA_FINAL.sql`](../SCHEMA_FINAL.sql)) y formalizados en Go en [`pkg/models`](../pkg/models).
2. **Consultas Estrictamente Tipadas (`sqlc`) y Servicios SOA (`:8000`, `:8001`, `:8002`):** Cero strings SQL sueltos y 100% de lógica de negocio verificada en el backend (Go 1.24 + Chi v5).
3. **Puente Nativo IPC de Wails V3 (`local-app/app.go`):** El proceso Go de escritorio importa `pkg/models`, custodia el `sessionToken` JWT en memoria privada (`sync.RWMutex`) y expone los métodos de los 7 módulos SOA + subida multimedia.
4. **Autogeneración de Tipos y Contrato Polimórfico (`IPalabrarioAPI`):** Wails genera `local-app/frontend/bindings/` (accesible en todo el monorepo vía alias universal `@bindings/*`), de donde se deriva `export type IPalabrarioAPI = typeof WailsAppBridge` para gobernar tanto a `WailsDesktopAdapter` (Desktop IPC), como a `CapacitorMobileAdapter` (Móvil / Navegador Web SPA con `fetch` + `WebSocket` nativo a `:8001`) y `AstroSSRAdapter` (SSR en Fase 6).
5. **Paquete Compartido de Diseño y Componentes (`packages/ui` $\to$ `@palabrario/ui`):** Aloja la configuración canónica de UnoCSS ([`.context/PALETTE.md`](./PALETTE.md)), la tipografía `Satoshi` y **todos los componentes `.vue` reutilizables y agnósticos del router (SSR-Safe)**, importando tipos directamente desde `@bindings/palabrario/pkg/models`.
6. **Shells Multiplataforma (`local-app/frontend` en Fase 5 y `web-portal` con Astro en Fase 6):** Consumen `@palabrario/ui` (`workspace:*`) y `IPalabrarioAPI`, cumpliendo las rutas de [`.context/ROUTES_MAP.md`](./ROUTES_MAP.md) y los flujos de [`.context/UI_FLOWS.md`](./UI_FLOWS.md).

---

## 1. Diagnóstico Ejecutivo e Histórico: ¿Dónde Estamos Hoy?

El proyecto ha completado y verificado al **100% toda la arquitectura de Backend (Fases 1 a 4)** cubriendo los **16 Casos de Uso (`CU01–CU16`)** y las **13 Reglas de Negocio (`RN01–RN13`)**. Actualmente nos encontramos en el arranque de la **Fase 5 (Monorepo `@palabrario/ui` + Cliente Wails V3 / Capacitor PWA)**, con la **Fase 6 (Plataforma Web 100% Híbrida: Astro SSR SEO + 23 Rutas Interactivas con WebSockets)** y la **Fase 7 (Eventos Culturales de Temporada)** planificadas en la hoja de ruta.

```mermaid
flowchart LR
    F1["✅ Fase 1\nEsquema SQL + Triggers\n(100%)"] --> F2["✅ Fase 2\nMonorepo Go + Contratos\n(100%)"]
    F2 --> F3["✅ Fase 3 (3.1 a 3.5)\nCore SOA Completo\n(100%)"]
    F3 --> F4["✅ Fase 4 (4.1 a 4.5)\nArena 3NF + WS + Media\n(100%)"]
    F4 --> F5["🚧 Fase 5 (5.1 a 5.9)\n@palabrario/ui + Wails V3\n+ Capacitor SPA"]
    F5 --> F6["🌐 Fase 6\nPlataforma Web 100%\n(Astro SSR SEO + SPA WS)"]
    F6 --> F7["🌱 Fase 7 (Largo Plazo)\nEventos de Temporada"]
```

### 1.1. Breakdown Histórico de Fases Completadas (Backend 100% Cerrado)

- **✅ Fase 1 & Fase 2 — Cimientos Relacionales y Monorepo Go:**
  - Esquema PostgreSQL 17 ([`SCHEMA_FINAL.sql`](../SCHEMA_FINAL.sql)) con 24 tablas en 3NF, triggers `RN06`/`RN07`, búsqueda FTS + `pg_trgm` y paquetes compartidos (`pkg/auth`, `pkg/models`, `pkg/response`, `pkg/utils`).
- **✅ Fase 3.1 a 3.3 — Core SOA Base, Diccionario y Curaduría por Quórum (`CU05–CU08`):**
  - `DictionaryService` (FTS, palabra del día, explorador paginado, autoría de lemas/sentidos/ejemplos/conceptos) y `CurationService` (cola de quórum tripartito `RN06`, anti-autoaprobación `RN07`, veto admin `RN08`).
- **✅ Fase 3.4 — Studio Personal, Motor Matemático SM-2 y Ledger de XP (`CU10–CU13`):**
  - `StudioService` (`RN03`, `RN04`, `RN05`) y `MasteryService` (`CalculateSM2` puro `RN09`, transacciones ACID con libro mayor inmutable `xp_transactions` `RN11` y rachas).
- **✅ Fase 3.5 — Identidad, Perfiles Sociales, Gobernanza y Moderación (`CU01–CU04, CU09`):**
  - `IdentityService` (perfiles público/privado duales), `GovernanceService` (`RN02`, promoción a `EDITOR`, protección del último `ADMIN`, boletos con `FOR UPDATE`) y resolución de reportes comunitarios (`CU09`).
- **✅ Fase 4 (Sub-fases 4.1 a 4.5) — Arena Multijugador 3NF, Seguridad S2S, `realtime-service` e `image-service` (`CU14–CU16`):**
  - Sub-llave criptográfica S2S (`sha256`, TTL 30s), `ArenaCatalogService` (piso de 8 sentidos, desempate determinista, Batch Flush ACID y Replay 3NF), `realtime-service` (`:8001`, `ArenaHub` en RAM, reloj monotónico `RN10`, reconexión 10s `RNF-CON-02`, sucesión de Host), `image-service` (`:8002`, magic bytes, limpieza EXIF, moderación heurística NSFW `RN12`) y modalidad visual `WORD_TO_IMAGE` con degradación elegante (`RN13`).

### 1.2. Matriz de Estado Real por Capa y Servicio

| Componente / Módulo | Estado Actual | Resumen de Verificación |
| :--- | :---: | :--- |
| **1. Base de Datos ([`SCHEMA_FINAL.sql`](../SCHEMA_FINAL.sql)) y `sqlc`** | **100% Completo** | 24 tablas en 3NF, triggers `RN06`/`RN07` y código Go autogenerado en `core-service/internal/db/`. |
| **2. Paquetes Compartidos ([`pkg/`](../pkg))** | **100% Completo** | `pkg/auth`, `pkg/response`, `pkg/utils` y todos los DTOs canónicos en `pkg/models/`. |
| **3. `core-service` (`:8000` — `CU01–CU14, CU16`)** | **100% Completo** | Los 7 servicios SOA implementados, enrutados en `router.go` y con suite de pruebas pasando al 100%. |
| **4. `realtime-service` (`:8001` — `CU15, CU16`)** | **100% Completo** | `ArenaHub` WebSocket concurrente en RAM con tests pasando con `-race`. |
| **5. `image-service` (`:8002` — `RN12 / RNF-SEG-04`)** | **100% Completo** | Sanitización EXIF, redimensionado y filtro NSFW en Go puro con tests pasando con `-race`. |
| **6. Monorepo UI (`packages/ui`) + Cliente (`local-app/`)** | **🚧 En Construcción (Fase 5)** | Inicializado con Wails V3 + Bun; listo para estructurar en **9 sub-fases atómicas (5.1 a 5.9)** bajo la arquitectura de [`.context/UI_SPECIFICATION.md`](./UI_SPECIFICATION.md), [`.context/ROUTES_MAP.md`](./ROUTES_MAP.md), [`.context/UI_FLOWS.md`](./UI_FLOWS.md) y [`.context/PALETTE.md`](./PALETTE.md). |
| **7. Plataforma Web Híbrida (`web-portal/` — Astro SSR + SPA)** | **⏳ Planificada (Fase 6)** | Consumirá `@palabrario/ui` y `IPalabrarioAPI` para ofrecer SSR/SEO + OpenGraph en rutas públicas (`#3–#10`) y el 100% de las 23 rutas interactivas (incluyendo WebSockets hacia `:8001`) en el navegador. |

---

## 2. Plan de Implementación Detallado de la Fase 5 (9 Sub-Fases Atómicas)

**Metodología de Ejecución:** Siguiendo [`AGENT_WORKFLOW.md`](../AGENT_WORKFLOW.md), hemos granularizado la Fase 5 en **9 sub-fases atómicas, independientes y fáciles de auditar**.  
Cada sub-fase corresponde a **1 sesión limpia con un agente constructor** y debe cerrar obligatoriamente con:
```powershell
wails3 task client:check # (o powershell -File .\scripts\client.ps1 check)
```
con código de salida `0` (cadena automatizada de bindings Go -> TS, chequeo de tipos de `@palabrario/ui` + `local-app/frontend`, empaquetado Vite y `go build .`).

```mermaid
flowchart TD
    F51["Sub-Fase 5.1\nBun Workspaces + @palabrario/ui\n(UnoCSS PALETTE.md) + Base Wails V3"] --> F52["Sub-Fase 5.2\nPuente Go (app.go) + Bindings TS\n+ Adaptador Dual (Wails/Capacitor)"]
    F52 --> F53["Sub-Fase 5.3\nLayouts, Router (Placeholders),\nTopBar Dropdown + Popover Credenciales\ny Auth/Perfil (Rutas #1–#3)"]
    F53 --> F54["Sub-Fase 5.4\nPortal Home, Buscador FTS y\nFicha Vocabulary.com (Rutas #4–#6)"]
    F54 --> F55["Sub-Fase 5.5\nCatálogos Públicos de Colecciones\ny Barajas + Clonación (Rutas #7–#10)"]
    F55 --> F56["Sub-Fase 5.6\nStudio Personal: Gestores y Editores\nde Barajas y Colecciones (Rutas #11, #13–#16)"]
    F56 --> F57["Sub-Fase 5.7\nSesión Inmersiva SM-2 y Tablero\nde Maestría / Ledger XP (Rutas #12, #17)"]
    F57 --> F58["Sub-Fase 5.8\nArena Multijugador WS en Vivo,\nWORD_TO_IMAGE y Replay 3NF (Rutas #18–#20)"]
    F58 --> F59["Sub-Fase 5.9\nPortal Editorial (Quórum 3, Autoría)\ny Gobernanza Admin (Rutas #21–#23)"]
```

---

### 📍 Sub-Fase 5.1: Monorepo Bun Workspaces (`packages/ui` $\to$ `@palabrario/ui`) + Inicialización de `local-app` (Wails V3 + Capacitor + UnoCSS) — **[SIGUIENTE EN COLA]**

- **Objetivo:** Plantar los cimientos del monorepo frontend con **Bun Workspaces**, creando el paquete compartido `packages/ui` (`@palabrario/ui`) con el sistema de diseño *"Imprenta Táctil Ilustrada"*, el alias universal `@bindings/*` en TypeScript y configurando `local-app/` (Wails V3 + Capacitor + Vue 3) para consumirlo sin fricción.
- **Entregables y Checklist:**
  - [ ] **Raíz del Monorepo Frontend (`package.json` y `tsconfig.base.json`):**
    - `package.json` en la raíz del repositorio con `"workspaces": ["packages/*", "local-app/frontend"]` (listo para añadir `"web-portal"` en Fase 6).
    - `tsconfig.base.json` con configuración TypeScript estricta compartida y mapeo de alias `"paths": { "@bindings/*": ["local-app/frontend/bindings/*"] }`.
  - [ ] **Paquete Compartido `packages/ui/` (`@palabrario/ui`):**
    - `packages/ui/package.json` (`name: "@palabrario/ui"`, `exports` para `.`, `./uno.config`, `./styles.css`) con dependencias `vue`, `reka-ui`, `lucide-vue-next`, `unocss`, `@unocss/preset-uno`, `@unocss/preset-icons` y `@iconify-json/lucide`.
    - `packages/ui/src/uno.config.ts`: Configuración canónica exacta de [`.context/PALETTE.md`](./PALETTE.md) (tokens `paper`, `surface`, `ink`, `brand`, `brand-soft`, `success`, `warning`, `info`, sombras `tactile-*` y shortcuts `card-editorial`, `btn-brand`, `btn-surface`, `badge-pill`, `input-editorial`) incluyendo `content.filesystem` apuntando a `packages/ui/src/components/**/*.{vue,ts}` para extracción cross-workspace garantizada.
    - `packages/ui/src/assets/main.css`: Tipografía `Satoshi`, variables editoriales y `.tabular-nums`.
    - Estructura de carpetas en `packages/ui/src/components/`: `ui/`, `dictionary/`, `studio/`, `arena/`, `editorial/` y `packages/ui/src/index.ts` (barrel export de componentes). Los componentes importan directamente sus tipos desde `@bindings/palabrario/pkg/models`.
  - [ ] **Capa Go de Escritorio (`local-app/go.mod`, `local-app/main.go`, `local-app/app.go` base y `local-app/Taskfile.yml`):**
    - Aplicación de escritorio Go Wails V3 vinculada a `palabrario/pkg/models`, ventana nativa (`1280x820`, min `390x680`, color `#F6F3EC` `paper`) y estructura `App` con `sessionToken` privado protegido por `sync.RWMutex` (reemplazando el `GreetService` de plantilla).
    - `local-app/Taskfile.yml` actualizado en la tarea `check` para validar `bindings` $\to$ `type-check` en `../packages/ui` $\to$ `type-check` en `frontend` $\to$ `build` en `frontend` $\to$ `go build .`.
  - [ ] **Shell Frontend (`local-app/frontend/`):**
    - `package.json` dependiendo de `"@palabrario/ui": "workspace:*"`, `vue`, `vue-router`, `@wailsio/runtime`, `@capacitor/core`, `@capacitor/cli`, `@capacitor/preferences`, `@vitejs/plugin-vue` y `unocss`.
    - `local-app/frontend/uno.config.ts` re-exportando `@palabrario/ui/uno.config`.
    - `local-app/frontend/vite.config.ts` con alias `@` -> `./src` y `@bindings` -> `./bindings`.
    - `capacitor.config.ts` apuntando a `dist`.
    - Estructura de directorios en `local-app/frontend/src/`: `api/`, `composables/`, `components/ui/`, `layouts/`, `router/`, `views/` (con subcarpetas `auth/`, `profile/`, `community/`, `studio/`, `arena/`, `editorial/`) y un `App.vue` mínimo verificable importando estilos de `@palabrario/ui`.
  - [ ] **`local-app/README.md`:** Documentación clara de la arquitectura del monorepo (`@palabrario/ui` + `local-app/frontend`), flujo dual (Wails Desktop vs Capacitor/Web) y comandos de verificación.
- **Criterio de Cierre:** `wails3 task client:check` (o `.\scripts\client.ps1 check`) pasa con 0 errores.

---

### 📍 Sub-Fase 5.2: Puente Nativo Go (`local-app/app.go`), Bindings TypeScript y Adaptador Polimórfico (`IPalabrarioAPI`)

- **Objetivo:** Construir la capa de comunicación 100% tipada que conecta el cliente con `core-service (:8000)`, `realtime-service (:8001)` e `image-service (:8002)` y disponibilizar los bindings para todos los paquetes mediante `@bindings/*`.
- **Entregables y Checklist:**
  - [ ] **`local-app/app.go` Completo:**
    - Gestión de sesión en memoria (`HasActiveSession`, `SignOut`, `GetMyProfile` con auto-hidratación).
    - Métodos IPC que envuelven los **7 módulos de `core-service (:8000)`** retornando tipos exactos de `pkg/models` (`auth.go`, `identity.go`, `governance.go`, `dictionary.go`, `curation.go`, `studio.go`, `arena.go`).
    - Método `UploadImage(category, filename string, fileBytes []byte) (*models.ImageUploadResponse, error)` hacia `image-service (:8002)`.
    - Cliente WebSocket Go (`ConnectArenaRoom`, `SubmitArenaAnswer`, `LeaveArenaRoom`) hacia `realtime-service (:8001)` que inyecta el `sessionToken` privado y reenvía cada `models.WSMessage` como evento IPC (`arena:ws_message`, `arena:ws_disconnected`).
  - [ ] **Generación de Bindings TypeScript (`local-app/frontend/bindings/`):**
    - Generados automáticamente con `wails3 task client:bindings` y consumibles directamente por `@palabrario/ui` y `local-app/frontend` vía `@bindings/palabrario/pkg/models`.
  - [ ] **Adaptador Polimórfico Dual (`local-app/frontend/src/api/`):**
    - `client.ts`: Define `import * as WailsAppBridge from '@bindings/palabrario/local-app/app'; export type IPalabrarioAPI = typeof WailsAppBridge;` y exporta el singleton `api` que selecciona automáticamente entre `WailsDesktopAdapter` y `CapacitorMobileAdapter`.
    - `wails.adapter.ts`: Delega 1:1 en `WailsAppBridge`.
    - `capacitor.adapter.ts`: Implementa exactamente `IPalabrarioAPI` usando `fetch` tipado hacia `:8000` y `:8002`, custodiando el JWT con `@capacitor/preferences` (funcional tanto en Capacitor móvil como en Navegador Web SPA).
    - `arena.events.ts`: Implementa `ArenaEventBus` unificando `Events.On("arena:ws_message")` en Wails y `new WebSocket("ws://.../ws/arena/{roomCode}?token={jwt}")` nativo en Capacitor/Web Browser.
- **Criterio de Cierre:** `wails3 task client:check` pasa con 0 errores, verificando que `WailsDesktopAdapter` y `CapacitorMobileAdapter` satisfacen al 100% el contrato `IPalabrarioAPI`.

---

### 📍 Sub-Fase 5.3: Shell de Navegación (`Layouts` + `Router` con Placeholders), `TopBar` + `CredentialsRolePopover` y Autenticación/Perfil (`Rutas #1, #2, #3` — `CU01, CU02, CU04`)

- **Objetivo:** Levantar el esqueleto completo de enrutamiento de las 23 rutas, los componentes de navegación en `@palabrario/ui`, el flujo de credenciales/rol en la imprenta, y cerrar las vistas de Autenticación y Perfil Social.
- **Entregables y Checklist:**
  - [ ] **Componentes Compartidos en `packages/ui/src/components/ui/` (SSR-Safe, exportados en `packages/ui/src/index.ts`):**
    - `TopBar.vue`: Recibe `currentUser`, `activeSpace` y contadores por props; incluye conmutador de espacio, racha (`🔥`), nivel, disparador `Ctrl + K` y **`DropdownMenu` de Reka UI** en el avatar/estado del usuario emitiendo `@navigate(path)`, `@open-search`, `@open-credentials` y `@sign-out`.
    - `BottomNav.vue`: Barra inferior responsiva para móvil/ventana estrecha que emite `@navigate(path)`.
    - `MediaDropzone.vue`: Dropzone reutilizable (`category: 'avatar' | 'credentials' | 'general'`, validación máx 5 MB y formatos `JPG/PNG/WEBP`).
  - [ ] **Estado, Layouts y Enrutador en `local-app/frontend/src/`:**
    - `src/composables/useAuth.ts`: Estado reactivo global (`currentUser`, `isAuthenticated`, `isEditor`, `isAdmin`, `hydrateSession`, `signIn`, `signUp`, `signOut`, `refreshProfile`) siguiendo el **Flujo 1 de [`.context/UI_FLOWS.md`](./UI_FLOWS.md)**.
    - `src/components/ui/CredentialsRolePopover.vue`: Popover/Dialog invocado desde `AppLayout.vue` al recibir `@open-credentials` de `TopBar.vue` siguiendo el **Flujo 2 de [`.context/UI_FLOWS.md`](./UI_FLOWS.md)** (muestra rol actual, pestaña *Canjear Boleto VIP* `CU04`, pestaña *Solicitar Ascenso a Editor* `CU02` con bloqueo `RN02` si `xp < 500 && !is_verified` y subida de credencial `RN12`, e historial `ListMyEditorRequests`).
    - `src/layouts/AppLayout.vue`, `AuthLayout.vue` e `ImmersiveLayout.vue`.
    - Vistas Placeholder (`Rutas #4 a #23`) en `local-app/frontend/src/views/` y `src/router/index.ts` con `createWebHashHistory()` y guardias (`OptionalAuth`, `guestOnly`, `requireAuth`, `requireRole`).
  - [ ] **Vistas Funcionales en `local-app/frontend/src/views/` (`Rutas #1, #2 y #3`):**
    - `views/auth/SignInView.vue` (`/#/auth/sign-in` — `RF01`).
    - `views/auth/SignUpView.vue` (`/#/auth/sign-up` — selector `ListRegions`).
    - `views/profile/UserProfileView.vue` (`/#/u/:username` — `CU01`): Ficha pública + panel de edición de perfil (`UpdateProfile` + `MediaDropzone(category='avatar')`) y cambio de contraseña (`ChangePassword`) cuando `is_owner == true`.
- **Criterio de Cierre:** Navegación funcional entre las 23 rutas, login/registro/logout operativos, Popover de Credenciales funcional desde `TopBar`, y `wails3 task client:check` en verde (0 errores).

---

### 📍 Sub-Fase 5.4: Portal Comunitario, Buscador FTS y Ficha Léxica Estilo Vocabulary.com (`Rutas #4, #5, #6` + Modales — `CU05, CU09, CU10`)

- **Objetivo:** Construir en `@palabrario/ui` los componentes del diccionario panhispánico (100% SSR-Safe, importando directamente `@bindings/palabrario/pkg/models`) y reemplazar en `local-app/frontend` los placeholders de las rutas `#4, #5 y #6`.
- **Entregables y Checklist:**
  - [ ] **Componentes de Diccionario y Modales en `packages/ui/src/components/` (SSR-Safe, exportados en `packages/ui/src/index.ts`):**
    - `dictionary/WordCard.vue`: Tarjeta táctil editorial (`models.WordSummaryCard`) que emite `@select(idWord)`.
    - `dictionary/FriendlyBlurb.vue`: Apertura narrativa estilo Vocabulary.com (Párrafo 1 con `word.friendly_blurb` en `text-lg/xl leading-relaxed`, Párrafo 2 con `word.etymology` e ilustración opcional `word.image_url` que emite `@open-image`).
    - `dictionary/HeteronymPills.vue`: Píldoras de equivalentes dialectales (`models.HeteronymItem[]`) que emiten `@select-word(idWord)` (sin acoplarse a `vue-router`).
    - `dictionary/SenseItem.vue`: Ficha numerada dentro del contenedor con listón `bg-brand` (*"Acepciones de {lemma}"*), pills de región/categoría/CEFR (`bg-success/15 text-success border-success`), definición, `HeteronymPills` y desplegable `Collapsible` de Reka UI (*"Ver más ⌄"*) con citas de uso (`examples`), `usage_note`, miniatura clicable y botones que emiten `@add-to-deck(sense)` y `@report(sense)`.
    - `dictionary/TabooWarningDialog.vue`: `AlertDialog` preventivo cuando `word.is_taboo == true`.
    - `dictionary/ReportModal.vue` (`CU09`): Modal Reka `Dialog` para denunciar `WORD` o `SENSE` (`reason_category`: `SPAM`, `OFFENSIVE`, `INCORRECT`, `DUPLICATE`, `OTHER`) emitiendo `@submit(payload)`.
    - `studio/AddToDeckModal.vue` (`CU10`): Modal rápido que recibe `decks` por prop y emite `@add` o `@create-deck` respetando `RN04`.
    - `ui/ImageLightboxModal.vue` y `ui/SearchModal.vue` (`Ctrl + K` con debounce de 250ms).
  - [ ] **Vistas Funcionales en `local-app/frontend/src/views/community/` (`Rutas #4, #5 y #6`):**
    - `HomeView.vue` (`/#/`): Palabra del Día (`GetWordOfTheDay`), accesos rápidos y escaparates de colecciones/decks públicos.
    - `DictionaryView.vue` (`/#/dictionary`): Buscador FTS + `pg_trgm` (`SearchWords`) y explorador paginado (`ExploreWords`) con filtros por país (`ListRegions`), categoría gramatical y nivel CEFR.
    - `WordDetailView.vue` (`/#/dictionary/:id`): Ficha completa integrando `FriendlyBlurb`, lista de `SenseItem` y modales según el **Flujo 3 de [`.context/UI_FLOWS.md`](./UI_FLOWS.md)**.
- **Criterio de Cierre:** `wails3 task client:check` pasa con 0 errores; trazabilidad completa de `CU05`, `CU09` y `CU10`.

---

### 📍 Sub-Fase 5.5: Catálogos Públicos de Colecciones y Barajas + Clonación (`Rutas #7, #8, #9, #10` — `CU10, CU11, RN03`)

- **Objetivo:** Crear las tarjetas reutilizables de barajas/colecciones en `@palabrario/ui` y reemplazar los placeholders de las rutas `#7 a #10` en `local-app/frontend`.
- **Entregables y Checklist:**
  - [ ] **Componentes en `packages/ui/src/components/studio/`:** `DeckCard.vue` y `CollectionCard.vue` (agnósticos del router, emitiendo `@select`).
  - [ ] **`local-app/frontend/src/views/community/PublicCollectionsView.vue` (`/#/collections` — Ruta #7):** Catálogo paginado de colecciones públicas (`ExplorePublicCollections`) con buscador y filtro de oficiales (`is_official`).
  - [ ] **`local-app/frontend/src/views/community/CollectionDetailView.vue` (`/#/collections/:id` — Ruta #8):** Portada editorial de la colección (`GetCollectionDetail`), descripción y grilla de barajas que la componen con enlace a cada una.
  - [ ] **`local-app/frontend/src/views/community/PublicDecksView.vue` (`/#/decks` — Ruta #9):** Catálogo paginado de barajas comunitarias (`ExplorePublicDecks`) mostrando autor, conteo de tarjetas (`word_count`) y contador de clones (`clone_count`).
  - [ ] **`local-app/frontend/src/views/community/PublicDeckDetailView.vue` (`/#/decks/:id` — Ruta #10):** Vista previa de las tarjetas (`SenseSummary`) de una baraja pública (`GetDeckDetail`) y botón táctil **"Clonar a Mi Studio"** (`ClonePublicDeck` — `RN03`) que redirige al nuevo deck clonado o inicia repaso.
- **Criterio de Cierre:** `wails3 task client:check` limpio (0 errores); clonación de decks públicos verificada contra `IPalabrarioAPI`.

---

### 📍 Sub-Fase 5.6: Studio Personal — Gestores y Editores de Barajas y Colecciones (`Rutas #11, #13, #14, #15, #16` — `CU10, CU11`)

- **Objetivo:** Construir en `local-app/frontend/src/views/studio/` el taller privado del aprendiz para administrar sus barajas de estudio (`RN04`) y colecciones (`RN05`).
- **Entregables y Checklist:**
  - [ ] **`views/studio/StudioDashboardView.vue` (`/#/studio` — Ruta #11):** Resumen ejecutivo del aprendiz (`GetMasterySummary` con conteo `due_for_review_today`, botón directo a `/#/studio/session`, racha activa y carruseles de mis decks y colecciones recientes usando `@palabrario/ui`).
  - [ ] **`views/studio/StudioDecksView.vue` (`/#/studio/decks` — Ruta #13):** Listado de barajas propias y clonadas (`ListMyDecks`), modal para crear baraja (`CreateDeck`) y eliminación con confirmación (`DeleteDeck`).
  - [ ] **`views/studio/StudioDeckEditorView.vue` (`/#/studio/decks/:id` — Ruta #14):** Editor de metadatos (`UpdateDeck`: nombre, descripción, switch `is_public`), buscador integrado para añadir acepciones (`AddSenseToDeck` respetando tope de 100 cartas `RN04`), remoción de cartas (`RemoveSenseFromDeck`) y botón *"Estudiar este Deck"* (`/#/studio/session?deckId=:id`).
  - [ ] **`views/studio/StudioCollectionsView.vue` (`/#/studio/collections` — Ruta #15):** Listado de colecciones propias (`ListMyCollections`), creación (`CreateCollection`) y borrado (`DeleteCollection`).
  - [ ] **`views/studio/StudioCollectionEditorView.vue` (`/#/studio/collections/:id` — Ruta #16):** Editor de colección (`UpdateCollection`), vinculación ordenada de barajas (`AddDeckToCollection` validando `RN05`: colecciones públicas solo admiten decks públicos) y desvinculación (`RemoveDeckFromCollection`).
- **Criterio de Cierre:** `wails3 task client:check` limpio (0 errores); reglas `RN04` y `RN05` reflejadas en la UI.

---

### 📍 Sub-Fase 5.7: Sesión Inmersiva de Repaso SM-2 y Tablero de Maestría (`Rutas #12, #17` — `CU12, CU13, RN09, RN11`)

- **Objetivo:** Implementar en `@palabrario/ui` los componentes de flashcards SM-2 y barra de maestría, y construir en `local-app/frontend` la sesión inmersiva y el libro mayor de XP.
- **Entregables y Checklist:**
  - [ ] **Componentes de Studio en `packages/ui/src/components/studio/` (exportados en `packages/ui/src/index.ts`):**
    - `StudyCard.vue`: Tarjeta interactiva con volteo anverso/reverso y botonera ergonómica de calificación SM-2 con atajos de teclado (`1: Olvidado`, `2: Difícil`, `3: Bien`, `4: Fácil`).
    - `MasteryProgressBar.vue`: Barra segmentada visualizando la proporción de palabras en estados `NEW`, `LEARNING`, `REVIEW` y `MASTERED`.
  - [ ] **`local-app/frontend/src/views/studio/StudySessionView.vue` (`/#/studio/session` — Ruta #12, Layout `immersive`):**
    - Carga cola de estudio (`GetStudySessionQueue(deckId, limit)`) siguiendo el **Flujo 4 de [`.context/UI_FLOWS.md`](./UI_FLOWS.md)**.
    - Procesa cada calificación con `SubmitCardReview`, muestra micro-feedback de XP ganado (`+10 XP`) y nuevo intervalo en días (`interval_days`), y al terminar presenta el resumen de sesión con bonificación por mazo completado (`+25 XP`) o racha diaria (`+15 XP`).
  - [ ] **`local-app/frontend/src/views/studio/StudioMasteryView.vue` (`/#/studio/mastery` — Ruta #17):**
    - Estadísticas globales de retención (`GetMasterySummary` + `MasteryProgressBar`).
    - Tabla del libro mayor auditable de experiencia (`GetMyXPTransactions` — `RN11`) y vitrina de medallas/logros desbloqueados (`GetMyAchievements`).
- **Criterio de Cierre:** `wails3 task client:check` limpio (0 errores); sesión SM-2 operable por click y teclado (`1..4`, `Espacio`).

---

### 📍 Sub-Fase 5.8: Arena Multijugador en Vivo (WebSockets), `WORD_TO_IMAGE` y Replay 3NF (`Rutas #18, #19, #20` — `CU14, CU15, CU16`)

- **Objetivo:** Construir en `@palabrario/ui` los componentes competitivos de Arena y en `local-app/frontend` la experiencia síncrona en tiempo real conectada a `ArenaEventBus` junto al visor de Replay 3NF.
- **Entregables y Checklist:**
  - [ ] **Componentes de Arena en `packages/ui/src/components/arena/` (exportados en `packages/ui/src/index.ts`):**
    - `ArenaTimer.vue`: Cronómetro regresivo con `font-variant-numeric: tabular-nums` y alerta visual en los últimos 3 segundos.
    - `ImageChoiceGrid.vue`: Cuadrícula táctil 2x2 de opciones ilustradas para rondas `WORD_TO_IMAGE` (`RN13`).
    - `PodiumCard.vue`: Tarjeta de podio final (`1°`, `2°`, `3°`) con desglose de puntos, aciertos, tiempo total (`TotalTimeSpentMs`) y XP otorgado.
  - [ ] **`local-app/frontend/src/views/arena/ArenaLobbyView.vue` (`/#/arena` — Ruta #18):**
    - Formulario para crear partida (`CreateArenaMatch`: selección de Deck o Colección con validación de mínimo 8 sentidos, rondas, tiempo por ronda, modos permitidos incluyendo `WORD_TO_IMAGE`) o unirse con código PIN de 6 caracteres (`JoinArenaMatch`).
    - Tabla de historial de partidas recientes (`GetMyMatchHistory`) con acceso directo al Replay.
  - [ ] **`local-app/frontend/src/views/arena/ArenaGameView.vue` (`/#/arena/room/:code` — Ruta #19, Layout `immersive`):**
    - Implementa el **Flujo 5 de [`.context/UI_FLOWS.md`](./UI_FLOWS.md)** suscribiéndose a `ArenaEventBus`.
    - Maneja todas las fases de la sala: `WAITING`, `COUNTDOWN`, `PLAYING` (renderizado condicional de opciones de texto o `ImageChoiceGrid` 2x2 + `SubmitArenaAnswer`), `ROUND_RESULT` y `FINISHED` (`PodiumCard` + enlace a `/#/arena/matches/:matchId`).
    - Maneja desconexión temporal / reconexión en ventana de gracia de 10s (`RNF-CON-02`) y salida limpia en `onUnmounted` (`LeaveArenaRoom`).
  - [ ] **`local-app/frontend/src/views/arena/ArenaReplayView.vue` (`/#/arena/matches/:matchId` — Ruta #20):**
    - Consume `GetMatchReplayDetail(matchId)` para mostrar el podio histórico y un acordeón/pestañas ronda por ronda (`MatchRoundReplay`) detallando la pregunta, respuesta correcta y qué contestó cada jugador con su tiempo en milisegundos (`time_spent_ms`).
- **Criterio de Cierre:** `wails3 task client:check` limpio (0 errores); manejo completo de las tramas WS (`ROOM_STATE`, `COUNTDOWN`, `QUESTION`, `ROUND_RESULT`, `PODIUM`, `ERROR`).

---

### 📍 Sub-Fase 5.9: Portal Editorial de Quórum Tripartito y Gobernanza Institucional (`Rutas #21, #22, #23` — `CU02–CU04, CU06–CU09`)

- **Objetivo:** Construir en `@palabrario/ui` los componentes de inspección de quórum y cerrar en `local-app/frontend` las herramientas de curaduría lexicográfica y gobernanza administrativa.
- **Entregables y Checklist:**
  - [ ] **Componentes Editoriales en `packages/ui/src/components/editorial/` (exportados en `packages/ui/src/index.ts`):**
    - `QuorumStepper.vue`: Indicador visual de los 3 slots de revisión (`RN06`) mostrando votos emitidos (`APPROVE`, `REQUEST_CHANGES`, `REJECT`), bloqueo de auto-aprobación (`RN07` cuando `currentUser.id_user == proposal.created_by`) y badge de veto administrativo (`RN08`).
    - `ProposalDiffCard.vue`: Inspector de propuesta léxica mostrando lema, etimología, sentidos, notas culturales, ejemplos e imágenes con atribuciones (`RN12`).
  - [ ] **`local-app/frontend/src/views/editorial/EditorDashboardView.vue` (`/#/editor` — Ruta #21, `EDITOR`/`ADMIN`):**
    - Cola paginada de propuestas en `PENDING_QUORUM` (`ListPendingQuorumQueue`), pestaña de *"Mis Propuestas"* (`ListMyLexicalProposals`), inspección con `ProposalDiffCard` + `QuorumStepper`, formulario de dictamen (`SubmitProposalReview`) y botón de resolución forzada exclusivo para `ADMIN` (`AdminForceResolveProposal` — `RN08`).
  - [ ] **`local-app/frontend/src/views/editorial/EditorWordFormView.vue` (`/#/editor/words/new` y `/#/editor/words/:id` — Ruta #22, `EDITOR`/`ADMIN`):**
    - Formulario asistido para crear/editar propuestas de lemas (`CreateWordProposal`, `UpdateWordProposal`), agregar/editar acepciones (`CreateSenseProposal`, `UpdateSenseProposal`), añadir ejemplos contextuales (`AddSenseExample`), vincular conceptos semánticos para heterónimos (`ListSemanticConcepts`, `CreateSemanticConcept`, `LinkSenseToConcept`) y subir ilustraciones con `MediaDropzone(category='general')` (`RN12`), siguiendo el **Flujo 6 de [`.context/UI_FLOWS.md`](./UI_FLOWS.md)**.
  - [ ] **`local-app/frontend/src/views/editorial/GovernanceDashboardView.vue` (`/#/editor/governance` — Ruta #23, exclusivo `ADMIN`):**
    - **Tab 1 — Solicitudes de Editor (`CU02`):** Lista `ListPendingEditorRequests`, visor de documento institucional (`credentials_url`) y botones de dictamen (`ResolveEditorRequest`).
    - **Tab 2 — Boletos de Invitación (`CU04`):** Generador de boletos VIP (`GenerateInvitationTickets`) y listado con estado (`ACTIVE`, `REDEEMED`, `EXPIRED`).
    - **Tab 3 — Moderación de Reportes (`CU09`):** Cola de denuncias (`ListCommunityReports`) y resolución (`ResolveCommunityReport`: `DISMISSED`, `RESOLVED` o `CONTENT_ARCHIVED` con archivado en cascada).
    - **Tab 4 — Directorio de Usuarios (`CU03`):** Buscador de cuentas (`ListUsersForGovernance`), cambio de rol (`UpdateUserRole`) y suspensión (`UpdateUserStatus`) respetando protección del último `ADMIN`.
- **Criterio de Cierre:** `wails3 task test` y `wails3 task client:check` 100% en verde con 0 errores; Fase 5 completada en su totalidad.

---

## 🌐 3. Fase 6: Plataforma Web 100% Híbrida (`web-portal/` con Astro SSR SEO + 23 Rutas Interactivas con WebSockets)

> **Objetivo:** Incorporar el workspace `web-portal/` con **Astro + `@astrojs/vue`** consumiendo `@palabrario/ui` (`workspace:*`) y `IPalabrarioAPI` para ofrecer una presencia web pública completa con **100/100 en SEO/Core Web Vitals** y el **100% de las 23 rutas de Palabrario funcionales en el navegador**.

### Arquitectura y Alcance de la Fase 6:
1. **Renderizado en Servidor (SSR) + OpenGraph / JSON-LD para Rutas Públicas (`Rutas #3 a #10`):**
   - `AstroSSRAdapter` consume `core-service (:8000)` en el servidor para entregar HTML semántico instantáneo, meta-tags OpenGraph dinámicas (`og:title`, `og:description`, `og:image`), Twitter Cards, `sitemap.xml` y datos estructurados `DefinedTerm` / `DefinedTermSet` de Schema.org en `/`, `/dictionary`, `/dictionary/[id]`, `/collections`, `/collections/[id]`, `/decks`, `/decks/[id]` y `/u/[username]`.
   - Reutiliza directamente `FriendlyBlurb`, `SenseItem`, `HeteronymPills`, `WordCard`, `DeckCard` y `CollectionCard` desde `@palabrario/ui` con cero duplicación.
2. **Aplicación Interactiva Completa en el Navegador (`100% de las 23 Rutas`):**
   - Integra las áreas autenticadas e interactivas (**Auth `#1–#2`**, **Studio Personal y Repaso SM-2 `#11–#17`**, **Arena Multijugador en Tiempo Real `#18–#20`** y **Portal Editorial / Gobernanza `#21–#23`**) en el navegador web utilizando el adaptador Web de `IPalabrarioAPI` (`fetch` hacia `:8000`/`:8002` y `new WebSocket("ws://.../ws/arena/{roomCode}?token={jwt}")` nativo hacia `realtime-service :8001`).

---

## 🌱 4. Fase 7 (Incremento Post-MVP a Largo Plazo): Sistema de Eventos Culturales de Temporada

> **Objetivo:** Incorporar festivales estacionales (*Día de Muertos*, *Feria de las Flores*, *San Juan*) y galería de arte hispanoamericano según [`proposals/SEASONAL_EVENTS_SPECIFICATION.md`](../proposals/SEASONAL_EVENTS_SPECIFICATION.md), únicamente tras haber cerrado y verificado la Fase 5 y la Fase 6.
