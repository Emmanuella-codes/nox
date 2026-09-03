package notification

import "github.com/jackc/pgx/v5/pgxpool"

type pgRepository struct {
	db *pgxpool.Pool
}

func newPgRepository(db *pgxpool.Pool) NotificationRepository {
	return &pgRepository{db: db}
}
