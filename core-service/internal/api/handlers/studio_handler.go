package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"palabrario/core-service/internal/service"
	"palabrario/pkg/auth"
	"palabrario/pkg/models"
	"palabrario/pkg/response"
)

type StudioHandler struct {
	studioService service.StudioService
}

func NewStudioHandler(studioService service.StudioService) *StudioHandler {
	return &StudioHandler{studioService: studioService}
}

// -----------------------------------------------------------------------------
// BARAJAS (CU10)
// -----------------------------------------------------------------------------

// ListDecks maneja GET /api/v1/decks?q=&official_only=false
func (h *StudioHandler) ListDecks(w http.ResponseWriter, r *http.Request) {
	var query *string
	if q := strings.TrimSpace(r.URL.Query().Get("q")); q != "" {
		query = &q
	}
	var officialOnly *bool
	if off := r.URL.Query().Get("official_only"); off != "" {
		val := off == "true" || off == "1"
		officialOnly = &val
	}

	decks, err := h.studioService.ListPublicDecks(r.Context(), query, officialOnly)
	if err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al listar barajas públicas")
		response.Error(w, status, msg)
		return
	}
	response.JSON(w, http.StatusOK, decks)
}

// ListStudioDecks maneja GET /api/v1/studio/decks
func (h *StudioHandler) ListStudioDecks(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	if !ok || claims == nil {
		response.Error(w, http.StatusUnauthorized, "No autenticado")
		return
	}

	summaries, err := h.studioService.ListStudioDecks(r.Context(), claims.UserID)
	if err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al listar barajas del Studio")
		response.Error(w, status, msg)
		return
	}
	response.JSON(w, http.StatusOK, summaries)
}

// GetDeck maneja GET /api/v1/decks/{id} (OptionalAuth)
func (h *StudioHandler) GetDeck(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	deckID, err := strconv.Atoi(idStr)
	if err != nil || deckID <= 0 {
		response.Error(w, http.StatusBadRequest, "Identificador de baraja inválido")
		return
	}

	var viewerID *string
	var viewerRole *models.AccountRole
	if claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims); ok && claims != nil {
		viewerID = &claims.UserID
		viewerRole = &claims.Role
	}

	deck, err := h.studioService.GetDeckDetail(r.Context(), viewerID, viewerRole, deckID)
	if err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al obtener detalle de la baraja")
		response.Error(w, status, msg)
		return
	}
	response.JSON(w, http.StatusOK, deck)
}

// CreateDeck maneja POST /api/v1/decks
func (h *StudioHandler) CreateDeck(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	if !ok || claims == nil {
		response.Error(w, http.StatusUnauthorized, "No autenticado")
		return
	}

	var req models.CreateDeckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Cuerpo JSON inválido")
		return
	}

	deck, err := h.studioService.CreateDeck(r.Context(), claims.UserID, claims.Role, req)
	if err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al crear la baraja")
		response.Error(w, status, msg)
		return
	}
	response.JSON(w, http.StatusCreated, deck)
}

// UpdateDeck maneja PUT /api/v1/decks/{id}
func (h *StudioHandler) UpdateDeck(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	if !ok || claims == nil {
		response.Error(w, http.StatusUnauthorized, "No autenticado")
		return
	}

	idStr := chi.URLParam(r, "id")
	deckID, err := strconv.Atoi(idStr)
	if err != nil || deckID <= 0 {
		response.Error(w, http.StatusBadRequest, "Identificador de baraja inválido")
		return
	}

	var req models.UpdateDeckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Cuerpo JSON inválido")
		return
	}

	deck, err := h.studioService.UpdateDeck(r.Context(), claims.UserID, claims.Role, deckID, req)
	if err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al actualizar la baraja")
		response.Error(w, status, msg)
		return
	}
	response.JSON(w, http.StatusOK, deck)
}

