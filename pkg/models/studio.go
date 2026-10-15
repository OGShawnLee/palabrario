package models

import "time"

type MasteryStage string

const (
	StageNew       MasteryStage = "NEW"
	StageLearning  MasteryStage = "LEARNING"
	StageReviewing MasteryStage = "REVIEWING"
	StageMastered  MasteryStage = "MASTERED"
)

type DeckItem struct {
	IDDeck    int          `json:"id_deck"`
	IDSense   int          `json:"id_sense"`
	SortOrder int          `json:"sort_order"`
	Sense     *SenseDetail `json:"sense,omitempty"`
}

type Deck struct {
	IDDeck      int        `json:"id_deck"`
	IDUser      string     `json:"id_user"`
	Title       string     `json:"title"`
	Description *string    `json:"description,omitempty"`
	IsPublic    bool       `json:"is_public"`
	IsOfficial  bool       `json:"is_official"`
	ItemCount   int        `json:"item_count"`
	Items       []DeckItem `json:"items,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type Collection struct {
	IDCollection int       `json:"id_collection"`
	IDUser       string    `json:"id_user"`
	Title        string    `json:"title"`
	Description  *string   `json:"description,omitempty"`
	IsPublic     bool      `json:"is_public"`
	IsFeatured   bool      `json:"is_featured"`
	IsOfficial   bool      `json:"is_official"`
	Decks        []Deck    `json:"decks,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserMastery struct {
	IDUser              string       `json:"id_user"`
	IDSense             int          `json:"id_sense"`
	ProgressStage       MasteryStage `json:"progress_stage"`
	MasteryPercentage   float32      `json:"mastery_percentage"`
	ConsecutiveCorrect  int          `json:"consecutive_correct"`
	TimesReviewed       int          `json:"times_reviewed"`
	EasinessFactor      float32      `json:"easiness_factor"`
	CurrentIntervalDays int          `json:"current_interval_days"`
	NextReviewAt        time.Time    `json:"next_review_at"`
	LastReviewedAt      time.Time    `json:"last_reviewed_at"`
}
