package pipes

import (
	"context"

	set_repo "github.com/emmanuella-codes/nox/repositories/set"
	"github.com/emmanuella-codes/nox/set/messages"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

func (p *SetPipe) DeleteSetPipe(ctx context.Context, userID uuid.UUID, setID uuid.UUID) *shared.PipeRes[any] {
	set, err := p.setRepo.FindSetByID(ctx, setID)
	if err != nil {
		if err == set_repo.ErrSetNotFound {
			return shared.PipeError[any](messages.Set_Not_Found)
		}
		return pipeInternalError[any](err, "set.find_for_delete")
	}
	if set.AuthorUserID != userID {
		return shared.PipeError[any](messages.Forbidden)
	}
	if err := p.setRepo.DeleteSet(ctx, setID); err != nil {
		if err == set_repo.ErrSetNotFound {
			return shared.PipeError[any](messages.Set_Not_Found)
		}
		return pipeInternalError[any](err, "set.delete")
	}
	return shared.PipeSuccess[any](messages.Set_Deleted, nil)
}