// DeleteDeck maneja DELETE /api/v1/decks/{id}
func (h *StudioHandler) DeleteDeck(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	if !ok || claims == nil {
		response.Error(w, http.StatusUnauthorized, "No autenticado")
		return
	}

	idStr := chi.URLParam(r, "id")
	deckID, err := strconv.Atoi(idStr)
	if err != nil || deckID <= 0 {
		response.Error(w, http.StatusBadRequest, "Identificador de baraja inválido")
		return
	}

	if err := h.studioService.DeleteDeck(r.Context(), claims.UserID, claims.Role, deckID); err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al eliminar la baraja")
		response.Error(w, status, msg)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// AddDeckItem maneja POST /api/v1/decks/{id}/items
func (h *StudioHandler) AddDeckItem(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	if !ok || claims == nil {
		response.Error(w, http.StatusUnauthorized, "No autenticado")
		return
	}

	idStr := chi.URLParam(r, "id")
	deckID, err := strconv.Atoi(idStr)
	if err != nil || deckID <= 0 {
		response.Error(w, http.StatusBadRequest, "Identificador de baraja inválido")
		return
	}

	var req models.AddDeckItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Cuerpo JSON inválido")
		return
	}

	item, err := h.studioService.AddDeckItem(r.Context(), claims.UserID, claims.Role, deckID, req)
	if err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al agregar tarjeta a la baraja")
		response.Error(w, status, msg)
		return
	}
	response.JSON(w, http.StatusCreated, item)
}

// RemoveDeckItem maneja DELETE /api/v1/decks/{id}/items/{senseId}
func (h *StudioHandler) RemoveDeckItem(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	if !ok || claims == nil {
		response.Error(w, http.StatusUnauthorized, "No autenticado")
		return
	}

	idStr := chi.URLParam(r, "id")
	deckID, err := strconv.Atoi(idStr)
	if err != nil || deckID <= 0 {
		response.Error(w, http.StatusBadRequest, "Identificador de baraja inválido")
		return
	}

	senseStr := chi.URLParam(r, "senseId")
	senseID, err := strconv.Atoi(senseStr)
	if err != nil || senseID <= 0 {
		response.Error(w, http.StatusBadRequest, "Identificador de sentido inválido")
		return
	}

	if err := h.studioService.RemoveDeckItem(r.Context(), claims.UserID, claims.Role, deckID, senseID); err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al remover tarjeta de la baraja")
		response.Error(w, status, msg)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// SaveDeck maneja POST /api/v1/decks/{id}/save
func (h *StudioHandler) SaveDeck(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	if !ok || claims == nil {
		response.Error(w, http.StatusUnauthorized, "No autenticado")
		return
	}

	idStr := chi.URLParam(r, "id")
	deckID, err := strconv.Atoi(idStr)
	if err != nil || deckID <= 0 {
		response.Error(w, http.StatusBadRequest, "Identificador de baraja inválido")
		return
	}

	if err := h.studioService.SaveDeck(r.Context(), claims.UserID, claims.Role, deckID); err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al guardar baraja en Studio")
		response.Error(w, status, msg)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// UnsaveDeck maneja DELETE /api/v1/decks/{id}/save
func (h *StudioHandler) UnsaveDeck(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	if !ok || claims == nil {
		response.Error(w, http.StatusUnauthorized, "No autenticado")
		return
	}

	idStr := chi.URLParam(r, "id")
	deckID, err := strconv.Atoi(idStr)
	if err != nil || deckID <= 0 {
		response.Error(w, http.StatusBadRequest, "Identificador de baraja inválido")
		return
	}

	if err := h.studioService.UnsaveDeck(r.Context(), claims.UserID, deckID); err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al dejar de guardar baraja")
		response.Error(w, status, msg)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// -----------------------------------------------------------------------------
// COLECCIONES TEMÁTICAS (CU11)
// -----------------------------------------------------------------------------

// ListCollections maneja GET /api/v1/collections?featured=true
func (h *StudioHandler) ListCollections(w http.ResponseWriter, r *http.Request) {
	var featuredOnly *bool
	if feat := r.URL.Query().Get("featured"); feat != "" {
		val := feat == "true" || feat == "1"
		featuredOnly = &val
	}

	cols, err := h.studioService.ListPublicCollections(r.Context(), featuredOnly)
	if err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al listar colecciones públicas")
		response.Error(w, status, msg)
		return
	}
	response.JSON(w, http.StatusOK, cols)
}

// ListStudioCollections maneja GET /api/v1/studio/collections
func (h *StudioHandler) ListStudioCollections(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	if !ok || claims == nil {
		response.Error(w, http.StatusUnauthorized, "No autenticado")
		return
	}

	cols, err := h.studioService.ListStudioCollections(r.Context(), claims.UserID)
	if err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al listar colecciones del Studio")
		response.Error(w, status, msg)
		return
	}
	response.JSON(w, http.StatusOK, cols)
}

// GetCollection maneja GET /api/v1/collections/{id} (OptionalAuth)
func (h *StudioHandler) GetCollection(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	colID, err := strconv.Atoi(idStr)
	if err != nil || colID <= 0 {
		response.Error(w, http.StatusBadRequest, "Identificador de colección inválido")
		return
	}

	var viewerID *string
	var viewerRole *models.AccountRole
	if claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims); ok && claims != nil {
		viewerID = &claims.UserID
		viewerRole = &claims.Role
	}

	col, err := h.studioService.GetCollectionDetail(r.Context(), viewerID, viewerRole, colID)
	if err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al obtener detalle de la colección")
		response.Error(w, status, msg)
		return
	}
	response.JSON(w, http.StatusOK, col)
}

