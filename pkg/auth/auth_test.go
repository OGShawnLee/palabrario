package auth_test

import (
	"crypto/sha256"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"palabrario/pkg/auth"
	"palabrario/pkg/models"
)

func sampleUser() *models.UserProfile {
	return &models.UserProfile{
		IDUser:        "b1111111-0000-0000-0000-000000000002",
		Email:         "silas@palabrario.org",
		UserName:      "silas_editor",
		AnchorDialect: "MX",
		Role:          models.RoleEditor,
	}
}

func TestPasswordHashingAndJWT(t *testing.T) {
	hash, err := auth.HashPassword("password123")
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if !auth.CheckPassword("password123", hash) {
		t.Fatal("CheckPassword returned false for valid password")
	}
	if auth.CheckPassword("wrong-password", hash) {
		t.Fatal("CheckPassword returned true for wrong password")
	}

	svc := auth.NewTokenService("test-secret-key-minimum-32-chars", 24*time.Hour)
	token, err := svc.GenerateToken(sampleUser())
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	claims, err := svc.VerifyToken(token)
	if err != nil {
		t.Fatalf("VerifyToken failed: %v", err)
	}
	if claims.AnchorDialect != "MX" || claims.Role != models.RoleEditor {
		t.Fatalf("unexpected claims: %+v", claims)
	}
	if claims.TokenVersion != 1 {
		t.Fatalf("expected default TokenVersion 1, got %d", claims.TokenVersion)
	}
	if claims.TokenType != auth.TokenTypeAccess {
		t.Fatalf("expected TokenType %q, got %q", auth.TokenTypeAccess, claims.TokenType)
	}
}

func TestTokenServiceDefaultTTLAndAccessToken(t *testing.T) {
	svc := auth.NewTokenService("test-secret-key-minimum-32-chars", 0)
	if svc.AccessTTL() != 15*time.Minute {
		t.Fatalf("expected default TTL 15m, got %v", svc.AccessTTL())
	}

	tok, err := svc.GenerateAccessToken(sampleUser(), 4)
	if err != nil {
		t.Fatalf("GenerateAccessToken failed: %v", err)
	}

	claims, err := svc.VerifyToken(tok)
	if err != nil {
		t.Fatalf("VerifyToken failed: %v", err)
	}
	if claims.TokenVersion != 4 {
		t.Fatalf("expected TokenVersion 4, got %d", claims.TokenVersion)
	}
	if claims.TokenType != "access" {
		t.Fatalf("expected TokenType access, got %q", claims.TokenType)
	}
	ttl := claims.ExpiresAt.Time.Sub(claims.IssuedAt.Time)
	if ttl != 15*time.Minute {
		t.Fatalf("expected 15m token lifetime, got %v", ttl)
	}
}

