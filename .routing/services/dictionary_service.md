# Dictionary Service (`dictionary_service.md`)

> **Domain service file — self-contained.** Combines the domain entities and business rules (from the former `DOMAIN.md`) with the technical contracts: DTOs, SOA interface, `sqlc` queries, and HTTP/WS endpoints (from the former `INTERFACES.md`).
> **Service:** `DictionaryService` · **Use Cases:** `CU05`, `CU06` · **Domain:** Lexicon, Senses, Heteronyms (`SemanticConcept`/`ConceptMapping`), Full-Text Search · **Routes:** `/dictionary` (see [`ROUTES.md`](../ROUTES.md))
> **Global rules:** the contract-first rules, SOA architecture, and domain-error-to-HTTP mapping shared by every service live in [`DOMAIN.md`](../DOMAIN.md) (sections 3 and 4).
> **Canonical persistence:** PostgreSQL 17 ([`SCHEMA_FINAL.sql`](../../SCHEMA_FINAL.sql)) · **Shared DTOs:** [`pkg/models/`](../../pkg/models)

---

## Part A — Domain Model

### Entity-Relationship Fragment

```mermaid
erDiagram
    REGIONS ||--o{ WORDS : "endemic origin"
    REGIONS ||--o{ SENSES : "contextualizes dialect"
    WORDS ||--|{ SENSES : "groups senses"
    SENSES ||--o{ EXAMPLES : "exemplifies in context"
    USERS ||--o{ SENSES : "proposes (id_editor)"
    USERS ||--o{ EXAMPLES : "drafts (id_editor)"

    SEMANTIC_CONCEPTS ||--|{ CONCEPT_MAPPINGS : "unifies meaning"
    SENSES ||--o{ CONCEPT_MAPPINGS : "regional heteronym"
```

### 3.2. Context 2: Pan-Hispanic Lexicographical Core and Heteronymy (`Lexicon & Dialectology`)
This is the linguistic heart of Palabrario (**CU05, CU06**). It models regional polysemy, speech register, taboo weight, and transdialectal equivalences.

#### 5. `Region` (Table: `regions`)
Represents a Spanish-speaking country or recognized dialectal macro-zone (`MX`, `ES`, `AR`, `CO`, `CL`, `CU`, `CAM`, `PAN`).
* **Identifier:** `id_region` (`UUID`, PK) and natural key `code` (`VARCHAR(10)`, unique).
* **Key Attributes:** `name` (`VARCHAR(100)`), `parent_zone` (`Norteamérica`, `Rioplatense`, `Andina`, `Caribe`, `Peninsular`, `General`), `country_code` (`ISO-3166-1 alpha-2`), `flag_emoji` (`VARCHAR(10)`).

#### 6. `Word` (Table: `words`)
Represents the **Canonical Lemma** (singular noun, infinitive verb, or set phrase). Acts as the Aggregate Root for a word's senses.
* **Identifier:** `id_word` (`SERIAL / INT`, PK).
* **Key Attributes:**
  * `lemma` (`VARCHAR(100)`): Canonical orthographic form with accents and diacritics (e.g., `"apapachar"`, `"tianguis"`).
  * `normalized_lemma` (`VARCHAR(100)`): Unaccented lowercase lemma, indexed with `pg_trgm` (`gin_trgm_ops`) and `varchar_pattern_ops` for instant autocomplete and typo tolerance (**RF06**).
  * `part_of_speech` (`VARCHAR(50)`): Part of speech (`sustantivo`, `verbo`, `adjetivo`, `modismo`).
  * `gender` (`word_gender`: `MALE | FEMALE | NEUTRAL | AMBIGUOUS`).
  * `etymology` (`TEXT`): Historical or indigenous root (Nahuatl, Quechua, Maya, Latin, Lunfardo, etc.).
  * `friendly_blurb` (`TEXT`): **Engaging Cultural Blurb**. Brief narrative written with humor and wit explaining the soul of the word and its contrasts across Spanish-speaking countries.
  * `image_url` (`TEXT`, optional): URL to conceptual or representative image generated and sanitized by `image-service` (**RN12**).
  * `id_origin_region` (`UUID`, optional FK $\to$ `regions`): Identifies the region if the term is endemic to a specific region; if `NULL` or `PAN`, it belongs to general Spanish.
  * `search_vector` (`tsvector` generated and stored): Weights `lemma` with weight `'A'` and `friendly_blurb` with weight `'B'` under the `spanish_unaccent` configuration.

