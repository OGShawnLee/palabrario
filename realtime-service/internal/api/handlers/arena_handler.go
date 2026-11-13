package handlers

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"palabrario/pkg/auth"
	"palabrario/pkg/response"
	"palabrario/realtime-service/internal/service"
)

// ArenaHandler manages WebSocket authentication and subscription for Arena rooms.
type ArenaHandler struct {
	hub          *service.ArenaHub
	tokenService *auth.TokenService
}

// NewArenaHandler creates a new ArenaHandler with its injected dependencies.
func NewArenaHandler(hub *service.ArenaHub, tokenService *auth.TokenService) *ArenaHandler {
	return &ArenaHandler{
		hub:          hub,
		tokenService: tokenService,
	}
}

// HandleWS validates the player's JWT credentials (Authorization header primary, ?token= fallback)
// and upgrades the HTTP connection to a WebSocket for the requested room.
func (h *ArenaHandler) HandleWS(w http.ResponseWriter, r *http.Request) {
	roomCode := strings.TrimSpace(chi.URLParam(r, "roomCode"))
	if roomCode == "" {
		response.ErrorWithCode(w, http.StatusBadRequest, "MISSING_ROOM_CODE", "Room code is required")
		return
	}

	tokenStr, present, err := auth.ExtractBearerToken(r)
	if present {
		if err != nil {
			response.ErrorWithCode(w, http.StatusUnauthorized, "INVALID_TOKEN_FORMAT", "Invalid authorization header format")
			return
		}
	} else {
		tokenStr = strings.TrimSpace(r.URL.Query().Get("token"))
	}

	if tokenStr == "" {
		response.ErrorWithCode(w, http.StatusUnauthorized, "MISSING_TOKEN", "JWT token is required in Authorization header or ?token= query parameter")
		return
	}

	claims, err := h.tokenService.VerifyToken(tokenStr)
	if err != nil {
		response.ErrorWithCode(w, http.StatusUnauthorized, "INVALID_TOKEN", "Invalid or expired token")
		return
	}

	h.hub.HandleWS(w, r, roomCode, claims)
}
