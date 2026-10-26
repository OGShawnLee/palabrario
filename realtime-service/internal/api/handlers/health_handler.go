package handlers

import (
	"net/http"
	"time"

	"palabrario/pkg/response"
	"palabrario/realtime-service/internal/service"
)

// HealthResponse estructura la respuesta del probe de salud de realtime-service.
type HealthResponse struct {
	Status      string    `json:"status"`
	Service     string    `json:"service"`
	ActiveRooms int       `json:"active_rooms"`
	Timestamp   time.Time `json:"timestamp"`
}

// HealthHandler administra los endpoints de salud y bienvenida del servicio.
type HealthHandler struct {
	hub *service.ArenaHub
}

// NewHealthHandler inicializa un nuevo HealthHandler con el ArenaHub inyectado.
func NewHealthHandler(hub *service.ArenaHub) *HealthHandler {
	return &HealthHandler{hub: hub}
}

// HealthCheck responde con el estado operacional y conteo de salas activas en RAM.
func (h *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	activeRooms := 0
	if h.hub != nil {
		activeRooms = h.hub.ActiveRoomsCount()
	}

	response.JSON(w, http.StatusOK, HealthResponse{
		Status:      "ok",
		Service:     "realtime-service",
		ActiveRooms: activeRooms,
		Timestamp:   time.Now().UTC(),
	})
}

// Root responde con el mensaje canónico de bienvenida del servicio.
func (h *HealthHandler) Root(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]string{
		"message": "Bienvenido a Palabrario Realtime Service (Golang)",
	})
}
