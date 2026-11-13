package middleware

import (
	"context"
	"errors"
	"net/http"
	"slices"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"palabrario/core-service/internal/db"
	"palabrario/pkg/auth"
	"palabrario/pkg/models"
	"palabrario/pkg/response"
)

type AuthStateStore interface {
	GetAuthAccountStateByID(ctx context.Context, idAccount uuid.UUID) (db.GetAuthAccountStateByIDRow, error)
}

type AuthMiddleware struct {
	tokenService *auth.TokenService
	store        AuthStateStore
}

func NewAuthMiddleware(tokenService *auth.TokenService, store AuthStateStore) *AuthMiddleware {
	return &AuthMiddleware{
		tokenService: tokenService,
		store:        store,
	}
}

func (m *AuthMiddleware) verifyAccessToken(r *http.Request) (*auth.Claims, string, string) {
	tokenStr, present, err := auth.ExtractBearerToken(r)
	if !present {
		return nil, "MISSING_TOKEN", "Authorization header is required."
	}
	if err != nil {
		return nil, "INVALID_TOKEN_FORMAT", "Invalid authorization header format. Expected Bearer <token>."
	}

	claims, err := m.tokenService.VerifyToken(tokenStr)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, "TOKEN_EXPIRED", "Access token has expired."
		}
		return nil, "INVALID_TOKEN", "Invalid or expired token."
	}

	return claims, "", ""
}

func requiresStrictDBCheck(r *http.Request) bool {
	if r.Method != http.MethodGet && r.Method != http.MethodHead && r.Method != http.MethodOptions {
		return true
	}
	path := r.URL.Path
	return strings.HasPrefix(path, "/api/v1/auth/") ||
		path == "/api/v1/auth" ||
		strings.HasPrefix(path, "/api/v1/governance/") ||
		path == "/api/v1/governance"
}

func (m *AuthMiddleware) enforceActiveAccount(w http.ResponseWriter, r *http.Request, claims *auth.Claims) bool {
	if m.store == nil {
		return true
	}

	userUUID, err := uuid.Parse(strings.TrimSpace(claims.UserID))
	if err != nil {
		response.ErrorWithCode(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid user session.")
		return false
	}

	state, err := m.store.GetAuthAccountStateByID(r.Context(), userUUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.ErrorWithCode(w, http.StatusUnauthorized, "UNAUTHORIZED", "User session is no longer valid.")
			return false
		}
		response.ErrorWithCode(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to verify account state.")
		return false
	}

	if !state.IsActive {
		response.ErrorWithCode(w, http.StatusForbidden, "ACCOUNT_SUSPENDED", "This account has been suspended by an administrator.")
		return false
	}

	if int(state.TokenVersion) != claims.TokenVersion {
		response.ErrorWithCode(w, http.StatusUnauthorized, "TOKEN_REVOKED", "Session is no longer valid. Please sign in again.")
		return false
	}

	// Sync authoritative role from DB before downstream RequireRole executes
	claims.Role = models.AccountRole(state.Role)
	return true
}

func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, code, msg := m.verifyAccessToken(r)
		if claims == nil {
			response.ErrorWithCode(w, http.StatusUnauthorized, code, msg)
			return
		}

		if requiresStrictDBCheck(r) {
			if !m.enforceActiveAccount(w, r, claims) {
				return
			}
		}

		ctx := context.WithValue(r.Context(), auth.UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// OptionalAuth injects user claims if a valid Bearer token is provided,
// continues anonymously if no Authorization header is sent, and rejects
// malformed or expired Authorization headers with 401 Unauthorized.
func (m *AuthMiddleware) OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr, present, err := auth.ExtractBearerToken(r)
		if !present {
			next.ServeHTTP(w, r)
			return
		}
		if err != nil {
			response.ErrorWithCode(w, http.StatusUnauthorized, "INVALID_TOKEN_FORMAT", "Invalid authorization header format. Expected Bearer <token>.")
			return
		}

		claims, err := m.tokenService.VerifyToken(tokenStr)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				response.ErrorWithCode(w, http.StatusUnauthorized, "TOKEN_EXPIRED", "Access token has expired.")
				return
			}
			response.ErrorWithCode(w, http.StatusUnauthorized, "INVALID_TOKEN", "Invalid or expired token.")
			return
		}

		ctx := context.WithValue(r.Context(), auth.UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *AuthMiddleware) RequireRole(roles ...models.AccountRole) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
			if !ok || claims == nil {
				response.ErrorWithCode(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required.")
				return
			}

			if !slices.Contains(roles, claims.Role) {
				response.ErrorWithCode(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions to access this resource.")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireInternalService validates that the request carries an internal S2S JWT.
func (m *AuthMiddleware) RequireInternalService(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr, present, err := auth.ExtractBearerToken(r)
		if !present {
			response.ErrorWithCode(w, http.StatusUnauthorized, "MISSING_TOKEN", "Authorization header is required.")
			return
		}
		if err != nil {
			response.ErrorWithCode(w, http.StatusUnauthorized, "INVALID_TOKEN_FORMAT", "Invalid authorization header format. Expected Bearer <token>.")
			return
		}

		claims, err := m.tokenService.VerifyInternalToken(tokenStr)
		if err != nil {
			response.ErrorWithCode(w, http.StatusForbidden, "FORBIDDEN", "Access restricted to authorized internal services.")
			return
		}

		ctx := context.WithValue(r.Context(), auth.InternalServiceContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
