package pipes

import (
	"context"
	"strings"

	adminservices "github.com/emmanuella-codes/nox/admin/services"
	authservices "github.com/emmanuella-codes/nox/auth/services"
	"github.com/emmanuella-codes/nox/config"
	"github.com/emmanuella-codes/nox/models"
	adminrepo "github.com/emmanuella-codes/nox/repositories/admin"
	userrepo "github.com/emmanuella-codes/nox/repositories/user"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type auditContextKey string

const (
	requestIDContextKey auditContextKey = "request_id"
	ipAddressContextKey auditContextKey = "ip_address"
	userAgentContextKey auditContextKey = "user_agent"
)

type AdminPipe struct {
	adminRepo    adminrepo.AdminRepository
	userRepo     userrepo.UserRepository
	hashService  *authservices.HashService
	otpService   *authservices.OTPService
	emailService *authservices.EmailService
	tokenService *adminservices.TokenService
	db           *pgxpool.Pool
	redis        *redis.Client
	cfg          *config.Config
}

type AdminPipeDeps struct {
	AdminRepo    adminrepo.AdminRepository
	UserRepo     userrepo.UserRepository
	HashService  *authservices.HashService
	OTPService   *authservices.OTPService
	EmailService *authservices.EmailService
	TokenService *adminservices.TokenService
	DB           *pgxpool.Pool
	Redis        *redis.Client
	Config       *config.Config
}

type AuthResponse struct {
	Tokens adminservices.TokenPair `json:"tokens"`
}

type MeResponse struct {
	UserID   uuid.UUID        `json:"user_id"`
	Fullname string           `json:"fullname"`
	Email    string           `json:"email"`
	Role     models.AdminRole `json:"role"`
	IsActive bool             `json:"is_active"`
}

func NewAdminPipe(deps AdminPipeDeps) *AdminPipe {
	return &AdminPipe{
		adminRepo:    deps.AdminRepo,
		userRepo:     deps.UserRepo,
		hashService:  deps.HashService,
		otpService:   deps.OTPService,
		emailService: deps.EmailService,
		tokenService: deps.TokenService,
		db:           deps.DB,
		redis:        deps.Redis,
		cfg:          deps.Config,
	}
}

func (p *AdminPipe) issueTokenPair(ctx context.Context, userID uuid.UUID) (*adminservices.TokenPair, error) {
	tokens, err := p.tokenService.IssuePair(userID)
	if err != nil {
		return nil, err
	}

	if err := p.storeRefreshSession(ctx, userID, tokens.RefreshTokenID); err != nil {
		return nil, err
	}

	return tokens, nil
}

func (p *AdminPipe) storeRefreshSession(ctx context.Context, userID uuid.UUID, tokenID string) error {
	pipe := p.redis.TxPipeline()
	pipe.Set(ctx, refreshSessionKey(tokenID), userID.String(), p.cfg.AdminJWTRefreshTTL)
	pipe.SAdd(ctx, userRefreshSessionsKey(userID), tokenID)
	pipe.Expire(ctx, userRefreshSessionsKey(userID), p.cfg.AdminJWTRefreshTTL)
	_, err := pipe.Exec(ctx)
	return err
}

func (p *AdminPipe) deleteRefreshSession(ctx context.Context, userID uuid.UUID, tokenID string) error {
	pipe := p.redis.TxPipeline()
	pipe.Del(ctx, refreshSessionKey(tokenID))
	pipe.SRem(ctx, userRefreshSessionsKey(userID), tokenID)
	_, err := pipe.Exec(ctx)
	return err
}

func (p *AdminPipe) revokeUserRefreshSessions(ctx context.Context, userID uuid.UUID) error {
	sessionIDs, err := p.redis.SMembers(ctx, userRefreshSessionsKey(userID)).Result()
	if err != nil {
		return err
	}

	pipe := p.redis.TxPipeline()
	for _, sessionID := range sessionIDs {
		pipe.Del(ctx, refreshSessionKey(sessionID))
	}
	pipe.Del(ctx, userRefreshSessionsKey(userID))
	_, err = pipe.Exec(ctx)
	return err
}

func (p *AdminPipe) audit(ctx context.Context, adminUserID uuid.UUID, action string, meta map[string]any) {
	if err := p.adminRepo.CreateAuditLog(ctx, adminrepo.CreateAuditLogParams{
		AdminUserID: adminUserID,
		Action:      action,
		RequestID:   stringValue(ctx, requestIDContextKey),
		IPAddress:   stringValue(ctx, ipAddressContextKey),
		UserAgent:   stringValue(ctx, userAgentContextKey),
		Metadata:    meta,
	}); err != nil {
		log.Error().Err(err).Str("action", action).Msg("admin audit log failed")
	}
}

func authResponse(tokens *adminservices.TokenPair) *AuthResponse {
	return &AuthResponse{Tokens: *tokens}
}

func meResponse(identity *models.AdminIdentity) *MeResponse {
	return &MeResponse{
		UserID:   identity.UserID,
		Fullname: identity.Fullname,
		Email:    identity.Email,
		Role:     identity.Role,
		IsActive: identity.IsActive,
	}
}

func logInternalError(err error, operation string) {
	if err == nil {
		return
	}
	log.Error().Err(err).Str("operation", operation).Msg("admin auth internal error")
}

func logRefreshTokenReuse(userID uuid.UUID, tokenID string, operation string) {
	log.Warn().
		Str("user_id", userID.String()).
		Str("token_id", tokenID).
		Str("operation", operation).
		Msg("possible admin refresh token reuse detected")
}

func refreshSessionKey(tokenID string) string {
	return "admin:session:" + tokenID
}

func userRefreshSessionsKey(userID uuid.UUID) string {
	return "admin:sessions:user:" + userID.String()
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func appRefreshSessionKey(tokenID string) string {
	return "session:" + tokenID
}

func appUserRefreshSessionsKey(userID uuid.UUID) string {
	return "sessions:user:" + userID.String()
}

func appEmailVerificationKey(userID string) string {
	return "email_verify:" + userID
}

func appEmailVerificationAttemptsKey(userID string) string {
	return "email_verify_attempts:" + userID
}

func RequestIDContextKey() any {
	return requestIDContextKey
}

func IPAddressContextKey() any {
	return ipAddressContextKey
}

func UserAgentContextKey() any {
	return userAgentContextKey
}

func stringValue(ctx context.Context, key auditContextKey) string {
	value, _ := ctx.Value(key).(string)
	return value
}
