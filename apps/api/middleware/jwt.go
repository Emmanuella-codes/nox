package middleware

import (
	"strings"

	"github.com/emmanuella-codes/nox/config"
	"github.com/emmanuella-codes/nox/shared"
	sharedtoken "github.com/emmanuella-codes/nox/shared/token"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const userIDLocalKey = "user_id"

func JWT(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, ok := bearerToken(c.Get(fiber.HeaderAuthorization))
		if !ok {
			return unauthorized(c)
		}

		claims, err := sharedtoken.VerifyWithOptions(token, cfg.JWTAccessSecret, sharedtoken.AccessTokenType, cfg.JWTIssuer, cfg.JWTAudience)
		if err != nil {
			return unauthorized(c)
		}

		c.Locals(userIDLocalKey, claims.UserID)
		return c.Next()
	}
}

func CurrentUserID(c *fiber.Ctx) (uuid.UUID, bool) {
	userID, ok := c.Locals(userIDLocalKey).(uuid.UUID)
	return userID, ok
}

func bearerToken(header string) (string, bool) {
	scheme, token, ok := strings.Cut(header, " ")
	if !ok || !strings.EqualFold(scheme, "Bearer") || strings.TrimSpace(token) == "" {
		return "", false
	}
	return strings.TrimSpace(token), true
}

func unauthorized(c *fiber.Ctx) error {
	return c.Status(fiber.StatusUnauthorized).JSON(shared.PipeRes[any]{
		Success: false,
		Message: "invalid_token",
	})
}
