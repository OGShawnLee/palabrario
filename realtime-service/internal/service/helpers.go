package service

import (
	"strings"
)

// DialectToFlag convierte un código de dialecto de 2 letras (ISO 3166-1 alpha-2, ej. "MX", "ES")
// en su respectivo emoji de bandera panhispánica utilizando Regional Indicator Symbols.
func DialectToFlag(code string) string {
	code = strings.ToUpper(strings.TrimSpace(code))
	if len(code) == 2 && code[0] >= 'A' && code[0] <= 'Z' && code[1] >= 'A' && code[1] <= 'Z' {
		r1 := rune(0x1F1E6 + int(code[0]-'A'))
		r2 := rune(0x1F1E6 + int(code[1]-'A'))
		return string([]rune{r1, r2})
	}
	return "🌐"
}
