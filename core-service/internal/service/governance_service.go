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

// GovernanceService define el contrato SOA del Módulo 2 (Gobernanza Institucional,
// Roles y Boletos — CU02, CU03, CU04) según la sección 2.2 de .context/INTERFACES.md.
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
// CU02: Editor Promotion Requests (RN02)
// ------------------------------------------------------------------------------

func (s *governanceService) CreateEditorRequest(ctx context.Context, learnerID string, req models.CreateEditorRequestDTO) (*models.EditorRequest, error) {
	learnerUUID, err := ParseDomainUUID(learnerID, "Invalid user identifier.")
	if err != nil {
		return nil, err
	}

	req.InstitutionName = strings.TrimSpace(req.InstitutionName)
	if len(req.InstitutionName) < 3 || len(req.InstitutionName) > 150 {
		return nil, NewDomainErrorWithCode(ErrInvalidInput, "VALIDATION_ERROR", "Institution name must be between 3 and 150 characters.")
	}

	req.Justification = strings.TrimSpace(req.Justification)
	if len(req.Justification) < 20 {
		return nil, NewDomainErrorWithCode(ErrInvalidInput, "VALIDATION_ERROR", "Justification must be at least 20 characters long.")
	}

	// RN02: Check if user already holds a higher role
	account, err := s.queries.GetAccountByID(ctx, learnerUUID)
	if err != nil {
		return nil, MapDBError(err, "User account not found.")
	}
	if account.Role == db.AccountRoleEDITOR || account.Role == db.AccountRoleADMIN {
		return nil, NewDomainErrorWithCode(ErrConflict, "ALREADY_PRIVILEGED", "User already has the Editor or Administrator role.")
	}

	// RN12: Validate credential media URL if supplied
	if err := ValidateInternalMediaURL(req.CredentialsURL); err != nil {
		return nil, err
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
		return nil, MapDBError(err, "You already have a pending promotion request.")
	}

	return s.getEditorRequestDetail(ctx, createdReq.IDRequest)
}

func (s *governanceService) ListMyEditorRequests(ctx context.Context, learnerID string) ([]models.EditorRequest, error) {
	learnerUUID, err := ParseDomainUUID(learnerID, "Invalid user identifier.")
	if err != nil {
		return nil, err
	}

	rows, err := s.queries.ListMyEditorRequests(ctx, learnerUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to list user editor requests: %w", err)
	}

	results := make([]models.EditorRequest, len(rows))
	for i, r := range rows {
		results[i] = mapEditorRequestRow(r.IDRequest, r.IDUser, r.UserName, r.UserFullName, r.UserEmail, r.InstitutionName, r.Justification, r.CredentialsUrl, r.Status, r.ReviewedBy, r.ReviewNotes, r.CreatedAt, r.ResolvedAt)
	}
	return results, nil
}

// ------------------------------------------------------------------------------
// CU03: Request & User Administration (RN08, Last Admin Protection)
// ------------------------------------------------------------------------------

func (s *governanceService) ListEditorRequests(ctx context.Context, statusFilter *string) ([]models.EditorRequest, error) {
	var st pgtype.Text
	if statusFilter != nil {
		st = StringToText(strings.ToUpper(*statusFilter))
	}

	rows, err := s.queries.ListEditorRequests(ctx, st)
	if err != nil {
		return nil, fmt.Errorf("failed to list editor requests: %w", err)
	}

	results := make([]models.EditorRequest, len(rows))
	for i, r := range rows {
		results[i] = mapEditorRequestRow(r.IDRequest, r.IDUser, r.UserName, r.UserFullName, r.UserEmail, r.InstitutionName, r.Justification, r.CredentialsUrl, r.Status, r.ReviewedBy, r.ReviewNotes, r.CreatedAt, r.ResolvedAt)
	}
	return results, nil
}

