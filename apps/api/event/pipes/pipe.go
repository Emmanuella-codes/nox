package pipes

import (
	"context"
	"strings"

	"github.com/emmanuella-codes/nox/event/dtos"
	"github.com/emmanuella-codes/nox/event/messages"
	"github.com/emmanuella-codes/nox/models"
	event_repo "github.com/emmanuella-codes/nox/repositories/event"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type EventPipe struct {
	eventRepo   event_repo.EventRepository
	personaRepo persona_repo.PersonaRepository
}

func NewEventPipe(eventRepo event_repo.EventRepository, personaRepo persona_repo.PersonaRepository) *EventPipe {
	return &EventPipe{eventRepo: eventRepo, personaRepo: personaRepo}
}

func (p *EventPipe) CreateEventPipe(ctx context.Context, userID uuid.UUID, dto dtos.CreateEventDTO) *shared.PipeRes[models.Event] {
	dto.Title = strings.TrimSpace(dto.Title)
	dto.Venue = strings.TrimSpace(dto.Venue)
	dto.Location = strings.TrimSpace(dto.Location)
	dto.Description = strings.TrimSpace(dto.Description)
	dto.CoverURL = strings.TrimSpace(dto.CoverURL)
	dto.TicketURL = strings.TrimSpace(dto.TicketURL)

	persona, err := p.personaRepo.FindPersonaByID(ctx, dto.OrganizerID)
	if err != nil {
		if err == persona_repo.ErrPersonaNotFound {
			return pipeError[models.Event](messages.Persona_Not_Found)
		}
		return pipeInternalError[models.Event](err, "event.find_organizer")
	}
	if persona.UserID != userID || persona.PersonaType != models.VisiblePersonaType {
		return pipeError[models.Event](messages.Forbidden)
	}

	event, err := p.eventRepo.CreateEvent(ctx, dto)
	if err != nil {
		return pipeInternalError[models.Event](err, "event.create")
	}
	return pipeSuccess(messages.Event_Created, event)
}

func (p *EventPipe) GetEventPipe(ctx context.Context, eventID uuid.UUID) *shared.PipeRes[models.Event] {
	event, err := p.eventRepo.FindEventByID(ctx, eventID)
	if err != nil {
		if err == event_repo.ErrEventNotFound {
			return pipeError[models.Event](messages.Event_Not_Found)
		}
		return pipeInternalError[models.Event](err, "event.get")
	}
	return pipeSuccess(messages.Event_Fetched, event)
}

func (p *EventPipe) ListEventsPipe(ctx context.Context, limit int) *shared.PipeRes[[]*models.Event] {
	events, err := p.eventRepo.FindEvents(ctx, limit)
	if err != nil {
		return pipeInternalError[[]*models.Event](err, "event.list")
	}
	return pipeSuccess(messages.Events_Listed, &events)
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
