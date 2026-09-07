package admin

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/emmanuella-codes/nox/models"
	"github.com/emmanuella-codes/nox/repositories/hashtag"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *pgRepository) ListAdminHashtags(ctx context.Context, params ListAdminHashtagsParams) (*models.AdminHashtagPage, error) {
	limit, offset := reportPage(params.Limit, params.Offset)
	where := []string{"h.post_count > 0"}
	args := []any{}
	if params.Query != "" {
		args = append(args, "%"+strings.ToLower(strings.TrimSpace(params.Query))+"%")
		where = append(where, "h.tag ILIKE $"+fmt.Sprint(len(args)))
	}
	if params.Status != nil {
		args = append(args, params.Status)
		where = append(where, "COALESCE(hm.status, 'active') = $"+fmt.Sprint(len(args)))
	}
	args = append(args, limit+1, offset)
	rows, err := r.db.Query(ctx, `SELECT h.id, h.tag, h.post_count, h.created_at,
		COALESCE(hm.status, 'active'), COALESCE(hm.reason, ''), hm.moderated_by_user_id, hm.moderated_at,
		hs.id IS NOT NULL, COALESCE(hs.reason, ''), hs.expires_at
		FROM hashtags h
		LEFT JOIN hashtag_moderation hm ON hm.hashtag_id = h.id
		LEFT JOIN hashtag_suppressions hs ON hs.hashtag_id = h.id AND (hs.expires_at IS NULL OR hs.expires_at > now())
		WHERE `+strings.Join(where, " AND ")+`
		ORDER BY h.post_count DESC, h.tag ASC
		LIMIT $`+fmt.Sprint(len(args)-1)+` OFFSET $`+fmt.Sprint(len(args)), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]models.AdminHashtag, 0, limit)
	for rows.Next() {
		item, err := scanAdminHashtag(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	page := &models.AdminHashtagPage{Items: items, Limit: limit, Offset: offset}
	if len(items) > limit {
		page.HasMore = true
		page.Items = items[:limit]
	}
	return page, nil
}

func (r *pgRepository) ModerateHashtag(ctx context.Context, params ModerateHashtagParams) (*models.AdminHashtag, error) {
	tag := hashtag.NormalizeTag(params.Tag)
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	var id uuid.UUID
	if err := tx.QueryRow(ctx, `SELECT id FROM hashtags WHERE tag = $1`, tag).Scan(&id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrHashtagNotFound
		}
		return nil, err
	}
	_, err = tx.Exec(ctx, `INSERT INTO hashtag_moderation (hashtag_id, status, reason, moderated_by_user_id, moderated_at)
		VALUES ($1, $2, $3, $4, now()) ON CONFLICT (hashtag_id) DO UPDATE SET status = EXCLUDED.status,
		reason = EXCLUDED.reason, moderated_by_user_id = EXCLUDED.moderated_by_user_id, moderated_at = EXCLUDED.moderated_at,
		updated_at = now()`, id, params.Status, params.Reason, params.AdminID)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return r.findAdminHashtag(ctx, id)
}

func (r *pgRepository) SuppressHashtag(ctx context.Context, params SuppressHashtagParams) (*models.AdminHashtag, error) {
	tag := hashtag.NormalizeTag(params.Tag)
	var id uuid.UUID
	if err := r.db.QueryRow(ctx, `SELECT id FROM hashtags WHERE tag = $1`, tag).Scan(&id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrHashtagNotFound
		}
		return nil, err
	}
	_, err := r.db.Exec(ctx, `INSERT INTO hashtag_suppressions (hashtag_id, reason, created_by_user_id, expires_at)
		VALUES ($1, $2, $3, $4) ON CONFLICT (hashtag_id) DO UPDATE SET reason = EXCLUDED.reason,
		created_by_user_id = EXCLUDED.created_by_user_id, created_at = now(), expires_at = EXCLUDED.expires_at`, id, params.Reason, params.AdminID, params.ExpiresAt)
	if err != nil {
		return nil, err
	}
	return r.findAdminHashtag(ctx, id)
}

