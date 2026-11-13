package service

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"palabrario/core-service/internal/database"
	"palabrario/core-service/internal/db"
	"palabrario/pkg/auth"
	"palabrario/pkg/models"
)

const defaultRefreshTokenTTL = 168 * time.Hour

var (
	usernameRegex = regexp.MustCompile(`^[a-z0-9][a-z0-9_-]{2,49}$`)
	validDialects = map[string]bool{
		"MX":  true,
		"ES":  true,
		"AR":  true,
		"CO":  true,
		"CL":  true,
		"CU":  true,
		"CAM": true,
		"PAN": true,
	}
)

// IdentityService defines the SOA contract for Module 1 (Authentication and Profile — CU01).
type IdentityService interface {
	SignUp(ctx context.Context, req models.SignUpRequest) (*models.AuthResponse, error)
	SignIn(ctx context.Context, req models.SignInRequest) (*models.AuthResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*models.AuthResponse, error)
	SignOut(ctx context.Context, refreshToken string) error
	GetMe(ctx context.Context, userID string) (*models.UserProfile, error)
	UpdateProfile(ctx context.Context, userID string, req models.UpdateProfileRequest) (*models.UserProfile, error)
	ChangePassword(ctx context.Context, userID string, req models.ChangePasswordRequest) error
	GetPublicProfileByUsername(ctx context.Context, viewerID *string, username string) (*models.PublicUserProfile, error)
}

type identityService struct {
	db              *database.DB
	queries         *db.Queries
	tokenService    *auth.TokenService
	refreshTokenTTL time.Duration
}

func NewIdentityService(databaseInstance *database.DB, tokenService *auth.TokenService, refreshTokenTTL ...time.Duration) IdentityService {
	ttl := defaultRefreshTokenTTL
	if len(refreshTokenTTL) > 0 && refreshTokenTTL[0] > 0 {
		ttl = refreshTokenTTL[0]
	}
	return &identityService{
		db:              databaseInstance,
		queries:         db.New(databaseInstance.Pool),
		tokenService:    tokenService,
		refreshTokenTTL: ttl,
	}
}

func validateEmail(email string) error {
	if email == "" || len(email) > 255 || strings.ContainsAny(email, " \t\r\n") {
		return NewDomainErrorWithCode(ErrInvalidInput, "VALIDATION_ERROR", "Please enter a valid email address.")
	}
	addr, err := mail.ParseAddress(email)
	if err != nil || addr.Address != email {
		return NewDomainErrorWithCode(ErrInvalidInput, "VALIDATION_ERROR", "Please enter a valid email address.")
	}
	local, domain, ok := strings.Cut(email, "@")
	if !ok || local == "" || domain == "" || !strings.Contains(domain, ".") || strings.HasPrefix(domain, ".") || strings.HasSuffix(domain, ".") {
		return NewDomainErrorWithCode(ErrInvalidInput, "VALIDATION_ERROR", "Please enter a valid email address.")
	}
	return nil
}

func validatePassword(password string) error {
	byteLen := len([]byte(password))
	if byteLen < 8 {
		return NewDomainErrorWithCode(ErrInvalidInput, "VALIDATION_ERROR", "Password must be at least 8 characters long.")
	}
	if byteLen > auth.MaxBcryptPasswordLen {
		return NewDomainErrorWithCode(ErrInvalidInput, "VALIDATION_ERROR", fmt.Sprintf("Password must not exceed %d bytes.", auth.MaxBcryptPasswordLen))
	}
	return nil
}

func validateAnchorDialect(raw string) (string, error) {
	dialect := strings.TrimSpace(strings.ToUpper(raw))
	if dialect == "" {
		return "MX", nil
	}
	if !validDialects[dialect] {
		return "", NewDomainErrorWithCode(ErrInvalidInput, "VALIDATION_ERROR", "Invalid anchor dialect. Allowed values: MX, ES, AR, CO, CL, CU, CAM, PAN.")
	}
	return dialect, nil
}

