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

func (p *AuthPipe) SignupPipe(ctx context.Context, dto dtos.SignupDTO) *shared.PipeRes[VerificationResponse] {
	email := normalizeEmail(dto.Email)

	userExists, err := p.userRepo.FindUserByEmail(ctx, email)
	if err != nil {
		logInternalError(err, "signup.find_user_by_email")
		return pipeInternalError[VerificationResponse]()
	}
	if userExists != nil {
		return pipeError[VerificationResponse](messages.User_Already_Exists)
	}

	passwordHash, err := p.hashService.HashPassword(dto.Password)
	if err != nil {
		logInternalError(err, "signup.hash_password")
		return pipeInternalError[VerificationResponse]()
	}

	fullname := strings.TrimSpace(dto.Firstname + " " + dto.Lastname)
	createdUser, err := p.userRepo.CreateUser(ctx, fullname, email, passwordHash)
	if err != nil {
		if errors.Is(err, user.ErrUserAlreadyExists) {
			return pipeError[VerificationResponse](messages.User_Already_Exists)
		}
		logInternalError(err, "signup.create_user")
		return pipeInternalError[VerificationResponse]()
	}

	if err := p.sendVerificationOTP(ctx, createdUser); err != nil {
		logInternalError(err, "signup.send_verification_otp")
		return pipeInternalError[VerificationResponse]()
	}

	return pipeSuccess(messages.Verification_Sent, &VerificationResponse{
		User: UserResponse{
			ID:    createdUser.ID,
			Email: createdUser.Email,
		},
		ExpiresInSeconds: int64(p.cfg.EmailOTPTTL.Seconds()),
	})
}
