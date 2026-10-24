package utils

// ToPtr retorna un puntero al valor suministrado (útil para literales en DTOs y tests).
func ToPtr[T any](v T) *T {
	return &v
}
