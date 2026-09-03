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
		log.Fatal().Err(err).Msg("failed to bootstrap story worker")
	}
	defer app.Close()

	worker := NewWorker(app.Config, app.Repos.Story)
	if err := worker.Run(app.Context); err != nil {
		log.Fatal().Err(err).Msg("story cleanup worker failed")
	}
}
