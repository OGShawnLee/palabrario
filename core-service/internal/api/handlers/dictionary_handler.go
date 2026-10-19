package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"palabrario/core-service/internal/service"
	"palabrario/pkg/auth"
	"palabrario/pkg/models"
	"palabrario/pkg/response"
)

type DictionaryHandler struct {
	dictService service.DictionaryService
}

func NewDictionaryHandler(dictService service.DictionaryService) *DictionaryHandler {
	return &DictionaryHandler{dictService: dictService}
}

// Search maneja GET /api/v1/words/search?q={query}&dialect={MX}&limit={10}
func (h *DictionaryHandler) Search(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.URL.Query().Get("q"))
	if query == "" {
		response.JSON(w, http.StatusOK, []models.WordSearchResult{})
		return
	}

	// Si no se especifica dialecto en query string, priorizar el dialecto ancla del JWT (vía OptionalAuth)
	dialect := strings.TrimSpace(strings.ToUpper(r.URL.Query().Get("dialect")))
	if dialect == "" {
		if claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims); ok && claims != nil && claims.AnchorDialect != "" {
			dialect = claims.AnchorDialect
		}
		if dialect == "" {
			dialect = "MX"
		}
	}

	limit := 10
	if limitStr := strings.TrimSpace(r.URL.Query().Get("limit")); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	results, err := h.dictService.SearchDropdown(r.Context(), query, dialect, limit)
	if err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al procesar la búsqueda de vocabulario")
		response.Error(w, status, msg)
		return
	}

	response.JSON(w, http.StatusOK, results)
}

// GetByID maneja GET /api/v1/words/{id}
func (h *DictionaryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	wordID, err := strconv.Atoi(idStr)
	if err != nil || wordID <= 0 {
		response.Error(w, http.StatusBadRequest, "El identificador de palabra debe ser un número entero positivo")
		return
	}

	detail, err := h.dictService.GetWordDetail(r.Context(), wordID)
	if err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al obtener el detalle de la palabra")
		response.Error(w, status, msg)
		return
	}

	response.JSON(w, http.StatusOK, detail)
}
