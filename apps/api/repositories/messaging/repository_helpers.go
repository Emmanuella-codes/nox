package messaging

import (
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// orderedPersonaIDs sorts direct conversation participant ids into a stable pair.
func orderedPersonaIDs(a uuid.UUID, b uuid.UUID) (uuid.UUID, uuid.UUID) {
	if a.String() > b.String() {
		return b, a
	}
	return a, b
}

// mapNotFound maps pgx not-found errors into repository-level errors.
func mapNotFound(err error, fallback error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return fallback
	}
	return err
}
