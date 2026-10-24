package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"palabrario/core-service/internal/database"
	"palabrario/core-service/internal/db"
	"palabrario/pkg/models"
)

// GovernanceService define el contrato SOA del Módulo 2 (Gobernanza Institucional, Roles y Boletos — CU02, CU03, CU04)
// según la sección 2.2.2 de .context/INTERFACES.md.
type GovernanceService interface {
	CreateEditorRequest(ctx context.Context, learnerID string, req models.CreateEditorRequestDTO) (*models.EditorRequest, error)
	ListMyEditorRequests(ctx context.Context, learnerID string) ([]models.EditorRequest, error)
	ListEditorRequests(ctx context.Context, statusFilter *string) ([]models.EditorRequest, error)
	ResolveEditorRequest(ctx context.Context, adminID string, requestID string, req models.ResolveEditorRequestDTO) (*models.EditorRequest, error)

	ListUsers(ctx context.Context, search string, roleFilter *string, limit, offset int) ([]models.UserProfile, int, error)
	UpdateUserRoleOrStatus(ctx context.Context, adminID string, targetUserID string, req models.UpdateUserRoleOrStatusRequest) (*models.UserProfile, error)

	CreateInvitationTicket(ctx context.Context, adminID string, req models.CreateTicketRequest) (*models.InvitationTicket, error)
	ListInvitationTickets(ctx context.Context, activeOnly bool) ([]models.InvitationTicket, error)
	RevokeInvitationTicket(ctx context.Context, adminID string, ticketID string) error
	RedeemInvitationTicket(ctx context.Context, userID string, code string) (*models.UserProfile, error)
}

type governanceService struct {
	databaseInstance *database.DB
	queries          *db.Queries
}

func NewGovernanceService(databaseInstance *database.DB) GovernanceService {
	return &governanceService{
		databaseInstance: databaseInstance,
		queries:          db.New(databaseInstance.Pool),
	}
}

// ------------------------------------------------------------------------------
// CU02: Solicitud de Ascenso a Editor (RN02)
// ------------------------------------------------------------------------------

func (s *governanceService) CreateEditorRequest(ctx context.Context, learnerID string, req models.CreateEditorRequestDTO) (*models.EditorRequest, error) {
	learnerUUID, err := ParseDomainUUID(learnerID, "Identificador de usuario inválido")
	if err != nil {
		return nil, err
	}

	req.InstitutionName = strings.TrimSpace(req.InstitutionName)
	if len(req.InstitutionName) < 3 || len(req.InstitutionName) > 150 {
		return nil, NewDomainError(ErrInvalidInput, "El nombre de la institución debe tener entre 3 y 150 caracteres")
	}

	req.Justification = strings.TrimSpace(req.Justification)
	if len(req.Justification) < 20 {
		return nil, NewDomainError(ErrInvalidInput, "La justificación debe tener al menos 20 caracteres")
	}

	// RN02: Verificar si ya cuenta con rol superior
	account, err := s.queries.GetAccountByID(ctx, learnerUUID)
	if err != nil {
		return nil, MapDBError(err, "Cuenta de usuario no encontrada")
	}
	if account.Role == db.AccountRoleEDITOR || account.Role == db.AccountRoleADMIN {
		return nil, NewDomainError(ErrConflict, "El usuario ya cuenta con rol de Editor o Administrador")
	}

	credURL := ""
	if req.CredentialsURL != nil {
		credURL = strings.TrimSpace(*req.CredentialsURL)
	}

	createdReq, err := s.queries.CreateEditorRequest(ctx, db.CreateEditorRequestParams{
		IDUser:          learnerUUID,
		InstitutionName: req.InstitutionName,
		Justification:   req.Justification,
		CredentialsUrl:  credURL,
	})
	if err != nil {
		return nil, MapDBError(err, "Ya tienes una solicitud de promoción pendiente de revisión")
	}

	// Obtener la solicitud completa con datos de usuario
	return s.getEditorRequestDetail(ctx, createdReq.IDRequest)
}

