package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"palabrario/core-service/internal/service"
	"palabrario/pkg/auth"
	"palabrario/pkg/models"
	"palabrario/pkg/response"
)

type UserHandler struct {
	identityService service.IdentityService
}

func NewUserHandler(identityService service.IdentityService) *UserHandler {
	return &UserHandler{
		identityService: identityService,
	}
}

// UpdateProfile maneja PUT /api/v1/users/me (CU01, RN01)
func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	if !ok || claims == nil {
		response.Error(w, http.StatusUnauthorized, "No autorizado")
		return
	}

	var req models.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Cuerpo de petición inválido")
		return
	}

	profile, err := h.identityService.UpdateProfile(r.Context(), claims.UserID, req)
	if err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al actualizar el perfil")
		response.Error(w, status, msg)
		return
	}

	response.JSON(w, http.StatusOK, profile)
}

// ChangePassword maneja PUT /api/v1/users/me/password (CU01, RN01)
func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	if !ok || claims == nil {
		response.Error(w, http.StatusUnauthorized, "No autorizado")
		return
	}

	var req models.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Cuerpo de petición inválido")
		return
	}

	if err := h.identityService.ChangePassword(r.Context(), claims.UserID, req); err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al cambiar la contraseña")
		response.Error(w, status, msg)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetPublicProfile maneja GET /api/v1/users/u/{username} (CU01 / UI)
func (h *UserHandler) GetPublicProfile(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	if username == "" {
		response.Error(w, http.StatusBadRequest, "Nombre de usuario requerido")
		return
	}

	var viewerID *string
	if claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims); ok && claims != nil {
		viewerID = &claims.UserID
	}

	profile, err := h.identityService.GetPublicProfileByUsername(r.Context(), viewerID, username)
	if err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error consultando perfil público")
		response.Error(w, status, msg)
		return
	}

	response.JSON(w, http.StatusOK, profile)
}
