package pipes

import (
	"github.com/emmanuella-codes/nox/event/messages"
	event_repo "github.com/emmanuella-codes/nox/repositories/event"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	preference_repo "github.com/emmanuella-codes/nox/repositories/preference"
	"github.com/emmanuella-codes/nox/shared"
)

type EventPipe struct {
	eventRepo      event_repo.EventRepository
	personaRepo    persona_repo.PersonaRepository
	preferenceRepo preference_repo.PreferenceRepository
}

func NewEventPipe(eventRepo event_repo.EventRepository, personaRepo persona_repo.PersonaRepository, deps ...any) *EventPipe {
	pipe := &EventPipe{eventRepo: eventRepo, personaRepo: personaRepo}
	for _, dep := range deps {
		if preferenceRepo, ok := dep.(preference_repo.PreferenceRepository); ok {
			pipe.preferenceRepo = preferenceRepo
		}
	}
	return pipe
}

func pipeInternalError[T any](err error, operation string) *shared.PipeRes[T] {
	return shared.PipeInternalError[T](err, "event", operation, messages.Internal_Error)
}
