package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
)

func (p *AdminPipe) revokeAllUserSessions(ctx context.Context, userID uuid.UUID) error {
	if err := p.revokeAppRefreshSessions(ctx, userID); err != nil {
		return err
	}

	membership, err := p.adminRepo.FindMembershipByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if membership == nil {
		return nil
	}

	return p.revokeUserRefreshSessions(ctx, userID)
}

func (p *AdminPipe) revokeAppRefreshSessions(ctx context.Context, userID uuid.UUID) error {
	sessionIDs, err := p.redis.SMembers(ctx, appUserRefreshSessionsKey(userID)).Result()
	if err != nil {
		return err
	}

	pipe := p.redis.TxPipeline()
	for _, sessionID := range sessionIDs {
		pipe.Del(ctx, appRefreshSessionKey(sessionID))
	}
	pipe.Del(ctx, appUserRefreshSessionsKey(userID))
	_, err = pipe.Exec(ctx)
	return err
}

func (p *AdminPipe) sendAppVerificationOTP(ctx context.Context, user *models.User) error {
	otp, err := p.otpService.Generate()
	if err != nil {
		return err
	}

	otpHash, err := p.otpService.Hash(otp)
	if err != nil {
		return err
	}

	pipe := p.redis.TxPipeline()
	pipe.Set(ctx, appEmailVerificationKey(user.ID.String()), otpHash, p.cfg.EmailOTPTTL)
	pipe.Del(ctx, appEmailVerificationAttemptsKey(user.ID.String()))
	if _, err := pipe.Exec(ctx); err != nil {
		return err
	}

	return p.emailService.SendVerificationOTP(ctx, user.Email, otp)
}

func (p *AdminPipe) deleteAppEmailVerification(ctx context.Context, userID string) error {
	_, err := p.redis.Del(ctx, appEmailVerificationKey(userID), appEmailVerificationAttemptsKey(userID)).Result()
	return err
}
