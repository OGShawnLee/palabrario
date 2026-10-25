package models

import "time"

type ReviewDecision string

const (
	DecisionApprove        ReviewDecision = "APPROVE"
	DecisionRequestChanges ReviewDecision = "REQUEST_CHANGES"
	DecisionReject         ReviewDecision = "REJECT"
)

type SenseApprovalVote struct {
	IDSense        int            `json:"id_sense"`
	IDEditor       string         `json:"id_editor"`
	EditorUserName string         `json:"editor_user_name"`
	EditorFullName string         `json:"editor_full_name"`
	EditorDialect  string         `json:"editor_dialect"`
	Decision       ReviewDecision `json:"decision"`
	FeedbackNotes  *string        `json:"feedback_notes,omitempty"`
	ApprovedAt     time.Time      `json:"approved_at"`
}

type SenseQuorumItem struct {
	IDSense         int       `json:"id_sense"`
	IDWord          int       `json:"id_word"`
	Lemma           string    `json:"lemma"`
	PartOfSpeech    string    `json:"part_of_speech"`
	RegionCode      string    `json:"region_code"`
	RegionName      string    `json:"region_name"`
	FlagEmoji       string    `json:"flag_emoji"`
	ShortDefinition string    `json:"short_definition"`
	Definition      string    `json:"definition"`
	CulturalNote    *string   `json:"cultural_note,omitempty"`
	ImageURL        *string   `json:"image_url,omitempty"`
	CEFRLevel       string    `json:"cefr_level"`
	UsageFrequency  string    `json:"usage_frequency"`
	IsTabooOrVulgar bool      `json:"is_taboo_or_vulgar"`
	State           string    `json:"state"` // 'DRAFT', 'APPROVED', 'HIDDEN'
	AuthorID        string    `json:"author_id"`
	AuthorUserName  string    `json:"author_user_name"`
	AuthorFullName  string    `json:"author_full_name"`
	ApprovalsCount  int       `json:"approvals_count"` // Votos con decision = 'APPROVE'
	RequiredQuorum  int       `json:"required_quorum"` // Constante: 3 (RN06)
	CreatedAt       time.Time `json:"created_at"`
}

// PendingSenseItem es un alias canónico para los ítems de la cola de curaduría
type PendingSenseItem = SenseQuorumItem

type SenseQuorumDetail struct {
	Proposal SenseQuorumItem     `json:"proposal"`
	Examples []Example           `json:"examples"`
	Votes    []SenseApprovalVote `json:"votes"`
}

type SubmitVoteRequest struct {
	Decision      ReviewDecision `json:"decision"`                 // 'APPROVE', 'REQUEST_CHANGES', 'REJECT'
	FeedbackNotes *string        `json:"feedback_notes,omitempty"` // Obligatorio si Decision != 'APPROVE'
}

type AdminOverrideStateRequest struct {
	State  string  `json:"state"`            // 'APPROVED' o 'HIDDEN'
	Reason *string `json:"reason,omitempty"` // Nota de auditoría
}

// --- DTOs DE MODERACIÓN DE REPORTES (CU09) ---

type ContentReport struct {
	IDReport        string     `json:"id_report"`
	IDReporter      string     `json:"id_reporter"`
	ReporterName    string     `json:"reporter_user_name"`
	TargetType      string     `json:"target_type"` // 'WORD', 'SENSE', 'EXAMPLE'
	IDWord          *int       `json:"id_word,omitempty"`
	IDSense         *int       `json:"id_sense,omitempty"`
	IDExample       *int       `json:"id_example,omitempty"`
	TargetPreview   string     `json:"target_preview"` // Lema o texto reportado para inspección rápida
	Reason          string     `json:"reason"`         // 'INACCURATE_DEFINITION', 'OFFENSIVE_UNMARKED', 'SPAM_OR_VANDALISM', 'DUPLICATE', 'OTHER'
	Details         string     `json:"details"`
	Status          string     `json:"status"` // 'PENDING', 'RESOLVED_HIDDEN', 'RESOLVED_AMENDED', 'DISMISSED'
	ResolvedBy      *string    `json:"resolved_by,omitempty"`
	ResolutionNotes *string    `json:"resolution_notes,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	ResolvedAt      *time.Time `json:"resolved_at,omitempty"`
}

type CreateReportRequest struct {
	TargetType string `json:"target_type"` // 'WORD', 'SENSE', 'EXAMPLE'
	IDWord     *int   `json:"id_word,omitempty"`
	IDSense    *int   `json:"id_sense,omitempty"`
	IDExample  *int   `json:"id_example,omitempty"`
	Reason     string `json:"reason"`
	Details    string `json:"details"` // Obligatorio, mín. 10 chars
}

type ResolveReportRequest struct {
	Status          string  `json:"status"` // 'RESOLVED_HIDDEN', 'RESOLVED_AMENDED', 'DISMISSED'
	ResolutionNotes *string `json:"resolution_notes,omitempty"`
}
