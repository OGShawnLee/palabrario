package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"palabrario/core-service/internal/service"
	"palabrario/pkg/auth"
	"palabrario/pkg/models"
	"palabrario/pkg/response"
)

type AuthHandler struct {
	identityService service.IdentityService
}

func NewAuthHandler(identityService service.IdentityService) *AuthHandler {
	return &AuthHandler{
		identityService: identityService,
	}
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var req models.SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ErrorWithCode(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body.")
		return
	}

	authResp, err := h.identityService.SignUp(r.Context(), req)
	if err != nil {
		status, code, msg := service.HTTPStatusCodeAndMessage(err, "SIGN_UP_FAILED", "Failed to register account.")
		response.ErrorWithCode(w, status, code, msg)
		return
	}

	response.JSON(w, http.StatusCreated, authResp)
}

func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var req models.SignInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ErrorWithCode(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body.")
		return
	}

	authResp, err := h.identityService.SignIn(r.Context(), req)
	if err != nil {
		status, code, msg := service.HTTPStatusCodeAndMessage(err, "SIGN_IN_FAILED", "Failed to sign in.")
		response.ErrorWithCode(w, status, code, msg)
		return
	}

	response.JSON(w, http.StatusOK, authResp)
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req models.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ErrorWithCode(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body.")
		return
	}

	authResp, err := h.identityService.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		status, code, msg := service.HTTPStatusCodeAndMessage(err, "REFRESH_FAILED", "Failed to refresh session.")
		response.ErrorWithCode(w, status, code, msg)
		return
	}

	response.JSON(w, http.StatusOK, authResp)
}

func (h *AuthHandler) SignOut(w http.ResponseWriter, r *http.Request) {
	var req models.SignOutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil && !errors.Is(err, io.EOF) {
		response.ErrorWithCode(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body.")
		return
	}

	if err := h.identityService.SignOut(r.Context(), req.RefreshToken); err != nil {
		status, code, msg := service.HTTPStatusCodeAndMessage(err, "SIGN_OUT_FAILED", "Failed to sign out.")
		response.ErrorWithCode(w, status, code, msg)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	if !ok || claims == nil {
		response.ErrorWithCode(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required.")
		return
	}

	profile, err := h.identityService.GetMe(r.Context(), claims.UserID)
	if err != nil {
		status, code, msg := service.HTTPStatusCodeAndMessage(err, "INTERNAL_ERROR", "Failed to load user profile.")
		response.ErrorWithCode(w, status, code, msg)
		return
	}

	response.JSON(w, http.StatusOK, profile)
}
