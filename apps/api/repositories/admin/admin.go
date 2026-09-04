package admin

import (
	"context"

	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateAuditLogParams struct {
	AdminUserID  uuid.UUID
	Action       string
	TargetUserID *uuid.UUID
	RequestID    string
	IPAddress    string
	UserAgent    string
	Metadata     map[string]any
}

type AdminRepository interface {
	FindMembershipByUserID(ctx context.Context, userID uuid.UUID) (*models.AdminMembership, error)
	FindIdentityByUserID(ctx context.Context, userID uuid.UUID) (*models.AdminIdentity, error)
	CreateAuditLog(ctx context.Context, params CreateAuditLogParams) error
}

func NewAdminRepository(db *pgxpool.Pool) AdminRepository {
	return newPgRepository(db)
}
