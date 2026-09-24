package ticket

import (
	"context"
	"encoding/json"
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
	CreateOrder(context.Context, uuid.UUID, uuid.UUID) (*models.TicketOrder, error)
	FindOrder(context.Context, uuid.UUID, uuid.UUID) (*models.TicketOrder, error)
	SettlePayment(context.Context, string, string, int64, string, json.RawMessage) error
}

func NewTicketRepository(db *pgxpool.Pool) TicketRepository { return &pgRepository{db: db} }
