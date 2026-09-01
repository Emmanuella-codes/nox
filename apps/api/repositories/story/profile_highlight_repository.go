package story

import (
	"context"

	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
)

// AddProfileStoryHighlight adds one story to one persona's profile highlight list.
func (r *pgRepository) AddProfileStoryHighlight(ctx context.Context, ownerPersonaID uuid.UUID, storyID uuid.UUID) (*models.ProfileStoryHighlight, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `SELECT id FROM profile_story_highlights WHERE owner_persona_id = $1 FOR UPDATE`, ownerPersonaID); err != nil {
		return nil, err
	}
	var position int
	if err := tx.QueryRow(ctx, `SELECT COALESCE(MAX(position), 0) + 1 FROM profile_story_highlights WHERE owner_persona_id = $1`, ownerPersonaID).Scan(&position); err != nil {
		return nil, err
	}
	row := tx.QueryRow(ctx, `
		INSERT INTO profile_story_highlights (owner_persona_id, story_id, position)
		VALUES ($1, $2, $3)
		ON CONFLICT (owner_persona_id, story_id) DO UPDATE
		SET story_id = EXCLUDED.story_id
		RETURNING id, owner_persona_id, story_id, position, created_at
	`, ownerPersonaID, storyID, position)
	highlight, err := scanProfileStoryHighlight(row)
	if err != nil {
		return nil, err
	}
	return highlight, tx.Commit(ctx)
}

// FindProfileStoryHighlights lists one persona's highlighted stories in position order.
func (r *pgRepository) FindProfileStoryHighlights(ctx context.Context, ownerPersonaID uuid.UUID) ([]*models.ProfileStoryHighlight, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, owner_persona_id, story_id, position, created_at
		FROM profile_story_highlights
		WHERE owner_persona_id = $1
		ORDER BY position ASC
	`, ownerPersonaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanProfileStoryHighlights(rows)
}

// RemoveProfileStoryHighlight removes one story from one persona's highlight list.
func (r *pgRepository) RemoveProfileStoryHighlight(ctx context.Context, ownerPersonaID uuid.UUID, storyID uuid.UUID) error {
	tag, err := r.db.Exec(ctx, `
		DELETE FROM profile_story_highlights
		WHERE owner_persona_id = $1 AND story_id = $2
	`, ownerPersonaID, storyID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrStoryNotFound
	}
	return nil
}

// HasProfileStoryHighlight reports whether one story is highlighted on one persona profile.
func (r *pgRepository) HasProfileStoryHighlight(ctx context.Context, ownerPersonaID uuid.UUID, storyID uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM profile_story_highlights
			WHERE owner_persona_id = $1 AND story_id = $2
		)
	`, ownerPersonaID, storyID).Scan(&exists)
	return exists, err
}

// scanProfileStoryHighlight scans one profile highlight row into the model shape.
func scanProfileStoryHighlight(scanner storyScanner) (*models.ProfileStoryHighlight, error) {
	var highlight models.ProfileStoryHighlight
	if err := scanner.Scan(&highlight.ID, &highlight.OwnerPersonaID, &highlight.StoryID, &highlight.Position, &highlight.CreatedAt); err != nil {
		return nil, err
	}
	return &highlight, nil
}

// scanProfileStoryHighlights scans many profile highlight rows into model values.
func scanProfileStoryHighlights(rows storyRows) ([]*models.ProfileStoryHighlight, error) {
	var highlights []*models.ProfileStoryHighlight
	for rows.Next() {
		highlight, err := scanProfileStoryHighlight(rows)
		if err != nil {
			return nil, err
		}
		highlights = append(highlights, highlight)
	}
	return highlights, rows.Err()
}
