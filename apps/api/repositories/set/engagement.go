package set

import (
	"context"

	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
)

func (r *pgRepository) LikeSet(ctx context.Context, personaID uuid.UUID, setID uuid.UUID) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	commandTag, err := tx.Exec(ctx, `
		INSERT INTO set_likes (persona_id, set_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`, personaID, setID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() > 0 {
		if _, err := tx.Exec(ctx, `UPDATE sets SET like_count = like_count + 1 WHERE id = $1`, setID); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (r *pgRepository) UnlikeSet(ctx context.Context, personaID uuid.UUID, setID uuid.UUID) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	commandTag, err := tx.Exec(ctx, `DELETE FROM set_likes WHERE persona_id = $1 AND set_id = $2`, personaID, setID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() > 0 {
		if _, err := tx.Exec(ctx, `UPDATE sets SET like_count = GREATEST(like_count - 1, 0) WHERE id = $1`, setID); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (r *pgRepository) HasSetLike(ctx context.Context, personaID uuid.UUID, setID uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM set_likes WHERE persona_id = $1 AND set_id = $2
		)
	`, personaID, setID).Scan(&exists)
	return exists, err
}

func (r *pgRepository) FindLikedSetIDs(ctx context.Context, personaID uuid.UUID, setIDs []uuid.UUID) (map[uuid.UUID]bool, error) {
	liked := make(map[uuid.UUID]bool, len(setIDs))
	if len(setIDs) == 0 {
		return liked, nil
	}
	rows, err := r.db.Query(ctx, `
		SELECT set_id
		FROM set_likes
		WHERE persona_id = $1 AND set_id = ANY($2)
	`, personaID, setIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var setID uuid.UUID
		if err := rows.Scan(&setID); err != nil {
			return nil, err
		}
		liked[setID] = true
	}
	return liked, rows.Err()
}

func (r *pgRepository) IncrementPlayCount(ctx context.Context, setID uuid.UUID) error {
	commandTag, err := r.db.Exec(ctx, `UPDATE sets SET play_count = play_count + 1 WHERE id = $1`, setID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return ErrSetNotFound
	}
	return nil
}

func (r *pgRepository) CreateSetComment(ctx context.Context, personaID uuid.UUID, setID uuid.UUID, body string, parentID uuid.UUID) (*models.SetComment, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	row := tx.QueryRow(ctx, `
		INSERT INTO set_comments (persona_id, set_id, body, parent_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, persona_id, set_id, body, COALESCE(parent_id, '00000000-0000-0000-0000-000000000000'::uuid),
		          like_count, created_at
	`, personaID, setID, body, uuidToNil(parentID))
	comment, err := scanSetComment(row)
	if err != nil {
		return nil, err
	}
	if _, err := tx.Exec(ctx, `UPDATE sets SET comment_count = comment_count + 1 WHERE id = $1`, setID); err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return comment, nil
}

func (r *pgRepository) FindSetComments(ctx context.Context, setID uuid.UUID, limit int, offset int) ([]*models.SetComment, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, persona_id, set_id, body, COALESCE(parent_id, '00000000-0000-0000-0000-000000000000'::uuid),
		       like_count, created_at
		FROM set_comments
		WHERE set_id = $1
		ORDER BY created_at ASC
		LIMIT $2 OFFSET $3
	`, setID, normalizeLimit(limit), normalizeOffset(offset))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanSetComments(rows)
}
