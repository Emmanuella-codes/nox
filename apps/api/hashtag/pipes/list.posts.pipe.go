package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/hashtag/messages"
	post_pipes "github.com/emmanuella-codes/nox/post/pipes"
	hashtag_repo "github.com/emmanuella-codes/nox/repositories/hashtag"
	"github.com/emmanuella-codes/nox/shared"
)

type HashtagPostsResponse struct {
	Tag   string                    `json:"tag"`
	Posts []post_pipes.PostResponse `json:"posts"`
}

func (p *HashtagPipe) PostsByTagPipe(ctx context.Context, tag string, limit int) *shared.PipeRes[HashtagPostsResponse] {
	normalized := hashtag_repo.NormalizeTag(tag)
	if normalized == "" {
		return shared.PipeError[HashtagPostsResponse](messages.Invalid_Tag)
	}

	posts, err := p.repo.FindPostsByTag(ctx, normalized, limit)
	if err != nil {
		return pipeInternalError[HashtagPostsResponse](err, "hashtag.posts")
	}
	responses, err := p.postResponses(ctx, posts)
	if err != nil {
		return pipeInternalError[HashtagPostsResponse](err, "hashtag.post_responses")
	}
	response := HashtagPostsResponse{Tag: normalized, Posts: responses}
	return shared.PipeSuccess(messages.Hashtag_Posts, &response)
}
