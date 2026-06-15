package media

import (
	"context"
	"time"

	"github.com/emmanuella-codes/nox/media/dtos"
	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MediaRepository interface {
	CreateMediaAsset(ctx context.Context, ownerUserID uuid.UUID, dto dtos.CreateMediaAssetDTO) (*models.MediaAsset, error)
	CreatePendingMediaAsset(ctx context.Context, ownerUserID uuid.UUID, storageKey string, playbackURL string, dto dtos.InitiateSetVideoUploadDTO) (*models.MediaAsset, error)
	CreatePendingStoryMediaAsset(ctx context.Context, ownerUserID uuid.UUID, storageKey string, playbackURL string, dto dtos.InitiateStoryVideoUploadDTO) (*models.MediaAsset, error)
	FindMediaAssetByID(ctx context.Context, mediaAssetID uuid.UUID) (*models.MediaAsset, error)
	MarkMediaAssetReady(ctx context.Context, mediaAssetID uuid.UUID, dto dtos.CompleteMediaProcessingDTO) (*models.MediaAsset, error)
	MarkMediaAssetFailed(ctx context.Context, mediaAssetID uuid.UUID) (*models.MediaAsset, error)
	DeleteOrphanedMediaAssets(ctx context.Context, olderThan time.Time, limit int) (int64, error)
}

func NewMediaRepository(db *pgxpool.Pool) MediaRepository {
	return newPgRepository(db)
}
