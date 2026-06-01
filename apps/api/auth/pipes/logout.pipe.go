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
		return shared.PipeError[any](messages.Invalid_Token)
	}

	claims, err := p.tokenService.VerifyRefreshToken(refreshToken)
	if err != nil {
		return shared.PipeError[any](messages.Invalid_Token)
	}

	sessionUserID, err := p.redis.Get(ctx, refreshSessionKey(claims.ID)).Result()
	if errors.Is(err, redis.Nil) {
		logRefreshTokenReuse(claims.UserID, claims.ID, "logout.missing_session")
		_ = p.revokeUserRefreshSessions(ctx, claims.UserID)
		return shared.PipeError[any](messages.Invalid_Token)
	}
	if err != nil {
		logInternalError(err, "logout.get_session")
		return shared.PipeError[any](messages.Internal_Error)
	}
	if sessionUserID != claims.UserID.String() {
		logRefreshTokenReuse(claims.UserID, claims.ID, "logout.session_user_mismatch")
		_ = p.revokeUserRefreshSessions(ctx, claims.UserID)
		return shared.PipeError[any](messages.Invalid_Token)
	}

	if err := p.deleteRefreshSession(ctx, claims.UserID, claims.ID); err != nil {
		logInternalError(err, "logout.delete_session")
		return shared.PipeError[any](messages.Internal_Error)
	}

	return shared.PipeSuccess[any](messages.User_Logged_Out, nil)
}
