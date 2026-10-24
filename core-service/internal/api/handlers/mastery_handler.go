package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"palabrario/core-service/internal/service"
	"palabrario/pkg/auth"
	"palabrario/pkg/models"
	"palabrario/pkg/response"
)

type MasteryHandler struct {
	masteryService service.MasteryService
}

func NewMasteryHandler(masteryService service.MasteryService) *MasteryHandler {
	return &MasteryHandler{masteryService: masteryService}
}

// GetStudySession maneja GET /api/v1/studio/study-session?deck_id=1&limit=15
func (h *MasteryHandler) GetStudySession(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	if !ok || claims == nil {
		response.Error(w, http.StatusUnauthorized, "No autenticado")
		return
	}

	var deckID *int
	if dStr := r.URL.Query().Get("deck_id"); dStr != "" {
		id, err := strconv.Atoi(dStr)
		if err == nil && id > 0 {
			deckID = &id
		}
	}

	limit := queryIntOrDefault(r, "limit", 15, 1, 0)

	session, err := h.masteryService.GetStudySession(r.Context(), claims.UserID, claims.Role, deckID, limit)
	if err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al obtener la sesión de estudio")
		response.Error(w, status, msg)
		return
	}
	response.JSON(w, http.StatusOK, session)
}

// SubmitReview maneja POST /api/v1/studio/study-session/review
func (h *MasteryHandler) SubmitReview(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	if !ok || claims == nil {
		response.Error(w, http.StatusUnauthorized, "No autenticado")
		return
	}

	var req models.SubmitReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Cuerpo JSON inválido")
		return
	}

	result, err := h.masteryService.SubmitReview(r.Context(), claims.UserID, req)
	if err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al registrar la revisión de estudio")
		response.Error(w, status, msg)
		return
	}
	response.JSON(w, http.StatusOK, result)
}

// GetMasteryDashboard maneja GET /api/v1/studio/mastery
func (h *MasteryHandler) GetMasteryDashboard(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	if !ok || claims == nil {
		response.Error(w, http.StatusUnauthorized, "No autenticado")
		return
	}

	dashboard, err := h.masteryService.GetMasteryDashboard(r.Context(), claims.UserID)
	if err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al obtener el tablero de maestría")
		response.Error(w, status, msg)
		return
	}
	response.JSON(w, http.StatusOK, dashboard)
}

// ListMasteryItems maneja GET /api/v1/studio/mastery/items?stage=MASTERED
func (h *MasteryHandler) ListMasteryItems(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	if !ok || claims == nil {
		response.Error(w, http.StatusUnauthorized, "No autenticado")
		return
	}

	stage := optionalQueryString(r, "stage")

	items, err := h.masteryService.ListMasteryItems(r.Context(), claims.UserID, stage)
	if err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al listar las tarjetas de maestría")
		response.Error(w, status, msg)
		return
	}
	response.JSON(w, http.StatusOK, items)
}

// ListXPLedger maneja GET /api/v1/studio/xp-ledger?limit=50
func (h *MasteryHandler) ListXPLedger(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	if !ok || claims == nil {
		response.Error(w, http.StatusUnauthorized, "No autenticado")
		return
	}

	limit := queryIntOrDefault(r, "limit", 50, 1, 0)

	ledger, err := h.masteryService.ListXPLedger(r.Context(), claims.UserID, limit)
	if err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al obtener el libro mayor de XP")
		response.Error(w, status, msg)
		return
	}
	response.JSON(w, http.StatusOK, ledger)
}
