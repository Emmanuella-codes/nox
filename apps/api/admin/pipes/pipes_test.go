package pipes

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	admindtos "github.com/emmanuella-codes/nox/admin/dtos"
	adminmessages "github.com/emmanuella-codes/nox/admin/messages"
	adminservices "github.com/emmanuella-codes/nox/admin/services"
	authservices "github.com/emmanuella-codes/nox/auth/services"
	"github.com/emmanuella-codes/nox/config"
	"github.com/emmanuella-codes/nox/models"
	adminrepo "github.com/emmanuella-codes/nox/repositories/admin"
	"github.com/emmanuella-codes/nox/shared/mail"
	workerruntime "github.com/emmanuella-codes/nox/workers/runtime"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func TestLoginPipeIssuesAdminTokensForActiveMembership(t *testing.T) {
	ctx := context.Background()
	user := adminTestUser(t, "ada@example.com", "password123", true)
	adminStore := &adminTestRepo{
		membership: &models.AdminMembership{UserID: user.ID, Role: models.AdminRoleSupport, IsActive: true},
		identity:   &models.AdminIdentity{UserID: user.ID, Fullname: user.Fullname, Email: user.Email, Role: models.AdminRoleSupport, IsActive: true},
	}
	pipe, _ := newAdminTestPipe(t, &adminTestUserRepo{userByEmail: user}, adminStore)

	res := pipe.LoginPipe(ctx, admindtos.LoginDTO{Email: "ADA@example.com", Password: "password123"})
	if !res.Success {
		t.Fatalf("expected login success, got %q", res.Message)
	}
	if res.Data == nil || res.Data.Tokens.RefreshTokenID == "" {
		t.Fatal("expected token pair")
	}
	if pipe.redis.Exists(ctx, refreshSessionKey(res.Data.Tokens.RefreshTokenID)).Val() != 1 {
		t.Fatal("expected admin refresh session to be stored")
	}
	if len(adminStore.auditActions) != 1 || adminStore.auditActions[0] != "admin.auth.login" {
		t.Fatalf("expected login audit log, got %#v", adminStore.auditActions)
	}
}

func TestLoginPipeRejectsNonAdminUser(t *testing.T) {
	user := adminTestUser(t, "ada@example.com", "password123", true)
	pipe, _ := newAdminTestPipe(t, &adminTestUserRepo{userByEmail: user}, &adminTestRepo{})

	res := pipe.LoginPipe(context.Background(), admindtos.LoginDTO{Email: "ada@example.com", Password: "password123"})
	if res.Message != adminmessages.Invalid_Credentials {
		t.Fatalf("expected invalid credentials, got %q", res.Message)
	}
}

func TestRefreshPipeRejectsInactiveMembership(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	adminStore := &adminTestRepo{
		membership: &models.AdminMembership{UserID: userID, Role: models.AdminRoleModerator, IsActive: false},
	}
	pipe, _ := newAdminTestPipe(t, &adminTestUserRepo{}, adminStore)
	tokenID := "inactive-admin-refresh"
	refreshToken, err := pipe.tokenService.SignRefreshToken(userID, tokenID)
	if err != nil {
		t.Fatalf("sign refresh token: %v", err)
	}
	if err := pipe.storeRefreshSession(ctx, userID, tokenID); err != nil {
		t.Fatalf("store refresh session: %v", err)
	}

	res := pipe.RefreshPipe(ctx, refreshToken)
	if res.Message != adminmessages.Admin_Access_Denied {
		t.Fatalf("expected access denied, got %q", res.Message)
	}
}

func TestMePipeReturnsAdminIdentity(t *testing.T) {
	userID := uuid.New()
	adminStore := &adminTestRepo{
		identity: &models.AdminIdentity{
			UserID:   userID,
			Fullname: "Ada Lovelace",
			Email:    "ada@example.com",
			Role:     models.AdminRoleOps,
			IsActive: true,
		},
	}
	pipe, _ := newAdminTestPipe(t, &adminTestUserRepo{}, adminStore)

	res := pipe.MePipe(context.Background(), userID)
	if !res.Success {
		t.Fatalf("expected me success, got %q", res.Message)
	}
	if res.Data == nil || res.Data.Role != models.AdminRoleOps {
		t.Fatal("expected admin identity")
	}
}

