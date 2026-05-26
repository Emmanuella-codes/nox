package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/event/messages"
	"github.com/emmanuella-codes/nox/models"
	"github.com/emmanuella-codes/nox/shared"
)

func (p *EventPipe) ListEventsPipe(ctx context.Context, limit int) *shared.PipeRes[[]*models.Event] {
	events, err := p.eventRepo.FindEvents(ctx, limit)
	if err != nil {
		return pipeInternalError[[]*models.Event](err, "event.list")
	}
	return shared.PipeSuccess(messages.Events_Listed, &events)
}
