package main

import (
	media_repo "github.com/emmanuella-codes/nox/repositories/media"
	workerruntime "github.com/emmanuella-codes/nox/workers/runtime"
	"github.com/rs/zerolog/log"
)

func main() {
	ctx, stop := workerruntime.SignalContext()
	defer stop()

	app, err := workerruntime.Bootstrap(ctx, workerruntime.Options{ConnectRedis: true, RunMigrations: true})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to bootstrap media worker")
	}
	defer app.Close()

	worker := NewWorker(app.Config, media_repo.NewCleanupRepository(app.DB))
	if err := worker.Run(app.Context); err != nil {
		log.Fatal().Err(err).Msg("media cleanup worker failed")
	}
}
