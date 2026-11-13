package main

import (
	"errors"
	"net/http"
	"net/url"

	"palabrario/pkg/models"
)

// =============================================================================
// MÓDULO 1: Sesión & Autenticación (CU01)
// =============================================================================

// SignIn authenticates user credentials, stores the access token in RAM, and persists the refresh token on disk.
func (a *App) SignIn(req models.SignInRequest) (*models.UserProfile, error) {
	var authResp models.AuthResponse
	if err := a.doRequest(a.context(), http.MethodPost, a.coreURL+"/api/v1/auth/sign-in", req, &authResp); err != nil {
		return nil, err
	}
	a.setTokens(authResp.Token, authResp.RefreshToken)
	return &authResp.User, nil
}

// SignUp registers a new learner, stores the access token in RAM, and persists the refresh token on disk.
func (a *App) SignUp(req models.SignUpRequest) (*models.UserProfile, error) {
	var authResp models.AuthResponse
	if err := a.doRequest(a.context(), http.MethodPost, a.coreURL+"/api/v1/auth/sign-up", req, &authResp); err != nil {
		return nil, err
	}
	a.setTokens(authResp.Token, authResp.RefreshToken)
	return &authResp.User, nil
}

// SignOut revokes the refresh token on the server (best-effort), clears local tokens, and closes active WS connections.
func (a *App) SignOut() error {
	refresh := a.getRefreshToken()
	a.clearSession()
	_ = a.LeaveArenaRoom()
	if refresh != "" {
		_ = a.doRequest(a.context(), http.MethodPost, a.coreURL+"/api/v1/auth/sign-out", models.SignOutRequest{RefreshToken: refresh}, nil)
	}
	return nil
}

// GetMe returns the authenticated user's profile, exchanging any persisted refresh token on startup if needed.
// Returns (nil, nil) when unauthenticated or when the session has expired/been suspended.
func (a *App) GetMe() (*models.UserProfile, error) {
	access, refresh := a.getTokens()
	if access == "" && refresh == "" {
		return nil, nil
	}

	if access == "" && refresh != "" {
		profile, refreshed, err := a.refreshSessionIfStale(a.context(), "")
		if err != nil {
			var apiErr *APIError
			if errors.As(err, &apiErr) && (apiErr.Status == http.StatusUnauthorized || apiErr.Status == http.StatusForbidden) {
				return nil, nil
			}
			return nil, err
		}
		if refreshed && profile != nil {
			return profile, nil
		}
	}

	var user models.UserProfile
	if err := a.doRequest(a.context(), http.MethodGet, a.coreURL+"/api/v1/auth/me", nil, &user); err != nil {
		var apiErr *APIError
		if errors.As(err, &apiErr) && (apiErr.Status == http.StatusUnauthorized || apiErr.Status == http.StatusForbidden) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// =============================================================================
// MÓDULO 2: Identidad & Perfil Social (CU01)
// =============================================================================

// GetPublicProfileByUsername fetches a public profile with owner enrichment when authenticated.
func (a *App) GetPublicProfileByUsername(username string) (*models.PublicUserProfile, error) {
	var profile models.PublicUserProfile
	targetURL := a.coreURL + "/api/v1/users/u/" + url.PathEscape(username)
	if err := a.doRequest(a.context(), http.MethodGet, targetURL, nil, &profile); err != nil {
		return nil, err
	}
	return &profile, nil
}

// UpdateProfile updates the authenticated user's editable profile fields and refreshes the access token
// so claims such as AnchorDialect stay synchronized.
func (a *App) UpdateProfile(req models.UpdateProfileRequest) (*models.UserProfile, error) {
	var profile models.UserProfile
	if err := a.doRequest(a.context(), http.MethodPut, a.coreURL+"/api/v1/users/me", req, &profile); err != nil {
		return nil, err
	}
	if a.getRefreshToken() != "" {
		_, _ = a.refreshSession(a.context())
	}
	return &profile, nil
}

// ChangePassword updates the authenticated user's password.
func (a *App) ChangePassword(req models.ChangePasswordRequest) error {
	return a.doRequest(a.context(), http.MethodPut, a.coreURL+"/api/v1/users/me/password", req, nil)
}
