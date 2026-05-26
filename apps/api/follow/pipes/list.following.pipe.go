package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/follow/messages"
	"github.com/emmanuella-codes/nox/models"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

func (p *FollowPipe) FollowingPipe(ctx context.Context, personaID uuid.UUID, limit int) *shared.PipeRes[[]*models.Persona] {
	if res := p.validateVisiblePersona(ctx, personaID); res != nil {
		return &shared.PipeRes[[]*models.Persona]{Success: false, Message: res.Message}
	}

	following, err := p.followRepo.FindFollowing(ctx, personaID, limit)
	if err != nil {
		return pipeInternalError[[]*models.Persona](err, "follow.following")
	}

	return shared.PipeSuccess(messages.Following_Listed, &following)
}
