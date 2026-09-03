package media

import (
	"context"
	"time"

	"github.com/emmanuella-codes/nox/models"
)

func (r *pgRepository) deleteOrphanedMediaAssetsByStatus(ctx context.Context, olderThan time.Time, limit int, status models.MediaProcessingStatus) (int64, error) {
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
			LEFT JOIN message_attachments ON message_attachments.media_asset_id = media_assets.id
			WHERE sets.id IS NULL
			  AND story_items.id IS NULL
			  AND post_media_assets.post_id IS NULL
			  AND messages.id IS NULL
			  AND message_attachments.message_id IS NULL
			  AND media_assets.processing_status = $3
			  AND media_assets.created_at < $1
			ORDER BY media_assets.created_at ASC
			LIMIT $2
		)
	`, olderThan, limit, status)
	if err != nil {
		return 0, err
	}
	return commandTag.RowsAffected(), nil
}
