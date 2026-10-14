package event

import (
	"context"
	"errors"
	"time"

	"github.com/emmanuella-codes/nox/event/dtos"
	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrEventNotFound = errors.New("event not found")
var ErrInvalidEventTransition = errors.New("invalid event transition")

type EventRepository interface {
	CreateEvent(ctx context.Context, dto dtos.CreateEventDTO) (*models.Event, error)
	UpdateEvent(ctx context.Context, userID, eventID uuid.UUID, dto dtos.UpdateEventDTO) (*models.Event, error)
	PublishEvent(ctx context.Context, userID, eventID uuid.UUID) (*models.Event, error)
	CloseEvent(ctx context.Context, userID, eventID uuid.UUID) (*models.Event, error)
	ArchiveEvent(ctx context.Context, userID, eventID uuid.UUID) (*models.Event, error)
	AdvanceLifecycle(ctx context.Context, now time.Time, limit int) (int, error)
	FindEventByID(ctx context.Context, eventID uuid.UUID) (*models.Event, error)
	FindEvents(ctx context.Context, limit int) ([]*models.Event, error)
}

var EventRepo EventRepository

func InitEventRepo(db *pgxpool.Pool) {
	EventRepo = newPgRepository(db)
}
