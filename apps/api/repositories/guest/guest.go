package guest

import (
	"context"
	"github.com/emmanuella-codes/nox/guest/dtos"
	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
	"time"
)

type GuestRepository interface {
	Configure(context.Context, uuid.UUID, uuid.UUID, dtos.ConfigureDTO) error
	CreateQuota(context.Context, uuid.UUID, uuid.UUID, dtos.QuotaDTO) (*models.GuestQuota, error)
	AddGuest(context.Context, uuid.UUID, uuid.UUID, uuid.UUID, dtos.GuestDTO, []string) ([]*models.GuestEntry, error)
	Claim(context.Context, string, *uuid.UUID) (*models.GuestEntry, error)
	Decline(context.Context, string) error
	Revoke(context.Context, uuid.UUID, uuid.UUID) error
	Attach(context.Context, uuid.UUID, uuid.UUID) error
	CheckIn(context.Context, uuid.UUID, uuid.UUID, string) (*models.GuestEntry, error)
	Report(context.Context, uuid.UUID, uuid.UUID) (*models.GuestListReport, error)
	ExpireInvitations(context.Context, time.Time) (int64, error)
	MarkNoShows(context.Context, time.Time) (int64, error)
}
