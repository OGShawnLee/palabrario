package models

import "time"

type MasteryStage string

const (
	StageNew       MasteryStage = "NEW"
	StageLearning  MasteryStage = "LEARNING"
	StageReviewing MasteryStage = "REVIEWING"
	StageMastered  MasteryStage = "MASTERED"
)

const (
	XPReasonQuizCorrect         = "QUIZ_CORRECT"
	XPReasonMatchWin            = "MATCH_WIN"
	XPReasonDailyStreak         = "DAILY_STREAK"
	XPReasonDeckCompleted       = "DECK_COMPLETED"
	XPReasonAchievementUnlocked = "ACHIEVEMENT_UNLOCKED"
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

// --- DTOs STUDIO PERSONAL (Sección 2.5.1 de INTERFACES.md) ---

type StudioDeckSummary struct {
	Deck
	OwnerUserName    string  `json:"owner_user_name"`
	IsSavedFromComm  bool    `json:"is_saved_from_community"`
	DueCardsToday    int     `json:"due_cards_today"`
	RetentionPercent float32 `json:"retention_percent"`
}

type CreateDeckRequest struct {
	Title       string  `json:"title"` // Obligatorio, 3..150 chars
	Description *string `json:"description,omitempty"`
	IsPublic    bool    `json:"is_public"`
	IsOfficial  bool    `json:"is_official"`         // Solo respetado si role IN ('EDITOR', 'ADMIN') (RN05)
	SenseIDs    []int   `json:"sense_ids,omitempty"` // Tarjetas iniciales opcionales
}

type UpdateDeckRequest struct {
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	IsPublic    bool    `json:"is_public"`
	IsOfficial  bool    `json:"is_official"`
}

type AddDeckItemRequest struct {
	IDSense   int  `json:"id_sense"`
	SortOrder *int `json:"sort_order,omitempty"`
}

type CreateCollectionRequest struct {
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	IsPublic    bool    `json:"is_public"`
	IsFeatured  bool    `json:"is_featured"` // Solo ADMIN/EDITOR
	IsOfficial  bool    `json:"is_official"` // Solo ADMIN/EDITOR (RN05)
	DeckIDs     []int   `json:"deck_ids"`    // Ordenados según el arreglo
}

type UpdateCollectionRequest struct {
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	IsPublic    bool    `json:"is_public"`
	IsFeatured  bool    `json:"is_featured"`
	IsOfficial  bool    `json:"is_official"`
	DeckIDs     []int   `json:"deck_ids,omitempty"`
}

// --- DTOs MOTOR SM-2 Y REPETICIÓN ESPACIADA (Sección 2.6.1 de INTERFACES.md) ---

type StudyCardItem struct {
	IDSense         int          `json:"id_sense"`
	IDWord          int          `json:"id_word"`
	Lemma           string       `json:"lemma"`
	PartOfSpeech    string       `json:"part_of_speech"`
	RegionCode      string       `json:"region_code"`
	RegionName      string       `json:"region_name"`
	FlagEmoji       string       `json:"flag_emoji"`
	ShortDefinition string       `json:"short_definition"`
	Definition      string       `json:"definition"`
	CulturalNote    *string      `json:"cultural_note,omitempty"`
	Examples        []Example    `json:"examples"`
	CurrentMastery  *UserMastery `json:"current_mastery,omitempty"` // nil si es tarjeta NEW
}

type StudySessionPayload struct {
	DeckID       *int            `json:"deck_id,omitempty"`
	TotalDue     int             `json:"total_due"`
	SessionCards []StudyCardItem `json:"session_cards"`
}

type SubmitReviewRequest struct {
	IDSense       int `json:"id_sense"`       // Obligatorio
	QualityRating int `json:"quality_rating"` // 1: Otra vez (Fallo), 2: Difícil, 3: Bueno, 4: Fácil
	TimeSpentMs   int `json:"time_spent_ms"`  // Tiempo de respuesta en milisegundos
}

type CardReviewResult struct {
	UpdatedMastery     UserMastery         `json:"updated_mastery"`
	XPAwarded          int                 `json:"xp_awarded"`
	NewTotalXP         int                 `json:"new_total_xp"`
	NewLevel           int                 `json:"new_level"`
	CurrentStreak      int                 `json:"current_streak"`
	NewlyUnlockedBadge []AchievementStatus `json:"newly_unlocked_badges,omitempty"`
}

type XPTransactionItem struct {
	IDTransaction string    `json:"id_transaction"`
	AmountXP      int       `json:"amount_xp"`
	Reason        string    `json:"reason"` // 'QUIZ_CORRECT', 'MATCH_WIN', 'DAILY_STREAK', 'DECK_COMPLETED', 'ACHIEVEMENT_UNLOCKED'
	MetadataID    *string   `json:"metadata_id,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

type AchievementStatus struct {
	IDAchievement string     `json:"id_achievement"`
	Code          string     `json:"code"`
	Name          string     `json:"name"`
	Description   string     `json:"description"`
	Category      string     `json:"category"`
	Tier          string     `json:"tier"` // 'BRONZE', 'SILVER', 'GOLD', 'PLATINUM', 'DIAMOND'
	IconName      string     `json:"icon_name"`
	XPReward      int        `json:"xp_reward"`
	IsUnlocked    bool       `json:"is_unlocked"`
	UnlockedAt    *time.Time `json:"unlocked_at,omitempty"`
}

type MasteryDashboard struct {
	TotalWordsStudied int                 `json:"total_words_studied"`
	DueTodayCount     int                 `json:"due_today_count"`
	CountNew          int                 `json:"count_new"`
	CountLearning     int                 `json:"count_learning"`
	CountReviewing    int                 `json:"count_reviewing"`
	CountMastered     int                 `json:"count_mastered"`
	GlobalRetention   float32             `json:"global_retention"`
	TotalXP           int                 `json:"total_xp"`
	CurrentLevel      int                 `json:"current_level"`
	CurrentStreak     int                 `json:"current_streak"`
	LongestStreak     int                 `json:"longest_streak"`
	RecentXP          []XPTransactionItem `json:"recent_xp_transactions"`
	Achievements      []AchievementStatus `json:"achievements"`
}
