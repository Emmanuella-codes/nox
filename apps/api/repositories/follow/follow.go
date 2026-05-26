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
	FindFollowers(ctx context.Context, personaID uuid.UUID, limit int) ([]*models.Persona, error)
	FindFollowing(ctx context.Context, personaID uuid.UUID, limit int) ([]*models.Persona, error)
}

func NewFollowRepository(db *pgxpool.Pool) FollowRepository {
	return newPgRepository(db)
}
