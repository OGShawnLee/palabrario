package main

import (
	"net/http"
	"net/url"
	"strconv"

	"palabrario/pkg/models"
)

// =============================================================================
// MÓDULO 3: Gobernanza Institucional, Roles y Boletos (CU02, CU03, CU04)
// =============================================================================

// RedeemInvitationTicket canjea un boleto institucional para ascenso directo de rol.
func (a *App) RedeemInvitationTicket(req models.RedeemTicketRequest) (*models.UserProfile, error) {
	var profile models.UserProfile
	if err := a.doRequest(a.context(), http.MethodPost, a.coreURL+"/api/v1/governance/tickets/redeem", req, &profile); err != nil {
		return nil, err
	}
	return &profile, nil
}

// CreateEditorRequest solicita ascenso a Editor con comprobante institucional.
func (a *App) CreateEditorRequest(req models.CreateEditorRequestDTO) (*models.EditorRequest, error) {
	var editorReq models.EditorRequest
	if err := a.doRequest(a.context(), http.MethodPost, a.coreURL+"/api/v1/governance/editor-requests", req, &editorReq); err != nil {
		return nil, err
	}
	return &editorReq, nil
}

// ListMyEditorRequests lista el historial de solicitudes del usuario actual.
func (a *App) ListMyEditorRequests() ([]models.EditorRequest, error) {
	var list []models.EditorRequest
	if err := a.doRequest(a.context(), http.MethodGet, a.coreURL+"/api/v1/governance/editor-requests/me", nil, &list); err != nil {
		return nil, err
	}
	if list == nil {
		list = []models.EditorRequest{}
	}
	return list, nil
}

// ListPendingEditorRequests lista las solicitudes de editor en estado PENDING para administradores.
func (a *App) ListPendingEditorRequests() ([]models.EditorRequest, error) {
	return a.ListEditorRequests("PENDING")
}

// ListEditorRequests lista solicitudes con filtro opcional de estado.
func (a *App) ListEditorRequests(status string) ([]models.EditorRequest, error) {
	targetURL := a.coreURL + "/api/v1/governance/editor-requests"
	if status != "" {
		targetURL += "?status=" + url.QueryEscape(status)
	}
	var list []models.EditorRequest
	if err := a.doRequest(a.context(), http.MethodGet, targetURL, nil, &list); err != nil {
		return nil, err
	}
	if list == nil {
		list = []models.EditorRequest{}
	}
	return list, nil
}

