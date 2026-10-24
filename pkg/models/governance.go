package models

import "time"

type RequestStatus string

const (
	RequestPending  RequestStatus = "PENDING"
	RequestApproved RequestStatus = "APPROVED"
	RequestRejected RequestStatus = "REJECTED"
)

type CreateEditorRequestDTO struct {
	InstitutionName string  `json:"institution_name"`          // Obligatorio, 3..150 chars
	Justification   string  `json:"justification"`             // Obligatorio, mín. 20 chars
	CredentialsURL  *string `json:"credentials_url,omitempty"` // Opcional, URL del comprobante
}

type ResolveEditorRequestDTO struct {
	Status      RequestStatus `json:"status"`                 // 'APPROVED' o 'REJECTED'
	ReviewNotes *string       `json:"review_notes,omitempty"` // Motivo del dictamen administrativo
}

type EditorRequest struct {
	IDRequest       string        `json:"id_request"`
	IDUser          string        `json:"id_user"`
	UserName        string        `json:"user_name"`
	UserFullName    string        `json:"user_full_name"`
	UserEmail       string        `json:"user_email"`
	InstitutionName string        `json:"institution_name"`
	Justification   string        `json:"justification"`
	CredentialsURL  *string       `json:"credentials_url,omitempty"`
	Status          RequestStatus `json:"status"`
	ReviewedBy      *string       `json:"reviewed_by,omitempty"`
	ReviewNotes     *string       `json:"review_notes,omitempty"`
	CreatedAt       time.Time     `json:"created_at"`
	ResolvedAt      *time.Time    `json:"resolved_at,omitempty"`
}

type UpdateUserRoleOrStatusRequest struct {
	Role     *AccountRole `json:"role,omitempty"`      // 'LEARNER', 'EDITOR', 'ADMIN'
	IsActive *bool        `json:"is_active,omitempty"` // Suspensión / Reactivación de cuenta
}

type CreateTicketRequest struct {
	Code            *string     `json:"code,omitempty"`             // Si se omite, el backend genera ej. "PAL-UNAM-2026-X9A2"
	GrantedRole     AccountRole `json:"granted_role"`               // 'LEARNER' o 'EDITOR' (default 'LEARNER')
	InstitutionName *string     `json:"institution_name,omitempty"` // Ej. "Facultad de Filosofía y Letras UNAM"
	MaxUses         int         `json:"max_uses"`                   // >= 1 (1 = uniuso, >1 = multiuso)
	ExpiresAt       *time.Time  `json:"expires_at,omitempty"`       // Opcional
}

type InvitationTicket struct {
	IDTicket        string      `json:"id_ticket"`
	Code            string      `json:"code"`
	CreatedBy       string      `json:"created_by"`
	GrantedRole     AccountRole `json:"granted_role"`
	InstitutionName *string     `json:"institution_name,omitempty"`
	MaxUses         int         `json:"max_uses"`
	CurrentUses     int         `json:"current_uses"`
	IsActive        bool        `json:"is_active"`
	ExpiresAt       *time.Time  `json:"expires_at,omitempty"`
	CreatedAt       time.Time   `json:"created_at"`
}

type RedeemTicketRequest struct {
	Code string `json:"code"` // Obligatorio
}
