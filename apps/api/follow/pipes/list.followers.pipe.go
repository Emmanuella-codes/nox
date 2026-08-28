package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/follow/messages"
	follow_repo "github.com/emmanuella-codes/nox/repositories/follow"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

// FollowersPipe lists followers for a public profile.
func (p *FollowPipe) FollowersPipe(ctx context.Context, personaID uuid.UUID, options follow_repo.ListOptions) *shared.PipeRes[FollowListResponse] {
	return p.followersPipe(ctx, personaID, options, nil)
}

// FollowersForViewerPipe lists followers and annotates viewer follow state.
func (p *FollowPipe) FollowersForViewerPipe(ctx context.Context, userID uuid.UUID, personaID uuid.UUID, viewerPersonaID uuid.UUID, options follow_repo.ListOptions) *shared.PipeRes[FollowListResponse] {
	viewerPersona, res := p.validateOwnedProfile(ctx, userID, viewerPersonaID)
	if res != nil {
		return &shared.PipeRes[FollowListResponse]{Success: false, Message: res.Message}
	}
	return p.followersPipe(ctx, personaID, options, &viewerPersona.ID)
}

// followersPipe loads followers and optional viewer follow state.
func (p *FollowPipe) followersPipe(ctx context.Context, personaID uuid.UUID, options follow_repo.ListOptions, viewerPersonaID *uuid.UUID) *shared.PipeRes[FollowListResponse] {
	if res := p.validateProfile(ctx, personaID); res != nil {
		return &shared.PipeRes[FollowListResponse]{Success: false, Message: res.Message}
	}

	options = follow_repo.NormalizeListOptions(options)
	followers, err := p.followRepo.FindFollowers(ctx, personaID, options)
	if err != nil {
		return pipeInternalError[FollowListResponse](err, "follow.followers")
	}

	following := map[uuid.UUID]bool{}
	if viewerPersonaID != nil {
		var err error
		following, err = p.followRepo.FindFollowingIDs(ctx, *viewerPersonaID, personaIDs(followers))
		if err != nil {
			return pipeInternalError[FollowListResponse](err, "follow.followers_status")
		}
	}
	response := followListResponse(followers, options, following)
	return shared.PipeSuccess(messages.Followers_Listed, &response)
}
