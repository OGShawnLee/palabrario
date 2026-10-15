package models

type MatchStatus string

const (
	MatchStatusWaiting    MatchStatus = "WAITING"
	MatchStatusInProgress MatchStatus = "IN_PROGRESS"
	MatchStatusFinished   MatchStatus = "FINISHED"
)

type ArenaParticipant struct {
	IDUser      string `json:"id_user"`
	UserName    string `json:"user_name"`
	Name        string `json:"name"`
	AvatarURL   string `json:"avatar_url"`
	Score       int    `json:"score"`
	Rank        int    `json:"rank"`
	IsHost      bool   `json:"is_host"`
	IsReady     bool   `json:"is_ready"`
	DialectCode string `json:"dialect_code"`
	FlagEmoji   string `json:"flag_emoji"`
}

type ArenaQuestion struct {
	RoundNumber   int      `json:"round_number"`
	Prompt        string   `json:"prompt"`
	Options       []string `json:"options"`
	Explanation   string   `json:"explanation"`
	TimeLimitSecs int      `json:"time_limit_secs"`
}

type ArenaRoomState struct {
	RoomCode     string             `json:"room_code"`
	DeckID       int                `json:"deck_id"`
	DeckTitle    string             `json:"deck_title"`
	Status       MatchStatus        `json:"status"`
	TotalRounds  int                `json:"total_rounds"`
	CurrentRound int                `json:"current_round"`
	RoundTimeSec int                `json:"round_time_sec"`
	Participants []ArenaParticipant `json:"participants"`
}

// WebSocket Event Wrapper
type WSEvent struct {
	Type    string      `json:"type"` // e.g. "ROOM_STATE", "COUNTDOWN", "QUESTION", "ROUND_RESULT", "PODIUM"
	Payload interface{} `json:"payload,omitempty"`
}
