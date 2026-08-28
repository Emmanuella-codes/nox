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

// NewPersonaPipe builds the public-profile orchestration layer from the repository.
func NewPersonaPipe(repo persona_repo.PersonaRepository) *PersonaPipe {
	return &PersonaPipe{repo: repo}
}

// pipeInternalError maps internal profile errors to pipe responses.
func pipeInternalError[T any](err error, operation string) *shared.PipeRes[T] {
	return shared.PipeInternalError[T](err, "persona", operation, messages.Internal_Error)
}

// validPersonaType accepts the single supported public profile type for now.
func validPersonaType(personaType models.PersonaType) bool {
	return personaType == "" || personaType == models.VisiblePersonaType
}

// validPersonaCategory validates the supported public profile categories.
func validPersonaCategory(category models.PersonaCategory) bool {
	switch category {
	case models.PatronPersonaCategory, models.DJPersonaCategory, models.OrganizerPersonaCategory:
		return true
	default:
		return false
	}
}

// ensureOwnedPersona fetches one public profile and verifies ownership.
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
