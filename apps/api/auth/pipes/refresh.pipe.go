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
		logRefreshTokenReuse(claims.UserID, claims.ID, "refresh.missing_session")
		_ = p.revokeUserRefreshSessions(ctx, claims.UserID)
		return pipeError[services.TokenPair](messages.Invalid_Token)
	}
	if err != nil {
		logInternalError(err, "refresh.get_session")
		return pipeInternalError[services.TokenPair]()
	}
	if sessionUserID != claims.UserID.String() {
		logRefreshTokenReuse(claims.UserID, claims.ID, "refresh.session_user_mismatch")
		_ = p.revokeUserRefreshSessions(ctx, claims.UserID)
		return pipeError[services.TokenPair](messages.Invalid_Token)
	}

	if err := p.deleteRefreshSession(ctx, claims.UserID, claims.ID); err != nil {
		logInternalError(err, "refresh.delete_session")
		return pipeInternalError[services.TokenPair]()
	}

	tokens, err := p.issueTokenPair(ctx, claims.UserID)
	if err != nil {
		logInternalError(err, "refresh.issue_token_pair")
		return pipeInternalError[services.TokenPair]()
	}

	return pipeSuccess(messages.Token_Refreshed, tokens)
}
