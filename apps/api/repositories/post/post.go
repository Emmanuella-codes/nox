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
	// CreatePost inserts one post row.
	CreatePost(ctx context.Context, authorUserID uuid.UUID, dto dtos.CreatePostDTO) (*models.Post, error)
	// CreatePostWithHashtags inserts a post and associated hashtags.
	CreatePostWithHashtags(ctx context.Context, authorUserID uuid.UUID, dto dtos.CreatePostDTO, tags []string) (*models.Post, error)
	// CreatePostWithHashtagsAndMedia inserts a post and associated hashtags and media.
	CreatePostWithHashtagsAndMedia(ctx context.Context, authorUserID uuid.UUID, dto dtos.CreatePostDTO, tags []string, mediaAssetIDs []uuid.UUID) (*models.Post, error)
	// FindPostByID fetches one post by id.
	FindPostByID(ctx context.Context, postID uuid.UUID) (*models.Post, error)
	// FindPostsByPersonaID fetches public posts for a public profile.
	FindPostsByPersonaID(ctx context.Context, personaID uuid.UUID, limit int) ([]*models.Post, error)
	// FindPostsByAuthorUserID fetches all posts for one owner, including anonymous posts.
	FindPostsByAuthorUserID(ctx context.Context, authorUserID uuid.UUID, limit int) ([]*models.Post, error)
	// FindFeedPosts fetches the mixed public and anonymous feed.
	FindFeedPosts(ctx context.Context, personaID uuid.UUID, limit int) ([]*models.Post, error)
	// FindFollowingFeedPosts fetches followed public-profile posts.
	FindFollowingFeedPosts(ctx context.Context, personaID uuid.UUID, limit int) ([]*models.Post, error)
	// EnsureAnonymousThreadIdentity creates or reuses one anonymous identity per thread owner.
	EnsureAnonymousThreadIdentity(ctx context.Context, threadID uuid.UUID, userID uuid.UUID, personaID uuid.UUID, anonymousHandle string, anonymousAvatarKey string) (*models.AnonymousThreadIdentity, error)
	// FindAnonymousThreadIdentities fetches thread identities for the supplied users.
	FindAnonymousThreadIdentities(ctx context.Context, threadID uuid.UUID, userIDs []uuid.UUID) (map[uuid.UUID]*models.AnonymousThreadIdentity, error)
	// FindMediaAssetsByPostIDs fetches media grouped by post id.
	FindMediaAssetsByPostIDs(ctx context.Context, postIDs []uuid.UUID) (map[uuid.UUID][]*models.MediaAsset, error)
	// DeletePost removes one post row.
	DeletePost(ctx context.Context, postID uuid.UUID) error
	// DeletePostWithHashtags removes one post and its hashtag links.
	DeletePostWithHashtags(ctx context.Context, postID uuid.UUID) error
}

// NewPostRepository builds the Postgres-backed post repository.
func NewPostRepository(db *pgxpool.Pool) PostRepository {
	return newPgRepository(db)
}
