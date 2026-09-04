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

	return []api.RouterSchema{
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
	}
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