func (r *pgRepository) UnsuppressHashtag(ctx context.Context, tag string) error {
	result, err := r.db.Exec(ctx, `DELETE FROM hashtag_suppressions hs USING hashtags h WHERE hs.hashtag_id = h.id AND h.tag = $1`, hashtag.NormalizeTag(tag))
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrHashtagNotFound
	}
	return nil
}

func (r *pgRepository) findAdminHashtag(ctx context.Context, id uuid.UUID) (*models.AdminHashtag, error) {
	row := r.db.QueryRow(ctx, `SELECT h.id, h.tag, h.post_count, h.created_at,
		COALESCE(hm.status, 'active'), COALESCE(hm.reason, ''), hm.moderated_by_user_id, hm.moderated_at,
		hs.id IS NOT NULL, COALESCE(hs.reason, ''), hs.expires_at
		FROM hashtags h LEFT JOIN hashtag_moderation hm ON hm.hashtag_id = h.id
		LEFT JOIN hashtag_suppressions hs ON hs.hashtag_id = h.id AND (hs.expires_at IS NULL OR hs.expires_at > now())
		WHERE h.id = $1`, id)
	return scanAdminHashtag(row)
}

func scanAdminHashtag(row interface{ Scan(...any) error }) (*models.AdminHashtag, error) {
	item := &models.AdminHashtag{}
	err := row.Scan(&item.ID, &item.Tag, &item.PostCount, &item.CreatedAt, &item.ModerationStatus,
		&item.ModerationReason, &item.ModeratedByUserID, &item.ModeratedAt, &item.IsSuppressed,
		&item.SuppressionReason, &item.SuppressionExpires)
	return item, err
}

func (r *pgRepository) ListSearchSuppressions(ctx context.Context, params ListSearchSuppressionsParams) (*models.SearchSuppressionPage, error) {
	limit, offset := reportPage(params.Limit, params.Offset)
	query := strings.TrimSpace(params.Query)
	rows, err := r.db.Query(ctx, `SELECT id, normalized_query, reason, created_by_user_id, created_at, expires_at
		FROM search_suppressions
		WHERE ($1 = '' OR normalized_query ILIKE '%' || lower($1) || '%')
		ORDER BY created_at DESC, id DESC LIMIT $2 OFFSET $3`, query, limit+1, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]models.SearchSuppression, 0, limit)
	for rows.Next() {
		var item models.SearchSuppression
		if err := rows.Scan(&item.ID, &item.NormalizedQuery, &item.Reason, &item.CreatedBy, &item.CreatedAt, &item.ExpiresAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	page := &models.SearchSuppressionPage{Items: items, Limit: limit, Offset: offset}
	if len(items) > limit {
		page.HasMore = true
		page.Items = items[:limit]
	}
	return page, nil
}

func (r *pgRepository) CreateSearchSuppression(ctx context.Context, params CreateSearchSuppressionParams) (*models.SearchSuppression, error) {
	query := normalizeSearchSuppression(params.Query)
	var item models.SearchSuppression
	err := r.db.QueryRow(ctx, `INSERT INTO search_suppressions (normalized_query, reason, created_by_user_id, expires_at)
		VALUES ($1, $2, $3, $4) ON CONFLICT (normalized_query) DO UPDATE SET reason = EXCLUDED.reason,
		created_by_user_id = EXCLUDED.created_by_user_id, created_at = now(), expires_at = EXCLUDED.expires_at
		RETURNING id, normalized_query, reason, created_by_user_id, created_at, expires_at`, query, params.Reason, params.AdminID, params.ExpiresAt).
		Scan(&item.ID, &item.NormalizedQuery, &item.Reason, &item.CreatedBy, &item.CreatedAt, &item.ExpiresAt)
	return &item, err
}

func (r *pgRepository) DeleteSearchSuppression(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.Exec(ctx, `DELETE FROM search_suppressions WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrSearchSuppressionNotFound
	}
	return nil
}

func normalizeSearchSuppression(value string) string {
	return strings.ToLower(strings.TrimSpace(strings.TrimLeft(value, "#")))
}
