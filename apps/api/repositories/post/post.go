package post

import (
	"context"
	"errors"

	"github.com/emmanuella-codes/nox/models"
	"github.com/emmanuella-codes/nox/post/dtos"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrPostNotFound = errors.New("post not found")

type PostRepository interface {
	CreatePost(ctx context.Context, personaID uuid.UUID, dto dtos.CreatePostDTO) (*models.Post, error)
	FindPostByID(ctx context.Context, postID uuid.UUID) (*models.Post, error)
	FindPostsByPersonaID(ctx context.Context, personaID uuid.UUID, limit int) ([]*models.Post, error)
	FindFeedPosts(ctx context.Context, personaID uuid.UUID, limit int) ([]*models.Post, error)
	DeletePost(ctx context.Context, postID uuid.UUID) error
}

func NewPostRepository(db *pgxpool.Pool) PostRepository {
	return newPgRepository(db)
}
