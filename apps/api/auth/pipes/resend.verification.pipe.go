package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/auth/dtos"
	"github.com/emmanuella-codes/nox/auth/messages"
	"github.com/emmanuella-codes/nox/shared"
)

func (p *AuthPipe) ResendVerificationPipe(ctx context.Context, dto dtos.ResendVerificationDTO) *shared.PipeRes[VerificationResponse] {
	foundUser, err := p.userRepo.FindUserByEmail(ctx, normalizeEmail(dto.Email))
	if err != nil {
		logInternalError(err, "resend_verification.find_user_by_email")
		return pipeInternalError[VerificationResponse]()
	}
	if foundUser == nil {
		return pipeError[VerificationResponse](messages.Invalid_Credentials)
	}
	if foundUser.EmailVerified {
		return pipeError[VerificationResponse](messages.Email_Already_Verified)
	}

	if err := p.sendVerificationOTP(ctx, foundUser); err != nil {
		logInternalError(err, "resend_verification.send_otp")
		return pipeInternalError[VerificationResponse]()
	}

	return pipeSuccess(messages.Verification_Sent, &VerificationResponse{
		User: UserResponse{
			ID:    foundUser.ID,
			Email: foundUser.Email,
		},
		ExpiresInSeconds: int64(p.cfg.EmailOTPTTL.Seconds()),
	})
}
