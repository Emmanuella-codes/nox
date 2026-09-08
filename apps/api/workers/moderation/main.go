package main

import (
	workerruntime "github.com/emmanuella-codes/nox/workers/runtime"
	"github.com/rs/zerolog/log"
)

func main() {
	ctx, stop := workerruntime.SignalContext()
	defer stop()
	app, err := workerruntime.Bootstrap(ctx, workerruntime.Options{ConnectRedis: true, RunMigrations: true, ServiceName: "nox-moderation-worker"})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to bootstrap moderation worker")
	}
	defer app.Close()
	if err := NewWorker(app.Config, app.DB, app.Redis).Run(app.Context); err != nil {
		log.Fatal().Err(err).Msg("moderation worker failed")
	}
}
