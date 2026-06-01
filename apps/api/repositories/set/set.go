package set

import (
	"context"
	"errors"

	"github.com/emmanuella-codes/nox/models"
	setdtos "github.com/emmanuella-codes/nox/set/dtos"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrSetNotFound = errors.New("set not found")

type SetRepository interface {
	CreateSet(ctx context.Context, authorUserID uuid.UUID, durationSeconds int, dto setdtos.CreateSetDTO) (*models.Set, error)
	FindSetByID(ctx context.Context, setID uuid.UUID) (*models.Set, error)
	FindSets(ctx context.Context, limit int, offset int) ([]*models.Set, error)
	FindSetsByPersonaID(ctx context.Context, personaID uuid.UUID, limit int, offset int) ([]*models.Set, error)
	DeleteSet(ctx context.Context, setID uuid.UUID) error
}

func NewSetRepository(db *pgxpool.Pool) SetRepository {
	return newPgRepository(db)
}
