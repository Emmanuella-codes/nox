package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	AccessTokenType  = "access"
	RefreshTokenType = "refresh"
)

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Type   string    `json:"type"`
	jwt.RegisteredClaims
}

// Sign creates a token with a generated jti.
func Sign(userID uuid.UUID, tokenType, secret string, ttl time.Duration) (string, error) {
	return signWithID(userID, uuid.NewString(), tokenType, secret, ttl)
}

// SignWithID creates a token with an explicit jti.
func SignWithID(userID uuid.UUID, tokenID, tokenType, secret string, ttl time.Duration) (string, error) {
	return signWithID(userID, tokenID, tokenType, secret, ttl)
}

func signWithID(userID uuid.UUID, tokenID, tokenType, secret string, ttl time.Duration) (string, error) {
	if secret == "" {
		return "", errors.New("jwt secret is empty")
	}
	if tokenID == "" {
		return "", errors.New("token id is empty")
	}

	now := time.Now()
	claims := Claims{
		UserID: userID,
		Type:   tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			ID:        tokenID,
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(secret))
}

// Verify parses and validates a token, checking type and required claims.
func Verify(rawToken, secret, expectedType string) (*Claims, error) {
	if secret == "" {
		return nil, errors.New("jwt secret is empty")
	}

	t, err := jwt.ParseWithClaims(rawToken, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %s", t.Method.Alg())
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := t.Claims.(*Claims)
	if !ok || !t.Valid {
		return nil, errors.New("invalid token claims")
	}
	if claims.Type != expectedType {
		return nil, errors.New("invalid token type")
	}
	if claims.UserID == uuid.Nil {
		return nil, errors.New("missing user id claim")
	}
	if claims.ID == "" {
		return nil, errors.New("missing token id claim")
	}

	return claims, nil
}
