package set

import (
	"context"
	"errors"

	"github.com/emmanuella-codes/nox/models"
	setdtos "github.com/emmanuella-codes/nox/set/dtos"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrSetNotFound        = errors.New("set not found")
	ErrSetCommentNotFound = errors.New("set comment not found")
	ErrSetMediaInUse      = errors.New("set media asset already used")
)

type SetRepository interface {
	CreateSet(ctx context.Context, authorUserID uuid.UUID, durationSeconds int, dto setdtos.CreateSetDTO) (*models.Set, error)
	FindSetByID(ctx context.Context, setID uuid.UUID) (*models.Set, error)
	FindSets(ctx context.Context, limit int, offset int) ([]*models.Set, error)
	FindSetsWithFilters(ctx context.Context, genreTag string, sort string, limit int, offset int) ([]*models.Set, error)
	FindSetsByPersonaID(ctx context.Context, personaID uuid.UUID, limit int, offset int) ([]*models.Set, error)
	DeleteSet(ctx context.Context, setID uuid.UUID) error
	LikeSet(ctx context.Context, personaID uuid.UUID, setID uuid.UUID) error
	UnlikeSet(ctx context.Context, personaID uuid.UUID, setID uuid.UUID) error
	HasSetLike(ctx context.Context, personaID uuid.UUID, setID uuid.UUID) (bool, error)
	FindLikedSetIDs(ctx context.Context, personaID uuid.UUID, setIDs []uuid.UUID) (map[uuid.UUID]bool, error)
	IncrementPlayCount(ctx context.Context, setID uuid.UUID) error
	CreateSetComment(ctx context.Context, personaID uuid.UUID, setID uuid.UUID, body string, parentID uuid.UUID) (*models.SetComment, error)
	FindSetComments(ctx context.Context, setID uuid.UUID, limit int, offset int) ([]*models.SetComment, error)
}

func NewSetRepository(db *pgxpool.Pool) SetRepository {
	return newPgRepository(db)
}
