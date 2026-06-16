package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/models"
	"github.com/emmanuella-codes/nox/post/messages"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	post_repo "github.com/emmanuella-codes/nox/repositories/post"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

func (p *PostPipe) GetPostPipe(ctx context.Context, postID uuid.UUID) *shared.PipeRes[PostResponse] {
	post, err := p.postRepo.FindPostByID(ctx, postID)
	if err != nil {
		if err == post_repo.ErrPostNotFound {
			return shared.PipeError[PostResponse](messages.Post_Not_Found)
		}
		return pipeInternalError[PostResponse](err, "post.get")
	}

	persona, pipeErr := p.publicPostPersona(ctx, post)
	if pipeErr != nil {
		return pipeErr
	}
	identity, err := p.anonymousPostIdentity(ctx, post)
	if err != nil {
		return pipeInternalError[PostResponse](err, "post.anonymous_identity")
	}

	response := postResponse(post, persona, identity)
	if p.hashtagRepo != nil {
		tags, err := p.hashtagRepo.FindTagsByPostIDs(ctx, []uuid.UUID{post.ID})
		if err != nil {
			return pipeInternalError[PostResponse](err, "post.get_hashtags")
		}
		response.Hashtags = tags[post.ID]
	}
	return shared.PipeSuccess(messages.Post_Fetched, &response)
}

func (p *PostPipe) GetPostForViewerPipe(ctx context.Context, postID uuid.UUID, viewerPersonaID uuid.UUID) *shared.PipeRes[PostResponse] {
	res := p.GetPostPipe(ctx, postID)
	if !res.Success || res.Data == nil || p.likeRepo == nil {
		return res
	}

	liked, err := p.likeRepo.HasPostLike(ctx, viewerPersonaID, postID)
	if err != nil {
		return pipeInternalError[PostResponse](err, "post.viewer_like_status")
	}
	res.Data.IsLiked = liked
	return res
}

func (p *PostPipe) GetPersonaPostsPipe(ctx context.Context, personaID uuid.UUID, limit int) *shared.PipeRes[[]PostResponse] {
	if _, err := p.personaRepo.FindPersonaByID(ctx, personaID); err != nil {
		if err == persona_repo.ErrPersonaNotFound {
			return shared.PipeError[[]PostResponse](messages.Persona_Not_Found)
		}
		return pipeInternalError[[]PostResponse](err, "post.find_persona")
	}

	posts, err := p.postRepo.FindPostsByPersonaID(ctx, personaID, limit)
	if err != nil {
		return pipeInternalError[[]PostResponse](err, "post.list_by_persona")
	}

	personas, pipeErr := p.publicPostPersonas(ctx, posts)
	if pipeErr != nil {
		return pipeErr
	}

	identities, err := p.anonymousPostIdentities(ctx, posts)
	if err != nil {
		return pipeInternalError[[]PostResponse](err, "post.anonymous_identities")
	}

	responses := postResponses(posts, personas, identities)
	if err := p.hydrateHashtags(ctx, responses); err != nil {
		return pipeInternalError[[]PostResponse](err, "post.persona_posts_hashtags")
	}
	return shared.PipeSuccess(messages.Posts_Listed, &responses)
}

func (p *PostPipe) GetPersonaPostsForViewerPipe(ctx context.Context, personaID uuid.UUID, viewerPersonaID uuid.UUID, limit int) *shared.PipeRes[[]PostResponse] {
	res := p.GetPersonaPostsPipe(ctx, personaID, limit)
	if !res.Success || res.Data == nil || p.likeRepo == nil {
		return res
	}

	if err := p.hydrateLikedState(ctx, viewerPersonaID, *res.Data); err != nil {
		return pipeInternalError[[]PostResponse](err, "post.persona_posts_like_status")
	}
	return res
}

func (p *PostPipe) GetFeedPipe(ctx context.Context, personaID uuid.UUID, limit int) *shared.PipeRes[[]PostResponse] {
	if _, err := p.personaRepo.FindPersonaByID(ctx, personaID); err != nil {
		if err == persona_repo.ErrPersonaNotFound {
			return shared.PipeError[[]PostResponse](messages.Persona_Not_Found)
		}
		return pipeInternalError[[]PostResponse](err, "post.find_persona_for_feed")
	}

	posts, err := p.postRepo.FindFeedPosts(ctx, personaID, limit)
	if err != nil {
		return pipeInternalError[[]PostResponse](err, "post.feed")
	}

	personas, pipeErr := p.publicPostPersonas(ctx, posts)
	if pipeErr != nil {
		return pipeErr
	}

	identities, err := p.anonymousPostIdentities(ctx, posts)
	if err != nil {
		return pipeInternalError[[]PostResponse](err, "post.feed_anonymous_identities")
	}

	responses := postResponses(posts, personas, identities)
	if err := p.hydrateHashtags(ctx, responses); err != nil {
		return pipeInternalError[[]PostResponse](err, "post.feed_hashtags")
	}
	if p.likeRepo != nil {
		if err := p.hydrateLikedState(ctx, personaID, responses); err != nil {
			return pipeInternalError[[]PostResponse](err, "post.feed_like_status")
		}
	}
	return shared.PipeSuccess(messages.Feed_Listed, &responses)
}

