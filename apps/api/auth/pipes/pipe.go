package pipes

import (
	"context"
	"strings"

	"github.com/emmanuella-codes/nox/auth/services"
	"github.com/emmanuella-codes/nox/config"
	"github.com/emmanuella-codes/nox/models"
	"github.com/emmanuella-codes/nox/repositories/user"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type AuthPipe struct {
	userRepo     user.UserRepository
	hashService  *services.HashService
	tokenService *services.TokenService
	redis        *redis.Client
	cfg          *config.Config
}

type AuthPipeDeps struct {
	UserRepo     user.UserRepository
	HashService  *services.HashService
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
		tokenService: deps.TokenService,
		redis:        deps.Redis,
		cfg:          deps.Config,
	}
}

func (p *AuthPipe) issueTokenPair(ctx context.Context, userID uuid.UUID) (*services.TokenPair, error) {
	tokens, err := p.tokenService.IssuePair(userID)
	if err != nil {
		return nil, err
	}

	if err := p.redis.Set(ctx, refreshSessionKey(tokens.RefreshTokenID), userID.String(), p.cfg.JWTRefreshTTL).Err(); err != nil {
		return nil, err
	}

	return tokens, nil
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

func refreshSessionKey(tokenID string) string {
	return "session:" + tokenID
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
