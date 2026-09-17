package ticket

import (
	"context"
	"errors"
	"github.com/emmanuella-codes/nox/models"
	"github.com/emmanuella-codes/nox/ticket/dtos"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

var ErrNotFound = errors.New("ticket resource not found")
var ErrUnavailable = errors.New("ticket inventory unavailable")

type TicketRepository interface {
	CreateWave(context.Context, uuid.UUID, uuid.UUID, dtos.CreateWaveDTO) (*models.TicketWave, error)
	CreateTier(context.Context, uuid.UUID, uuid.UUID, dtos.CreateTierDTO) (*models.TicketTier, error)
	FindProducts(context.Context, uuid.UUID) ([]*models.TicketWave, error)
	ReleaseWave(context.Context, uuid.UUID, uuid.UUID) (*models.TicketWave, error)
	AdvanceSales(context.Context, time.Time, int) (int64, error)
	Reserve(context.Context, uuid.UUID, uuid.UUID, dtos.ReserveDTO, time.Duration) (*models.TicketReservation, error)
	ReleaseExpired(context.Context, time.Time, int) (int64, error)
}

func NewTicketRepository(db *pgxpool.Pool) TicketRepository { return &pgRepository{db: db} }
