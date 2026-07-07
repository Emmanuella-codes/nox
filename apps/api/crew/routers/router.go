package routers

import (
	"strings"
	"time"

	"github.com/emmanuella-codes/nox/config"
	"github.com/emmanuella-codes/nox/crew/controllers"
	"github.com/emmanuella-codes/nox/middleware"
	"github.com/emmanuella-codes/nox/shared/api"
	"github.com/emmanuella-codes/nox/typings"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func CrewRoutes(controller *controllers.CrewController, cfg *config.Config, redisClient *redis.Client) []api.RouterSchema {
	auth := []typings.FiberMiddleware{middleware.JWT(cfg)}
	joinLimit := redisRateLimiter(redisClient, rateLimitConfig{
		Prefix: "rate:crew:join",
		Max:    10,
		Window: time.Minute,
		Key:    joinRateLimitKey,
	})
	return []api.RouterSchema{
		{RouteMethod: "POST", Path: "/events/:eventID/crews", Middlewares: auth, Handler: controller.CreateCrew},
		{RouteMethod: "GET", Path: "/events/:eventID/crews/me", Middlewares: auth, Handler: controller.ListMyEventCrews},
		{RouteMethod: "POST", Path: "/crews/join", Middlewares: append(auth, joinLimit), Handler: controller.JoinCrew},
		{RouteMethod: "GET", Path: "/crews/:crewID", Middlewares: auth, Handler: controller.GetCrew},
		{RouteMethod: "POST", Path: "/crews/:crewID/leave", Middlewares: auth, Handler: controller.LeaveCrew},
		{RouteMethod: "POST", Path: "/crews/:crewID/end", Middlewares: auth, Handler: controller.EndCrew},
		{RouteMethod: "PATCH", Path: "/crews/:crewID/location-sharing", Middlewares: auth, Handler: controller.UpdateLocationSharing},
		{RouteMethod: "PUT", Path: "/crews/:crewID/location", Middlewares: auth, Handler: controller.UpdateLocation},
		{RouteMethod: "GET", Path: "/crews/:crewID/locations", Middlewares: auth, Handler: controller.ListLocations},
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
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "internal_error"})
		}
		if count == 1 {
			if err := redisClient.Expire(c.Context(), key, cfg.Window).Err(); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "internal_error"})
			}
		}
		if count > cfg.Max {
			ttl, err := redisClient.TTL(c.Context(), key).Result()
			if err == nil && ttl > 0 {
				c.Set(fiber.HeaderRetryAfter, ttl.Truncate(time.Second).String())
			}
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"success": false, "message": "too_many_requests"})
		}
		return c.Next()
	}
}

func joinRateLimitKey(c *fiber.Ctx) string {
	var payload struct {
		JoinCode string `json:"join_code"`
	}
	if err := c.BodyParser(&payload); err != nil {
		return c.IP()
	}
	joinCode := strings.ToUpper(strings.TrimSpace(payload.JoinCode))
	if joinCode == "" {
		return c.IP()
	}
	return c.IP() + ":" + joinCode
}
