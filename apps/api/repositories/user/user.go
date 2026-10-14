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
	FindUserByID(ctx context.Context, userID string) (*models.User, error)
	MarkEmailVerified(ctx context.Context, userID string) error
}

var UserRepo UserRepository

func InitUserRepo(db *pgxpool.Pool) {
	UserRepo = newPgRepository(db)
}
