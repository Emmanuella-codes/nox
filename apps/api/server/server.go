package server

import (
	"context"
	"time"

	"github.com/emmanuella-codes/nox/config"
	"github.com/emmanuella-codes/nox/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

func RunServer(ctx context.Context, cfg *config.Config, pool *pgxpool.Pool) {
	app := fiber.New(fiber.Config{
		AppName:      "nox-api",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	})

	app.Use(middleware.Logger())

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := app.ShutdownWithContext(shutdownCtx); err != nil {
			log.Fatal().Err(err).Msg("failed to shutdown server")
		}
		log.Info().Msg("server shutdown complete")
	}()

	log.Info().Str("port", cfg.Port).Str("env", cfg.Environment).Msg("server starting")
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatal().Err(err).Msg("server error")
	}
	log.Info().Msg("server stopped")
}
