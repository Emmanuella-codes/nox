package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/models"
	"github.com/emmanuella-codes/nox/persona/messages"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type PersonaPipe struct {
	repo persona_repo.PersonaRepository
}

func NewPersonaPipe(repo persona_repo.PersonaRepository) *PersonaPipe {
	return &PersonaPipe{repo: repo}
}

func pipeSuccess[T any](message shared.PipeMessage, data *T) *shared.PipeRes[T] {
	return &shared.PipeRes[T]{Success: true, Message: message, Data: data}
}

func pipeError[T any](message shared.PipeMessage) *shared.PipeRes[T] {
	return &shared.PipeRes[T]{Success: false, Message: message}
}

func pipeInternalError[T any](err error, operation string) *shared.PipeRes[T] {
	if err != nil {
		log.Error().Err(err).Str("operation", operation).Msg("persona internal error")
	}
	return pipeError[T](messages.Internal_Error)
}

func validPersonaType(personaType models.PersonaType) bool {
	return personaType == models.VisiblePersonaType || personaType == models.GhostPersonaType
}

func (p *PersonaPipe) ensureOwnedPersona(ctx context.Context, userID uuid.UUID, personaID uuid.UUID) (*models.Persona, *shared.PipeRes[models.Persona]) {
	persona, err := p.repo.FindPersonaByID(ctx, personaID)
	if err != nil {
		if err == persona_repo.ErrPersonaNotFound {
			return nil, pipeError[models.Persona](messages.Persona_Not_Found)
		}
		return nil, pipeInternalError[models.Persona](err, "persona.find_by_id")
	}
	if persona.UserID != userID {
		return nil, pipeError[models.Persona](messages.Forbidden)
	}
	return persona, nil
}
