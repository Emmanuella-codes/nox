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
		return shared.PipeError[VerificationResponse](messages.Internal_Error)
	}
	if foundUser == nil {
		return shared.PipeError[VerificationResponse](messages.Invalid_Credentials)
	}
	if foundUser.EmailVerified {
		return shared.PipeError[VerificationResponse](messages.Email_Already_Verified)
	}

	if err := p.sendVerificationOTP(ctx, foundUser); err != nil {
		logInternalError(err, "resend_verification.send_otp")
		return shared.PipeError[VerificationResponse](messages.Internal_Error)
	}

	return shared.PipeSuccess(messages.Verification_Sent, &VerificationResponse{
		ExpiresInSeconds: int64(p.cfg.EmailOTPTTL.Seconds()),
	})
}
