package post

import (
	"context"

	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
)

// FindMediaAssetsByPostIDs fetches media grouped by post id.
func (r *pgRepository) FindMediaAssetsByPostIDs(ctx context.Context, postIDs []uuid.UUID) (map[uuid.UUID][]*models.MediaAsset, error) {
	assetsByPost := make(map[uuid.UUID][]*models.MediaAsset)
	if len(postIDs) == 0 {
		return assetsByPost, nil
	}
	rows, err := r.db.Query(ctx, `
		SELECT pma.post_id, ma.id, ma.owner_user_id, ma.owner_persona_id, ma.media_kind, ma.storage_key,
		       ma.playback_url, COALESCE(ma.thumbnail_url, ''), ma.mime_type, ma.duration_seconds,
		       ma.size_bytes, ma.processing_status, ma.created_at, ma.updated_at
		FROM post_media_assets pma
		INNER JOIN media_assets ma ON ma.id = pma.media_asset_id
		WHERE pma.post_id = ANY($1)
		ORDER BY pma.position ASC
	`, postIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var postID uuid.UUID
		var asset models.MediaAsset
		if err := rows.Scan(
			&postID,
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
		); err != nil {
			return nil, err
		}
		assetsByPost[postID] = append(assetsByPost[postID], &asset)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return assetsByPost, nil
}
