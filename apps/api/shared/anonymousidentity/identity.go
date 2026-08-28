package anonymousidentity

import (
	"fmt"

	"github.com/google/uuid"
)

var avatarKeys = [...]string{
	"mask_01",
	"mask_02",
	"mask_03",
	"mask_04",
	"mask_05",
	"mask_06",
}

// GenerateHandle returns a display-only anonymous handle for a thread identity.
func GenerateHandle() string {
	id := uuid.NewString()
	return fmt.Sprintf("ghost_%s", id[:8])
}

// GenerateAvatarKey returns a fallback avatar key for a thread identity.
func GenerateAvatarKey() string {
	id := uuid.New()
	return avatarKeys[int(id[0])%len(avatarKeys)]
}
