package pipes

import (
	"context"
	"strings"

	"github.com/emmanuella-codes/nox/models"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	searchrepo "github.com/emmanuella-codes/nox/repositories/search"
	searchmessages "github.com/emmanuella-codes/nox/search/messages"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/emmanuella-codes/nox/shared/anonymousidentity"
	"github.com/google/uuid"
)

type SearchOptions = searchrepo.Options

// Search runs search without viewer-specific state hydration.
func (p *SearchPipe) Search(ctx context.Context, query string, options SearchOptions) *shared.PipeRes[SearchResponse] {
	return p.search(ctx, query, options, nil)
}

// SearchForViewer runs search with viewer-specific like/follow state hydration.
func (p *SearchPipe) SearchForViewer(ctx context.Context, query string, options SearchOptions, viewerPersonaID uuid.UUID) *shared.PipeRes[SearchResponse] {
	return p.search(ctx, query, options, &viewerPersonaID)
}

// FindViewerPersona verifies that a public profile belongs to the current user.
func (p *SearchPipe) FindViewerPersona(ctx context.Context, userID uuid.UUID, personaID uuid.UUID) (*models.Persona, shared.PipeMessage) {
	if p.personaRepo == nil {
		return nil, searchmessages.Internal_Error
	}

	persona, err := p.personaRepo.FindPersonaByID(ctx, personaID)
	if err != nil {
		if err == persona_repo.ErrPersonaNotFound {
			return nil, searchmessages.Persona_Not_Found
		}
		return nil, searchmessages.Internal_Error
	}
	if persona.UserID != userID {
		return nil, searchmessages.Forbidden
	}
	return persona, ""
}

// search executes repository search and hydrates response-only state.
func (p *SearchPipe) search(ctx context.Context, query string, options SearchOptions, viewerPersonaID *uuid.UUID) *shared.PipeRes[SearchResponse] {
	query = strings.TrimSpace(query)
	if len(query) < 2 || len(query) > 80 {
		return shared.PipeError[SearchResponse](searchmessages.Invalid_Query)
	}
	options = searchrepo.NormalizeOptions(options)
	viewerKey := "public"
	if viewerPersonaID != nil {
		viewerKey = viewerPersonaID.String()
	}
	var cached SearchResponse
	hit, err := p.loadSearchCache(ctx, query, options, viewerKey, &cached)
	if err != nil {
		return pipeInternalError[SearchResponse](err, "search.cache_load")
	}
	if hit {
		return shared.PipeSuccess(searchmessages.Search_Listed, &cached)
	}

	results, err := p.repo.Search(ctx, query, options)
	if err != nil {
		return pipeInternalError[SearchResponse](err, "search.query")
	}

	response := SearchResponse{
		Query:      query,
		Limit:      options.Limit,
		Offset:     options.Offset,
		Scope:      options.Scope,
		HasMore:    results.HasMore,
		NextOffset: nextOffset(options, results.HasMore),
		Personas:   personaResponses(results.Personas),
		Posts:      postResponses(results.Posts),
		Events:     eventResponses(results.Events),
		Hashtags:   hashtagResponses(results.Hashtags),
		Sets:       setResponses(results.Sets),
	}
	if viewerPersonaID != nil && p.likeRepo != nil {
		if err := p.hydrateLikedState(ctx, *viewerPersonaID, response.Posts); err != nil {
			return pipeInternalError[SearchResponse](err, "search.like_status")
		}
	}
	if viewerPersonaID != nil && p.followRepo != nil {
		if err := p.hydrateFollowingState(ctx, *viewerPersonaID, response.Personas); err != nil {
			return pipeInternalError[SearchResponse](err, "search.follow_status")
		}
	}
	if viewerPersonaID != nil {
		if err := p.hydrateSetLikedState(ctx, *viewerPersonaID, response.Sets); err != nil {
			return pipeInternalError[SearchResponse](err, "search.set_like_status")
		}
	}
	if err := p.hydrateHashtags(ctx, response.Posts); err != nil {
		return pipeInternalError[SearchResponse](err, "search.hashtags")
	}
	if err := p.hydrateMedia(ctx, response.Posts); err != nil {
		return pipeInternalError[SearchResponse](err, "search.media")
	}
	if err := p.hydrateAnonymousAuthors(ctx, results.Posts, response.Posts); err != nil {
		return pipeInternalError[SearchResponse](err, "search.anonymous_authors")
	}
	if err := p.storeSearchCache(ctx, query, options, viewerKey, &response); err != nil {
		return pipeInternalError[SearchResponse](err, "search.cache_store")
	}
	return shared.PipeSuccess(searchmessages.Search_Listed, &response)
}

