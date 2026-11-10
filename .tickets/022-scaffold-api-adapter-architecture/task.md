# 022 — Scaffold API Adapter Architecture and Core Domains (Auth, Dictionary, Studio)

| Field | Value |
| :--- | :--- |
| **ID** | 022 |
| **Status** | completed |
| **Model** | Claude Sonnet 5.5 (`claude-sonnet-5-5`, Antigravity provider) |
| **Role** | Frontend Architect & Implementer |
| **Spec** | [spec.md](spec.md) |

> **Architectural Note:**  
> This ticket scaffolds the new `app/frontend/src/api/` layer according to `.routing/API.md`. It establishes pure TypeScript domain contracts returning native `Promise<T>`, builds the static `MockAdapter` with rich Spanish fixtures for offline instant exploration, and implements the live `WailsAdapter` delegating to `@bindings/palabrario/app/app`.

## Scope

May create or edit:
- `app/frontend/src/api/contract/index.ts`
- `app/frontend/src/api/contract/auth.contract.ts`
- `app/frontend/src/api/contract/dictionary.contract.ts`
- `app/frontend/src/api/contract/studio.contract.ts`
- `app/frontend/src/api/mock/mock.adapter.ts`
- `app/frontend/src/api/mock/fixtures/words.fixture.ts`
- `app/frontend/src/api/mock/fixtures/decks.fixture.ts`
- `app/frontend/src/api/mock/fixtures/user.fixture.ts`
- `app/frontend/src/api/live/wails/wails.adapter.ts`
- `app/frontend/src/api/live/wails/auth.wails.ts`
- `app/frontend/src/api/live/wails/dictionary.wails.ts`
- `app/frontend/src/api/live/wails/studio.wails.ts`
- `app/frontend/src/api/index.ts`
- `app/frontend/src/api/config.ts`
- `.tickets/022-scaffold-api-adapter-architecture/changelog.md` (create)

Must not touch:
- `app/frontend/uno.config.ts`
- Backend Go code (`app/*.go`, `core-service/`, etc.)
- Existing route views (views will be wired in their respective route/integration tickets)
- Any file outside the above scope

## Inputs

- [.routing/API.md](../../.routing/API.md): Canonical API architecture, 3 pillars of synchronization, directory structure, mock rules.
- [.routing/CONTEXT.md](../../.routing/CONTEXT.md): Ground rules (strictly NO `flag_emoji`, NO regional codes like `CU`, `CL`, `MX`; full region names only).
- [app/frontend/bindings/palabrario/app/app.ts](../../app/frontend/bindings/palabrario/app/app.ts): Live Wails bridge methods (`SignIn`, `SignUp`, `GetMe`, `HasActiveSession`, `SignOut`, `SearchWords`, `GetWordDetail`, `SearchDecks`, `SearchCollections`).
- [app/frontend/bindings/palabrario/pkg/models/index.ts](../../app/frontend/bindings/palabrario/pkg/models/index.ts): Canonical models (`UserProfile`, `WordDetail`, `WordSearchResult`, `DeckSearchResult`, `CollectionSearchResult`, etc.).
- [app/frontend/src/data/mockWordOfTheDay.ts](../../app/frontend/src/data/mockWordOfTheDay.ts): Existing Guagua exemplar and heteronym fixture data to incorporate into `words.fixture.ts`.
- [app/frontend/src/data/regions.ts](../../app/frontend/src/data/regions.ts): Canonical region definitions.

## Task

