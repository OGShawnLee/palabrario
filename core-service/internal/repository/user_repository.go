package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"palabrario/core-service/internal/database"
	"palabrario/core-service/internal/db"
	"palabrario/pkg/models"
)

type UserRepository struct {
	db      *database.DB
	queries *db.Queries
}

func NewUserRepository(databaseInstance *database.DB) *UserRepository {
	return &UserRepository{
		db:      databaseInstance,
		queries: db.New(databaseInstance.Pool),
	}
}

func UUIDToString(u pgtype.UUID) string {
	if !u.Valid {
		return ""
	}
	src := u.Bytes
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		src[0:4], src[4:6], src[6:8], src[8:10], src[10:16])
}

func parseUUID(s string) (pgtype.UUID, error) {
	var u pgtype.UUID
	err := u.Scan(s)
	return u, err
}

func textToPtr(t pgtype.Text) *string {
	if !t.Valid {
		return nil
	}
	s := t.String
	return &s
}

func (r *UserRepository) CreateUser(ctx context.Context, email, passwordHash, name, userName, anchorDialect string) (*models.UserProfile, error) {
	if anchorDialect == "" {
		anchorDialect = "MX"
	}

	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("error iniciando transacción: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	qtx := r.queries.WithTx(tx)

	// 1. Crear cuenta tipada con sqlc
	account, err := qtx.CreateAccount(ctx, db.CreateAccountParams{
		Email:        email,
		PasswordHash: passwordHash,
	})
	if err != nil {
		return nil, fmt.Errorf("error al crear cuenta: %w", err)
	}

	// 2. Crear perfil tipado con sqlc
	userProfile, err := qtx.CreateUserProfile(ctx, db.CreateUserProfileParams{
		IDUser:        account.IDAccount,
		Name:          name,
		UserName:      userName,
		AnchorDialect: anchorDialect,
	})
	if err != nil {
		return nil, fmt.Errorf("error al crear perfil: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("error confirmando transacción: %w", err)
	}

	return &models.UserProfile{
		IDUser:        UUIDToString(account.IDAccount),
		Email:         account.Email,
		Name:          userProfile.Name,
		UserName:      userProfile.UserName,
		AvatarURL:     textToPtr(userProfile.AvatarUrl),
		AnchorDialect: userProfile.AnchorDialect,
		TotalXP:       int(userProfile.TotalXp),
		CurrentLevel:  int(userProfile.CurrentLevel),
		CurrentStreak: int(userProfile.CurrentStreak),
		LongestStreak: int(userProfile.LongestStreak),
		Role:          models.AccountRole(account.Role),
		CreatedAt:     account.CreatedAt.Time,
		UpdatedAt:     account.UpdatedAt.Time,
	}, nil
}

func (r *UserRepository) GetAccountByEmail(ctx context.Context, email string) (*db.Account, error) {
	acc, err := r.queries.GetAccountByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("credenciales inválidas")
		}
		return nil, fmt.Errorf("error consultando cuenta: %w", err)
	}
	return &acc, nil
}

func (r *UserRepository) GetProfileByID(ctx context.Context, userID string) (*models.UserProfile, error) {
	uuidVal, err := parseUUID(userID)
	if err != nil {
		return nil, fmt.Errorf("id de usuario inválido: %w", err)
	}

	row, err := r.queries.GetUserProfileByID(ctx, uuidVal)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("usuario no encontrado")
		}
		return nil, fmt.Errorf("error consultando perfil: %w", err)
	}

	return &models.UserProfile{
		IDUser:        UUIDToString(row.IDUser),
		Email:         row.Email,
		Name:          row.Name,
		UserName:      row.UserName,
		AvatarURL:     textToPtr(row.AvatarUrl),
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
