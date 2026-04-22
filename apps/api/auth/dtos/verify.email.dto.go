package dtos

type VerifyEmailDTO struct {
	Email string `json:"email" validate:"required,email"`
	OTP   string `json:"otp" validate:"required,len=6"`
}

type ResendVerificationDTO struct {
	Email string `json:"email" validate:"required,email"`
}
