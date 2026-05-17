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
	CreateComment(ctx context.Context, postID uuid.UUID, dto dtos.CreateCommentDTO) (*models.Comment, error)
	FindCommentsByPostID(ctx context.Context, postID uuid.UUID, limit int) ([]*models.Comment, error)
	FindCommentByID(ctx context.Context, commentID uuid.UUID) (*models.Comment, error)
	DeleteComment(ctx context.Context, commentID uuid.UUID) error
}

func NewCommentRepository(db *pgxpool.Pool) CommentRepository {
	return newPgRepository(db)
}
