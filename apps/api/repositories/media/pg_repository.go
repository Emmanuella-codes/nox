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
	return r.createPendingVideoAsset(ctx, ownerUserID, dto.OwnerPersonaID, storageKey, playbackURL, dto.MimeType, dto.SizeBytes)
}

func (r *pgRepository) CreatePendingStoryMediaAsset(ctx context.Context, ownerUserID uuid.UUID, storageKey string, playbackURL string, dto dtos.InitiateStoryVideoUploadDTO) (*models.MediaAsset, error) {
	return r.createPendingVideoAsset(ctx, ownerUserID, dto.OwnerPersonaID, storageKey, playbackURL, dto.MimeType, dto.SizeBytes)
}

func (r *pgRepository) createPendingVideoAsset(ctx context.Context, ownerUserID uuid.UUID, ownerPersonaID uuid.UUID, storageKey string, playbackURL string, mimeType string, sizeBytes int64) (*models.MediaAsset, error) {
	row := r.db.QueryRow(ctx, `
		INSERT INTO media_assets (
			owner_user_id, owner_persona_id, media_kind, storage_key, playback_url, thumbnail_url,
			mime_type, duration_seconds, size_bytes, processing_status
		) VALUES ($1, $2, 'video', $3, $4, NULL, $5, 1, $6, 'pending')
		RETURNING id, owner_user_id, owner_persona_id, media_kind, storage_key, playback_url,
		          COALESCE(thumbnail_url, ''), mime_type, duration_seconds, size_bytes,
		          processing_status, created_at, updated_at
	`, ownerUserID, ownerPersonaID, storageKey, playbackURL, mimeType, sizeBytes)

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
	if limit <= 0 {
		limit = 100
	}
	if limit > 1000 {
		limit = 1000
	}
	commandTag, err := r.db.Exec(ctx, `
		DELETE FROM media_assets
		WHERE id IN (
			SELECT media_assets.id
			FROM media_assets
			LEFT JOIN sets ON sets.media_asset_id = media_assets.id
			LEFT JOIN story_items ON story_items.media_asset_id = media_assets.id
			LEFT JOIN post_media_assets ON post_media_assets.media_asset_id = media_assets.id
			LEFT JOIN messages ON messages.media_asset_id = media_assets.id
			WHERE sets.id IS NULL
			  AND story_items.id IS NULL
			  AND post_media_assets.post_id IS NULL
			  AND messages.id IS NULL
			  AND media_assets.processing_status IN ('pending', 'failed')
			  AND media_assets.created_at < $1
			ORDER BY media_assets.created_at ASC
			LIMIT $2
		)
	`, olderThan, limit)
	if err != nil {
		return 0, err
	}
	return commandTag.RowsAffected(), nil
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
