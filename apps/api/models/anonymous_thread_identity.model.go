package models

import (
	"time"

	"github.com/google/uuid"
)

type AnonymousThreadIdentity struct {
	ID                 uuid.UUID `json:"id"`
	ThreadID           uuid.UUID `json:"thread_id"`
	UserID             uuid.UUID `json:"-"`
	PersonaID          uuid.UUID `json:"-"`
	AnonymousHandle    string    `json:"anonymous_handle"`
	AnonymousAvatarKey string    `json:"anonymous_avatar_key"`
	CreatedAt          time.Time `json:"created_at"`
}
