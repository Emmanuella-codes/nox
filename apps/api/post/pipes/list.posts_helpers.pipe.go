package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/models"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

// postsResponse hydrates profiles, anonymous identities, media, and hashtags for a post list.
func (p *PostPipe) postsResponse(ctx context.Context, posts []*models.Post, message shared.PipeMessage) *shared.PipeRes[[]PostResponse] {
	personas, pipeErr := p.publicPostPersonas(ctx, posts)
	if pipeErr != nil {
		return pipeErr
	}
	identities, err := p.anonymousPostIdentities(ctx, posts)
	if err != nil {
		return pipeInternalError[[]PostResponse](err, "post.anonymous_identities")
	}
	mediaByPost, err := p.postRepo.FindMediaAssetsByPostIDs(ctx, postIDs(posts))
	if err != nil {
		return pipeInternalError[[]PostResponse](err, "post.media")
	}
	responses := make([]PostResponse, 0, len(posts))
	for _, post := range posts {
		var persona *models.Persona
		if post.PersonaID != nil {
			persona = personas[post.PersonaID.String()]
		}
		response := postResponse(post, publicAuthorProfile(post, persona), identities[post.ID])
		response.Media = mediaByPost[post.ID]
		responses = append(responses, response)
	}
	if err := p.hydrateHashtags(ctx, responses); err != nil {
		return pipeInternalError[[]PostResponse](err, "post.hashtags")
	}
	return shared.PipeSuccess(message, &responses)
}

// hydrateLikedState adds viewer-specific like status to each post response.
func (p *PostPipe) hydrateLikedState(ctx context.Context, viewerPersonaID uuid.UUID, responses []PostResponse) error {
	if p.likeRepo == nil || len(responses) == 0 {
		return nil
	}
	postIDs := make([]uuid.UUID, 0, len(responses))
	for _, response := range responses {
		postID, err := uuid.Parse(response.ID)
		if err != nil {
			return err
		}
		postIDs = append(postIDs, postID)
	}
	liked, err := p.likeRepo.FindLikedPostIDs(ctx, viewerPersonaID, postIDs)
	if err != nil {
		return err
	}
	for i := range responses {
		responses[i].IsLiked = liked[postIDs[i]]
	}
	return nil
}

// hydrateHashtags adds hydrated hashtag lists to each post response.
func (p *PostPipe) hydrateHashtags(ctx context.Context, responses []PostResponse) error {
	if p.hashtagRepo == nil || len(responses) == 0 {
		return nil
	}
	ids := make([]uuid.UUID, 0, len(responses))
	for _, response := range responses {
		postID, err := uuid.Parse(response.ID)
		if err != nil {
			return err
		}
		ids = append(ids, postID)
	}
	tags, err := p.hashtagRepo.FindTagsByPostIDs(ctx, ids)
	if err != nil {
		return err
	}
	for i := range responses {
		responses[i].Hashtags = tags[ids[i]]
	}
	return nil
}
