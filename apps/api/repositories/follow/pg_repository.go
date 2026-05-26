package follow

import (
	"context"

	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgRepository struct {
	db *pgxpool.Pool
}

func newPgRepository(db *pgxpool.Pool) *pgRepository {
	return &pgRepository{db: db}
}

func (r *pgRepository) Follow(ctx context.Context, followerID, followingID uuid.UUID) error {}

func (r *pgRepository) Unfollow(ctx context.Context, followerID, followingID uuid.UUID) error {}

func (r *pgRepository) IsFollowing(ctx context.Context, followerID, followingID uuid.UUID) (bool, error) {}

func (r *pgRepository) FindFollowers(ctx context.Context, personaID uuid.UUID, limit int) ([]*models.Persona, error) {}
