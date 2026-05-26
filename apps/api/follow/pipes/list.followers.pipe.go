package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/follow/messages"
	"github.com/emmanuella-codes/nox/models"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

func (p *FollowPipe) FollowersPipe(ctx context.Context, personaID uuid.UUID, limit int) *shared.PipeRes[[]*models.Persona] {
	if res := p.validateVisiblePersona(ctx, personaID); res != nil {
		return &shared.PipeRes[[]*models.Persona]{Success: false, Message: res.Message}
	}

	followers, err := p.followRepo.FindFollowers(ctx, personaID, limit)
	if err != nil {
		return pipeInternalError[[]*models.Persona](err, "follow.followers")
	}

	return shared.PipeSuccess(messages.Followers_Listed, &followers)
}
