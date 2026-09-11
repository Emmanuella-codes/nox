package admin

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const mediaReferenceCTE = `WITH media_refs AS (
	SELECT media_asset_id, 'set' AS entity_type, id AS entity_id FROM sets
	UNION ALL SELECT media_asset_id, 'story_item', id FROM story_items
	UNION ALL SELECT media_asset_id, 'post', post_id FROM post_media_assets
	UNION ALL SELECT media_asset_id, 'message', id FROM messages WHERE media_asset_id IS NOT NULL
	UNION ALL SELECT media_asset_id, 'message_attachment', message_id FROM message_attachments
), reference_counts AS (
	SELECT media_asset_id, COUNT(*) AS reference_count FROM media_refs GROUP BY media_asset_id
) `

const mediaSelect = `m.id, m.owner_user_id, m.owner_persona_id, m.media_kind, m.storage_key,
	m.playback_url, COALESCE(m.thumbnail_url, ''), m.mime_type, m.duration_seconds, m.size_bytes,
	m.processing_status, m.processing_attempts, m.last_processing_error, m.last_processing_started_at,
	m.created_at, m.updated_at, COALESCE(mm.status, 'active'), COALESCE(mm.reason, ''),
	mm.moderated_by_user_id, mm.moderated_at, COALESCE(rc.reference_count, 0), rc.reference_count IS NULL`

func (r *pgRepository) ListAdminMedia(ctx context.Context, params ListAdminMediaParams) (*models.AdminMediaPage, error) {
	limit, offset := reportPage(params.Limit, params.Offset)
	where := []string{"TRUE"}
	args := []any{}
	add := func(value any) string {
		args = append(args, value)
		return fmt.Sprintf("$%d", len(args))
	}
	if params.Status != nil {
		where = append(where, "m.processing_status = "+add(params.Status))
	}
	if params.Kind != nil {
		where = append(where, "m.media_kind = "+add(params.Kind))
	}
	if params.Moderation != nil {
		where = append(where, "COALESCE(mm.status, 'active') = "+add(params.Moderation))
	}
	if params.OwnerID != nil {
		where = append(where, "m.owner_user_id = "+add(params.OwnerID))
	}
	if params.OlderThan != nil {
		where = append(where, "m.updated_at < "+add(params.OlderThan))
	}
	if params.Orphaned != nil {
		value := "rc.reference_count IS NULL"
		if !*params.Orphaned {
			value = "rc.reference_count IS NOT NULL"
		}
		where = append(where, value)
	}
	args = append(args, limit+1, offset)
	query := mediaReferenceCTE + `SELECT ` + mediaSelect + `
		FROM media_assets m
		LEFT JOIN media_moderation mm ON mm.media_asset_id = m.id
		LEFT JOIN reference_counts rc ON rc.media_asset_id = m.id
		WHERE ` + strings.Join(where, " AND ") + `
		ORDER BY m.created_at DESC, m.id DESC
		LIMIT $` + fmt.Sprint(len(args)-1) + ` OFFSET $` + fmt.Sprint(len(args))
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]models.AdminMediaAsset, 0, limit)
	for rows.Next() {
		item, err := scanAdminMedia(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	page := &models.AdminMediaPage{Items: items, Limit: limit, Offset: offset}
	if len(items) > limit {
		page.HasMore = true
		page.Items = items[:limit]
	}
	return page, nil
}

func (r *pgRepository) FindAdminMedia(ctx context.Context, id uuid.UUID) (*models.AdminMediaAsset, error) {
	row := r.db.QueryRow(ctx, mediaReferenceCTE+`SELECT `+mediaSelect+`
		FROM media_assets m
		LEFT JOIN media_moderation mm ON mm.media_asset_id = m.id
		LEFT JOIN reference_counts rc ON rc.media_asset_id = m.id
		WHERE m.id = $1`, id)
	item, err := scanAdminMedia(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrMediaNotFound
	}
	if err != nil {
		return nil, err
	}
	rows, err := r.db.Query(ctx, `SELECT entity_type, entity_id FROM (
		SELECT media_asset_id, 'set' AS entity_type, id AS entity_id FROM sets
		UNION ALL SELECT media_asset_id, 'story_item', id FROM story_items
		UNION ALL SELECT media_asset_id, 'post', post_id FROM post_media_assets
		UNION ALL SELECT media_asset_id, 'message', id FROM messages WHERE media_asset_id IS NOT NULL
		UNION ALL SELECT media_asset_id, 'message_attachment', message_id FROM message_attachments
	) refs WHERE media_asset_id = $1 ORDER BY entity_type, entity_id`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var reference models.AdminMediaReference
		if err := rows.Scan(&reference.EntityType, &reference.EntityID); err != nil {
			return nil, err
		}
		item.References = append(item.References, reference)
	}
	return item, rows.Err()
}

func (r *pgRepository) RetryMedia(ctx context.Context, id uuid.UUID, adminID uuid.UUID) (*models.AdminMediaAsset, error) {
	_, err := r.db.Exec(ctx, `UPDATE media_assets SET processing_status = 'pending',
		processing_attempts = processing_attempts + 1, last_processing_error = '',
		last_processing_started_at = now(), updated_at = now() WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}
	return r.FindAdminMedia(ctx, id)
}

func (r *pgRepository) CorrectMedia(ctx context.Context, params CorrectMediaParams) (*models.AdminMediaAsset, error) {
	row := r.db.QueryRow(ctx, `UPDATE media_assets SET processing_status = $2,
		playback_url = CASE WHEN $2 = 'ready' AND $3 <> '' THEN $3 ELSE playback_url END,
		thumbnail_url = CASE WHEN $4 <> '' THEN $4 ELSE thumbnail_url END,
		mime_type = CASE WHEN $5 <> '' THEN $5 ELSE mime_type END,
		duration_seconds = CASE WHEN $6 > 0 THEN $6 ELSE duration_seconds END,
		size_bytes = CASE WHEN $7 > 0 THEN $7 ELSE size_bytes END,
		last_processing_error = CASE WHEN $2 = 'failed' THEN $8 ELSE '' END,
		updated_at = now() WHERE id = $1 RETURNING id`, params.MediaAssetID, params.Status,
		params.PlaybackURL, params.ThumbnailURL, params.MimeType, params.Duration, params.SizeBytes, params.Reason)
	var id uuid.UUID
	if err := row.Scan(&id); errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrMediaNotFound
	} else if err != nil {
		return nil, err
	}
	return r.FindAdminMedia(ctx, id)
}

func (r *pgRepository) ModerateMedia(ctx context.Context, params ModerateMediaParams) (*models.AdminMediaAsset, error) {
	var exists bool
	if err := r.db.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM media_assets WHERE id = $1)`, params.MediaAssetID).Scan(&exists); err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrMediaNotFound
	}
	_, err := r.db.Exec(ctx, `INSERT INTO media_moderation (media_asset_id, status, reason, moderated_by_user_id, moderated_at)
		VALUES ($1, $2, $3, $4, now()) ON CONFLICT (media_asset_id) DO UPDATE SET status = EXCLUDED.status,
		reason = EXCLUDED.reason, moderated_by_user_id = EXCLUDED.moderated_by_user_id, moderated_at = EXCLUDED.moderated_at,
		updated_at = now()`, params.MediaAssetID, params.Status, params.Reason, params.AdminID)
	if err != nil {
		return nil, err
	}
	return r.FindAdminMedia(ctx, params.MediaAssetID)
}

