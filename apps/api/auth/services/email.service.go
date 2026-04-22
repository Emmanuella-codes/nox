package services

import "github.com/rs/zerolog/log"

type EmailService struct{}

func NewEmailService() *EmailService {
	return &EmailService{}
}

func (s *EmailService) SendVerificationOTP(email string, otp string) error {
	log.Info().Str("email", email).Str("otp", otp).Msg("email verification otp generated")
	return nil
}