func (p *PostPipe) GetFollowingFeedPipe(ctx context.Context, personaID uuid.UUID, limit int) *shared.PipeRes[[]PostResponse] {
	if _, err := p.personaRepo.FindPersonaByID(ctx, personaID); err != nil {
		if err == persona_repo.ErrPersonaNotFound {
			return shared.PipeError[[]PostResponse](messages.Persona_Not_Found)
		}
		return pipeInternalError[[]PostResponse](err, "post.find_persona_for_following_feed")
	}

	posts, err := p.postRepo.FindFollowingFeedPosts(ctx, personaID, limit)
	if err != nil {
		return pipeInternalError[[]PostResponse](err, "post.following_feed")
	}

	personas, pipeErr := p.publicPostPersonas(ctx, posts)
	if pipeErr != nil {
		return pipeErr
	}

	identities, err := p.anonymousPostIdentities(ctx, posts)
	if err != nil {
		return pipeInternalError[[]PostResponse](err, "post.following_feed_anonymous_identities")
	}

	responses := postResponses(posts, personas, identities)
	if err := p.hydrateHashtags(ctx, responses); err != nil {
		return pipeInternalError[[]PostResponse](err, "post.following_feed_hashtags")
	}
	if p.likeRepo != nil {
		if err := p.hydrateLikedState(ctx, personaID, responses); err != nil {
			return pipeInternalError[[]PostResponse](err, "post.following_feed_like_status")
		}
	}
	return shared.PipeSuccess(messages.Feed_Listed, &responses)
}

func (p *PostPipe) FindViewerPersona(ctx context.Context, userID uuid.UUID, personaID uuid.UUID) (*models.Persona, shared.PipeMessage) {
	persona, err := p.personaRepo.FindPersonaByID(ctx, personaID)
	if err != nil {
		if err == persona_repo.ErrPersonaNotFound {
			return nil, messages.Persona_Not_Found
		}
		return nil, messages.Internal_Error
	}
	if persona.UserID != userID || persona.PersonaType != models.VisiblePersonaType {
		return nil, messages.Forbidden
	}
	return persona, ""
}

func (p *PostPipe) publicPostPersona(ctx context.Context, post *models.Post) (*models.Persona, *shared.PipeRes[PostResponse]) {
	if post.PostingMode != models.PublicPostingMode || post.PersonaID == nil {
		return nil, nil
	}

	persona, err := p.personaRepo.FindPersonaByID(ctx, *post.PersonaID)
	if err != nil {
		if err == persona_repo.ErrPersonaNotFound {
			return nil, shared.PipeError[PostResponse](messages.Persona_Not_Found)
		}
		return nil, pipeInternalError[PostResponse](err, "post.find_public_persona")
	}

	return persona, nil
}

func (p *PostPipe) publicPostPersonas(ctx context.Context, posts []*models.Post) (map[string]*models.Persona, *shared.PipeRes[[]PostResponse]) {
	personas := make(map[string]*models.Persona)
	for _, post := range posts {
		if post.PostingMode != models.PublicPostingMode || post.PersonaID == nil {
			continue
		}
		personaID := post.PersonaID.String()
		if _, ok := personas[personaID]; ok {
			continue
		}

		persona, err := p.personaRepo.FindPersonaByID(ctx, *post.PersonaID)
		if err != nil {
			if err == persona_repo.ErrPersonaNotFound {
				return nil, shared.PipeError[[]PostResponse](messages.Persona_Not_Found)
			}
			return nil, pipeInternalError[[]PostResponse](err, "post.find_public_persona")
		}
		personas[personaID] = persona
	}

	return personas, nil
}

func (p *PostPipe) anonymousPostIdentity(ctx context.Context, post *models.Post) (*models.AnonymousThreadIdentity, error) {
	if post.PostingMode != models.AnonymousPostingMode || post.PersonaID == nil {
		return nil, nil
	}
	return p.postRepo.EnsureAnonymousThreadIdentity(ctx, post.ID, post.AuthorUserID, *post.PersonaID, anonymousHandle())
}

func (p *PostPipe) anonymousPostIdentities(ctx context.Context, posts []*models.Post) (map[uuid.UUID]*models.AnonymousThreadIdentity, error) {
	identities := make(map[uuid.UUID]*models.AnonymousThreadIdentity)
	for _, post := range posts {
		identity, err := p.anonymousPostIdentity(ctx, post)
		if err != nil {
			return nil, err
		}
		if identity != nil {
			identities[post.ID] = identity
		}
	}
	return identities, nil
}

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
		postID := postIDs[i]
		responses[i].IsLiked = liked[postID]
	}
	return nil
}

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