func TestVerifyToken_RejectsInvalidTokens(t *testing.T) {
	secret := "test-secret-key-minimum-32-chars"
	svc := auth.NewTokenService(secret, 15*time.Minute)

	t.Run("expired JWT rejected", func(t *testing.T) {
		past := time.Now().UTC().Add(-2 * time.Minute)
		claims := auth.Claims{
			UserID:        "u-1",
			Email:         "u1@palabrario.org",
			UserName:      "u1",
			AnchorDialect: "MX",
			Role:          models.RoleLearner,
			TokenVersion:  1,
			TokenType:     "access",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(past.Add(time.Minute)),
				IssuedAt:  jwt.NewNumericDate(past),
				NotBefore: jwt.NewNumericDate(past),
				Issuer:    "palabrario-core-service",
				Subject:   "u-1",
			},
		}
		tok, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
		if err != nil {
			t.Fatalf("sign failed: %v", err)
		}
		if _, err := svc.VerifyToken(tok); err == nil {
			t.Fatal("expected expired token to fail verification")
		}
	})

	t.Run("future nbf rejected", func(t *testing.T) {
		future := time.Now().UTC().Add(5 * time.Minute)
		claims := auth.Claims{
			UserID:        "u-1",
			Email:         "u1@palabrario.org",
			UserName:      "u1",
			AnchorDialect: "MX",
			Role:          models.RoleLearner,
			TokenVersion:  1,
			TokenType:     "access",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(future.Add(15 * time.Minute)),
				IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
				NotBefore: jwt.NewNumericDate(future),
				Issuer:    "palabrario-core-service",
				Subject:   "u-1",
			},
		}
		tok, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
		if err != nil {
			t.Fatalf("sign failed: %v", err)
		}
		if _, err := svc.VerifyToken(tok); err == nil {
			t.Fatal("expected future nbf token to fail verification")
		}
	})

	t.Run("alg=none rejected", func(t *testing.T) {
		now := time.Now().UTC()
		claims := auth.Claims{
			UserID:   "u-1",
			UserName: "u1",
			Role:     models.RoleAdmin,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
				IssuedAt:  jwt.NewNumericDate(now),
				NotBefore: jwt.NewNumericDate(now),
				Issuer:    "palabrario-core-service",
			},
		}
		tok, err := jwt.NewWithClaims(jwt.SigningMethodNone, claims).SignedString(jwt.UnsafeAllowNoneSignatureType)
		if err != nil {
			t.Fatalf("sign none failed: %v", err)
		}
		if _, err := svc.VerifyToken(tok); err == nil {
			t.Fatal("expected alg=none token to fail verification")
		}
	})

	t.Run("tampered signature rejected", func(t *testing.T) {
		tok, err := svc.GenerateToken(sampleUser())
		if err != nil {
			t.Fatalf("GenerateToken failed: %v", err)
		}
		parts := strings.Split(tok, ".")
		if len(parts) != 3 {
			t.Fatalf("unexpected JWT parts: %d", len(parts))
		}
		// Flip last character of signature
		sig := parts[2]
		last := sig[len(sig)-1]
		replacement := "A"
		if last == 'A' {
			replacement = "B"
		}
		tampered := parts[0] + "." + parts[1] + "." + sig[:len(sig)-1] + replacement
		if _, err := svc.VerifyToken(tampered); err == nil {
			t.Fatal("expected tampered token to fail verification")
		}
	})

	t.Run("wrong secret rejected", func(t *testing.T) {
		otherSvc := auth.NewTokenService("another-secret-key-minimum-32-chars", 15*time.Minute)
		tok, err := otherSvc.GenerateToken(sampleUser())
		if err != nil {
			t.Fatalf("GenerateToken failed: %v", err)
		}
		if _, err := svc.VerifyToken(tok); err == nil {
			t.Fatal("expected token signed with wrong secret to fail verification")
		}
	})

	t.Run("wrong issuer rejected", func(t *testing.T) {
		now := time.Now().UTC()
		claims := auth.Claims{
			UserID:       "u-1",
			UserName:     "u1",
			Role:         models.RoleLearner,
			TokenVersion: 1,
			TokenType:    "access",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)),
				IssuedAt:  jwt.NewNumericDate(now),
				NotBefore: jwt.NewNumericDate(now),
				Issuer:    "evil-issuer",
			},
		}
		tok, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
		if err != nil {
			t.Fatalf("sign failed: %v", err)
		}
		if _, err := svc.VerifyToken(tok); err == nil {
			t.Fatal("expected wrong issuer token to fail verification")
		}
	})

	t.Run("non-access token type rejected", func(t *testing.T) {
		now := time.Now().UTC()
		claims := auth.Claims{
			UserID:       "u-1",
			UserName:     "u1",
			Role:         models.RoleLearner,
			TokenVersion: 1,
			TokenType:    "refresh",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)),
				IssuedAt:  jwt.NewNumericDate(now),
				NotBefore: jwt.NewNumericDate(now),
				Issuer:    "palabrario-core-service",
			},
		}
		tok, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
		if err != nil {
			t.Fatalf("sign failed: %v", err)
		}
		if _, err := svc.VerifyToken(tok); err == nil {
			t.Fatal("expected non-access token_type to fail verification")
		}
	})
}

