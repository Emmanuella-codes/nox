package main

import (
	"context"
	"time"

	"github.com/emmanuella-codes/nox/config"
	media_repo "github.com/emmanuella-codes/nox/repositories/media"
	workerruntime "github.com/emmanuella-codes/nox/workers/runtime"
	"github.com/rs/zerolog/log"
)

type Worker struct {
	cfg  *config.Config
	repo media_repo.CleanupRepository
}

func NewWorker(cfg *config.Config, repo media_repo.CleanupRepository) *Worker {
	return &Worker{cfg: cfg, repo: repo}
}

func (w *Worker) Run(ctx context.Context) error {
	return workerruntime.RunLoop(ctx, w.cfg.MediaCleanupInterval, func(ctx context.Context) error {
		if err := w.tick(ctx); err != nil {
			log.Error().Err(err).Msg("media cleanup tick failed")
		}
		return nil
	})
}

func (w *Worker) tick(ctx context.Context) error {
	pendingBefore := time.Now().Add(-w.cfg.MediaPendingRetention)
	deletedPending, err := w.repo.DeleteOrphanedPendingMediaAssets(ctx, pendingBefore, w.cfg.MediaCleanupBatchSize)
	if err != nil {
		return err
	}
	failedBefore := time.Now().Add(-w.cfg.MediaFailedRetention)
	deletedFailed, err := w.repo.DeleteOrphanedFailedMediaAssets(ctx, failedBefore, w.cfg.MediaCleanupBatchSize)
	if err != nil {
		return err
	}
	if deletedPending > 0 || deletedFailed > 0 {
		log.Info().
			Int64("deleted_pending", deletedPending).
			Int64("deleted_failed", deletedFailed).
			Msg("media cleanup tick complete")
	}
	return nil
}
