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
			return shared.PipeError[models.Event](messages.Event_Not_Found)
		}
		return pipeInternalError[models.Event](err, "event.get")
	}
	return shared.PipeSuccess(messages.Event_Fetched, event)
}

func (p *EventPipe) GetEventForViewerPipe(ctx context.Context, eventID uuid.UUID, viewerPersonaID uuid.UUID) *shared.PipeRes[models.Event] {
	res := p.GetEventPipe(ctx, eventID)
	if !res.Success || res.Data == nil {
		return res
	}
	visible, err := p.viewerCanSeeEvent(ctx, viewerPersonaID, res.Data)
	if err != nil {
		return pipeInternalError[models.Event](err, "event.visibility")
	}
	if !visible {
		return shared.PipeError[models.Event](messages.Event_Not_Found)
	}
	return res
}
