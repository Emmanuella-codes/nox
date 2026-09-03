package pipes

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	post_repo "github.com/emmanuella-codes/nox/repositories/post"
	"github.com/google/uuid"
)

type encodedFeedCursor struct {
	CreatedAt string `json:"created_at"`
	ID        string `json:"id"`
}

// decodeFeedCursor parses one opaque feed cursor into repository feed options.
func decodeFeedCursor(value string) (*post_repo.FeedCursor, error) {
	if value == "" {
		return nil, nil
	}
	raw, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return nil, fmt.Errorf("decode feed cursor: %w", err)
	}
	var payload encodedFeedCursor
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, fmt.Errorf("unmarshal feed cursor: %w", err)
	}
	createdAt, err := time.Parse(time.RFC3339Nano, payload.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("parse feed cursor time: %w", err)
	}
	id, err := uuid.Parse(payload.ID)
	if err != nil {
		return nil, fmt.Errorf("parse feed cursor id: %w", err)
	}
	return &post_repo.FeedCursor{CreatedAt: createdAt, ID: id}, nil
}

// encodeFeedCursor serializes one post into the next-page cursor value.
func encodeFeedCursor(post PostResponse) (string, error) {
	payload := encodedFeedCursor{CreatedAt: post.CreatedAt.Format(time.RFC3339Nano), ID: post.ID}
	raw, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(raw), nil
}
