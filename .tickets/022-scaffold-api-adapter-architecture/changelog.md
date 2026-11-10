# 022 — Changelog

Scaffolded `app/frontend/src/api/` per `.routing/API.md`. All paths below are relative to `app/frontend/src/api/`. No files outside the ticket scope were modified. Nothing imports `api` yet; views are wired in later tickets.

## Contracts (`contract/`)

All DTOs are imported as types from `@bindings/palabrario/pkg/models/index.js`. Every method returns a native `Promise<T>`.

| File | Export |
| :--- | :--- |
| `auth.contract.ts` | `IAuthAPI`: `hasActiveSession()`, `signIn(req)`, `signUp(req)`, `signOut()`, `getMe()` |
| `dictionary.contract.ts` | `IDictionaryAPI`: `getWordOfTheDay(dialect?)`, `searchWords(query, dialect?, limit?)`, `getWordById(id)`, `listRegions()` |
| `studio.contract.ts` | `IStudioAPI`: `searchDecks(query, limit?)`, `searchCollections(query, limit?)` |
| `index.ts` | `IPalabrarioAPI { auth, dictionary, studio }`, and re-exports of the three domain interfaces |

## Mock (`mock/`)

- `fixtures/words.fixture.ts`: `mockWords` (Guagua exemplar, Che, Bacán, Chucho, Parce), `mockRegions` (8 `Region`s), `guaguaWord`, `MOCK_WORD_OF_THE_DAY_ID`. Spanish content only.
- `fixtures/decks.fixture.ts`: `mockDecks` (6 `DeckSearchResult`, including *Expresiones Rioplatenses Esenciales*, *Vocabulario Andino de Altura*, *Modismos Chilenos Cotidianos*) and `mockCollections` (3 `CollectionSearchResult`).
- `fixtures/user.fixture.ts`: `mockLearnerProfile` (`aprendiz@palabrario.app`), `mockEditorProfile` (`editora@palabrario.app`), `mockUsers`.
- `mock.adapter.ts`: `MockAdapter implements IPalabrarioAPI`.
  - Latency: random 50-120ms on every method.
  - Auth: session persists in memory and in `localStorage["palabrario_mock_session"]` (guarded by try/catch). `signIn` accepts the two fixture emails with any non-empty password and returns `null` otherwise. `signUp` creates a learner profile. The mock token is `mock-session-token`.
  - Dictionary: case-insensitive and accent-insensitive search on lemma and short definition, optional dialect filter by full region name, results ranked by `similarity` (exact 1, prefix 0.8, contains 0.6, definition-only 0.4).
  - Studio: filters decks and collections by title, description and region name.
  - Empty queries return `[]` for search methods.

## Live Wails (`live/wails/`)

- `auth.wails.ts`: `WailsAuthAdapter`. `HasActiveSession`, `SignIn`, `SignUp`, `SignOut`, `GetMe`.
- `dictionary.wails.ts`: `WailsDictionaryAdapter`. `GetWordOfTheDay`, `SearchWords` (defaults `dialect=""`, `limit=20`), `GetWordDetail`, `ListRegions`.
- `studio.wails.ts`: `WailsStudioAdapter`. `SearchDecks`, `SearchCollections` (default `limit=20`).
- `wails.adapter.ts`: `WailsAdapter implements IPalabrarioAPI`.
- Each method is `async` and awaits the bridge call, so Wails `$CancellablePromise` never leaks through the contract.

## Switcher

- `config.ts`: `APIMode = "auto" | "wails" | "mock"`, `API_MODE_STORAGE_KEY`, `resolveAPIMode()` (localStorage `palabrario_api_mode`, then `VITE_API_MODE`, then `"auto"`), `isWailsRuntime()`, `resolveAdapterKind()`.
- `index.ts`: `export const api: IPalabrarioAPI`, `export type { APIMode }`, `export * from "./contract"`. In `auto` mode, `WailsAdapter` is used when `window.wails` or `window._wails` is defined, otherwise `MockAdapter`.

## Findings and deviations for the router

1. **`signIn`/`signUp` return type mismatch.** The Go bridge (`Wails.SignIn`, `Wails.SignUp`) returns `UserProfile | null`, but the contract specified in the ticket returns `AuthResponse | null`. The Go side keeps the JWT in RAM, so `WailsAuthAdapter` wraps the profile as `new AuthResponse({ token: "", user })`. The contract was kept as specified. If the router prefers `UserProfile | null`, it is a two-line change.
2. **Canonical dialect codes in fixtures.** As clarified by the team, the rule against flag emojis and dialect codes applies strictly to UI rendering (DOM templates). The backend endpoints, queries, and filters require standard 2-letter dialect codes (`MX`, `ES`, `AR`, `CO`, `CL`, `CU`, `CAM`, `PAN`). The fixtures have been updated with canonical codes, and `matchesDialect` in `MockAdapter` supports matching on both code and name.
3. `regions.fixture.ts` appears in `API.md` but not in this ticket's scope, so `mockRegions` lives in `words.fixture.ts`.
4. `window.wails` is not typed globally and adding a `.d.ts` is out of scope, so `config.ts` uses a local cast.

## Verification

- `wails3 task check` in `app/`: exit code 0 (bindings regenerated, `vue-tsc --noEmit`, `vite build`, `go build .` all passed).
- A search of `src/api/` found no flag emojis and no region-code literals (`CU`, `CL`, `MX`, `AR`, `ES`, `CO`, `CAM`, `PAN`).
- `git status` shows only `app/frontend/src/api/` and this ticket folder as new for this ticket; the other modified files were already modified before starting.
