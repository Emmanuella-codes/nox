package pipes

import (
	"time"

	"github.com/google/uuid"
)

// memberKey builds one stable member lookup key.
func memberKey(conversationID uuid.UUID, personaID uuid.UUID) string {
	return conversationID.String() + ":" + personaID.String()
}

// followKey builds one stable follow lookup key.
func followKey(followerID uuid.UUID, followingID uuid.UUID) string {
	return followerID.String() + ":" + followingID.String()
}

// timePtr returns a pointer to the supplied timestamp.
func timePtr(value time.Time) *time.Time {
	return &value
}