func TestHealthPipeReturnsOfflineWorkersWithoutHeartbeats(t *testing.T) {
	pipe, _ := newAdminTestPipe(t, &adminTestUserRepo{}, &adminTestRepo{})

	res := pipe.HealthPipe(context.Background())
	if !res.Success {
		t.Fatalf("expected health success, got %q", res.Message)
	}
	if res.Data == nil || res.Data.Redis.Status != "ok" {
		t.Fatal("expected redis health to be ok")
	}
	if res.Data.Workers[workerruntime.WorkerMediaKey].Status != "offline" {
		t.Fatal("expected media worker to report offline without heartbeat")
	}
}

func TestListAdminUsersPipeRequiresSuperAdmin(t *testing.T) {
	actorID := uuid.New()
	pipe, _ := newAdminTestPipe(t, &adminTestUserRepo{}, &adminTestRepo{
		membership: &models.AdminMembership{UserID: actorID, Role: models.AdminRoleSupport, IsActive: true},
	})

	res := pipe.ListAdminUsersPipe(context.Background(), actorID)
	if res.Message != adminmessages.Admin_Access_Denied {
		t.Fatalf("expected access denied, got %q", res.Message)
	}
}

func TestCreateAdminUserPipeCreatesMembership(t *testing.T) {
	actorID := uuid.New()
	targetUserID := uuid.New()
	pipe, _ := newAdminTestPipe(t, &adminTestUserRepo{
		userByID: &models.User{ID: targetUserID, Fullname: "Target User", Email: "target@example.com"},
	}, &adminTestRepo{
		membership:         &models.AdminMembership{UserID: actorID, Role: models.AdminRoleSuperAdmin, IsActive: true},
		membershipByUserID: map[uuid.UUID]*models.AdminMembership{},
		identityByUserID:   map[uuid.UUID]*models.AdminIdentity{},
	})

	res := pipe.CreateAdminUserPipe(context.Background(), actorID, admindtos.CreateAdminUserDTO{
		UserID: targetUserID.String(),
		Role:   models.AdminRoleSupport,
	})
	if !res.Success {
		t.Fatalf("expected create success, got %q", res.Message)
	}
	if res.Data == nil || res.Data.UserID != targetUserID || res.Data.Role != models.AdminRoleSupport {
		t.Fatal("expected created admin identity")
	}
}

func TestUpdateAdminUserRolePipeProtectsLastSuperAdmin(t *testing.T) {
	actorID := uuid.New()
	targetUserID := uuid.New()
	pipe, _ := newAdminTestPipe(t, &adminTestUserRepo{}, &adminTestRepo{
		membership: &models.AdminMembership{UserID: actorID, Role: models.AdminRoleSuperAdmin, IsActive: true},
		membershipByUserID: map[uuid.UUID]*models.AdminMembership{
			targetUserID: {UserID: targetUserID, Role: models.AdminRoleSuperAdmin, IsActive: true},
		},
		identityByUserID: map[uuid.UUID]*models.AdminIdentity{
			targetUserID: {UserID: targetUserID, Fullname: "Target User", Email: "target@example.com", Role: models.AdminRoleSuperAdmin, IsActive: true},
		},
		activeCountByRole: map[models.AdminRole]int{models.AdminRoleSuperAdmin: 1},
	})

	res := pipe.UpdateAdminUserRolePipe(context.Background(), actorID, targetUserID, admindtos.UpdateAdminRoleDTO{
		Role: models.AdminRoleSupport,
	})
	if res.Message != adminmessages.Last_Super_Admin_Required {
		t.Fatalf("expected last super admin protection, got %q", res.Message)
	}
}

