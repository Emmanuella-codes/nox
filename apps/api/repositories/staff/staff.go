package staff

import (
	"context"
	"time"

	"github.com/emmanuella-codes/nox/models"
	"github.com/emmanuella-codes/nox/staff/dtos"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StaffRepository interface {
	Invite(context.Context, uuid.UUID, uuid.UUID, uuid.UUID) (*models.EventStaffMembership, error)
	AssignRoles(context.Context, uuid.UUID, uuid.UUID, []string) error
	Accept(context.Context, uuid.UUID, uuid.UUID) (*models.EventStaffMembership, error)
	Revoke(context.Context, uuid.UUID, uuid.UUID) error
	CreateShift(context.Context, uuid.UUID, uuid.UUID, dtos.ShiftDTO) (*models.EventStaffShift, error)
	List(context.Context, uuid.UUID, uuid.UUID) ([]*models.EventStaffMembership, error)
	Authorize(context.Context, uuid.UUID, uuid.UUID, string, time.Time) error
}

var StaffRepo StaffRepository

func InitStaffRepo(db *pgxpool.Pool) {
	StaffRepo = newPgRepository(db)
}
