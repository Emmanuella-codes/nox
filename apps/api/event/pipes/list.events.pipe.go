package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/event/messages"
	"github.com/emmanuella-codes/nox/models"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

func (p *EventPipe) ListEventsPipe(ctx context.Context, limit int) *shared.PipeRes[[]*models.Event] {
	events, err := p.eventRepo.FindEvents(ctx, limit)
	if err != nil {
		return pipeInternalError[[]*models.Event](err, "event.list")
	}
	return shared.PipeSuccess(messages.Events_Listed, &events)
}

func (p *EventPipe) ListEventsForViewerPipe(ctx context.Context, viewerPersonaID uuid.UUID, limit int) *shared.PipeRes[[]*models.Event] {
	events, err := p.eventRepo.FindEvents(ctx, limit)
	if err != nil {
		return pipeInternalError[[]*models.Event](err, "event.list")
	}
	filtered := make([]*models.Event, 0, len(events))
	for _, event := range events {
		visible, err := p.viewerCanSeeEvent(ctx, viewerPersonaID, event)
		if err != nil {
			return pipeInternalError[[]*models.Event](err, "event.list_visibility")
		}
		if visible {
			filtered = append(filtered, event)
		}
	}
	return shared.PipeSuccess(messages.Events_Listed, &filtered)
}
