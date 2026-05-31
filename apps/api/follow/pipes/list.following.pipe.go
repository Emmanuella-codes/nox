package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/follow/messages"
	follow_repo "github.com/emmanuella-codes/nox/repositories/follow"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

func (p *FollowPipe) FollowingPipe(ctx context.Context, personaID uuid.UUID, options follow_repo.ListOptions) *shared.PipeRes[FollowListResponse] {
	return p.followingPipe(ctx, personaID, options, nil)
}

func (p *FollowPipe) FollowingForViewerPipe(ctx context.Context, userID uuid.UUID, personaID uuid.UUID, viewerPersonaID uuid.UUID, options follow_repo.ListOptions) *shared.PipeRes[FollowListResponse] {
	viewerPersona, res := p.validateOwnedVisiblePersona(ctx, userID, viewerPersonaID)
	if res != nil {
		return &shared.PipeRes[FollowListResponse]{Success: false, Message: res.Message}
	}
	return p.followingPipe(ctx, personaID, options, &viewerPersona.ID)
}

func (p *FollowPipe) followingPipe(ctx context.Context, personaID uuid.UUID, options follow_repo.ListOptions, viewerPersonaID *uuid.UUID) *shared.PipeRes[FollowListResponse] {
	if res := p.validateVisiblePersona(ctx, personaID); res != nil {
		return &shared.PipeRes[FollowListResponse]{Success: false, Message: res.Message}
	}

	options = follow_repo.NormalizeListOptions(options)
	following, err := p.followRepo.FindFollowing(ctx, personaID, options)
	if err != nil {
		return pipeInternalError[FollowListResponse](err, "follow.following")
	}

	followingState := map[uuid.UUID]bool{}
	if viewerPersonaID != nil {
		var err error
		followingState, err = p.followRepo.FindFollowingIDs(ctx, *viewerPersonaID, personaIDs(following))
		if err != nil {
			return pipeInternalError[FollowListResponse](err, "follow.following_status")
		}
	}
	response := followListResponse(following, options, followingState)
	return shared.PipeSuccess(messages.Following_Listed, &response)
}
