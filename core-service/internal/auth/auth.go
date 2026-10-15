package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"palabrario/pkg/models"
)

type contextKey string

const UserContextKey contextKey = "user_claims"

type Claims struct {
	UserID   string             `json:"user_id"`
	Email    string             `json:"email"`
	UserName string             `json:"user_name"`
	Role     models.AccountRole `json:"role"`
	jwt.RegisteredClaims
}

type TokenService struct {
	secretKey []byte
	duration  time.Duration
}

func NewTokenService(secretKey string) *TokenService {
	return &TokenService{
		secretKey: []byte(secretKey),
		duration:  72 * time.Hour, // 3 días de vigencia
	}
}

func (s *TokenService) GenerateToken(user *models.UserProfile) (string, error) {
	now := time.Now().UTC()
	claims := Claims{
		UserID:   user.IDUser,
		Email:    user.Email,
		UserName: user.UserName,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.duration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "palabrario-core-service",
			Subject:   user.IDUser,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", fmt.Errorf("error firmando token JWT: %w", err)
	}

	return tokenString, nil
}

func (s *TokenService) VerifyToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de firma no válido: %v", t.Header["alg"])
		}
		return s.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("token inválido o expirado: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("reclamaciones de token inválidas")
	}

	return claims, nil
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error generando hash de contraseña: %w", err)
	}
	return string(hash), nil
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