func TestUpdateAdminUserStatusPipeDeactivatesNonFinalSuperAdmin(t *testing.T) {
	actorID := uuid.New()
	targetUserID := uuid.New()
	pipe, _ := newAdminTestPipe(t, &adminTestUserRepo{}, &adminTestRepo{
		membership: &models.AdminMembership{UserID: actorID, Role: models.AdminRoleSuperAdmin, IsActive: true},
		membershipByUserID: map[uuid.UUID]*models.AdminMembership{
			targetUserID: {UserID: targetUserID, Role: models.AdminRoleSuperAdmin, IsActive: true},
		},
		identityByUserID: map[uuid.UUID]*models.AdminIdentity{
			targetUserID: {UserID: targetUserID, Fullname: "Target User", Email: "target@example.com", Role: models.AdminRoleSuperAdmin, IsActive: true},
		},
		activeCountByRole: map[models.AdminRole]int{models.AdminRoleSuperAdmin: 2},
	})

	res := pipe.UpdateAdminUserStatusPipe(context.Background(), actorID, targetUserID, admindtos.UpdateAdminStatusDTO{
		IsActive: false,
	})
	if !res.Success {
		t.Fatalf("expected status update success, got %q", res.Message)
	}
	if res.Data == nil || res.Data.IsActive {
		t.Fatal("expected admin user to be inactive")
	}
}

func TestGetUserPipeReturnsManagedUser(t *testing.T) {
	actorID := uuid.New()
	targetUserID := uuid.New()
	pipe, _ := newAdminTestPipe(t, &adminTestUserRepo{}, &adminTestRepo{
		membership: &models.AdminMembership{UserID: actorID, Role: models.AdminRoleSuperAdmin, IsActive: true},
		managedUsersByID: map[uuid.UUID]*models.AdminManagedUser{
			targetUserID: {
				ID:            targetUserID,
				Fullname:      "Ada Lovelace",
				Email:         "ada@example.com",
				EmailVerified: true,
				Status:        models.UserStatusActive,
			},
		},
	})

	res := pipe.GetUserPipe(context.Background(), actorID, targetUserID)
	if !res.Success {
		t.Fatalf("expected user lookup success, got %q", res.Message)
	}
	if res.Data == nil || res.Data.UserID != targetUserID {
		t.Fatal("expected managed user response")
	}
}

func TestUpdateUserStatusPipeSuspendsAndRevokesSessions(t *testing.T) {
	ctx := context.Background()
	actorID := uuid.New()
	targetUserID := uuid.New()
	pipe, redisServer := newAdminTestPipe(t, &adminTestUserRepo{}, &adminTestRepo{
		membership: &models.AdminMembership{UserID: actorID, Role: models.AdminRoleSuperAdmin, IsActive: true},
		managedUsersByID: map[uuid.UUID]*models.AdminManagedUser{
			targetUserID: {
				ID:            targetUserID,
				Fullname:      "Target User",
				Email:         "target@example.com",
				EmailVerified: true,
				Status:        models.UserStatusActive,
			},
		},
	})
	redisServer.Set(appRefreshSessionKey("user-session-1"), targetUserID.String())
	redisServer.SAdd(appUserRefreshSessionsKey(targetUserID), "user-session-1")

	res := pipe.UpdateUserStatusPipe(ctx, actorID, targetUserID, models.UserStatusSuspended)
	if !res.Success {
		t.Fatalf("expected status update success, got %q", res.Message)
	}
	if res.Data == nil || res.Data.Status != models.UserStatusSuspended {
		t.Fatal("expected suspended status")
	}
	if pipe.redis.Exists(ctx, appRefreshSessionKey("user-session-1")).Val() != 0 {
		t.Fatal("expected user sessions to be revoked")
	}
}

