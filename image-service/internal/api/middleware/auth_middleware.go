package middleware

import (
	"context"
	"net/http"
	"strings"

	"palabrario/pkg/auth"
	"palabrario/pkg/response"
)

// AuthMiddleware encapsula la validación de tokens JWT para los endpoints de image-service.
type AuthMiddleware struct {
	tokenService *auth.TokenService
}

// NewAuthMiddleware inicializa el middleware de autenticación.
func NewAuthMiddleware(tokenService *auth.TokenService) *AuthMiddleware {
	return &AuthMiddleware{tokenService: tokenService}
}

// RequireAuth valida que la petición contenga un encabezado 'Authorization: Bearer <token>' válido.
func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			response.Error(w, http.StatusUnauthorized, "token de autorización no proporcionado o con formato inválido")
			return
		}

		tokenStr := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		if tokenStr == "" {
			response.Error(w, http.StatusUnauthorized, "token de autorización vacío")
			return
		}

		claims, err := m.tokenService.VerifyToken(tokenStr)
		if err != nil {
			response.Error(w, http.StatusUnauthorized, "token inválido o expirado: "+err.Error())
			return
		}

		ctx := context.WithValue(r.Context(), auth.UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
