# 🔄 Palabrario — Flujos Interactivos y Diagramas de Secuencia UI $\leftrightarrow$ Puente $\leftrightarrow$ Backend (`.context/UI_FLOWS.md`)

> **Propósito de este documento (SSOT de Comportamiento y Estado Reactivo en Frontend):**  
> Este documento especifica mediante **6 Diagramas de Secuencia y Flujo de Estado** las interacciones más críticas del cliente multiplataforma (`Vue 3 + Reka UI` $\to$ `IPalabrarioAPI` $\to$ `Wails Go BFF / Capacitor` $\to$ `:8000 / :8001 / :8002`).  
> Ningún agente constructor debe improvisar cómo se abren los menús, cómo se hidrata la sesión o cómo fluyen las tramas WebSocket: debe seguir estrictamente estos 6 flujos canónicos.

---

## Flujo 1: Autenticación e Hidratación de Sesión Dual (`Wails RAM` vs `Capacitor Keystore` — `CU01`)

Al iniciar la aplicación o al autenticarse en `/#/auth/sign-in` y `/#/auth/sign-up`, el composable global `useAuth()` coordina la sesión sin exponer jamás el token JWT a `localStorage` en escritorio:

```mermaid
sequenceDiagram
    autonumber
    actor U as Usuario
    participant V as SignInView / App.vue
    participant Auth as useAuth() Store
    participant API as IPalabrarioAPI (Adapter)
    participant BFF as Wails app.go (RAM) / Capacitor Keystore
    participant Core as core-service (:8000)

    Note over V,Core: A. Inicio de Sesión (SignIn / SignUp)
    U->>V: Ingresa credenciales y envía formulario
    V->>Auth: signIn({ email, password })
    Auth->>API: SignIn(SignInRequest)
    API->>BFF: Invoca método de transporte
    BFF->>Core: POST /api/v1/auth/sign-in
    Core-->>BFF: 200 OK -> AuthResponse { token, user: UserProfile }
    BFF->>BFF: Custodia JWT en memoria privada Go (sync.RWMutex)<br/>o en @capacitor/preferences (Móvil)
    BFF-->>Auth: Retorna UserProfile (sin exponer JWT al DOM en Desktop)
    Auth->>Auth: Actualiza currentUser.value = user
    Auth-->>V: Redirige a /#/ (o ruta destino protegida)

    Note over V,Core: B. Hidratación Reactiva y Cierre de Sesión
    V->>Auth: hydrateSession() (al montar AppLayout / TopBar)
    Auth->>API: GetMe()
    API->>BFF: Adjunta Bearer JWT si existe en bóveda
    alt Sesión válida
        BFF->>Core: GET /api/v1/auth/me
        Core-->>Auth: 200 OK -> UserProfile actualizado (rol, XP, racha)
    else Sin sesión o JWT expirado (401)
        BFF-->>Auth: Error 401 -> Limpia estado local (currentUser = null)
    end
```

---

## Flujo 2: Menú de Usuario en `TopBar` (`Reka UI DropdownMenu`) + `Popover` de *"Credenciales y Rol en la Imprenta"* (`CU02, CU04, RN02, RN12`)

> **Decisión de Diseño UX:**  
> En lugar de esconder la promoción editorial en el fondo de una página, al hacer clic sobre el avatar/indicador del usuario en el extremo superior derecho de **`TopBar.vue`**, se despliega un **`DropdownMenu` de Reka UI** con:
> 1. **Mi Perfil Público** (`/#/u/:username`)
> 2. **Configuración de Cuenta** (abre modal/drawer de ajustes de perfil y contraseña `CU01`)
> 3. **Credenciales y Rol en la Imprenta** $\to$ Al seleccionarlo, despliega un **`Popover` / `Dialog` interactivo de Reka UI (`CredentialsRolePopover.vue`)** desde la barra superior que gestiona tanto el **Canje de Boleto (`CU04`)** como la **Solicitud de Ascenso a Editor (`CU02`)**.
> 4. **Panel Editorial (`/#/editor`)** / **Gobernanza Admin (`/#/admin/governance`)** *(visibles dinámicamente según `user.role`)*.
> 5. **Cerrar Sesión (`SignOut`)**.