func (s *identityService) SignUp(ctx context.Context, req models.SignUpRequest) (*models.AuthResponse, error) {
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.UserName = strings.TrimPrefix(strings.TrimSpace(strings.ToLower(req.UserName)), "@")
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		req.Name = req.UserName
	}

	if err := validateEmail(req.Email); err != nil {
		return nil, err
	}
	if !usernameRegex.MatchString(req.UserName) {
		return nil, NewDomainErrorWithCode(ErrInvalidInput, "VALIDATION_ERROR", "Username must be 3 to 50 characters, start with a letter or number, and contain only lowercase letters, numbers, underscores, or hyphens.")
	}
	nameLen := utf8.RuneCountInString(req.Name)
	if nameLen < 2 || nameLen > 100 {
		return nil, NewDomainErrorWithCode(ErrInvalidInput, "VALIDATION_ERROR", "Name must be between 2 and 100 characters.")
	}
	dialect, err := validateAnchorDialect(req.AnchorDialect)
	if err != nil {
		return nil, err
	}
	req.AnchorDialect = dialect

	if err := validatePassword(req.Password); err != nil {
		return nil, err
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	rawRefresh, refreshHash, err := auth.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	tx, err := s.db.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	qtx := s.queries.WithTx(tx)

	account, err := qtx.CreateAccount(ctx, db.CreateAccountParams{
		Email:        req.Email,
		PasswordHash: hash,
	})
	if err != nil {
		return nil, MapDBError(err, "Failed to create account.")
	}

	userProfile, err := qtx.CreateUserProfile(ctx, db.CreateUserProfileParams{
		IDUser:        account.IDAccount,
		Name:          req.Name,
		UserName:      req.UserName,
		AnchorDialect: req.AnchorDialect,
	})
	if err != nil {
		return nil, MapDBError(err, "Failed to create user profile.")
	}

	if _, err := qtx.CreateRefreshToken(ctx, db.CreateRefreshTokenParams{
		IDAccount: account.IDAccount,
		TokenHash: refreshHash,
		ExpiresAt: pgtype.Timestamptz{Time: time.Now().UTC().Add(s.refreshTokenTTL), Valid: true},
	}); err != nil {
		return nil, fmt.Errorf("failed to persist refresh token: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
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

	token, err := s.tokenService.GenerateAccessToken(&profile, int(account.TokenVersion))
	if err != nil {
		return nil, fmt.Errorf("failed to issue access token: %w", err)
	}

	return &models.AuthResponse{
		Token:        token,
		RefreshToken: rawRefresh,
		User:         profile,
	}, nil
}

func (s *identityService) SignIn(ctx context.Context, req models.SignInRequest) (*models.AuthResponse, error) {
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if req.Email == "" || req.Password == "" {
		return nil, NewDomainErrorWithCode(ErrInvalidInput, "INVALID_INPUT", "Email and password are required.")
	}

	row, err := s.queries.GetAccountWithProfileByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, NewDomainErrorWithCode(ErrUnauthorized, "INVALID_CREDENTIALS", "Invalid email or password.")
		}
		return nil, fmt.Errorf("failed to query account: %w", err)
	}

	if !auth.CheckPassword(req.Password, row.PasswordHash) {
		return nil, NewDomainErrorWithCode(ErrUnauthorized, "INVALID_CREDENTIALS", "Invalid email or password.")
	}

	if !row.IsActive {
		return nil, NewDomainErrorWithCode(ErrForbidden, "ACCOUNT_SUSPENDED", "This account has been suspended by an administrator.")
	}

	rawRefresh, refreshHash, err := auth.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	if _, err := s.queries.CreateRefreshToken(ctx, db.CreateRefreshTokenParams{
		IDAccount: row.IDAccount,
		TokenHash: refreshHash,
		ExpiresAt: pgtype.Timestamptz{Time: time.Now().UTC().Add(s.refreshTokenTTL), Valid: true},
	}); err != nil {
		return nil, fmt.Errorf("failed to persist refresh token: %w", err)
	}

	profile := models.UserProfile{
		IDUser:        row.IDAccount.String(),
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
	}

	token, err := s.tokenService.GenerateAccessToken(&profile, int(row.TokenVersion))
	if err != nil {
		return nil, fmt.Errorf("failed to issue access token: %w", err)
	}

	return &models.AuthResponse{
		Token:        token,
		RefreshToken: rawRefresh,
		User:         profile,
	}, nil
}

