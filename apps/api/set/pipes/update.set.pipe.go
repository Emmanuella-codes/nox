package pipes

import (
	"context"
	"strings"

	"github.com/emmanuella-codes/nox/models"
	set_repo "github.com/emmanuella-codes/nox/repositories/set"
	"github.com/emmanuella-codes/nox/set/dtos"
	"github.com/emmanuella-codes/nox/set/messages"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

func (p *SetPipe) UpdateSetPipe(ctx context.Context, userID, setID uuid.UUID, dto dtos.UpdateSetDTO) *shared.PipeRes[models.Set] {
	dto.Title, dto.Description = strings.TrimSpace(dto.Title), strings.TrimSpace(dto.Description)
	tags, valid := normalizeGenreTags(dto.GenreTags)
	if dto.Title == "" || !valid {
		return shared.PipeError[models.Set](messages.Invalid_Set)
	}
	current, err := p.setRepo.FindSetByID(ctx, setID)
	if err != nil {
		if err == set_repo.ErrSetNotFound {
			return shared.PipeError[models.Set](messages.Set_Not_Found)
		}
		return pipeInternalError[models.Set](err, "set.update_find")
	}
	if current.AuthorUserID != userID {
		return shared.PipeError[models.Set](messages.Forbidden)
	}
	dto.GenreTags = tags
	duration := current.DurationSeconds
	if dto.MediaAssetID != nil {
		asset, findErr := p.mediaRepo.FindMediaAssetByID(ctx, *dto.MediaAssetID)
		if findErr != nil {
			return shared.PipeError[models.Set](messages.Media_Not_Found)
		}
		if asset.OwnerUserID != userID || asset.OwnerPersonaID != current.PersonaID || !validSetMedia(asset) {
			return shared.PipeError[models.Set](messages.Invalid_Set)
		}
		duration = asset.DurationSeconds
	}
	updated, err := p.setRepo.UpdateSet(ctx, userID, setID, duration, dto)
	if err != nil {
		if err == set_repo.ErrSetMediaInUse {
			return shared.PipeError[models.Set](messages.Media_In_Use)
		}
		if err == set_repo.ErrSetNotFound {
			return shared.PipeError[models.Set](messages.Set_Not_Found)
		}
		return pipeInternalError[models.Set](err, "set.update")
	}
	if err := p.hydrateSet(ctx, updated); err != nil {
		return pipeInternalError[models.Set](err, "set.update_hydrate")
	}
	return shared.PipeSuccess(messages.Set_Updated, updated)
}