#### 7. `Sense` (Table: `senses`)
Represents a specific **Dialectal Sense or Meaning** of a `Word` in a given `Region`. All study mechanics (`DeckItem`, `UserMastery`, `Question`) point to `Sense` and **not** to `Word`, because the learner studies a concrete meaning within a specific regional context.
* **Identifier:** `id_sense` (`SERIAL / INT`, PK).
* **Key Attributes:**
  * `id_word` (`INT`, FK $\to$ `words`), `id_region` (`UUID`, FK $\to$ `regions`), `id_editor` (`UUID`, FK $\to$ `users`, proposal author).
  * `short_definition` (`VARCHAR(255)`): Concise gloss for the Master-Detail search view and trivia options.
  * `definition` (`TEXT`): Full formal lexicographical definition.
  * `is_primary` (`BOOLEAN`): Indicates whether it is the primary sense of the lemma to prioritize it in searches and trivias.
  * `is_taboo_or_vulgar` (`BOOLEAN`) and `cultural_note` (`TEXT`): Essential pragmatic warning for interdialectal false friends (e.g., `"coger"` in Mexico/Argentina vs Spain).
  * `image_url` (`TEXT`, optional): URL to image specific to the concrete dialectal sense processed by `image-service` (**RN12**).
  * `cefr_level` (`cefr_level`: `A1 | A2 | B1 | B2 | C1 | C2 | NATIVO`) and `usage_frequency` (`EVERYDAY | COMMON | RARE | SLANG | ARCHAIC`).
  * `state` (`review_state`: `DRAFT | APPROVED | HIDDEN`): Editorial lifecycle. Only senses in `APPROVED` state are visible in public dictionary queries (**CU05, CU08**).
  * `sense_search_vector` (`tsvector` generated): Indexes `short_definition` (`A`), `definition` (`B`), and `cultural_note` (`C`).

#### 8. `Example` (Table: `examples`)
Real-world contextual sentence (literary, journalistic, or colloquial) illustrating the pragmatic usage of a `Sense`, serving as input for sentence-completion prompts (`SENTENCE_COMPLETION`).
* **Identifier:** `id_example` (`SERIAL / INT`, PK).
* **Key Attributes:** `id_sense` (`INT`, FK $\to$ `senses`), `id_editor` (`UUID`, FK $\to$ `users`), `sentence` (`TEXT`), `source` (`VARCHAR(255)`).

#### 9. `SemanticConcept` and `ConceptMapping` (Tables: `semantic_concepts`, `concept_mappings`)
Model the phenomenon of **Transdialectal Heteronymy** (*"How do they say it over there?"*).
* **Problem Solved:** Directly associating synonyms word-to-word ($N \times N$) would yield an unmanageable and ambiguous graph whenever a word is polysemous (such as *guagua* = bus and *guagua* = baby).
* **Solution:** `SemanticConcept` defines an unambiguous universal concept (`concept_key = 'public_bus'` vs `concept_key = 'infant_baby'`), and `ConceptMapping` links each specific regional `id_sense` to that concept. Thus, when looking up *guagua* in its Cuban sense (`id_sense = 103`), the system automatically discovers its heteronyms for `public_bus` (*camión*, *colectivo*, *micro*) without conflating them with those for `infant_baby`.

---

## Part B — Technical Contracts

