package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/emmanuella-codes/nox/config"
	"github.com/emmanuella-codes/nox/db"
	"github.com/emmanuella-codes/nox/repositories"
	"github.com/rs/zerolog/log"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}
	if err := db.RunMigrations(cfg.DatabaseURL); err != nil {
		log.Fatal().Err(err).Msg("failed to run database migrations")
	}
	pool, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect database")
	}
	defer pool.Close()

	repos := repositories.Init(pool)
	provider := newProvider(cfg)
	worker := NewWorker(cfg, repos.Notification, provider)
	if err := worker.Run(ctx); err != nil {
		log.Fatal().Err(err).Msg("notification worker failed")
	}
}
