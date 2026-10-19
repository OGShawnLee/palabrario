package handlers

import (
	"net/http"
	"time"

	"palabrario/core-service/internal/database"
	"palabrario/pkg/response"
)

type HealthResponse struct {
	Status    string    `json:"status"`
	Service   string    `json:"service"`
	Database  string    `json:"database"`
	Timestamp time.Time `json:"timestamp"`
}

type HealthHandler struct {
	db *database.DB
}

func NewHealthHandler(db *database.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	dbStatus := "connected"
	statusCode := http.StatusOK
	if err := h.db.Pool.Ping(r.Context()); err != nil {
		dbStatus = "disconnected"
		statusCode = http.StatusServiceUnavailable
	}

	response.JSON(w, statusCode, HealthResponse{
		Status:    "ok",
		Service:   "core-service",
		Database:  dbStatus,
		Timestamp: time.Now().UTC(),
	})
}
