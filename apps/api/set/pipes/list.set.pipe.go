package pipes

import (
	"context"
	"strings"

	"github.com/emmanuella-codes/nox/models"
	"github.com/emmanuella-codes/nox/set/messages"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

func (p *SetPipe) ListSetsPipe(ctx context.Context, limit int, offset int, genreTag string, sort string, viewerPersonaID *uuid.UUID) *shared.PipeRes[SetListResponse] {
	limit = normalizeLimit(limit)
	offset = normalizeOffset(offset)
	genreTag = strings.ToLower(strings.TrimSpace(strings.TrimPrefix(genreTag, "#")))
	if genreTag != "" && !genreTagPattern.MatchString(genreTag) {
		return shared.PipeError[SetListResponse](messages.Invalid_Set)
	}
	sort = strings.TrimSpace(sort)
	sets, err := p.setRepo.FindSetsWithFilters(ctx, genreTag, sort, limit+1, offset)
	if err != nil {
		return pipeInternalError[SetListResponse](err, "set.list")
	}
	if err := p.hydrateSets(ctx, trimForHydration(limit, sets)); err != nil {
		return pipeInternalError[SetListResponse](err, "set.hydrate_list")
	}
	response, err := p.listResponse(ctx, limit, offset, sets, viewerPersonaID)
	if err != nil {
		return pipeInternalError[SetListResponse](err, "set.viewer_list")
	}
	return shared.PipeSuccess(messages.Sets_Listed, response)
}

func (p *SetPipe) ListPersonaSetsPipe(ctx context.Context, personaID uuid.UUID, limit int, offset int, viewerPersonaID *uuid.UUID) *shared.PipeRes[SetListResponse] {
	limit = normalizeLimit(limit)
	offset = normalizeOffset(offset)
	sets, err := p.setRepo.FindSetsByPersonaID(ctx, personaID, limit+1, offset)
	if err != nil {
		return pipeInternalError[SetListResponse](err, "set.list_persona")
	}
	if err := p.hydrateSets(ctx, trimForHydration(limit, sets)); err != nil {
		return pipeInternalError[SetListResponse](err, "set.hydrate_persona_list")
	}
	response, err := p.listResponse(ctx, limit, offset, sets, viewerPersonaID)
	if err != nil {
		return pipeInternalError[SetListResponse](err, "set.viewer_persona_list")
	}
	return shared.PipeSuccess(messages.Sets_Listed, response)
}

func trimForHydration(limit int, sets []*models.Set) []*models.Set {
	if len(sets) > limit {
		return sets[:limit]
	}
	return sets
}

func (p *SetPipe) listResponse(ctx context.Context, limit int, offset int, sets []*models.Set, viewerPersonaID *uuid.UUID) (*SetListResponse, error) {
	hasMore := len(sets) > limit
	if hasMore {
		sets = sets[:limit]
	}
	liked := map[uuid.UUID]bool{}
	if viewerPersonaID != nil {
		setIDs := make([]uuid.UUID, 0, len(sets))
		for _, set := range sets {
			setIDs = append(setIDs, set.ID)
		}
		var err error
		liked, err = p.setRepo.FindLikedSetIDs(ctx, *viewerPersonaID, setIDs)
		if err != nil {
			return nil, err
		}
	}
	responses := make([]SetResponse, 0, len(sets))
	for _, set := range sets {
		responses = append(responses, setResponse(set, liked[set.ID]))
	}
	return &SetListResponse{
		Limit:      limit,
		Offset:     offset,
		HasMore:    hasMore,
		NextOffset: nextOffset(limit, offset, hasMore),
		Sets:       responses,
	}, nil
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
