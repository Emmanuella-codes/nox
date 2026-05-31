package follow

import (
	"context"

	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgRepository struct {
	db *pgxpool.Pool
}

func newPgRepository(db *pgxpool.Pool) *pgRepository {
	return &pgRepository{db: db}
}

func (r *pgRepository) Follow(ctx context.Context, followerID, followingID uuid.UUID) error {
	if followerID == followingID {
		return ErrSelfFollow
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	commandTag, err := tx.Exec(ctx, `
		INSERT INTO persona_follows (follower_id, following_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`, followerID, followingID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return ErrAlreadyFollowing
	}

	if _, err := tx.Exec(ctx, `UPDATE personas SET following_count = following_count + 1 WHERE id = $1`, followerID); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `UPDATE personas SET follower_count = follower_count + 1 WHERE id = $1`, followingID); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *pgRepository) Unfollow(ctx context.Context, followerID, followingID uuid.UUID) error {
	if followerID == followingID {
		return ErrSelfFollow
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	commandTag, err := tx.Exec(ctx, `
		DELETE FROM persona_follows
		WHERE follower_id = $1 AND following_id = $2
	`, followerID, followingID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return ErrNotFollowing
	}

	if _, err := tx.Exec(ctx, `UPDATE personas SET following_count = GREATEST(following_count - 1, 0) WHERE id = $1`, followerID); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `UPDATE personas SET follower_count = GREATEST(follower_count - 1, 0) WHERE id = $1`, followingID); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *pgRepository) IsFollowing(ctx context.Context, followerID, followingID uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM persona_follows WHERE follower_id = $1 AND following_id = $2
		)
	`, followerID, followingID).Scan(&exists)
	return exists, err
}

func (r *pgRepository) FindFollowingIDs(ctx context.Context, followerID uuid.UUID, followingIDs []uuid.UUID) (map[uuid.UUID]bool, error) {
	following := make(map[uuid.UUID]bool, len(followingIDs))
	if len(followingIDs) == 0 {
		return following, nil
	}

	rows, err := r.db.Query(ctx, `
		SELECT following_id
		FROM persona_follows
		WHERE follower_id = $1 AND following_id = ANY($2)
	`, followerID, followingIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var followingID uuid.UUID
		if err := rows.Scan(&followingID); err != nil {
			return nil, err
		}
		following[followingID] = true
	}
	return following, rows.Err()
}

func (r *pgRepository) FindFollowers(ctx context.Context, personaID uuid.UUID, options ListOptions) ([]*models.Persona, error) {
	options = NormalizeListOptions(options)
	rows, err := r.db.Query(ctx, `
		SELECT p.id, p.user_id, p.handle, p.display_name, p.bio, p.avatar_url, p.cover_url, p.persona_type,
		       p.genre_tags, p.follower_count, p.following_count, p.post_count, p.created_at, p.updated_at
		FROM persona_follows pf
		INNER JOIN personas p ON p.id = pf.follower_id
		WHERE pf.following_id = $1 AND p.persona_type = 'visible'
		ORDER BY pf.created_at DESC
		LIMIT $2
		OFFSET $3
	`, personaID, options.Limit+1, options.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanPersonas(rows)
}

func (r *pgRepository) FindFollowing(ctx context.Context, personaID uuid.UUID, options ListOptions) ([]*models.Persona, error) {
	options = NormalizeListOptions(options)
	rows, err := r.db.Query(ctx, `
		SELECT p.id, p.user_id, p.handle, p.display_name, p.bio, p.avatar_url, p.cover_url, p.persona_type,
		       p.genre_tags, p.follower_count, p.following_count, p.post_count, p.created_at, p.updated_at
		FROM persona_follows pf
		INNER JOIN personas p ON p.id = pf.following_id
		WHERE pf.follower_id = $1 AND p.persona_type = 'visible'
		ORDER BY pf.created_at DESC
		LIMIT $2
		OFFSET $3
	`, personaID, options.Limit+1, options.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanPersonas(rows)
}

type personaScanner interface {
	Scan(dest ...any) error
}

type personaRows interface {
	personaScanner
	Next() bool
	Err() error
}

func scanPersonas(rows personaRows) ([]*models.Persona, error) {
	var personas []*models.Persona
	for rows.Next() {
		persona, err := scanPersona(rows)
		if err != nil {
			return nil, err
		}
		personas = append(personas, persona)
	}
	return personas, rows.Err()
}

func scanPersona(scanner personaScanner) (*models.Persona, error) {
	var persona models.Persona
	err := scanner.Scan(
		&persona.ID,
		&persona.UserID,
		&persona.Handle,
		&persona.DisplayName,
		&persona.Bio,
		&persona.AvatarURL,
		&persona.CoverURL,
		&persona.PersonaType,
		&persona.GenreTags,
		&persona.FollowerCount,
		&persona.FollowingCount,
		&persona.PostCount,
		&persona.CreatedAt,
		&persona.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &persona, nil
}