```mermaid
sequenceDiagram
    autonumber
    actor U as Usuario Autenticado
    participant TB as TopBar.vue (Reka DropdownMenu)
    participant Pop as CredentialsRolePopover.vue (Reka Popover)
    participant API as IPalabrarioAPI / useAuth
    participant Img as image-service (:8002)
    participant Core as core-service (:8000)

    U->>TB: Clic en Avatar / Barra de Estado de Usuario
    TB-->>U: Despliega Reka UI DropdownMenu (Perfil, Config, Credenciales y Rol, Cerrar Sesión)
    U->>TB: Clic en ítem "Credenciales y Rol en la Imprenta"
    TB->>Pop: Abre Popover anclado a TopBar y consulta estado actual
    Pop->>API: ListMyEditorRequests()
    API->>Core: GET /api/v1/governance/editor-requests/me
    Core-->>Pop: EditorRequest[] (Verifica si existe 1 PENDING según RN02)

    alt Caso A: El Aprendiz tiene un Boleto de Invitación (CU04 — Ascenso Inmediato)
        U->>Pop: Pestaña "Canjear Boleto" -> Ingresa código (ej. "PAL-UNAM-2026")
        Pop->>API: RedeemInvitationTicket({ code })
        API->>Core: POST /api/v1/governance/tickets/redeem (Bloqueo FOR UPDATE)
        Core-->>Pop: 200 OK (Rol promovido a EDITOR)
        Pop->>API: GetMe() (Refresca useAuth().currentUser)
        API->>Core: GET /api/v1/auth/me
        Core-->>TB: UserProfile { role: 'EDITOR' }
        TB-->>U: Badge en TopBar cambia a EDITOR y habilita acceso a "/#/editor"
    else Caso B: Solicitud Formal de Ascenso a Editor con Comprobante (CU02 + RN02 + RN12)
        alt Ya tiene una solicitud en estado PENDING (RN02)
            Pop-->>U: Muestra tarjeta de estado "Solicitud en Revisión Institucional" (Formulario bloqueado)
        else Sin solicitud pendiente
            U->>Pop: Pestaña "Solicitar Acreditación" -> Arrastra comprobante en MediaDropzone
            Pop->>API: UploadImage(file, 'credentials')
            API->>Img: POST :8002/api/v1/images/upload (Multipart + Bearer JWT)
            Img-->>Pop: 201 Created -> ImageUploadResponse { url: "/uploads/credentials/..." }
            U->>Pop: Completa institution_name + justification (>= 20 chars) y envía
            Pop->>API: CreateEditorRequest({ institution_name, justification, credentials_url })
            API->>Core: POST /api/v1/governance/editor-requests
            Core-->>Pop: 201 Created -> EditorRequest (status = 'PENDING')
            Pop-->>U: Transiciona en vivo a vista "Solicitud Enviada (PENDING)"
        end
    end
```

---

## Flujo 3: Exploración Léxica, Ficha Estilo Vocabulary.com (`FriendlyBlurb` + `SenseItem`), Guardado en Baraja y Reporte (`CU05, CU09, CU10`)

Coordina el buscador omnibox (`Ctrl + K`), la lectura de una palabra con advertencia de tabú y las acciones rápidas sobre cada acepción:

```mermaid
sequenceDiagram
    autonumber
    actor U as Usuario
    participant Omni as SearchModal.vue (Ctrl+K)
    participant WD as WordDetailView.vue
    participant DeckM as AddToDeckModal.vue
    participant RepM as ReportModal.vue
    participant API as IPalabrarioAPI
    participant Core as core-service (:8000)

    U->>Omni: Presiona Ctrl+K y escribe término dialectal
    Omni->>API: SearchWords(q, region, 10)
    API->>Core: GET /api/v1/words/search?q=...
    Core-->>Omni: WordSearchResult[] (FTS + Trigramas)
    U->>Omni: Selecciona resultado -> Navega a /#/dictionary/word/:id

    WD->>API: GetWordByID(id)
    API->>Core: GET /api/v1/words/{id}
    Core-->>WD: WordDetail (friendly_blurb, etymology, senses[], heteronyms[])
    WD-->>U: Renderiza FriendlyBlurb (2 párrafos conversacionales) + Lista numerada SenseItem

    opt Guardar Acepción en Baraja Personal (CU10)
        U->>WD: Clic en "+ Guardar en mi baraja" dentro de SenseItem (id_sense)
        WD->>DeckM: Abre modal pasando id_sense
        DeckM->>API: ListStudioDecks()
        API->>Core: GET /api/v1/studio/decks
        Core-->>DeckM: StudioDeckSummary[] (Barajas propias del usuario)
        U->>DeckM: Selecciona baraja destino (id_deck)
        DeckM->>API: AddDeckItem(id_deck, { id_sense })
        API->>Core: POST /api/v1/decks/{id}/items
        Core-->>DeckM: 201 Created -> Toast de confirmación
    end

    opt Reportar Contenido Inexacto u Ofensivo (CU09)
        U->>WD: Clic en "Reportar" sobre lema, sentido o ejemplo
        WD->>RepM: Abre ReportModal con TargetType ('WORD'|'SENSE'|'EXAMPLE') e ID
        U->>RepM: Selecciona Reason y escribe Details (>= 10 chars)
        RepM->>API: CreateReport(CreateReportRequest)
        API->>Core: POST /api/v1/reports
        Core-->>RepM: 201 Created -> Notificación de reporte recibido
    end
```

