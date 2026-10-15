package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"palabrario/core-service/internal/auth"
	"palabrario/core-service/internal/repository"
	"palabrario/pkg/models"
)

type AuthHandler struct {
	userRepo     *repository.UserRepository
	tokenService *auth.TokenService
}

func NewAuthHandler(userRepo *repository.UserRepository, tokenService *auth.TokenService) *AuthHandler {
	return &AuthHandler{
		userRepo:     userRepo,
		tokenService: tokenService,
	}
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var req models.SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Cuerpo de petición inválido"}`, http.StatusBadRequest)
		return
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.UserName = strings.TrimSpace(strings.ToLower(req.UserName))
	req.Name = strings.TrimSpace(req.Name)

	if req.Email == "" || req.Password == "" || req.Name == "" || req.UserName == "" {
		http.Error(w, `{"error":"Todos los campos son obligatorios"}`, http.StatusBadRequest)
		return
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		http.Error(w, `{"error":"Error al procesar la contraseña"}`, http.StatusInternalServerError)
		return
	}

	profile, err := h.userRepo.CreateUser(r.Context(), req.Email, hash, req.Name, req.UserName, req.AnchorDialect)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique") {
			http.Error(w, `{"error":"El correo o nombre de usuario ya está registrado"}`, http.StatusConflict)
			return
		}
		http.Error(w, `{"error":"No se pudo registrar la cuenta"}`, http.StatusInternalServerError)
		return
	}

	token, err := h.tokenService.GenerateToken(profile)
	if err != nil {
		http.Error(w, `{"error":"Error al emitir el token de sesión"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(models.AuthResponse{
		Token: token,
		User:  *profile,
	})
}

func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var req models.SignInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Cuerpo de petición inválido"}`, http.StatusBadRequest)
		return
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if req.Email == "" || req.Password == "" {
		http.Error(w, `{"error":"Correo y contraseña requeridos"}`, http.StatusBadRequest)
		return
	}

	account, err := h.userRepo.GetAccountByEmail(r.Context(), req.Email)
	if err != nil {
		http.Error(w, `{"error":"Credenciales incorrectas"}`, http.StatusUnauthorized)
		return
	}

	if !auth.CheckPassword(req.Password, account.PasswordHash) {
		http.Error(w, `{"error":"Credenciales incorrectas"}`, http.StatusUnauthorized)
		return
	}

	profile, err := h.userRepo.GetProfileByID(r.Context(), repository.UUIDToString(account.IDAccount))
	if err != nil {
		http.Error(w, `{"error":"Error recuperando perfil de usuario"}`, http.StatusInternalServerError)
		return
	}

	token, err := h.tokenService.GenerateToken(profile)
	if err != nil {
		http.Error(w, `{"error":"Error emitiendo credenciales"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(models.AuthResponse{
		Token: token,
		User:  *profile,
	})
}

func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	if !ok || claims == nil {
		http.Error(w, `{"error":"No autorizado"}`, http.StatusUnauthorized)
		return
	}

	profile, err := h.userRepo.GetProfileByID(r.Context(), claims.UserID)
	if err != nil {
		http.Error(w, `{"error":"Usuario no encontrado"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(profile)
}
