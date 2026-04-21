package pipes

import (
	"context"
	"errors"
	"strings"

	"github.com/emmanuella-codes/nox/auth/messages"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/redis/go-redis/v9"
)

func (p *AuthPipe) LogoutPipe(ctx context.Context, refreshToken string) *shared.PipeRes[any] {
	if strings.TrimSpace(refreshToken) == "" {
		return pipeError[any](messages.Invalid_Token)
	}

	claims, err := p.tokenService.VerifyRefreshToken(refreshToken)
	if err != nil {
		return pipeError[any](messages.Invalid_Token)
	}

	sessionUserID, err := p.redis.Get(ctx, refreshSessionKey(claims.ID)).Result()
	if errors.Is(err, redis.Nil) {
		return pipeError[any](messages.Invalid_Token)
	}
	if err != nil {
		return pipeError[any](shared.CreatePipeMessage(err.Error()))
	}
	if sessionUserID != claims.UserID.String() {
		return pipeError[any](messages.Invalid_Token)
	}

	if err := p.redis.Del(ctx, refreshSessionKey(claims.ID)).Err(); err != nil {
		return pipeError[any](shared.CreatePipeMessage(err.Error()))
	}

	return pipeSuccess[any](messages.User_Logged_Out, nil)
}
