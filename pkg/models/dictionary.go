package models

import "time"

type WordGender string

const (
	GenderMale      WordGender = "MALE"
	GenderFemale    WordGender = "FEMALE"
	GenderNeutral   WordGender = "NEUTRAL"
	GenderAmbiguous WordGender = "AMBIGUOUS"
)

type WordSearchResult struct {
	IDWord          int     `json:"id_word"`
	Lemma           string  `json:"lemma"`
	ShortDefinition string  `json:"short_definition"`
	FlagEmoji       string  `json:"flag_emoji"`
	CountryCode     string  `json:"country_code"`
	Similarity      float32 `json:"similarity"`
}

type Example struct {
	IDExample int     `json:"id_example"`
	IDSense   int     `json:"id_sense"`
	Sentence  string  `json:"sentence"`
	Source    *string `json:"source,omitempty"`
}

type HeteronymItem struct {
	RegionCode  string `json:"region_code"`
	RegionName  string `json:"region_name"`
	CountryCode string `json:"country_code"`
	FlagEmoji   string `json:"flag_emoji"`
	Term        string `json:"term"`
}

type SenseDetail struct {
	IDSense         int       `json:"id_sense"`
	IDWord          int       `json:"id_word"`
	RegionCode      string    `json:"region_code"`
	RegionName      string    `json:"region_name"`
	CountryCode     string    `json:"country_code"`
	FlagEmoji       string    `json:"flag_emoji"`
	ShortDefinition string    `json:"short_definition"`
	Definition      string    `json:"definition"`
	IsPrimary       bool      `json:"is_primary"`
	IsTabooOrVulgar bool      `json:"is_taboo_or_vulgar"`
	CulturalNote    *string   `json:"cultural_note,omitempty"`
	ImageURL        *string   `json:"image_url,omitempty"`
	CEFRLevel       string    `json:"cefr_level"`
	UsageFrequency  string    `json:"usage_frequency"`
	State           string    `json:"state"`
	Examples        []Example `json:"examples"`
}

type WordDetail struct {
	IDWord          int             `json:"id_word"`
	Lemma           string          `json:"lemma"`
	NormalizedLemma string          `json:"normalized_lemma"`
	PartOfSpeech    string          `json:"part_of_speech"`
	Gender          WordGender      `json:"gender"`
	Etymology       *string         `json:"etymology,omitempty"`
	FriendlyBlurb   *string         `json:"friendly_blurb,omitempty"`
	ImageURL        *string         `json:"image_url,omitempty"`
	OriginRegion    *Region         `json:"origin_region,omitempty"`
	Senses          []SenseDetail   `json:"senses"`
	Heteronyms      []HeteronymItem `json:"heteronyms"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

// Filtro de exploración paginada (GET /api/v1/words)
type WordExploreFilter struct {
	Query          string `json:"q"`
	RegionCode     string `json:"region"`        // Ej. 'MX', 'AR', 'ES'
	PartOfSpeech   string `json:"pos"`           // Ej. 'sustantivo', 'verbo', 'adjetivo', 'modismo'
	CEFRLevel      string `json:"cefr"`          // 'A1'..'C2', 'NATIVO'
	UsageFrequency string `json:"frequency"`     // 'EVERYDAY', 'COMMON', 'RARE', 'SLANG', 'ARCHAIC'
	IncludeTaboo   bool   `json:"include_taboo"` // Default true
	Page           int    `json:"page"`          // Default 1
	Limit          int    `json:"limit"`         // Default 20, máx 50
}

type WordSummaryCard struct {
	IDWord          int        `json:"id_word"`
	Lemma           string     `json:"lemma"`
	PartOfSpeech    string     `json:"part_of_speech"`
	Gender          WordGender `json:"gender"`
	FriendlyBlurb   *string    `json:"friendly_blurb,omitempty"`
	ImageURL        *string    `json:"image_url,omitempty"`
	PrimarySenseID  int        `json:"primary_sense_id"`
	ShortDefinition string     `json:"short_definition"`
	RegionCode      string     `json:"region_code"`
	RegionName      string     `json:"region_name"`
	FlagEmoji       string     `json:"flag_emoji"`
	CEFRLevel       string     `json:"cefr_level"`
	UsageFrequency  string     `json:"usage_frequency"`
	IsTabooOrVulgar bool       `json:"is_taboo_or_vulgar"`
	ApprovedSenses  int        `json:"approved_senses_count"`
}

type PaginatedWords struct {
	Items      []WordSummaryCard `json:"items"`
	TotalCount int               `json:"total_count"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
}

// --- DTOs DE AUTORÍA EDITORIAL (CU06) ---

type CreateExampleInput struct {
	Sentence string  `json:"sentence"`         // Obligatorio, mín. 5 chars
	Source   *string `json:"source,omitempty"` // Opcional
}

type CreateSenseRequest struct {
	RegionCode      string               `json:"region_code"`      // Obligatorio ('MX', 'AR', etc.)
	ShortDefinition string               `json:"short_definition"` // Obligatorio, <= 255 chars
	Definition      string               `json:"definition"`       // Obligatorio
	IsPrimary       bool                 `json:"is_primary"`
	IsTabooOrVulgar bool                 `json:"is_taboo_or_vulgar"`
	CulturalNote    *string              `json:"cultural_note,omitempty"`
	ImageURL        *string              `json:"image_url,omitempty"`   // Opcional, URL procesada por image-service (RN12)
	CEFRLevel       string               `json:"cefr_level"`            // Default 'B1'
	UsageFrequency  string               `json:"usage_frequency"`       // Default 'COMMON'
	ConceptIDs      []int                `json:"concept_ids,omitempty"` // Opcional: vincular heterónimos
	Examples        []CreateExampleInput `json:"examples,omitempty"`    // Ejemplos iniciales
}

type CreateWordProposalRequest struct {
	Lemma            string             `json:"lemma"`          // Obligatorio
	PartOfSpeech     string             `json:"part_of_speech"` // Obligatorio
	Gender           WordGender         `json:"gender"`         // Default 'NEUTRAL'
	Etymology        *string            `json:"etymology,omitempty"`
	FriendlyBlurb    *string            `json:"friendly_blurb,omitempty"`
	ImageURL         *string            `json:"image_url,omitempty"` // Opcional, URL procesada por image-service (RN12)
	OriginRegionCode *string            `json:"origin_region_code,omitempty"`
	InitialSense     CreateSenseRequest `json:"initial_sense"` // Obligatorio
}

type UpdateWordRequest struct {
	Lemma            string     `json:"lemma"`
	PartOfSpeech     string     `json:"part_of_speech"`
	Gender           WordGender `json:"gender"`
	Etymology        *string    `json:"etymology,omitempty"`
	FriendlyBlurb    *string    `json:"friendly_blurb,omitempty"`
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
	ConceptKey           string `json:"concept_key"`           // Ej. 'street_market'
	UniversalDescription string `json:"universal_description"` // Ej. 'Mercado ambulante tradicional al aire libre'
}
