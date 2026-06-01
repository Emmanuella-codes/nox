package pipes

import (
	"strings"

	"github.com/emmanuella-codes/nox/media/messages"
	"github.com/emmanuella-codes/nox/models"
	media_repo "github.com/emmanuella-codes/nox/repositories/media"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	"github.com/emmanuella-codes/nox/shared"
)

type MediaPipe struct {
	mediaRepo   media_repo.MediaRepository
	personaRepo persona_repo.PersonaRepository
}

func NewMediaPipe(mediaRepo media_repo.MediaRepository, personaRepo persona_repo.PersonaRepository) *MediaPipe {
	return &MediaPipe{mediaRepo: mediaRepo, personaRepo: personaRepo}
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

func isDJPersona(persona *models.Persona) bool {
	return persona.PersonaType == models.VisiblePersonaType && persona.Category == models.DJPersonaCategory
}
