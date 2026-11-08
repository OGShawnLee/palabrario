# Changelog — Ticket 002: Translate INTERFACES.md

## 1. Line Counts and Spec Verification Metrics

| Metric | Reference (`.context/INTERFACES.md`) | Translated (`.routing/INTERFACES.md`) | Difference / Status |
| :--- | :---: | :---: | :--- |
| **Total lines** | 1050 | 1050 | Diff = 0 lines (Within allowable limit of $\pm 100$) |
| **Headings (`#`)** | 39 | 39 | Exact match (all levels and hierarchy preserved) |
| **Code fences (`````)** | 24 | 24 | Exact match |
| **Markdown links (`](`)** | 13 | 13 | Exact match (all link targets identical) |

## 2. Terms and Translation Decisions
- **Bóveda de sesión en RAM** $\to$ Translated as `Session vault in RAM`.
- **Dictamen filológico / administrativo** $\to$ Translated as `Philological ruling` / `Administrative ruling`.
- **Boleto de invitación / Boletos** $\to$ Translated as `Invitation ticket` / `Tickets`.
- **Reactivo / Reactivo interno** $\to$ Translated as `Prompt` / `Internal prompt`.
- **Distractores** $\to$ Translated as `Distractors`.
- **Retrospección y replay granular ronda por ronda** $\to$ Translated as `Retrospection and granular round-by-round replay`.
- **Desglose ronda por ronda** $\to$ Translated as `Round-by-round breakdown`.
- **Asientos en xp_transactions** $\to$ Translated as `Entries in xp_transactions`.
- **Unchanged Terms:** Maintained `RBAC`, `FTS`, `SM-2`, `CEFR`, `XP`, `JWT`, `S2S`, `WebSocket`, `PostgreSQL`, `sqlc`, `pgxpool`, `Chi`, `Wails`, `Vue`, `UnoCSS`, `Docker`, `Vagrant`, and roles `LEARNER`, `EDITOR`, `ADMIN` strictly as instructed.

## 3. Untranslated Content
- None. All Spanish prose, diagram annotations, directory tree comments, table descriptions, and code comments were translated into English. Spanish remains only in code identifiers, Go struct tags, endpoints, database entities, and backend-emitted localized error payloads (e.g., `{"error": "..."}`).

## 4. Proposed Glossary Entries
- **Invitation Ticket (`Boleto de invitación`):** Institutional token with single/multi-use limits used to redeem access or grant elevated roles (`EDITOR`).
- **Internal Prompt (`Reactivo interno`):** Complete server-side question representation including the correct answer and distractors for authoritative grading.
- **Batch Flush (`Descarga por lotes`):** End-of-match state consolidation payload sent from `realtime-service` to `core-service`.
- **Session Vault (`Bóveda de sesión`):** In-memory Go RAM store inside the Wails desktop process isolating user JWT tokens from the WebView layer to eliminate XSS risks.