// nextOffset calculates the next offset when more results are available.
func nextOffset(options SearchOptions, hasMore bool) *int {
	if !hasMore {
		return nil
	}
	next := options.Offset + options.Limit
	return &next
}

// hydrateLikedState attaches viewer-specific post like state.
func (p *SearchPipe) hydrateLikedState(ctx context.Context, viewerPersonaID uuid.UUID, posts []SearchPostResponse) error {
	if len(posts) == 0 {
		return nil
	}

	postIDs := make([]uuid.UUID, 0, len(posts))
	for _, post := range posts {
		postID, err := uuid.Parse(post.ID)
		if err != nil {
			return err
		}
		postIDs = append(postIDs, postID)
	}

	liked, err := p.likeRepo.FindLikedPostIDs(ctx, viewerPersonaID, postIDs)
	if err != nil {
		return err
	}
	for i := range posts {
		posts[i].IsLiked = liked[postIDs[i]]
	}
	return nil
}

// hydrateHashtags attaches hashtags to searched posts.
func (p *SearchPipe) hydrateHashtags(ctx context.Context, posts []SearchPostResponse) error {
	if p.hashtagRepo == nil || len(posts) == 0 {
		return nil
	}

	postIDs := make([]uuid.UUID, 0, len(posts))
	for _, post := range posts {
		postID, err := uuid.Parse(post.ID)
		if err != nil {
			return err
		}
		postIDs = append(postIDs, postID)
	}

	tags, err := p.hashtagRepo.FindTagsByPostIDs(ctx, postIDs)
	if err != nil {
		return err
	}
	for i := range posts {
		posts[i].Hashtags = tags[postIDs[i]]
	}
	return nil
}

// hydrateMedia attaches media assets to searched posts.
func (p *SearchPipe) hydrateMedia(ctx context.Context, posts []SearchPostResponse) error {
	if p.postRepo == nil || len(posts) == 0 {
		return nil
	}

	postIDs := make([]uuid.UUID, 0, len(posts))
	for _, post := range posts {
		postID, err := uuid.Parse(post.ID)
		if err != nil {
			return err
		}
		postIDs = append(postIDs, postID)
	}

	mediaByPost, err := p.postRepo.FindMediaAssetsByPostIDs(ctx, postIDs)
	if err != nil {
		return err
	}
	for i := range posts {
		posts[i].Media = mediaByPost[postIDs[i]]
		if posts[i].Media == nil {
			posts[i].Media = []*models.MediaAsset{}
		}
	}
	return nil
}

// hydrateAnonymousAuthors attaches thread-scoped anonymous display identities to searched posts.
func (p *SearchPipe) hydrateAnonymousAuthors(ctx context.Context, results []*searchrepo.PostResult, posts []SearchPostResponse) error {
	if p.postRepo == nil || len(results) == 0 {
		return nil
	}
	for i, result := range results {
		if i >= len(posts) || result.Post.PostingMode != models.AnonymousPostingMode || result.Post.PersonaID == nil {
			continue
		}
		identity, err := p.postRepo.EnsureAnonymousThreadIdentity(
			ctx,
			result.Post.ID,
			result.Post.AuthorUserID,
			*result.Post.PersonaID,
			anonymousidentity.GenerateHandle(),
			anonymousidentity.GenerateAvatarKey(),
		)
		if err != nil {
			return err
		}
		if identity != nil && identity.AnonymousHandle != "" {
			posts[i].Author.Anonymous = &SearchPostAnonymous{
				Handle:    identity.AnonymousHandle,
				AvatarKey: identity.AnonymousAvatarKey,
			}
		}
	}
	return nil
}

// hydrateFollowingState attaches viewer-specific follow state to searched profiles.
func (p *SearchPipe) hydrateFollowingState(ctx context.Context, viewerPersonaID uuid.UUID, personas []SearchPersonaResponse) error {
	if len(personas) == 0 {
		return nil
	}

	personaIDs := make([]uuid.UUID, 0, len(personas))
	for _, persona := range personas {
		personaID, err := uuid.Parse(persona.ID)
		if err != nil {
			return err
		}
		personaIDs = append(personaIDs, personaID)
	}

	following, err := p.followRepo.FindFollowingIDs(ctx, viewerPersonaID, personaIDs)
	if err != nil {
		return err
	}
	for i := range personas {
		personas[i].IsFollowing = following[personaIDs[i]]
	}
	return nil
}
