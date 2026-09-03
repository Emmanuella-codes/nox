package story

import (
	"context"
	"errors"

	"github.com/emmanuella-codes/nox/models"
	storydtos "github.com/emmanuella-codes/nox/story/dtos"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// AddStoryItem inserts one direct owner-managed item into a story.
func (r *pgRepository) AddStoryItem(ctx context.Context, storyID uuid.UUID, contributorUserID uuid.UUID, durationSeconds int, anonymousLabel string, dto storydtos.AddStoryItemDTO) (*models.StoryItem, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var currentTotal int
	if err := tx.QueryRow(ctx, `SELECT total_duration_seconds FROM stories WHERE id = $1 FOR UPDATE`, storyID).Scan(&currentTotal); err != nil {
		return nil, mapStoryError(err)
	}
	if currentTotal+durationSeconds > maxStoryDurationSeconds {
		return nil, ErrStoryDurationLimitExceeded
	}
	acceptsAdditions, err := storyAcceptsAdditions(ctx, tx, storyID)
	if err != nil {
		return nil, err
	}
	if !acceptsAdditions {
		return nil, ErrStoryClosedForAdditions
	}

	var position int
	if err := tx.QueryRow(ctx, `SELECT COALESCE(MAX(position), 0) + 1 FROM story_items WHERE story_id = $1`, storyID).Scan(&position); err != nil {
		return nil, err
	}
	item, err := insertStoryItemInTx(ctx, tx, storyID, dto.MediaAssetID, contributorUserID, dto.ContributorPersonaID, dto.PostingMode, anonymousLabel, durationSeconds, position)
	if err != nil {
		return nil, err
	}
	if err := recalculateStoryState(ctx, tx, storyID); err != nil {
		return nil, err
	}
	return item, tx.Commit(ctx)
}

// ReorderStoryItem moves one story item to a new position within the same story.
func (r *pgRepository) ReorderStoryItem(ctx context.Context, storyID uuid.UUID, itemID uuid.UUID, position int) (*models.StoryItem, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var currentPosition int
	if err := tx.QueryRow(ctx, `SELECT position FROM story_items WHERE story_id = $1 AND id = $2 FOR UPDATE`, storyID, itemID).Scan(&currentPosition); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrStoryItemNotFound
		}
		return nil, err
	}
	var maxPosition int
	if err := tx.QueryRow(ctx, `SELECT COUNT(*) FROM story_items WHERE story_id = $1`, storyID).Scan(&maxPosition); err != nil {
		return nil, err
	}
	position = boundedPosition(position, maxPosition)
	if position == currentPosition {
		item, err := r.storyItemInTx(ctx, tx, storyID, itemID)
		if err != nil {
			return nil, err
		}
		return item, tx.Commit(ctx)
	}
	if _, err := tx.Exec(ctx, `UPDATE story_items SET position = 0 WHERE story_id = $1 AND id = $2`, storyID, itemID); err != nil {
		return nil, err
	}
	if err := shiftPositions(ctx, tx, "story_items", "story_id", storyID, position, currentPosition); err != nil {
		return nil, err
	}
	row := tx.QueryRow(ctx, `
		UPDATE story_items
		SET position = $3
		WHERE story_id = $1 AND id = $2
		RETURNING id, story_id, media_asset_id, contributor_user_id, contributor_persona_id,
		          posting_mode, COALESCE(anonymous_label, ''), duration_seconds, position, expires_at, created_at
	`, storyID, itemID, position)
	item, err := scanStoryItem(row)
	if err != nil {
		return nil, err
	}
	return item, tx.Commit(ctx)
}

