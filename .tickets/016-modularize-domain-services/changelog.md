# Changelog — Ticket 016: Modularize Domain Services

Split the monolithic `.routing/DOMAIN.md` (271 lines) and `.routing/INTERFACES.md` (1051 lines) into 8 self-contained `*_service.md` files under `.routing/services/`. Rewrote `DOMAIN.md` as the root service dispatcher, deleted `INTERFACES.md`, and repointed the repository docs and active tickets to the new structure. Documentation only. Nothing was committed.

The first pass of this ticket placed the service files loose in `.routing/`. The updated ticket moved them into `.routing/services/`. This changelog describes the final state.

## Files

| Action | File | Lines | Sources |
| :--- | :--- | ---: | :--- |
| created | `.routing/services/identity_service.md` | 126 | DOMAIN §3.1 (Account, User) + ER line · INTERFACES Module 1 (`CU01`) |
| created | `.routing/services/governance_service.md` | 168 | DOMAIN §3.1 (EditorRequest, InvitationTicket/TicketRedemption) + ER lines · INTERFACES Module 2 (`CU02–CU04`) |
| created | `.routing/services/dictionary_service.md` | 243 | DOMAIN §3.2 (Region, Word, Sense, Example, SemanticConcept/ConceptMapping) + ER lines · INTERFACES Module 3 (`CU05–CU06`) |
| created | `.routing/services/curation_service.md` | 178 | DOMAIN §3.3 (state machine, SenseApproval, ContentReport) + ER lines · INTERFACES Module 4 (`CU07–CU09`) |
| created | `.routing/services/studio_service.md` | 129 | DOMAIN §3.4 (Deck/DeckItem/SavedDeck, Collection/CollectionDeck/SavedCollection) + ER lines · INTERFACES Module 5 (`CU10–CU11`) |
| created | `.routing/services/mastery_service.md` | 172 | DOMAIN §3.4 (UserMastery SM-2, XPTransaction/Achievement/UserAchievement) + ER lines · INTERFACES Module 6 (`CU12–CU13`) |
| created | `.routing/services/arena_service.md` | 262 | DOMAIN §3.5 (Question/MatchRound/QuizAttempt, Match/MatchParticipant) + ER lines · INTERFACES Module 7 (`CU14–CU16`) |
| created | `.routing/services/media_service.md` | 42 | `RN12` references collected from DOMAIN · INTERFACES Module 8 (`image-service`) |
| rewritten | `.routing/DOMAIN.md` | 191 (was 271) | Dispatcher table linking into `services/`, `RN01–RN13` index, philosophy, shared cross-cutting content |
| deleted | `.routing/INTERFACES.md` | 0 (was 1051) | Fully distributed |
| edited | `.routing/CONTEXT.md` | 48 (was 49) | `## Files`: `DOMAIN.md` entry reworded as dispatcher and pointing to `services/`; `INTERFACES.md` line removed |
| edited | `README.md` | 2 table rows changed | Blueprint table: `DOMAIN.md` is the dispatcher; a "Modular Domain Services" row lists `.routing/services/*_service.md` |
| edited | `.claude/skills/router/SKILL.md` | 1 line changed | Step 4 references `DOMAIN.md` (dispatcher) and `.routing/services/*_service.md`; `INTERFACES.md` removed |
| edited | `.tickets/009-build-dictionary-route/task.md` | — | Input → `.routing/services/dictionary_service.md` |
| edited | `.tickets/010-build-collections-route/task.md` | — | Input → `.routing/services/studio_service.md` |
| edited | `.tickets/011-build-decks-route/task.md` | — | Input → `.routing/services/studio_service.md` |
| edited | `.tickets/014-backend-search-fts/task.md` | — | Inputs → `.routing/services/dictionary_service.md` and `.routing/services/studio_service.md` |
| edited | `.routing/PROGRESS.md` | 1 line added vs. HEAD | Ticket 016 tracked as `in-progress` with the `services/` layout; no reference to ticket 017 |
| created | `.tickets/016-modularize-domain-services/changelog.md` | — | This file |

`.routing/` root now holds only `CONTEXT.md`, `DESIGN.md`, `DOMAIN.md`, `INFRASTRUCTURE.md`, `PROGRESS.md`, `ROUTES.md`, and `services/`. No loose `*_service.md` files remain.

Each service file has the same layout: header with service, use cases, rules, and routes; **Part A — Domain Model** (ER fragment + entities and rules); **Part B — Technical Contracts** (DTOs, SOA interface, `sqlc` queries, endpoints).

## Relative-link fixes after the move

The service files live one directory deeper, so their relative links were rewritten:

| Link target | Before (`.routing/`) | After (`.routing/services/`) |
| :--- | :--- | :--- |
| `SCHEMA_FINAL.sql`, `pkg/…`, `go.mod`, `go.work`, `local-app/…` | `../…` | `../../…` |
| `DOMAIN.md` (including the `#3-global-contract-first-rules-cross-cutting` anchor) | `DOMAIN.md` | `../DOMAIN.md` |
| `ROUTES.md` | `ROUTES.md` | `../ROUTES.md` |
| Sibling service files | `x_service.md` | `x_service.md` (unchanged) |

In `DOMAIN.md` every dispatcher and rule-index link now reads `services/<name>_service.md`.

## Where cross-cutting content went

Some source content belongs to no single service. It was kept, not dropped:

