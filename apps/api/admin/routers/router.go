package routers

import (
	"strings"
	"time"

	"github.com/emmanuella-codes/nox/admin/controllers"
	"github.com/emmanuella-codes/nox/config"
	"github.com/emmanuella-codes/nox/middleware"
	adminrepo "github.com/emmanuella-codes/nox/repositories/admin"
	"github.com/emmanuella-codes/nox/shared/api"
	"github.com/emmanuella-codes/nox/typings"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func AdminRoutes(controller *controllers.AdminController, cfg *config.Config, redisClient *redis.Client, adminRepo adminrepo.AdminRepository) []api.RouterSchema {
	ipLimit := redisRateLimiter(redisClient, rateLimitConfig{
		Prefix: "rate:admin_auth:ip",
		Max:    60,
		Window: time.Minute,
		Key: func(c *fiber.Ctx) string {
			return c.IP()
		},
	})
	emailLimit := redisRateLimiter(redisClient, rateLimitConfig{
		Prefix: "rate:admin_auth:email",
		Max:    12,
		Window: time.Minute,
		Key:    emailRateLimitKey,
	})

	routes := []api.RouterSchema{
		{
			RouteMethod: api.RouteMethod("POST"),
			Path:        "/auth/login",
			Middlewares: []typings.FiberMiddleware{ipLimit, emailLimit},
			Handler:     controller.Login,
		},
		{
			RouteMethod: api.RouteMethod("POST"),
			Path:        "/auth/refresh",
			Middlewares: []typings.FiberMiddleware{ipLimit},
			Handler:     controller.Refresh,
		},
		{
			RouteMethod: api.RouteMethod("POST"),
			Path:        "/auth/logout",
			Middlewares: []typings.FiberMiddleware{ipLimit},
			Handler:     controller.Logout,
		},
		{
			RouteMethod: api.RouteMethod("GET"),
			Path:        "/me",
			Middlewares: []typings.FiberMiddleware{middleware.AdminJWT(cfg, adminRepo)},
			Handler:     controller.Me,
		},
		{
			RouteMethod: api.RouteMethod("GET"),
			Path:        "/health",
			Middlewares: []typings.FiberMiddleware{middleware.AdminJWT(cfg, adminRepo)},
			Handler:     controller.Health,
		},
		{
			RouteMethod: api.RouteMethod("GET"),
			Path:        "/dashboard",
			Middlewares: []typings.FiberMiddleware{middleware.AdminJWT(cfg, adminRepo)},
			Handler:     controller.Dashboard,
		},
		{
			RouteMethod: api.RouteMethod("GET"),
			Path:        "/staff",
			Middlewares: []typings.FiberMiddleware{middleware.AdminJWT(cfg, adminRepo)},
			Handler:     controller.ListAdminUsers,
		},
		{
			RouteMethod: api.RouteMethod("GET"),
			Path:        "/staff/:userID",
			Middlewares: []typings.FiberMiddleware{middleware.AdminJWT(cfg, adminRepo)},
			Handler:     controller.GetAdminUser,
		},
		{
			RouteMethod: api.RouteMethod("POST"),
			Path:        "/staff",
			Middlewares: []typings.FiberMiddleware{middleware.AdminJWT(cfg, adminRepo)},
			Handler:     controller.CreateAdminUser,
		},
		{
			RouteMethod: api.RouteMethod("PATCH"),
			Path:        "/staff/:userID/role",
			Middlewares: []typings.FiberMiddleware{middleware.AdminJWT(cfg, adminRepo)},
			Handler:     controller.UpdateAdminUserRole,
		},
		{
			RouteMethod: api.RouteMethod("PATCH"),
			Path:        "/staff/:userID/status",
			Middlewares: []typings.FiberMiddleware{middleware.AdminJWT(cfg, adminRepo)},
			Handler:     controller.UpdateAdminUserStatus,
		},
		{
			RouteMethod: api.RouteMethod("GET"),
			Path:        "/users",
			Middlewares: []typings.FiberMiddleware{middleware.AdminJWT(cfg, adminRepo)},
			Handler:     controller.ListUsers,
		},
		{
			RouteMethod: api.RouteMethod("GET"),
			Path:        "/moderation/actions",
			Middlewares: []typings.FiberMiddleware{middleware.AdminJWT(cfg, adminRepo)},
			Handler:     controller.ListModerationActions,
		},
		{
			RouteMethod: api.RouteMethod("GET"),
			Path:        "/reports",
			Middlewares: []typings.FiberMiddleware{middleware.AdminJWT(cfg, adminRepo)},
			Handler:     controller.ListReports,
		},
		{
			RouteMethod: api.RouteMethod("GET"),
			Path:        "/reports/:reportID",
			Middlewares: []typings.FiberMiddleware{middleware.AdminJWT(cfg, adminRepo)},
			Handler:     controller.GetReport,
		},
		{
			RouteMethod: api.RouteMethod("PATCH"),
			Path:        "/reports/:reportID",
			Middlewares: []typings.FiberMiddleware{middleware.AdminJWT(cfg, adminRepo)},
			Handler:     controller.UpdateReport,
		},
		{
			RouteMethod: api.RouteMethod("POST"),
			Path:        "/reports/:reportID/actions",
			Middlewares: []typings.FiberMiddleware{middleware.AdminJWT(cfg, adminRepo)},
			Handler:     controller.ResolveReport,
		},
		{
			RouteMethod: api.RouteMethod("PATCH"),
			Path:        "/moderation/:entityType/:entityID",
			Middlewares: []typings.FiberMiddleware{middleware.AdminJWT(cfg, adminRepo)},
			Handler:     controller.ModerateContent,
		},
		{
			RouteMethod: api.RouteMethod("GET"),
			Path:        "/users/:userID",
			Middlewares: []typings.FiberMiddleware{middleware.AdminJWT(cfg, adminRepo)},
			Handler:     controller.GetUser,
		},
		{
			RouteMethod: api.RouteMethod("PATCH"),
			Path:        "/users/:userID/status",
			Middlewares: []typings.FiberMiddleware{middleware.AdminJWT(cfg, adminRepo)},
			Handler:     controller.UpdateUserStatus,
		},
		{
			RouteMethod: api.RouteMethod("POST"),
			Path:        "/users/:userID/revoke-sessions",
			Middlewares: []typings.FiberMiddleware{middleware.AdminJWT(cfg, adminRepo)},
			Handler:     controller.RevokeUserSessions,
		},
		{
			RouteMethod: api.RouteMethod("POST"),
			Path:        "/users/:userID/resend-verification",
			Middlewares: []typings.FiberMiddleware{middleware.AdminJWT(cfg, adminRepo)},
			Handler:     controller.ResendUserVerification,
		},
		{
			RouteMethod: api.RouteMethod("POST"),
			Path:        "/users/:userID/mark-email-verified",
			Middlewares: []typings.FiberMiddleware{middleware.AdminJWT(cfg, adminRepo)},
			Handler:     controller.MarkUserEmailVerified,
		},
		{
			RouteMethod: api.RouteMethod("GET"),
			Path:        "/hashtags",
			Middlewares: []typings.FiberMiddleware{middleware.AdminJWT(cfg, adminRepo)},
			Handler:     controller.ListAdminHashtags,
		},
		{
			RouteMethod: api.RouteMethod("PATCH"),
			Path:        "/hashtags/:tag/moderation",
			Middlewares: []typings.FiberMiddleware{middleware.AdminJWT(cfg, adminRepo)},
			Handler:     controller.ModerateHashtag,
		},
		{
			RouteMethod: api.RouteMethod("POST"),
			Path:        "/hashtags/:tag/suppress",
			Middlewares: []typings.FiberMiddleware{middleware.AdminJWT(cfg, adminRepo)},
			Handler:     controller.SuppressHashtag,
		},
		{
			RouteMethod: api.RouteMethod("DELETE"),
			Path:        "/hashtags/:tag/suppress",
			Middlewares: []typings.FiberMiddleware{middleware.AdminJWT(cfg, adminRepo)},
			Handler:     controller.UnsuppressHashtag,
		},
		{
			RouteMethod: api.RouteMethod("GET"),
			Path:        "/search/suppressions",
			Middlewares: []typings.FiberMiddleware{middleware.AdminJWT(cfg, adminRepo)},
			Handler:     controller.ListSearchSuppressions,
		},
		{
			RouteMethod: api.RouteMethod("POST"),
			Path:        "/search/suppressions",
			Middlewares: []typings.FiberMiddleware{middleware.AdminJWT(cfg, adminRepo)},
			Handler:     controller.CreateSearchSuppression,
		},
		{
			RouteMethod: api.RouteMethod("DELETE"),
			Path:        "/search/suppressions/:suppressionID",
			Middlewares: []typings.FiberMiddleware{middleware.AdminJWT(cfg, adminRepo)},
			Handler:     controller.DeleteSearchSuppression,
		},
		{
			RouteMethod: api.RouteMethod("POST"),
			Path:        "/search/reindex",
			Middlewares: []typings.FiberMiddleware{middleware.AdminJWT(cfg, adminRepo)},
			Handler:     controller.ReindexSearch,
		},
		{
			RouteMethod: api.RouteMethod("GET"),
			Path:        "/media-assets",
			Middlewares: []typings.FiberMiddleware{middleware.AdminJWT(cfg, adminRepo)},
			Handler:     controller.ListAdminMedia,
		},
		{
			RouteMethod: api.RouteMethod("GET"),
			Path:        "/media-assets/:mediaAssetID",
			Middlewares: []typings.FiberMiddleware{middleware.AdminJWT(cfg, adminRepo)},
			Handler:     controller.GetAdminMedia,
		},
		{
			RouteMethod: api.RouteMethod("POST"),
			Path:        "/media-assets/:mediaAssetID/retry",
			Middlewares: []typings.FiberMiddleware{middleware.AdminJWT(cfg, adminRepo)},
			Handler:     controller.RetryAdminMedia,
		},
		{
			RouteMethod: api.RouteMethod("PATCH"),
			Path:        "/media-assets/:mediaAssetID/status",
			Middlewares: []typings.FiberMiddleware{middleware.AdminJWT(cfg, adminRepo)},
			Handler:     controller.CorrectAdminMedia,
		},
		{
			RouteMethod: api.RouteMethod("PATCH"),
			Path:        "/media-assets/:mediaAssetID/moderation",
			Middlewares: []typings.FiberMiddleware{middleware.AdminJWT(cfg, adminRepo)},
			Handler:     controller.ModerateAdminMedia,
		},
		{
			RouteMethod: api.RouteMethod("POST"),
			Path:        "/media-assets/orphans/cleanup",
			Middlewares: []typings.FiberMiddleware{middleware.AdminJWT(cfg, adminRepo)},
			Handler:     controller.CleanupAdminMedia,
		},
	}
	return append(routes, adminContentRoutes(controller, cfg, adminRepo)...)
}

type rateLimitConfig struct {
	Prefix string
	Max    int64
	Window time.Duration
	Key    func(*fiber.Ctx) string
}

func redisRateLimiter(redisClient *redis.Client, cfg rateLimitConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		keyPart := cfg.Key(c)
		if keyPart == "" {
			return c.Next()
		}

		key := cfg.Prefix + ":" + keyPart
		count, err := redisClient.Incr(c.Context(), key).Result()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "internal_error",
			})
		}
		if count == 1 {
			if err := redisClient.Expire(c.Context(), key, cfg.Window).Err(); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"success": false,
					"message": "internal_error",
				})
			}
		}

		if count > cfg.Max {
			ttl, err := redisClient.TTL(c.Context(), key).Result()
			if err == nil && ttl > 0 {
				c.Set(fiber.HeaderRetryAfter, ttl.Truncate(time.Second).String())
			}
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"message": "too_many_requests",
			})
		}

		return c.Next()
	}
}

func emailRateLimitKey(c *fiber.Ctx) string {
	var payload struct {
		Email string `json:"email"`
	}
	if err := c.BodyParser(&payload); err != nil {
		return ""
	}
	return strings.ToLower(strings.TrimSpace(payload.Email))
}
