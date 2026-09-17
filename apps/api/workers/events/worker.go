package main

import (
	"context"
	"time"

	"github.com/emmanuella-codes/nox/config"
	eventrepo "github.com/emmanuella-codes/nox/repositories/event"
	workerruntime "github.com/emmanuella-codes/nox/workers/runtime"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type Worker struct {
	cfg   *config.Config
	repo  eventrepo.EventRepository
	redis *redis.Client
}

func NewWorker(cfg *config.Config, repo eventrepo.EventRepository, redisClient *redis.Client) *Worker {
	return &Worker{cfg: cfg, repo: repo, redis: redisClient}
}

func (w *Worker) Run(ctx context.Context) error {
	return workerruntime.RunLoop(ctx, w.cfg.EventLifecycleInterval, func(ctx context.Context) error {
		if err := workerruntime.RecordHeartbeat(ctx, w.redis, "events", workerruntime.WorkerHeartbeatTTL(w.cfg.EventLifecycleInterval)); err != nil {
			log.Error().Err(err).Msg("event worker heartbeat failed")
		}
		if _, err := w.repo.AdvanceLifecycle(ctx, time.Now().UTC(), w.cfg.EventLifecycleBatchSize); err != nil {
			log.Error().Err(err).Msg("event lifecycle tick failed")
		}
		return nil
	})
}
