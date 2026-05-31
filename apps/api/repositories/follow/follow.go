package follow

import (
	"context"
	"errors"

	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrAlreadyFollowing = errors.New("already following")
	ErrNotFollowing     = errors.New("not following")
	ErrSelfFollow       = errors.New("self follow")
)

type FollowRepository interface {
	Follow(ctx context.Context, followerID, followingID uuid.UUID) error
	Unfollow(ctx context.Context, followerID, followingID uuid.UUID) error
	IsFollowing(ctx context.Context, followerID, followingID uuid.UUID) (bool, error)
	FindFollowingIDs(ctx context.Context, followerID uuid.UUID, followingIDs []uuid.UUID) (map[uuid.UUID]bool, error)
	FindFollowers(ctx context.Context, personaID uuid.UUID, options ListOptions) ([]*models.Persona, error)
	FindFollowing(ctx context.Context, personaID uuid.UUID, options ListOptions) ([]*models.Persona, error)
}

func NewFollowRepository(db *pgxpool.Pool) FollowRepository {
	return newPgRepository(db)
}

type ListOptions struct {
	Limit  int
	Offset int
}

func NormalizeListOptions(options ListOptions) ListOptions {
	options.Limit = normalizeLimit(options.Limit)
	if options.Offset < 0 {
		options.Offset = 0
	}
	return options
}

func normalizeLimit(limit int) int {
	if limit <= 0 {
		return 20
	}
	if limit > 50 {
		return 50
	}
	return limit
}
