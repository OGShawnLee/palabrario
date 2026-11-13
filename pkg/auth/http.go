package auth

import (
	"errors"
	"net/http"
	"strings"
)

// ExtractBearerToken extracts a Bearer token from the Authorization header.
// It treats the "Bearer" scheme case-insensitively.
// If no Authorization header is present, it returns ("", false, nil).
// If an Authorization header is present but malformed or missing its token, it returns ("", true, err).
func ExtractBearerToken(r *http.Request) (token string, present bool, err error) {
	if r == nil {
		return "", false, errors.New("request is nil")
	}

	raw := strings.TrimSpace(r.Header.Get("Authorization"))
	if raw == "" {
		return "", false, nil
	}

	parts := strings.SplitN(raw, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", true, errors.New("invalid authorization header format; expected Bearer <token>")
	}

	token = strings.TrimSpace(parts[1])
	if token == "" || strings.ContainsAny(token, " \t\r\n") {
		return "", true, errors.New("authorization bearer token is empty or malformed")
	}

	return token, true, nil
}

// DefaultDevCORSOrigins returns the default allowed origins for local development, desktop Wails, and mobile Capacitor WebViews.
func DefaultDevCORSOrigins() []string {
	return []string{
		"http://localhost:*",
		"http://127.0.0.1:*",
		"wails://wails.localhost",
		"http://wails.localhost",
		"capacitor://localhost",
	}
}

// ParseCORSAllowedOrigins parses a comma-separated CORS_ALLOWED_ORIGINS value.
// When raw is empty, it returns the default local development/desktop/mobile origins.
func ParseCORSAllowedOrigins(raw string) []string {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return DefaultDevCORSOrigins()
	}

	parts := strings.Split(trimmed, ",")
	origins := make([]string, 0, len(parts))
	for _, part := range parts {
		origin := strings.TrimSpace(part)
		if origin != "" {
			origins = append(origins, origin)
		}
	}
	if len(origins) == 0 {
		return DefaultDevCORSOrigins()
	}
	return origins
}

// IsOriginAllowed checks whether an HTTP Origin header matches the configured allowlist.
// An empty Origin header is allowed for non-browser native clients (such as the Wails Go BFF).
func IsOriginAllowed(origin string, allowedOrigins []string) bool {
	trimmed := strings.TrimSpace(origin)
	if trimmed == "" {
		return true
	}

	lowerOrigin := strings.ToLower(trimmed)
	for _, pattern := range allowedOrigins {
		p := strings.ToLower(strings.TrimSpace(pattern))
		if p == "" {
			continue
		}
		if p == "*" || lowerOrigin == p {
			return true
		}
		if strings.Contains(p, "*") {
			prefix, suffix, _ := strings.Cut(p, "*")
			if len(lowerOrigin) < len(prefix)+len(suffix) {
				continue
			}
			if !strings.HasPrefix(lowerOrigin, prefix) || !strings.HasSuffix(lowerOrigin, suffix) {
				continue
			}
			// Guard trailing port wildcards like "http://localhost:*" so only digits match the port segment.
			if strings.HasSuffix(prefix, ":") && suffix == "" {
				portPart := strings.TrimPrefix(lowerOrigin, prefix)
				if !isAllDigits(portPart) {
					continue
				}
			}
			return true
		}
	}
	return false
}

func isAllDigits(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}
