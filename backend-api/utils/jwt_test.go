package utils

import (
	"testing"
	"time"

	"idash-backend-api/config"
)

func TestGenerateAndValidateJWT(t *testing.T) {
	config.AppConfig.JWTSecret = "test-secret"
	config.AppConfig.JWTTTL = 1

	roles := []string{"executive", "analyst"}
	permissions := []string{"view_dashboard"}
	companyIDs := []uint{1, 2, 3}

	token, err := GenerateJWT(42, roles, permissions, nil, companyIDs)
	if err != nil {
		t.Fatalf("unexpected error generating token: %v", err)
	}

	claims, err := ValidateJWT(token)
	if err != nil {
		t.Fatalf("unexpected error validating token: %v", err)
	}

	if claims.UserID != 42 {
		t.Errorf("expected user id 42, got %d", claims.UserID)
	}

	if len(claims.Roles) != len(roles) {
		t.Fatalf("expected %d roles, got %d", len(roles), len(claims.Roles))
	}

	if claims.Roles[0] != "executive" {
		t.Errorf("expected executive role, got %s", claims.Roles[0])
	}

	if len(claims.CompanyIDs) != 3 {
		t.Errorf("expected 3 company ids, got %d", len(claims.CompanyIDs))
	}
}

func TestExpiredJWT(t *testing.T) {
	config.AppConfig.JWTSecret = "test-secret-expired"
	config.AppConfig.JWTTTL = 0

	token, err := GenerateJWT(1, nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	time.Sleep(time.Millisecond * 10)
	_, err = ValidateJWT(token)
	if err == nil {
		t.Fatal("expected validation error for expired token")
	}
}
