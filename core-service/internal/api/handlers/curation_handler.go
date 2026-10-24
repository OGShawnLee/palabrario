package handlers

import (
	"encoding/json"
	"net/http"

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
	regionFilter := optionalQueryString(r, "region")

	queue, err := h.curationService.ListPendingQueue(r.Context(), regionFilter)
	if err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al listar la cola de curaduría")
		response.Error(w, status, msg)
		return
	}

	if queue == nil {
		queue = []models.SenseQuorumItem{}
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

	stateFilter := optionalQueryString(r, "state")

	proposals, err := h.curationService.ListMyProposals(r.Context(), claims.UserID, stateFilter)
	if err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al listar las propuestas del editor")
		response.Error(w, status, msg)
		return
	}

	if proposals == nil {
		proposals = []models.SenseQuorumItem{}
	}

	response.JSON(w, http.StatusOK, proposals)
}

// GetProposalDetail maneja GET /api/v1/curation/senses/{senseId}
func (h *CurationHandler) GetProposalDetail(w http.ResponseWriter, r *http.Request) {
	senseID, ok := parsePositiveIntParam(w, r, "senseId", "El identificador de acepción debe ser un número entero positivo")
	if !ok {
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

	senseID, valid := parsePositiveIntParam(w, r, "senseId", "El identificador de acepción debe ser un número entero positivo")
	if !valid {
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

	senseID, valid := parsePositiveIntParam(w, r, "senseId", "El identificador de acepción debe ser un número entero positivo")
	if !valid {
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

// CreateReport maneja POST /api/v1/reports (CU09, RF11)
func (h *CurationHandler) CreateReport(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	if !ok || claims == nil {
		response.Error(w, http.StatusUnauthorized, "No autorizado")
		return
	}

	var req models.CreateReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Cuerpo de petición inválido")
		return
	}

	report, err := h.curationService.CreateReport(r.Context(), claims.UserID, req)
	if err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al emitir reporte de contenido")
		response.Error(w, status, msg)
		return
	}

	response.JSON(w, http.StatusCreated, report)
}

// ListReports maneja GET /api/v1/reports (CU09)
func (h *CurationHandler) ListReports(w http.ResponseWriter, r *http.Request) {
	statusFilter := optionalQueryString(r, "status")

	reports, err := h.curationService.ListOpenReports(r.Context(), statusFilter)
	if err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error listando reportes de contenido")
		response.Error(w, status, msg)
		return
	}

	if reports == nil {
		reports = []models.ContentReport{}
	}

	response.JSON(w, http.StatusOK, reports)
}

// ResolveReport maneja PUT /api/v1/reports/{reportId}/resolve (CU09, RN08)
func (h *CurationHandler) ResolveReport(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	if !ok || claims == nil {
		response.Error(w, http.StatusUnauthorized, "No autorizado")
		return
	}

	reportID := chi.URLParam(r, "reportId")
	if reportID == "" {
		response.Error(w, http.StatusBadRequest, "Identificador de reporte requerido")
		return
	}

	var req models.ResolveReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Cuerpo de petición inválido")
		return
	}

	resolved, err := h.curationService.ResolveReport(r.Context(), claims.UserID, reportID, req)
	if err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al resolver reporte de contenido")
		response.Error(w, status, msg)
		return
	}

	response.JSON(w, http.StatusOK, resolved)
}
