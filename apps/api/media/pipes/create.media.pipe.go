package pipes

import (
	"context"
	"strings"

	"github.com/emmanuella-codes/nox/media/dtos"
	"github.com/emmanuella-codes/nox/media/messages"
	"github.com/emmanuella-codes/nox/models"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

func (p *MediaPipe) CreateSetVideoAssetPipe(ctx context.Context, userID uuid.UUID, dto dtos.CreateMediaAssetDTO) *shared.PipeRes[models.MediaAsset] {
	dto.StorageKey = strings.TrimSpace(dto.StorageKey)
	dto.PlaybackURL = strings.TrimSpace(dto.PlaybackURL)
	dto.ThumbnailURL = strings.TrimSpace(dto.ThumbnailURL)
	dto.MimeType = strings.TrimSpace(dto.MimeType)

	if !validSetVideo(dto.MimeType, dto.DurationSeconds) {
		return shared.PipeError[models.MediaAsset](messages.Invalid_Media)
	}

	persona, err := p.personaRepo.FindPersonaByID(ctx, dto.OwnerPersonaID)
	if err != nil {
		if err == persona_repo.ErrPersonaNotFound {
			return shared.PipeError[models.MediaAsset](messages.Persona_Not_Found)
		}
		return pipeInternalError[models.MediaAsset](err, "media.find_persona")
	}
	if persona.UserID != userID || !canOwnSetMedia(persona) {
		return shared.PipeError[models.MediaAsset](messages.Forbidden)
	}

	asset, err := p.mediaRepo.CreateMediaAsset(ctx, userID, dto)
	if err != nil {
		return pipeInternalError[models.MediaAsset](err, "media.create")
	}
	return shared.PipeSuccess(messages.Media_Asset_Created, asset)
}
