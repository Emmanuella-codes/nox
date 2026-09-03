package story

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const storyAdditionWindow = 20 * time.Hour

type storyQueryRower interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// Checks whether a story can still accept new direct items or contributions.
func (r *pgRepository) StoryAcceptsAdditions(ctx context.Context, storyID uuid.UUID) (bool, error) {
	return storyAcceptsAdditions(ctx, r.db, storyID)
}

// Rejects pending contribution requests for stories that can no longer accept additions.
func (r *pgRepository) RejectPendingContributionRequestsForClosedStories(ctx context.Context, limit int) (int64, error) {
	commandTag, err := r.db.Exec(ctx, `
		WITH candidates AS (
			SELECT scr.id
			FROM story_contribution_requests scr
			JOIN stories s ON s.id = scr.story_id
			WHERE scr.status = 'pending'
			  AND (
				COALESCE((
					SELECT MIN(si.created_at)
					FROM story_items si
					WHERE si.story_id = s.id AND si.expires_at > now()
				), s.created_at) <= now() - interval '20 hours'
			  )
			ORDER BY scr.created_at ASC, scr.id ASC
			LIMIT $1
		)
		UPDATE story_contribution_requests scr
		SET status = 'rejected',
		    reviewed_at = now()
		FROM candidates
		WHERE scr.id = candidates.id
	`, normalizeLimit(limit))
	return commandTag.RowsAffected(), err
}

// Deletes expired items after the retention window when the story is not highlighted anywhere.
func (r *pgRepository) DeleteExpiredNonHighlightedStoryItems(ctx context.Context, olderThan time.Time, limit int) (int64, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, `
		WITH candidates AS (
			SELECT si.id, si.story_id
			FROM story_items si
			WHERE si.expires_at <= $1
			  AND NOT EXISTS (SELECT 1 FROM event_highlight_stories ehs WHERE ehs.story_id = si.story_id)
			  AND NOT EXISTS (SELECT 1 FROM profile_story_highlights psh WHERE psh.story_id = si.story_id)
			ORDER BY si.expires_at ASC, si.id ASC
			LIMIT $2
		),
		deleted AS (
			DELETE FROM story_items si
			USING candidates
			WHERE si.id = candidates.id
			RETURNING si.story_id
		)
		SELECT DISTINCT story_id
		FROM deleted
	`, olderThan, normalizeLimit(limit))
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	storyIDs, err := scanStoryUUIDs(rows)
	if err != nil {
		return 0, err
	}
	for _, storyID := range storyIDs {
		if err := recalculateStoryState(ctx, tx, storyID); err != nil {
			return 0, err
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}
	return int64(len(storyIDs)), nil
}

// Deletes old empty stories once their expired items have already been cleaned up.
func (r *pgRepository) DeleteRetainedEmptyStories(ctx context.Context, olderThan time.Time, limit int) (int64, error) {
	commandTag, err := r.db.Exec(ctx, `
		DELETE FROM stories
		WHERE id IN (
			SELECT s.id
			FROM stories s
			WHERE s.expires_at <= $1
			  AND NOT EXISTS (SELECT 1 FROM story_items si WHERE si.story_id = s.id)
			  AND NOT EXISTS (SELECT 1 FROM event_highlight_stories ehs WHERE ehs.story_id = s.id)
			  AND NOT EXISTS (SELECT 1 FROM profile_story_highlights psh WHERE psh.story_id = s.id)
			ORDER BY s.expires_at ASC, s.id ASC
			LIMIT $2
		)
	`, olderThan, normalizeLimit(limit))
	return commandTag.RowsAffected(), err
}

// Recomputes story duration and latest expiry from the currently attached items.
func recalculateStoryState(ctx context.Context, tx pgx.Tx, storyID uuid.UUID) error {
	_, err := tx.Exec(ctx, `
		UPDATE stories s
		SET total_duration_seconds = COALESCE((
				SELECT SUM(si.duration_seconds)
				FROM story_items si
				WHERE si.story_id = s.id
			), 0),
		    expires_at = COALESCE((
				SELECT MAX(si.expires_at)
				FROM story_items si
				WHERE si.story_id = s.id
			), s.created_at + interval '24 hours'),
		    updated_at = now()
		WHERE s.id = $1
	`, storyID)
	return err
}

// Uses the oldest still-active item, or the story creation time when empty, as the add window reference.
func storyAcceptsAdditions(ctx context.Context, rower storyQueryRower, storyID uuid.UUID) (bool, error) {
	var storyCreatedAt time.Time
	if err := rower.QueryRow(ctx, `SELECT created_at FROM stories WHERE id = $1`, storyID).Scan(&storyCreatedAt); err != nil {
		return false, mapStoryError(err)
	}
	var oldestActiveItemCreatedAt *time.Time
	if err := rower.QueryRow(ctx, `
		SELECT MIN(created_at)
		FROM story_items
		WHERE story_id = $1 AND expires_at > now()
	`, storyID).Scan(&oldestActiveItemCreatedAt); err != nil {
		return false, err
	}
	reference := storyCreatedAt
	if oldestActiveItemCreatedAt != nil {
		reference = *oldestActiveItemCreatedAt
	}
	return time.Since(reference) < storyAdditionWindow, nil
}
