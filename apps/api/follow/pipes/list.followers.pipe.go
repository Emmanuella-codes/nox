package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/follow/messages"
	follow_repo "github.com/emmanuella-codes/nox/repositories/follow"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

func (p *FollowPipe) FollowersPipe(ctx context.Context, personaID uuid.UUID, options follow_repo.ListOptions) *shared.PipeRes[FollowListResponse] {
	if res := p.validateVisiblePersona(ctx, personaID); res != nil {
		return &shared.PipeRes[FollowListResponse]{Success: false, Message: res.Message}
	}

	options = follow_repo.NormalizeListOptions(options)
	followers, err := p.followRepo.FindFollowers(ctx, personaID, options)
	if err != nil {
		return pipeInternalError[FollowListResponse](err, "follow.followers")
	}

	response := followListResponse(followers, options)
	return shared.PipeSuccess(messages.Followers_Listed, &response)
}
