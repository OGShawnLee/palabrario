package middleware

import (
	"context"
	"net/http"
	"strings"

	"palabrario/core-service/internal/auth"
	"palabrario/pkg/models"
)

type AuthMiddleware struct {
	tokenService *auth.TokenService
}

func NewAuthMiddleware(tokenService *auth.TokenService) *AuthMiddleware {
	return &AuthMiddleware{tokenService: tokenService}
}

func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"error":"Encabezado de autorización ausente"}`, http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			http.Error(w, `{"error":"Formato de token inválido. Use Bearer <token>"}`, http.StatusUnauthorized)
			return
		}

		claims, err := m.tokenService.VerifyToken(parts[1])
		if err != nil {
			http.Error(w, `{"error":"Token inválido o expirado"}`, http.StatusUnauthorized)
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
				http.Error(w, `{"error":"No autenticado"}`, http.StatusUnauthorized)
				return
			}

			allowed := false
			for _, role := range roles {
				if claims.Role == role {
					allowed = true
					break
				}
			}

			if !allowed {
				http.Error(w, `{"error":"Permisos insuficientes para acceder a este recurso"}`, http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
