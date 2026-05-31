package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/follow/messages"
	follow_repo "github.com/emmanuella-codes/nox/repositories/follow"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

func (p *FollowPipe) FollowingPipe(ctx context.Context, personaID uuid.UUID, options follow_repo.ListOptions) *shared.PipeRes[FollowListResponse] {
	if res := p.validateVisiblePersona(ctx, personaID); res != nil {
		return &shared.PipeRes[FollowListResponse]{Success: false, Message: res.Message}
	}

	options = follow_repo.NormalizeListOptions(options)
	following, err := p.followRepo.FindFollowing(ctx, personaID, options)
	if err != nil {
		return pipeInternalError[FollowListResponse](err, "follow.following")
	}

	response := followListResponse(following, options)
	return shared.PipeSuccess(messages.Following_Listed, &response)
}
