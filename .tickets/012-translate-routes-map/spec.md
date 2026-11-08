# 012 — Spec checklist

Run by the router before accepting the ticket. All items must pass.

- [x] `.routing/ROUTES.md` exists and is completely in English.
- [x] The Master Route Table contains all 9 routes listed in the task with exact columns.
- [x] Layouts section documents `AppLayout` and `AuthLayout`.
- [x] Guards section accurately notes that no guards are active in code yet.
- [x] Ownership rules specify workers do not edit the file directly.
- [x] `.routing/CONTEXT.md` contains the link `- [ROUTES.md](ROUTES.md): route table, layouts, and access rules. Read only when a ticket touches a route.` under `## Files`.
- [x] No files outside `.routing/ROUTES.md`, `.routing/CONTEXT.md`, and `.tickets/012-translate-routes-map/` were modified.
- [x] `.tickets/012-translate-routes-map/changelog.md` exists and is populated.
- [x] Commit (if made): docs-only, Conventional Commits, English, with a detailed body.

