package notification

import "github.com/jackc/pgx/v5/pgxpool"

type pgRepository struct {
	db *pgxpool.Pool
}

// newPgRepository builds the Postgres-backed notification repository.
func newPgRepository(db *pgxpool.Pool) NotificationRepository {
	return &pgRepository{db: db}
}
