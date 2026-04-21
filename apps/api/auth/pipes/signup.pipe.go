package pipes

import (
	"context"
	"errors"
	"strings"

	"github.com/emmanuella-codes/nox/auth/dtos"
	"github.com/emmanuella-codes/nox/auth/messages"
	"github.com/emmanuella-codes/nox/repositories/user"
	"github.com/emmanuella-codes/nox/shared"
)

func (p *AuthPipe) SignupPipe(ctx context.Context, dto dtos.SignupDTO) *shared.PipeRes[AuthResponse] {
	email := normalizeEmail(dto.Email)

	userExists, err := p.userRepo.FindUserByEmail(ctx, email)
	if err != nil {
		return pipeError[AuthResponse](shared.CreatePipeMessage(err.Error()))
	}
	if userExists != nil {
		return pipeError[AuthResponse](messages.User_Already_Exists)
	}

	passwordHash, err := p.hashService.HashPassword(dto.Password)
	if err != nil {
		return pipeError[AuthResponse](shared.CreatePipeMessage(err.Error()))
	}

	fullname := strings.TrimSpace(dto.Firstname + " " + dto.Lastname)
	createdUser, err := p.userRepo.CreateUser(ctx, fullname, email, passwordHash)
	if err != nil {
		if errors.Is(err, user.ErrUserAlreadyExists) {
			return pipeError[AuthResponse](messages.User_Already_Exists)
		}
		return pipeError[AuthResponse](shared.CreatePipeMessage(err.Error()))
	}

	tokens, err := p.issueTokenPair(ctx, createdUser.ID)
	if err != nil {
		return pipeError[AuthResponse](shared.CreatePipeMessage(err.Error()))
	}

	return pipeSuccess(messages.User_Created, authResponse(createdUser, tokens))
}
