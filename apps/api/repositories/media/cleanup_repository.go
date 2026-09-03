package media

import (
	"context"
	"time"

	"github.com/emmanuella-codes/nox/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CleanupRepository interface {
	DeleteOrphanedPendingMediaAssets(ctx context.Context, olderThan time.Time, limit int) (int64, error)
	DeleteOrphanedFailedMediaAssets(ctx context.Context, olderThan time.Time, limit int) (int64, error)
}

func NewCleanupRepository(db *pgxpool.Pool) CleanupRepository {
	return newPgRepository(db)
}

func (r *pgRepository) DeleteOrphanedPendingMediaAssets(ctx context.Context, olderThan time.Time, limit int) (int64, error) {
	return r.deleteOrphanedMediaAssetsByStatus(ctx, olderThan, limit, models.PendingMediaStatus)
}

func (r *pgRepository) DeleteOrphanedFailedMediaAssets(ctx context.Context, olderThan time.Time, limit int) (int64, error) {
	return r.deleteOrphanedMediaAssetsByStatus(ctx, olderThan, limit, models.FailedMediaStatus)
}
