package routers

import (
	"sync"
	"time"

	"github.com/emmanuella-codes/nox/auth/controllers"
	"github.com/emmanuella-codes/nox/shared/api"
	"github.com/emmanuella-codes/nox/typings"
	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(controller *controllers.AuthController) []api.RouterSchema {
	rateLimit := authRateLimiter()

	return []api.RouterSchema{
		{
			RouteMethod: api.RouteMethod("POST"),
			Path:        "/register",
			Middlewares: []typings.FiberMiddleware{rateLimit},
			Handler:     controller.Signup,
		},
		{
			RouteMethod: api.RouteMethod("POST"),
			Path:        "/login",
			Middlewares: []typings.FiberMiddleware{rateLimit},
			Handler:     controller.Login,
		},
		{
			RouteMethod: api.RouteMethod("POST"),
			Path:        "/refresh",
			Middlewares: []typings.FiberMiddleware{rateLimit},
			Handler:     controller.Refresh,
		},
		{
			RouteMethod: api.RouteMethod("POST"),
			Path:        "/logout",
			Middlewares: []typings.FiberMiddleware{rateLimit},
			Handler:     controller.Logout,
		},
	}
}

func authRateLimiter() fiber.Handler {
	const maxRequests = 20
	window := time.Minute

	type clientWindow struct {
		count   int
		expires time.Time
	}

	var mu sync.Mutex
	clients := make(map[string]clientWindow)

	return func(c *fiber.Ctx) error {
		now := time.Now()
		key := c.IP()

		mu.Lock()
		current := clients[key]
		if now.After(current.expires) {
			current = clientWindow{expires: now.Add(window)}
		}
		current.count++
		clients[key] = current
		allowed := current.count <= maxRequests
		retryAfter := time.Until(current.expires)
		mu.Unlock()

		if !allowed {
			c.Set(fiber.HeaderRetryAfter, retryAfter.Truncate(time.Second).String())
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"message": "too_many_requests",
			})
		}

		return c.Next()
	}
}
