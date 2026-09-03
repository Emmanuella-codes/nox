package main

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/emmanuella-codes/nox/config"
)

type cleanupRepoStub struct {
	pendingOlderThan time.Time
	failedOlderThan  time.Time
	pendingLimit     int
	failedLimit      int
	pendingDeleted   int64
	failedDeleted    int64
	pendingErr       error
	failedErr        error
}

func (r *cleanupRepoStub) DeleteOrphanedPendingMediaAssets(ctx context.Context, olderThan time.Time, limit int) (int64, error) {
	r.pendingOlderThan = olderThan
	r.pendingLimit = limit
	return r.pendingDeleted, r.pendingErr
}

func (r *cleanupRepoStub) DeleteOrphanedFailedMediaAssets(ctx context.Context, olderThan time.Time, limit int) (int64, error) {
	r.failedOlderThan = olderThan
	r.failedLimit = limit
	return r.failedDeleted, r.failedErr
}

func TestTickUsesConfiguredMediaRetentionWindows(t *testing.T) {
	repo := &cleanupRepoStub{}
	cfg := &config.Config{
		MediaCleanupBatchSize: 12,
		MediaPendingRetention: 24 * time.Hour,
		MediaFailedRetention:  7 * 24 * time.Hour,
	}
	worker := NewWorker(cfg, repo)
	before := time.Now()

	if err := worker.tick(context.Background()); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	after := time.Now()
	if repo.pendingLimit != 12 || repo.failedLimit != 12 {
		t.Fatalf("expected batch size 12, got pending=%d failed=%d", repo.pendingLimit, repo.failedLimit)
	}
	assertBetween(t, repo.pendingOlderThan, before.Add(-24*time.Hour), after.Add(-24*time.Hour))
	assertBetween(t, repo.failedOlderThan, before.Add(-7*24*time.Hour), after.Add(-7*24*time.Hour))
}

func TestTickStopsOnPendingCleanupError(t *testing.T) {
	repo := &cleanupRepoStub{pendingErr: errors.New("pending failed")}
	worker := NewWorker(&config.Config{MediaCleanupBatchSize: 5, MediaPendingRetention: time.Hour, MediaFailedRetention: 2 * time.Hour}, repo)

	if err := worker.tick(context.Background()); err == nil {
		t.Fatal("expected pending cleanup error")
	}
	if !repo.failedOlderThan.IsZero() {
		t.Fatal("expected failed cleanup not to run after pending error")
	}
}

func TestTickReturnsFailedCleanupError(t *testing.T) {
	repo := &cleanupRepoStub{failedErr: errors.New("failed cleanup")}
	worker := NewWorker(&config.Config{MediaCleanupBatchSize: 5, MediaPendingRetention: time.Hour, MediaFailedRetention: 2 * time.Hour}, repo)

	if err := worker.tick(context.Background()); err == nil {
		t.Fatal("expected failed cleanup error")
	}
}

func assertBetween(t *testing.T, value time.Time, earliest time.Time, latest time.Time) {
	t.Helper()
	if value.Before(earliest) || value.After(latest) {
		t.Fatalf("expected %s between %s and %s", value, earliest, latest)
	}
}
