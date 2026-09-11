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

func (r *pgRepository) ListFeaturedSets(ctx context.Context, params ListFeaturedSetsParams) ([]models.AdminFeaturedSet, error) {
	limit, offset := reportPage(params.Limit, params.Offset)
	rows, err := r.db.Query(ctx, `SELECT sf.set_id, sf.position, sf.featured_by_user_id, sf.created_at, sf.expires_at
		FROM set_features sf JOIN sets s ON s.id = sf.set_id
		WHERE s.moderation_status = 'active' AND (sf.expires_at IS NULL OR sf.expires_at > now())
		ORDER BY sf.position ASC, sf.created_at DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]models.AdminFeaturedSet, 0)
	for rows.Next() {
		var item models.AdminFeaturedSet
		if err := rows.Scan(&item.SetID, &item.Position, &item.FeaturedByUserID, &item.CreatedAt, &item.ExpiresAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *pgRepository) FeatureSet(ctx context.Context, setID uuid.UUID, adminID uuid.UUID, expiresAt *time.Time) (*models.AdminFeaturedSet, error) {
	var exists bool
	if err := r.db.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM sets WHERE id = $1 AND moderation_status = 'active')`, setID).Scan(&exists); err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrModerationEntityNotFound
	}
	var position int
	if err := r.db.QueryRow(ctx, `SELECT COALESCE(MAX(position), 0) + 1 FROM set_features`).Scan(&position); err != nil {
		return nil, err
	}
	row := r.db.QueryRow(ctx, `INSERT INTO set_features (set_id, position, featured_by_user_id, expires_at)
		VALUES ($1, $2, $3, $4) ON CONFLICT (set_id) DO UPDATE SET featured_by_user_id = EXCLUDED.featured_by_user_id,
		expires_at = EXCLUDED.expires_at RETURNING set_id, position, featured_by_user_id, created_at, expires_at`, setID, position, adminID, expiresAt)
	var item models.AdminFeaturedSet
	err := row.Scan(&item.SetID, &item.Position, &item.FeaturedByUserID, &item.CreatedAt, &item.ExpiresAt)
	return &item, err
}

func (r *pgRepository) UnfeatureSet(ctx context.Context, setID uuid.UUID) error {
	result, err := r.db.Exec(ctx, `DELETE FROM set_features WHERE set_id = $1`, setID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrModerationEntityNotFound
	}
	return nil
}

func (r *pgRepository) ListStoryContributions(ctx context.Context, params ListStoryContributionsParams) ([]models.AdminStoryContribution, error) {
	limit, offset := reportPage(params.Limit, params.Offset)
	rows, err := r.db.Query(ctx, `SELECT id, story_id, media_asset_id, contributor_user_id, contributor_persona_id,
		status, reviewed_by_persona_id, story_item_id, created_at, reviewed_at
		FROM story_contribution_requests WHERE ($1::text IS NULL OR status = $1)
		ORDER BY created_at DESC, id DESC LIMIT $2 OFFSET $3`, params.Status, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]models.AdminStoryContribution, 0)
	for rows.Next() {
		var item models.AdminStoryContribution
		if err := rows.Scan(&item.ID, &item.StoryID, &item.MediaAssetID, &item.ContributorUserID, &item.ContributorPersonaID,
			&item.Status, &item.ReviewedByPersonaID, &item.StoryItemID, &item.CreatedAt, &item.ReviewedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *pgRepository) ReviewStoryContribution(ctx context.Context, params ReviewStoryContributionParams) (*models.AdminStoryContribution, error) {
	if params.Action != "reject" && params.Action != "remove" {
		return nil, ErrModerationEntityNotFound
	}
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	var item models.AdminStoryContribution
	err = tx.QueryRow(ctx, `SELECT id, story_id, media_asset_id, contributor_user_id, contributor_persona_id,
		status, reviewed_by_persona_id, story_item_id, created_at, reviewed_at
		FROM story_contribution_requests WHERE id = $1 FOR UPDATE`, params.RequestID).Scan(&item.ID, &item.StoryID, &item.MediaAssetID,
		&item.ContributorUserID, &item.ContributorPersonaID, &item.Status, &item.ReviewedByPersonaID, &item.StoryItemID, &item.CreatedAt, &item.ReviewedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrModerationEntityNotFound
	}
	if err != nil {
		return nil, err
	}
	if params.Action == "reject" {
		if item.Status != models.PendingStoryContributionRequestStatus {
			return nil, ErrModerationEntityNotFound
		}
		if err := tx.QueryRow(ctx, `UPDATE story_contribution_requests SET status = 'rejected', reviewed_at = now()
			WHERE id = $1 RETURNING status, reviewed_at`, item.ID).Scan(&item.Status, &item.ReviewedAt); err != nil {
			return nil, err
		}
	} else {
		if item.StoryItemID == nil {
			return nil, ErrModerationEntityNotFound
		}
		if _, err := tx.Exec(ctx, `UPDATE story_items SET moderation_status = 'removed', moderation_reason = $3,
			moderated_at = now(), moderated_by_user_id = $2 WHERE id = $1`, item.StoryItemID, params.AdminID, params.Reason); err != nil {
			return nil, err
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *pgRepository) ListAdminHighlights(ctx context.Context, params ListAdminHighlightsParams) ([]models.AdminHighlight, error) {
	limit, offset := reportPage(params.Limit, params.Offset)
	where := "TRUE"
	args := []any{}
	if params.HighlightType != "" {
		where = "highlight_type = $1"
		args = append(args, params.HighlightType)
	}
	args = append(args, limit, offset)
	rows, err := r.db.Query(ctx, `SELECT id, highlight_type, event_id, owner_persona_id, story_id,
		added_by_persona_id, position, created_at FROM (
		SELECT id, 'event' AS highlight_type, event_id, NULL::uuid AS owner_persona_id, story_id,
			added_by_persona_id, position, created_at FROM event_highlight_stories
		UNION ALL
		SELECT id, 'profile', NULL::uuid, owner_persona_id, story_id, NULL::uuid, position, created_at FROM profile_story_highlights
	) highlights WHERE `+where+` ORDER BY created_at DESC, id DESC LIMIT $`+fmt.Sprint(len(args)-1)+` OFFSET $`+fmt.Sprint(len(args)), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]models.AdminHighlight, 0)
	for rows.Next() {
		var item models.AdminHighlight
		if err := rows.Scan(&item.ID, &item.HighlightType, &item.EventID, &item.OwnerPersonaID, &item.StoryID, &item.AddedByPersonaID, &item.Position, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *pgRepository) RemoveAdminHighlight(ctx context.Context, highlightType string, highlightID uuid.UUID) error {
	var table string
	switch strings.ToLower(strings.TrimSpace(highlightType)) {
	case "event":
		table = "event_highlight_stories"
	case "profile":
		table = "profile_story_highlights"
	default:
		return ErrModerationEntityNotFound
	}
	result, err := r.db.Exec(ctx, `DELETE FROM `+table+` WHERE id = $1`, highlightID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrModerationEntityNotFound
	}
	return nil
}
