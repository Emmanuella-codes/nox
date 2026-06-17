package comment

import (
	"context"
	"errors"

	"github.com/emmanuella-codes/nox/comment/dtos"
	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgRepository struct {
	db *pgxpool.Pool
}

func newPgRepository(db *pgxpool.Pool) *pgRepository {
	return &pgRepository{db: db}
}

func (r *pgRepository) CreateComment(ctx context.Context, postID uuid.UUID, dto dtos.CreateCommentDTO) (*models.Comment, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	row := tx.QueryRow(ctx, `
		INSERT INTO comments (persona_id, post_id, posting_mode, body, parent_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, persona_id, post_id, posting_mode, body, parent_id, like_count, created_at
	`, dto.PersonaID, postID, dto.PostingMode, dto.Body, dto.ParentID)

	comment, err := scanComment(row)
	if err != nil {
		return nil, mapCommentError(err)
	}

	if _, err := tx.Exec(ctx, `UPDATE posts SET comment_count = comment_count + 1 WHERE id = $1`, postID); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return comment, nil
}

func (r *pgRepository) FindCommentsByPostID(ctx context.Context, postID uuid.UUID, limit int) ([]*models.Comment, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, persona_id, post_id, posting_mode, body, parent_id, like_count, created_at
		FROM comments
		WHERE post_id = $1
		ORDER BY created_at ASC
		LIMIT $2
	`, postID, normalizeLimit(limit))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*models.Comment
	for rows.Next() {
		comment, err := scanComment(rows)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, rows.Err()
}

func (r *pgRepository) FindCommentByID(ctx context.Context, commentID uuid.UUID) (*models.Comment, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, persona_id, post_id, posting_mode, body, parent_id, like_count, created_at
		FROM comments
		WHERE id = $1
	`, commentID)

	comment, err := scanComment(row)
	if err != nil {
		return nil, mapCommentError(err)
	}
	return comment, nil
}

func (r *pgRepository) DeleteComment(ctx context.Context, commentID uuid.UUID) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var postID uuid.UUID
	if err := tx.QueryRow(ctx, `DELETE FROM comments WHERE id = $1 RETURNING post_id`, commentID).Scan(&postID); err != nil {
		return mapCommentError(err)
	}

	if _, err := tx.Exec(ctx, `UPDATE posts SET comment_count = GREATEST(comment_count - 1, 0) WHERE id = $1`, postID); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

type commentScanner interface {
	Scan(dest ...any) error
}

func scanComment(scanner commentScanner) (*models.Comment, error) {
	var comment models.Comment
	var parentID uuid.NullUUID
	err := scanner.Scan(
		&comment.ID,
		&comment.PersonaID,
		&comment.PostID,
		&comment.PostingMode,
		&comment.Body,
		&parentID,
		&comment.LikeCount,
		&comment.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	if parentID.Valid {
		comment.ParentID = &parentID.UUID
	}
	return &comment, nil
}

func mapCommentError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrCommentNotFound
	}
	return err
}

func normalizeLimit(limit int) int {
	if limit <= 0 {
		return 50
	}
	if limit > 100 {
		return 100
	}
	return limit
}
