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
	OriginRegion    *Region         `json:"origin_region,omitempty"`
	Senses          []SenseDetail   `json:"senses"`
	Heteronyms      []HeteronymItem `json:"heteronyms"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}
