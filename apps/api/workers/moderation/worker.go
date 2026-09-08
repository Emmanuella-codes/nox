package main

import (
	"context"

	"github.com/emmanuella-codes/nox/config"
	workerruntime "github.com/emmanuella-codes/nox/workers/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type Worker struct {
	cfg   *config.Config
	db    *pgxpool.Pool
	redis *redis.Client
}

func NewWorker(cfg *config.Config, db *pgxpool.Pool, redisClient *redis.Client) *Worker {
	return &Worker{cfg: cfg, db: db, redis: redisClient}
}

func (w *Worker) Run(ctx context.Context) error {
	return workerruntime.RunLoop(ctx, w.cfg.ModerationCleanupInterval, func(ctx context.Context) error {
		if err := workerruntime.RecordHeartbeat(ctx, w.redis, workerruntime.WorkerModerationKey, workerruntime.WorkerHeartbeatTTL(w.cfg.ModerationCleanupInterval)); err != nil {
			log.Error().Err(err).Msg("moderation worker heartbeat failed")
		}
		if _, err := w.db.Exec(ctx, `DELETE FROM message_report_evidence WHERE expires_at <= now()`); err != nil {
			log.Error().Err(err).Msg("message evidence cleanup tick failed")
		}
		return nil
	})
}
