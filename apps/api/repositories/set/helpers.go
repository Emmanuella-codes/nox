package set

import "github.com/google/uuid"

func uuidToNil(id uuid.UUID) any {
	if id == uuid.Nil {
		return nil
	}
	return id
}

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
