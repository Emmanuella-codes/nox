package post

import "github.com/jackc/pgx/v5/pgxpool"

type pgRepository struct {
	db *pgxpool.Pool
}

// newPgRepository builds the Postgres-backed post repository.
func newPgRepository(db *pgxpool.Pool) *pgRepository {
	return &pgRepository{db: db}
}
