package routers

import (
	"strings"
	"time"

	"github.com/emmanuella-codes/nox/auth/controllers"
	"github.com/emmanuella-codes/nox/shared/api"
	"github.com/emmanuella-codes/nox/typings"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func AuthRoutes(controller *controllers.AuthController, redisClient *redis.Client) []api.RouterSchema {
	ipLimit := redisRateLimiter(redisClient, rateLimitConfig{
		Prefix: "rate:auth:ip",
		Max:    100,
		Window: time.Minute,
		Key: func(c *fiber.Ctx) string {
			return c.IP()
		},
	})
	emailLimit := redisRateLimiter(redisClient, rateLimitConfig{
		Prefix: "rate:auth:email",
		Max:    20,
		Window: time.Minute,
		Key:    emailRateLimitKey,
	})
	verifyEmailLimit := redisRateLimiter(redisClient, rateLimitConfig{
		Prefix: "rate:auth:verify_email",
		Max:    10,
		Window: time.Minute,
		Key:    emailRateLimitKey,
	})
	resendEmailLimit := redisRateLimiter(redisClient, rateLimitConfig{
		Prefix: "rate:auth:resend_verification",
		Max:    3,
		Window: 15 * time.Minute,
		Key:    emailRateLimitKey,
	})

	return []api.RouterSchema{
		{
			RouteMethod: api.RouteMethod("POST"),
			Path:        "/register",
			Middlewares: []typings.FiberMiddleware{ipLimit},
			Handler:     controller.Signup,
		},
		{
			RouteMethod: api.RouteMethod("POST"),
			Path:        "/login",
			Middlewares: []typings.FiberMiddleware{ipLimit, emailLimit},
			Handler:     controller.Login,
		},
		{
			RouteMethod: api.RouteMethod("POST"),
			Path:        "/verify-email",
			Middlewares: []typings.FiberMiddleware{ipLimit, verifyEmailLimit},
			Handler:     controller.VerifyEmail,
		},
		{
			RouteMethod: api.RouteMethod("POST"),
			Path:        "/resend-verification",
			Middlewares: []typings.FiberMiddleware{ipLimit, resendEmailLimit},
			Handler:     controller.ResendVerification,
		},
		{
			RouteMethod: api.RouteMethod("POST"),
			Path:        "/refresh",
			Middlewares: []typings.FiberMiddleware{ipLimit},
			Handler:     controller.Refresh,
		},
		{
			RouteMethod: api.RouteMethod("POST"),
			Path:        "/logout",
			Middlewares: []typings.FiberMiddleware{ipLimit},
			Handler:     controller.Logout,
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
