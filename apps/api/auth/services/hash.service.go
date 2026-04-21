package services

import "golang.org/x/crypto/bcrypt"

type HashService struct{}

func NewHashService() *HashService {
	return &HashService{}
}

func (s *HashService) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (s *HashService) ComparePassword(hash string, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
