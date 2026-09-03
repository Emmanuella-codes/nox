package pipes

import "time"

// Sums the visible durations for the current story response.
func storyItemsDuration(items []StoryItemResponse) int {
	total := 0
	for _, item := range items {
		total += item.DurationSeconds
	}
	return total
}

// Returns the latest visible item expiry for the current story response.
func latestStoryItemExpiry(items []StoryItemResponse) (time.Time, bool) {
	var latest time.Time
	for _, item := range items {
		if item.ExpiresAt.After(latest) {
			latest = item.ExpiresAt
		}
	}
	return latest, !latest.IsZero()
}
