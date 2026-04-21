package user

import (
	"context"
	"errors"

	"github.com/emmanuella-codes/nox/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrUserAlreadyExists = errors.New("user already exists")

type UserRepository interface {
	CreateUser(ctx context.Context, fullname string, email string, passwordHash string) (*models.User, error)
	FindUserByEmail(ctx context.Context, email string) (*models.User, error)
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return newPgRepository(db)
}
