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
	if !res.Success || res.Data == nil {
		return res
	}
	post, err := p.postRepo.FindPostByID(ctx, postID)
	if err != nil {
		if err == post_repo.ErrPostNotFound {
			return shared.PipeError[PostResponse](messages.Post_Not_Found)
		}
		return pipeInternalError[PostResponse](err, "post.viewer_lookup")
	}
	blocked, err := p.blockedForViewer(ctx, viewerPersonaID, post.AuthorUserID)
	if err != nil {
		return pipeInternalError[PostResponse](err, "post.viewer_visibility")
	}
	if blocked {
		return shared.PipeError[PostResponse](messages.Post_Not_Found)
	}
	if p.likeRepo == nil {
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
	blocked, err := p.blockedForViewer(ctx, viewerPersonaID, targetProfile.UserID)
	if err != nil {
		return pipeInternalError[[]PostResponse](err, "post.profile_visibility")
	}
	if targetProfile.UserID != viewerProfile.UserID && blocked {
		return shared.PipeError[[]PostResponse](messages.Persona_Not_Found)
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

// GetFeedPipe fetches the base feed with cursor pagination and liked state.
func (p *PostPipe) GetFeedPipe(ctx context.Context, personaID uuid.UUID, limit int, cursor string) *shared.PipeRes[PostListResponse] {
	if _, err := p.personaRepo.FindPersonaByID(ctx, personaID); err != nil {
		if err == persona_repo.ErrPersonaNotFound {
			return shared.PipeError[PostListResponse](messages.Persona_Not_Found)
		}
		return pipeInternalError[PostListResponse](err, "post.find_persona_for_feed")
	}
	options, err := p.feedOptions(limit, cursor)
	if err != nil {
		return shared.PipeError[PostListResponse](messages.Invalid_Payload)
	}
	if cached := p.cachedFeedResponse(ctx, "all", personaID, options, messages.Feed_Listed); cached != nil {
		return cached
	}
	posts, err := p.postRepo.FindFeedPosts(ctx, personaID, options)
	if err != nil {
		return pipeInternalError[PostListResponse](err, "post.feed")
	}
	response, pipeRes := p.feedResponse(ctx, posts, personaID, options.Limit, messages.Feed_Listed)
	if pipeRes != nil {
		return pipeRes
	}
	if err := p.storeFeedCache(ctx, "all", personaID.String(), options.Limit, cursor, response); err != nil {
		return pipeInternalError[PostListResponse](err, "post.feed_cache_store")
	}
	return shared.PipeSuccess(messages.Feed_Listed, response)
}

// GetFollowingFeedPipe fetches the following feed with cursor pagination and liked state.
func (p *PostPipe) GetFollowingFeedPipe(ctx context.Context, personaID uuid.UUID, limit int, cursor string) *shared.PipeRes[PostListResponse] {
	if _, err := p.personaRepo.FindPersonaByID(ctx, personaID); err != nil {
		if err == persona_repo.ErrPersonaNotFound {
			return shared.PipeError[PostListResponse](messages.Persona_Not_Found)
		}
		return pipeInternalError[PostListResponse](err, "post.find_persona_for_following_feed")
	}
	options, err := p.feedOptions(limit, cursor)
	if err != nil {
		return shared.PipeError[PostListResponse](messages.Invalid_Payload)
	}
	if cached := p.cachedFeedResponse(ctx, "following", personaID, options, messages.Feed_Listed); cached != nil {
		return cached
	}
	posts, err := p.postRepo.FindFollowingFeedPosts(ctx, personaID, options)
	if err != nil {
		return pipeInternalError[PostListResponse](err, "post.following_feed")
	}
	response, pipeRes := p.feedResponse(ctx, posts, personaID, options.Limit, messages.Feed_Listed)
	if pipeRes != nil {
		return pipeRes
	}
	if err := p.storeFeedCache(ctx, "following", personaID.String(), options.Limit, cursor, response); err != nil {
		return pipeInternalError[PostListResponse](err, "post.following_feed_cache_store")
	}
	return shared.PipeSuccess(messages.Feed_Listed, response)
}

// feedOptions builds validated feed options from incoming request parameters.
func (p *PostPipe) feedOptions(limit int, cursor string) (post_repo.FeedOptions, error) {
	decodedCursor, err := decodeFeedCursor(cursor)
	if err != nil {
		return post_repo.FeedOptions{}, err
	}
	return post_repo.NormalizeFeedOptions(post_repo.FeedOptions{Limit: limit, Cursor: decodedCursor}), nil
}

// cachedFeedResponse returns one cached feed response when available.
func (p *PostPipe) cachedFeedResponse(ctx context.Context, kind string, personaID uuid.UUID, options post_repo.FeedOptions, message shared.PipeMessage) *shared.PipeRes[PostListResponse] {
	var cached PostListResponse
	cursor := ""
	if options.Cursor != nil {
		var err error
		cursor, err = encodeFeedCursor(PostResponse{ID: options.Cursor.ID.String(), CreatedAt: options.Cursor.CreatedAt})
		if err != nil {
			return pipeInternalError[PostListResponse](err, "post.feed_cache_cursor")
		}
	}
	hit, err := p.loadFeedCache(ctx, kind, personaID.String(), options.Limit, cursor, &cached)
	if err != nil {
		return pipeInternalError[PostListResponse](err, "post.feed_cache_load")
	}
	if !hit {
		return nil
	}
	return shared.PipeSuccess(message, &cached)
}

// feedResponse maps one feed page into a cursor-paginated response body.
func (p *PostPipe) feedResponse(ctx context.Context, posts []*models.Post, personaID uuid.UUID, limit int, message shared.PipeMessage) (*PostListResponse, *shared.PipeRes[PostListResponse]) {
	hasMore := len(posts) > limit
	if hasMore {
		posts = posts[:limit]
	}
	res := p.postsResponse(ctx, posts, message)
	if !res.Success || res.Data == nil {
		if res.Success {
			return nil, shared.PipeSuccess(message, &PostListResponse{Limit: limit, HasMore: false, Posts: []PostResponse{}})
		}
		return nil, &shared.PipeRes[PostListResponse]{Success: false, Message: res.Message}
	}
	if err := p.hydrateLikedState(ctx, personaID, *res.Data); err != nil {
		return nil, pipeInternalError[PostListResponse](err, "post.feed_like_status")
	}
	response := &PostListResponse{Limit: limit, HasMore: hasMore, Posts: *res.Data}
	if hasMore && len(response.Posts) > 0 {
		cursor, err := encodeFeedCursor(response.Posts[len(response.Posts)-1])
		if err != nil {
			return nil, pipeInternalError[PostListResponse](err, "post.feed_next_cursor")
		}
		response.NextCursor = &cursor
	}
	return response, nil
}
