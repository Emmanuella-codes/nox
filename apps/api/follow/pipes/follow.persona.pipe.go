package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/follow/messages"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

func (p *FollowPipe) FollowPersonaPipe(ctx context.Context, userID uuid.UUID, followerPersonaID uuid.UUID, targetPersonaID uuid.UUID) *shared.PipeRes[any] {
	if res := p.validateFollowAction(ctx, userID, followerPersonaID, targetPersonaID); res != nil {
		return res
	}

	if err := p.followRepo.Follow(ctx, followerPersonaID, targetPersonaID); err != nil {
		return p.mapFollowError(err, "follow.create")
	}

	return shared.PipeSuccess[any](messages.Followed_Successfully, nil)
}
