package set

import (
	"context"

	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

func (r *pgRepository) RecordSetPlay(ctx context.Context, setID uuid.UUID, userID uuid.UUID) (bool, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer tx.Rollback(ctx)
	var exists bool
	if err := tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM sets WHERE id = $1 AND moderation_status = 'active')`, setID).Scan(&exists); err != nil {
		return false, err
	}
	if !exists {
		return false, ErrSetNotFound
	}
	result, err := tx.Exec(ctx, `INSERT INTO set_plays (set_id, user_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`, setID, userID)
	if err != nil {
		return false, err
	}
	if result.RowsAffected() == 0 {
		if err := tx.Commit(ctx); err != nil {
			return false, err
		}
		return false, nil
	}
	if _, err := tx.Exec(ctx, `UPDATE sets SET play_count = play_count + 1 WHERE id = $1`, setID); err != nil {
		return false, err
	}
	if err := tx.Commit(ctx); err != nil {
		return false, err
	}
	return true, nil
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
		          like_count, created_at, updated_at, deleted_at
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
		       like_count, created_at, updated_at, deleted_at
		FROM set_comments
		WHERE set_id = $1 AND deleted_at IS NULL
		ORDER BY created_at ASC
		LIMIT $2 OFFSET $3
	`, setID, normalizeLimit(limit), normalizeOffset(offset))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanSetComments(rows)
}

func (r *pgRepository) UpdateSetComment(ctx context.Context, personaID, commentID uuid.UUID, body string) (*models.SetComment, error) {
	row := r.db.QueryRow(ctx, `
		UPDATE set_comments
		SET body = $3, updated_at = now()
		WHERE id = $1 AND persona_id = $2 AND deleted_at IS NULL
		RETURNING id, persona_id, set_id, body, COALESCE(parent_id, '00000000-0000-0000-0000-000000000000'::uuid), like_count, created_at, updated_at, deleted_at
	`, commentID, personaID, body)
	comment, err := scanSetComment(row)
	if err != nil {
		if err == ErrSetCommentNotFound {
			return nil, err
		}
		return nil, mapSetCommentError(err)
	}
	return comment, nil
}

func (r *pgRepository) DeleteSetComment(ctx context.Context, personaID, commentID uuid.UUID) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	var setID uuid.UUID
	if err := tx.QueryRow(ctx, `UPDATE set_comments SET deleted_at = now(), updated_at = now() WHERE id = $1 AND persona_id = $2 AND deleted_at IS NULL RETURNING set_id`, commentID, personaID).Scan(&setID); err != nil {
		return mapSetCommentError(err)
	}
	if _, err := tx.Exec(ctx, `UPDATE sets SET comment_count = GREATEST(comment_count - 1, 0) WHERE id = $1`, setID); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *pgRepository) LikeSetComment(ctx context.Context, personaID, commentID uuid.UUID) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	created, err := tx.Exec(ctx, `INSERT INTO set_comment_likes (persona_id, comment_id) SELECT $1, id FROM set_comments WHERE id = $2 AND deleted_at IS NULL ON CONFLICT DO NOTHING`, personaID, commentID)
	if err != nil {
		return err
	}
	if created.RowsAffected() == 0 {
		return tx.Commit(ctx)
	}
	if _, err := tx.Exec(ctx, `UPDATE set_comments SET like_count = like_count + 1 WHERE id = $1`, commentID); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *pgRepository) UnlikeSetComment(ctx context.Context, personaID, commentID uuid.UUID) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	removed, err := tx.Exec(ctx, `DELETE FROM set_comment_likes WHERE persona_id = $1 AND comment_id = $2`, personaID, commentID)
	if err != nil {
		return err
	}
	if removed.RowsAffected() > 0 {
		if _, err := tx.Exec(ctx, `UPDATE set_comments SET like_count = GREATEST(like_count - 1, 0) WHERE id = $1`, commentID); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (r *pgRepository) HasSetCommentLike(ctx context.Context, personaID, commentID uuid.UUID) (bool, error) {
	var liked bool
	err := r.db.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM set_comment_likes WHERE persona_id = $1 AND comment_id = $2)`, personaID, commentID).Scan(&liked)
	return liked, err
}

func mapSetCommentError(err error) error {
	if err == pgx.ErrNoRows {
		return ErrSetCommentNotFound
	}
	return err
}
