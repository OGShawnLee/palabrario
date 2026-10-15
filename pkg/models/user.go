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
