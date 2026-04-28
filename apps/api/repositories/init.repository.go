package repositories

import (
	"github.com/emmanuella-codes/nox/repositories/persona"
	"github.com/emmanuella-codes/nox/repositories/post"
	"github.com/emmanuella-codes/nox/repositories/user"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repositories struct {
	User user.UserRepository
	Persona persona.PersonaRepository
	Post post.PostRepository
}

func Init(pool *pgxpool.Pool) *Repositories {
	return &Repositories{
		User: user.NewUserRepository(pool),
		Persona: persona.NewPersonaRepository(pool),
		Post: post.NewPostRepository(pool),
	}
}
