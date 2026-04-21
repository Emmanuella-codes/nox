package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/emmanuella-codes/nox/config"
	"github.com/emmanuella-codes/nox/db"
	"github.com/emmanuella-codes/nox/repositories"
	"github.com/emmanuella-codes/nox/server"
	"github.com/rs/zerolog/log"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	pool, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect database")
	}
	defer pool.Close()

	redisClient, err := db.ConnectRedis(ctx, cfg.RedisURL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect redis")
	}
	defer func() {
		if err := redisClient.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close redis")
		}
	}()

	repos := repositories.Init(pool)

	server.RunServer(ctx, cfg, redisClient, repos)
}
