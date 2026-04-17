package middleware

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()

		log.Info().
			Str("method", c.Method()).
			Str("path", c.Path()).
			Str("ip", c.IP()).
			Str("status", strconv.Itoa(c.Response().StatusCode())).
			Int("bytes", len(c.Response().Body())).
			Dur("duration", time.Since(start)).
			Err(err).
			Msg("HTTP request")

		return err
	}
}
