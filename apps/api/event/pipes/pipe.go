package pipes

import (
	"github.com/emmanuella-codes/nox/event/messages"
	event_repo "github.com/emmanuella-codes/nox/repositories/event"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	"github.com/emmanuella-codes/nox/shared"
)

type EventPipe struct {
	eventRepo   event_repo.EventRepository
	personaRepo persona_repo.PersonaRepository
}

func NewEventPipe(eventRepo event_repo.EventRepository, personaRepo persona_repo.PersonaRepository) *EventPipe {
	return &EventPipe{eventRepo: eventRepo, personaRepo: personaRepo}
}

func pipeInternalError[T any](err error, operation string) *shared.PipeRes[T] {
	return shared.PipeInternalError[T](err, "event", operation, messages.Internal_Error)
}