func TestRevokeUserSessionsPipeRevokesRegularAndAdminSessions(t *testing.T) {
	ctx := context.Background()
	actorID := uuid.New()
	targetUserID := uuid.New()
	pipe, redisServer := newAdminTestPipe(t, &adminTestUserRepo{}, &adminTestRepo{
		membership: &models.AdminMembership{UserID: actorID, Role: models.AdminRoleSuperAdmin, IsActive: true},
		membershipByUserID: map[uuid.UUID]*models.AdminMembership{
			targetUserID: {UserID: targetUserID, Role: models.AdminRoleSupport, IsActive: true},
		},
		managedUsersByID: map[uuid.UUID]*models.AdminManagedUser{
			targetUserID: {
				ID:            targetUserID,
				Fullname:      "Target User",
				Email:         "target@example.com",
				EmailVerified: true,
				Status:        models.UserStatusActive,
			},
		},
	})
	redisServer.Set(appRefreshSessionKey("user-session-1"), targetUserID.String())
	redisServer.SAdd(appUserRefreshSessionsKey(targetUserID), "user-session-1")
	redisServer.Set(refreshSessionKey("admin-session-1"), targetUserID.String())
	redisServer.SAdd(userRefreshSessionsKey(targetUserID), "admin-session-1")

	res := pipe.RevokeUserSessionsPipe(ctx, actorID, targetUserID)
	if !res.Success {
		t.Fatalf("expected revoke success, got %q", res.Message)
	}
	if pipe.redis.Exists(ctx, appRefreshSessionKey("user-session-1"), refreshSessionKey("admin-session-1")).Val() != 0 {
		t.Fatal("expected all sessions to be revoked")
	}
}

func TestMarkUserEmailVerifiedPipeMarksUserVerified(t *testing.T) {
	ctx := context.Background()
	actorID := uuid.New()
	targetUserID := uuid.New()
	userRepo := &adminTestUserRepo{
		userByID: &models.User{
			ID:            targetUserID,
			Fullname:      "Target User",
			Email:         "target@example.com",
			EmailVerified: false,
			Status:        models.UserStatusActive,
		},
	}
	pipe, redisServer := newAdminTestPipe(t, userRepo, &adminTestRepo{
		membership: &models.AdminMembership{UserID: actorID, Role: models.AdminRoleSuperAdmin, IsActive: true},
		managedUsersByID: map[uuid.UUID]*models.AdminManagedUser{
			targetUserID: {
				ID:            targetUserID,
				Fullname:      "Target User",
				Email:         "target@example.com",
				EmailVerified: true,
				Status:        models.UserStatusActive,
			},
		},
	})
	redisServer.Set(appEmailVerificationKey(targetUserID.String()), "otp")
	redisServer.Set(appEmailVerificationAttemptsKey(targetUserID.String()), "2")

	res := pipe.MarkUserEmailVerifiedPipe(ctx, actorID, targetUserID)
	if !res.Success {
		t.Fatalf("expected mark verified success, got %q", res.Message)
	}
	if userRepo.markedVerifiedUserID != targetUserID.String() {
		t.Fatal("expected user repo to mark email verified")
	}
	if pipe.redis.Exists(ctx, appEmailVerificationKey(targetUserID.String()), appEmailVerificationAttemptsKey(targetUserID.String())).Val() != 0 {
		t.Fatal("expected verification state to be cleared")
	}
}

func newAdminTestPipe(t *testing.T, userStore *adminTestUserRepo, adminStore *adminTestRepo) (*AdminPipe, *miniredis.Miniredis) {
	t.Helper()

	redisServer, err := miniredis.Run()
	if err != nil {
		t.Fatalf("start miniredis: %v", err)
	}

	redisClient := redis.NewClient(&redis.Options{Addr: redisServer.Addr()})
	cfg := &config.Config{
		AdminJWTAccessSecret:  "admin-access-secret",
		AdminJWTRefreshSecret: "admin-refresh-secret",
		AdminJWTIssuer:        "nox-api",
		AdminJWTAudience:      "nox-admin",
		AdminJWTAccessTTL:     time.Minute,
		AdminJWTRefreshTTL:    time.Hour,
	}

	pipe := NewAdminPipe(AdminPipeDeps{
		AdminRepo:    adminStore,
		UserRepo:     userStore,
		HashService:  authservices.NewHashService(),
		OTPService:   authservices.NewOTPService(),
		EmailService: authservices.NewEmailService(&adminTestMailProvider{}),
		TokenService: adminservices.NewTokenService(cfg),
		Redis:        redisClient,
		Config:       cfg,
	})

	return pipe, redisServer
}

