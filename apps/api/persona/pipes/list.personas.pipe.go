package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/models"
	"github.com/emmanuella-codes/nox/persona/messages"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

func (p *PersonaPipe) GetPersonaPipe(ctx context.Context, personaID uuid.UUID) *shared.PipeRes[models.Persona] {
	persona, err := p.repo.FindPersonaByID(ctx, personaID)
	if err != nil {
		if err == persona_repo.ErrPersonaNotFound {
			return shared.PipeError[models.Persona](messages.Persona_Not_Found)
		}
		return pipeInternalError[models.Persona](err, "persona.get")
	}

	return shared.PipeSuccess(messages.Persona_Fetched, persona)
}

func (p *PersonaPipe) GetMyPersonasPipe(ctx context.Context, userID uuid.UUID) *shared.PipeRes[[]*models.Persona] {
	personas, err := p.repo.FindPersonasByUserID(ctx, userID)
	if err != nil {
		return pipeInternalError[[]*models.Persona](err, "persona.list_by_user")
	}

	return shared.PipeSuccess(messages.Personas_Listed, &personas)
}
