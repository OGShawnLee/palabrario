package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"palabrario/core-service/internal/auth"
	"palabrario/core-service/internal/repository"
)

type DictionaryHandler struct {
	repo *repository.DictionaryRepository
}

func NewDictionaryHandler(repo *repository.DictionaryRepository) *DictionaryHandler {
	return &DictionaryHandler{repo: repo}
}

// Search maneja GET /api/v1/words/search?q={query}&dialect={MX}&limit={10}
func (h *DictionaryHandler) Search(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.URL.Query().Get("q"))
	if query == "" {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]interface{}{})
		return
	}

	// Si el usuario está autenticado, intentar usar su dialecto ancla por defecto
	dialect := strings.TrimSpace(r.URL.Query().Get("dialect"))
	if dialect == "" {
		if claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims); ok && claims != nil {
			// Si hay contexto de auth, se puede enriquecer aquí
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

	results, err := h.repo.SearchWords(r.Context(), query, dialect, limit)
	if err != nil {
		http.Error(w, `{"error":"Error al procesar la búsqueda de vocabulario"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(results)
}

// GetByID maneja GET /api/v1/words/{id}
func (h *DictionaryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	wordID, err := strconv.Atoi(idStr)
	if err != nil || wordID <= 0 {
		http.Error(w, `{"error":"El identificador de palabra debe ser un número entero positivo"}`, http.StatusBadRequest)
		return
	}

	detail, err := h.repo.GetWordDetail(r.Context(), wordID)
	if err != nil {
		if errors.Is(err, repository.ErrWordNotFound) {
			http.Error(w, `{"error":"Palabra no encontrada en el diccionario"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error":"Error al obtener el detalle de la palabra"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(detail)
}
