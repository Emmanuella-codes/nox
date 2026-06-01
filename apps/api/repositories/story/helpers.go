package story

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