func TestBcryptMaxLength72Bytes(t *testing.T) {
	exact72 := strings.Repeat("a", 72)
	hash, err := auth.HashPassword(exact72)
	if err != nil {
		t.Fatalf("expected 72-byte password to succeed, got: %v", err)
	}
	if !auth.CheckPassword(exact72, hash) {
		t.Fatal("expected CheckPassword to succeed for 72-byte password")
	}

	over72ASCII := strings.Repeat("a", 73)
	if _, err := auth.HashPassword(over72ASCII); err == nil {
		t.Fatal("expected >72-byte ASCII password to be rejected")
	}
	if auth.CheckPassword(over72ASCII, hash) {
		t.Fatal("expected CheckPassword to return false for >72-byte password")
	}

	// Multi-byte UTF-8: 25 'á' chars = 50 bytes (allowed); 37 'á' chars = 74 bytes (rejected)
	utf8Under := strings.Repeat("á", 36) // 72 bytes
	if _, err := auth.HashPassword(utf8Under); err != nil {
		t.Fatalf("expected 72-byte UTF-8 password to succeed, got: %v", err)
	}
	utf8Over := strings.Repeat("á", 37) // 74 bytes
	if _, err := auth.HashPassword(utf8Over); err == nil {
		t.Fatal("expected 74-byte UTF-8 password to be rejected")
	}
}

func TestRefreshTokenGenerationAndHashing(t *testing.T) {
	raw1, hash1, err := auth.GenerateRefreshToken()
	if err != nil {
		t.Fatalf("GenerateRefreshToken failed: %v", err)
	}
	raw2, hash2, err := auth.GenerateRefreshToken()
	if err != nil {
		t.Fatalf("second GenerateRefreshToken failed: %v", err)
	}

	if raw1 == "" || raw2 == "" || raw1 == raw2 {
		t.Fatalf("expected distinct non-empty raw refresh tokens, got %q and %q", raw1, raw2)
	}
	if len(hash1) != 64 || len(hash2) != 64 || hash1 == hash2 {
		t.Fatalf("expected distinct 64-char SHA-256 hex hashes, got %q and %q", hash1, hash2)
	}
	if got := auth.HashRefreshToken(raw1); got != hash1 {
		t.Fatalf("HashRefreshToken not deterministic: expected %s, got %s", hash1, got)
	}
	if gotAgain := auth.HashRefreshToken(raw1); gotAgain != hash1 {
		t.Fatalf("HashRefreshToken second call mismatch: expected %s, got %s", hash1, gotAgain)
	}
}

func TestInternalS2STokenLifecycleAndCrossUseRejection(t *testing.T) {
	secret := "test-secret-key-minimum-32-chars"
	svc := auth.NewTokenService(secret, 15*time.Minute)

	s2sTok, err := svc.GenerateInternalToken("realtime-service", "MX1234")
	if err != nil {
		t.Fatalf("GenerateInternalToken failed: %v", err)
	}

	s2sClaims, err := svc.VerifyInternalToken(s2sTok)
	if err != nil {
		t.Fatalf("VerifyInternalToken failed: %v", err)
	}
	if s2sClaims.ServiceName != "realtime-service" || s2sClaims.MatchScope != "MX1234" {
		t.Fatalf("unexpected S2S claims: %+v", s2sClaims)
	}

	// Cross-use: S2S token must NEVER pass VerifyToken
	if _, err := svc.VerifyToken(s2sTok); err == nil {
		t.Fatal("expected S2S token to be rejected by VerifyToken")
	}

	// Cross-use: User token must NEVER pass VerifyInternalToken
	userTok, err := svc.GenerateToken(sampleUser())
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}
	if _, err := svc.VerifyInternalToken(userTok); err == nil {
		t.Fatal("expected user token to be rejected by VerifyInternalToken")
	}

	// Expired S2S token must be rejected
	s2sKey := sha256.Sum256([]byte(secret + ":palabrario-s2s-internal"))
	past := time.Now().UTC().Add(-2 * time.Minute)
	expiredS2S := auth.InternalServiceClaims{
		ServiceName: "realtime-service",
		MatchScope:  "MX1234",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(past.Add(30 * time.Second)),
			IssuedAt:  jwt.NewNumericDate(past),
			NotBefore: jwt.NewNumericDate(past),
			Issuer:    "palabrario-internal-s2s",
			Subject:   "realtime-service",
		},
	}
	expiredTok, err := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredS2S).SignedString(s2sKey[:])
	if err != nil {
		t.Fatalf("sign expired S2S token failed: %v", err)
	}
	if _, err := svc.VerifyInternalToken(expiredTok); err == nil {
		t.Fatal("expected expired S2S token to fail verification")
	}
}

