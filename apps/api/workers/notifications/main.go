package main

import (
	workerruntime "github.com/emmanuella-codes/nox/workers/runtime"
	"github.com/rs/zerolog/log"
)

func main() {
	ctx, stop := workerruntime.SignalContext()
	defer stop()

	app, err := workerruntime.Bootstrap(ctx, workerruntime.Options{ConnectRedis: true, RunMigrations: true})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to bootstrap notification worker")
	}
	defer app.Close()

	provider := newProvider(app.Config)
	worker := NewWorker(app.Config, app.Repos.Notification, provider)
	if err := worker.Run(app.Context); err != nil {
		log.Fatal().Err(err).Msg("notification worker failed")
	}
}