> [!IMPORTANT]
> Contract-first rules apply: no invented JSON fields, strict `snake_case`, typed `sqlc` queries only, and `response.Error(w, status, message)` for every error. See [`DOMAIN.md`](../DOMAIN.md#3-global-contract-first-rules-cross-cutting).

### MODULE 3: Panhispanic Dictionary, Lexical Authorship, and Heteronymy (`DictionaryService` — `CU05, CU06`)

#### 2.3.1. Complete Data Structures (`pkg/models/dictionary.go`)

In addition to already existing structures (`WordGender`, `WordSearchResult`, `Example`, `HeteronymItem`, `SenseDetail`, `WordDetail`), the following are defined in `pkg/models/dictionary.go`:

```go
// Paginated exploration filter (GET /api/v1/words)
type WordExploreFilter struct {
	Query          string `json:"q"`
	RegionCode     string `json:"region"`          // E.g. 'MX', 'AR', 'ES'
	PartOfSpeech   string `json:"pos"`             // E.g. 'sustantivo', 'verbo', 'adjetivo', 'modismo'
	CEFRLevel      string `json:"cefr"`            // 'A1'..'C2', 'NATIVO'
	UsageFrequency string `json:"frequency"`       // 'EVERYDAY', 'COMMON', 'RARE', 'SLANG', 'ARCHAIC'
	IncludeTaboo   bool   `json:"include_taboo"`   // Default true
	Page           int    `json:"page"`            // Default 1
	Limit          int    `json:"limit"`           // Default 20, max 50
}

type WordSummaryCard struct {
	IDWord           int        `json:"id_word"`
	Lemma            string     `json:"lemma"`
	PartOfSpeech     string     `json:"part_of_speech"`
	Gender           WordGender `json:"gender"`
	FriendlyBlurb    *string    `json:"friendly_blurb,omitempty"`
	ImageURL         *string    `json:"image_url,omitempty"`
	ImageURL         *string    `json:"image_url,omitempty"`
	PrimarySenseID   int        `json:"primary_sense_id"`
	ShortDefinition  string     `json:"short_definition"`
	RegionCode       string     `json:"region_code"`
	RegionName       string     `json:"region_name"`
	FlagEmoji        string     `json:"flag_emoji"`
	CEFRLevel        string     `json:"cefr_level"`
	UsageFrequency   string     `json:"usage_frequency"`
	IsTabooOrVulgar  bool       `json:"is_taboo_or_vulgar"`
	ApprovedSenses   int        `json:"approved_senses_count"`
}

type PaginatedWords struct {
	Items      []WordSummaryCard `json:"items"`
	TotalCount int               `json:"total_count"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
}

// --- EDITORIAL AUTHORSHIP DTOs (CU06) ---

type CreateExampleInput struct {
	Sentence string  `json:"sentence"`          // Required, min. 5 chars
	Source   *string `json:"source,omitempty"`  // Optional
}

type CreateSenseRequest struct {
	RegionCode      string               `json:"region_code"`               // Required ('MX', 'AR', etc.)
	ShortDefinition string               `json:"short_definition"`          // Required, <= 255 chars
	Definition      string               `json:"definition"`                // Required
	IsPrimary       bool                 `json:"is_primary"`
	IsTabooOrVulgar bool                 `json:"is_taboo_or_vulgar"`
	CulturalNote    *string              `json:"cultural_note,omitempty"`
	ImageURL        *string              `json:"image_url,omitempty"`        // Optional, URL processed by image-service (RN12)
	CEFRLevel       string               `json:"cefr_level"`                // Default 'B1'
	UsageFrequency  string               `json:"usage_frequency"`           // Default 'COMMON'
	ConceptIDs      []int                `json:"concept_ids,omitempty"`     // Optional: link heteronyms
	Examples        []CreateExampleInput `json:"examples,omitempty"`        // Initial examples
}

type CreateWordProposalRequest struct {
	Lemma            string             `json:"lemma"`                      // Required
	PartOfSpeech     string             `json:"part_of_speech"`             // Required
	Gender           WordGender         `json:"gender"`                     // Default 'NEUTRAL'
	Etymology        *string            `json:"etymology,omitempty"`
	FriendlyBlurb    *string            `json:"friendly_blurb,omitempty"`
	ImageURL         *string            `json:"image_url,omitempty"`         // Optional, URL processed by image-service (RN12)
	OriginRegionCode *string            `json:"origin_region_code,omitempty"`
	InitialSense     CreateSenseRequest `json:"initial_sense"`              // Required
}