func TestExtractBearerToken(t *testing.T) {
	cases := []struct {
		name        string
		header      string
		wantToken   string
		wantPresent bool
		wantErr     bool
	}{
		{name: "missing header", header: "", wantToken: "", wantPresent: false, wantErr: false},
		{name: "standard Bearer", header: "Bearer abc.def.ghi", wantToken: "abc.def.ghi", wantPresent: true, wantErr: false},
		{name: "lowercase bearer", header: "bearer abc.def.ghi", wantToken: "abc.def.ghi", wantPresent: true, wantErr: false},
		{name: "mixed case bEaReR", header: "bEaReR   abc.def.ghi  ", wantToken: "abc.def.ghi", wantPresent: true, wantErr: false},
		{name: "empty token after Bearer", header: "Bearer   ", wantToken: "", wantPresent: true, wantErr: true},
		{name: "wrong scheme Basic", header: "Basic dXNlcjpwYXNz", wantToken: "", wantPresent: true, wantErr: true},
		{name: "token with internal spaces", header: "Bearer part1 part2", wantToken: "", wantPresent: true, wantErr: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tc.header != "" {
				req.Header.Set("Authorization", tc.header)
			}
			gotToken, gotPresent, err := auth.ExtractBearerToken(req)
			if gotPresent != tc.wantPresent {
				t.Fatalf("present: got %v, want %v", gotPresent, tc.wantPresent)
			}
			if (err != nil) != tc.wantErr {
				t.Fatalf("err: got %v, wantErr %v", err, tc.wantErr)
			}
			if gotToken != tc.wantToken {
				t.Fatalf("token: got %q, want %q", gotToken, tc.wantToken)
			}
		})
	}
}

func TestCORSAllowedOriginsAndMatching(t *testing.T) {
	defaults := auth.ParseCORSAllowedOrigins("")
	if len(defaults) != 5 {
		t.Fatalf("expected 5 default dev origins, got %v", defaults)
	}

	custom := auth.ParseCORSAllowedOrigins(" https://app.palabrario.org , , https://admin.palabrario.org ")
	if len(custom) != 2 || custom[0] != "https://app.palabrario.org" || custom[1] != "https://admin.palabrario.org" {
		t.Fatalf("unexpected parsed custom origins: %v", custom)
	}

	if !auth.IsOriginAllowed("", custom) {
		t.Fatal("expected empty Origin to be allowed for native desktop clients")
	}
	if !auth.IsOriginAllowed("http://localhost:5173", defaults) {
		t.Fatal("expected http://localhost:5173 to match default dev origins")
	}
	if !auth.IsOriginAllowed("http://127.0.0.1:8080", defaults) {
		t.Fatal("expected http://127.0.0.1:8080 to match default dev origins")
	}
	if !auth.IsOriginAllowed("wails://wails.localhost", defaults) {
		t.Fatal("expected wails://wails.localhost to match default dev origins")
	}
	if !auth.IsOriginAllowed("capacitor://localhost", defaults) {
		t.Fatal("expected capacitor://localhost to match default dev origins")
	}
	if auth.IsOriginAllowed("http://localhost:5173.evil.com", defaults) {
		t.Fatal("expected spoofed port-suffix origin to be rejected")
	}
	if auth.IsOriginAllowed("https://evil.example.com", defaults) {
		t.Fatal("expected unknown external origin to be rejected")
	}
}
