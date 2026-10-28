package main

import (
	"net/http"
	"net/url"

	"palabrario/pkg/models"
)

// =============================================================================
// MÓDULO 1: Sesión & Autenticación (CU01)
// =============================================================================

// SignIn autentica credenciales y almacena el JWT en RAM segura.
func (a *App) SignIn(req models.SignInRequest) (*models.UserProfile, error) {
	var authResp models.AuthResponse
	if err := a.doRequest(a.context(), http.MethodPost, a.coreURL+"/api/v1/auth/sign-in", req, &authResp); err != nil {
		return nil, err
	}
	a.setSessionToken(authResp.Token)
	return &authResp.User, nil
}

// SignUp registra un nuevo aprendiz y custodia su JWT en RAM segura.
func (a *App) SignUp(req models.SignUpRequest) (*models.UserProfile, error) {
	var authResp models.AuthResponse
	if err := a.doRequest(a.context(), http.MethodPost, a.coreURL+"/api/v1/auth/sign-up", req, &authResp); err != nil {
		return nil, err
	}
	a.setSessionToken(authResp.Token)
	return &authResp.User, nil
}

// SignOut limpia el JWT en memoria y cierra conexiones activas.
func (a *App) SignOut() error {
	a.setSessionToken("")
	_ = a.LeaveArenaRoom()
	return nil
}

// GetMe obtiene el perfil del usuario autenticado.
func (a *App) GetMe() (*models.UserProfile, error) {
	var user models.UserProfile
	if err := a.doRequest(a.context(), http.MethodGet, a.coreURL+"/api/v1/auth/me", nil, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// GetMyProfile alias canónico para obtener el perfil propio.
func (a *App) GetMyProfile() (*models.UserProfile, error) {
	return a.GetMe()
}

// =============================================================================
// MÓDULO 2: Identidad & Perfil Social (CU01)
// =============================================================================

// GetPublicProfileByUsername obtiene el perfil público con enriquecimiento de propietario si aplica.
func (a *App) GetPublicProfileByUsername(username string) (*models.PublicUserProfile, error) {
	var profile models.PublicUserProfile
	targetURL := a.coreURL + "/api/v1/users/u/" + url.PathEscape(username)
	if err := a.doRequest(a.context(), http.MethodGet, targetURL, nil, &profile); err != nil {
		return nil, err
	}
	return &profile, nil
}

// UpdateProfile actualiza los metadatos editables del perfil propio.
func (a *App) UpdateProfile(req models.UpdateProfileRequest) (*models.UserProfile, error) {
	var profile models.UserProfile
	if err := a.doRequest(a.context(), http.MethodPut, a.coreURL+"/api/v1/users/me", req, &profile); err != nil {
		return nil, err
	}
	return &profile, nil
}

// ChangePassword actualiza la contraseña del usuario autenticado.
func (a *App) ChangePassword(req models.ChangePasswordRequest) error {
	return a.doRequest(a.context(), http.MethodPut, a.coreURL+"/api/v1/users/me/password", req, nil)
}
