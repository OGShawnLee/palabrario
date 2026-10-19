package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"palabrario/core-service/internal/database"
	"palabrario/core-service/internal/db"
	"palabrario/pkg/auth"
	"palabrario/pkg/models"
)

// IdentityService define el contrato SOA del Módulo 1 (Autenticación y Perfil — CU01)
// según la sección 2.1.2 de .context/INTERFACES.md.
// Los métodos adicionales de CU01 (UpdateProfile, ChangePassword, GetPublicProfileByUsername)
// se añadirán en la Fase 3.5 siguiendo este mismo molde.
type IdentityService interface {
	SignUp(ctx context.Context, req models.SignUpRequest) (*models.AuthResponse, error)
	SignIn(ctx context.Context, req models.SignInRequest) (*models.AuthResponse, error)
	GetMe(ctx context.Context, userID string) (*models.UserProfile, error)
}

type identityService struct {
	db           *database.DB
	queries      *db.Queries
	tokenService *auth.TokenService
}

func NewIdentityService(databaseInstance *database.DB, tokenService *auth.TokenService) IdentityService {
	return &identityService{
		db:           databaseInstance,
		queries:      db.New(databaseInstance.Pool),
		tokenService: tokenService,
	}
}

func (s *identityService) SignUp(ctx context.Context, req models.SignUpRequest) (*models.AuthResponse, error) {
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.UserName = strings.TrimSpace(strings.ToLower(req.UserName))
	req.Name = strings.TrimSpace(req.Name)
	req.AnchorDialect = strings.TrimSpace(strings.ToUpper(req.AnchorDialect))

	if req.Email == "" || req.Password == "" || req.Name == "" || req.UserName == "" {
		return nil, NewDomainError(ErrInvalidInput, "Todos los campos son obligatorios")
	}
	if req.AnchorDialect == "" {
		req.AnchorDialect = "MX"
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("error al procesar la contraseña: %w", err)
	}

	tx, err := s.db.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("error iniciando transacción: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	qtx := s.queries.WithTx(tx)

	account, err := qtx.CreateAccount(ctx, db.CreateAccountParams{
		Email:        req.Email,
		PasswordHash: hash,
	})
	if err != nil {
		return nil, MapDBError(err, "El correo o nombre de usuario ya está registrado")
	}

	userProfile, err := qtx.CreateUserProfile(ctx, db.CreateUserProfileParams{
		IDUser:        account.IDAccount,
		Name:          req.Name,
		UserName:      req.UserName,
		AnchorDialect: req.AnchorDialect,
	})
	if err != nil {
		return nil, MapDBError(err, "El correo o nombre de usuario ya está registrado")
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("error confirmando transacción: %w", err)
	}

	profile := models.UserProfile{
		IDUser:        account.IDAccount.String(),
		Email:         account.Email,
		Name:          userProfile.Name,
		UserName:      userProfile.UserName,
		AvatarURL:     TextToPtr(userProfile.AvatarUrl),
		AnchorDialect: userProfile.AnchorDialect,
		TotalXP:       int(userProfile.TotalXp),
		CurrentLevel:  int(userProfile.CurrentLevel),
		CurrentStreak: int(userProfile.CurrentStreak),
		LongestStreak: int(userProfile.LongestStreak),
		Role:          models.AccountRole(account.Role),
		CreatedAt:     account.CreatedAt.Time,
		UpdatedAt:     account.UpdatedAt.Time,
	}

	token, err := s.tokenService.GenerateToken(&profile)
	if err != nil {
		return nil, fmt.Errorf("error al emitir el token de sesión: %w", err)
	}

	return &models.AuthResponse{
		Token: token,
		User:  profile,
	}, nil
}

func (s *identityService) SignIn(ctx context.Context, req models.SignInRequest) (*models.AuthResponse, error) {
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if req.Email == "" || req.Password == "" {
		return nil, NewDomainError(ErrInvalidInput, "Correo y contraseña requeridos")
	}

	account, err := s.queries.GetAccountByEmail(ctx, req.Email)
	if err != nil {
		return nil, NewDomainError(ErrUnauthorized, "Credenciales incorrectas")
	}

	if !auth.CheckPassword(req.Password, account.PasswordHash) {
		return nil, NewDomainError(ErrUnauthorized, "Credenciales incorrectas")
	}

	if !account.IsActive {
		return nil, NewDomainError(ErrForbidden, "Esta cuenta ha sido suspendida por un administrador")
	}

	profile, err := s.GetMe(ctx, account.IDAccount.String())
	if err != nil {
		return nil, err
	}

	token, err := s.tokenService.GenerateToken(profile)
	if err != nil {
		return nil, fmt.Errorf("error emitiendo credenciales: %w", err)
	}

	return &models.AuthResponse{
		Token: token,
		User:  *profile,
	}, nil
}

func (s *identityService) GetMe(ctx context.Context, userID string) (*models.UserProfile, error) {
	parsedUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, NewDomainError(ErrInvalidInput, "Identificador de usuario inválido")
	}

	row, err := s.queries.GetUserProfileByID(ctx, parsedUUID)
	if err != nil {
		return nil, MapDBError(err, "Usuario no encontrado")
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
