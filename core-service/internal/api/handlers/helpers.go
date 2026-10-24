package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"palabrario/pkg/response"
)

// parsePositiveIntParam extrae un parámetro entero positivo de la ruta Chi.
// Si es inválido o <= 0, escribe automáticamente 400 Bad Request y retorna (0, false).
func parsePositiveIntParam(w http.ResponseWriter, r *http.Request, key string, errMsg string) (int, bool) {
	raw := chi.URLParam(r, key)
	val, err := strconv.Atoi(raw)
	if err != nil || val <= 0 {
		response.Error(w, http.StatusBadRequest, errMsg)
		return 0, false
	}
	return val, true
}

// queryIntOrDefault extrae un entero de los query params respetando límites [minVal, maxVal].
// Si maxVal <= 0, no aplica cota superior.
func queryIntOrDefault(r *http.Request, key string, defaultVal, minVal, maxVal int) int {
	raw := strings.TrimSpace(r.URL.Query().Get(key))
	if raw == "" {
		return defaultVal
	}
	val, err := strconv.Atoi(raw)
	if err != nil || val < minVal {
		return defaultVal
	}
	if maxVal > 0 && val > maxVal {
		return defaultVal
	}
	return val
}

// optionalQueryString retorna un *string si el query param no está vacío, o nil en caso contrario.
func optionalQueryString(r *http.Request, key string) *string {
	val := strings.TrimSpace(r.URL.Query().Get(key))
	if val == "" {
		return nil
	}
	return &val
}
