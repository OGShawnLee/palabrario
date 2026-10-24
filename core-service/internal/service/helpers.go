package service

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// StringToText convierte un string a pgtype.Text aplicando TrimSpace.
// Si la cadena resultante está vacía, retorna Valid: false.
func StringToText(s string) pgtype.Text {
	trimmed := strings.TrimSpace(s)
	if trimmed == "" {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: trimmed, Valid: true}
}

// PtrToText convierte un *string anulable a pgtype.Text saneando espacios en blanco.
func PtrToText(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return StringToText(*s)
}

// TextToPtr convierte un pgtype.Text anulable a *string para los DTOs de pkg/models.
func TextToPtr(t pgtype.Text) *string {
	if !t.Valid {
		return nil
	}
	s := t.String
	return &s
}

// UUIDToPgtype convierte un uuid.UUID a pgtype.UUID válido.
func UUIDToPgtype(u uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: u, Valid: true}
}

// ParseDomainUUID valida y parsea un string UUID o retorna un DomainError(ErrInvalidInput, errMsg).
func ParseDomainUUID(raw string, errMsg string) (uuid.UUID, error) {
	parsed, err := uuid.Parse(strings.TrimSpace(raw))
	if err != nil {
		return uuid.Nil, NewDomainError(ErrInvalidInput, errMsg)
	}
	return parsed, nil
}

// InterfaceToString convierte valores dinámicos de columnas SQL/JSONB a string.
func InterfaceToString(val any, defaultVal string) string {
	if val == nil {
		return defaultVal
	}
	switch v := val.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	default:
		return fmt.Sprintf("%v", v)
	}
}