func (s *governanceService) ListMyEditorRequests(ctx context.Context, learnerID string) ([]models.EditorRequest, error) {
	learnerUUID, err := ParseDomainUUID(learnerID, "Identificador de usuario inválido")
	if err != nil {
		return nil, err
	}

	rows, err := s.queries.ListMyEditorRequests(ctx, learnerUUID)
	if err != nil {
		return nil, fmt.Errorf("error listando solicitudes del usuario: %w", err)
	}

	results := make([]models.EditorRequest, len(rows))
	for i, r := range rows {
		results[i] = mapEditorRequestRow(r.IDRequest, r.IDUser, r.UserName, r.UserFullName, r.UserEmail, r.InstitutionName, r.Justification, r.CredentialsUrl, r.Status, r.ReviewedBy, r.ReviewNotes, r.CreatedAt, r.ResolvedAt)
	}
	return results, nil
}

// ------------------------------------------------------------------------------
// CU03: Administración de Solicitudes y Usuarios (RN08, Protección Último Admin)
// ------------------------------------------------------------------------------

func (s *governanceService) ListEditorRequests(ctx context.Context, statusFilter *string) ([]models.EditorRequest, error) {
	var st pgtype.Text
	if statusFilter != nil {
		st = StringToText(strings.ToUpper(*statusFilter))
	}

	rows, err := s.queries.ListEditorRequests(ctx, st)
	if err != nil {
		return nil, fmt.Errorf("error listando solicitudes de acreditación: %w", err)
	}

	results := make([]models.EditorRequest, len(rows))
	for i, r := range rows {
		results[i] = mapEditorRequestRow(r.IDRequest, r.IDUser, r.UserName, r.UserFullName, r.UserEmail, r.InstitutionName, r.Justification, r.CredentialsUrl, r.Status, r.ReviewedBy, r.ReviewNotes, r.CreatedAt, r.ResolvedAt)
	}
	return results, nil
}

