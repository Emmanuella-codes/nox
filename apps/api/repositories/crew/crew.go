package crew

import (
	"context"
	"errors"
	"time"

	"github.com/emmanuella-codes/nox/crew/dtos"
	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrCrewNotFound       = errors.New("crew not found")
	ErrCrewCodeTaken      = errors.New("crew join code already taken")
	ErrCrewMemberNotFound = errors.New("crew member not found")
	ErrCrewFull           = errors.New("crew is full")
)

type CrewRepository interface {
	CreateCrew(ctx context.Context, ownerUserID uuid.UUID, eventID uuid.UUID, joinCode string, expiresAt time.Time, dto dtos.CreateCrewDTO) (*models.EventCrew, error)
	FindCrewByID(ctx context.Context, crewID uuid.UUID) (*models.EventCrew, error)
	FindCrewByJoinCode(ctx context.Context, joinCode string) (*models.EventCrew, error)
	FindEventCrewsForPersona(ctx context.Context, eventID uuid.UUID, personaID uuid.UUID, limit int, offset int) ([]*models.EventCrew, error)
	FindCrewMembers(ctx context.Context, crewID uuid.UUID) ([]*models.EventCrewMember, error)
	FindCrewMember(ctx context.Context, crewID uuid.UUID, personaID uuid.UUID) (*models.EventCrewMember, error)
	JoinCrew(ctx context.Context, crew *models.EventCrew, persona *models.Persona) (*models.EventCrewMember, error)
	LeaveCrew(ctx context.Context, crewID uuid.UUID, personaID uuid.UUID) error
	EndCrew(ctx context.Context, crewID uuid.UUID) (*models.EventCrew, error)
	UpdateLocationSharing(ctx context.Context, crewID uuid.UUID, personaID uuid.UUID, enabled bool) (*models.EventCrewMember, error)
	UpsertCrewLocation(ctx context.Context, crewID uuid.UUID, userID uuid.UUID, dto dtos.UpdateLocationDTO, expiresAt time.Time) (*models.EventCrewLocation, error)
	FindActiveCrewLocations(ctx context.Context, crewID uuid.UUID) ([]*models.EventCrewLocation, error)
	DeleteCrewLocation(ctx context.Context, crewID uuid.UUID, personaID uuid.UUID) error
	DeleteExpiredCrewLocations(ctx context.Context, now time.Time) (int64, error)
}

func NewCrewRepository(db *pgxpool.Pool) CrewRepository {
	return newPgRepository(db)
}
