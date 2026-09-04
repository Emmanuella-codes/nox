package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/admin/dtos"
	"github.com/emmanuella-codes/nox/admin/messages"
	"github.com/emmanuella-codes/nox/models"
	"github.com/emmanuella-codes/nox/shared"
)

func (p *AdminPipe) LoginPipe(ctx context.Context, dto dtos.LoginDTO) *shared.PipeRes[AuthResponse] {
	user, err := p.userRepo.FindUserByEmail(ctx, normalizeEmail(dto.Email))
	if err != nil {
		logInternalError(err, "login.find_user_by_email")
		return shared.PipeError[AuthResponse](messages.Internal_Error)
	}
	if user == nil || !user.EmailVerified {
		return shared.PipeError[AuthResponse](messages.Invalid_Credentials)
	}

	if !p.hashService.ComparePassword(user.Password, dto.Password) {
		return shared.PipeError[AuthResponse](messages.Invalid_Credentials)
	}
	if models.NormalizeUserStatus(user.Status) != models.UserStatusActive {
		return shared.PipeError[AuthResponse](messages.Admin_Access_Denied)
	}

	membership, err := p.adminRepo.FindMembershipByUserID(ctx, user.ID)
	if err != nil {
		logInternalError(err, "login.find_membership")
		return shared.PipeError[AuthResponse](messages.Internal_Error)
	}
	if membership == nil || !membership.IsActive {
		return shared.PipeError[AuthResponse](messages.Invalid_Credentials)
	}

	tokens, err := p.issueTokenPair(ctx, user.ID)
	if err != nil {
		logInternalError(err, "login.issue_token_pair")
		return shared.PipeError[AuthResponse](messages.Internal_Error)
	}

	p.audit(ctx, user.ID, "admin.auth.login", map[string]any{"role": membership.Role})
	return shared.PipeSuccess(messages.Admin_Logged_In, authResponse(tokens))
}
