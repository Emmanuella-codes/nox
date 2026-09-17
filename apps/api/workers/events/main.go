package main

import (
	workerruntime "github.com/emmanuella-codes/nox/workers/runtime"
	"github.com/rs/zerolog/log"
)

func main() {
	ctx, stop := workerruntime.SignalContext()
	defer stop()
	app, err := workerruntime.Bootstrap(ctx, workerruntime.Options{ConnectRedis: true, RunMigrations: true, ServiceName: "nox-event-worker"})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to bootstrap event worker")
	}
	defer app.Close()
	if err := NewWorker(app.Config, app.Repos.Event, app.Redis).Run(app.Context); err != nil {
		log.Fatal().Err(err).Msg("event lifecycle worker failed")
	}
}
