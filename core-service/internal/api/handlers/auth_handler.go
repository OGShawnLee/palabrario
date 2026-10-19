package handlers

import (
	"encoding/json"
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
		response.Error(w, http.StatusBadRequest, "Cuerpo de petición inválido")
		return
	}

	authResp, err := h.identityService.SignUp(r.Context(), req)
	if err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "No se pudo registrar la cuenta")
		response.Error(w, status, msg)
		return
	}

	response.JSON(w, http.StatusCreated, authResp)
}

func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var req models.SignInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Cuerpo de petición inválido")
		return
	}

	authResp, err := h.identityService.SignIn(r.Context(), req)
	if err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al iniciar sesión")
		response.Error(w, status, msg)
		return
	}

	response.JSON(w, http.StatusOK, authResp)
}

func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	if !ok || claims == nil {
		response.Error(w, http.StatusUnauthorized, "No autorizado")
		return
	}

	profile, err := h.identityService.GetMe(r.Context(), claims.UserID)
	if err != nil {
		status, msg := service.HTTPStatusAndMessage(err, "Error al obtener el perfil del usuario")
		response.Error(w, status, msg)
		return
	}

	response.JSON(w, http.StatusOK, profile)
}
