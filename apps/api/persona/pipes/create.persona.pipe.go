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
	dto.Handle = strings.ToLower(strings.TrimSpace(dto.Handle))
	dto.DisplayName = strings.TrimSpace(dto.DisplayName)
	dto.Bio = strings.TrimSpace(dto.Bio)
	dto.AvatarURL = strings.TrimSpace(dto.AvatarURL)
	dto.CoverURL = strings.TrimSpace(dto.CoverURL)

	if !validPersonaType(dto.PersonaType) {
		return pipeError[models.Persona](messages.Invalid_Persona_Type)
	}

	persona, err := p.repo.CreatePersona(ctx, userID, dto)
	if err != nil {
		switch err {
		case persona_repo.ErrHandleAlreadyTaken:
			return pipeError[models.Persona](messages.Handle_Already_Taken)
		case persona_repo.ErrGhostPersonaAlreadyUsed:
			return pipeError[models.Persona](messages.Ghost_Persona_Already_Set)
		default:
			return pipeInternalError[models.Persona](err, "persona.create")
		}
	}

	return pipeSuccess(messages.Persona_Created, persona)
}