1. **Define Core TypeScript Contracts (`src/api/contract/`):**
   - Import all DTO types directly from `@bindings/palabrario/pkg/models/index.js` (never duplicate DTO classes).
   - Create `auth.contract.ts`:
     ```ts
     export interface IAuthAPI {
       hasActiveSession(): Promise<boolean>
       signIn(req: SignInRequest): Promise<AuthResponse | null>
       signUp(req: SignUpRequest): Promise<AuthResponse | null>
       signOut(): Promise<void>
       getMe(): Promise<UserProfile | null>
     }
     ```
   - Create `dictionary.contract.ts`:
     ```ts
     export interface IDictionaryAPI {
       getWordOfTheDay(dialect?: string): Promise<WordDetail | null>
       searchWords(query: string, dialect?: string, limit?: number): Promise<WordSearchResult[]>
       getWordById(id: number): Promise<WordDetail | null>
       listRegions(): Promise<Region[]>
     }
     ```
   - Create `studio.contract.ts`:
     ```ts
     export interface IStudioAPI {
       searchDecks(query: string, limit?: number): Promise<DeckSearchResult[]>
       searchCollections(query: string, limit?: number): Promise<CollectionSearchResult[]>
     }
     ```
   - Create `index.ts`:
     ```ts
     export interface IPalabrarioAPI {
       readonly auth: IAuthAPI
       readonly dictionary: IDictionaryAPI
       readonly studio: IStudioAPI
     }
     ```
   - All contract methods must return native `Promise<T>` (never Wails `$CancellablePromise<T>`).

2. **Build Spanish Mock Fixtures (`src/api/mock/fixtures/`):**
   - `words.fixture.ts`: Incorporate the canonical *Guagua* exemplar from `mockWordOfTheDay.ts` (pure Spanish text for senses, definitions, etymology, and blurb) plus searchable word entries (`che`, `bacán`, `chucho`, `parce`, `guagua`). Strictly NO flag emojis and NO dialect codes (full region names only).
   - `decks.fixture.ts`: Realistic study decks and curated collections (e.g. *Expresiones Rioplatenses Esenciales*, *Vocabulario Andino de Altura*, *Modismos Chilenos Cotidianos*) conforming to `DeckSearchResult` and `CollectionSearchResult`.
   - `user.fixture.ts`: Mock profiles for learner and editor roles.

3. **Implement `MockAdapter` (`src/api/mock/mock.adapter.ts`):**
   - Implement `IPalabrarioAPI`:
     - `auth`: maintains simple in-memory / localStorage state (`palabrario_mock_session`) so signing in/out persists during the session.
     - `dictionary`: returns `words.fixture.ts` data, with case-insensitive search matching against lemma and short definition.
     - `studio`: returns filtered decks and collections matching query.
     - Realistic async delay (50–120ms) using `new Promise(resolve => setTimeout(resolve, delay))`.

4. **Implement Live `WailsAdapter` (`src/api/live/wails/`):**
   - Implement `IAuthAPI` (`auth.wails.ts`), `IDictionaryAPI` (`dictionary.wails.ts`), and `IStudioAPI` (`studio.wails.ts`).
   - Delegates directly to `@bindings/palabrario/app/app.js` bridge functions.
   - Compose into `WailsAdapter` in `wails.adapter.ts` implementing `IPalabrarioAPI`.

5. **Implement Active Export & Switcher (`src/api/index.ts` and `src/api/config.ts`):**
   - Define `export type APIMode = "auto" | "wails" | "mock"`.
   - Resolution order:
     1. `localStorage.getItem("palabrario_api_mode")` if set to `"wails"` or `"mock"`.
     2. `import.meta.env.VITE_API_MODE` if defined.
     3. Auto-detect: if `typeof window !== "undefined" && window.wails !== undefined`, use `WailsAdapter`; otherwise fallback to `MockAdapter`.
   - Export singleton:
     ```ts
     export const api: IPalabrarioAPI = createAPIInstance()
     export * from "./contract"
     ```

6. **Validation:**
   - Run `wails3 task check` from `app/`. It must exit with code 0 (`vue-tsc --noEmit` and production build must pass).

7. **Change log:**
   - Document all created files, contracts, and export signatures in `.tickets/022-scaffold-api-adapter-architecture/changelog.md`.

## Verification

All checks are listed in [spec.md](spec.md).