func (s *identityService) RefreshToken(ctx context.Context, refreshToken string) (*models.AuthResponse, error) {
	raw := strings.TrimSpace(refreshToken)
	if raw == "" {
		return nil, NewDomainErrorWithCode(ErrUnauthorized, "TOKEN_REVOKED", "Refresh token is required.")
	}

	tokenHash := auth.HashRefreshToken(raw)

	tx, err := s.db.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	qtx := s.queries.WithTx(tx)

	stored, err := qtx.GetRefreshTokenForUpdate(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, NewDomainErrorWithCode(ErrUnauthorized, "TOKEN_REVOKED", "Refresh token is invalid or has been revoked.")
		}
		return nil, fmt.Errorf("failed to look up refresh token: %w", err)
	}

	if stored.RevokedAt.Valid {
		return nil, NewDomainErrorWithCode(ErrUnauthorized, "TOKEN_REVOKED", "Refresh token has already been revoked.")
	}

	if !stored.ExpiresAt.Valid || !stored.ExpiresAt.Time.After(time.Now().UTC()) {
		return nil, NewDomainErrorWithCode(ErrUnauthorized, "TOKEN_EXPIRED", "Refresh token has expired. Please sign in again.")
	}

	profileRow, err := qtx.GetUserProfileByID(ctx, stored.IDAccount)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, NewDomainErrorWithCode(ErrUnauthorized, "UNAUTHORIZED", "User session is no longer valid.")
		}
		return nil, fmt.Errorf("failed to load user profile: %w", err)
	}

	if !profileRow.IsActive {
		if err := qtx.RevokeAllRefreshTokensForAccount(ctx, stored.IDAccount); err != nil {
			return nil, fmt.Errorf("failed to revoke refresh tokens for suspended account: %w", err)
		}
		if err := tx.Commit(ctx); err != nil {
			return nil, fmt.Errorf("failed to commit revocation for suspended account: %w", err)
		}
		return nil, NewDomainErrorWithCode(ErrForbidden, "ACCOUNT_SUSPENDED", "This account has been suspended by an administrator.")
	}

	if err := qtx.RevokeRefreshTokenByID(ctx, stored.IDToken); err != nil {
		return nil, fmt.Errorf("failed to revoke used refresh token: %w", err)
	}

	newRaw, newHash, err := auth.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	if _, err := qtx.CreateRefreshToken(ctx, db.CreateRefreshTokenParams{
		IDAccount: stored.IDAccount,
		TokenHash: newHash,
		ExpiresAt: pgtype.Timestamptz{Time: time.Now().UTC().Add(s.refreshTokenTTL), Valid: true},
	}); err != nil {
		return nil, fmt.Errorf("failed to persist rotated refresh token: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit refresh transaction: %w", err)
	}

	profile := models.UserProfile{
		IDUser:        profileRow.IDUser.String(),
		Email:         profileRow.Email,
		Name:          profileRow.Name,
		UserName:      profileRow.UserName,
		AvatarURL:     TextToPtr(profileRow.AvatarUrl),
		AnchorDialect: profileRow.AnchorDialect,
		TotalXP:       int(profileRow.TotalXp),
		CurrentLevel:  int(profileRow.CurrentLevel),
		CurrentStreak: int(profileRow.CurrentStreak),
		LongestStreak: int(profileRow.LongestStreak),
		Role:          models.AccountRole(profileRow.Role),
		CreatedAt:     profileRow.CreatedAt.Time,
		UpdatedAt:     profileRow.UpdatedAt.Time,
	}

	accessToken, err := s.tokenService.GenerateAccessToken(&profile, int(profileRow.TokenVersion))
	if err != nil {
		return nil, fmt.Errorf("failed to issue access token: %w", err)
	}

	return &models.AuthResponse{
		Token:        accessToken,
		RefreshToken: newRaw,
		User:         profile,
	}, nil
}

func (s *identityService) SignOut(ctx context.Context, refreshToken string) error {
	raw := strings.TrimSpace(refreshToken)
	if raw == "" {
		return nil
	}

	tokenHash := auth.HashRefreshToken(raw)
	if err := s.queries.RevokeRefreshTokenByHash(ctx, tokenHash); err != nil {
		return fmt.Errorf("failed to revoke refresh token: %w", err)
	}
	return nil
}

