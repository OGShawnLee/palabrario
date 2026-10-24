package middleware

import (
	"context"
	"net/http"
	"slices"
	"strings"

	"palabrario/pkg/auth"
	"palabrario/pkg/models"
	"palabrario/pkg/response"
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
			response.Error(w, http.StatusUnauthorized, "Encabezado de autorización ausente")
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			response.Error(w, http.StatusUnauthorized, "Formato de token inválido. Use Bearer <token>")
			return
		}

		claims, err := m.tokenService.VerifyToken(parts[1])
		if err != nil {
			response.Error(w, http.StatusUnauthorized, "Token inválido o expirado")
			return
		}

		ctx := context.WithValue(r.Context(), auth.UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// OptionalAuth inyecta los Claims en el contexto si el cliente envía un token Bearer válido,
// pero permite continuar la petición como usuario anónimo si no hay token.
func (m *AuthMiddleware) OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
				if claims, err := m.tokenService.VerifyToken(parts[1]); err == nil {
					ctx := context.WithValue(r.Context(), auth.UserContextKey, claims)
					r = r.WithContext(ctx)
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}

func (m *AuthMiddleware) RequireRole(roles ...models.AccountRole) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
			if !ok || claims == nil {
				response.Error(w, http.StatusUnauthorized, "No autenticado")
				return
			}

			if !slices.Contains(roles, claims.Role) {
				response.Error(w, http.StatusForbidden, "Permisos insuficientes para acceder a este recurso")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireInternalService valida que la petición contenga un JWT firmado con la sub-llave S2S
// emitido para comunicación interna entre servicios de Palabrario.
func (m *AuthMiddleware) RequireInternalService(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response.Error(w, http.StatusUnauthorized, "Encabezado de autorización ausente")
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			response.Error(w, http.StatusUnauthorized, "Formato de token inválido. Use Bearer <token>")
			return
		}

		claims, err := m.tokenService.VerifyInternalToken(parts[1])
		if err != nil {
			response.Error(w, http.StatusForbidden, "Acceso restringido exclusivamente a servicios internos autorizados")
			return
		}

		ctx := context.WithValue(r.Context(), auth.InternalServiceContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
