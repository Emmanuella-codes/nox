package comment

import (
	"context"
	"errors"

	"github.com/emmanuella-codes/nox/comment/dtos"
	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrCommentNotFound = errors.New("comment not found")

type CommentRepository interface {
	// CreateComment inserts one comment row.
	CreateComment(ctx context.Context, postID uuid.UUID, dto dtos.CreateCommentDTO) (*models.Comment, error)
	// FindCommentsByPostID fetches comments for one post.
	FindCommentsByPostID(ctx context.Context, postID uuid.UUID, limit int) ([]*models.Comment, error)
	// FindCommentByID fetches one comment by id.
	FindCommentByID(ctx context.Context, commentID uuid.UUID) (*models.Comment, error)
	// DeleteComment removes one comment row.
	DeleteComment(ctx context.Context, commentID uuid.UUID) error
}

// NewCommentRepository builds the Postgres-backed comment repository.
func NewCommentRepository(db *pgxpool.Pool) CommentRepository {
	return newPgRepository(db)
}
