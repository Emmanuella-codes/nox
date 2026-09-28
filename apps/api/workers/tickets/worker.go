package main

import (
	"context"
	"github.com/emmanuella-codes/nox/config"
	guestrepo "github.com/emmanuella-codes/nox/repositories/guest"
	ticketrepo "github.com/emmanuella-codes/nox/repositories/ticket"
	workerruntime "github.com/emmanuella-codes/nox/workers/runtime"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"time"
)

type Worker struct {
	cfg   *config.Config
	repo  ticketrepo.TicketRepository
	redis *redis.Client
	guest guestrepo.GuestRepository
}

func NewWorker(cfg *config.Config, repo ticketrepo.TicketRepository, redis *redis.Client, guests ...guestrepo.GuestRepository) *Worker {
	var guest guestrepo.GuestRepository
	if len(guests) > 0 {
		guest = guests[0]
	}
	return &Worker{cfg: cfg, repo: repo, redis: redis, guest: guest}
}
func (w *Worker) Run(ctx context.Context) error {
	return workerruntime.RunLoop(ctx, w.cfg.TicketInventoryInterval, func(ctx context.Context) error {
		if err := workerruntime.RecordHeartbeat(ctx, w.redis, "tickets", workerruntime.WorkerHeartbeatTTL(w.cfg.TicketInventoryInterval)); err != nil {
			log.Error().Err(err).Msg("ticket worker heartbeat failed")
		}
		if _, err := w.repo.ReleaseExpired(ctx, time.Now().UTC(), w.cfg.TicketInventoryBatchSize); err != nil {
			log.Error().Err(err).Msg("ticket inventory cleanup failed")
		}
		if _, err := w.repo.AdvanceSales(ctx, time.Now().UTC(), w.cfg.TicketInventoryBatchSize); err != nil {
			log.Error().Err(err).Msg("ticket sales advancement failed")
		}
		if w.guest != nil {
			if _, err := w.guest.ExpireInvitations(ctx, time.Now().UTC()); err != nil {
				log.Error().Err(err).Msg("guest invitation expiry failed")
			}
			if _, err := w.guest.MarkNoShows(ctx, time.Now().UTC()); err != nil {
				log.Error().Err(err).Msg("guest no-show reconciliation failed")
			}
		}
		return nil
	})
}
