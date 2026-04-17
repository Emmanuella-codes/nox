package models

import (
	"time"

	"github.com/google/uuid"
)

type Like struct {
	PersonaID uuid.UUID `json:"persona_id"`
	PostID    uuid.UUID `json:"post_id"`
	CreatedAt time.Time `json:"created_at"`
}