// ResolveEditorRequest dictamina formalmente una solicitud de editor.
func (a *App) ResolveEditorRequest(id string, req models.ResolveEditorRequestDTO) (*models.EditorRequest, error) {
	var res models.EditorRequest
	targetURL := a.coreURL + "/api/v1/governance/editor-requests/" + url.PathEscape(id) + "/resolve"
	if err := a.doRequest(a.context(), http.MethodPut, targetURL, req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// GenerateInvitationTickets alias para emisión de boletos institucionales.
func (a *App) GenerateInvitationTickets(req models.CreateTicketRequest) (*models.InvitationTicket, error) {
	return a.CreateInvitationTicket(req)
}

// CreateInvitationTicket genera un boleto institucional de acceso/ascenso.
func (a *App) CreateInvitationTicket(req models.CreateTicketRequest) (*models.InvitationTicket, error) {
	var ticket models.InvitationTicket
	if err := a.doRequest(a.context(), http.MethodPost, a.coreURL+"/api/v1/governance/tickets", req, &ticket); err != nil {
		return nil, err
	}
	return &ticket, nil
}

// ListInvitationTickets lista boletos de invitación emitidos.
func (a *App) ListInvitationTickets(activeOnly bool) ([]models.InvitationTicket, error) {
	targetURL := a.coreURL + "/api/v1/governance/tickets"
	if activeOnly {
		targetURL += "?active_only=true"
	}
	var list []models.InvitationTicket
	if err := a.doRequest(a.context(), http.MethodGet, targetURL, nil, &list); err != nil {
		return nil, err
	}
	if list == nil {
		list = []models.InvitationTicket{}
	}
	return list, nil
}

// RevokeInvitationTicket revoca un boleto de invitación.
func (a *App) RevokeInvitationTicket(id string) error {
	return a.doRequest(a.context(), http.MethodDelete, a.coreURL+"/api/v1/governance/tickets/"+url.PathEscape(id), nil, nil)
}

// ListUsersForGovernance lista usuarios para moderación y gobierno de roles.
func (a *App) ListUsersForGovernance(q string, role string, limit int, offset int) ([]models.UserProfile, error) {
	return a.ListUsers(q, role, limit, offset)
}

// ListUsers busca y lista usuarios con paginación.
func (a *App) ListUsers(q string, role string, limit int, offset int) ([]models.UserProfile, error) {
	params := url.Values{}
	if q != "" {
		params.Set("q", q)
	}
	if role != "" {
		params.Set("role", role)
	}
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}
	if offset > 0 {
		params.Set("offset", strconv.Itoa(offset))
	}

	targetURL := a.coreURL + "/api/v1/governance/users"
	if len(params) > 0 {
		targetURL += "?" + params.Encode()
	}

	var users []models.UserProfile
	if err := a.doRequest(a.context(), http.MethodGet, targetURL, nil, &users); err != nil {
		return nil, err
	}
	if users == nil {
		users = []models.UserProfile{}
	}
	return users, nil
}

// UpdateUserRole modifica el rol de un usuario.
func (a *App) UpdateUserRole(id string, role models.AccountRole) (*models.UserProfile, error) {
	return a.UpdateUserRoleOrStatus(id, models.UpdateUserRoleOrStatusRequest{Role: &role})
}

// UpdateUserStatus activa o suspende la cuenta de un usuario.
func (a *App) UpdateUserStatus(id string, isActive bool) (*models.UserProfile, error) {
	return a.UpdateUserRoleOrStatus(id, models.UpdateUserRoleOrStatusRequest{IsActive: &isActive})
}

// UpdateUserRoleOrStatus actualiza rol y/o estado de activación de usuario.
func (a *App) UpdateUserRoleOrStatus(id string, req models.UpdateUserRoleOrStatusRequest) (*models.UserProfile, error) {
	var user models.UserProfile
	targetURL := a.coreURL + "/api/v1/governance/users/" + url.PathEscape(id) + "/access"
	if err := a.doRequest(a.context(), http.MethodPut, targetURL, req, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateCommunityReport emite un reporte comunitario sobre léxico o sentidos.
func (a *App) CreateCommunityReport(req models.CreateReportRequest) (*models.ContentReport, error) {
	return a.CreateReport(req)
}

// CreateReport emite un reporte de moderación.
func (a *App) CreateReport(req models.CreateReportRequest) (*models.ContentReport, error) {
	var rep models.ContentReport
	if err := a.doRequest(a.context(), http.MethodPost, a.coreURL+"/api/v1/reports", req, &rep); err != nil {
		return nil, err
	}
	return &rep, nil
}

// ListCommunityReports lista reportes abiertos de la comunidad.
func (a *App) ListCommunityReports(status string) ([]models.ContentReport, error) {
	return a.ListReports(status)
}

// ListReports lista reportes con filtro de estado.
func (a *App) ListReports(status string) ([]models.ContentReport, error) {
	targetURL := a.coreURL + "/api/v1/reports"
	if status != "" {
		targetURL += "?status=" + url.QueryEscape(status)
	}
	var list []models.ContentReport
	if err := a.doRequest(a.context(), http.MethodGet, targetURL, nil, &list); err != nil {
		return nil, err
	}
	if list == nil {
		list = []models.ContentReport{}
	}
	return list, nil
}

// ResolveCommunityReport resuelve formalmente un reporte comunitario.
func (a *App) ResolveCommunityReport(reportId string, req models.ResolveReportRequest) (*models.ContentReport, error) {
	return a.ResolveReport(reportId, req)
}

// ResolveReport resuelve un reporte de contenido.
func (a *App) ResolveReport(reportId string, req models.ResolveReportRequest) (*models.ContentReport, error) {
	var rep models.ContentReport
	targetURL := a.coreURL + "/api/v1/reports/" + url.PathEscape(reportId) + "/resolve"
	if err := a.doRequest(a.context(), http.MethodPut, targetURL, req, &rep); err != nil {
		return nil, err
	}
	return &rep, nil
}
