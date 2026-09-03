package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/models"
	"github.com/emmanuella-codes/nox/preference/messages"
	event_repo "github.com/emmanuella-codes/nox/repositories/event"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	post_repo "github.com/emmanuella-codes/nox/repositories/post"
	preference_repo "github.com/emmanuella-codes/nox/repositories/preference"
	set_repo "github.com/emmanuella-codes/nox/repositories/set"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

type PreferencePipe struct {
	preferenceRepo preference_repo.PreferenceRepository
	personaRepo    persona_repo.PersonaRepository
	postRepo       post_repo.PostRepository
	eventRepo      event_repo.EventRepository
	setRepo        set_repo.SetRepository
}

func NewPreferencePipe(preferenceRepo preference_repo.PreferenceRepository, personaRepo persona_repo.PersonaRepository, deps ...any) *PreferencePipe {
	pipe := &PreferencePipe{preferenceRepo: preferenceRepo, personaRepo: personaRepo}
	for _, dep := range deps {
		switch typed := dep.(type) {
		case post_repo.PostRepository:
			pipe.postRepo = typed
		case event_repo.EventRepository:
			pipe.eventRepo = typed
		case set_repo.SetRepository:
			pipe.setRepo = typed
		}
	}
	return pipe
}

func pipeInternalError[T any](err error, operation string) *shared.PipeRes[T] {
	return shared.PipeInternalError[T](err, "preference", operation, messages.Internal_Error)
}

func (p *PreferencePipe) ownedPersona(ctx context.Context, userID uuid.UUID, personaID uuid.UUID) (*models.Persona, *shared.PipeRes[any]) {
	persona, err := p.personaRepo.FindPersonaByID(ctx, personaID)
	if err != nil {
		if err == persona_repo.ErrPersonaNotFound {
			return nil, shared.PipeError[any](messages.Persona_Not_Found)
		}
		return nil, pipeInternalError[any](err, "preference.find_persona")
	}
	if persona.UserID != userID {
		return nil, shared.PipeError[any](messages.Forbidden)
	}
	return persona, nil
}

func (p *PreferencePipe) targetPersona(ctx context.Context, personaID uuid.UUID) (*models.Persona, *shared.PipeRes[any]) {
	persona, err := p.personaRepo.FindPersonaByID(ctx, personaID)
	if err != nil {
		if err == persona_repo.ErrPersonaNotFound {
			return nil, shared.PipeError[any](messages.Persona_Not_Found)
		}
		return nil, pipeInternalError[any](err, "preference.find_target_persona")
	}
	return persona, nil
}

func validDiscoveryTargetType(targetType models.DiscoverySuppressionTargetType) bool {
	switch targetType {
	case models.PersonaSuppressionTargetType, models.PostSuppressionTargetType, models.EventSuppressionTargetType, models.SetSuppressionTargetType:
		return true
	default:
		return false
	}
}
