package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/hashtag/messages"
	post_pipes "github.com/emmanuella-codes/nox/post/pipes"
	hashtag_repo "github.com/emmanuella-codes/nox/repositories/hashtag"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

type HashtagPostsResponse struct {
	Tag        string                    `json:"tag"`
	Limit      int                       `json:"limit"`
	Offset     int                       `json:"offset"`
	HasMore    bool                      `json:"has_more"`
	NextOffset *int                      `json:"next_offset,omitempty"`
	Posts      []post_pipes.PostResponse `json:"posts"`
}

func (p *HashtagPipe) PostsByTagPipe(ctx context.Context, tag string, limit int, offset int) *shared.PipeRes[HashtagPostsResponse] {
	return p.postsByTag(ctx, tag, limit, offset, nil)
}

func (p *HashtagPipe) PostsByTagForViewerPipe(ctx context.Context, tag string, limit int, offset int, viewerPersonaID uuid.UUID) *shared.PipeRes[HashtagPostsResponse] {
	return p.postsByTag(ctx, tag, limit, offset, &viewerPersonaID)
}

func (p *HashtagPipe) postsByTag(ctx context.Context, tag string, limit int, offset int, viewerPersonaID *uuid.UUID) *shared.PipeRes[HashtagPostsResponse] {
	normalized := hashtag_repo.NormalizeTag(tag)
	if normalized == "" {
		return shared.PipeError[HashtagPostsResponse](messages.Invalid_Tag)
	}

	limit = normalizeLimit(limit)
	offset = normalizeOffset(offset)
	fetchLimit := limit + 1

	posts, err := p.repo.FindPostsByTag(ctx, normalized, fetchLimit, offset)
	if err != nil {
		return pipeInternalError[HashtagPostsResponse](err, "hashtag.posts")
	}
	hasMore := len(posts) > limit
	if hasMore {
		posts = posts[:limit]
	}

	responses, err := p.postResponses(ctx, posts)
	if err != nil {
		return pipeInternalError[HashtagPostsResponse](err, "hashtag.post_responses")
	}
	if viewerPersonaID != nil {
		if err := p.hydrateLikedState(ctx, *viewerPersonaID, responses); err != nil {
			return pipeInternalError[HashtagPostsResponse](err, "hashtag.like_status")
		}
	}

	response := HashtagPostsResponse{
		Tag:        normalized,
		Limit:      limit,
		Offset:     offset,
		HasMore:    hasMore,
		NextOffset: nextOffset(limit, offset, hasMore),
		Posts:      responses,
	}
	return shared.PipeSuccess(messages.Hashtag_Posts, &response)
}

func normalizeLimit(limit int) int {
	if limit <= 0 {
		return 20
	}
	if limit > 30 {
		return 30
	}
	return limit
}

func normalizeOffset(offset int) int {
	if offset < 0 {
		return 0
	}
	return offset
}

func nextOffset(limit int, offset int, hasMore bool) *int {
	if !hasMore {
		return nil
	}
	next := offset + limit
	return &next
}
