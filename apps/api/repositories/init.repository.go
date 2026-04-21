package repositories

import (
	"github.com/emmanuella-codes/nox/repositories/user"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repositories struct {
	User user.UserRepository
}

func Init(pool *pgxpool.Pool) *Repositories {
	return &Repositories{
		User: user.NewUserRepository(pool),
	}
}