type adminTestUserRepo struct {
	userByEmail          *models.User
	userByID             *models.User
	markedVerifiedUserID string
}

func (r *adminTestUserRepo) CreateUser(ctx context.Context, fullname string, email string, passwordHash string) (*models.User, error) {
	return nil, errors.New("not implemented")
}

func (r *adminTestUserRepo) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	if r.userByEmail != nil && r.userByEmail.Email == email {
		return r.userByEmail, nil
	}
	return nil, nil
}

func (r *adminTestUserRepo) FindUserByID(ctx context.Context, userID string) (*models.User, error) {
	if r.userByID != nil && r.userByID.ID.String() == userID {
		return r.userByID, nil
	}
	return nil, nil
}

func (r *adminTestUserRepo) MarkEmailVerified(ctx context.Context, userID string) error {
	r.markedVerifiedUserID = userID
	return nil
}

type adminTestRepo struct {
	membership         *models.AdminMembership
	membershipByUserID map[uuid.UUID]*models.AdminMembership
	identity           *models.AdminIdentity
	identityByUserID   map[uuid.UUID]*models.AdminIdentity
	identities         []models.AdminIdentity
	managedUsersByID   map[uuid.UUID]*models.AdminManagedUser
	activeCountByRole  map[models.AdminRole]int
	auditActions       []string
}

func (r *adminTestRepo) ListIdentities(ctx context.Context) ([]models.AdminIdentity, error) {
	return r.identities, nil
}

func (r *adminTestRepo) FindMembershipByUserID(ctx context.Context, userID uuid.UUID) (*models.AdminMembership, error) {
	if membership, ok := r.membershipByUserID[userID]; ok {
		return membership, nil
	}
	if r.membership == nil || r.membership.UserID != userID {
		return nil, nil
	}
	return r.membership, nil
}

func (r *adminTestRepo) FindIdentityByUserID(ctx context.Context, userID uuid.UUID) (*models.AdminIdentity, error) {
	if identity, ok := r.identityByUserID[userID]; ok {
		return identity, nil
	}
	if r.identity == nil || r.identity.UserID != userID {
		return nil, nil
	}
	return r.identity, nil
}

func (r *adminTestRepo) CreateMembership(ctx context.Context, userID uuid.UUID, role models.AdminRole) (*models.AdminIdentity, error) {
	if _, ok := r.membershipByUserID[userID]; ok {
		return nil, adminrepo.ErrMembershipAlreadyExists
	}
	if r.membershipByUserID == nil {
		r.membershipByUserID = map[uuid.UUID]*models.AdminMembership{}
	}
	if r.identityByUserID == nil {
		r.identityByUserID = map[uuid.UUID]*models.AdminIdentity{}
	}

	r.membershipByUserID[userID] = &models.AdminMembership{UserID: userID, Role: role, IsActive: true}
	r.identityByUserID[userID] = &models.AdminIdentity{UserID: userID, Fullname: "Target User", Email: "target@example.com", Role: role, IsActive: true}
	return r.identityByUserID[userID], nil
}

func (r *adminTestRepo) UpdateMembershipRole(ctx context.Context, userID uuid.UUID, role models.AdminRole) (*models.AdminIdentity, error) {
	membership, ok := r.membershipByUserID[userID]
	if !ok {
		return nil, adminrepo.ErrMembershipNotFound
	}
	identity, ok := r.identityByUserID[userID]
	if !ok {
		return nil, adminrepo.ErrMembershipNotFound
	}
	membership.Role = role
	identity.Role = role
	return identity, nil
}

