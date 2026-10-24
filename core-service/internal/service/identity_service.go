package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"palabrario/core-service/internal/database"
	"palabrario/core-service/internal/db"
	"palabrario/pkg/auth"
	"palabrario/pkg/models"
)

// IdentityService define el contrato SOA del Módulo 1 (Autenticación y Perfil — CU01)
// según la sección 2.1.2 de .context/INTERFACES.md.
type IdentityService interface {
	SignUp(ctx context.Context, req models.SignUpRequest) (*models.AuthResponse, error)
	SignIn(ctx context.Context, req models.SignInRequest) (*models.AuthResponse, error)
	GetMe(ctx context.Context, userID string) (*models.UserProfile, error)
	UpdateProfile(ctx context.Context, userID string, req models.UpdateProfileRequest) (*models.UserProfile, error)
	ChangePassword(ctx context.Context, userID string, req models.ChangePasswordRequest) error
	GetPublicProfileByUsername(ctx context.Context, viewerID *string, username string) (*models.PublicUserProfile, error)
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
	parsedUUID, err := ParseDomainUUID(userID, "Identificador de usuario inválido")
	if err != nil {
		return nil, err
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

func (s *identityService) UpdateProfile(ctx context.Context, userID string, req models.UpdateProfileRequest) (*models.UserProfile, error) {
	parsedUUID, err := ParseDomainUUID(userID, "Identificador de usuario inválido")
	if err != nil {
		return nil, err
	}

	req.Name = strings.TrimSpace(req.Name)
	if len(req.Name) < 2 || len(req.Name) > 100 {
		return nil, NewDomainError(ErrInvalidInput, "El nombre debe tener entre 2 y 100 caracteres")
	}

	req.AnchorDialect = strings.TrimSpace(strings.ToUpper(req.AnchorDialect))
	if req.AnchorDialect == "" {
		req.AnchorDialect = "MX"
	}
	validDialects := map[string]bool{"MX": true, "ES": true, "AR": true, "CO": true, "CL": true, "CU": true, "CAM": true, "PAN": true}
	if !validDialects[req.AnchorDialect] {
		return nil, NewDomainError(ErrInvalidInput, "Dialecto ancla no válido (valores permitidos: 'MX', 'ES', 'AR', 'CO', 'CL', 'CU', 'CAM', 'PAN')")
	}

	row, err := s.queries.UpdateUserProfile(ctx, db.UpdateUserProfileParams{
		IDUser:        parsedUUID,
		Name:          req.Name,
		AnchorDialect: req.AnchorDialect,
		AvatarUrl:     PtrToText(req.AvatarURL),
	})
	if err != nil {
		return nil, MapDBError(err, "Error actualizando el perfil de usuario")
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

func (s *identityService) ChangePassword(ctx context.Context, userID string, req models.ChangePasswordRequest) error {
	parsedUUID, err := ParseDomainUUID(userID, "Identificador de usuario inválido")
	if err != nil {
		return err
	}

	if req.CurrentPassword == "" {
		return NewDomainError(ErrInvalidInput, "La contraseña actual es obligatoria")
	}
	if len(req.NewPassword) < 8 {
		return NewDomainError(ErrInvalidInput, "La nueva contraseña debe tener al menos 8 caracteres")
	}

	account, err := s.queries.GetAccountByID(ctx, parsedUUID)
	if err != nil {
		return MapDBError(err, "Cuenta no encontrada")
	}

	if !auth.CheckPassword(req.CurrentPassword, account.PasswordHash) {
		return NewDomainError(ErrUnauthorized, "La contraseña actual no es correcta")
	}

	newHash, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("error al cifrar nueva contraseña: %w", err)
	}

	if err := s.queries.UpdateAccountPassword(ctx, db.UpdateAccountPasswordParams{
		IDAccount:    parsedUUID,
		PasswordHash: newHash,
	}); err != nil {
		return MapDBError(err, "Error actualizando la contraseña")
	}

	return nil
}

func (s *identityService) GetPublicProfileByUsername(ctx context.Context, viewerID *string, username string) (*models.PublicUserProfile, error) {
	username = strings.TrimSpace(strings.ToLower(username))
	if username == "" {
		return nil, NewDomainError(ErrInvalidInput, "Nombre de usuario requerido")
	}

	userRow, err := s.queries.GetUserProfileByUsername(ctx, username)
	if err != nil {
		return nil, MapDBError(err, "Usuario no encontrado")
	}

	isOwner := false
	if viewerID != nil && strings.TrimSpace(*viewerID) != "" {
		if strings.EqualFold(userRow.IDUser.String(), strings.TrimSpace(*viewerID)) {
			isOwner = true
		}
	}

	// 1. Logros
	achRows, err := s.queries.ListUserAchievements(ctx, userRow.IDUser)
	unlockedBadges := make([]models.AchievementStatus, 0)
	if err == nil {
		for _, a := range achRows {
			var unlAt *time.Time
			if a.UnlockedAt.Valid {
				t := a.UnlockedAt.Time
				unlAt = &t
			}
			unlockedBadges = append(unlockedBadges, models.AchievementStatus{
				IDAchievement: fmt.Sprintf("%d", a.IDAchievement),
				Code:          a.Slug,
				Name:          a.Title,
				Description:   a.Description,
				Category:      a.Category,
				Tier:          a.Tier,
				IconName:      a.BadgeIconUrl,
				XPReward:      int(a.XpReward),
				IsUnlocked:    a.IsUnlocked,
				UnlockedAt:    unlAt,
			})
		}
	}

	// 2. Barajas públicas
	deckRows, err := s.queries.ListUserPublicDecks(ctx, userRow.IDUser)
	publicDecks := make([]models.Deck, 0)
	if err == nil {
		for _, d := range deckRows {
			publicDecks = append(publicDecks, models.Deck{
				IDDeck:      int(d.IDDeck),
				IDUser:      d.IDUser.String(),
				Title:       d.Title,
				Description: TextToPtr(d.Description),
				IsPublic:    d.IsPublic,
				IsOfficial:  d.IsOfficial,
				ItemCount:   int(d.ItemCount),
				CreatedAt:   d.CreatedAt.Time,
				UpdatedAt:   d.UpdatedAt.Time,
			})
		}
	}

	// 3. Colecciones públicas
	colRows, err := s.queries.ListUserPublicCollections(ctx, userRow.IDUser)
	publicCollections := make([]models.Collection, 0)
	if err == nil {
		for _, c := range colRows {
			publicCollections = append(publicCollections, models.Collection{
				IDCollection: int(c.IDCollection),
				IDUser:       c.IDUser.String(),
				Title:        c.Title,
				Description:  TextToPtr(c.Description),
				IsPublic:     c.IsPublic,
				IsFeatured:   c.IsFeatured,
				IsOfficial:   c.IsOfficial,
				CreatedAt:    c.CreatedAt.Time,
				UpdatedAt:    c.UpdatedAt.Time,
			})
		}
	}

	// 4. Estadísticas de Arena
	matchesWon := 0
	podiumCount := 0
	stats, err := s.queries.GetUserArenaStats(ctx, UUIDToPgtype(userRow.IDUser))
	if err == nil {
		matchesWon = int(stats.MatchesWon)
		podiumCount = int(stats.PodiumCount)
	}

	publicProfile := &models.PublicUserProfile{
		IDUser:            userRow.IDUser.String(),
		Name:              userRow.Name,
		UserName:          userRow.UserName,
		AvatarURL:         TextToPtr(userRow.AvatarUrl),
		AnchorDialect:     userRow.AnchorDialect,
		Role:              models.AccountRole(userRow.Role),
		TotalXP:           int(userRow.TotalXp),
		CurrentLevel:      int(userRow.CurrentLevel),
		CurrentStreak:     int(userRow.CurrentStreak),
		LongestStreak:     int(userRow.LongestStreak),
		UnlockedBadges:    unlockedBadges,
		PublicDecks:       publicDecks,
		PublicCollections: publicCollections,
		MatchesWon:        matchesWon,
		PodiumCount:       podiumCount,
		IsOwner:           isOwner,
		CreatedAt:         userRow.CreatedAt.Time,
	}

	// Campos exclusivos para el propietario (RN01)
	if isOwner {
		email := userRow.Email
		publicProfile.Email = &email

		privCount, err := s.queries.CountUserPrivateDecks(ctx, userRow.IDUser)
		if err == nil {
			privInt := int(privCount)
			publicProfile.PrivateDecksCount = &privInt
		}

		dueCount, err := s.queries.CountDueStudySessionCardsAll(ctx, userRow.IDUser)
		if err == nil {
			dueInt := int(dueCount)
			publicProfile.DueReviewsToday = &dueInt
		}
	}

	return publicProfile, nil
}
