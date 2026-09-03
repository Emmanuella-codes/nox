package media

import (
	"context"
	"time"

	"github.com/emmanuella-codes/nox/media/dtos"
	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgRepository struct {
	db *pgxpool.Pool
}

func newPgRepository(db *pgxpool.Pool) *pgRepository {
	return &pgRepository{db: db}
}

func (r *pgRepository) CreateMediaAsset(ctx context.Context, ownerUserID uuid.UUID, dto dtos.CreateMediaAssetDTO) (*models.MediaAsset, error) {
	mediaKind := dto.MediaKind
	if mediaKind == "" {
		mediaKind = models.VideoMediaKind
	}
	row := r.db.QueryRow(ctx, `
		INSERT INTO media_assets (
			owner_user_id, owner_persona_id, media_kind, storage_key, playback_url, thumbnail_url,
			mime_type, duration_seconds, size_bytes, processing_status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, 'ready')
		RETURNING id, owner_user_id, owner_persona_id, media_kind, storage_key, playback_url,
		          COALESCE(thumbnail_url, ''), mime_type, duration_seconds, size_bytes,
		          processing_status, created_at, updated_at
	`, ownerUserID, dto.OwnerPersonaID, mediaKind, dto.StorageKey, dto.PlaybackURL, emptyToNil(dto.ThumbnailURL), dto.MimeType, dto.DurationSeconds, dto.SizeBytes)

	return scanMediaAsset(row)
}

func (r *pgRepository) CreatePostMediaAsset(ctx context.Context, ownerUserID uuid.UUID, dto dtos.ConfirmPostMediaUploadDTO) (*models.MediaAsset, error) {
	row := r.db.QueryRow(ctx, `
		INSERT INTO media_assets (
			owner_user_id, owner_persona_id, media_kind, storage_key, playback_url, thumbnail_url,
			mime_type, duration_seconds, size_bytes, processing_status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, 'ready')
		RETURNING id, owner_user_id, owner_persona_id, media_kind, storage_key, playback_url,
		          COALESCE(thumbnail_url, ''), mime_type, duration_seconds, size_bytes,
		          processing_status, created_at, updated_at
	`, ownerUserID, dto.OwnerPersonaID, dto.MediaKind, dto.PublicID, dto.SecureURL, emptyToNil(dto.ThumbnailURL), dto.MimeType, dto.DurationSeconds, dto.SizeBytes)

	return scanMediaAsset(row)
}

func (r *pgRepository) CreatePendingMediaAsset(ctx context.Context, ownerUserID uuid.UUID, storageKey string, playbackURL string, dto dtos.InitiateSetVideoUploadDTO) (*models.MediaAsset, error) {
	return r.createPendingAsset(ctx, ownerUserID, dto.OwnerPersonaID, models.VideoMediaKind, storageKey, playbackURL, dto.MimeType, 1, dto.SizeBytes)
}

func (r *pgRepository) CreatePendingStoryMediaAsset(ctx context.Context, ownerUserID uuid.UUID, storageKey string, playbackURL string, dto dtos.InitiateStoryMediaUploadDTO) (*models.MediaAsset, error) {
	return r.createPendingAsset(
		ctx,
		ownerUserID,
		dto.OwnerPersonaID,
		dto.MediaKind,
		storageKey,
		playbackURL,
		dto.MimeType,
		storyMediaDuration(dto.MediaKind, 1),
		dto.SizeBytes,
	)
}

func storyMediaDuration(kind models.MediaKind, durationSeconds int) int {
	if kind == models.ImageMediaKind {
		return 5
	}
	return durationSeconds
}

func (r *pgRepository) createPendingAsset(ctx context.Context, ownerUserID uuid.UUID, ownerPersonaID uuid.UUID, mediaKind models.MediaKind, storageKey string, playbackURL string, mimeType string, durationSeconds int, sizeBytes int64) (*models.MediaAsset, error) {
	row := r.db.QueryRow(ctx, `
		INSERT INTO media_assets (
			owner_user_id, owner_persona_id, media_kind, storage_key, playback_url, thumbnail_url,
			mime_type, duration_seconds, size_bytes, processing_status
		) VALUES ($1, $2, $3, $4, $5, NULL, $6, $7, $8, 'pending')
		RETURNING id, owner_user_id, owner_persona_id, media_kind, storage_key, playback_url,
		          COALESCE(thumbnail_url, ''), mime_type, duration_seconds, size_bytes,
		          processing_status, created_at, updated_at
	`, ownerUserID, ownerPersonaID, mediaKind, storageKey, playbackURL, mimeType, durationSeconds, sizeBytes)

	return scanMediaAsset(row)
}

func (r *pgRepository) FindMediaAssetByID(ctx context.Context, mediaAssetID uuid.UUID) (*models.MediaAsset, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, owner_user_id, owner_persona_id, media_kind, storage_key, playback_url,
		       COALESCE(thumbnail_url, ''), mime_type, duration_seconds, size_bytes,
		       processing_status, created_at, updated_at
		FROM media_assets
		WHERE id = $1
	`, mediaAssetID)

	return scanMediaAsset(row)
}

func (r *pgRepository) MarkMediaAssetReady(ctx context.Context, mediaAssetID uuid.UUID, dto dtos.CompleteMediaProcessingDTO) (*models.MediaAsset, error) {
	row := r.db.QueryRow(ctx, `
		UPDATE media_assets
		SET playback_url = $2,
		    thumbnail_url = $3,
		    mime_type = $4,
		    duration_seconds = $5,
		    size_bytes = $6,
		    processing_status = 'ready',
		    updated_at = now()
		WHERE id = $1
		RETURNING id, owner_user_id, owner_persona_id, media_kind, storage_key, playback_url,
		          COALESCE(thumbnail_url, ''), mime_type, duration_seconds, size_bytes,
		          processing_status, created_at, updated_at
	`, mediaAssetID, dto.PlaybackURL, emptyToNil(dto.ThumbnailURL), dto.MimeType, dto.DurationSeconds, dto.SizeBytes)

	return scanMediaAsset(row)
}

func (r *pgRepository) MarkMediaAssetFailed(ctx context.Context, mediaAssetID uuid.UUID) (*models.MediaAsset, error) {
	row := r.db.QueryRow(ctx, `
		UPDATE media_assets
		SET processing_status = 'failed',
		    updated_at = now()
		WHERE id = $1
		RETURNING id, owner_user_id, owner_persona_id, media_kind, storage_key, playback_url,
		          COALESCE(thumbnail_url, ''), mime_type, duration_seconds, size_bytes,
		          processing_status, created_at, updated_at
	`, mediaAssetID)

	return scanMediaAsset(row)
}

func (r *pgRepository) DeleteOrphanedMediaAssets(ctx context.Context, olderThan time.Time, limit int) (int64, error) {
	pendingDeleted, err := r.deleteOrphanedMediaAssetsByStatus(ctx, olderThan, limit, models.PendingMediaStatus)
	if err != nil {
		return 0, err
	}
	failedDeleted, err := r.deleteOrphanedMediaAssetsByStatus(ctx, olderThan, limit, models.FailedMediaStatus)
	if err != nil {
		return pendingDeleted, err
	}
	return pendingDeleted + failedDeleted, nil
}

type mediaAssetScanner interface {
	Scan(dest ...any) error
}

func scanMediaAsset(scanner mediaAssetScanner) (*models.MediaAsset, error) {
	var asset models.MediaAsset
	err := scanner.Scan(
		&asset.ID,
		&asset.OwnerUserID,
		&asset.OwnerPersonaID,
		&asset.MediaKind,
		&asset.StorageKey,
		&asset.PlaybackURL,
		&asset.ThumbnailURL,
		&asset.MimeType,
		&asset.DurationSeconds,
		&asset.SizeBytes,
		&asset.ProcessingStatus,
		&asset.CreatedAt,
		&asset.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

func emptyToNil(value string) any {
	if value == "" {
		return nil
	}
	return value
}