---

## Flujo 4: Sesión de Estudio Inmersiva con Motor Matemático SuperMemo SM-2 y Ledger de XP (`CU12, CU13, RN09, RN11`)

Ejecuta el ciclo de tarjetas en `/#/studio/study` (`layout: immersive`), midiendo el tiempo de respuesta en cliente y consolidando cada repaso en una transacción ACID en el backend:

```mermaid
sequenceDiagram
    autonumber
    actor U as Aprendiz
    participant SV as StudySessionView.vue
    participant Card as StudyCard.vue
    participant API as IPalabrarioAPI
    participant Core as core-service (:8000)

    SV->>API: GetStudySession(deckId?, 20)
    API->>Core: GET /api/v1/studio/study-session
    Core-->>SV: StudySessionPayload { total_due, session_cards: StudyCardItem[] }

    loop Por cada StudyCardItem en session_cards
        SV->>Card: Muestra Anverso (Lema + Bandera + Categoría) e inicia cronómetro local (t0)
        U->>Card: Pulsa "Revelar Respuesta" (o barra espaciadora)
        Card-->>U: Muestra Reverso (Definición, Nota Cultural, Ejemplos) y 4 botones SM-2
        U->>Card: Califica recuerdo: 1 (Otra vez), 2 (Difícil), 3 (Bien) o 4 (Fácil)
        Card->>SV: Emite { id_sense, quality_rating, time_spent_ms: now - t0 }
        SV->>API: SubmitStudyReview(SubmitReviewRequest)
        API->>Core: POST /api/v1/studio/study-session/review (Tx ACID: SM-2 + xp_transactions + racha)
        Core-->>SV: CardReviewResult { updated_mastery, xp_awarded, new_total_xp, newly_unlocked_badges }
        SV-->>U: Animación "+XP", actualiza barra de progreso y muestra medalla si newly_unlocked_badges != empty
    end

    SV-->>U: Muestra Resumen de Sesión Completada y botón "Volver al Studio" (/#/studio)
```

---

## Flujo 5: Arena Multijugador en Tiempo Real (`WebSockets`), Modalidad `WORD_TO_IMAGE` y Replay 3NF (`CU14, CU15, CU16, RN10, RN13`)

Coordina la creación de la sala en `core-service` (`:8000`), la partida síncrona en `realtime-service` (`:8001`) a través de `ArenaEventBus` y el Replay final en 3NF:

