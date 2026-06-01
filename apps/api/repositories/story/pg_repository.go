package story

import (
	"context"
	"errors"

	"github.com/emmanuella-codes/nox/models"
	storydtos "github.com/emmanuella-codes/nox/story/dtos"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const maxStoryDurationSeconds = 300

type pgRepository struct {
	db *pgxpool.Pool
}

func newPgRepository(db *pgxpool.Pool) *pgRepository {
	return &pgRepository{db: db}
}

func (r *pgRepository) CreateStory(ctx context.Context, ownerUserID uuid.UUID, dto storydtos.CreateStoryDTO) (*models.Story, error) {
	row := r.db.QueryRow(ctx, `
		INSERT INTO stories (event_id, owner_user_id, owner_persona_id, title, contribution_mode)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, event_id, owner_user_id, owner_persona_id, title, contribution_mode,
		          total_duration_seconds, created_at, updated_at
	`, dto.EventID, ownerUserID, dto.OwnerPersonaID, dto.Title, dto.ContributionMode)
	return scanStory(row)
}

func (r *pgRepository) FindStoryByID(ctx context.Context, storyID uuid.UUID) (*models.Story, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, event_id, owner_user_id, owner_persona_id, title, contribution_mode,
		       total_duration_seconds, created_at, updated_at
		FROM stories
		WHERE id = $1
	`, storyID)
	story, err := scanStory(row)
	if err != nil {
		return nil, mapStoryError(err)
	}
	return story, nil
}

func (r *pgRepository) FindStoriesByEventID(ctx context.Context, eventID uuid.UUID, limit int, offset int) ([]*models.Story, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, event_id, owner_user_id, owner_persona_id, title, contribution_mode,
		       total_duration_seconds, created_at, updated_at
		FROM stories
		WHERE event_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`, eventID, normalizeLimit(limit), normalizeOffset(offset))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanStories(rows)
}

func (r *pgRepository) FindStoriesByOwnerPersonaID(ctx context.Context, personaID uuid.UUID, limit int, offset int) ([]*models.Story, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, event_id, owner_user_id, owner_persona_id, title, contribution_mode,
		       total_duration_seconds, created_at, updated_at
		FROM stories
		WHERE owner_persona_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`, personaID, normalizeLimit(limit), normalizeOffset(offset))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanStories(rows)
}

func (r *pgRepository) DeleteStory(ctx context.Context, storyID uuid.UUID) error {
	commandTag, err := r.db.Exec(ctx, `DELETE FROM stories WHERE id = $1`, storyID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return ErrStoryNotFound
	}
	return nil
}

func (r *pgRepository) AddStoryItem(ctx context.Context, storyID uuid.UUID, contributorUserID uuid.UUID, durationSeconds int, anonymousLabel string, dto storydtos.AddStoryItemDTO) (*models.StoryItem, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var currentTotal int
	if err := tx.QueryRow(ctx, `
		SELECT total_duration_seconds
		FROM stories
		WHERE id = $1
		FOR UPDATE
	`, storyID).Scan(&currentTotal); err != nil {
		return nil, mapStoryError(err)
	}
	if currentTotal+durationSeconds > maxStoryDurationSeconds {
		return nil, ErrStoryDurationLimitExceeded
	}

	var position int
	if err := tx.QueryRow(ctx, `
		SELECT COALESCE(MAX(position), 0) + 1
		FROM story_items
		WHERE story_id = $1
	`, storyID).Scan(&position); err != nil {
		return nil, err
	}

	row := tx.QueryRow(ctx, `
		INSERT INTO story_items (
			story_id, media_asset_id, contributor_user_id, contributor_persona_id,
			posting_mode, anonymous_label, duration_seconds, position
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, story_id, media_asset_id, contributor_user_id, contributor_persona_id,
		          posting_mode, COALESCE(anonymous_label, ''), duration_seconds, position, created_at
	`, storyID, dto.MediaAssetID, contributorUserID, dto.ContributorPersonaID, dto.PostingMode, emptyToNil(anonymousLabel), durationSeconds, position)

	item, err := scanStoryItem(row)
	if err != nil {
		return nil, err
	}
	if _, err := tx.Exec(ctx, `
		UPDATE stories
		SET total_duration_seconds = total_duration_seconds + $2,
		    updated_at = now()
		WHERE id = $1
	`, storyID, durationSeconds); err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return item, nil
}

func (r *pgRepository) FindStoryItems(ctx context.Context, storyID uuid.UUID) ([]*models.StoryItem, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, story_id, media_asset_id, contributor_user_id, contributor_persona_id,
		       posting_mode, COALESCE(anonymous_label, ''), duration_seconds, position, created_at
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
		          posting_mode, COALESCE(anonymous_label, ''), duration_seconds, position, created_at
	`, storyID, itemID)
	item, err := scanStoryItem(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrStoryItemNotFound
		}
		return nil, err
	}
	if _, err := tx.Exec(ctx, `
		UPDATE stories
		SET total_duration_seconds = GREATEST(total_duration_seconds - $2, 0),
		    updated_at = now()
		WHERE id = $1
	`, storyID, item.DurationSeconds); err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return item, nil
}

func (r *pgRepository) AddEventHighlightStory(ctx context.Context, eventID uuid.UUID, dto storydtos.AddEventHighlightStoryDTO) (*models.EventHighlightStory, error) {
	var position int
	if err := r.db.QueryRow(ctx, `
		SELECT COALESCE(MAX(position), 0) + 1
		FROM event_highlight_stories
		WHERE event_id = $1
	`, eventID).Scan(&position); err != nil {
		return nil, err
	}
	row := r.db.QueryRow(ctx, `
		INSERT INTO event_highlight_stories (event_id, story_id, added_by_persona_id, position)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (event_id, story_id) DO UPDATE SET story_id = EXCLUDED.story_id
		RETURNING id, event_id, story_id, added_by_persona_id, position, created_at
	`, eventID, dto.StoryID, dto.AddedByPersonaID, position)
	return scanEventHighlightStory(row)
}

func (r *pgRepository) FindEventHighlightStories(ctx context.Context, eventID uuid.UUID) ([]*models.EventHighlightStory, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, event_id, story_id, added_by_persona_id, position, created_at
		FROM event_highlight_stories
		WHERE event_id = $1
		ORDER BY position ASC
	`, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanEventHighlightStories(rows)
}

func (r *pgRepository) RemoveEventHighlightStory(ctx context.Context, eventID uuid.UUID, storyID uuid.UUID) error {
	commandTag, err := r.db.Exec(ctx, `
		DELETE FROM event_highlight_stories
		WHERE event_id = $1 AND story_id = $2
	`, eventID, storyID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return ErrEventHighlightNotFound
	}
	return nil
}
