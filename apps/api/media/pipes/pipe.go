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
	cloudinaryclient "github.com/emmanuella-codes/nox/shared/cloudinary/client"
	"github.com/google/uuid"
)

type MediaPipe struct {
	mediaRepo        media_repo.MediaRepository
	personaRepo      persona_repo.PersonaRepository
	config           *config.Config
	cloudinaryClient *cloudinaryclient.Client
}

// NewMediaPipe builds the media orchestration layer from repositories and config.
func NewMediaPipe(mediaRepo media_repo.MediaRepository, personaRepo persona_repo.PersonaRepository, cfg *config.Config) *MediaPipe {
	var cloudinary *cloudinaryclient.Client
	if cfg != nil {
		cloudinary = cloudinaryclient.New(cloudinaryclient.Config{
			CloudName:    cfg.CloudinaryCloudName,
			APIKey:       cfg.CloudinaryAPIKey,
			APISecret:    cfg.CloudinaryAPISecret,
			UploadFolder: cfg.CloudinaryUploadFolder,
		})
	}
	return &MediaPipe{mediaRepo: mediaRepo, personaRepo: personaRepo, config: cfg, cloudinaryClient: cloudinary}
}

// pipeInternalError maps internal media errors to pipe responses.
func pipeInternalError[T any](err error, operation string) *shared.PipeRes[T] {
	return shared.PipeInternalError[T](err, "media", operation, messages.Internal_Error)
}

// validSetVideo validates uploaded set video metadata.
func validSetVideo(mimeType string, durationSeconds int) bool {
	switch strings.ToLower(strings.TrimSpace(mimeType)) {
	case "video/mp4", "video/webm", "video/quicktime":
	default:
		return false
	}
	return durationSeconds > 0 && durationSeconds <= 900
}

// validStoryVideo validates uploaded story video metadata.
func validStoryVideo(mimeType string, durationSeconds int) bool {
	return validSetVideoMime(mimeType) && durationSeconds > 0 && durationSeconds <= 300
}

// validSetVideoMime validates supported set and story video mime types.
func validSetVideoMime(mimeType string) bool {
	switch strings.ToLower(strings.TrimSpace(mimeType)) {
	case "video/mp4", "video/webm", "video/quicktime":
		return true
	default:
		return false
	}
}

// validPostMedia validates supported uploaded media kinds for posts and messaging attachments.
func validPostMedia(kind models.MediaKind, mimeType string, sizeBytes int64, durationSeconds int) bool {
	if sizeBytes <= 0 {
		return false
	}
	switch kind {
	case models.ImageMediaKind:
		switch strings.ToLower(strings.TrimSpace(mimeType)) {
		case "image/jpeg", "image/png", "image/webp", "image/gif":
			return true
		default:
			return false
		}
	case models.VideoMediaKind:
		return validSetVideoMime(mimeType) && durationSeconds <= 300
	case models.AudioMediaKind:
		switch strings.ToLower(strings.TrimSpace(mimeType)) {
		case "audio/mpeg", "audio/mp3", "audio/wav", "audio/x-wav", "audio/aac", "audio/mp4", "audio/ogg", "audio/webm":
			return durationSeconds > 0 && durationSeconds <= 900
		default:
			return false
		}
	default:
		return false
	}
}

// cloudinaryResourceType maps media kinds into Cloudinary upload resource types.
func cloudinaryResourceType(kind models.MediaKind) string {
	if kind == models.VideoMediaKind || kind == models.AudioMediaKind {
		return "video"
	}
	return "image"
}

// uploadURL builds one upload target URL for pending assets.
func (p *MediaPipe) uploadURL(storageKey string) string {
	if p.config == nil || strings.TrimSpace(p.config.MediaUploadBaseURL) == "" {
		return ""
	}
	return strings.TrimRight(p.config.MediaUploadBaseURL, "/") + "/" + storageKey
}

// playbackURL builds one public playback URL for pending assets.
func (p *MediaPipe) playbackURL(storageKey string) string {
	if p.config == nil || strings.TrimSpace(p.config.MediaPublicBaseURL) == "" {
		return ""
	}
	return strings.TrimRight(p.config.MediaPublicBaseURL, "/") + "/" + storageKey
}

// setVideoStorageKey builds one storage key for uploaded set media.
func setVideoStorageKey(ownerPersonaID string) string {
	return fmt.Sprintf("sets/%s/%s", ownerPersonaID, uuid.NewString())
}

// storyVideoStorageKey builds one storage key for uploaded story media.
func storyVideoStorageKey(ownerPersonaID string) string {
	return fmt.Sprintf("stories/%s/%s", ownerPersonaID, uuid.NewString())
}

// canOwnSetMedia checks whether one public profile can own set media.
func canOwnSetMedia(persona *models.Persona) bool {
	return persona.PersonaType == models.VisiblePersonaType &&
		(persona.Category == models.PatronPersonaCategory ||
			persona.Category == models.DJPersonaCategory ||
			persona.Category == models.OrganizerPersonaCategory)
}
