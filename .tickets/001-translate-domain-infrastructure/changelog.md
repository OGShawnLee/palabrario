# Changelog — Ticket 001: Translate DOMAIN.md and INFRASTRUCTURE.md

## 1. Section 3.6 Omission
- Confirmed: Section 3.6 (`### 3.6. Incremento Arquitectónico a Largo Plazo: Eventos Culturales de Temporada (Seasonal Events)`) and its associated content in `.context/DOMAIN.md` have been completely omitted from `.routing/DOMAIN.md`. No placeholder or heading remains.
- All preceding section numbers (1, 2, 3, 3.1, 3.2, 3.3, 3.4, 3.5) and entity numbers (1 to 17) have been preserved without modification.

## 2. Line Counts and Spec Verification Metrics

| Metric | Reference (.context/) | Translated (.routing/) | Difference / Status |
| :--- | :--- | :--- | :--- |
| **DOMAIN.md lines** | 273 (excluding 3.6) | 270 | Diff = 3 lines (Within allowable limit of $\pm 100$) |
| **DOMAIN.md headings (`#`)** | 26 | 26 | Exact match |
| **DOMAIN.md code fences (`````)** | 4 | 4 | Exact match |
| **DOMAIN.md markdown links (`](`)** | 1 | 1 | Exact match |
| **INFRASTRUCTURE.md lines** | 149 | 149 | Diff = 0 lines (Within allowable limit of $\pm 100$) |
| **INFRASTRUCTURE.md headings (`#`)** | 16 (11 MD + 5 comments) | 16 (11 MD + 5 comments) | Exact match |
| **INFRASTRUCTURE.md code fences (`````)** | 12 | 12 | Exact match |
| **INFRASTRUCTURE.md markdown links (`](`)** | 5 | 5 | Exact match |

## 3. Terms and Translation Decisions
- **Lema Canónico** $\to$ Translated as `Canonical Lemma`.
- **Acepción o Sentido Dialectal** $\to$ Translated as `Dialectal Sense or Meaning` / `Sense`.
- **Dialecto Ancla** $\to$ Translated as `Anchor Dialect`.
- **Reseña Cultural Amena** $\to$ Translated as `Engaging Cultural Blurb`.
- **Heteronimia Transdialectal** $\to$ Translated as `Transdialectal Heteronymy`.
- **Reactivo** $\to$ Translated as `Prompt` / `Curricular prompt`.
- **Baraja** $\to$ Translated as `Deck`.
- **Libro mayor** $\to$ Translated as `Ledger` / `Auditable ledger`.
- **Quórum** $\to$ Kept as `Quorum`.
- **Unchanged Terms:** Maintained `RBAC`, `FTS`, `SM-2`, `CEFR`, `XP`, `JWT`, `S2S`, `WebSocket`, `PostgreSQL`, `sqlc`, `pgxpool`, `Chi`, `Wails`, `Vue`, `UnoCSS`, `Docker`, `Vagrant`, and roles `LEARNER`, `EDITOR`, `ADMIN` strictly as instructed.

## 4. Untranslated Content
- None. All Spanish prose, diagram annotations, and table descriptions were translated into English. Spanish remains only in proper nouns, language samples / linguistic examples (e.g., `"guagua"`, `"coger"`, `"apapachar"`), database identifiers, and code.

## 5. Proposed Glossary Entries
- **Sense (`Acepción / Sentido`):** Concrete dialectal meaning of a lemma within a specific region.
- **Anchor Dialect (`Dialecto Ancla`):** The primary dialect zone chosen by a user against which searches and cultural contrasts are prioritized.
- **Engaging Cultural Blurb (`Reseña Cultural Amena`):** Concise cultural narrative illustrating the soul and regional contrasts of a lemma.
- **Transdialectal Heteronymy (`Heteronimia Transdialectal`):** Distinct regional lexical terms referring to the same underlying concept.
