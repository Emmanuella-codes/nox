package pipes

import (
	"context"
	"strings"

	"github.com/emmanuella-codes/nox/models"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	"github.com/emmanuella-codes/nox/set/dtos"
	"github.com/emmanuella-codes/nox/set/messages"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (p *SetPipe) CreateSetPipe(ctx context.Context, userID uuid.UUID, dto dtos.CreateSetDTO) *shared.PipeRes[models.Set] {
	dto.Title = strings.TrimSpace(dto.Title)
	dto.Description = strings.TrimSpace(dto.Description)
	if dto.Title == "" {
		return shared.PipeError[models.Set](messages.Invalid_Set)
	}

	persona, err := p.personaRepo.FindPersonaByID(ctx, dto.PersonaID)
	if err != nil {
		if err == persona_repo.ErrPersonaNotFound {
			return shared.PipeError[models.Set](messages.Persona_Not_Found)
		}
		return pipeInternalError[models.Set](err, "set.find_persona")
	}
	if persona.UserID != userID || !isDJPersona(persona) {
		return shared.PipeError[models.Set](messages.Forbidden)
	}

	asset, err := p.mediaRepo.FindMediaAssetByID(ctx, dto.MediaAssetID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return shared.PipeError[models.Set](messages.Media_Not_Found)
		}
		return pipeInternalError[models.Set](err, "set.find_media")
	}
	if asset.OwnerUserID != userID || asset.OwnerPersonaID != dto.PersonaID || !validSetMedia(asset) {
		return shared.PipeError[models.Set](messages.Invalid_Set)
	}

	set, err := p.setRepo.CreateSet(ctx, userID, asset.DurationSeconds, dto)
	if err != nil {
		return pipeInternalError[models.Set](err, "set.create")
	}
	if err := p.hydrateSet(ctx, set); err != nil {
		return pipeInternalError[models.Set](err, "set.hydrate")
	}
	set.Persona = persona
	set.MediaAsset = asset
	return shared.PipeSuccess(messages.Set_Created, set)
}
