package service

import (
	"fmt"
	"net/url"
	"path"
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

// ValidateInternalMediaURL valida conforme a la regla de negocio RN12 que la URL
// multimedia suministrada pertenezca al dominio seguro gestionado por image-service (/uploads/).
// Si rawURL es nil o está vacío, se considera válido (campo opcional).
// Rechaza esquemas maliciosos (javascript:, data:), traversals (..), fragmentos inyectados
// o dominios no autorizados que no apunten a /uploads/{category}/*.
func ValidateInternalMediaURL(rawURL *string) error {
	if rawURL == nil {
		return nil
	}
	u := strings.TrimSpace(*rawURL)
	if u == "" {
		return nil
	}

	lower := strings.ToLower(u)
	if strings.Contains(lower, "javascript:") || strings.Contains(lower, "data:") || strings.Contains(lower, "..") {
		return NewDomainError(ErrValidation, "RN12: solo se admiten URLs multimedia procesadas por image-service")
	}

	parsed, err := url.Parse(u)
	if err != nil {
		return NewDomainError(ErrValidation, "RN12: URL multimedia con formato inválido")
	}

	// Normalizar ruta URL
	cleanPath := path.Clean(parsed.Path)
	if !strings.HasPrefix(cleanPath, "/uploads/") {
		return NewDomainError(ErrValidation, "RN12: solo se admiten URLs multimedia procesadas por image-service")
	}

	// Verificar categoría permitida en el path (/uploads/{category}/...)
	parts := strings.Split(strings.TrimPrefix(cleanPath, "/uploads/"), "/")
	if len(parts) < 2 || parts[0] == "" || parts[1] == "" {
		return NewDomainError(ErrValidation, "RN12: la ruta de imagen debe incluir categoría y nombre de archivo")
	}

	validCategories := map[string]bool{
		"avatar":      true,
		"avatars":     true,
		"credentials": true,
		"word":        true,
		"sense":       true,
	}
	if !validCategories[parts[0]] {
		return NewDomainError(ErrValidation, "RN12: categoría de almacenamiento multimedia no permitida")
	}

	return nil
}