// CreateCollection maneja POST /api/v1/collections
func (h *StudioHandler) CreateCollection(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	if !ok || claims == nil {
		response.Error(w, http.StatusUnauthorized, "No autenticado")
		return
	}

	var req models.CreateCollectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Cuerpo JSON inválido")
		return
	}

	col, err := h.studioService.CreateCollection(r.Context(), claims.UserID, claims.Role, req)
	if err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al crear colección")
		response.Error(w, status, msg)
		return
	}
	response.JSON(w, http.StatusCreated, col)
}

// UpdateCollection maneja PUT /api/v1/collections/{id}
func (h *StudioHandler) UpdateCollection(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	if !ok || claims == nil {
		response.Error(w, http.StatusUnauthorized, "No autenticado")
		return
	}

	idStr := chi.URLParam(r, "id")
	colID, err := strconv.Atoi(idStr)
	if err != nil || colID <= 0 {
		response.Error(w, http.StatusBadRequest, "Identificador de colección inválido")
		return
	}

	var req models.UpdateCollectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Cuerpo JSON inválido")
		return
	}

	col, err := h.studioService.UpdateCollection(r.Context(), claims.UserID, claims.Role, colID, req)
	if err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al actualizar colección")
		response.Error(w, status, msg)
		return
	}
	response.JSON(w, http.StatusOK, col)
}

// DeleteCollection maneja DELETE /api/v1/collections/{id}
func (h *StudioHandler) DeleteCollection(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	if !ok || claims == nil {
		response.Error(w, http.StatusUnauthorized, "No autenticado")
		return
	}

	idStr := chi.URLParam(r, "id")
	colID, err := strconv.Atoi(idStr)
	if err != nil || colID <= 0 {
		response.Error(w, http.StatusBadRequest, "Identificador de colección inválido")
		return
	}

	if err := h.studioService.DeleteCollection(r.Context(), claims.UserID, claims.Role, colID); err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al eliminar colección")
		response.Error(w, status, msg)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// SaveCollection maneja POST /api/v1/collections/{id}/save
func (h *StudioHandler) SaveCollection(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	if !ok || claims == nil {
		response.Error(w, http.StatusUnauthorized, "No autenticado")
		return
	}

	idStr := chi.URLParam(r, "id")
	colID, err := strconv.Atoi(idStr)
	if err != nil || colID <= 0 {
		response.Error(w, http.StatusBadRequest, "Identificador de colección inválido")
		return
	}

	if err := h.studioService.SaveCollection(r.Context(), claims.UserID, claims.Role, colID); err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al guardar colección en Studio")
		response.Error(w, status, msg)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// UnsaveCollection maneja DELETE /api/v1/collections/{id}/save
func (h *StudioHandler) UnsaveCollection(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	if !ok || claims == nil {
		response.Error(w, http.StatusUnauthorized, "No autenticado")
		return
	}

	idStr := chi.URLParam(r, "id")
	colID, err := strconv.Atoi(idStr)
	if err != nil || colID <= 0 {
		response.Error(w, http.StatusBadRequest, "Identificador de colección inválido")
		return
	}

	if err := h.studioService.UnsaveCollection(r.Context(), claims.UserID, colID); err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al dejar de guardar colección")
		response.Error(w, status, msg)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
