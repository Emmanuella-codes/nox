package pipes

import (
	"context"
	"fmt"
	"time"

	"github.com/emmanuella-codes/nox/shared/cache"
)

const searchCacheTTL = 30 * time.Second

// loads one cached search response when Redis caching is configured.
func (p *SearchPipe) loadSearchCache(ctx context.Context, query string, options SearchOptions, viewerKey string, out *SearchResponse) (bool, error) {
	key := fmt.Sprintf("search:%s:%s:%d:%d:%s", viewerKey, options.Scope, options.Limit, options.Offset, query)
	return cache.GetJSON(ctx, p.cacheClient, key, out)
}

// stores one search response in Redis when caching is configured.
func (p *SearchPipe) storeSearchCache(ctx context.Context, query string, options SearchOptions, viewerKey string, value *SearchResponse) error {
	key := fmt.Sprintf("search:%s:%s:%d:%d:%s", viewerKey, options.Scope, options.Limit, options.Offset, query)
	return cache.SetJSON(ctx, p.cacheClient, key, value, searchCacheTTL)
}
