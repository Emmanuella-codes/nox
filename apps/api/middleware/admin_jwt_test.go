package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/emmanuella-codes/nox/config"
	"github.com/emmanuella-codes/nox/models"
	adminrepo "github.com/emmanuella-codes/nox/repositories/admin"
	sharedtoken "github.com/emmanuella-codes/nox/shared/token"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func TestAdminJWTAcceptsActiveAdminToken(t *testing.T) {
	cfg := &config.Config{
		AdminJWTAccessSecret: "admin-access-secret",
		AdminJWTIssuer:       "nox-api",
		AdminJWTAudience:     "nox-admin",
	}
	userID := uuid.New()
	rawToken, err := sharedtoken.SignWithOptions(userID, "admin-access-jti", sharedtoken.AccessTokenType, cfg.AdminJWTAccessSecret, time.Minute, cfg.AdminJWTIssuer, cfg.AdminJWTAudience)
	if err != nil {
		t.Fatalf("sign admin access token: %v", err)
	}

	app := fiber.New()
	app.Get("/admin", AdminJWT(cfg, adminMiddlewareRepo{
		membership: &models.AdminMembership{UserID: userID, Role: models.AdminRoleSuperAdmin, IsActive: true},
	}), func(c *fiber.Ctx) error {
		currentUserID, ok := CurrentAdminUserID(c)
		if !ok || currentUserID != userID {
			t.Fatal("expected current admin user id")
		}
		return c.SendStatus(fiber.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	req.Header.Set(fiber.HeaderAuthorization, "Bearer "+rawToken)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("execute request: %v", err)
	}
	if res.StatusCode != fiber.StatusNoContent {
		t.Fatalf("expected status %d, got %d", fiber.StatusNoContent, res.StatusCode)
	}
}

func TestAdminJWTRejectsUserAudienceToken(t *testing.T) {
	cfg := &config.Config{
		AdminJWTAccessSecret: "admin-access-secret",
		AdminJWTIssuer:       "nox-api",
		AdminJWTAudience:     "nox-admin",
	}
	rawToken, err := sharedtoken.SignWithOptions(uuid.New(), "user-access-jti", sharedtoken.AccessTokenType, cfg.AdminJWTAccessSecret, time.Minute, cfg.AdminJWTIssuer, "nox-client")
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}

	app := fiber.New()
	app.Get("/admin", AdminJWT(cfg, adminMiddlewareRepo{}), func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	req.Header.Set(fiber.HeaderAuthorization, "Bearer "+rawToken)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("execute request: %v", err)
	}
	if res.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", fiber.StatusUnauthorized, res.StatusCode)
	}
}

func TestAdminJWTRejectsInactiveMembership(t *testing.T) {
	cfg := &config.Config{
		AdminJWTAccessSecret: "admin-access-secret",
		AdminJWTIssuer:       "nox-api",
		AdminJWTAudience:     "nox-admin",
	}
	userID := uuid.New()
	rawToken, err := sharedtoken.SignWithOptions(userID, "admin-access-jti", sharedtoken.AccessTokenType, cfg.AdminJWTAccessSecret, time.Minute, cfg.AdminJWTIssuer, cfg.AdminJWTAudience)
	if err != nil {
		t.Fatalf("sign admin access token: %v", err)
	}

	app := fiber.New()
	app.Get("/admin", AdminJWT(cfg, adminMiddlewareRepo{
		membership: &models.AdminMembership{UserID: userID, Role: models.AdminRoleSupport, IsActive: false},
	}), func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	req.Header.Set(fiber.HeaderAuthorization, "Bearer "+rawToken)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("execute request: %v", err)
	}
	if res.StatusCode != fiber.StatusForbidden {
		t.Fatalf("expected status %d, got %d", fiber.StatusForbidden, res.StatusCode)
	}
}

type adminMiddlewareRepo struct {
	membership *models.AdminMembership
}

func (r adminMiddlewareRepo) FindMembershipByUserID(ctx context.Context, userID uuid.UUID) (*models.AdminMembership, error) {
	if r.membership == nil || r.membership.UserID != userID {
		return nil, nil
	}
	return r.membership, nil
}

func (r adminMiddlewareRepo) FindIdentityByUserID(ctx context.Context, userID uuid.UUID) (*models.AdminIdentity, error) {
	return nil, nil
}

func (r adminMiddlewareRepo) CreateAuditLog(ctx context.Context, params adminrepo.CreateAuditLogParams) error {
	return nil
}
