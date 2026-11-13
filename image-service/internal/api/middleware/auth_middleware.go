package middleware

import (
	"context"
	"net/http"

	"palabrario/pkg/auth"
	"palabrario/pkg/response"
)

// AuthMiddleware validates user JWT tokens for image-service endpoints.
type AuthMiddleware struct {
	tokenService *auth.TokenService
}

// NewAuthMiddleware initializes the authentication middleware.
func NewAuthMiddleware(tokenService *auth.TokenService) *AuthMiddleware {
	return &AuthMiddleware{tokenService: tokenService}
}

// RequireAuth validates that the request includes a valid 'Authorization: Bearer <token>' header.
func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr, present, err := auth.ExtractBearerToken(r)
		if !present {
			response.ErrorWithCode(w, http.StatusUnauthorized, "MISSING_TOKEN", "Authorization header is required")
			return
		}
		if err != nil {
			response.ErrorWithCode(w, http.StatusUnauthorized, "INVALID_TOKEN_FORMAT", "Invalid authorization header format")
			return
		}

		claims, err := m.tokenService.VerifyToken(tokenStr)
		if err != nil {
			response.ErrorWithCode(w, http.StatusUnauthorized, "INVALID_TOKEN", "Invalid or expired token")
			return
		}

		ctx := context.WithValue(r.Context(), auth.UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
