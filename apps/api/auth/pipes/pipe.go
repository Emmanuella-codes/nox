package pipes

import (
	"context"
	"strings"

	"github.com/emmanuella-codes/nox/auth/messages"
	"github.com/emmanuella-codes/nox/auth/services"
	"github.com/emmanuella-codes/nox/config"
	"github.com/emmanuella-codes/nox/models"
	"github.com/emmanuella-codes/nox/repositories/user"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type AuthPipe struct {
	userRepo     user.UserRepository
	hashService  *services.HashService
	otpService   *services.OTPService
	emailService *services.EmailService
	tokenService *services.TokenService
	redis        *redis.Client
	cfg          *config.Config
}

type AuthPipeDeps struct {
	UserRepo     user.UserRepository
	HashService  *services.HashService
	OTPService   *services.OTPService
	EmailService *services.EmailService
	TokenService *services.TokenService
	Redis        *redis.Client
	Config       *config.Config
}

type AuthResponse struct {
	User   UserResponse       `json:"user"`
	Tokens services.TokenPair `json:"tokens"`
}

type UserResponse struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
}

func NewAuthPipe(deps AuthPipeDeps) *AuthPipe {
	return &AuthPipe{
		userRepo:     deps.UserRepo,
		hashService:  deps.HashService,
		otpService:   deps.OTPService,
		emailService: deps.EmailService,
		tokenService: deps.TokenService,
		redis:        deps.Redis,
		cfg:          deps.Config,
	}
}

type VerificationResponse struct {
	User             UserResponse `json:"user"`
	ExpiresInSeconds int64        `json:"expires_in_seconds"`
}

func (p *AuthPipe) issueTokenPair(ctx context.Context, userID uuid.UUID) (*services.TokenPair, error) {
	tokens, err := p.tokenService.IssuePair(userID)
	if err != nil {
		return nil, err
	}

	if err := p.storeRefreshSession(ctx, userID, tokens.RefreshTokenID); err != nil {
		return nil, err
	}

	return tokens, nil
}

func (p *AuthPipe) storeRefreshSession(ctx context.Context, userID uuid.UUID, tokenID string) error {
	pipe := p.redis.TxPipeline()
	pipe.Set(ctx, refreshSessionKey(tokenID), userID.String(), p.cfg.JWTRefreshTTL)
	pipe.SAdd(ctx, userRefreshSessionsKey(userID), tokenID)
	pipe.Expire(ctx, userRefreshSessionsKey(userID), p.cfg.JWTRefreshTTL)
	_, err := pipe.Exec(ctx)
	return err
}

func (p *AuthPipe) deleteRefreshSession(ctx context.Context, userID uuid.UUID, tokenID string) error {
	pipe := p.redis.TxPipeline()
	pipe.Del(ctx, refreshSessionKey(tokenID))
	pipe.SRem(ctx, userRefreshSessionsKey(userID), tokenID)
	_, err := pipe.Exec(ctx)
	return err
}

func (p *AuthPipe) revokeUserRefreshSessions(ctx context.Context, userID uuid.UUID) error {
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

func authResponse(user *models.User, tokens *services.TokenPair) *AuthResponse {
	return &AuthResponse{
		User: UserResponse{
			ID:    user.ID,
			Email: user.Email,
		},
		Tokens: *tokens,
	}
}

func pipeSuccess[T any](message shared.PipeMessage, data *T) *shared.PipeRes[T] {
	return &shared.PipeRes[T]{
		Success: true,
		Message: message,
		Data:    data,
	}
}

func pipeError[T any](message shared.PipeMessage) *shared.PipeRes[T] {
	return &shared.PipeRes[T]{
		Success: false,
		Message: message,
	}
}

func pipeInternalError[T any]() *shared.PipeRes[T] {
	return pipeError[T](messages.Internal_Error)
}

func (p *AuthPipe) sendVerificationOTP(ctx context.Context, user *models.User) error {
	otp, err := p.otpService.Generate()
	if err != nil {
		return err
	}

	otpHash, err := p.otpService.Hash(otp)
	if err != nil {
		return err
	}

	if err := p.redis.Set(ctx, emailVerificationKey(user.ID.String()), otpHash, p.cfg.EmailOTPTTL).Err(); err != nil {
		return err
	}

	return p.emailService.SendVerificationOTP(user.Email, otp)
}

func logInternalError(err error, operation string) {
	if err == nil {
		return
	}
	log.Error().Err(err).Str("operation", operation).Msg("auth internal error")
}

func logRefreshTokenReuse(userID uuid.UUID, tokenID string, operation string) {
	log.Warn().
		Str("user_id", userID.String()).
		Str("token_id", tokenID).
		Str("operation", operation).
		Msg("possible refresh token reuse detected")
}

func refreshSessionKey(tokenID string) string {
	return "session:" + tokenID
}

func userRefreshSessionsKey(userID uuid.UUID) string {
	return "sessions:user:" + userID.String()
}

func emailVerificationKey(userID string) string {
	return "email_verify:" + userID
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
