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
	Modality      string   `json:"modality"` // 'DEFINITION_TO_WORD', 'WORD_TO_DEFINITION', 'WORD_TO_IMAGE', etc.
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

// Solicitud del cliente a core-service para preparar una nueva sala (CU14)
type CreateMatchRequest struct {
	IDDeck       int   `json:"id_deck"`             // Obligatorio: Baraja base (mín. 8 sentidos APPROVED y >= TotalRounds)
	TotalRounds  int   `json:"total_rounds"`        // 8..20 (default 8)
	RoundTimeSec int   `json:"round_time_sec"`      // 10..30 segundos por ronda (default 15)
	IsPublic     *bool `json:"is_public,omitempty"` // Opcional (default true)
}

// Reactivo interno completo (incluye CorrectOption para calificación en el servidor y persistencia en match_rounds;
// cuando realtime-service emite el evento WS "QUESTION" a los clientes, envía ArenaQuestion omitiendo CorrectOption)
type InternalArenaQuestion struct {
	RoundNumber   int      `json:"round_number"`
	IDSense       int      `json:"id_sense"`
	IDQuestion    *int     `json:"id_question,omitempty"`
	Lemma         string   `json:"lemma,omitempty"`
	Modality      string   `json:"modality"` // 'DEFINITION_TO_WORD', 'WORD_TO_DEFINITION', 'WORD_TO_IMAGE', etc.
	Prompt        string   `json:"prompt"`
	Options       []string `json:"options"`        // 4 opciones barajadas (option_a..option_d: 1 correcta + 3 distractores)
	CorrectOption string   `json:"correct_option"` // NUNCA enviado en el evento QUESTION, revelado en ROUND_RESULT
	Explanation   string   `json:"explanation"`
	TimeLimitSecs int      `json:"time_limit_secs"`
}

// Paquete devuelto al crear/preparar una sala en core-service y al hidratar realtime-service
type MatchPackage struct {
	IDMatch      string                  `json:"id_match"`
	RoomCode     string                  `json:"room_code"` // PIN único de 6 caracteres, ej. "MX8921"
	HostUserID   string                  `json:"host_user_id"`
	IDDeck       int                     `json:"id_deck"`
	DeckTitle    string                  `json:"deck_title"`
	IsPublic     bool                    `json:"is_public"`
	TotalRounds  int                     `json:"total_rounds"`
	RoundTimeSec int                     `json:"round_time_sec"`
	Questions    []InternalArenaQuestion `json:"questions"`
}

// Trama enviada por el jugador en el evento WS "SUBMIT_ANSWER"
type SubmitAnswerPayload struct {
	RoundNumber     int    `json:"round_number"`
	SelectedOption  string `json:"selected_option"`
	ClientElapsedMs int    `json:"client_elapsed_ms"` // Solo telemetría; el servidor califica con su reloj monotónico
}

// Trama difundida por el servidor al cerrar una ronda en el evento WS "ROUND_RESULT"
type RoundResultPayload struct {
	RoundNumber   int                `json:"round_number"`
	IDSense       int                `json:"id_sense"`
	CorrectOption string             `json:"correct_option"`
	Explanation   string             `json:"explanation"`
	Leaderboard   []ArenaParticipant `json:"leaderboard"` // Ordenado por Score DESC, CorrectAnswers DESC, TotalTimeSpentMs ASC
}

// Detalle de respuesta de un participante en una ronda específica para persistir en quiz_attempts (CU16 Replay)
type ParticipantRoundAttempt struct {
	RoundNumber    int    `json:"round_number"`
	IDSense        int    `json:"id_sense"`
	IDQuestion     *int   `json:"id_question,omitempty"`
	SelectedOption string `json:"selected_option"` // "" si no respondió (timeout)
	IsCorrect      bool   `json:"is_correct"`
	TimeSpentMs    int    `json:"time_spent_ms"` // ElapsedMs del servidor (o RoundTimeSec*1000 en timeout)
	PointsEarned   int    `json:"points_earned"`
}

// Payload de Batch Flush enviado de realtime-service a core-service al concluir la partida (CU16)
type ParticipantFinalResult struct {
	IDUser           string                    `json:"id_user"`
	FinalScore       int                       `json:"final_score"`
	AccuracyRate     float32                   `json:"accuracy_rate"`
	FinalRank        int                       `json:"final_rank"`
	TotalTimeSpentMs int                       `json:"total_time_spent_ms"`
	CorrectSenseIDs  []int                     `json:"correct_sense_ids"`
	WrongSenseIDs    []int                     `json:"wrong_sense_ids"`
	RoundAttempts    []ParticipantRoundAttempt `json:"round_attempts"`
}

type ConsolidateMatchRequest struct {
	IDMatch      string                   `json:"id_match"`
	RoomCode     string                   `json:"room_code"`
	Participants []ParticipantFinalResult `json:"participants"`
}

type MatchSummary struct {
	IDMatch      string             `json:"id_match"`
	RoomCode     string             `json:"room_code"`
	IDDeck       int                `json:"id_deck"`
	DeckTitle    string             `json:"deck_title"`
	WinnerUserID *string            `json:"winner_user_id,omitempty"`
	Status       MatchStatus        `json:"status"`
	TotalRounds  int                `json:"total_rounds"`
	RoundTimeSec int                `json:"round_time_sec"`
	Podium       []ArenaParticipant `json:"podium"`
	FinishedAt   string             `json:"finished_at"`
}

// --- DTOs DE RETROSPECCIÓN Y REPLAY GRANULAR RONDA POR RONDA (CU16 Replay 3NF) ---

type ParticipantRoundAnswer struct {
	IDUser         string `json:"id_user"`
	UserName       string `json:"user_name"`
	Name           string `json:"name"`
	AvatarURL      string `json:"avatar_url,omitempty"`
	DialectCode    string `json:"dialect_code"`
	FlagEmoji      string `json:"flag_emoji"`
	SelectedOption string `json:"selected_option"`
	IsCorrect      bool   `json:"is_correct"`
	TimeSpentMs    int    `json:"time_spent_ms"`
	PointsEarned   int    `json:"points_earned"`
}

type MatchRoundReplay struct {
	RoundNumber   int                      `json:"round_number"`
	IDSense       int                      `json:"id_sense"`
	Lemma         string                   `json:"lemma"`
	Modality      string                   `json:"modality"`
	Prompt        string                   `json:"prompt"`
	Options       []string                 `json:"options"`
	CorrectOption string                   `json:"correct_option"`
	Explanation   string                   `json:"explanation"`
	Answers       []ParticipantRoundAnswer `json:"answers"`
}

type MatchReplayDetail struct {
	MatchSummary
	Rounds []MatchRoundReplay `json:"rounds"`
}
