# 001 — Spec checklist

Run by the router before accepting the ticket. All items must pass.

- [x] `.routing/DOMAIN.md` exists and is not empty.
- [x] `.routing/INFRASTRUCTURE.md` exists and is not empty.
- [x] `.routing/DOMAIN.md` has no heading or content from section 3.6.
- [x] Line count of `.routing/DOMAIN.md` is within 100 of `.context/DOMAIN.md` minus section 3.6.
- [x] Line count of `.routing/INFRASTRUCTURE.md` is within 100 of `.context/INFRASTRUCTURE.md`.
- [x] Heading count (`#` lines) matches each reference (DOMAIN excludes 3.6).
- [x] Fenced code block delimiter count matches each reference.
- [x] Markdown link count (`](`) matches each reference.
- [x] Output is English. Spanish remains only in identifiers, code, and glossary terms.
- [x] Changes are limited to `.routing/DOMAIN.md`, `.routing/INFRASTRUCTURE.md`, and this ticket folder.
- [x] `changelog.md` exists in this ticket folder and records the line counts used.
- [x] Commit (if made): docs-only, Conventional Commits, English, with a detailed body.
