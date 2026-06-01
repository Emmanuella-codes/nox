package pipes

import (
	"context"

	set_repo "github.com/emmanuella-codes/nox/repositories/set"
	"github.com/emmanuella-codes/nox/set/messages"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

func (p *SetPipe) GetSetPipe(ctx context.Context, setID uuid.UUID, viewerPersonaID *uuid.UUID) *shared.PipeRes[SetResponse] {
	set, err := p.setRepo.FindSetByID(ctx, setID)
	if err != nil {
		if err == set_repo.ErrSetNotFound {
			return shared.PipeError[SetResponse](messages.Set_Not_Found)
		}
		return pipeInternalError[SetResponse](err, "set.get")
	}
	if err := p.hydrateSet(ctx, set); err != nil {
		return pipeInternalError[SetResponse](err, "set.hydrate")
	}
	response, err := p.setResponse(ctx, set, viewerPersonaID)
	if err != nil {
		return pipeInternalError[SetResponse](err, "set.viewer_state")
	}
	return shared.PipeSuccess(messages.Set_Fetched, response)
}
