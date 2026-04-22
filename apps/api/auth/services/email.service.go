package services

import (
	"context"
	"fmt"

	"github.com/emmanuella-codes/nox/shared/mail"
)

type EmailService struct {
	provider mail.Provider
}

func NewEmailService(provider mail.Provider) *EmailService {
	return &EmailService{provider: provider}
}

func (s *EmailService) SendVerificationOTP(ctx context.Context, email string, otp string) error {
	return s.provider.Send(ctx, mail.Message{
		ToEmail: email,
		Subject: "Verify your Nox email",
		TextContent: fmt.Sprintf(
			"Your Nox verification code is %s. This code expires soon. If you did not request this, you can ignore this email.",
			otp,
		),
		HTMLContent: fmt.Sprintf(
			"<p>Your Nox verification code is <strong>%s</strong>.</p><p>This code expires soon. If you did not request this, you can ignore this email.</p>",
			otp,
		),
	})
}
