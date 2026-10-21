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

type CurationHandler struct {
	curationService service.CurationService
}

func NewCurationHandler(curationService service.CurationService) *CurationHandler {
	return &CurationHandler{curationService: curationService}
}

// ListPendingQueue maneja GET /api/v1/curation/queue?region={MX}
func (h *CurationHandler) ListPendingQueue(w http.ResponseWriter, r *http.Request) {
	var regionFilter *string
	if reg := strings.TrimSpace(r.URL.Query().Get("region")); reg != "" {
		regionFilter = &reg
	}

	queue, err := h.curationService.ListPendingQueue(r.Context(), regionFilter)
	if err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al listar la cola de curaduría")
		response.Error(w, status, msg)
		return
	}

	response.JSON(w, http.StatusOK, queue)
}

// ListMyProposals maneja GET /api/v1/curation/my-proposals?state={DRAFT}
func (h *CurationHandler) ListMyProposals(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	if !ok || claims == nil {
		response.Error(w, http.StatusUnauthorized, "No autenticado")
		return
	}

	var stateFilter *string
	if st := strings.TrimSpace(r.URL.Query().Get("state")); st != "" {
		stateFilter = &st
	}

	proposals, err := h.curationService.ListMyProposals(r.Context(), claims.UserID, stateFilter)
	if err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al listar las propuestas del editor")
		response.Error(w, status, msg)
		return
	}

	response.JSON(w, http.StatusOK, proposals)
}

// GetProposalDetail maneja GET /api/v1/curation/senses/{senseId}
func (h *CurationHandler) GetProposalDetail(w http.ResponseWriter, r *http.Request) {
	senseIdStr := chi.URLParam(r, "senseId")
	senseID, err := strconv.Atoi(senseIdStr)
	if err != nil || senseID <= 0 {
		response.Error(w, http.StatusBadRequest, "El identificador de acepción debe ser un número entero positivo")
		return
	}

	detail, err := h.curationService.GetProposalDetail(r.Context(), senseID)
	if err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al consultar detalle de propuesta")
		response.Error(w, status, msg)
		return
	}

	response.JSON(w, http.StatusOK, detail)
}

// SubmitVote maneja POST /api/v1/curation/senses/{senseId}/vote
func (h *CurationHandler) SubmitVote(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	if !ok || claims == nil {
		response.Error(w, http.StatusUnauthorized, "No autenticado")
		return
	}

	senseIdStr := chi.URLParam(r, "senseId")
	senseID, err := strconv.Atoi(senseIdStr)
	if err != nil || senseID <= 0 {
		response.Error(w, http.StatusBadRequest, "El identificador de acepción debe ser un número entero positivo")
		return
	}

	var req models.SubmitVoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Cuerpo JSON inválido")
		return
	}

	updatedDetail, err := h.curationService.SubmitVote(r.Context(), claims.UserID, senseID, req)
	if err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al emitir voto de curaduría")
		response.Error(w, status, msg)
		return
	}

	response.JSON(w, http.StatusOK, updatedDetail)
}

// AdminOverrideState maneja PUT /api/v1/curation/senses/{senseId}/state
func (h *CurationHandler) AdminOverrideState(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	if !ok || claims == nil {
		response.Error(w, http.StatusUnauthorized, "No autenticado")
		return
	}

	senseIdStr := chi.URLParam(r, "senseId")
	senseID, err := strconv.Atoi(senseIdStr)
	if err != nil || senseID <= 0 {
		response.Error(w, http.StatusBadRequest, "El identificador de acepción debe ser un número entero positivo")
		return
	}

	var req models.AdminOverrideStateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Cuerpo JSON inválido")
		return
	}

	senseDetail, err := h.curationService.AdminOverrideState(r.Context(), claims.UserID, senseID, req)
	if err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al actualizar estado administrativo del sentido")
		response.Error(w, status, msg)
		return
	}

	response.JSON(w, http.StatusOK, senseDetail)
}
