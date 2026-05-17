package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/emmanuella-codes/nox/config"
	sharedtoken "github.com/emmanuella-codes/nox/shared/token"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func TestJWTMiddlewareAcceptsValidAccessToken(t *testing.T) {
	cfg := testJWTConfig()
	userID := uuid.New()
	rawToken, err := sharedtoken.SignWithOptions(userID, "access-jti", sharedtoken.AccessTokenType, cfg.JWTAccessSecret, time.Minute, cfg.JWTIssuer, cfg.JWTAudience)
	if err != nil {
		t.Fatalf("sign access token: %v", err)
	}

	app := fiber.New()
	app.Get("/private", JWT(cfg), func(c *fiber.Ctx) error {
		currentUserID, ok := CurrentUserID(c)
		if !ok {
			t.Fatal("expected current user id to be available")
		}
		if currentUserID != userID {
			t.Fatalf("expected user id %s, got %s", userID, currentUserID)
		}
		return c.SendStatus(fiber.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/private", nil)
	req.Header.Set(fiber.HeaderAuthorization, "Bearer "+rawToken)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("execute request: %v", err)
	}
	if res.StatusCode != fiber.StatusNoContent {
		t.Fatalf("expected status %d, got %d", fiber.StatusNoContent, res.StatusCode)
	}
}

func TestJWTMiddlewareRejectsRefreshToken(t *testing.T) {
	cfg := testJWTConfig()
	rawToken, err := sharedtoken.SignWithOptions(uuid.New(), "refresh-jti", sharedtoken.RefreshTokenType, cfg.JWTAccessSecret, time.Minute, cfg.JWTIssuer, cfg.JWTAudience)
	if err != nil {
		t.Fatalf("sign refresh token: %v", err)
	}

	app := fiber.New()
	app.Get("/private", JWT(cfg), func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/private", nil)
	req.Header.Set(fiber.HeaderAuthorization, "Bearer "+rawToken)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("execute request: %v", err)
	}
	if res.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", fiber.StatusUnauthorized, res.StatusCode)
	}
}

func TestJWTMiddlewareRejectsWrongAudience(t *testing.T) {
	cfg := testJWTConfig()
	rawToken, err := sharedtoken.SignWithOptions(uuid.New(), "access-jti", sharedtoken.AccessTokenType, cfg.JWTAccessSecret, time.Minute, cfg.JWTIssuer, "other-client")
	if err != nil {
		t.Fatalf("sign access token: %v", err)
	}

	app := fiber.New()
	app.Get("/private", JWT(cfg), func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/private", nil)
	req.Header.Set(fiber.HeaderAuthorization, "Bearer "+rawToken)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("execute request: %v", err)
	}
	if res.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", fiber.StatusUnauthorized, res.StatusCode)
	}
}

func testJWTConfig() *config.Config {
	return &config.Config{
		JWTAccessSecret: "access-secret",
		JWTIssuer:       "nox-api",
		JWTAudience:     "nox-client",
	}
}
