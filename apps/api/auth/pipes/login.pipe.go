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
		return pipeError[AuthResponse](shared.CreatePipeMessage(err.Error()))
	}
	if foundUser == nil || !p.hashService.ComparePassword(foundUser.Password, dto.Password) {
		return pipeError[AuthResponse](messages.Invalid_Credentials)
	}

	tokens, err := p.issueTokenPair(ctx, foundUser.ID)
	if err != nil {
		return pipeError[AuthResponse](shared.CreatePipeMessage(err.Error()))
	}

	return pipeSuccess(messages.User_Logged_In, authResponse(foundUser, tokens))
}