type UpdateWordRequest struct {
	Lemma            string     `json:"lemma"`
	PartOfSpeech     string     `json:"part_of_speech"`
	Gender           WordGender `json:"gender"`
	Etymology        *string    `json:"etymology,omitempty"`
	FriendlyBlurb    *string    `json:"friendly_blurb,omitempty"`
	ImageURL         *string    `json:"image_url,omitempty"`
	ImageURL         *string    `json:"image_url,omitempty"`
	OriginRegionCode *string    `json:"origin_region_code,omitempty"`
}

type UpdateSenseRequest struct {
	RegionCode      string  `json:"region_code"`
	ShortDefinition string  `json:"short_definition"`
	Definition      string  `json:"definition"`
	IsPrimary       bool    `json:"is_primary"`
	IsTabooOrVulgar bool    `json:"is_taboo_or_vulgar"`
	CulturalNote    *string `json:"cultural_note,omitempty"`
	ImageURL        *string `json:"image_url,omitempty"`
	CEFRLevel       string  `json:"cefr_level"`
	UsageFrequency  string  `json:"usage_frequency"`
}

type SemanticConcept struct {
	IDConcept            int    `json:"id_concept"`
	ConceptKey           string `json:"concept_key"`
	UniversalDescription string `json:"universal_description"`
	LinkedSensesCount    int    `json:"linked_senses_count"`
}

