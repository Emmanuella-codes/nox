package admission

import (
	"context"

	"github.com/emmanuella-codes/nox/models"
	staff_repo "github.com/emmanuella-codes/nox/repositories/staff"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AdmissionRepository interface {
	CheckIn(context.Context, uuid.UUID, uuid.UUID, string, string) (*models.CheckInResult, error)
	ManualCheckIn(context.Context, uuid.UUID, uuid.UUID, uuid.UUID) (*models.CheckInResult, error)
	OverrideCheckIn(context.Context, uuid.UUID, uuid.UUID, string) (*models.CheckInResult, error)
	Search(context.Context, uuid.UUID, uuid.UUID, string) ([]*models.AdmissionSearchResult, error)
	Reconcile(context.Context, uuid.UUID, uuid.UUID) (map[string]int, error)
	EnsureGuest(context.Context, uuid.UUID, uuid.UUID, string, string, string, string) error
	Create(context.Context, uuid.UUID, uuid.UUID, string, uuid.UUID, uuid.UUID, string, string, string, string) (*models.Admission, error)
}

var AdmissionRepo AdmissionRepository

func InitAdmissionRepo(db *pgxpool.Pool, staffRepo staff_repo.StaffRepository) {
	AdmissionRepo = newPgRepository(db, staffRepo)
}
