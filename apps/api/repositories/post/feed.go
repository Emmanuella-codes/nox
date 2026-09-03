package post

import (
	"time"

	"github.com/google/uuid"
)

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
