package pipes

import (
	"context"
	"strings"

	searchmessages "github.com/emmanuella-codes/nox/search/messages"
	"github.com/emmanuella-codes/nox/shared"
)

func (p *SearchPipe) Search(ctx context.Context, query string, limit int) *shared.PipeRes[SearchResponse] {
	query = strings.TrimSpace(query)
	if len(query) < 2 || len(query) > 80 {
		return pipeError[SearchResponse](searchmessages.Invalid_Query)
	}

	results, err := p.repo.Search(ctx, query, limit)
	if err != nil {
		return pipeInternalError[SearchResponse](err, "search.query")
	}

	response := SearchResponse{
		Query:    query,
		Personas: personaResponses(results.Personas),
		Posts:    postResponses(results.Posts),
		Events:   eventResponses(results.Events),
	}
	return pipeSuccess(searchmessages.Search_Listed, &response)
}
