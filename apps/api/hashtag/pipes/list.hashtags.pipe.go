package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/hashtag/messages"
	"github.com/emmanuella-codes/nox/models"
	"github.com/emmanuella-codes/nox/shared"
)

func (p *HashtagPipe) TrendingPipe(ctx context.Context, limit int) *shared.PipeRes[[]*models.Hashtag] {
	hashtags, err := p.repo.FindTrending(ctx, limit)
	if err != nil {
		return pipeInternalError[[]*models.Hashtag](err, "hashtag.trending")
	}
	return shared.PipeSuccess(messages.Hashtags_Listed, &hashtags)
}
