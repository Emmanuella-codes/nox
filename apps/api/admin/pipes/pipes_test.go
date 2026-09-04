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
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func TestLoginPipeIssuesAdminTokensForActiveMembership(t *testing.T) {
	ctx := context.Background()
	user := adminTestUser(t, "ada@example.com", "password123", true)
	adminStore := &adminTestRepo{
		membership: &models.AdminMembership{UserID: user.ID, Role: models.AdminRoleSupport, IsActive: true},
		identity:   &models.AdminIdentity{UserID: user.ID, Fullname: user.Fullname, Email: user.Email, Role: models.AdminRoleSupport, IsActive: true, Scopes: models.AdminScopesForRole(models.AdminRoleSupport)},
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
			Scopes:   models.AdminScopesForRole(models.AdminRoleOps),
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
		TokenService: adminservices.NewTokenService(cfg),
		Redis:        redisClient,
		Config:       cfg,
	})

	return pipe, redisServer
}

type adminTestUserRepo struct {
	userByEmail *models.User
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
	return nil, nil
}

func (r *adminTestUserRepo) MarkEmailVerified(ctx context.Context, userID string) error {
	return nil
}

type adminTestRepo struct {
	membership   *models.AdminMembership
	identity     *models.AdminIdentity
	auditActions []string
}

func (r *adminTestRepo) FindMembershipByUserID(ctx context.Context, userID uuid.UUID) (*models.AdminMembership, error) {
	if r.membership == nil || r.membership.UserID != userID {
		return nil, nil
	}
	return r.membership, nil
}

func (r *adminTestRepo) FindIdentityByUserID(ctx context.Context, userID uuid.UUID) (*models.AdminIdentity, error) {
	if r.identity == nil || r.identity.UserID != userID {
		return nil, nil
	}
	return r.identity, nil
}

func (r *adminTestRepo) CreateAuditLog(ctx context.Context, params adminrepo.CreateAuditLogParams) error {
	r.auditActions = append(r.auditActions, params.Action)
	return nil
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
	}
}
