package services

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/emmanuella-codes/nox/config"
	"github.com/emmanuella-codes/nox/shared/mail"
	sharedtoken "github.com/emmanuella-codes/nox/shared/token"
	"github.com/google/uuid"
)

func TestHashServiceHashesAndComparesPasswords(t *testing.T) {
	service := NewHashService()

	hash, err := service.HashPassword("password123")
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	if hash == "password123" {
		t.Fatal("expected stored password to be hashed")
	}
	if !service.ComparePassword(hash, "password123") {
		t.Fatal("expected password to match hash")
	}
	if service.ComparePassword(hash, "wrong-password") {
		t.Fatal("expected wrong password not to match hash")
	}
}

func TestOTPServiceGeneratesSixDigitCodesAndComparesHashes(t *testing.T) {
	service := NewOTPService()

	otp, err := service.Generate()
	if err != nil {
		t.Fatalf("generate otp: %v", err)
	}
	if len(otp) != 6 {
		t.Fatalf("expected 6 digit OTP, got %q", otp)
	}
	for _, char := range otp {
		if char < '0' || char > '9' {
			t.Fatalf("expected numeric OTP, got %q", otp)
		}
	}

	hash, err := service.Hash(otp)
	if err != nil {
		t.Fatalf("hash otp: %v", err)
	}
	if !service.Compare(hash, otp) {
		t.Fatal("expected OTP to match hash")
	}
	if service.Compare(hash, "000000") && otp != "000000" {
		t.Fatal("expected wrong OTP not to match hash")
	}
}

func TestTokenServiceIssuesAndVerifiesTokenPair(t *testing.T) {
	cfg := testTokenConfig()
	service := NewTokenService(cfg)
	userID := uuid.New()

	pair, err := service.IssuePair(userID)
	if err != nil {
		t.Fatalf("issue token pair: %v", err)
	}
	if pair.AccessToken == "" || pair.RefreshToken == "" || pair.RefreshTokenID == "" {
		t.Fatal("expected access token, refresh token, and refresh token id")
	}

	accessClaims, err := service.VerifyAccessToken(pair.AccessToken)
	if err != nil {
		t.Fatalf("verify access token: %v", err)
	}
	if accessClaims.UserID != userID || accessClaims.Type != sharedtoken.AccessTokenType {
		t.Fatalf("unexpected access claims: %+v", accessClaims)
	}

	refreshClaims, err := service.VerifyRefreshToken(pair.RefreshToken)
	if err != nil {
		t.Fatalf("verify refresh token: %v", err)
	}
	if refreshClaims.UserID != userID || refreshClaims.ID != pair.RefreshTokenID || refreshClaims.Type != sharedtoken.RefreshTokenType {
		t.Fatalf("unexpected refresh claims: %+v", refreshClaims)
	}

	if _, err := service.VerifyAccessToken(pair.RefreshToken); err == nil {
		t.Fatal("expected refresh token to be rejected by access verifier")
	}
}

func TestTokenServiceRejectsWrongIssuerAudienceAndSecret(t *testing.T) {
	cfg := testTokenConfig()
	service := NewTokenService(cfg)
	userID := uuid.New()
	token, err := service.SignAccessToken(userID)
	if err != nil {
		t.Fatalf("sign access token: %v", err)
	}

	wrongAudience := NewTokenService(&config.Config{
		JWTAccessSecret:  cfg.JWTAccessSecret,
		JWTRefreshSecret: cfg.JWTRefreshSecret,
		JWTIssuer:        cfg.JWTIssuer,
		JWTAudience:      "other-client",
		JWTAccessTTL:     cfg.JWTAccessTTL,
		JWTRefreshTTL:    cfg.JWTRefreshTTL,
	})
	if _, err := wrongAudience.VerifyAccessToken(token); err == nil {
		t.Fatal("expected wrong audience to be rejected")
	}

	wrongSecret := NewTokenService(&config.Config{
		JWTAccessSecret:  "wrong-secret",
		JWTRefreshSecret: cfg.JWTRefreshSecret,
		JWTIssuer:        cfg.JWTIssuer,
		JWTAudience:      cfg.JWTAudience,
		JWTAccessTTL:     cfg.JWTAccessTTL,
		JWTRefreshTTL:    cfg.JWTRefreshTTL,
	})
	if _, err := wrongSecret.VerifyAccessToken(token); err == nil {
		t.Fatal("expected wrong secret to be rejected")
	}
}

func TestEmailServiceSendsVerificationOTPMessage(t *testing.T) {
	provider := &serviceTestMailProvider{}
	service := NewEmailService(provider)

	if err := service.SendVerificationOTP(context.Background(), "ada@example.com", "123456"); err != nil {
		t.Fatalf("send verification otp: %v", err)
	}
	if len(provider.messages) != 1 {
		t.Fatalf("expected one message, got %d", len(provider.messages))
	}

	message := provider.messages[0]
	if message.ToEmail != "ada@example.com" {
		t.Fatalf("expected recipient ada@example.com, got %q", message.ToEmail)
	}
	if message.Subject != "Verify your Nox email" {
		t.Fatalf("unexpected subject %q", message.Subject)
	}
	if !strings.Contains(message.TextContent, "123456") || !strings.Contains(message.HTMLContent, "123456") {
		t.Fatal("expected OTP to be included in email content")
	}
}

func testTokenConfig() *config.Config {
	return &config.Config{
		JWTAccessSecret:  "access-secret",
		JWTRefreshSecret: "refresh-secret",
		JWTIssuer:        "nox-api",
		JWTAudience:      "nox-client",
		JWTAccessTTL:     time.Minute,
		JWTRefreshTTL:    time.Hour,
	}
}

type serviceTestMailProvider struct {
	messages []mail.Message
}

func (p *serviceTestMailProvider) Send(ctx context.Context, message mail.Message) error {
	p.messages = append(p.messages, message)
	return nil
}