type CreateConceptRequest struct {
	ConceptKey           string `json:"concept_key"`           // E.g. 'street_market'
	UniversalDescription string `json:"universal_description"` // E.g. 'Mercado ambulante tradicional al aire libre'
}
```

#### 2.3.2. Exact `sqlc` Queries in `sql/queries/dictionary.sql`

In addition to the 5 base queries (`SearchVocabularyDropdown`, `GetWordByID`, `GetApprovedSensesByWordID`, `GetExamplesBySenseIDs`, `GetHeteronymsByWordID`), it includes:

1. `-- name: ListRegions :many`
   Selects all rows from `regions` ordered by `name ASC`.
2. `-- name: GetRegionByCode :one`
   Retrieves `id_region, code, name, parent_zone, country_code, flag_emoji` where `code = $1`.
3. `-- name: GetWordOfTheDayID :one`
   Deterministically selects an `id_word` that has at least one `APPROVED` sense (prioritizing region `$1` if supplied, using `MOD(EXTRACT(DOY FROM CURRENT_DATE)::INT, COUNT(*))` or `friendly_blurb IS NOT NULL`) to subsequently assemble its complete `WordDetail` using existing queries.
4. `-- name: ExploreWordsPaginated :many` and `-- name: CountExploreWords :one`
   Filters `words w JOIN senses s ON s.id_word = w.id_word AND s.state = 'APPROVED' JOIN regions r ON r.id_region = s.id_region` applying optional filters (`sqlc.narg`): `region_code`, `part_of_speech`, `cefr_level`, `usage_frequency`, `include_taboo`, and search query `q`, paginated with `LIMIT $1 OFFSET $2`.
5. `-- name: GetWordByNormalizedLemmaAndPOS :one`
   Checks if a lemma already exists with `(normalized_lemma = $1 AND part_of_speech = $2)`. If it already exists during `CreateWordProposal`, reuses the existing `id_word` and adds the new `Sense` in `DRAFT` (**CU06 Main Flow Step 2**).
6. `-- name: CreateWord :one`
   Inserts into `words (lemma, normalized_lemma, part_of_speech, gender, etymology, friendly_blurb, id_origin_region)` and returns the created row.
7. `-- name: UpdateWordMetadata :one`
   Updates metadata in `words` and `updated_at = now()` by `id_word`.
8. `-- name: CreateSense :one`
   Inserts into `senses (id_word, id_region, id_editor, short_definition, definition, is_primary, is_taboo_or_vulgar, cultural_note, cefr_level, usage_frequency, state)`.
   - **Business Rule in `CreateSense`:** If the actor has the `ADMIN` role, they may choose to create directly in `'APPROVED'` (`RN08`); if `EDITOR`, it is always inserted with `state = 'DRAFT'` (`CU06 Step 5`).
9. `-- name: GetSenseByID :one`
   Retrieves the sense (any `state`) along with its `id_editor` to verify authorship before edit/delete (`CU06 2a`).
10. `-- name: UpdateSense :one` and `-- name: DeleteSense :exec`
    Updates or deletes a sense. If the actor is `EDITOR` (not `ADMIN`), verifies that `sense.id_editor == actorID` and `sense.state == 'DRAFT'`.
11. `-- name: CreateExample :one`, `-- name: UpdateExample :one`, `-- name: DeleteExample :exec`
    CRUD of example sentences in `examples`.
12. `-- name: ListSemanticConcepts :many`, `-- name: CreateSemanticConcept :one`, `-- name: LinkSenseToConcept :exec`, `-- name: UnlinkSenseFromConcept :exec`
    Management of the heteronym graph in `semantic_concepts` and `concept_mappings`.

#### 2.3.3. Dictionary HTTP Endpoints Specification (`/api/v1/regions`, `/api/v1/words`, `/api/v1/senses`, `/api/v1/concepts`)

| Method | Path | Middleware | Use Case | Status | Successful Response |
| :---: | :--- | :---: | :---: | :---: | :--- |
| `GET` | `/api/v1/regions` | Public | CU05 | ✅ Ready | `200 OK` $\to$ `[]models.Region` |
| `GET` | `/api/v1/words/word-of-the-day` | Public | CU05 | ✅ Ready | `200 OK` $\to$ `models.WordDetail` (`?dialect=MX`) |
| `GET` | `/api/v1/words/search` | Public | CU05 (RF06) | ✅ Ready | `200 OK` $\to$ `[]models.WordSearchResult` (`?q=&dialect=&limit=8`) |
| `GET` | `/api/v1/words` | Public | CU05 (RF07) | ✅ Ready | `200 OK` $\to$ `models.PaginatedWords` |
| `GET` | `/api/v1/words/{id}` | Public | CU05 (RF07) | ✅ Ready | `200 OK` $\to$ `models.WordDetail` |
| `POST` | `/api/v1/words` | `RequireRole(EDITOR, ADMIN)` | CU06 (RF08) | ✅ Ready | `201 Created` $\to$ `models.WordDetail` |
| `PUT` | `/api/v1/words/{id}` | `RequireRole(EDITOR, ADMIN)` | CU06 | ✅ Ready | `200 OK` $\to$ `models.WordDetail` |
| `POST` | `/api/v1/words/{id}/senses` | `RequireRole(EDITOR, ADMIN)` | CU06 (RF08) | ✅ Ready | `201 Created` $\to$ `models.SenseDetail` |
| `PUT` | `/api/v1/senses/{senseId}` | `RequireRole(EDITOR, ADMIN)` | CU06 (2a) | ✅ Ready | `200 OK` $\to$ `models.SenseDetail` (`403` if another editor or already approved) |
| `DELETE` | `/api/v1/senses/{senseId}` | `RequireRole(EDITOR, ADMIN)` | CU06 | ✅ Ready | `204 No Content` |
| `POST` | `/api/v1/senses/{senseId}/examples` | `RequireRole(EDITOR, ADMIN)` | CU06 | ✅ Ready | `201 Created` $\to$ `models.Example` |
| `PUT` | `/api/v1/examples/{exampleId}` | `RequireRole(EDITOR, ADMIN)` | CU06 | ✅ Ready | `200 OK` $\to$ `models.Example` |
| `DELETE` | `/api/v1/examples/{exampleId}` | `RequireRole(EDITOR, ADMIN)` | CU06 | ✅ Ready | `204 No Content` |
| `GET` | `/api/v1/concepts` | Public | CU05 / CU06 | ✅ Ready | `200 OK` $\to$ `[]models.SemanticConcept` (`?q=`) |
| `POST` | `/api/v1/concepts` | `RequireRole(EDITOR, ADMIN)` | CU06 | ✅ Ready | `201 Created` $\to$ `models.SemanticConcept` |
| `POST` | `/api/v1/senses/{senseId}/concepts/{conceptId}` | `RequireRole(EDITOR, ADMIN)` | CU06 | ✅ Ready | `204 No Content` |
| `DELETE` | `/api/v1/senses/{senseId}/concepts/{conceptId}` | `RequireRole(EDITOR, ADMIN)` | CU06 | ✅ Ready | `204 No Content` |
