package media

import (
	"context"

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
	row := r.db.QueryRow(ctx, `
		INSERT INTO media_assets (
			owner_user_id, owner_persona_id, media_kind, storage_key, playback_url, thumbnail_url,
			mime_type, duration_seconds, size_bytes, processing_status
		) VALUES ($1, $2, 'video', $3, $4, $5, $6, $7, $8, 'ready')
		RETURNING id, owner_user_id, owner_persona_id, media_kind, storage_key, playback_url,
		          COALESCE(thumbnail_url, ''), mime_type, duration_seconds, size_bytes,
		          processing_status, created_at, updated_at
	`, ownerUserID, dto.OwnerPersonaID, dto.StorageKey, dto.PlaybackURL, emptyToNil(dto.ThumbnailURL), dto.MimeType, dto.DurationSeconds, dto.SizeBytes)

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
