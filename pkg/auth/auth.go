package auth

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"palabrario/pkg/models"
)

type contextKey string

const (
	UserContextKey            contextKey = "user_claims"
	InternalServiceContextKey contextKey = "internal_service_claims"
)

type Claims struct {
	UserID        string             `json:"user_id"`
	Email         string             `json:"email"`
	UserName      string             `json:"user_name"`
	AnchorDialect string             `json:"anchor_dialect"`
	Role          models.AccountRole `json:"role"`
	jwt.RegisteredClaims
}

type InternalServiceClaims struct {
	ServiceName string `json:"service_name"`
	MatchScope  string `json:"match_scope"`
	jwt.RegisteredClaims
}

type TokenService struct {
	secretKey    []byte
	s2sSecretKey []byte
	duration     time.Duration
}

func NewTokenService(secretKey string, duration time.Duration) *TokenService {
	if duration <= 0 {
		duration = 72 * time.Hour
	}
	sum := sha256.Sum256([]byte(secretKey + ":palabrario-s2s-internal"))
	return &TokenService{
		secretKey:    []byte(secretKey),
		s2sSecretKey: sum[:],
		duration:     duration,
	}
}

func (s *TokenService) GenerateToken(user *models.UserProfile) (string, error) {
	now := time.Now().UTC()
	dialect := user.AnchorDialect
	if dialect == "" {
		dialect = "MX"
	}

	claims := Claims{
		UserID:        user.IDUser,
		Email:         user.Email,
		UserName:      user.UserName,
		AnchorDialect: dialect,
		Role:          user.Role,
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

	if claims.Issuer != "palabrario-core-service" {
		return nil, errors.New("emisor de token de usuario inválido")
	}

	return claims, nil
}

func (s *TokenService) GenerateInternalToken(serviceName string, matchScope string) (string, error) {
	now := time.Now().UTC()

	claims := InternalServiceClaims{
		ServiceName: serviceName,
		MatchScope:  matchScope,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(30 * time.Second)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "palabrario-internal-s2s",
			Subject:   serviceName,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.s2sSecretKey)
	if err != nil {
		return "", fmt.Errorf("error firmando token interno S2S: %w", err)
	}

	return tokenString, nil
}

func (s *TokenService) VerifyInternalToken(tokenString string) (*InternalServiceClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &InternalServiceClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de firma no válido: %v", t.Header["alg"])
		}
		return s.s2sSecretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("token interno inválido o expirado: %w", err)
	}

	claims, ok := token.Claims.(*InternalServiceClaims)
	if !ok || !token.Valid {
		return nil, errors.New("reclamaciones de token interno inválidas")
	}

	if claims.Issuer != "palabrario-internal-s2s" {
		return nil, errors.New("emisor de token interno inválido")
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
