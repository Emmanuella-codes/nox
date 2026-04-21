package pipes

import (
	"context"
	"errors"

	"github.com/emmanuella-codes/nox/auth/messages"
	"github.com/emmanuella-codes/nox/auth/services"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/redis/go-redis/v9"
)

func (p *AuthPipe) RefreshPipe(ctx context.Context, refreshToken string) *shared.PipeRes[services.TokenPair] {
	claims, err := p.tokenService.VerifyRefreshToken(refreshToken)
	if err != nil {
		return pipeError[services.TokenPair](messages.Invalid_Token)
	}

	sessionUserID, err := p.redis.Get(ctx, refreshSessionKey(claims.ID)).Result()
	if errors.Is(err, redis.Nil) {
		return pipeError[services.TokenPair](messages.Invalid_Token)
	}
	if err != nil {
		return pipeError[services.TokenPair](shared.CreatePipeMessage(err.Error()))
	}
	if sessionUserID != claims.UserID.String() {
		return pipeError[services.TokenPair](messages.Invalid_Token)
	}

	if err := p.redis.Del(ctx, refreshSessionKey(claims.ID)).Err(); err != nil {
		return pipeError[services.TokenPair](shared.CreatePipeMessage(err.Error()))
	}

	tokens, err := p.issueTokenPair(ctx, claims.UserID)
	if err != nil {
		return pipeError[services.TokenPair](shared.CreatePipeMessage(err.Error()))
	}

	return pipeSuccess(messages.Token_Refreshed, tokens)
}
