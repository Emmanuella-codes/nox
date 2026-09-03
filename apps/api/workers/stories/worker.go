package main

import (
	"context"
	"time"

	"github.com/emmanuella-codes/nox/config"
	story_repo "github.com/emmanuella-codes/nox/repositories/story"
	workerruntime "github.com/emmanuella-codes/nox/workers/runtime"
	"github.com/rs/zerolog/log"
)

type Worker struct {
	cfg  *config.Config
	repo story_repo.StoryRepository
}

// Builds one story cleanup worker from the shared config and repositories.
func NewWorker(cfg *config.Config, repo story_repo.StoryRepository) *Worker {
	return &Worker{cfg: cfg, repo: repo}
}

// Polls story cleanup work until shutdown.
func (w *Worker) Run(ctx context.Context) error {
	return workerruntime.RunLoop(ctx, w.cfg.StoryCleanupInterval, func(ctx context.Context) error {
		if err := w.tick(ctx); err != nil {
			log.Error().Err(err).Msg("story cleanup tick failed")
		}
		return nil
	})
}

// Rejects stale requests and purges expired non-highlighted story data in batches.
func (w *Worker) tick(ctx context.Context) error {
	rejected, err := w.repo.RejectPendingContributionRequestsForClosedStories(ctx, w.cfg.StoryCleanupBatchSize)
	if err != nil {
		return err
	}
	expiredBefore := time.Now().Add(-w.cfg.StoryExpiryRetention)
	deletedItems, err := w.repo.DeleteExpiredNonHighlightedStoryItems(ctx, expiredBefore, w.cfg.StoryCleanupBatchSize)
	if err != nil {
		return err
	}
	deletedStories, err := w.repo.DeleteRetainedEmptyStories(ctx, expiredBefore, w.cfg.StoryCleanupBatchSize)
	if err != nil {
		return err
	}
	if rejected > 0 || deletedItems > 0 || deletedStories > 0 {
		log.Info().
			Int64("rejected_requests", rejected).
			Int64("deleted_items", deletedItems).
			Int64("deleted_stories", deletedStories).
			Msg("story cleanup tick complete")
	}
	return nil
}
