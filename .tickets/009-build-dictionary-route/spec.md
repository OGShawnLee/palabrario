# 009 — Spec checklist

Run by the router before accepting the ticket. All items must pass.

- [ ] `/dictionary` renders under `AppLayout` without console errors.
- [ ] Mobile-first styling: base styles target mobile screens, with responsive scaling.
- [ ] Monochrome rules strictly observed: no `rounded-*` (except `rounded-none`), no drop shadows, no `hover:` or `active:` styles unless explicitly defined.
- [ ] `cd app && wails3 task check` exits with code 0.
- [ ] No files modified outside specified scope.
- [ ] `changelog.md` exists and is populated.
- [ ] Commit (if made): Conventional Commits, English, detailed body.
