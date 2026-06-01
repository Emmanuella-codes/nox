package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/models"
	"github.com/emmanuella-codes/nox/persona/messages"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

type PersonaPipe struct {
	repo persona_repo.PersonaRepository
}

func NewPersonaPipe(repo persona_repo.PersonaRepository) *PersonaPipe {
	return &PersonaPipe{repo: repo}
}

func pipeInternalError[T any](err error, operation string) *shared.PipeRes[T] {
	return shared.PipeInternalError[T](err, "persona", operation, messages.Internal_Error)
}

func validPersonaType(personaType models.PersonaType) bool {
	return personaType == models.VisiblePersonaType
}

func (p *PersonaPipe) ensureOwnedPersona(ctx context.Context, userID uuid.UUID, personaID uuid.UUID) (*models.Persona, *shared.PipeRes[models.Persona]) {
	persona, err := p.repo.FindPersonaByID(ctx, personaID)
	if err != nil {
		if err == persona_repo.ErrPersonaNotFound {
			return nil, shared.PipeError[models.Persona](messages.Persona_Not_Found)
		}
		return nil, pipeInternalError[models.Persona](err, "persona.find_by_id")
	}
	if persona.UserID != userID {
		return nil, shared.PipeError[models.Persona](messages.Forbidden)
	}
	return persona, nil
}
