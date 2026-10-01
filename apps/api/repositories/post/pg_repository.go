package post

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgRepository struct {
	db *pgxpool.Pool
}

func newPgRepository(db *pgxpool.Pool) *pgRepository {
	return &pgRepository{db: db}
}

type FeedCursor struct {
	CreatedAt time.Time
	ID        uuid.UUID
}

type FeedOptions struct {
	Limit  int
	Cursor *FeedCursor
}

// clamps one feed options payload to supported bounds.
func NormalizeFeedOptions(options FeedOptions) FeedOptions {
	options.Limit = normalizeLimit(options.Limit)
	return options
}
