package admin

import (
	"context"
	"errors"

	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrMembershipAlreadyExists = errors.New("admin membership already exists")
	ErrMembershipNotFound      = errors.New("admin membership not found")
	ErrUserNotFound            = errors.New("user not found")
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
	ListIdentities(ctx context.Context) ([]models.AdminIdentity, error)
	FindMembershipByUserID(ctx context.Context, userID uuid.UUID) (*models.AdminMembership, error)
	FindIdentityByUserID(ctx context.Context, userID uuid.UUID) (*models.AdminIdentity, error)
	CreateMembership(ctx context.Context, userID uuid.UUID, role models.AdminRole) (*models.AdminIdentity, error)
	UpdateMembershipRole(ctx context.Context, userID uuid.UUID, role models.AdminRole) (*models.AdminIdentity, error)
	UpdateMembershipStatus(ctx context.Context, userID uuid.UUID, isActive bool) (*models.AdminIdentity, error)
	CountActiveMembershipsByRole(ctx context.Context, role models.AdminRole) (int, error)
	ListManagedUsers(ctx context.Context) ([]models.AdminManagedUser, error)
	FindManagedUserByID(ctx context.Context, userID uuid.UUID) (*models.AdminManagedUser, error)
	UpdateUserStatus(ctx context.Context, userID uuid.UUID, status models.UserStatus) (*models.AdminManagedUser, error)
	CreateAuditLog(ctx context.Context, params CreateAuditLogParams) error
}

func NewAdminRepository(db *pgxpool.Pool) AdminRepository {
	return newPgRepository(db)
}