func (s *identityService) GetMe(ctx context.Context, userID string) (*models.UserProfile, error) {
	parsedUUID, err := ParseDomainUUID(userID, "Invalid user identifier.")
	if err != nil {
		return nil, err
	}

	row, err := s.queries.GetUserProfileByID(ctx, parsedUUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, NewDomainErrorWithCode(ErrUnauthorized, "UNAUTHORIZED", "User session is no longer valid.")
		}
		return nil, MapDBError(err, "Failed to load user profile.")
	}

	if !row.IsActive {
		return nil, NewDomainErrorWithCode(ErrForbidden, "ACCOUNT_SUSPENDED", "This account has been suspended by an administrator.")
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
	parsedUUID, err := ParseDomainUUID(userID, "Invalid user identifier.")
	if err != nil {
		return nil, err
	}

	req.Name = strings.TrimSpace(req.Name)
	nameLen := utf8.RuneCountInString(req.Name)
	if nameLen < 2 || nameLen > 100 {
		return nil, NewDomainErrorWithCode(ErrInvalidInput, "VALIDATION_ERROR", "Name must be between 2 and 100 characters.")
	}

	dialect, err := validateAnchorDialect(req.AnchorDialect)
	if err != nil {
		return nil, err
	}
	req.AnchorDialect = dialect

	// RN12: Validate safe avatar URL format
	if err := ValidateInternalMediaURL(req.AvatarURL); err != nil {
		return nil, err
	}

	state, err := s.queries.GetAuthAccountStateByID(ctx, parsedUUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, NewDomainErrorWithCode(ErrUnauthorized, "UNAUTHORIZED", "User session is no longer valid.")
		}
		return nil, MapDBError(err, "Failed to verify account status.")
	}
	if !state.IsActive {
		return nil, NewDomainErrorWithCode(ErrForbidden, "ACCOUNT_SUSPENDED", "This account has been suspended by an administrator.")
	}

	row, err := s.queries.UpdateUserProfile(ctx, db.UpdateUserProfileParams{
		IDUser:        parsedUUID,
		Name:          req.Name,
		AnchorDialect: req.AnchorDialect,
		AvatarUrl:     PtrToText(req.AvatarURL),
	})
	if err != nil {
		return nil, MapDBError(err, "Failed to update user profile.")
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
	parsedUUID, err := ParseDomainUUID(userID, "Invalid user identifier.")
	if err != nil {
		return err
	}

	if req.CurrentPassword == "" {
		return NewDomainErrorWithCode(ErrInvalidInput, "VALIDATION_ERROR", "Current password is required.")
	}
	if err := validatePassword(req.NewPassword); err != nil {
		return err
	}

	account, err := s.queries.GetAccountByID(ctx, parsedUUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return NewDomainErrorWithCode(ErrUnauthorized, "UNAUTHORIZED", "User session is no longer valid.")
		}
		return MapDBError(err, "Failed to look up account.")
	}

	if !account.IsActive {
		return NewDomainErrorWithCode(ErrForbidden, "ACCOUNT_SUSPENDED", "This account has been suspended by an administrator.")
	}

	if !auth.CheckPassword(req.CurrentPassword, account.PasswordHash) {
		return NewDomainErrorWithCode(ErrUnauthorized, "INVALID_CREDENTIALS", "Current password is incorrect.")
	}

	newHash, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	tx, err := s.db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin password update transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	qtx := s.queries.WithTx(tx)

	if err := qtx.UpdateAccountPassword(ctx, db.UpdateAccountPasswordParams{
		IDAccount:    parsedUUID,
		PasswordHash: newHash,
	}); err != nil {
		return MapDBError(err, "Failed to update password.")
	}

	if _, err := qtx.IncrementAccountTokenVersion(ctx, parsedUUID); err != nil {
		return fmt.Errorf("failed to increment token version: %w", err)
	}

	if err := qtx.RevokeAllRefreshTokensForAccount(ctx, parsedUUID); err != nil {
		return fmt.Errorf("failed to revoke refresh tokens: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit password update: %w", err)
	}

	return nil
}

func (s *identityService) GetPublicProfileByUsername(ctx context.Context, viewerID *string, username string) (*models.PublicUserProfile, error) {
	username = strings.TrimSpace(strings.ToLower(username))
	if username == "" {
		return nil, NewDomainErrorWithCode(ErrInvalidInput, "INVALID_INPUT", "Username is required.")
	}

	userRow, err := s.queries.GetUserProfileByUsername(ctx, username)
	if err != nil {
		return nil, MapDBError(err, "User not found.")
	}

	isOwner := false
	if viewerID != nil && strings.TrimSpace(*viewerID) != "" {
		if strings.EqualFold(userRow.IDUser.String(), strings.TrimSpace(*viewerID)) {
			isOwner = true
		}
	}

	// 1. Achievements
	achRows, err := s.queries.ListUserAchievements(ctx, userRow.IDUser)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("failed to load user achievements: %w", err)
	}
	unlockedBadges := make([]models.AchievementStatus, 0, len(achRows))
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

	// 2. Public decks
	deckRows, err := s.queries.ListUserPublicDecks(ctx, userRow.IDUser)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("failed to load user public decks: %w", err)
	}
	publicDecks := make([]models.Deck, 0, len(deckRows))
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

	// 3. Public collections
	colRows, err := s.queries.ListUserPublicCollections(ctx, userRow.IDUser)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("failed to load user public collections: %w", err)
	}
	publicCollections := make([]models.Collection, 0, len(colRows))
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

	// 4. Arena stats
	matchesWon := 0
	podiumCount := 0
	stats, err := s.queries.GetUserArenaStats(ctx, UUIDToPgtype(userRow.IDUser))
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("failed to load user arena stats: %w", err)
		}
	} else {
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

	// Owner-only private fields (RN01)
	if isOwner {
		email := userRow.Email
		publicProfile.Email = &email

		privCount, err := s.queries.CountUserPrivateDecks(ctx, userRow.IDUser)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("failed to count private decks: %w", err)
		}
		privInt := int(privCount)
		publicProfile.PrivateDecksCount = &privInt

		dueCount, err := s.queries.CountDueStudySessionCardsAll(ctx, userRow.IDUser)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("failed to count due review cards: %w", err)
		}
		dueInt := int(dueCount)
		publicProfile.DueReviewsToday = &dueInt
	}

	return publicProfile, nil
}
