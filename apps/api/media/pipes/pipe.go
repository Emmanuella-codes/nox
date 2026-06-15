package pipes

import (
	"fmt"
	"strings"

	"github.com/emmanuella-codes/nox/config"
	"github.com/emmanuella-codes/nox/media/messages"
	"github.com/emmanuella-codes/nox/models"
	media_repo "github.com/emmanuella-codes/nox/repositories/media"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

type MediaPipe struct {
	mediaRepo   media_repo.MediaRepository
	personaRepo persona_repo.PersonaRepository
	config      *config.Config
}

func NewMediaPipe(mediaRepo media_repo.MediaRepository, personaRepo persona_repo.PersonaRepository, cfg *config.Config) *MediaPipe {
	return &MediaPipe{mediaRepo: mediaRepo, personaRepo: personaRepo, config: cfg}
}

func pipeInternalError[T any](err error, operation string) *shared.PipeRes[T] {
	return shared.PipeInternalError[T](err, "media", operation, messages.Internal_Error)
}

func validSetVideo(mimeType string, durationSeconds int) bool {
	switch strings.ToLower(strings.TrimSpace(mimeType)) {
	case "video/mp4", "video/webm", "video/quicktime":
	default:
		return false
	}
	return durationSeconds > 0 && durationSeconds <= 900
}

func validStoryVideo(mimeType string, durationSeconds int) bool {
	return validSetVideoMime(mimeType) && durationSeconds > 0 && durationSeconds <= 300
}

func validSetVideoMime(mimeType string) bool {
	switch strings.ToLower(strings.TrimSpace(mimeType)) {
	case "video/mp4", "video/webm", "video/quicktime":
		return true
	default:
		return false
	}
}

func (p *MediaPipe) uploadURL(storageKey string) string {
	if p.config == nil || strings.TrimSpace(p.config.MediaUploadBaseURL) == "" {
		return ""
	}
	return strings.TrimRight(p.config.MediaUploadBaseURL, "/") + "/" + storageKey
}

func (p *MediaPipe) playbackURL(storageKey string) string {
	if p.config == nil || strings.TrimSpace(p.config.MediaPublicBaseURL) == "" {
		return ""
	}
	return strings.TrimRight(p.config.MediaPublicBaseURL, "/") + "/" + storageKey
}

func setVideoStorageKey(ownerPersonaID string) string {
	return fmt.Sprintf("sets/%s/%s", ownerPersonaID, uuid.NewString())
}

func storyVideoStorageKey(ownerPersonaID string) string {
	return fmt.Sprintf("stories/%s/%s", ownerPersonaID, uuid.NewString())
}

func isDJPersona(persona *models.Persona) bool {
	return persona.PersonaType == models.VisiblePersonaType && persona.Category == models.DJPersonaCategory
}
