package main

import (
	"context"
	"github.com/emmanuella-codes/nox/config"
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
}

func NewWorker(cfg *config.Config, repo ticketrepo.TicketRepository, redis *redis.Client) *Worker {
	return &Worker{cfg: cfg, repo: repo, redis: redis}
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
		return nil
	})
}