func (s *governanceService) ResolveEditorRequest(ctx context.Context, adminID string, requestID string, req models.ResolveEditorRequestDTO) (*models.EditorRequest, error) {
	adminUUID, err := ParseDomainUUID(adminID, "Invalid administrator identifier.")
	if err != nil {
		return nil, err
	}
	requestUUID, err := ParseDomainUUID(requestID, "Invalid request identifier.")
	if err != nil {
		return nil, err
	}

	if req.Status != models.RequestApproved && req.Status != models.RequestRejected {
		return nil, NewDomainErrorWithCode(ErrInvalidInput, "VALIDATION_ERROR", "Status must be 'APPROVED' or 'REJECTED'.")
	}

	tx, err := s.databaseInstance.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	qtx := s.queries.WithTx(tx)

	currentReq, err := qtx.GetEditorRequestByID(ctx, requestUUID)
	if err != nil {
		return nil, MapDBError(err, "Editor promotion request not found.")
	}
	if currentReq.Status != "PENDING" {
		return nil, NewDomainErrorWithCode(ErrConflict, "REQUEST_ALREADY_RESOLVED", "This request has already been resolved.")
	}

	_, err = qtx.ResolveEditorRequest(ctx, db.ResolveEditorRequestParams{
		IDRequest:   requestUUID,
		Status:      db.RequestStatus(req.Status),
		ReviewedBy:  UUIDToPgtype(adminUUID),
		ReviewNotes: PtrToText(req.ReviewNotes),
	})
	if err != nil {
		return nil, MapDBError(err, "Failed to resolve editor promotion request.")
	}

	// Atomic promotion to EDITOR + token_version bump if approved
	if req.Status == models.RequestApproved {
		if err := qtx.UpdateAccountRole(ctx, db.UpdateAccountRoleParams{
			IDAccount: currentReq.IDUser,
			Role:      db.AccountRoleEDITOR,
		}); err != nil {
			return nil, MapDBError(err, "Failed to promote account to EDITOR.")
		}
		if _, err := qtx.IncrementAccountTokenVersion(ctx, currentReq.IDUser); err != nil {
			return nil, fmt.Errorf("failed to increment token version: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
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
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	total, err := s.queries.CountUsers(ctx, db.CountUsersParams{
		Search: sText,
		Role:   rText,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
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
	targetUUID, err := ParseDomainUUID(targetUserID, "Invalid user identifier.")
	if err != nil {
		return nil, err
	}

	tx, err := s.databaseInstance.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	qtx := s.queries.WithTx(tx)

	targetAcc, err := qtx.GetAccountByID(ctx, targetUUID)
	if err != nil {
		return nil, MapDBError(err, "User account not found.")
	}

	// CU03 4a: Protect the last active administrator
	if targetAcc.Role == db.AccountRoleADMIN {
		willDemote := req.Role != nil && *req.Role != models.RoleAdmin
		willDeactivate := req.IsActive != nil && !*req.IsActive
		if willDemote || willDeactivate {
			activeAdmins, err := qtx.CountActiveAdmins(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to check active administrators: %w", err)
			}
			if activeAdmins <= 1 {
				return nil, NewDomainErrorWithCode(ErrConflict, "LAST_ADMIN_PROTECTED", "Cannot demote or deactivate the only active administrator.")
			}
		}
	}

	var roleParam db.NullAccountRole
	if req.Role != nil {
		if *req.Role != models.RoleLearner && *req.Role != models.RoleEditor && *req.Role != models.RoleAdmin {
			return nil, NewDomainErrorWithCode(ErrInvalidInput, "VALIDATION_ERROR", "Role must be 'LEARNER', 'EDITOR', or 'ADMIN'.")
		}
		roleParam = db.NullAccountRole{AccountRole: db.AccountRole(*req.Role), Valid: true}
	}

	var activeParam pgtype.Bool
	if req.IsActive != nil {
		activeParam = pgtype.Bool{Bool: *req.IsActive, Valid: true}
	}

	row, err := qtx.UpdateUserRoleOrStatus(ctx, db.UpdateUserRoleOrStatusParams{
		IDUser:   targetUUID,
		Role:     roleParam,
		IsActive: activeParam,
	})
	if err != nil {
		return nil, MapDBError(err, "Failed to update user role or status.")
	}

	deactivating := req.IsActive != nil && !*req.IsActive
	if req.Role != nil || deactivating {
		if _, err := qtx.IncrementAccountTokenVersion(ctx, targetUUID); err != nil {
			return nil, fmt.Errorf("failed to increment token version: %w", err)
		}
	}
	if deactivating {
		if err := qtx.RevokeAllRefreshTokensForAccount(ctx, targetUUID); err != nil {
			return nil, fmt.Errorf("failed to revoke refresh tokens: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
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
// CU04: Institutional Invitation Tickets (Issue, Revoke & Redeem FOR UPDATE)
// ------------------------------------------------------------------------------

func (s *governanceService) CreateInvitationTicket(ctx context.Context, adminID string, req models.CreateTicketRequest) (*models.InvitationTicket, error) {
	adminUUID, err := ParseDomainUUID(adminID, "Invalid administrator identifier.")
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
		return nil, NewDomainErrorWithCode(ErrInvalidInput, "VALIDATION_ERROR", "Granted role must be 'LEARNER' or 'EDITOR'.")
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
		return nil, MapDBError(err, "An invitation ticket with that code already exists.")
	}

	return mapInvitationTicketRow(row.IDTicket, row.Code, row.CreatedBy, row.GrantedRole, row.InstitutionName, row.MaxUses, row.CurrentUses, row.IsActive, row.ExpiresAt, row.CreatedAt), nil
}

func (s *governanceService) ListInvitationTickets(ctx context.Context, activeOnly bool) ([]models.InvitationTicket, error) {
	rows, err := s.queries.ListInvitationTickets(ctx, pgtype.Bool{Bool: activeOnly, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to list invitation tickets: %w", err)
	}

	tickets := make([]models.InvitationTicket, len(rows))
	for i, r := range rows {
		tickets[i] = *mapInvitationTicketRow(r.IDTicket, r.Code, r.CreatedBy, r.GrantedRole, r.InstitutionName, r.MaxUses, r.CurrentUses, r.IsActive, r.ExpiresAt, r.CreatedAt)
	}
	return tickets, nil
}

func (s *governanceService) RevokeInvitationTicket(ctx context.Context, adminID string, ticketID string) error {
	ticketUUID, err := ParseDomainUUID(ticketID, "Invalid ticket identifier.")
	if err != nil {
		return err
	}

	_, err = s.queries.GetInvitationTicketByID(ctx, ticketUUID)
	if err != nil {
		return MapDBError(err, "Invitation ticket not found.")
	}

	if err := s.queries.RevokeInvitationTicket(ctx, ticketUUID); err != nil {
		return MapDBError(err, "Failed to revoke invitation ticket.")
	}

	return nil
}

func (s *governanceService) RedeemInvitationTicket(ctx context.Context, userID string, code string) (*models.UserProfile, error) {
	userUUID, err := ParseDomainUUID(userID, "Invalid user identifier.")
	if err != nil {
		return nil, err
	}

	code = strings.TrimSpace(strings.ToUpper(code))
	if code == "" {
		return nil, NewDomainErrorWithCode(ErrInvalidInput, "VALIDATION_ERROR", "Invitation ticket code is required.")
	}

	tx, err := s.databaseInstance.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	qtx := s.queries.WithTx(tx)

	// SELECT FOR UPDATE for strict concurrency serialization
	ticket, err := qtx.GetInvitationTicketByCodeForUpdate(ctx, code)
	if err != nil {
		return nil, MapDBError(err, "Invitation ticket not found.")
	}

	if !ticket.IsActive {
		return nil, NewDomainErrorWithCode(ErrConflict, "TICKET_REVOKED", "This invitation ticket has been revoked or is no longer active.")
	}

	if ticket.ExpiresAt.Valid && ticket.ExpiresAt.Time.Before(time.Now()) {
		return nil, NewDomainErrorWithCode(ErrConflict, "TICKET_EXPIRED", "This invitation ticket has expired.")
	}

	if ticket.CurrentUses >= ticket.MaxUses {
		return nil, NewDomainErrorWithCode(ErrConflict, "TICKET_EXHAUSTED", "This invitation ticket has reached its maximum number of uses.")
	}

	hasRedeemed, err := qtx.HasUserRedeemedTicket(ctx, db.HasUserRedeemedTicketParams{
		IDTicket: ticket.IDTicket,
		IDUser:   userUUID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to verify prior ticket redemption: %w", err)
	}
	if hasRedeemed {
		return nil, NewDomainErrorWithCode(ErrConflict, "TICKET_ALREADY_REDEEMED", "You have already redeemed this invitation ticket.")
	}

	if err := qtx.CreateTicketRedemption(ctx, db.CreateTicketRedemptionParams{
		IDTicket: ticket.IDTicket,
		IDUser:   userUUID,
	}); err != nil {
		return nil, MapDBError(err, "Failed to record ticket redemption.")
	}

	if _, err := qtx.IncrementTicketUses(ctx, ticket.IDTicket); err != nil {
		return nil, MapDBError(err, "Failed to update ticket usage count.")
	}

	// If ticket grants EDITOR and user is LEARNER, promote account + bump token_version
	currentAcc, err := qtx.GetAccountByID(ctx, userUUID)
	if err != nil {
		return nil, MapDBError(err, "User account not found.")
	}
	if ticket.GrantedRole == db.AccountRoleEDITOR && currentAcc.Role == db.AccountRoleLEARNER {
		if err := qtx.UpdateAccountRole(ctx, db.UpdateAccountRoleParams{
			IDAccount: userUUID,
			Role:      db.AccountRoleEDITOR,
		}); err != nil {
			return nil, MapDBError(err, "Failed to promote account to EDITOR.")
		}
		if _, err := qtx.IncrementAccountTokenVersion(ctx, userUUID); err != nil {
			return nil, fmt.Errorf("failed to increment token version: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	row, err := s.queries.GetUserProfileByID(ctx, userUUID)
	if err != nil {
		return nil, MapDBError(err, "Failed to load updated user profile.")
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
// Internal Helpers
// ------------------------------------------------------------------------------

func (s *governanceService) getEditorRequestDetail(ctx context.Context, id uuid.UUID) (*models.EditorRequest, error) {
	r, err := s.queries.GetEditorRequestByID(ctx, id)
	if err != nil {
		return nil, MapDBError(err, "Editor promotion request not found.")
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
