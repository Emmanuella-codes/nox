package pipes

import (
	"context"
	"strings"

	"github.com/emmanuella-codes/nox/admin/messages"
	"github.com/emmanuella-codes/nox/shared"
)

func (p *AdminPipe) LogoutPipe(ctx context.Context, refreshToken string) *shared.PipeRes[any] {
	refreshToken = strings.TrimSpace(refreshToken)
	if refreshToken == "" {
		return shared.PipeError[any](messages.Invalid_Token)
	}

	claims, err := p.tokenService.VerifyRefreshToken(refreshToken)
	if err != nil {
		return shared.PipeError[any](messages.Invalid_Token)
	}

	if err := p.deleteRefreshSession(ctx, claims.UserID, claims.ID); err != nil {
		logInternalError(err, "logout.delete_refresh_session")
		return shared.PipeError[any](messages.Internal_Error)
	}

	p.audit(ctx, claims.UserID, "admin.auth.logout", nil)
	return shared.PipeSuccess[any](messages.Admin_Logged_Out, nil)
}
