package middleware

import (
	"runtime/debug"

	"github.com/emmanuella-codes/nox/shared"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func Recover() fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		defer func() {
			if recovered := recover(); recovered != nil {
				log.Error().
					Any("panic", recovered).
					Str("request_id", requestID(c)).
					Bytes("stack", debug.Stack()).
					Msg("panic recovered")

				err = c.Status(fiber.StatusInternalServerError).JSON(shared.PipeRes[any]{
					Success: false,
					Message: "internal_error",
					Data:    nil,
				})
			}
		}()
		return c.Next()
	}
}

func requestID(c *fiber.Ctx) string {
	if value, ok := c.Locals(RequestIDLocalKey).(string); ok {
		return value
	}
	return ""
}
