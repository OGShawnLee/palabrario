package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

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

	dialect := strings.TrimSpace(strings.ToUpper(r.URL.Query().Get("dialect")))
	if dialect == "" {
		if claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims); ok && claims != nil && claims.AnchorDialect != "" {
			dialect = claims.AnchorDialect
		}
		if dialect == "" {
			dialect = "MX"
		}
	}

	limit := queryIntOrDefault(r, "limit", 10, 1, 0)

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
	wordID, ok := parsePositiveIntParam(w, r, "id", "El identificador de palabra debe ser un número entero positivo")
	if !ok {
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

// ListRegions maneja GET /api/v1/regions
func (h *DictionaryHandler) ListRegions(w http.ResponseWriter, r *http.Request) {
	regions, err := h.dictService.ListRegions(r.Context())
	if err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al listar regiones")
		response.Error(w, status, msg)
		return
	}
	response.JSON(w, http.StatusOK, regions)
}

// GetWordOfTheDay maneja GET /api/v1/words/word-of-the-day?dialect={MX}
func (h *DictionaryHandler) GetWordOfTheDay(w http.ResponseWriter, r *http.Request) {
	dialect := strings.TrimSpace(r.URL.Query().Get("dialect"))
	if dialect == "" {
		if claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims); ok && claims != nil && claims.AnchorDialect != "" {
			dialect = claims.AnchorDialect
		}
	}

	word, err := h.dictService.GetWordOfTheDay(r.Context(), dialect)
	if err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al obtener la palabra del día")
		response.Error(w, status, msg)
		return
	}
	response.JSON(w, http.StatusOK, word)
}

// ExploreWords maneja GET /api/v1/words
func (h *DictionaryHandler) ExploreWords(w http.ResponseWriter, r *http.Request) {
	page := queryIntOrDefault(r, "page", 1, 1, 0)
	limit := queryIntOrDefault(r, "limit", 20, 1, 50)

	includeTaboo := true
	if tabStr := r.URL.Query().Get("include_taboo"); tabStr != "" {
		if parsed, err := strconv.ParseBool(tabStr); err == nil {
			includeTaboo = parsed
		}
	}

	filter := models.WordExploreFilter{
		Query:          r.URL.Query().Get("q"),
		RegionCode:     r.URL.Query().Get("region"),
		PartOfSpeech:   r.URL.Query().Get("pos"),
		CEFRLevel:      r.URL.Query().Get("cefr"),
		UsageFrequency: r.URL.Query().Get("frequency"),
		IncludeTaboo:   includeTaboo,
		Page:           page,
		Limit:          limit,
	}

	paginated, err := h.dictService.ExploreWords(r.Context(), filter)
	if err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al explorar palabras")
		response.Error(w, status, msg)
		return
	}
	response.JSON(w, http.StatusOK, paginated)
}

// CreateWord maneja POST /api/v1/words
func (h *DictionaryHandler) CreateWord(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	if !ok || claims == nil {
		response.Error(w, http.StatusUnauthorized, "No autenticado")
		return
	}

	var req models.CreateWordProposalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Cuerpo JSON inválido")
		return
	}

	detail, err := h.dictService.CreateWordProposal(r.Context(), claims.UserID, claims.Role, req)
	if err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al registrar propuesta de palabra")
		response.Error(w, status, msg)
		return
	}

	response.JSON(w, http.StatusCreated, detail)
}

// UpdateWord maneja PUT /api/v1/words/{id}
func (h *DictionaryHandler) UpdateWord(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	if !ok || claims == nil {
		response.Error(w, http.StatusUnauthorized, "No autenticado")
		return
	}

	wordID, valid := parsePositiveIntParam(w, r, "id", "El identificador de palabra debe ser un número entero positivo")
	if !valid {
		return
	}

	var req models.UpdateWordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Cuerpo JSON inválido")
		return
	}

	detail, err := h.dictService.UpdateWord(r.Context(), claims.UserID, claims.Role, wordID, req)
	if err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al actualizar metadatos de la palabra")
		response.Error(w, status, msg)
		return
	}

	response.JSON(w, http.StatusOK, detail)
}

// CreateSense maneja POST /api/v1/words/{id}/senses
func (h *DictionaryHandler) CreateSense(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	if !ok || claims == nil {
		response.Error(w, http.StatusUnauthorized, "No autenticado")
		return
	}

	wordID, valid := parsePositiveIntParam(w, r, "id", "El identificador de palabra debe ser un número entero positivo")
	if !valid {
		return
	}

	var req models.CreateSenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Cuerpo JSON inválido")
		return
	}

	sense, err := h.dictService.CreateSense(r.Context(), claims.UserID, claims.Role, wordID, req)
	if err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al crear acepción")
		response.Error(w, status, msg)
		return
	}

	response.JSON(w, http.StatusCreated, sense)
}

