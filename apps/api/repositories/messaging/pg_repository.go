package messaging

import "github.com/jackc/pgx/v5/pgxpool"

type pgRepository struct {
	db *pgxpool.Pool
}

func newPgRepository(db *pgxpool.Pool) *pgRepository {
	return &pgRepository{db: db}
}
