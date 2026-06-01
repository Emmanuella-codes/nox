package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/models"
	"github.com/emmanuella-codes/nox/set/messages"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

type SetListResponse struct {
	Limit      int           `json:"limit"`
	Offset     int           `json:"offset"`
	HasMore    bool          `json:"has_more"`
	NextOffset *int          `json:"next_offset,omitempty"`
	Sets       []*models.Set `json:"sets"`
}

func (p *SetPipe) ListSetsPipe(ctx context.Context, limit int, offset int) *shared.PipeRes[SetListResponse] {
	limit = normalizeLimit(limit)
	offset = normalizeOffset(offset)
	sets, err := p.setRepo.FindSets(ctx, limit+1, offset)
	if err != nil {
		return pipeInternalError[SetListResponse](err, "set.list")
	}
	if err := p.hydrateSets(ctx, trimForHydration(limit, sets)); err != nil {
		return pipeInternalError[SetListResponse](err, "set.hydrate_list")
	}
	return shared.PipeSuccess(messages.Sets_Listed, listResponse(limit, offset, sets))
}

func (p *SetPipe) ListPersonaSetsPipe(ctx context.Context, personaID uuid.UUID, limit int, offset int) *shared.PipeRes[SetListResponse] {
	limit = normalizeLimit(limit)
	offset = normalizeOffset(offset)
	sets, err := p.setRepo.FindSetsByPersonaID(ctx, personaID, limit+1, offset)
	if err != nil {
		return pipeInternalError[SetListResponse](err, "set.list_persona")
	}
	if err := p.hydrateSets(ctx, trimForHydration(limit, sets)); err != nil {
		return pipeInternalError[SetListResponse](err, "set.hydrate_persona_list")
	}
	return shared.PipeSuccess(messages.Sets_Listed, listResponse(limit, offset, sets))
}

func trimForHydration(limit int, sets []*models.Set) []*models.Set {
	if len(sets) > limit {
		return sets[:limit]
	}
	return sets
}

func listResponse(limit int, offset int, sets []*models.Set) *SetListResponse {
	hasMore := len(sets) > limit
	if hasMore {
		sets = sets[:limit]
	}
	return &SetListResponse{
		Limit:      limit,
		Offset:     offset,
		HasMore:    hasMore,
		NextOffset: nextOffset(limit, offset, hasMore),
		Sets:       sets,
	}
}

func normalizeLimit(limit int) int {
	if limit <= 0 {
		return 20
	}
	if limit > 50 {
		return 50
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
