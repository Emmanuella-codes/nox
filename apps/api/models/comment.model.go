package models

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID        uuid.UUID `json:"id"`
	PersonaID uuid.UUID `json:"persona_id"`
	PostID    uuid.UUID `json:"post_id"`
	Body      string    `json:"body"` // 280 characters
	ParentID  uuid.UUID `json:"parent_id"`
	LikeCount int       `json:"like_count"`
	CreatedAt time.Time `json:"created_at"`
}
