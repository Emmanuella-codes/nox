package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/hashtag/messages"
	hashtag_repo "github.com/emmanuella-codes/nox/repositories/hashtag"
	"github.com/emmanuella-codes/nox/shared"
)

type HashtagDetailResponse struct {
	Tag       string `json:"tag"`
	PostCount int    `json:"post_count"`
}

func (p *HashtagPipe) GetHashtagPipe(ctx context.Context, tag string) *shared.PipeRes[HashtagDetailResponse] {
	normalized := hashtag_repo.NormalizeTag(tag)
	if normalized == "" {
		return shared.PipeError[HashtagDetailResponse](messages.Invalid_Tag)
	}

	hashtag, err := p.repo.FindByTag(ctx, normalized)
	if err != nil {
		return pipeInternalError[HashtagDetailResponse](err, "hashtag.get")
	}
	response := HashtagDetailResponse{Tag: normalized}
	if hashtag != nil {
		response.PostCount = hashtag.PostCount
	}
	return shared.PipeSuccess(messages.Hashtag_Fetched, &response)
}