```mermaid
sequenceDiagram
    autonumber
    actor Host as Anfitrión / Jugador
    participant Lobby as ArenaLobbyView.vue
    participant Game as ArenaGameView.vue
    participant Bus as IPalabrarioAPI + ArenaEventBus
    participant Core as core-service (:8000)
    participant RT as realtime-service (:8001)

    Note over Host,Core: 1. Preparación de Sala (CU14 — Validación de >= 8 sentidos y RN13)
    Host->>Lobby: Selecciona baraja, rondas (8..20) y tiempo (10..30s) -> "Crear Sala"
    Lobby->>Bus: PrepareMatch(CreateMatchRequest)
    Bus->>Core: POST /api/v1/arena/matches/prepare
    Core-->>Lobby: 201 Created -> MatchPackage { id_match, room_code: "MX8921" }
    Lobby->>Game: Redirige a /#/game/MX8921

    Note over Host,RT: 2. Conexión WebSocket y Ciclo de Rondas en Vivo (CU15 + RN10)
    Game->>Bus: ConnectArena("MX8921") + onArenaEvent(handler)
    Bus->>RT: WS /ws/arena/MX8921?token={jwt}
    RT-->>Bus: WSEvent { type: "ROOM_STATE", payload: ArenaRoomState }
    Host->>Game: Anfitrión pulsa "Iniciar Contienda"
    Game->>Bus: StartArenaMatch() -> WS { type: "START_MATCH" }

    loop Por cada Ronda (1..total_rounds)
        RT-->>Bus: WSEvent { type: "COUNTDOWN", payload: { seconds: 3 } }
        RT-->>Bus: WSEvent { type: "QUESTION", payload: ArenaQuestion (sin CorrectOption) }
        alt modality == 'WORD_TO_IMAGE' (RN13)
            Game-->>Host: Renderiza ImageChoiceGrid (2x2 imágenes) + ArenaTimer (tabular-nums)
        else Modalidad Textual ('WORD_TO_DEFINITION' / 'DEFINITION_TO_WORD')
            Game-->>Host: Renderiza 4 tarjetas de texto + ArenaTimer
        end
        Host->>Game: Selecciona opción (clic o tecla 1..4)
        Game->>Bus: SubmitArenaAnswer({ round_number, selected_option, client_elapsed_ms })
        Bus->>RT: WS { type: "SUBMIT_ANSWER", payload }
        RT-->>Bus: WSEvent { type: "ROUND_RESULT", payload: RoundResultPayload (revela correct_option + leaderboard) }
    end

    Note over Host,RT: 3. Consolidación Batch Flush S2S y Replay 3NF (CU16)
    RT->>Core: POST /api/v1/arena/matches/consolidate (S2S JWT sha256 sub-key)
    Core-->>RT: 200 OK -> MatchSummary
    RT-->>Bus: WSEvent { type: "PODIUM", payload: MatchSummary }
    Game-->>Host: Renderiza PodiumCard y botón "Ver Replay Detallado Ronda por Ronda"
    Host->>Game: Clic en "Ver Replay" -> Navega a /#/game/:matchId/replay
    Game->>Core: GET /api/v1/arena/matches/{matchId} -> MatchReplayDetail (3NF)
```

---

## Flujo 6: Autoría Léxica con Multimedia (`RN12`), Curaduría por Quórum Tripartito (`RN06, RN07`) y Gobernanza Admin (`CU03, CU06–CU09, RN08`)

Describe cómo los Editores proponen contenido, cómo se vota hasta alcanzar los **3 sellos de aprobación** (triggers `RN06`/`RN07`) y cómo el Administrador ejerce voto directo o resuelve reportes en cascada:

```mermaid
sequenceDiagram
    autonumber
    actor Ed as Editor Acreditado
    actor Adm as Administrador
    participant Form as EditorWordFormView.vue
    participant Dash as EditorDashboardView.vue
    participant Gov as GovernanceDashboardView.vue
    participant Img as image-service (:8002)
    participant Core as core-service (:8000)

    Note over Ed,Core: 1. Propuesta de Palabra o Acepción con Imagen Validada (CU06 + RN12)
    Ed->>Form: Arrastra ilustración en MediaDropzone (category='general')
    Form->>Img: POST :8002/api/v1/images/upload (Limpieza EXIF + Filtro NSFW + Resize 800px)
    Img-->>Form: 201 Created -> { url: "/uploads/general/uuid.webp" }
    Ed->>Form: Envía CreateWordProposalRequest (o CreateSenseRequest)
    Form->>Core: POST /api/v1/words (El sentido nace en state = 'DRAFT')

    Note over Ed,Core: 2. Votación por Quórum Tripartito (CU07, CU08 + Triggers RN06 y RN07)
    Ed->>Dash: Abre /#/editor -> ListPendingQueue()
    Dash->>Core: GET /api/v1/curation/queue
    Core-->>Dash: SenseQuorumItem[] (Muestra QuorumStepper: approvals_count / 3)
    Ed->>Dash: Inspecciona propuesta ajena (author_id != currentUser.id_user por RN07) y vota
    Dash->>Core: POST /api/v1/curation/senses/{senseId}/vote (SubmitVoteRequest)
    Core-->>Dash: 200 OK (Al llegar al 3er voto APPROVE, el trigger RN06 publica el sentido a APPROVED)

    Note over Adm,Core: 3. Veto/Publicación Directa (RN08) y Moderación en Cascada (CU09)
    Adm->>Dash: Ejecuta AdminOverrideState(senseId, { state: 'APPROVED' | 'HIDDEN', reason })
    Dash->>Core: PUT /api/v1/curation/senses/{senseId}/state
    Adm->>Gov: Abre /#/admin/governance -> Resuelve reporte comunitario (ResolveReportRequest)
    Gov->>Core: PUT /api/v1/reports/{reportId}/resolve (Oculta recurso en cascada si RESOLVED_HIDDEN)
```
