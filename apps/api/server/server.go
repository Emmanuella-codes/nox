package server

import (
	"context"
	"time"

	"github.com/emmanuella-codes/nox/auth/controllers"
	"github.com/emmanuella-codes/nox/auth/pipes"
	"github.com/emmanuella-codes/nox/auth/routers"
	"github.com/emmanuella-codes/nox/auth/services"
	comment_controllers "github.com/emmanuella-codes/nox/comment/controllers"
	comment_pipes "github.com/emmanuella-codes/nox/comment/pipes"
	comment_routers "github.com/emmanuella-codes/nox/comment/routers"
	"github.com/emmanuella-codes/nox/config"
	event_controllers "github.com/emmanuella-codes/nox/event/controllers"
	event_pipes "github.com/emmanuella-codes/nox/event/pipes"
	event_routers "github.com/emmanuella-codes/nox/event/routers"
	follow_controllers "github.com/emmanuella-codes/nox/follow/controllers"
	follow_pipes "github.com/emmanuella-codes/nox/follow/pipes"
	follow_routers "github.com/emmanuella-codes/nox/follow/routers"
	hashtag_controllers "github.com/emmanuella-codes/nox/hashtag/controllers"
	hashtag_pipes "github.com/emmanuella-codes/nox/hashtag/pipes"
	hashtag_routers "github.com/emmanuella-codes/nox/hashtag/routers"
	like_controllers "github.com/emmanuella-codes/nox/like/controllers"
	like_pipes "github.com/emmanuella-codes/nox/like/pipes"
	like_routers "github.com/emmanuella-codes/nox/like/routers"
	"github.com/emmanuella-codes/nox/middleware"
	persona_controllers "github.com/emmanuella-codes/nox/persona/controllers"
	persona_pipes "github.com/emmanuella-codes/nox/persona/pipes"
	persona_routers "github.com/emmanuella-codes/nox/persona/routers"
	post_controllers "github.com/emmanuella-codes/nox/post/controllers"
	post_pipes "github.com/emmanuella-codes/nox/post/pipes"
	post_routers "github.com/emmanuella-codes/nox/post/routers"
	"github.com/emmanuella-codes/nox/repositories"
	search_controllers "github.com/emmanuella-codes/nox/search/controllers"
	search_pipes "github.com/emmanuella-codes/nox/search/pipes"
	search_routers "github.com/emmanuella-codes/nox/search/routers"
	shared_api "github.com/emmanuella-codes/nox/shared/api"
	"github.com/emmanuella-codes/nox/shared/mail"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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

	app.Use(cors.New(cors.Config{
		AllowMethods: "GET,POST,PUT,PATCH,OPTIONS,DELETE",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

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
	personaController := persona_controllers.NewPersonaController(persona_pipes.NewPersonaPipe(repos.Persona))
	postController := post_controllers.NewPostController(post_pipes.NewPostPipe(repos.Post, repos.Persona, repos.Like, repos.Hashtag))
	commentController := comment_controllers.NewCommentController(comment_pipes.NewCommentPipe(repos.Comment, repos.Persona, repos.Post))
	likeController := like_controllers.NewLikeController(like_pipes.NewLikePipe(repos.Like, repos.Persona, repos.Post))
	eventController := event_controllers.NewEventController(event_pipes.NewEventPipe(repos.Event, repos.Persona))
	searchController := search_controllers.NewSearchController(search_pipes.NewSearchPipe(repos.Search, repos.Like, repos.Persona, repos.Hashtag, repos.Follow))
	followController := follow_controllers.NewFollowController(follow_pipes.NewFollowPipe(repos.Follow, repos.Persona))
	hashtagController := hashtag_controllers.NewHashtagController(hashtag_pipes.NewHashtagPipe(repos.Hashtag, repos.Persona, repos.Like))

	api := app.Group("/api/v1")

	shared_api.BaseRouter(api.Group("/auth"), routers.AuthRoutes(authController, redisClient))
	shared_api.BaseRouter(api.Group("/personas"), persona_routers.PersonaRoutes(personaController, cfg, repos.Persona))
	shared_api.BaseRouter(api.Group("/posts"), post_routers.PostRoutes(postController, cfg, repos.Persona))
	shared_api.BaseRouter(api, comment_routers.CommentRoutes(commentController, cfg))
	shared_api.BaseRouter(api, like_routers.LikeRoutes(likeController, cfg))
	shared_api.BaseRouter(api, follow_routers.FollowRoutes(followController, cfg))
	shared_api.BaseRouter(api.Group("/events"), event_routers.EventRoutes(eventController, cfg))
	shared_api.BaseRouter(api.Group("/search"), search_routers.SearchRoutes(searchController, cfg))
	shared_api.BaseRouter(api.Group("/hashtags"), hashtag_routers.HashtagRoutes(hashtagController))

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