- **Root `DOMAIN.md`, section 3:** the contract-first rules (former INTERFACES §0).
- **Root `DOMAIN.md`, section 4:** SOA and Wails IPC architecture with its diagram (former §1), the `core-service` directory tree (§1.1), and the domain-error → HTTP status table (§1.2).
- **Root `DOMAIN.md`, header:** the stack and pattern lines from the former INTERFACES header.
- **ER diagram:** the single diagram in the former DOMAIN §2 is split by relationship into each service's *Entity-Relationship Fragment*. Every one of the 39 relationship lines appears in exactly one fragment.
- **`RN12`:** the former DOMAIN had no standalone `RN12` section, only per-field mentions. `media_service.md` Part A is a table of those mentions, with owning services, plus the related contract fields. No new rules were invented.

## Zero-content-loss verification

The service files were assembled by script from exact line ranges of the original files, so the extracted text is verbatim. Verification compared each non-blank, non-`---` line of the two originals against the union of the new files. It was re-run after the move to `services/`.

- `DOMAIN.md` original: 217 non-trivial lines, 3 not found verbatim. All are intentionally retitled headings: file title, former §2 and former §3.
- `INTERFACES.md` original: 891 non-trivial lines, 8 not found verbatim:
  - 7 are intentional: file title, header line 7 (the obsolete `Synchronized with` links), and headings §0, §1, §1.1, §1.2 and §2, renumbered or retitled into root `DOMAIN.md` sections 3 and 4.
  - 1 is the link-depth fix: the `pkg/auth/auth.go` link in `arena_service.md` now starts with `../../`.
- All DTO structs, `sqlc` query lists, SOA interfaces, endpoint tables, algorithms (SM-2, speed scoring), triggers, and error mappings are therefore present.

| Set | Result |
| :--- | :--- |
| `CU01–CU16` | All 16 present. Primary homes: `CU01` identity · `CU02–04` governance · `CU05–06` dictionary · `CU07–09` curation · `CU10–11` studio · `CU12–13` mastery · `CU14–16` arena |
| `RN01–RN13` | All 13 present. Primary homes: `RN01` identity · `RN02` governance · `RN03–05` studio · `RN06–08` curation · `RN09` and `RN11` mastery · `RN10` and `RN13` arena · `RN12` media |
| Mermaid / code fences | Balanced in all files |
| Encoding | UTF-8, no mojibake. Line endings and BOM of edited existing files are preserved (BOM restored on files where an earlier edit had dropped it) |

## Link validation

A script resolved every relative Markdown link (outside code fences) in `README.md`, `.claude/skills/router/SKILL.md`, `.routing/**`, and `.tickets/**`: 268 links checked.

- **All links introduced or changed by this ticket resolve:** every `services/…` link in `DOMAIN.md`, every cross-link between service files, every link from service files back to `../DOMAIN.md`, `../ROUTES.md` and the repo root, and the links in `README.md` and tickets 009, 010, 011, and 014.
- **Remaining unresolved links, all outside this ticket's declared scope:**

| File | Link | Note |
| :--- | :--- | :--- |
| `.tickets/008-translate-readme/changelog.md` | `../../.routing/INTERFACES.md` | Completed ticket, historical record. The file no longer exists by design. |
| `.tickets/016-modularize-domain-services/task.md` | `../../.routing/INTERFACES.md`, `services/*.md`, `.routing/DOMAIN.md` | Router-owned ticket text. The `services/` and `.routing/` targets appear as table or example text, not as links from the ticket folder. |
| `.tickets/008-translate-readme/task.md`, `.tickets/002-translate-interfaces/*` | mentions of `INTERFACES.md` | Historical tickets, plain text. |
| `.routing/PROGRESS.md:47` | `.routing/INTERFACES.md` | Plain text in a completed-ticket history line, not a link. |
| `README.md` | `core-service/.env`, `realtime-service/.env`, `image-service/.env` | Pre-existing. Git-ignored env files. |
| `.tickets/001`, `002` | backtick-wrapped link text | Pre-existing markup. |
| `.tickets/006`, `012` | `WindowBar.md`, `ROUTES.md` | Pre-existing. |
| `.tickets/009`, `013` | screenshots under `.t3/userdata/attachments/` | Pre-existing, outside the repo. |

## Notes for the router

1. **Existing source typos preserved verbatim** (zero-loss). These are in the original text, not introduced here:
   - `WordSummaryCard` lists `ImageURL` twice (`dictionary_service.md`).
   - `UpdateWordRequest` lists `ImageURL` twice (`dictionary_service.md`).
2. **Entity numbering** (`#### 1.` … `#### 17.`) and heading levels of the sources are preserved. Numbers therefore look non-contiguous inside a service file (e.g. governance starts at `#### 3.`).
3. **`RN08` spans several services** (tickets, sense override, moderation, admin `CreateSense`). The root index points to `curation_service.md`; it is also cited in governance, dictionary, and studio.
4. **Ticket 017:** the `.tickets/017-update-docs-modular-services/` folder was no longer present when this pass started, and no `017` references remain in `README.md`, `.routing/`, `.claude/skills/router/SKILL.md`, or other tickets.
5. **Pre-existing edits in the working tree.** When this pass started, `README.md`, the router `SKILL.md`, `PROGRESS.md`, and tickets 009–011 and 014 already had edits pointing at the loose `.routing/*_service.md` paths, and several had lost their BOM. Those edits were updated to the `services/` paths. The BOMs were restored on `PROGRESS.md` and tickets 009, 010, 011, and 014. `.tickets/009-build-dictionary-route/task.md` and `spec.md` carry other, unrelated uncommitted edits that were left as they are.

## Scope compliance

No application code, `core-service/`, `realtime-service/`, `image-service/`, `pkg/`, `sql/`, or `.context/` files were modified.
