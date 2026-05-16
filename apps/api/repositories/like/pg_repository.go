package like

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgRepository struct {
	db *pgxpool.Pool
}

func newPgRepository(db *pgxpool.Pool) *pgRepository {
	return &pgRepository{db: db}
}

func (r *pgRepository) LikePost(ctx context.Context, personaID uuid.UUID, postID uuid.UUID) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	commandTag, err := tx.Exec(ctx, `
		INSERT INTO post_likes (persona_id, post_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`, personaID, postID)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() > 0 {
		if _, err := tx.Exec(ctx, `UPDATE posts SET like_count = like_count + 1 WHERE id = $1`, postID); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *pgRepository) UnlikePost(ctx context.Context, personaID uuid.UUID, postID uuid.UUID) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	commandTag, err := tx.Exec(ctx, `DELETE FROM post_likes WHERE persona_id = $1 AND post_id = $2`, personaID, postID)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() > 0 {
		if _, err := tx.Exec(ctx, `UPDATE posts SET like_count = GREATEST(like_count - 1, 0) WHERE id = $1`, postID); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
