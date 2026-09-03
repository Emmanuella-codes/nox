package story

import (
	"context"

	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// UpsertStoryItemView records one viewer's view for one story item.
func (r *pgRepository) UpsertStoryItemView(ctx context.Context, storyID uuid.UUID, itemID uuid.UUID, viewerUserID uuid.UUID, viewerPersonaID uuid.UUID) (*models.StoryItemView, bool, error) {
	row := r.db.QueryRow(ctx, `
		INSERT INTO story_item_views (story_id, story_item_id, viewer_user_id, viewer_persona_id)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (story_item_id, viewer_persona_id) DO UPDATE
		SET viewer_user_id = EXCLUDED.viewer_user_id
		RETURNING story_id, story_item_id, viewer_user_id, viewer_persona_id, created_at,
		          (xmax = 0) AS inserted
	`, storyID, itemID, viewerUserID, viewerPersonaID)
	var view models.StoryItemView
	var inserted bool
	if err := row.Scan(&view.StoryID, &view.StoryItemID, &view.ViewerUserID, &view.ViewerPersonaID, &view.CreatedAt, &inserted); err != nil {
		return nil, false, err
	}
	return &view, inserted, nil
}

// FindStoryItemViewCounts returns view counts for all items in one story.
func (r *pgRepository) FindStoryItemViewCounts(ctx context.Context, storyID uuid.UUID) (map[uuid.UUID]int, error) {
	rows, err := r.db.Query(ctx, `
		SELECT story_item_id, COUNT(*)::INT
		FROM story_item_views
		WHERE story_id = $1
		GROUP BY story_item_id
	`, storyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	counts := make(map[uuid.UUID]int)
	for rows.Next() {
		var itemID uuid.UUID
		var count int
		if err := rows.Scan(&itemID, &count); err != nil {
			return nil, err
		}
		counts[itemID] = count
	}
	return counts, rows.Err()
}

// FindStoryItemViewerPersonaIDs lists viewer personas for one story item in reverse chronological order.
func (r *pgRepository) FindStoryItemViewerPersonaIDs(ctx context.Context, storyID uuid.UUID, itemID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := r.db.Query(ctx, `
		SELECT viewer_persona_id
		FROM story_item_views
		WHERE story_id = $1 AND story_item_id = $2
		ORDER BY created_at DESC, viewer_persona_id DESC
	`, storyID, itemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanStoryUUIDs(rows)
}

// FindViewedStoryItemIDs lists the story items already viewed by one persona.
func (r *pgRepository) FindViewedStoryItemIDs(ctx context.Context, storyID uuid.UUID, viewerPersonaID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := r.db.Query(ctx, `
		SELECT story_item_id
		FROM story_item_views
		WHERE story_id = $1 AND viewer_persona_id = $2
	`, storyID, viewerPersonaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanStoryUUIDs(rows)
}

// UpsertStoryItemReaction stores one active reaction for one viewer on one story item.
func (r *pgRepository) UpsertStoryItemReaction(ctx context.Context, storyID uuid.UUID, itemID uuid.UUID, reactorUserID uuid.UUID, reactorPersonaID uuid.UUID, reactionType models.StoryReactionType) (*models.StoryItemReaction, error) {
	row := r.db.QueryRow(ctx, `
		INSERT INTO story_item_reactions (story_id, story_item_id, reactor_user_id, reactor_persona_id, reaction_type)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (story_item_id, reactor_persona_id) DO UPDATE
		SET reactor_user_id = EXCLUDED.reactor_user_id,
		    reaction_type = EXCLUDED.reaction_type,
		    updated_at = now()
		RETURNING story_id, story_item_id, reactor_user_id, reactor_persona_id, reaction_type, created_at, updated_at
	`, storyID, itemID, reactorUserID, reactorPersonaID, reactionType)
	return scanStoryItemReaction(row)
}

// DeleteStoryItemReaction removes one viewer reaction from one story item.
func (r *pgRepository) DeleteStoryItemReaction(ctx context.Context, storyID uuid.UUID, itemID uuid.UUID, reactorPersonaID uuid.UUID) error {
	tag, err := r.db.Exec(ctx, `
		DELETE FROM story_item_reactions
		WHERE story_id = $1 AND story_item_id = $2 AND reactor_persona_id = $3
	`, storyID, itemID, reactorPersonaID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrStoryItemNotFound
	}
	return nil
}

// FindStoryItemReactionCounts returns reaction counts grouped by item and reaction type.
func (r *pgRepository) FindStoryItemReactionCounts(ctx context.Context, storyID uuid.UUID) (map[uuid.UUID]map[models.StoryReactionType]int, error) {
	rows, err := r.db.Query(ctx, `
		SELECT story_item_id, reaction_type, COUNT(*)::INT
		FROM story_item_reactions
		WHERE story_id = $1
		GROUP BY story_item_id, reaction_type
	`, storyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	counts := make(map[uuid.UUID]map[models.StoryReactionType]int)
	for rows.Next() {
		var itemID uuid.UUID
		var reactionType models.StoryReactionType
		var count int
		if err := rows.Scan(&itemID, &reactionType, &count); err != nil {
			return nil, err
		}
		if counts[itemID] == nil {
			counts[itemID] = make(map[models.StoryReactionType]int)
		}
		counts[itemID][reactionType] = count
	}
	return counts, rows.Err()
}

// FindStoryItemReactionsByPersona returns the current viewer reaction per story item.
func (r *pgRepository) FindStoryItemReactionsByPersona(ctx context.Context, storyID uuid.UUID, personaID uuid.UUID) (map[uuid.UUID]models.StoryReactionType, error) {
	rows, err := r.db.Query(ctx, `
		SELECT story_item_id, reaction_type
		FROM story_item_reactions
		WHERE story_id = $1 AND reactor_persona_id = $2
	`, storyID, personaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	reactions := make(map[uuid.UUID]models.StoryReactionType)
	for rows.Next() {
		var itemID uuid.UUID
		var reactionType models.StoryReactionType
		if err := rows.Scan(&itemID, &reactionType); err != nil {
			return nil, err
		}
		reactions[itemID] = reactionType
	}
	return reactions, rows.Err()
}

// scanStoryItemReaction scans one story reaction row into the model shape.
func scanStoryItemReaction(scanner storyScanner) (*models.StoryItemReaction, error) {
	var reaction models.StoryItemReaction
	if err := scanner.Scan(&reaction.StoryID, &reaction.StoryItemID, &reaction.ReactorUserID, &reaction.ReactorPersonaID, &reaction.ReactionType, &reaction.CreatedAt, &reaction.UpdatedAt); err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrStoryItemNotFound
		}
		return nil, err
	}
	return &reaction, nil
}