func (s *governanceService) ResolveEditorRequest(ctx context.Context, adminID string, requestID string, req models.ResolveEditorRequestDTO) (*models.EditorRequest, error) {
	adminUUID, err := ParseDomainUUID(adminID, "Identificador de administrador inválido")
	if err != nil {
		return nil, err
	}
	requestUUID, err := ParseDomainUUID(requestID, "Identificador de solicitud inválido")
	if err != nil {
		return nil, err
	}

	if req.Status != models.RequestApproved && req.Status != models.RequestRejected {
		return nil, NewDomainError(ErrInvalidInput, "El estado debe ser 'APPROVED' o 'REJECTED'")
	}

	tx, err := s.databaseInstance.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("error iniciando transacción: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	qtx := s.queries.WithTx(tx)

	currentReq, err := qtx.GetEditorRequestByID(ctx, requestUUID)
	if err != nil {
		return nil, MapDBError(err, "Solicitud de acreditación no encontrada")
	}
	if currentReq.Status != "PENDING" {
		return nil, NewDomainError(ErrConflict, "La solicitud ya ha sido resuelta previamente")
	}

	_, err = qtx.ResolveEditorRequest(ctx, db.ResolveEditorRequestParams{
		IDRequest:   requestUUID,
		Status:      db.RequestStatus(req.Status),
		ReviewedBy:  UUIDToPgtype(adminUUID),
		ReviewNotes: PtrToText(req.ReviewNotes),
	})
	if err != nil {
		return nil, MapDBError(err, "Error resolviendo la solicitud de acreditación")
	}

	// Promoción atómica a EDITOR si fue aprobada
	if req.Status == models.RequestApproved {
		err = qtx.UpdateAccountRole(ctx, db.UpdateAccountRoleParams{
			IDAccount: currentReq.IDUser,
			Role:      db.AccountRoleEDITOR,
		})
		if err != nil {
			return nil, MapDBError(err, "Error actualizando rol a EDITOR")
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("error confirmando transacción: %w", err)
	}

	return s.getEditorRequestDetail(ctx, requestUUID)
}

func (s *governanceService) ListUsers(ctx context.Context, search string, roleFilter *string, limit, offset int) ([]models.UserProfile, int, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	sText := StringToText(search)
	var rText pgtype.Text
	if roleFilter != nil {
		rText = StringToText(strings.ToUpper(*roleFilter))
	}

	rows, err := s.queries.ListUsers(ctx, db.ListUsersParams{
		Search:    sText,
		Role:      rText,
		LimitVal:  int32(limit),
		OffsetVal: int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("error listando usuarios: %w", err)
	}

	total, err := s.queries.CountUsers(ctx, db.CountUsersParams{
		Search: sText,
		Role:   rText,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("error contando usuarios: %w", err)
	}

	users := make([]models.UserProfile, len(rows))
	for i, r := range rows {
		users[i] = models.UserProfile{
			IDUser:        r.IDUser.String(),
			Email:         r.Email,
			Name:          r.Name,
			UserName:      r.UserName,
			AvatarURL:     TextToPtr(r.AvatarUrl),
			AnchorDialect: r.AnchorDialect,
			TotalXP:       int(r.TotalXp),
			CurrentLevel:  int(r.CurrentLevel),
			CurrentStreak: int(r.CurrentStreak),
			LongestStreak: int(r.LongestStreak),
			Role:          models.AccountRole(r.Role),
			CreatedAt:     r.CreatedAt.Time,
			UpdatedAt:     r.UpdatedAt.Time,
		}
	}

	return users, int(total), nil
}

func (s *governanceService) UpdateUserRoleOrStatus(ctx context.Context, adminID string, targetUserID string, req models.UpdateUserRoleOrStatusRequest) (*models.UserProfile, error) {
	targetUUID, err := ParseDomainUUID(targetUserID, "Identificador de usuario inválido")
	if err != nil {
		return nil, err
	}

	targetAcc, err := s.queries.GetAccountByID(ctx, targetUUID)
	if err != nil {
		return nil, MapDBError(err, "Cuenta de usuario no encontrada")
	}

	// CU03 4a: Protección del último Administrador activo
	if targetAcc.Role == db.AccountRoleADMIN {
		willDemote := req.Role != nil && *req.Role != models.RoleAdmin
		willDeactivate := req.IsActive != nil && !*req.IsActive
		if willDemote || willDeactivate {
			activeAdmins, err := s.queries.CountActiveAdmins(ctx)
			if err != nil {
				return nil, fmt.Errorf("error verificando administradores activos: %w", err)
			}
			if activeAdmins <= 1 {
				return nil, NewDomainError(ErrConflict, "No se puede revocar ni desactivar al único Administrador activo del sistema")
			}
		}
	}

	var roleParam db.NullAccountRole
	if req.Role != nil {
		roleParam = db.NullAccountRole{AccountRole: db.AccountRole(*req.Role), Valid: true}
	}

	var activeParam pgtype.Bool
	if req.IsActive != nil {
		activeParam = pgtype.Bool{Bool: *req.IsActive, Valid: true}
	}

	row, err := s.queries.UpdateUserRoleOrStatus(ctx, db.UpdateUserRoleOrStatusParams{
		IDUser:   targetUUID,
		Role:     roleParam,
		IsActive: activeParam,
	})
	if err != nil {
		return nil, MapDBError(err, "Error actualizando rol o estado del usuario")
	}

	return &models.UserProfile{
		IDUser:        row.IDUser.String(),
		Email:         row.Email,
		Name:          row.Name,
		UserName:      row.UserName,
		AvatarURL:     TextToPtr(row.AvatarUrl),
		AnchorDialect: row.AnchorDialect,
		TotalXP:       int(row.TotalXp),
		CurrentLevel:  int(row.CurrentLevel),
		CurrentStreak: int(row.CurrentStreak),
		LongestStreak: int(row.LongestStreak),
		Role:          models.AccountRole(row.Role),
		CreatedAt:     row.CreatedAt.Time,
		UpdatedAt:     row.UpdatedAt.Time,
	}, nil
}

// ------------------------------------------------------------------------------
// CU04: Boletos de Invitación Institucionales (Emisión, Revocación y Redención FOR UPDATE)
// ------------------------------------------------------------------------------

func (s *governanceService) CreateInvitationTicket(ctx context.Context, adminID string, req models.CreateTicketRequest) (*models.InvitationTicket, error) {
	adminUUID, err := ParseDomainUUID(adminID, "Identificador de administrador inválido")
	if err != nil {
		return nil, err
	}

	code := ""
	if req.Code != nil && strings.TrimSpace(*req.Code) != "" {
		code = strings.ToUpper(strings.TrimSpace(*req.Code))
	} else {
		code = fmt.Sprintf("PAL-%s-%s", strings.ToUpper(uuid.New().String()[:4]), strings.ToUpper(uuid.New().String()[:4]))
	}

	if req.MaxUses < 1 {
		req.MaxUses = 1
	}

	if req.GrantedRole == "" {
		req.GrantedRole = models.RoleLearner
	}
	if req.GrantedRole != models.RoleLearner && req.GrantedRole != models.RoleEditor {
		return nil, NewDomainError(ErrInvalidInput, "El rol otorgado debe ser 'LEARNER' o 'EDITOR'")
	}

	exp := time.Now().Add(30 * 24 * time.Hour)
	if req.ExpiresAt != nil {
		exp = *req.ExpiresAt
	}

	row, err := s.queries.CreateInvitationTicket(ctx, db.CreateInvitationTicketParams{
		Code:            code,
		CreatedBy:       adminUUID,
		GrantedRole:     db.AccountRole(req.GrantedRole),
		InstitutionName: PtrToText(req.InstitutionName),
		MaxUses:         int32(req.MaxUses),
		ExpiresAt:       pgtype.Timestamptz{Time: exp, Valid: true},
	})
	if err != nil {
		return nil, MapDBError(err, "El código de boleto ya existe")
	}

	return mapInvitationTicketRow(row.IDTicket, row.Code, row.CreatedBy, row.GrantedRole, row.InstitutionName, row.MaxUses, row.CurrentUses, row.IsActive, row.ExpiresAt, row.CreatedAt), nil
}

func (s *governanceService) ListInvitationTickets(ctx context.Context, activeOnly bool) ([]models.InvitationTicket, error) {
	rows, err := s.queries.ListInvitationTickets(ctx, pgtype.Bool{Bool: activeOnly, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("error listando boletos: %w", err)
	}

	tickets := make([]models.InvitationTicket, len(rows))
	for i, r := range rows {
		tickets[i] = *mapInvitationTicketRow(r.IDTicket, r.Code, r.CreatedBy, r.GrantedRole, r.InstitutionName, r.MaxUses, r.CurrentUses, r.IsActive, r.ExpiresAt, r.CreatedAt)
	}
	return tickets, nil
}

func (s *governanceService) RevokeInvitationTicket(ctx context.Context, adminID string, ticketID string) error {
	ticketUUID, err := ParseDomainUUID(ticketID, "Identificador de boleto inválido")
	if err != nil {
		return err
	}

	_, err = s.queries.GetInvitationTicketByID(ctx, ticketUUID)
	if err != nil {
		return MapDBError(err, "Boleto de invitación no encontrado")
	}

	if err := s.queries.RevokeInvitationTicket(ctx, ticketUUID); err != nil {
		return MapDBError(err, "Error revocando boleto de invitación")
	}

	return nil
}

func (s *governanceService) RedeemInvitationTicket(ctx context.Context, userID string, code string) (*models.UserProfile, error) {
	userUUID, err := ParseDomainUUID(userID, "Identificador de usuario inválido")
	if err != nil {
		return nil, err
	}

	code = strings.TrimSpace(strings.ToUpper(code))
	if code == "" {
		return nil, NewDomainError(ErrInvalidInput, "Código de boleto requerido")
	}

	tx, err := s.databaseInstance.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("error iniciando transacción: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	qtx := s.queries.WithTx(tx)

	// SELECT FOR UPDATE para serialización estricta de concurrencia
	ticket, err := qtx.GetInvitationTicketByCodeForUpdate(ctx, code)
	if err != nil {
		return nil, MapDBError(err, "Boleto de invitación no encontrado")
	}

	if !ticket.IsActive {
		return nil, NewDomainError(ErrConflict, "El boleto de invitación ha sido revocado o ya no está activo")
	}

	if ticket.ExpiresAt.Valid && ticket.ExpiresAt.Time.Before(time.Now()) {
		return nil, NewDomainError(ErrConflict, "El boleto de invitación ha expirado")
	}

	if ticket.CurrentUses >= ticket.MaxUses {
		return nil, NewDomainError(ErrConflict, "El boleto de invitación ha alcanzado su límite máximo de usos")
	}

	// Verificar si el usuario ya redimió este boleto
	hasRedeemed, err := qtx.HasUserRedeemedTicket(ctx, db.HasUserRedeemedTicketParams{
		IDTicket: ticket.IDTicket,
		IDUser:   userUUID,
	})
	if err != nil {
		return nil, fmt.Errorf("error verificando redención previa: %w", err)
	}
	if hasRedeemed {
		return nil, NewDomainError(ErrConflict, "Ya has redimido este boleto de invitación previamente")
	}

	// Registrar redención
	if err := qtx.CreateTicketRedemption(ctx, db.CreateTicketRedemptionParams{
		IDTicket: ticket.IDTicket,
		IDUser:   userUUID,
	}); err != nil {
		return nil, MapDBError(err, "Error registrando redención del boleto")
	}

	// Incrementar contador de usos
	if _, err := qtx.IncrementTicketUses(ctx, ticket.IDTicket); err != nil {
		return nil, MapDBError(err, "Error actualizando usos del boleto")
	}

	// Si el boleto otorga rol EDITOR y el usuario es LEARNER, promover cuenta
	currentAcc, err := qtx.GetAccountByID(ctx, userUUID)
	if err != nil {
		return nil, MapDBError(err, "Cuenta de usuario no encontrada")
	}
	if ticket.GrantedRole == db.AccountRoleEDITOR && currentAcc.Role == db.AccountRoleLEARNER {
		if err := qtx.UpdateAccountRole(ctx, db.UpdateAccountRoleParams{
			IDAccount: userUUID,
			Role:      db.AccountRoleEDITOR,
		}); err != nil {
			return nil, MapDBError(err, "Error al promover cuenta a Editor")
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("error confirmando transacción: %w", err)
	}

	// Obtener perfil actualizado
	row, err := s.queries.GetUserProfileByID(ctx, userUUID)
	if err != nil {
		return nil, MapDBError(err, "Error obteniendo perfil de usuario actualizado")
	}

	return &models.UserProfile{
		IDUser:        row.IDUser.String(),
		Email:         row.Email,
		Name:          row.Name,
		UserName:      row.UserName,
		AvatarURL:     TextToPtr(row.AvatarUrl),
		AnchorDialect: row.AnchorDialect,
		TotalXP:       int(row.TotalXp),
		CurrentLevel:  int(row.CurrentLevel),
		CurrentStreak: int(row.CurrentStreak),
		LongestStreak: int(row.LongestStreak),
		Role:          models.AccountRole(row.Role),
		CreatedAt:     row.CreatedAt.Time,
		UpdatedAt:     row.UpdatedAt.Time,
	}, nil
}

// ------------------------------------------------------------------------------
// Helpers Internos
// ------------------------------------------------------------------------------

func (s *governanceService) getEditorRequestDetail(ctx context.Context, id uuid.UUID) (*models.EditorRequest, error) {
	r, err := s.queries.GetEditorRequestByID(ctx, id)
	if err != nil {
		return nil, MapDBError(err, "Solicitud de acreditación no encontrada")
	}
	res := mapEditorRequestRow(r.IDRequest, r.IDUser, r.UserName, r.UserFullName, r.UserEmail, r.InstitutionName, r.Justification, r.CredentialsUrl, r.Status, r.ReviewedBy, r.ReviewNotes, r.CreatedAt, r.ResolvedAt)
	return &res, nil
}

func mapEditorRequestRow(idReq, idUser uuid.UUID, userName, fullName, email, institution, justification, credURL, status, reviewedBy string, notes pgtype.Text, createdAt, resolvedAt pgtype.Timestamptz) models.EditorRequest {
	var credPtr *string
	if credURL != "" {
		credPtr = &credURL
	}
	var revByPtr *string
	if reviewedBy != "" {
		revByPtr = &reviewedBy
	}
	var resAtPtr *time.Time
	if resolvedAt.Valid {
		t := resolvedAt.Time
		resAtPtr = &t
	}

	return models.EditorRequest{
		IDRequest:       idReq.String(),
		IDUser:          idUser.String(),
		UserName:        userName,
		UserFullName:    fullName,
		UserEmail:       email,
		InstitutionName: institution,
		Justification:   justification,
		CredentialsURL:  credPtr,
		Status:          models.RequestStatus(status),
		ReviewedBy:      revByPtr,
		ReviewNotes:     TextToPtr(notes),
		CreatedAt:       createdAt.Time,
		ResolvedAt:      resAtPtr,
	}
}

func mapInvitationTicketRow(idTicket uuid.UUID, code string, createdBy uuid.UUID, grantedRole db.AccountRole, institution pgtype.Text, maxUses, currentUses int32, isActive bool, expiresAt, createdAt pgtype.Timestamptz) *models.InvitationTicket {
	var expPtr *time.Time
	if expiresAt.Valid {
		t := expiresAt.Time
		expPtr = &t
	}

	return &models.InvitationTicket{
		IDTicket:        idTicket.String(),
		Code:            code,
		CreatedBy:       createdBy.String(),
		GrantedRole:     models.AccountRole(grantedRole),
		InstitutionName: TextToPtr(institution),
		MaxUses:         int(maxUses),
		CurrentUses:     int(currentUses),
		IsActive:        isActive,
		ExpiresAt:       expPtr,
		CreatedAt:       createdAt.Time,
	}
}