// UpdateSense maneja PUT /api/v1/senses/{senseId}
func (h *DictionaryHandler) UpdateSense(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	if !ok || claims == nil {
		response.Error(w, http.StatusUnauthorized, "No autenticado")
		return
	}

	senseID, valid := parsePositiveIntParam(w, r, "senseId", "El identificador de acepción debe ser un número entero positivo")
	if !valid {
		return
	}

	var req models.UpdateSenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Cuerpo JSON inválido")
		return
	}

	sense, err := h.dictService.UpdateSense(r.Context(), claims.UserID, claims.Role, senseID, req)
	if err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al actualizar acepción")
		response.Error(w, status, msg)
		return
	}

	response.JSON(w, http.StatusOK, sense)
}

// DeleteSense maneja DELETE /api/v1/senses/{senseId}
func (h *DictionaryHandler) DeleteSense(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	if !ok || claims == nil {
		response.Error(w, http.StatusUnauthorized, "No autenticado")
		return
	}

	senseID, valid := parsePositiveIntParam(w, r, "senseId", "El identificador de acepción debe ser un número entero positivo")
	if !valid {
		return
	}

	if err := h.dictService.DeleteSense(r.Context(), claims.UserID, claims.Role, senseID); err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al eliminar acepción")
		response.Error(w, status, msg)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// CreateExample maneja POST /api/v1/senses/{senseId}/examples
func (h *DictionaryHandler) CreateExample(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	if !ok || claims == nil {
		response.Error(w, http.StatusUnauthorized, "No autenticado")
		return
	}

	senseID, valid := parsePositiveIntParam(w, r, "senseId", "El identificador de acepción debe ser un número entero positivo")
	if !valid {
		return
	}

	var req models.CreateExampleInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Cuerpo JSON inválido")
		return
	}

	ex, err := h.dictService.CreateExample(r.Context(), claims.UserID, claims.Role, senseID, req)
	if err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al crear ejemplo")
		response.Error(w, status, msg)
		return
	}

	response.JSON(w, http.StatusCreated, ex)
}

// UpdateExample maneja PUT /api/v1/examples/{exampleId}
func (h *DictionaryHandler) UpdateExample(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	if !ok || claims == nil {
		response.Error(w, http.StatusUnauthorized, "No autenticado")
		return
	}

	exampleID, valid := parsePositiveIntParam(w, r, "exampleId", "El identificador de ejemplo debe ser un número entero positivo")
	if !valid {
		return
	}

	var req models.CreateExampleInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Cuerpo JSON inválido")
		return
	}

	ex, err := h.dictService.UpdateExample(r.Context(), claims.UserID, claims.Role, exampleID, req)
	if err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al actualizar ejemplo")
		response.Error(w, status, msg)
		return
	}

	response.JSON(w, http.StatusOK, ex)
}

// DeleteExample maneja DELETE /api/v1/examples/{exampleId}
func (h *DictionaryHandler) DeleteExample(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	if !ok || claims == nil {
		response.Error(w, http.StatusUnauthorized, "No autenticado")
		return
	}

	exampleID, valid := parsePositiveIntParam(w, r, "exampleId", "El identificador de ejemplo debe ser un número entero positivo")
	if !valid {
		return
	}

	if err := h.dictService.DeleteExample(r.Context(), claims.UserID, claims.Role, exampleID); err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al eliminar ejemplo")
		response.Error(w, status, msg)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListConcepts maneja GET /api/v1/concepts?q=
func (h *DictionaryHandler) ListConcepts(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	concepts, err := h.dictService.ListSemanticConcepts(r.Context(), q)
	if err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al listar conceptos semánticos")
		response.Error(w, status, msg)
		return
	}
	response.JSON(w, http.StatusOK, concepts)
}

// CreateConcept maneja POST /api/v1/concepts
func (h *DictionaryHandler) CreateConcept(w http.ResponseWriter, r *http.Request) {
	var req models.CreateConceptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Cuerpo JSON inválido")
		return
	}

	concept, err := h.dictService.CreateSemanticConcept(r.Context(), req)
	if err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al crear concepto semántico")
		response.Error(w, status, msg)
		return
	}

	response.JSON(w, http.StatusCreated, concept)
}

// LinkConcept maneja POST /api/v1/senses/{senseId}/concepts/{conceptId}
func (h *DictionaryHandler) LinkConcept(w http.ResponseWriter, r *http.Request) {
	senseID, ok := parsePositiveIntParam(w, r, "senseId", "senseId inválido")
	if !ok {
		return
	}
	conceptID, ok := parsePositiveIntParam(w, r, "conceptId", "conceptId inválido")
	if !ok {
		return
	}

	if err := h.dictService.LinkSenseToConcept(r.Context(), senseID, conceptID); err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al vincular concepto")
		response.Error(w, status, msg)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UnlinkConcept maneja DELETE /api/v1/senses/{senseId}/concepts/{conceptId}
func (h *DictionaryHandler) UnlinkConcept(w http.ResponseWriter, r *http.Request) {
	senseID, ok := parsePositiveIntParam(w, r, "senseId", "senseId inválido")
	if !ok {
		return
	}
	conceptID, ok := parsePositiveIntParam(w, r, "conceptId", "conceptId inválido")
	if !ok {
		return
	}

	if err := h.dictService.UnlinkSenseFromConcept(r.Context(), senseID, conceptID); err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al desvincular concepto")
		response.Error(w, status, msg)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
