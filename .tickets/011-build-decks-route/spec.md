# 011 — Spec checklist

Run by the router before accepting the ticket. All items must pass.

- [ ] `/decks` route is accessible under `AppLayout` and renders without errors.
- [ ] `DecksView.vue` consumes decks from `api.studio.listDecks(...)` via `import { api } from "@/api"`. Zero hardcoded mock decks exist inside the view or subcomponents.
- [ ] `IStudioAPI` contract and adapters (Wails and Mock) properly support `listDecks(q?: string, officialOnly?: boolean)`.
- [ ] Mobile-first responsive layout implemented.
- [ ] Monochrome rules strictly observed: no `rounded-*`, no drop shadows, sharp solid borders.
- [ ] `cd app && wails3 task check` exits with code 0.
- [ ] Changes stay within declared scope.
- [ ] `changelog.md` exists and is populated.
- [ ] Commit (if made): Conventional Commits, English, detailed body.
