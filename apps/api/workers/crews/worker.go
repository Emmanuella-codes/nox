package main

import (
	"context"
	"time"

	"github.com/emmanuella-codes/nox/config"
	crewrepo "github.com/emmanuella-codes/nox/repositories/crew"
	workerruntime "github.com/emmanuella-codes/nox/workers/runtime"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type Worker struct {
	cfg   *config.Config
	repo  crewrepo.CrewRepository
	redis *redis.Client
}

func NewWorker(cfg *config.Config, repo crewrepo.CrewRepository, redisClient *redis.Client) *Worker {
	return &Worker{cfg: cfg, repo: repo, redis: redisClient}
}

func (w *Worker) Run(ctx context.Context) error {
	return workerruntime.RunLoop(ctx, w.cfg.CrewCleanupInterval, func(ctx context.Context) error {
		if err := workerruntime.RecordHeartbeat(ctx, w.redis, workerruntime.WorkerCrewsKey, workerruntime.WorkerHeartbeatTTL(w.cfg.CrewCleanupInterval)); err != nil {
			log.Error().Err(err).Msg("crew worker heartbeat failed")
		}
		if _, err := w.repo.DeleteExpiredCrewLocations(ctx, time.Now()); err != nil {
			log.Error().Err(err).Msg("crew location cleanup tick failed")
		}
		return nil
	})
}