func (r *adminTestRepo) UpdateMembershipStatus(ctx context.Context, userID uuid.UUID, isActive bool) (*models.AdminIdentity, error) {
	membership, ok := r.membershipByUserID[userID]
	if !ok {
		return nil, adminrepo.ErrMembershipNotFound
	}
	identity, ok := r.identityByUserID[userID]
	if !ok {
		return nil, adminrepo.ErrMembershipNotFound
	}
	membership.IsActive = isActive
	identity.IsActive = isActive
	return identity, nil
}

func (r *adminTestRepo) CountActiveMembershipsByRole(ctx context.Context, role models.AdminRole) (int, error) {
	if r.activeCountByRole == nil {
		return 0, nil
	}
	return r.activeCountByRole[role], nil
}

func (r *adminTestRepo) ListManagedUsers(ctx context.Context) ([]models.AdminManagedUser, error) {
	users := []models.AdminManagedUser{}
	for _, user := range r.managedUsersByID {
		users = append(users, *user)
	}
	return users, nil
}

func (r *adminTestRepo) FindManagedUserByID(ctx context.Context, userID uuid.UUID) (*models.AdminManagedUser, error) {
	if user, ok := r.managedUsersByID[userID]; ok {
		return user, nil
	}
	return nil, nil
}

func (r *adminTestRepo) UpdateUserStatus(ctx context.Context, userID uuid.UUID, status models.UserStatus) (*models.AdminManagedUser, error) {
	user, ok := r.managedUsersByID[userID]
	if !ok {
		return nil, adminrepo.ErrUserNotFound
	}
	user.Status = status
	return user, nil
}

func (r *adminTestRepo) CreateAuditLog(ctx context.Context, params adminrepo.CreateAuditLogParams) error {
	r.auditActions = append(r.auditActions, params.Action)
	return nil
}

func (r *adminTestRepo) Moderate(ctx context.Context, params adminrepo.ModerateParams) (*models.ModerationState, error) {
	return &models.ModerationState{EntityType: params.EntityType, EntityID: params.EntityID, Status: params.Status, Reason: params.Reason, ModeratedBy: &params.AdminID}, nil
}

func (r *adminTestRepo) ListModerationActions(ctx context.Context, params adminrepo.ListModerationActionsParams) ([]models.ModerationAction, error) {
	return []models.ModerationAction{}, nil
}

func (r *adminTestRepo) CreateReport(ctx context.Context, params adminrepo.CreateReportParams) (*models.Report, error) {
	return nil, nil
}

func (r *adminTestRepo) ListReports(ctx context.Context, params adminrepo.ListReportsParams) ([]models.Report, error) {
	return []models.Report{}, nil
}

func (r *adminTestRepo) FindReportByID(ctx context.Context, reportID uuid.UUID) (*models.Report, error) {
	return nil, adminrepo.ErrReportNotFound
}

func (r *adminTestRepo) UpdateReport(ctx context.Context, params adminrepo.UpdateReportParams) (*models.Report, error) {
	return nil, adminrepo.ErrReportNotFound
}

func (r *adminTestRepo) ResolveReport(ctx context.Context, params adminrepo.ResolveReportParams) (*models.Report, error) {
	return nil, adminrepo.ErrReportNotFound
}

func (r *adminTestRepo) Dashboard(ctx context.Context) (*models.AdminDashboard, error) {
	return &models.AdminDashboard{}, nil
}

func adminTestUser(t *testing.T, email string, password string, verified bool) *models.User {
	t.Helper()

	hashService := authservices.NewHashService()
	passwordHash, err := hashService.HashPassword(password)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}

	return &models.User{
		ID:            uuid.New(),
		Fullname:      "Ada Lovelace",
		Email:         email,
		Password:      passwordHash,
		EmailVerified: verified,
		Status:        models.UserStatusActive,
	}
}

type adminTestMailProvider struct{}

func (p *adminTestMailProvider) Send(ctx context.Context, message mail.Message) error {
	return nil
}
