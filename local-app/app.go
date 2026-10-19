package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"palabrario/pkg/models"
)

// App struct maneja el estado nativo y seguro de la aplicación de escritorio.
// El token JWT se almacena de forma PRIVADA y VOLÁTIL en la memoria de este proceso Go.
// JavaScript en el WebView NUNCA tiene acceso al token JWT crudo.
type App struct {
	ctx        context.Context
	apiBaseURL string
	httpClient *http.Client

	// Estado privado en memoria del proceso Go
	mu          sync.RWMutex
	sessionToken string
	currentUser  *models.UserProfile
}

// NewApp crea una nueva instancia de la aplicación Wails
func NewApp() *App {
	return &App{
		apiBaseURL: "http://localhost:8000/api/v1",
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// startup se ejecuta al iniciar la aplicación Wails
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet método de demostración inicial
func (a *App) Greet(name string) string {
	return fmt.Sprintf("¡Bienvenido a Palabrario, %s!", name)
}

// SignIn realiza la autenticación contra core-service mediante net/http en Go.
// Guarda el JWT en la memoria privada de Go y retorna únicamente el perfil del usuario a Vue.
func (a *App) SignIn(email, password string) (*models.UserProfile, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	password = strings.TrimSpace(password)

	if email == "" || password == "" {
		return nil, fmt.Errorf("correo y contraseña son requeridos")
	}

	payload, err := json.Marshal(models.SignInRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return nil, fmt.Errorf("error al codificar credenciales: %w", err)
	}

	resp, err := a.httpClient.Post(a.apiBaseURL+"/auth/sign-in", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("no se pudo conectar con el servidor central: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		var errResp struct {
			Error string `json:"error"`
		}
		if err := json.Unmarshal(bodyBytes, &errResp); err == nil && errResp.Error != "" {
			return nil, fmt.Errorf("%s", errResp.Error)
		}
		return nil, fmt.Errorf("error de autenticación (código %d)", resp.StatusCode)
	}

	var authResp models.AuthResponse
	if err := json.Unmarshal(bodyBytes, &authResp); err != nil {
		return nil, fmt.Errorf("error decodificando respuesta de sesión: %w", err)
	}

	// Guardar sesión de forma privada y segura en Go
	a.mu.Lock()
	a.sessionToken = authResp.Token
	a.currentUser = &authResp.User
	a.mu.Unlock()

	return &authResp.User, nil
}

// SignUp registra un nuevo aprendiz mediante net/http en Go.
// Almacena el JWT en memoria privada y devuelve el perfil del usuario creado a Vue.
func (a *App) SignUp(name, userName, email, password, anchorDialect string) (*models.UserProfile, error) {
	name = strings.TrimSpace(name)
	userName = strings.TrimSpace(strings.ToLower(userName))
	email = strings.TrimSpace(strings.ToLower(email))
	password = strings.TrimSpace(password)
	anchorDialect = strings.TrimSpace(strings.ToUpper(anchorDialect))
	if anchorDialect == "" {
		anchorDialect = "MX"
	}

	if name == "" || userName == "" || email == "" || password == "" {
		return nil, fmt.Errorf("todos los campos son obligatorios")
	}

	payload, err := json.Marshal(models.SignUpRequest{
		Name:          name,
		UserName:      userName,
		Email:         email,
		Password:      password,
		AnchorDialect: anchorDialect,
	})
	if err != nil {
		return nil, fmt.Errorf("error al codificar datos de registro: %w", err)
	}

	resp, err := a.httpClient.Post(a.apiBaseURL+"/auth/sign-up", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("no se pudo conectar con el servidor central: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		var errResp struct {
			Error string `json:"error"`
		}
		if err := json.Unmarshal(bodyBytes, &errResp); err == nil && errResp.Error != "" {
			return nil, fmt.Errorf("%s", errResp.Error)
		}
		return nil, fmt.Errorf("error al registrar cuenta (código %d)", resp.StatusCode)
	}

	var authResp models.AuthResponse
	if err := json.Unmarshal(bodyBytes, &authResp); err != nil {
		return nil, fmt.Errorf("error decodificando respuesta de registro: %w", err)
	}

	// Guardar sesión de forma privada y segura en Go
	a.mu.Lock()
	a.sessionToken = authResp.Token
	a.currentUser = &authResp.User
	a.mu.Unlock()

	return &authResp.User, nil
}

// GetCurrentUser retorna el perfil del usuario autenticado en la memoria de Go.
func (a *App) GetCurrentUser() *models.UserProfile {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.currentUser
}

// SignOut borra la sesión de la memoria privada de Go.
func (a *App) SignOut() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.sessionToken = ""
	a.currentUser = nil
}

// IsAuthenticated indica si existe una sesión activa en memoria.
func (a *App) IsAuthenticated() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.sessionToken != ""
}
