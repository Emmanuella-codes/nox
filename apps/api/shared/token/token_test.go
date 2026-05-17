package token

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestVerifyWithOptionsRequiresJWTID(t *testing.T) {
	_, err := SignWithOptions(uuid.New(), "", AccessTokenType, "secret", time.Minute, "nox-api", "nox-client")
	if err == nil {
		t.Fatal("expected signing without token id to fail")
	}
}

func TestVerifyWithOptionsRejectsWrongTokenType(t *testing.T) {
	rawToken, err := SignWithOptions(uuid.New(), "refresh-jti", RefreshTokenType, "secret", time.Minute, "nox-api", "nox-client")
	if err != nil {
		t.Fatalf("sign refresh token: %v", err)
	}

	if _, err := VerifyWithOptions(rawToken, "secret", AccessTokenType, "nox-api", "nox-client"); err == nil {
		t.Fatal("expected refresh token to be rejected when access token is required")
	}
}

func TestVerifyWithOptionsRejectsWrongIssuerOrAudience(t *testing.T) {
	rawToken, err := SignWithOptions(uuid.New(), "access-jti", AccessTokenType, "secret", time.Minute, "nox-api", "nox-client")
	if err != nil {
		t.Fatalf("sign access token: %v", err)
	}

	if _, err := VerifyWithOptions(rawToken, "secret", AccessTokenType, "other-issuer", "nox-client"); err == nil {
		t.Fatal("expected token with wrong issuer to be rejected")
	}

	if _, err := VerifyWithOptions(rawToken, "secret", AccessTokenType, "nox-api", "other-client"); err == nil {
		t.Fatal("expected token with wrong audience to be rejected")
	}

	claims, err := VerifyWithOptions(rawToken, "secret", AccessTokenType, "nox-api", "nox-client")
	if err != nil {
		t.Fatalf("expected token with matching issuer and audience to verify: %v", err)
	}
	if claims.ID != "access-jti" {
		t.Fatalf("expected jti to round trip, got %q", claims.ID)
	}
}
