package middleware

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		status := c.Response().StatusCode()

		eventForStatus(status).
			Str("request_id", requestID(c)).
			Str("method", c.Method()).
			Str("path", c.Path()).
			Str("ip", c.IP()).
			Str("user_agent", c.Get("User-Agent")).
			Str("status", strconv.Itoa(status)).
			Int("bytes", len(c.Response().Body())).
			Dur("duration", time.Since(start)).
			Err(err).
			Msg("HTTP request")

		return err
	}
}

func eventForStatus(status int) *zerolog.Event {
	switch {
	case status >= 500:
		return log.Error()
	case status >= 400:
		return log.Warn()
	default:
		return log.Info()
	}
}
