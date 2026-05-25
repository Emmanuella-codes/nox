package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/event/messages"
	"github.com/emmanuella-codes/nox/models"
	event_repo "github.com/emmanuella-codes/nox/repositories/event"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

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
