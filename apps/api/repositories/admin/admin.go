package admin

import (
	"context"
	"errors"
	"time"

	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrMembershipAlreadyExists   = errors.New("admin membership already exists")
	ErrMembershipNotFound        = errors.New("admin membership not found")
	ErrUserNotFound              = errors.New("user not found")
	ErrModerationEntityNotFound  = errors.New("moderation entity not found")
	ErrReportAlreadyExists       = errors.New("open report already exists")
	ErrReportNotFound            = errors.New("report not found")
	ErrSelfReport                = errors.New("users cannot report their own content")
	ErrHashtagNotFound           = errors.New("hashtag not found")
	ErrSearchSuppressionNotFound = errors.New("search suppression not found")
	ErrMediaNotFound             = errors.New("media asset not found")
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

type ListAdminContentParams struct {
	EntityType  models.ModerationEntityType
	Status      *models.ModerationStatus
	OwnerID     *uuid.UUID
	ParentID    *uuid.UUID
	CreatedFrom *time.Time
	CreatedTo   *time.Time
	Limit       int
	Offset      int
}

type ModerateParams struct {
	EntityType models.ModerationEntityType
	EntityID   uuid.UUID
	Status     models.ModerationStatus
	Reason     string
	AdminID    uuid.UUID
}

type ListAdminHashtagsParams struct {
	Query  string
	Status *models.ModerationStatus
	Limit  int
	Offset int
}

type ModerateHashtagParams struct {
	Tag     string
	Status  models.ModerationStatus
	Reason  string
	AdminID uuid.UUID
}

type SuppressHashtagParams struct {
	Tag       string
	Reason    string
	AdminID   uuid.UUID
	ExpiresAt *time.Time
}

type CreateSearchSuppressionParams struct {
	Query     string
	Reason    string
	AdminID   uuid.UUID
	ExpiresAt *time.Time
}

type ListSearchSuppressionsParams struct {
	Query  string
	Limit  int
	Offset int
}

type ListAdminMediaParams struct {
	Status     *models.MediaProcessingStatus
	Kind       *models.MediaKind
	Moderation *models.ModerationStatus
	OwnerID    *uuid.UUID
	Orphaned   *bool
	OlderThan  *time.Time
	Limit      int
	Offset     int
}

type CorrectMediaParams struct {
	MediaAssetID uuid.UUID
	Status       models.MediaProcessingStatus
	PlaybackURL  string
	ThumbnailURL string
	MimeType     string
	Duration     int
	SizeBytes    int64
	Reason       string
	AdminID      uuid.UUID
}

type ModerateMediaParams struct {
	MediaAssetID uuid.UUID
	Status       models.ModerationStatus
	Reason       string
	AdminID      uuid.UUID
}

type ListFeaturedSetsParams struct {
	Limit  int
	Offset int
}

type ListStoryContributionsParams struct {
	Status *models.StoryContributionRequestStatus
	Limit  int
	Offset int
}

type ReviewStoryContributionParams struct {
	RequestID uuid.UUID
	Action    string
	Reason    string
	AdminID   uuid.UUID
}

type ListAdminHighlightsParams struct {
	HighlightType string
	Limit         int
	Offset        int
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
	ListAdminContent(ctx context.Context, params ListAdminContentParams) (*models.AdminContentPage, error)
	FindAdminContent(ctx context.Context, entityType models.ModerationEntityType, entityID uuid.UUID) (*models.AdminContentRecord, error)
	ListAdminHashtags(ctx context.Context, params ListAdminHashtagsParams) (*models.AdminHashtagPage, error)
	ModerateHashtag(ctx context.Context, params ModerateHashtagParams) (*models.AdminHashtag, error)
	SuppressHashtag(ctx context.Context, params SuppressHashtagParams) (*models.AdminHashtag, error)
	UnsuppressHashtag(ctx context.Context, tag string) error
	ListSearchSuppressions(ctx context.Context, params ListSearchSuppressionsParams) (*models.SearchSuppressionPage, error)
	CreateSearchSuppression(ctx context.Context, params CreateSearchSuppressionParams) (*models.SearchSuppression, error)
	DeleteSearchSuppression(ctx context.Context, id uuid.UUID) error
	ListAdminMedia(ctx context.Context, params ListAdminMediaParams) (*models.AdminMediaPage, error)
	FindAdminMedia(ctx context.Context, id uuid.UUID) (*models.AdminMediaAsset, error)
	RetryMedia(ctx context.Context, id uuid.UUID, adminID uuid.UUID) (*models.AdminMediaAsset, error)
	CorrectMedia(ctx context.Context, params CorrectMediaParams) (*models.AdminMediaAsset, error)
	ModerateMedia(ctx context.Context, params ModerateMediaParams) (*models.AdminMediaAsset, error)
	CleanupOrphanedMedia(ctx context.Context, olderThan time.Time, limit int) (int64, error)
	ListFeaturedSets(ctx context.Context, params ListFeaturedSetsParams) ([]models.AdminFeaturedSet, error)
	FeatureSet(ctx context.Context, setID uuid.UUID, adminID uuid.UUID, expiresAt *time.Time) (*models.AdminFeaturedSet, error)
	UnfeatureSet(ctx context.Context, setID uuid.UUID) error
	ListStoryContributions(ctx context.Context, params ListStoryContributionsParams) ([]models.AdminStoryContribution, error)
	ReviewStoryContribution(ctx context.Context, params ReviewStoryContributionParams) (*models.AdminStoryContribution, error)
	ListAdminHighlights(ctx context.Context, params ListAdminHighlightsParams) ([]models.AdminHighlight, error)
	RemoveAdminHighlight(ctx context.Context, highlightType string, highlightID uuid.UUID) error
}

func NewAdminRepository(db *pgxpool.Pool) AdminRepository {
	return newPgRepository(db)
}
