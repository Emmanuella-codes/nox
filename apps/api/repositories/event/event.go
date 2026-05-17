package event

import (
	"context"
	"errors"

	"github.com/emmanuella-codes/nox/event/dtos"
	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrEventNotFound = errors.New("event not found")

type EventRepository interface {
	CreateEvent(ctx context.Context, dto dtos.CreateEventDTO) (*models.Event, error)
	FindEventByID(ctx context.Context, eventID uuid.UUID) (*models.Event, error)
	FindEvents(ctx context.Context, limit int) ([]*models.Event, error)
}

func NewEventRepository(db *pgxpool.Pool) EventRepository {
	return newPgRepository(db)
}
