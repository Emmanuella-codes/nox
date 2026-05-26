package pipes

import (
	"context"
	"errors"

	"github.com/emmanuella-codes/nox/auth/dtos"
	"github.com/emmanuella-codes/nox/auth/messages"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/redis/go-redis/v9"
)

func (p *AuthPipe) VerifyEmailPipe(ctx context.Context, dto dtos.VerifyEmailDTO) *shared.PipeRes[any] {
	foundUser, err := p.userRepo.FindUserByEmail(ctx, normalizeEmail(dto.Email))
	if err != nil {
		logInternalError(err, "verify_email.find_user_by_email")
		return shared.PipeError[any](messages.Internal_Error)
	}
	if foundUser == nil {
		return shared.PipeError[any](messages.Invalid_Credentials)
	}
	if foundUser.EmailVerified {
		return shared.PipeError[any](messages.Email_Already_Verified)
	}

	otpHash, err := p.redis.Get(ctx, emailVerificationKey(foundUser.ID.String())).Result()
	if errors.Is(err, redis.Nil) {
		return shared.PipeError[any](messages.OTP_Expired)
	}
	if err != nil {
		logInternalError(err, "verify_email.get_otp")
		return shared.PipeError[any](messages.Internal_Error)
	}
	if !p.otpService.Compare(otpHash, dto.OTP) {
		attempts, err := p.recordFailedVerificationAttempt(ctx, foundUser.ID.String())
		if err != nil {
			logInternalError(err, "verify_email.record_failed_attempt")
			return shared.PipeError[any](messages.Internal_Error)
		}
		if attempts >= maxEmailVerificationAttempts {
			if err := p.deleteEmailVerification(ctx, foundUser.ID.String()); err != nil {
				logInternalError(err, "verify_email.delete_locked_otp")
				return shared.PipeError[any](messages.Internal_Error)
			}
			return shared.PipeError[any](messages.OTP_Locked)
		}
		return shared.PipeError[any](messages.Invalid_OTP)
	}

	if err := p.userRepo.MarkEmailVerified(ctx, foundUser.ID.String()); err != nil {
		logInternalError(err, "verify_email.mark_verified")
		return shared.PipeError[any](messages.Internal_Error)
	}
	if err := p.deleteEmailVerification(ctx, foundUser.ID.String()); err != nil {
		logInternalError(err, "verify_email.delete_otp")
		return shared.PipeError[any](messages.Internal_Error)
	}

	return shared.PipeSuccess[any](messages.Email_Verified, nil)
}
