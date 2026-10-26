package handlers

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"palabrario/pkg/auth"
	"palabrario/pkg/response"
	"palabrario/realtime-service/internal/service"
)

// ArenaHandler administra la conexión y suscripción de clientes WebSocket a la Arena.
type ArenaHandler struct {
	hub          *service.ArenaHub
	tokenService *auth.TokenService
}

// NewArenaHandler crea una nueva instancia de ArenaHandler con sus dependencias inyectadas.
func NewArenaHandler(hub *service.ArenaHub, tokenService *auth.TokenService) *ArenaHandler {
	return &ArenaHandler{
		hub:          hub,
		tokenService: tokenService,
	}
}

// HandleWS valida las credenciales del jugador y eleva la conexión HTTP a WebSocket para la sala solicitada.
func (h *ArenaHandler) HandleWS(w http.ResponseWriter, r *http.Request) {
	roomCode := strings.TrimSpace(chi.URLParam(r, "roomCode"))
	if roomCode == "" {
		response.Error(w, http.StatusBadRequest, "El código de sala es requerido")
		return
	}

	tokenStr := strings.TrimSpace(r.URL.Query().Get("token"))
	if tokenStr == "" {
		authHeader := r.Header.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenStr = strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		}
	}

	if tokenStr == "" {
		response.Error(w, http.StatusUnauthorized, "Token JWT requerido en ?token= o Authorization")
		return
	}

	claims, err := h.tokenService.VerifyToken(tokenStr)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "Token inválido o expirado")
		return
	}

	h.hub.HandleWS(w, r, roomCode, claims)
}
