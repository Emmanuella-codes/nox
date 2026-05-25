package pipes

import (
	"github.com/emmanuella-codes/nox/event/messages"
	event_repo "github.com/emmanuella-codes/nox/repositories/event"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/rs/zerolog/log"
)

type EventPipe struct {
	eventRepo   event_repo.EventRepository
	personaRepo persona_repo.PersonaRepository
}

func NewEventPipe(eventRepo event_repo.EventRepository, personaRepo persona_repo.PersonaRepository) *EventPipe {
	return &EventPipe{eventRepo: eventRepo, personaRepo: personaRepo}
}

func pipeSuccess[T any](message shared.PipeMessage, data *T) *shared.PipeRes[T] {
	return &shared.PipeRes[T]{Success: true, Message: message, Data: data}
}

func pipeError[T any](message shared.PipeMessage) *shared.PipeRes[T] {
	return &shared.PipeRes[T]{Success: false, Message: message}
}

func pipeInternalError[T any](err error, operation string) *shared.PipeRes[T] {
	if err != nil {
		log.Error().Err(err).Str("operation", operation).Msg("event internal error")
	}
	return pipeError[T](messages.Internal_Error)
}
