package auth

import "testing"

func TestJWTServiceGeneratesAndValidatesToken(t *testing.T) {
	service := NewJWTService("test-secret")

	token, err := service.Generate("admin", []string{"users:read"})
	if err != nil {
		t.Fatalf("expected token generation to succeed: %v", err)
	}

	claims, err := service.Validate(token)
	if err != nil {
		t.Fatalf("expected token validation to succeed: %v", err)
	}

	if claims.Subject != "admin" {
		t.Fatalf("expected subject %q, got %q", "admin", claims.Subject)
	}

	if !claims.HasPermission("users:read") {
		t.Fatal("expected users:read permission to be present")
	}
}

func TestJWTServiceRejectsInvalidToken(t *testing.T) {
	service := NewJWTService("test-secret")

	if _, err := service.Validate("invalid-token"); err == nil {
		t.Fatal("expected invalid token to be rejected")
	}
}
