package admission

import (
	"context"
	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
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
