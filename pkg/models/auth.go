package models

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpRequest struct {
	Email         string `json:"email"`
	Password      string `json:"password"`
	Name          string `json:"name"`
	UserName      string `json:"user_name"`
	AnchorDialect string `json:"anchor_dialect"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type SignOutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type AuthResponse struct {
	Token        string      `json:"token"`
	RefreshToken string      `json:"refresh_token"`
	User         UserProfile `json:"user"`
}
