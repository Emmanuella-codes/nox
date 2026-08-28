package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/follow/messages"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

func (p *FollowPipe) FollowStatusPipe(ctx context.Context, userID uuid.UUID, followerPersonaID uuid.UUID, targetPersonaID uuid.UUID) *shared.PipeRes[FollowStatusResponse] {
	if _, _, res := p.validateFollowAction(ctx, userID, followerPersonaID, targetPersonaID); res != nil {
		return &shared.PipeRes[FollowStatusResponse]{Success: false, Message: res.Message}
	}

	isFollowing, err := p.followRepo.IsFollowing(ctx, followerPersonaID, targetPersonaID)
	if err != nil {
		return pipeInternalError[FollowStatusResponse](err, "follow.status")
	}

	return shared.PipeSuccess(messages.Follow_Status_Fetched, &FollowStatusResponse{IsFollowing: isFollowing})
}
