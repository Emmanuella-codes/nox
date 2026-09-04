package pipes

import (
	"context"
	"strings"

	"github.com/emmanuella-codes/nox/admin/messages"
	"github.com/emmanuella-codes/nox/shared"
)

func (p *AdminPipe) RefreshPipe(ctx context.Context, refreshToken string) *shared.PipeRes[AuthResponse] {
	refreshToken = strings.TrimSpace(refreshToken)
	if refreshToken == "" {
		return shared.PipeError[AuthResponse](messages.Invalid_Token)
	}

	claims, err := p.tokenService.VerifyRefreshToken(refreshToken)
	if err != nil {
		return shared.PipeError[AuthResponse](messages.Invalid_Token)
	}

	sessionUserID, err := p.redis.Get(ctx, refreshSessionKey(claims.ID)).Result()
	if err != nil || sessionUserID != claims.UserID.String() {
		if revokeErr := p.revokeUserRefreshSessions(ctx, claims.UserID); revokeErr != nil {
			logInternalError(revokeErr, "refresh.revoke_refresh_sessions")
		}
		logRefreshTokenReuse(claims.UserID, claims.ID, "refresh")
		return shared.PipeError[AuthResponse](messages.Invalid_Token)
	}

	membership, err := p.adminRepo.FindMembershipByUserID(ctx, claims.UserID)
	if err != nil {
		logInternalError(err, "refresh.find_membership")
		return shared.PipeError[AuthResponse](messages.Internal_Error)
	}
	if membership == nil || !membership.IsActive {
		return shared.PipeError[AuthResponse](messages.Admin_Access_Denied)
	}

	if err := p.deleteRefreshSession(ctx, claims.UserID, claims.ID); err != nil {
		logInternalError(err, "refresh.delete_session")
		return shared.PipeError[AuthResponse](messages.Internal_Error)
	}

	tokens, err := p.issueTokenPair(ctx, claims.UserID)
	if err != nil {
		logInternalError(err, "refresh.issue_token_pair")
		return shared.PipeError[AuthResponse](messages.Internal_Error)
	}

	p.audit(ctx, claims.UserID, "admin.auth.refresh", map[string]any{"role": membership.Role})
	return shared.PipeSuccess(messages.Token_Refreshed, authResponse(tokens))
}
