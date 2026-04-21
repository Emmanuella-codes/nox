package services

import "golang.org/x/crypto/bcrypt"

const dummyPassword = "nox-dummy-password"

type HashService struct {
	dummyHash string
}

func NewHashService() *HashService {
	hash, _ := bcrypt.GenerateFromPassword([]byte(dummyPassword), bcrypt.DefaultCost)
	return &HashService{dummyHash: string(hash)}
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

func (s *HashService) CompareDummyPassword(password string) {
	_ = bcrypt.CompareHashAndPassword([]byte(s.dummyHash), []byte(password))
}
