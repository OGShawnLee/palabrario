package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"palabrario/core-service/internal/service"
	"palabrario/pkg/auth"
	"palabrario/pkg/models"
	"palabrario/pkg/response"
)

type GovernanceHandler struct {
	governanceService service.GovernanceService
}

func NewGovernanceHandler(governanceService service.GovernanceService) *GovernanceHandler {
	return &GovernanceHandler{
		governanceService: governanceService,
	}
}

// CreateEditorRequest maneja POST /api/v1/governance/editor-requests (CU02, RN02)
func (h *GovernanceHandler) CreateEditorRequest(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	if !ok || claims == nil {
		response.Error(w, http.StatusUnauthorized, "No autorizado")
		return
	}

	var req models.CreateEditorRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Cuerpo de petición inválido")
		return
	}

	res, err := h.governanceService.CreateEditorRequest(r.Context(), claims.UserID, req)
	if err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al crear la solicitud de editor")
		response.Error(w, status, msg)
		return
	}

	response.JSON(w, http.StatusCreated, res)
}

// ListMyEditorRequests maneja GET /api/v1/governance/editor-requests/me (CU02)
func (h *GovernanceHandler) ListMyEditorRequests(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	if !ok || claims == nil {
		response.Error(w, http.StatusUnauthorized, "No autorizado")
		return
	}

	res, err := h.governanceService.ListMyEditorRequests(r.Context(), claims.UserID)
	if err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al listar las solicitudes del usuario")
		response.Error(w, status, msg)
		return
	}

	if res == nil {
		res = []models.EditorRequest{}
	}

	response.JSON(w, http.StatusOK, res)
}

// ListEditorRequests maneja GET /api/v1/governance/editor-requests (CU03)
func (h *GovernanceHandler) ListEditorRequests(w http.ResponseWriter, r *http.Request) {
	statusFilter := optionalQueryString(r, "status")

	res, err := h.governanceService.ListEditorRequests(r.Context(), statusFilter)
	if err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al listar solicitudes de editor")
		response.Error(w, status, msg)
		return
	}

	if res == nil {
		res = []models.EditorRequest{}
	}

	response.JSON(w, http.StatusOK, res)
}

// ResolveEditorRequest maneja PUT /api/v1/governance/editor-requests/{id}/resolve (CU03, RN08)
func (h *GovernanceHandler) ResolveEditorRequest(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	if !ok || claims == nil {
		response.Error(w, http.StatusUnauthorized, "No autorizado")
		return
	}

	requestID := chi.URLParam(r, "id")
	if requestID == "" {
		response.Error(w, http.StatusBadRequest, "Identificador de solicitud requerido")
		return
	}

	var req models.ResolveEditorRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Cuerpo de petición inválido")
		return
	}

	res, err := h.governanceService.ResolveEditorRequest(r.Context(), claims.UserID, requestID, req)
	if err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al resolver solicitud de editor")
		response.Error(w, status, msg)
		return
	}

	response.JSON(w, http.StatusOK, res)
}

// ListUsers maneja GET /api/v1/governance/users (CU03)
func (h *GovernanceHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	roleFilter := optionalQueryString(r, "role")
	limit := queryIntOrDefault(r, "limit", 20, 1, 0)
	offset := queryIntOrDefault(r, "offset", 0, 0, 0)

	users, total, err := h.governanceService.ListUsers(r.Context(), query, roleFilter, limit, offset)
	if err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error listando usuarios")
		response.Error(w, status, msg)
		return
	}

	if users == nil {
		users = []models.UserProfile{}
	}

	w.Header().Set("X-Total-Count", strconv.Itoa(total))
	response.JSON(w, http.StatusOK, users)
}

// UpdateUserRoleOrStatus maneja PUT /api/v1/governance/users/{id}/access (CU03, RN08)
func (h *GovernanceHandler) UpdateUserRoleOrStatus(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	if !ok || claims == nil {
		response.Error(w, http.StatusUnauthorized, "No autorizado")
		return
	}

	targetUserID := chi.URLParam(r, "id")
	if targetUserID == "" {
		response.Error(w, http.StatusBadRequest, "Identificador de usuario requerido")
		return
	}

	var req models.UpdateUserRoleOrStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Cuerpo de petición inválido")
		return
	}

	updated, err := h.governanceService.UpdateUserRoleOrStatus(r.Context(), claims.UserID, targetUserID, req)
	if err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al actualizar rol o estado del usuario")
		response.Error(w, status, msg)
		return
	}

	response.JSON(w, http.StatusOK, updated)
}

// CreateInvitationTicket maneja POST /api/v1/governance/tickets (CU04, RN08)
func (h *GovernanceHandler) CreateInvitationTicket(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	if !ok || claims == nil {
		response.Error(w, http.StatusUnauthorized, "No autorizado")
		return
	}

	var req models.CreateTicketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Cuerpo de petición inválido")
		return
	}

	ticket, err := h.governanceService.CreateInvitationTicket(r.Context(), claims.UserID, req)
	if err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al crear boleto de invitación")
		response.Error(w, status, msg)
		return
	}

	response.JSON(w, http.StatusCreated, ticket)
}

// ListInvitationTickets maneja GET /api/v1/governance/tickets (CU04)
func (h *GovernanceHandler) ListInvitationTickets(w http.ResponseWriter, r *http.Request) {
	activeOnly := false
	if a := r.URL.Query().Get("active_only"); a == "true" {
		activeOnly = true
	}

	tickets, err := h.governanceService.ListInvitationTickets(r.Context(), activeOnly)
	if err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error listando boletos de invitación")
		response.Error(w, status, msg)
		return
	}

	if tickets == nil {
		tickets = []models.InvitationTicket{}
	}

	response.JSON(w, http.StatusOK, tickets)
}

// RevokeInvitationTicket maneja DELETE /api/v1/governance/tickets/{id} (CU04)
func (h *GovernanceHandler) RevokeInvitationTicket(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	if !ok || claims == nil {
		response.Error(w, http.StatusUnauthorized, "No autorizado")
		return
	}

	ticketID := chi.URLParam(r, "id")
	if ticketID == "" {
		response.Error(w, http.StatusBadRequest, "Identificador de boleto requerido")
		return
	}

	if err := h.governanceService.RevokeInvitationTicket(r.Context(), claims.UserID, ticketID); err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error revocando boleto de invitación")
		response.Error(w, status, msg)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// RedeemInvitationTicket maneja POST /api/v1/governance/tickets/redeem (CU04)
func (h *GovernanceHandler) RedeemInvitationTicket(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	if !ok || claims == nil {
		response.Error(w, http.StatusUnauthorized, "No autorizado")
		return
	}

	var req models.RedeemTicketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Cuerpo de petición inválido")
		return
	}

	profile, err := h.governanceService.RedeemInvitationTicket(r.Context(), claims.UserID, req.Code)
	if err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al redimir boleto de invitación")
		response.Error(w, status, msg)
		return
	}

	response.JSON(w, http.StatusOK, profile)
}