// FindStoryItems lists all items for one story in position order.
func (r *pgRepository) FindStoryItems(ctx context.Context, storyID uuid.UUID) ([]*models.StoryItem, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, story_id, media_asset_id, contributor_user_id, contributor_persona_id,
		       posting_mode, COALESCE(anonymous_label, ''), duration_seconds, position, expires_at, created_at
		FROM story_items
		WHERE story_id = $1 AND expires_at > now()
		ORDER BY position ASC
	`, storyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanStoryItems(rows)
}

// Lists all items for one story regardless of expiry.
func (r *pgRepository) FindStoryItemsAny(ctx context.Context, storyID uuid.UUID) ([]*models.StoryItem, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, story_id, media_asset_id, contributor_user_id, contributor_persona_id,
		       posting_mode, COALESCE(anonymous_label, ''), duration_seconds, position, expires_at, created_at
		FROM story_items
		WHERE story_id = $1
		ORDER BY position ASC
	`, storyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanStoryItems(rows)
}

// FindStoryItemByID loads one story item scoped to one story.
func (r *pgRepository) FindStoryItemByID(ctx context.Context, storyID uuid.UUID, itemID uuid.UUID) (*models.StoryItem, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, story_id, media_asset_id, contributor_user_id, contributor_persona_id,
		       posting_mode, COALESCE(anonymous_label, ''), duration_seconds, position, expires_at, created_at
		FROM story_items
		WHERE story_id = $1 AND id = $2 AND expires_at > now()
	`, storyID, itemID)
	item, err := scanStoryItem(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrStoryItemNotFound
		}
		return nil, err
	}
	return item, nil
}

// DeleteStoryItem removes one story item and updates the story duration total.
func (r *pgRepository) DeleteStoryItem(ctx context.Context, storyID uuid.UUID, itemID uuid.UUID) (*models.StoryItem, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	row := tx.QueryRow(ctx, `
		DELETE FROM story_items
		WHERE story_id = $1 AND id = $2
		RETURNING id, story_id, media_asset_id, contributor_user_id, contributor_persona_id,
		          posting_mode, COALESCE(anonymous_label, ''), duration_seconds, position, expires_at, created_at
	`, storyID, itemID)
	item, err := scanStoryItem(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrStoryItemNotFound
		}
		return nil, err
	}
	if err := recalculateStoryState(ctx, tx, storyID); err != nil {
		return nil, err
	}
	return item, tx.Commit(ctx)
}

// storyItemInTx loads one story item inside an existing transaction.
func (r *pgRepository) storyItemInTx(ctx context.Context, tx pgx.Tx, storyID uuid.UUID, itemID uuid.UUID) (*models.StoryItem, error) {
	row := tx.QueryRow(ctx, `
		SELECT id, story_id, media_asset_id, contributor_user_id, contributor_persona_id,
		       posting_mode, COALESCE(anonymous_label, ''), duration_seconds, position, expires_at, created_at
		FROM story_items
		WHERE story_id = $1 AND id = $2
	`, storyID, itemID)
	return scanStoryItem(row)
}

// Inserts one story item row with a fresh twenty-four hour expiry.
func insertStoryItemInTx(ctx context.Context, tx pgx.Tx, storyID uuid.UUID, mediaAssetID uuid.UUID, contributorUserID uuid.UUID, contributorPersonaID uuid.UUID, postingMode models.PostingMode, anonymousLabel string, durationSeconds int, position int) (*models.StoryItem, error) {
	row := tx.QueryRow(ctx, `
		INSERT INTO story_items (
			story_id, media_asset_id, contributor_user_id, contributor_persona_id,
			posting_mode, anonymous_label, duration_seconds, position, expires_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, now() + interval '24 hours')
		RETURNING id, story_id, media_asset_id, contributor_user_id, contributor_persona_id,
		          posting_mode, COALESCE(anonymous_label, ''), duration_seconds, position, expires_at, created_at
	`, storyID, mediaAssetID, contributorUserID, contributorPersonaID, postingMode, emptyToNil(anonymousLabel), durationSeconds, position)
	item, err := scanStoryItem(row)
	if err != nil {
		if isUniqueViolation(err, "story_items_media_asset_id_key") {
			return nil, ErrStoryMediaInUse
		}
		return nil, err
	}
	return item, nil
}
