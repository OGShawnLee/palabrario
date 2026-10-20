package auth_test

import (
	"testing"
	"time"

	"palabrario/pkg/auth"
	"palabrario/pkg/models"
)

func TestPasswordHashingAndJWT(t *testing.T) {
	hash, err := auth.HashPassword("password123")
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}
	t.Logf("Valid bcrypt hash for 'password123': %s", hash)

	if !auth.CheckPassword("password123", hash) {
		t.Fatal("CheckPassword returned false for valid password")
	}

	svc := auth.NewTokenService("test-secret-key-minimum-32-chars", 24*time.Hour)
	token, err := svc.GenerateToken(&models.UserProfile{
		IDUser:        "b1111111-0000-0000-0000-000000000002",
		Email:         "silas@palabrario.org",
		UserName:      "silas_editor",
		AnchorDialect: "MX",
		Role:          models.RoleEditor,
	})
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
}
