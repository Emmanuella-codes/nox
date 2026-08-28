package persona

import (
	"context"
	"errors"

	"github.com/emmanuella-codes/nox/models"
	"github.com/emmanuella-codes/nox/persona/dtos"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrPersonaNotFound         = errors.New("persona not found")
	ErrHandleAlreadyTaken      = errors.New("handle already taken")
	ErrProfileAlreadyExists    = errors.New("profile already exists for user")
	ErrGhostPersonaAlreadyUsed = errors.New("ghost persona already exists for user")
)

type PersonaRepository interface {
	CreatePersona(ctx context.Context, userID uuid.UUID, dto dtos.CreatePersonaDTO) (*models.Persona, error)
	FindPersonaByID(ctx context.Context, personaID uuid.UUID) (*models.Persona, error)
	FindPersonasByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Persona, error)
	FindPersonaByHandle(ctx context.Context, handle string) (*models.Persona, error)
	UpdatePersona(ctx context.Context, personaID uuid.UUID, dto dtos.UpdatePersonaDTO) (*models.Persona, error)
}

func NewPersonaRepository(db *pgxpool.Pool) PersonaRepository {
	return newPgRepository(db)
}
