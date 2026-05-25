package like

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LikeRepository interface {
	LikePost(ctx context.Context, personaID uuid.UUID, postID uuid.UUID) error
	UnlikePost(ctx context.Context, personaID uuid.UUID, postID uuid.UUID) error
	HasPostLike(ctx context.Context, personaID uuid.UUID, postID uuid.UUID) (bool, error)
	FindLikedPostIDs(ctx context.Context, personaID uuid.UUID, postIDs []uuid.UUID) (map[uuid.UUID]bool, error)
}

func NewLikeRepository(db *pgxpool.Pool) LikeRepository {
	return newPgRepository(db)
}
