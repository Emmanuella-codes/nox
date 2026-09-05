package set

import (
	"context"

	"github.com/emmanuella-codes/nox/models"
	setdtos "github.com/emmanuella-codes/nox/set/dtos"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgRepository struct {
	db *pgxpool.Pool
}

func newPgRepository(db *pgxpool.Pool) *pgRepository {
	return &pgRepository{db: db}
}

func (r *pgRepository) CreateSet(ctx context.Context, authorUserID uuid.UUID, durationSeconds int, dto setdtos.CreateSetDTO) (*models.Set, error) {
	row := r.db.QueryRow(ctx, `
		INSERT INTO sets (author_user_id, persona_id, media_asset_id, title, description, genre_tags, duration_seconds)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, author_user_id, persona_id, media_asset_id, title, description, genre_tags,
		          duration_seconds, like_count, comment_count, play_count, created_at, updated_at
	`, authorUserID, dto.PersonaID, dto.MediaAssetID, dto.Title, dto.Description, dto.GenreTags, durationSeconds)

	set, err := scanSet(row)
	if err != nil {
		if isUniqueViolation(err, "sets_media_asset_id_key") {
			return nil, ErrSetMediaInUse
		}
		return nil, mapSetError(err)
	}
	return set, nil
}

func (r *pgRepository) FindSetByID(ctx context.Context, setID uuid.UUID) (*models.Set, error) {
	row := r.db.QueryRow(ctx, `
		SELECT s.id, s.author_user_id, s.persona_id, s.media_asset_id, s.title, s.description, s.genre_tags,
		       s.duration_seconds, s.like_count, s.comment_count, s.play_count, s.created_at, s.updated_at
		FROM sets s
		INNER JOIN personas p ON p.id = s.persona_id AND p.moderation_status = 'active'
		WHERE s.id = $1 AND s.moderation_status = 'active'
	`, setID)

	set, err := scanSet(row)
	if err != nil {
		return nil, mapSetError(err)
	}
	return set, nil
}

func (r *pgRepository) FindSets(ctx context.Context, limit int, offset int) ([]*models.Set, error) {
	return r.FindSetsWithFilters(ctx, "", "newest", limit, offset)
}

func (r *pgRepository) FindSetsWithFilters(ctx context.Context, genreTag string, sort string, limit int, offset int) ([]*models.Set, error) {
	orderBy := "s.created_at DESC"
	switch sort {
	case "most_played":
		orderBy = "s.play_count DESC, s.created_at DESC"
	case "most_liked":
		orderBy = "s.like_count DESC, s.created_at DESC"
	case "most_discussed":
		orderBy = "s.comment_count DESC, s.created_at DESC"
	}
	rows, err := r.db.Query(ctx, `
		SELECT s.id, s.author_user_id, s.persona_id, s.media_asset_id, s.title, s.description, s.genre_tags,
		       s.duration_seconds, s.like_count, s.comment_count, s.play_count, s.created_at, s.updated_at
		FROM sets s
		INNER JOIN personas p ON p.id = s.persona_id AND p.moderation_status = 'active'
		WHERE ($1 = '' OR $1 = ANY(s.genre_tags))
		  AND s.moderation_status = 'active'
		ORDER BY `+orderBy+`
		LIMIT $2 OFFSET $3
	`, genreTag, normalizeLimit(limit), normalizeOffset(offset))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanSets(rows)
}

func (r *pgRepository) FindSetsByPersonaID(ctx context.Context, personaID uuid.UUID, limit int, offset int) ([]*models.Set, error) {
	rows, err := r.db.Query(ctx, `
		SELECT s.id, s.author_user_id, s.persona_id, s.media_asset_id, s.title, s.description, s.genre_tags,
		       s.duration_seconds, s.like_count, s.comment_count, s.play_count, s.created_at, s.updated_at
		FROM sets s
		INNER JOIN personas p ON p.id = s.persona_id AND p.moderation_status = 'active'
		WHERE s.persona_id = $1 AND s.moderation_status = 'active'
		ORDER BY s.created_at DESC
		LIMIT $2 OFFSET $3
	`, personaID, normalizeLimit(limit), normalizeOffset(offset))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanSets(rows)
}

func (r *pgRepository) DeleteSet(ctx context.Context, setID uuid.UUID) error {
	commandTag, err := r.db.Exec(ctx, `DELETE FROM sets WHERE id = $1`, setID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return ErrSetNotFound
	}
	return nil
}
