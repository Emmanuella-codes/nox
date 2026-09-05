package admin

import (
	"context"
	"errors"

	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrMembershipAlreadyExists  = errors.New("admin membership already exists")
	ErrMembershipNotFound       = errors.New("admin membership not found")
	ErrUserNotFound             = errors.New("user not found")
	ErrModerationEntityNotFound = errors.New("moderation entity not found")
	ErrReportAlreadyExists      = errors.New("open report already exists")
	ErrReportNotFound           = errors.New("report not found")
	ErrSelfReport               = errors.New("users cannot report their own content")
)

type CreateReportParams struct {
	ReporterUserID    uuid.UUID
	ReporterPersonaID uuid.UUID
	TargetType        models.ReportTargetType
	TargetID          uuid.UUID
	Reason            string
	Description       string
}

type ListReportsParams struct {
	Status     *models.ReportStatus
	TargetType *models.ReportTargetType
	AssignedID *uuid.UUID
	Limit      int
	Offset     int
}

type UpdateReportParams struct {
	ReportID        uuid.UUID
	Status          models.ReportStatus
	AssignedAdminID *uuid.UUID
	ResolutionNote  string
}

type ResolveReportParams struct {
	ReportID   uuid.UUID
	Status     models.ModerationStatus
	Reason     string
	AdminID    uuid.UUID
	Resolution string
}

type ModerateParams struct {
	EntityType models.ModerationEntityType
	EntityID   uuid.UUID
	Status     models.ModerationStatus
	Reason     string
	AdminID    uuid.UUID
}

type ListModerationActionsParams struct {
	EntityType *models.ModerationEntityType
	Status     *models.ModerationStatus
	Limit      int
	Offset     int
}

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
	Moderate(ctx context.Context, params ModerateParams) (*models.ModerationState, error)
	ListModerationActions(ctx context.Context, params ListModerationActionsParams) ([]models.ModerationAction, error)
	CreateReport(ctx context.Context, params CreateReportParams) (*models.Report, error)
	ListReports(ctx context.Context, params ListReportsParams) ([]models.Report, error)
	FindReportByID(ctx context.Context, reportID uuid.UUID) (*models.Report, error)
	UpdateReport(ctx context.Context, params UpdateReportParams) (*models.Report, error)
	ResolveReport(ctx context.Context, params ResolveReportParams) (*models.Report, error)
	Dashboard(ctx context.Context) (*models.AdminDashboard, error)
}

func NewAdminRepository(db *pgxpool.Pool) AdminRepository {
	return newPgRepository(db)
}
