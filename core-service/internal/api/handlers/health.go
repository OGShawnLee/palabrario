package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"palabrario/core-service/internal/database"
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
	if err := h.db.Pool.Ping(r.Context()); err != nil {
		dbStatus = "disconnected"
	}

	w.Header().Set("Content-Type", "application/json")
	if dbStatus != "connected" {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	resp := HealthResponse{
		Status:    "ok",
		Service:   "core-service",
		Database:  dbStatus,
		Timestamp: time.Now().UTC(),
	}

	_ = json.NewEncoder(w).Encode(resp)
}
