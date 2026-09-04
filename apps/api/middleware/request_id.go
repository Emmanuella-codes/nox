package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const RequestIDLocalKey = "request_id"

func RequestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestID := c.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.NewString()
		}
		c.Locals(RequestIDLocalKey, requestID)
		c.Set("X-Request-ID", requestID)
		return c.Next()
	}
}
