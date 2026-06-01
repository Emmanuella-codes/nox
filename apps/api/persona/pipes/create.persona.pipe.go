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

func (p *PersonaPipe) CreatePersonaPipe(ctx context.Context, userID uuid.UUID, dto dtos.CreatePersonaDTO) *shared.PipeRes[models.Persona] {
	dto.DisplayName = strings.TrimSpace(dto.DisplayName)
	dto.Bio = strings.TrimSpace(dto.Bio)
	dto.AvatarURL = strings.TrimSpace(dto.AvatarURL)
	dto.CoverURL = strings.TrimSpace(dto.CoverURL)
	if dto.Category == "" {
		dto.Category = models.FanPersonaCategory
	}

	if !validPersonaType(dto.PersonaType) {
		return shared.PipeError[models.Persona](messages.Invalid_Persona_Type)
	}
	if !validPersonaCategory(dto.Category) {
		return shared.PipeError[models.Persona](messages.Invalid_Persona_Type)
	}

	dto.Handle = strings.ToLower(strings.TrimSpace(dto.Handle))
	if dto.Handle == "" {
		return shared.PipeError[models.Persona](messages.Handle_Required)
	}

	persona, err := p.repo.CreatePersona(ctx, userID, dto)
	if err != nil {
		switch err {
		case persona_repo.ErrHandleAlreadyTaken:
			return shared.PipeError[models.Persona](messages.Handle_Already_Taken)
		case persona_repo.ErrGhostPersonaAlreadyUsed:
			return shared.PipeError[models.Persona](messages.Ghost_Persona_Already_Set)
		default:
			return pipeInternalError[models.Persona](err, "persona.create")
		}
	}

	return shared.PipeSuccess(messages.Persona_Created, persona)
}