func (r *pgRepository) CleanupOrphanedMedia(ctx context.Context, olderThan time.Time, limit int) (int64, error) {
	if limit <= 0 || limit > 1000 {
		limit = 100
	}
	commandTag, err := r.db.Exec(ctx, `DELETE FROM media_assets WHERE id IN (
		SELECT m.id FROM media_assets m
		LEFT JOIN sets s ON s.media_asset_id = m.id
		LEFT JOIN story_items si ON si.media_asset_id = m.id
		LEFT JOIN post_media_assets pma ON pma.media_asset_id = m.id
		LEFT JOIN messages msg ON msg.media_asset_id = m.id
		LEFT JOIN message_attachments ma ON ma.media_asset_id = m.id
		WHERE s.id IS NULL AND si.id IS NULL AND pma.post_id IS NULL AND msg.id IS NULL AND ma.message_id IS NULL
		  AND m.processing_status IN ('pending', 'failed') AND m.created_at < $1
		ORDER BY m.created_at ASC LIMIT $2)`, olderThan, limit)
	return commandTag.RowsAffected(), err
}

func scanAdminMedia(row interface{ Scan(...any) error }) (*models.AdminMediaAsset, error) {
	item := &models.AdminMediaAsset{}
	err := row.Scan(&item.ID, &item.OwnerUserID, &item.OwnerPersonaID, &item.MediaKind, &item.StorageKey,
		&item.PlaybackURL, &item.ThumbnailURL, &item.MimeType, &item.DurationSeconds, &item.SizeBytes,
		&item.ProcessingStatus, &item.ProcessingAttempts, &item.LastProcessingError, &item.LastProcessingStartedAt,
		&item.CreatedAt, &item.UpdatedAt, &item.ModerationStatus, &item.ModerationReason, &item.ModeratedByUserID,
		&item.ModeratedAt, &item.ReferenceCount, &item.IsOrphaned)
	return item, err
}
