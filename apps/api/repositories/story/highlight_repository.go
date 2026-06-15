package story

import (
	"context"
	"errors"

	"github.com/emmanuella-codes/nox/models"
	storydtos "github.com/emmanuella-codes/nox/story/dtos"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *pgRepository) AddEventHighlightStory(ctx context.Context, eventID uuid.UUID, dto storydtos.AddEventHighlightStoryDTO) (*models.EventHighlightStory, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `SELECT id FROM event_highlight_stories WHERE event_id = $1 FOR UPDATE`, eventID); err != nil {
		return nil, err
	}
	var position int
	if err := tx.QueryRow(ctx, `SELECT COALESCE(MAX(position), 0) + 1 FROM event_highlight_stories WHERE event_id = $1`, eventID).Scan(&position); err != nil {
		return nil, err
	}
	row := tx.QueryRow(ctx, `
		INSERT INTO event_highlight_stories (event_id, story_id, added_by_persona_id, position)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (event_id, story_id) DO UPDATE
		SET added_by_persona_id = EXCLUDED.added_by_persona_id
		RETURNING id, event_id, story_id, added_by_persona_id, position, created_at
	`, eventID, dto.StoryID, dto.AddedByPersonaID, position)
	highlight, err := scanEventHighlightStory(row)
	if err != nil {
		return nil, err
	}
	return highlight, tx.Commit(ctx)
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

func (r *pgRepository) ReorderEventHighlightStory(ctx context.Context, eventID uuid.UUID, storyID uuid.UUID, position int) (*models.EventHighlightStory, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var currentPosition int
	if err := tx.QueryRow(ctx, `SELECT position FROM event_highlight_stories WHERE event_id = $1 AND story_id = $2 FOR UPDATE`, eventID, storyID).Scan(&currentPosition); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrEventHighlightNotFound
		}
		return nil, err
	}
	var maxPosition int
	if err := tx.QueryRow(ctx, `SELECT COUNT(*) FROM event_highlight_stories WHERE event_id = $1`, eventID).Scan(&maxPosition); err != nil {
		return nil, err
	}
	position = boundedPosition(position, maxPosition)
	if _, err := tx.Exec(ctx, `UPDATE event_highlight_stories SET position = 0 WHERE event_id = $1 AND story_id = $2`, eventID, storyID); err != nil {
		return nil, err
	}
	if err := shiftPositions(ctx, tx, "event_highlight_stories", "event_id", eventID, position, currentPosition); err != nil {
		return nil, err
	}
	row := tx.QueryRow(ctx, `
		UPDATE event_highlight_stories
		SET position = $3
		WHERE event_id = $1 AND story_id = $2
		RETURNING id, event_id, story_id, added_by_persona_id, position, created_at
	`, eventID, storyID, position)
	highlight, err := scanEventHighlightStory(row)
	if err != nil {
		return nil, err
	}
	return highlight, tx.Commit(ctx)
}

func (r *pgRepository) RemoveEventHighlightStory(ctx context.Context, eventID uuid.UUID, storyID uuid.UUID) error {
	commandTag, err := r.db.Exec(ctx, `DELETE FROM event_highlight_stories WHERE event_id = $1 AND story_id = $2`, eventID, storyID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return ErrEventHighlightNotFound
	}
	return nil
}
