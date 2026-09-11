package middleware

import (
	"strings"

	"github.com/emmanuella-codes/nox/config"
	"github.com/emmanuella-codes/nox/models"
	adminrepo "github.com/emmanuella-codes/nox/repositories/admin"
	"github.com/emmanuella-codes/nox/shared"
	sharedtoken "github.com/emmanuella-codes/nox/shared/token"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const userIDLocalKey = "user_id"
const adminIdentityLocalKey = "admin_identity"

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

func OptionalJWT(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, ok := bearerToken(c.Get(fiber.HeaderAuthorization))
		if !ok {
			return c.Next()
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

func AdminJWT(cfg *config.Config, adminRepo adminrepo.AdminRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, ok := bearerToken(c.Get(fiber.HeaderAuthorization))
		if !ok {
			return unauthorized(c)
		}

		claims, err := sharedtoken.VerifyWithOptions(token, cfg.AdminJWTAccessSecret, sharedtoken.AccessTokenType, cfg.AdminJWTIssuer, cfg.AdminJWTAudience)
		if err != nil {
			return unauthorized(c)
		}

		membership, err := adminRepo.FindMembershipByUserID(c.Context(), claims.UserID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(shared.PipeRes[any]{
				Success: false,
				Message: "internal_error",
			})
		}
		if membership == nil || !membership.IsActive || !models.ValidAdminRole(membership.Role) {
			return c.Status(fiber.StatusForbidden).JSON(shared.PipeRes[any]{
				Success: false,
				Message: "admin_access_denied",
			})
		}

		c.Locals(userIDLocalKey, claims.UserID)
		c.Locals(adminIdentityLocalKey, membership)
		return c.Next()
	}
}

func CurrentAdminUserID(c *fiber.Ctx) (uuid.UUID, bool) {
	return CurrentUserID(c)
}

func CurrentAdminMembership(c *fiber.Ctx) (*models.AdminMembership, bool) {
	membership, ok := c.Locals(adminIdentityLocalKey).(*models.AdminMembership)
	return membership, ok
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
