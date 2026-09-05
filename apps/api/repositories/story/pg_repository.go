package story

import (
	"context"

	"github.com/emmanuella-codes/nox/models"
	storydtos "github.com/emmanuella-codes/nox/story/dtos"
	"github.com/google/uuid"
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
		INSERT INTO stories (event_id, owner_user_id, owner_persona_id, title, contribution_mode, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, event_id, owner_user_id, owner_persona_id, title, contribution_mode,
		          total_duration_seconds, expires_at, created_at, updated_at
	`, dto.EventID, ownerUserID, dto.OwnerPersonaID, dto.Title, dto.ContributionMode, dto.ExpiresAt)
	return scanStory(row)
}

func (r *pgRepository) FindStoryByID(ctx context.Context, storyID uuid.UUID) (*models.Story, error) {
	row := r.db.QueryRow(ctx, `
		SELECT s.id, s.event_id, s.owner_user_id, s.owner_persona_id, s.title, s.contribution_mode,
		       s.total_duration_seconds, s.expires_at, s.created_at, s.updated_at
		FROM stories s
		INNER JOIN personas p ON p.id = s.owner_persona_id AND p.moderation_status = 'active'
		INNER JOIN events e ON e.id = s.event_id AND e.moderation_status = 'active'
		WHERE s.id = $1 AND s.expires_at > now() AND s.moderation_status = 'active'
	`, storyID)
	story, err := scanStory(row)
	if err != nil {
		return nil, mapStoryError(err)
	}
	return story, nil
}

// FindStoryByIDAny loads one story regardless of expiry state.
func (r *pgRepository) FindStoryByIDAny(ctx context.Context, storyID uuid.UUID) (*models.Story, error) {
	row := r.db.QueryRow(ctx, `
		SELECT s.id, s.event_id, s.owner_user_id, s.owner_persona_id, s.title, s.contribution_mode,
		       s.total_duration_seconds, s.expires_at, s.created_at, s.updated_at
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
		SELECT s.id, s.event_id, s.owner_user_id, s.owner_persona_id, s.title, s.contribution_mode,
		       s.total_duration_seconds, s.expires_at, s.created_at, s.updated_at
		FROM stories s
		INNER JOIN personas p ON p.id = s.owner_persona_id AND p.moderation_status = 'active'
		INNER JOIN events e ON e.id = s.event_id AND e.moderation_status = 'active'
		WHERE s.event_id = $1 AND s.expires_at > now() AND s.moderation_status = 'active'
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
		SELECT s.id, s.event_id, s.owner_user_id, s.owner_persona_id, s.title, s.contribution_mode,
		       s.total_duration_seconds, s.expires_at, s.created_at, s.updated_at
		FROM stories s
		INNER JOIN personas p ON p.id = s.owner_persona_id AND p.moderation_status = 'active'
		INNER JOIN events e ON e.id = s.event_id AND e.moderation_status = 'active'
		WHERE s.owner_persona_id = $1 AND s.expires_at > now() AND s.moderation_status = 'active'
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
