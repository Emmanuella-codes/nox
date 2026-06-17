package story

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func normalizeLimit(limit int) int {
	if limit <= 0 {
		return 20
	}
	if limit > 50 {
		return 50
	}
	return limit
}

func normalizeOffset(offset int) int {
	if offset < 0 {
		return 0
	}
	return offset
}

func emptyToNil(value string) any {
	if value == "" {
		return nil
	}
	return value
}

func boundedPosition(position int, maxPosition int) int {
	if position > maxPosition {
		return maxPosition
	}
	if position < 1 {
		return 1
	}
	return position
}

func shiftPositions(ctx context.Context, tx pgx.Tx, table string, scopeColumn string, scopeID uuid.UUID, position int, currentPosition int) error {
	if position < currentPosition {
		if _, err := tx.Exec(ctx, fmt.Sprintf(`
			UPDATE %s
			SET position = -position
			WHERE %s = $1 AND position >= $2 AND position < $3
		`, table, scopeColumn), scopeID, position, currentPosition); err != nil {
			return err
		}
		_, err := tx.Exec(ctx, fmt.Sprintf(`UPDATE %s SET position = -position + 1 WHERE %s = $1 AND position < 0`, table, scopeColumn), scopeID)
		return err
	}
	if position > currentPosition {
		if _, err := tx.Exec(ctx, fmt.Sprintf(`
			UPDATE %s
			SET position = -position
			WHERE %s = $1 AND position <= $2 AND position > $3
		`, table, scopeColumn), scopeID, position, currentPosition); err != nil {
			return err
		}
		_, err := tx.Exec(ctx, fmt.Sprintf(`UPDATE %s SET position = -position - 1 WHERE %s = $1 AND position < 0`, table, scopeColumn), scopeID)
		return err
	}
	return nil
}
