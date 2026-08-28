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

// GetPostPipe fetches one post and hydrates its visible or anonymous author state.
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
	mediaByPost, err := p.postRepo.FindMediaAssetsByPostIDs(ctx, []uuid.UUID{post.ID})
	if err != nil {
		return pipeInternalError[PostResponse](err, "post.media")
	}
	response.Media = mediaByPost[post.ID]
	if p.hashtagRepo != nil {
		tags, err := p.hashtagRepo.FindTagsByPostIDs(ctx, []uuid.UUID{post.ID})
		if err != nil {
			return pipeInternalError[PostResponse](err, "post.get_hashtags")
		}
		response.Hashtags = tags[post.ID]
	}
	return shared.PipeSuccess(messages.Post_Fetched, &response)
}

// GetPostForViewerPipe fetches one post and viewer-specific liked state.
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

// GetPersonaPostsPipe fetches the public profile view of a user's posts.
func (p *PostPipe) GetPersonaPostsPipe(ctx context.Context, personaID uuid.UUID, limit int) *shared.PipeRes[[]PostResponse] {
	persona, err := p.personaRepo.FindPersonaByID(ctx, personaID)
	if err != nil {
		if err == persona_repo.ErrPersonaNotFound {
			return shared.PipeError[[]PostResponse](messages.Persona_Not_Found)
		}
		return pipeInternalError[[]PostResponse](err, "post.find_persona")
	}
	posts, err := p.postRepo.FindPostsByPersonaID(ctx, persona.ID, limit)
	if err != nil {
		return pipeInternalError[[]PostResponse](err, "post.list_by_persona")
	}
	return p.postsResponse(ctx, posts, messages.Posts_Listed)
}

// GetPersonaPostsForViewerPipe fetches owner-aware profile posts, including anonymous posts for the owner.
func (p *PostPipe) GetPersonaPostsForViewerPipe(ctx context.Context, personaID uuid.UUID, viewerPersonaID uuid.UUID, limit int) *shared.PipeRes[[]PostResponse] {
	targetProfile, err := p.personaRepo.FindPersonaByID(ctx, personaID)
	if err != nil {
		if err == persona_repo.ErrPersonaNotFound {
			return shared.PipeError[[]PostResponse](messages.Persona_Not_Found)
		}
		return pipeInternalError[[]PostResponse](err, "post.find_persona")
	}
	viewerProfile, err := p.personaRepo.FindPersonaByID(ctx, viewerPersonaID)
	if err != nil {
		if err == persona_repo.ErrPersonaNotFound {
			return shared.PipeError[[]PostResponse](messages.Persona_Not_Found)
		}
		return pipeInternalError[[]PostResponse](err, "post.find_viewer_persona")
	}
	var posts []*models.Post
	if targetProfile.UserID == viewerProfile.UserID {
		posts, err = p.postRepo.FindPostsByAuthorUserID(ctx, targetProfile.UserID, limit)
	} else {
		posts, err = p.postRepo.FindPostsByPersonaID(ctx, targetProfile.ID, limit)
	}
	if err != nil {
		return pipeInternalError[[]PostResponse](err, "post.list_profile_posts")
	}
	res := p.postsResponse(ctx, posts, messages.Posts_Listed)
	if !res.Success || res.Data == nil || p.likeRepo == nil {
		return res
	}
	if err := p.hydrateLikedState(ctx, viewerPersonaID, *res.Data); err != nil {
		return pipeInternalError[[]PostResponse](err, "post.persona_posts_like_status")
	}
	return res
}

// GetFeedPipe fetches the base feed with liked state.
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
	res := p.postsResponse(ctx, posts, messages.Feed_Listed)
	if !res.Success || res.Data == nil || p.likeRepo == nil {
		return res
	}
	if err := p.hydrateLikedState(ctx, personaID, *res.Data); err != nil {
		return pipeInternalError[[]PostResponse](err, "post.feed_like_status")
	}
	return res
}

// GetFollowingFeedPipe fetches the following feed with liked state.
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
	res := p.postsResponse(ctx, posts, messages.Feed_Listed)
	if !res.Success || res.Data == nil || p.likeRepo == nil {
		return res
	}
	if err := p.hydrateLikedState(ctx, personaID, *res.Data); err != nil {
		return pipeInternalError[[]PostResponse](err, "post.following_feed_like_status")
	}
	return res
}
