package guest

import (
	"context"
	"github.com/emmanuella-codes/nox/guest/dtos"
	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
)

type GuestRepository interface {
	Configure(context.Context, uuid.UUID, uuid.UUID, dtos.ConfigureDTO) error
	CreateQuota(context.Context, uuid.UUID, uuid.UUID, dtos.QuotaDTO) (*models.GuestQuota, error)
	AddGuest(context.Context, uuid.UUID, uuid.UUID, uuid.UUID, dtos.GuestDTO, string) (*models.GuestEntry, error)
	Claim(context.Context, string, *uuid.UUID) (*models.GuestEntry, error)
	CheckIn(context.Context, uuid.UUID, uuid.UUID, string) (*models.GuestEntry, error)
}
