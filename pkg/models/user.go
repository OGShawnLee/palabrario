package models

import "time"

type AccountRole string

const (
	RoleLearner AccountRole = "LEARNER"
	RoleEditor  AccountRole = "EDITOR"
	RoleAdmin   AccountRole = "ADMIN"
)

type Region struct {
	IDRegion    string  `json:"id_region"`
	Code        string  `json:"code"`         // 'MX', 'ES', 'AR', 'CO', 'CL', 'CU', 'CAM', 'PAN'
	Name        string  `json:"name"`         // 'México', 'España'
	ParentZone  string  `json:"parent_zone"`  // 'Norteamérica', 'Rioplatense'
	CountryCode *string `json:"country_code"` // 'MX', 'ES'
	FlagEmoji   string  `json:"flag_emoji"`   // '🇲🇽', '🇪🇸'
}

type UserProfile struct {
	IDUser        string      `json:"id_user"`
	Email         string      `json:"email"`
	Name          string      `json:"name"`
	UserName      string      `json:"user_name"`
	AvatarURL     *string     `json:"avatar_url,omitempty"`
	AnchorDialect string      `json:"anchor_dialect"`
	TotalXP       int         `json:"total_xp"`
	CurrentLevel  int         `json:"current_level"`
	CurrentStreak int         `json:"current_streak"`
	LongestStreak int         `json:"longest_streak"`
	Role          AccountRole `json:"role"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
}

type UpdateProfileRequest struct {
	Name          string  `json:"name"`                     // Obligatorio, 2..100 chars
	AnchorDialect string  `json:"anchor_dialect"`           // Obligatorio: 'MX','ES','AR','CO','CL','CU','CAM','PAN'
	AvatarURL     *string `json:"avatar_url,omitempty"`     // Opcional, URL generada por image-service
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password"`        // Obligatorio
	NewPassword     string `json:"new_password"`            // Obligatorio, mín. 8 caracteres
}

type PublicUserProfile struct {
	IDUser            string              `json:"id_user"`
	Name              string              `json:"name"`
	UserName          string              `json:"user_name"`
	AvatarURL         *string             `json:"avatar_url,omitempty"`
	AnchorDialect     string              `json:"anchor_dialect"`
	Role              AccountRole         `json:"role"`
	TotalXP           int                 `json:"total_xp"`
	CurrentLevel      int                 `json:"current_level"`
	CurrentStreak     int                 `json:"current_streak"`
	LongestStreak     int                 `json:"longest_streak"`
	UnlockedBadges    []AchievementStatus `json:"unlocked_badges"`
	PublicDecks       []Deck              `json:"public_decks"`
	PublicCollections []Collection        `json:"public_collections"`
	MatchesWon        int                 `json:"matches_won"`
	PodiumCount       int                 `json:"podium_count"`
	IsOwner           bool                `json:"is_owner"`
	// Campos poblados ÚNICAMENTE cuando IsOwner == true:
	Email             *string             `json:"email,omitempty"`
	PrivateDecksCount *int                `json:"private_decks_count,omitempty"`
	DueReviewsToday   *int                `json:"due_reviews_today,omitempty"`
	CreatedAt         time.Time           `json:"created_at"`
}
