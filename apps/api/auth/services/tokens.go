package services

import (
	"github.com/emmanuella-codes/nox/config"
	sharedtoken "github.com/emmanuella-codes/nox/shared/token"
	"github.com/google/uuid"
)

type TokenPair struct {
	AccessToken    string `json:"access_token"`
	RefreshToken   string `json:"refresh_token"`
	RefreshTokenID string `json:"-"`
}

type TokenService struct {
	cfg *config.Config
}

func NewTokenService(cfg *config.Config) *TokenService {
	return &TokenService{cfg: cfg}
}

func (s *TokenService) IssuePair(userID uuid.UUID) (*TokenPair, error) {
	accessToken, err := s.SignAccessToken(userID)
	if err != nil {
		return nil, err
	}

	refreshTokenID := uuid.NewString()
	refreshToken, err := s.SignRefreshToken(userID, refreshTokenID)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:    accessToken,
		RefreshToken:   refreshToken,
		RefreshTokenID: refreshTokenID,
	}, nil
}

func (s *TokenService) SignAccessToken(userID uuid.UUID) (string, error) {
	return sharedtoken.Sign(userID, sharedtoken.AccessTokenType, s.cfg.JWTAccessSecret, s.cfg.JWTAccessTTL)
}

func (s *TokenService) SignRefreshToken(userID uuid.UUID, tokenID string) (string, error) {
	return sharedtoken.SignWithID(userID, tokenID, sharedtoken.RefreshTokenType, s.cfg.JWTRefreshSecret, s.cfg.JWTRefreshTTL)
}

func (s *TokenService) VerifyAccessToken(rawToken string) (*sharedtoken.Claims, error) {
	return sharedtoken.Verify(rawToken, s.cfg.JWTAccessSecret, sharedtoken.AccessTokenType)
}

func (s *TokenService) VerifyRefreshToken(rawToken string) (*sharedtoken.Claims, error) {
	return sharedtoken.Verify(rawToken, s.cfg.JWTRefreshSecret, sharedtoken.RefreshTokenType)
}
