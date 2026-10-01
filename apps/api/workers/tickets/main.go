package main

import (
	"github.com/emmanuella-codes/nox/repositories/guest"
	"github.com/emmanuella-codes/nox/repositories/ticket"
	workerruntime "github.com/emmanuella-codes/nox/workers/runtime"
	"github.com/rs/zerolog/log"
)

func main() {
	ctx, stop := workerruntime.SignalContext()
	defer stop()
	app, err := workerruntime.Bootstrap(ctx, workerruntime.Options{ConnectRedis: true, RunMigrations: true, ServiceName: "nox-ticket-worker"})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to bootstrap ticket worker")
	}
	defer app.Close()
	if err := NewWorker(app.Config, ticket.TicketRepo, app.Redis, guest.GuestRepo).Run(app.Context); err != nil {
		log.Fatal().Err(err).Msg("ticket worker failed")
	}
}
