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
)

func (p *EventPipe) UpdateEventPipe(ctx context.Context, userID, eventID uuid.UUID, dto dtos.UpdateEventDTO) *shared.PipeRes[models.Event] {
	dto.Title, dto.Venue, dto.Location, dto.Description = strings.TrimSpace(dto.Title), strings.TrimSpace(dto.Venue), strings.TrimSpace(dto.Location), strings.TrimSpace(dto.Description)
	if dto.EventEndDate.IsZero() || !dto.EventDate.Before(dto.EventEndDate) || dto.Title == "" || dto.Venue == "" || dto.Location == "" || dto.Description == "" {
		return shared.PipeError[models.Event](messages.Invalid_Payload)
	}
	current, err := p.eventRepo.FindEventByID(ctx, eventID)
	if err != nil {
		return shared.PipeError[models.Event](messages.Event_Not_Found)
	}
	organizer, err := p.personaRepo.FindPersonaByID(ctx, current.OrganizerID)
	if err != nil || organizer.UserID != userID {
		return shared.PipeError[models.Event](messages.Forbidden)
	}
	event, err := p.eventRepo.UpdateEvent(ctx, userID, eventID, dto)
	if err != nil {
		return pipeLifecycleError[models.Event](err)
	}
	return shared.PipeSuccess(messages.Event_Updated, event)
}

func (p *EventPipe) transitionEvent(ctx context.Context, userID, eventID uuid.UUID, action string) *shared.PipeRes[models.Event] {
	var event *models.Event
	var err error
	switch action {
	case "publish":
		event, err = p.eventRepo.PublishEvent(ctx, userID, eventID)
	case "close":
		event, err = p.eventRepo.CloseEvent(ctx, userID, eventID)
	case "archive":
		event, err = p.eventRepo.ArchiveEvent(ctx, userID, eventID)
	}
	if err != nil {
		return pipeLifecycleError[models.Event](err)
	}
	message := messages.Event_Published
	if action == "close" {
		message = messages.Event_Closed
	}
	if action == "archive" {
		message = messages.Event_Archived
	}
	return shared.PipeSuccess(message, event)
}

func (p *EventPipe) PublishEventPipe(ctx context.Context, userID, eventID uuid.UUID) *shared.PipeRes[models.Event] {
	return p.transitionEvent(ctx, userID, eventID, "publish")
}
func (p *EventPipe) CloseEventPipe(ctx context.Context, userID, eventID uuid.UUID) *shared.PipeRes[models.Event] {
	return p.transitionEvent(ctx, userID, eventID, "close")
}
func (p *EventPipe) ArchiveEventPipe(ctx context.Context, userID, eventID uuid.UUID) *shared.PipeRes[models.Event] {
	return p.transitionEvent(ctx, userID, eventID, "archive")
}

func pipeLifecycleError[T any](err error) *shared.PipeRes[T] {
	if err == event_repo.ErrEventNotFound || err == persona_repo.ErrPersonaNotFound {
		return shared.PipeError[T](messages.Event_Not_Found)
	}
	if err == event_repo.ErrInvalidEventTransition {
		return shared.PipeError[T](messages.Invalid_Transition)
	}
	return pipeInternalError[T](err, "event.lifecycle")
}
