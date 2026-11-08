# 015 — Frontend: Integrate SearchModal with Live Wails Backend Bindings

| Field | Value |
| :--- | :--- |
| **ID** | 015 |
| **Status** | draft |
| **Model** | Sonnet 5.5 (`claude-sonnet-5-5`, Antigravity provider) |
| **Role** | Frontend Implementer |
| **Spec** | [spec.md](spec.md) |

> **Draft Note:**  
> This ticket connects `SearchModal.vue` to the official Wails Go bindings produced by Ticket 014, replacing all mock data and local mock interfaces with canonical `@bindings/palabrario/pkg/models` types.

## Scope

May create or edit:
- `app/frontend/src/components/SearchModal.vue`
- `.tickets/015-integrate-search-modal/changelog.md` (create)

Must not touch:
- `app/frontend/uno.config.ts`
- Backend code, Go services, or `.context/`
- Any file outside the above scope

## Inputs

- [.tickets/014-backend-search-fts/task.md](../014-backend-search-fts/task.md): Backend contracts and bindings.
- [app/frontend/bindings/palabrario/pkg/models/models.ts](../../app/frontend/bindings/palabrario/pkg/models/models.ts): Canonical models.
- [app/frontend/src/components/SearchModal.vue](../../app/frontend/src/components/SearchModal.vue): Existing UI component.

## Task

1. **Remove Local Mock Interfaces:**
   - Remove `MockWordResult` and `MockStudyDeck` from `SearchModal.vue`.
   - Import canonical models from `@bindings/palabrario/pkg/models` (`WordSearchResult`, `DeckSearchResult`, etc.).
2. **Wire Reactive Search with Debounce:**
   - Use `@vueuse/core` (`useDebounceFn` or `refDebounced` with 200–250ms) to trigger backend search calls as the user types in the input.
   - When query is empty, display curated/default suggestions or clear results cleanly.
   - Call `api.SearchWords(query, userDialect, limit)` and `api.SearchDecks(query, limit)`.
3. **Preserve User Interactions & Ground Rules:**
   - Keep dynamic synchronization between highlighted item and the left Featured Entry card.
   - Ensure zero flag emojis and zero dialect codes are displayed (use `region_name` provided by backend).
   - Keep keyboard navigation (`ArrowUp`, `ArrowDown`, `Enter`, `Escape`) and tab filters (`All`, `Dictionary`, `Study Decks`).
4. **Validation:**
   - Run `cd app && wails3 task check`. Must exit with code 0.
5. **Change log:**
   - Document all changes in `.tickets/015-integrate-search-modal/changelog.md`.

## Verification

All checks are listed in [spec.md](spec.md).
