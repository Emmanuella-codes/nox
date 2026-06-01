package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/models"
	set_repo "github.com/emmanuella-codes/nox/repositories/set"
	"github.com/emmanuella-codes/nox/set/messages"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

func (p *SetPipe) GetSetPipe(ctx context.Context, setID uuid.UUID) *shared.PipeRes[models.Set] {
	set, err := p.setRepo.FindSetByID(ctx, setID)
	if err != nil {
		if err == set_repo.ErrSetNotFound {
			return shared.PipeError[models.Set](messages.Set_Not_Found)
		}
		return pipeInternalError[models.Set](err, "set.get")
	}
	if err := p.hydrateSet(ctx, set); err != nil {
		return pipeInternalError[models.Set](err, "set.hydrate")
	}
	return shared.PipeSuccess(messages.Set_Fetched, set)
}
