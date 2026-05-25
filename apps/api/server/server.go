package server

import (
	"context"
	"time"

	"github.com/emmanuella-codes/nox/auth/controllers"
	"github.com/emmanuella-codes/nox/auth/pipes"
	"github.com/emmanuella-codes/nox/auth/routers"
	"github.com/emmanuella-codes/nox/auth/services"
	commentcontrollers "github.com/emmanuella-codes/nox/comment/controllers"
	commentpipes "github.com/emmanuella-codes/nox/comment/pipes"
	commentrouters "github.com/emmanuella-codes/nox/comment/routers"
	"github.com/emmanuella-codes/nox/config"
	eventcontrollers "github.com/emmanuella-codes/nox/event/controllers"
	eventpipes "github.com/emmanuella-codes/nox/event/pipes"
	eventrouters "github.com/emmanuella-codes/nox/event/routers"
	likecontrollers "github.com/emmanuella-codes/nox/like/controllers"
	likepipes "github.com/emmanuella-codes/nox/like/pipes"
	likerouters "github.com/emmanuella-codes/nox/like/routers"
	"github.com/emmanuella-codes/nox/middleware"
	personacontrollers "github.com/emmanuella-codes/nox/persona/controllers"
	personapipes "github.com/emmanuella-codes/nox/persona/pipes"
	personarouters "github.com/emmanuella-codes/nox/persona/routers"
	postcontrollers "github.com/emmanuella-codes/nox/post/controllers"
	postpipes "github.com/emmanuella-codes/nox/post/pipes"
	postrouters "github.com/emmanuella-codes/nox/post/routers"
	"github.com/emmanuella-codes/nox/repositories"
	searchcontrollers "github.com/emmanuella-codes/nox/search/controllers"
	searchpipes "github.com/emmanuella-codes/nox/search/pipes"
	searchrouters "github.com/emmanuella-codes/nox/search/routers"
	shared_api "github.com/emmanuella-codes/nox/shared/api"
	"github.com/emmanuella-codes/nox/shared/mail"
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

	mailProvider := mail.NewBrevoProvider(mail.BrevoConfig{
		APIKey:      cfg.BrevoAPIKey,
		BaseURL:     cfg.BrevoBaseURL,
		SenderEmail: cfg.MailFromEmail,
		SenderName:  cfg.MailFromName,
	})

	authPipe := pipes.NewAuthPipe(pipes.AuthPipeDeps{
		UserRepo:     repos.User,
		HashService:  services.NewHashService(),
		OTPService:   services.NewOTPService(),
		EmailService: services.NewEmailService(mailProvider),
		TokenService: services.NewTokenService(cfg),
		Redis:        redisClient,
		Config:       cfg,
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "time": time.Now().Format(time.RFC3339)})
	})

	authController := controllers.NewAuthController(authPipe)
	personaController := personacontrollers.NewPersonaController(personapipes.NewPersonaPipe(repos.Persona))
	postController := postcontrollers.NewPostController(postpipes.NewPostPipe(repos.Post, repos.Persona, repos.Like))
	commentController := commentcontrollers.NewCommentController(commentpipes.NewCommentPipe(repos.Comment, repos.Persona, repos.Post))
	likeController := likecontrollers.NewLikeController(likepipes.NewLikePipe(repos.Like, repos.Persona, repos.Post))
	eventController := eventcontrollers.NewEventController(eventpipes.NewEventPipe(repos.Event, repos.Persona))
	searchController := searchcontrollers.NewSearchController(searchpipes.NewSearchPipe(repos.Search, repos.Like, repos.Persona))

	api := app.Group("/api/v1")

	shared_api.BaseRouter(api.Group("/auth"), routers.AuthRoutes(authController, redisClient))
	shared_api.BaseRouter(api.Group("/personas"), personarouters.PersonaRoutes(personaController, cfg, repos.Persona))
	shared_api.BaseRouter(api.Group("/posts"), postrouters.PostRoutes(postController, cfg, repos.Persona))
	shared_api.BaseRouter(api, commentrouters.CommentRoutes(commentController, cfg))
	shared_api.BaseRouter(api, likerouters.LikeRoutes(likeController, cfg))
	shared_api.BaseRouter(api.Group("/events"), eventrouters.EventRoutes(eventController, cfg))
	shared_api.BaseRouter(api.Group("/search"), searchrouters.SearchRoutes(searchController, cfg))

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
