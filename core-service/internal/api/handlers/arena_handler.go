package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"palabrario/core-service/internal/service"
	"palabrario/pkg/auth"
	"palabrario/pkg/models"
	"palabrario/pkg/response"
)

type ArenaHandler struct {
	arenaService service.ArenaCatalogService
}

func NewArenaHandler(arenaService service.ArenaCatalogService) *ArenaHandler {
	return &ArenaHandler{arenaService: arenaService}
}

// PrepareMatch maneja POST /api/v1/arena/matches/prepare (RequireAuth)
func (h *ArenaHandler) PrepareMatch(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	if !ok || claims == nil {
		response.Error(w, http.StatusUnauthorized, "No autenticado")
		return
	}

	var req models.CreateMatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Cuerpo de solicitud JSON inválido")
		return
	}

	pkg, err := h.arenaService.PrepareMatch(r.Context(), claims.UserID, claims.Role, req)
	if err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al preparar la partida de Arena")
		response.Error(w, status, msg)
		return
	}

	response.JSON(w, http.StatusCreated, pkg)
}

// ListHistory maneja GET /api/v1/arena/matches/history (RequireAuth)
func (h *ArenaHandler) ListHistory(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	if !ok || claims == nil {
		response.Error(w, http.StatusUnauthorized, "No autenticado")
		return
	}

	limit := queryIntOrDefault(r, "limit", 20, 1, 50)

	history, err := h.arenaService.ListMatchHistory(r.Context(), claims.UserID, limit)
	if err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al consultar el historial de partidas")
		response.Error(w, status, msg)
		return
	}

	response.JSON(w, http.StatusOK, history)
}

// GetReplay maneja GET /api/v1/arena/matches/{matchId} (RequireAuth)
func (h *ArenaHandler) GetReplay(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	if !ok || claims == nil {
		response.Error(w, http.StatusUnauthorized, "No autenticado")
		return
	}

	matchID := strings.TrimSpace(chi.URLParam(r, "matchId"))
	if matchID == "" {
		response.Error(w, http.StatusBadRequest, "El identificador de partida es requerido")
		return
	}

	replay, err := h.arenaService.GetMatchReplay(r.Context(), claims.UserID, claims.Role, matchID)
	if err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al consultar el replay de la partida")
		response.Error(w, status, msg)
		return
	}

	response.JSON(w, http.StatusOK, replay)
}

// GetPackageByRoomCode maneja GET /api/v1/arena/matches/{roomCode}/package (RequireInternalService)
func (h *ArenaHandler) GetPackageByRoomCode(w http.ResponseWriter, r *http.Request) {
	roomCode := strings.TrimSpace(chi.URLParam(r, "roomCode"))
	if roomCode == "" {
		response.Error(w, http.StatusBadRequest, "El código de sala es requerido")
		return
	}

	pkg, err := h.arenaService.GetMatchPackageByRoomCode(r.Context(), roomCode)
	if err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al consultar el paquete de partida")
		response.Error(w, status, msg)
		return
	}

	response.JSON(w, http.StatusOK, pkg)
}

// ConsolidateMatch maneja POST /api/v1/arena/matches/consolidate (RequireInternalService)
func (h *ArenaHandler) ConsolidateMatch(w http.ResponseWriter, r *http.Request) {
	var req models.ConsolidateMatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Cuerpo de solicitud JSON inválido")
		return
	}

	summary, err := h.arenaService.ConsolidateMatch(r.Context(), req)
	if err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al consolidar la partida")
		response.Error(w, status, msg)
		return
	}

	response.JSON(w, http.StatusOK, summary)
}
