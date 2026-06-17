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
	CreatePost(ctx context.Context, authorUserID uuid.UUID, dto dtos.CreatePostDTO) (*models.Post, error)
	CreatePostWithHashtags(ctx context.Context, authorUserID uuid.UUID, dto dtos.CreatePostDTO, tags []string) (*models.Post, error)
	CreatePostWithHashtagsAndMedia(ctx context.Context, authorUserID uuid.UUID, dto dtos.CreatePostDTO, tags []string, mediaAssetIDs []uuid.UUID) (*models.Post, error)
	FindPostByID(ctx context.Context, postID uuid.UUID) (*models.Post, error)
	FindPostsByPersonaID(ctx context.Context, personaID uuid.UUID, limit int) ([]*models.Post, error)
	FindFeedPosts(ctx context.Context, personaID uuid.UUID, limit int) ([]*models.Post, error)
	FindFollowingFeedPosts(ctx context.Context, personaID uuid.UUID, limit int) ([]*models.Post, error)
	EnsureAnonymousThreadIdentity(ctx context.Context, threadID uuid.UUID, userID uuid.UUID, personaID uuid.UUID, anonymousHandle string) (*models.AnonymousThreadIdentity, error)
	FindAnonymousThreadIdentities(ctx context.Context, threadID uuid.UUID, personaIDs []uuid.UUID) (map[uuid.UUID]*models.AnonymousThreadIdentity, error)
	FindMediaAssetsByPostIDs(ctx context.Context, postIDs []uuid.UUID) (map[uuid.UUID][]*models.MediaAsset, error)
	DeletePost(ctx context.Context, postID uuid.UUID) error
	DeletePostWithHashtags(ctx context.Context, postID uuid.UUID) error
}

func NewPostRepository(db *pgxpool.Pool) PostRepository {
	return newPgRepository(db)
}
