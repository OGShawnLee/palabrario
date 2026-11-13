package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
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

	DefaultAccessTokenTTL = 15 * time.Minute
	MaxBcryptPasswordLen  = 72
	TokenTypeAccess       = "access"
)

type Claims struct {
	UserID        string             `json:"user_id"`
	Email         string             `json:"email"`
	UserName      string             `json:"user_name"`
	AnchorDialect string             `json:"anchor_dialect"`
	Role          models.AccountRole `json:"role"`
	TokenVersion  int                `json:"token_version"`
	TokenType     string             `json:"token_type,omitempty"`
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
		duration = DefaultAccessTokenTTL
	}
	sum := sha256.Sum256([]byte(secretKey + ":palabrario-s2s-internal"))
	return &TokenService{
		secretKey:    []byte(secretKey),
		s2sSecretKey: sum[:],
		duration:     duration,
	}
}

// AccessTTL returns the configured access token lifetime.
func (s *TokenService) AccessTTL() time.Duration {
	return s.duration
}

// GenerateToken issues an access JWT with default TokenVersion = 1 for backwards compatibility.
func (s *TokenService) GenerateToken(user *models.UserProfile) (string, error) {
	return s.GenerateAccessToken(user, 1)
}

// GenerateAccessToken issues a short-lived user access JWT carrying TokenVersion and TokenType = "access".
func (s *TokenService) GenerateAccessToken(user *models.UserProfile, tokenVersion int) (string, error) {
	if user == nil {
		return "", errors.New("user profile is required")
	}

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
		TokenVersion:  tokenVersion,
		TokenType:     TokenTypeAccess,
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
		return "", fmt.Errorf("error signing JWT token: %w", err)
	}

	return tokenString, nil
}

func (s *TokenService) VerifyToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid or expired token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	if claims.Issuer != "palabrario-core-service" {
		return nil, errors.New("invalid user token issuer")
	}

	if claims.TokenType != "" && claims.TokenType != TokenTypeAccess {
		return nil, errors.New("invalid token type")
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
		return "", fmt.Errorf("error signing internal S2S token: %w", err)
	}

	return tokenString, nil
}

func (s *TokenService) VerifyInternalToken(tokenString string) (*InternalServiceClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &InternalServiceClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.s2sSecretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid or expired internal token: %w", err)
	}

	claims, ok := token.Claims.(*InternalServiceClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid internal token claims")
	}

	if claims.Issuer != "palabrario-internal-s2s" {
		return nil, errors.New("invalid internal token issuer")
	}

	return claims, nil
}

// GenerateRefreshToken generates a 32-byte cryptographically random opaque refresh token
// (base64url-encoded without padding) together with its 64-character SHA-256 hex digest.
func GenerateRefreshToken() (raw string, hash string, err error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", "", fmt.Errorf("error generating refresh token: %w", err)
	}
	raw = base64.RawURLEncoding.EncodeToString(buf)
	hash = HashRefreshToken(raw)
	return raw, hash, nil
}

// HashRefreshToken computes the deterministic 64-character lowercase SHA-256 hex digest of a raw refresh token.
func HashRefreshToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func HashPassword(password string) (string, error) {
	if len([]byte(password)) > MaxBcryptPasswordLen {
		return "", fmt.Errorf("password must not exceed %d bytes", MaxBcryptPasswordLen)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error hashing password: %w", err)
	}
	return string(hash), nil
}

func CheckPassword(password, hash string) bool {
	if len([]byte(password)) > MaxBcryptPasswordLen {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
