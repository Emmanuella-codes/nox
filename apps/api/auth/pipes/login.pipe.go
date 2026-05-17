package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/auth/dtos"
	"github.com/emmanuella-codes/nox/auth/messages"
	"github.com/emmanuella-codes/nox/shared"
)

func (p *AuthPipe) LoginPipe(ctx context.Context, dto dtos.LoginDTO) *shared.PipeRes[AuthResponse] {
	foundUser, err := p.userRepo.FindUserByEmail(ctx, normalizeEmail(dto.Email))
	if err != nil {
		logInternalError(err, "login.find_user_by_email")
		return pipeInternalError[AuthResponse]()
	}
	if foundUser == nil {
		p.hashService.CompareDummyPassword(dto.Password)
		return pipeError[AuthResponse](messages.Invalid_Credentials)
	}
	if !p.hashService.ComparePassword(foundUser.Password, dto.Password) {
		return pipeError[AuthResponse](messages.Invalid_Credentials)
	}
	if !foundUser.EmailVerified {
		return pipeError[AuthResponse](messages.Email_Not_Verified)
	}

	tokens, err := p.issueTokenPair(ctx, foundUser.ID)
	if err != nil {
		logInternalError(err, "login.issue_token_pair")
		return pipeInternalError[AuthResponse]()
	}

	return pipeSuccess(messages.User_Logged_In, authResponse(foundUser, tokens))
}
