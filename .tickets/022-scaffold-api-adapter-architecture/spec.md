# 022 — Spec checklist

Run by the router before accepting the ticket. All items must pass.

- [x] `app/frontend/src/api/contract/index.ts` exports `IPalabrarioAPI` composing `IAuthAPI`, `IDictionaryAPI`, and `IStudioAPI`.
- [x] Contracts import DTOs directly from `@bindings/palabrario/pkg/models/index.js` (zero duplicate model classes).
- [x] All contract methods return native `Promise<T>` (zero proprietary `$CancellablePromise<T>` in contracts).
- [x] `app/frontend/src/api/mock/fixtures/` contains authentic Spanish fixtures (`words.fixture.ts`, `decks.fixture.ts`, `user.fixture.ts`).
- [x] Strictly NO `flag_emoji` and NO regional codes (`CU`, `CL`, `MX`); all regions render as full names (`Cuba`, `Chile`, `Mexico`, etc.).
- [x] `app/frontend/src/api/mock/mock.adapter.ts` implements `IPalabrarioAPI` with realistic async latency (50–120ms) and session persistence.
- [x] `app/frontend/src/api/live/wails/` implements `IPalabrarioAPI` by delegating to `@bindings/palabrario/app/app.js`.
- [x] `app/frontend/src/api/index.ts` exports singleton `api` instance with runtime mode switching (`auto`, `wails`, `mock`).
- [x] `wails3 task check` in `app/` exits with code 0 (`vue-tsc --noEmit` and Vite build succeed).
- [x] No files modified outside specified scope.
- [x] `.tickets/022-scaffold-api-adapter-architecture/changelog.md` exists and is populated.
