package server

import (
	"context"
	"time"

	"github.com/emmanuella-codes/nox/auth/controllers"
	"github.com/emmanuella-codes/nox/auth/pipes"
	"github.com/emmanuella-codes/nox/auth/routers"
	"github.com/emmanuella-codes/nox/auth/services"
	"github.com/emmanuella-codes/nox/config"
	"github.com/emmanuella-codes/nox/middleware"
	"github.com/emmanuella-codes/nox/repositories"
	shared_api "github.com/emmanuella-codes/nox/shared/api"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

func RunServer(ctx context.Context, cfg *config.Config, redisClient *redis.Client, repos *repositories.Repositories) {
	app := fiber.New(fiber.Config{
		AppName:      "nox-api",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	})

	app.Use(middleware.Logger())

	authPipe := pipes.NewAuthPipe(pipes.AuthPipeDeps{
		UserRepo:     repos.User,
		HashService:  services.NewHashService(),
		TokenService: services.NewTokenService(cfg),
		Redis:        redisClient,
		Config:       cfg,
	})

	authController := controllers.NewAuthController(authPipe)

	api := app.Group("/api/v1")
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "time": time.Now().Format(time.RFC3339)})
	})

	shared_api.BaseRouter(api.Group("/auth"), routers.AuthRoutes(authController))

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
