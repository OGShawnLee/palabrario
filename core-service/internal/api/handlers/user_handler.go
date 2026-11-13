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

// UpdateProfile handles PUT /api/v1/users/me (CU01, RN01)
func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	if !ok || claims == nil {
		response.ErrorWithCode(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required.")
		return
	}

	var req models.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ErrorWithCode(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body.")
		return
	}

	profile, err := h.identityService.UpdateProfile(r.Context(), claims.UserID, req)
	if err != nil {
		status, code, msg := service.HTTPStatusCodeAndMessage(err, "INTERNAL_ERROR", "Failed to update profile.")
		response.ErrorWithCode(w, status, code, msg)
		return
	}

	response.JSON(w, http.StatusOK, profile)
}

// ChangePassword handles PUT /api/v1/users/me/password (CU01, RN01)
func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	if !ok || claims == nil {
		response.ErrorWithCode(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required.")
		return
	}

	var req models.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ErrorWithCode(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body.")
		return
	}

	if err := h.identityService.ChangePassword(r.Context(), claims.UserID, req); err != nil {
		status, code, msg := service.HTTPStatusCodeAndMessage(err, "INTERNAL_ERROR", "Failed to change password.")
		response.ErrorWithCode(w, status, code, msg)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetPublicProfile handles GET /api/v1/users/u/{username} (CU01 / UI)
func (h *UserHandler) GetPublicProfile(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	if username == "" {
		response.ErrorWithCode(w, http.StatusBadRequest, "INVALID_INPUT", "Username is required.")
		return
	}

	var viewerID *string
	if claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims); ok && claims != nil {
		viewerID = &claims.UserID
	}

	profile, err := h.identityService.GetPublicProfileByUsername(r.Context(), viewerID, username)
	if err != nil {
		status, code, msg := service.HTTPStatusCodeAndMessage(err, "INTERNAL_ERROR", "Failed to load public profile.")
		response.ErrorWithCode(w, status, code, msg)
		return
	}

	response.JSON(w, http.StatusOK, profile)
}
