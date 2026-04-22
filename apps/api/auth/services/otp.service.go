package services

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"golang.org/x/crypto/bcrypt"
)

const otpDigits = 6

type OTPService struct{}

func NewOTPService() *OTPService {
	return &OTPService{}
}

func (s *OTPService) Generate() (string, error) {
	max := 1
	for range otpDigits {
		max *= 10
	}

	n, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%0*d", otpDigits, n.Int64()), nil
}

func (s *OTPService) Hash(otp string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(otp), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (s *OTPService) Compare(hash string, otp string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(otp)) == nil
}
