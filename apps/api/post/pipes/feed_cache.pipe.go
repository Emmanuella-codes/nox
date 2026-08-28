package pipes

import (
	"context"
	"fmt"
	"time"

	"github.com/emmanuella-codes/nox/shared/cache"
)

const feedCacheTTL = 30 * time.Second

// loadFeedCache loads one cached feed response when Redis caching is configured.
func (p *PostPipe) loadFeedCache(ctx context.Context, kind string, personaID string, limit int, cursor string, out *PostListResponse) (bool, error) {
	key := fmt.Sprintf("feed:%s:%s:%d:%s", kind, personaID, limit, cursor)
	return cache.GetJSON(ctx, p.cacheClient, key, out)
}

// storeFeedCache stores one feed response in Redis when caching is configured.
func (p *PostPipe) storeFeedCache(ctx context.Context, kind string, personaID string, limit int, cursor string, value *PostListResponse) error {
	key := fmt.Sprintf("feed:%s:%s:%d:%s", kind, personaID, limit, cursor)
	return cache.SetJSON(ctx, p.cacheClient, key, value, feedCacheTTL)
}
