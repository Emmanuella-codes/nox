package pipes

import (
	"context"
	"strings"

	"github.com/emmanuella-codes/nox/models"
	"github.com/emmanuella-codes/nox/persona/dtos"
	"github.com/emmanuella-codes/nox/persona/messages"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

func (p *PersonaPipe) UpdatePersonaPipe(ctx context.Context, userID uuid.UUID, personaID uuid.UUID, dto dtos.UpdatePersonaDTO) *shared.PipeRes[models.Persona] {
	ownedPersona, res := p.ensureOwnedPersona(ctx, userID, personaID)
	if res != nil {
		return res
	}

	dto.DisplayName = strings.TrimSpace(dto.DisplayName)
	dto.Bio = strings.TrimSpace(dto.Bio)
	dto.AvatarURL = strings.TrimSpace(dto.AvatarURL)
	dto.CoverURL = strings.TrimSpace(dto.CoverURL)

	persona, err := p.repo.UpdatePersona(ctx, ownedPersona.ID, dto)
	if err != nil {
		if err == persona_repo.ErrPersonaNotFound {
			return shared.PipeError[models.Persona](messages.Persona_Not_Found)
		}
		return pipeInternalError[models.Persona](err, "persona.update")
	}

	return shared.PipeSuccess(messages.Persona_Updated, persona)
}
