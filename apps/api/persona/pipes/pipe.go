package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/models"
	"github.com/emmanuella-codes/nox/persona/messages"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	preference_repo "github.com/emmanuella-codes/nox/repositories/preference"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

type PersonaPipe struct {
	repo           persona_repo.PersonaRepository
	preferenceRepo preference_repo.PreferenceRepository
}

func NewPersonaPipe(repo persona_repo.PersonaRepository, deps ...any) *PersonaPipe {
	pipe := &PersonaPipe{repo: repo}
	for _, dep := range deps {
		if preferenceRepo, ok := dep.(preference_repo.PreferenceRepository); ok {
			pipe.preferenceRepo = preferenceRepo
		}
	}
	return pipe
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

func (p *PersonaPipe) FindViewerPersona(ctx context.Context, userID uuid.UUID, personaID uuid.UUID) (*models.Persona, shared.PipeMessage) {
	persona, err := p.repo.FindPersonaByID(ctx, personaID)
	if err != nil {
		if err == persona_repo.ErrPersonaNotFound {
			return nil, messages.Persona_Not_Found
		}
		return nil, messages.Internal_Error
	}
	if persona.UserID != userID {
		return nil, messages.Forbidden
	}
	return persona, ""
}

func (p *PersonaPipe) visibleToViewer(ctx context.Context, viewerPersonaID uuid.UUID, target *models.Persona) (bool, error) {
	if p.preferenceRepo == nil || target == nil {
		return true, nil
	}
	viewer, err := p.repo.FindPersonaByID(ctx, viewerPersonaID)
	if err != nil {
		return false, err
	}
	if viewer.UserID == target.UserID {
		return true, nil
	}
	return blockedVisibility(ctx, p.preferenceRepo, viewer.UserID, target.UserID)
}

func blockedVisibility(ctx context.Context, preferenceRepo preference_repo.PreferenceRepository, viewerUserID uuid.UUID, targetUserID uuid.UUID) (bool, error) {
	blocked, err := preferenceRepo.IsBlockedBetween(ctx, viewerUserID, targetUserID)
	if err != nil {
		return false, err
	}
	return !blocked, nil
}
