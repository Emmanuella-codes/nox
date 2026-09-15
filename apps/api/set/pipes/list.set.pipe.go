package pipes

import (
	"context"
	"strings"

	"github.com/emmanuella-codes/nox/models"
	"github.com/emmanuella-codes/nox/set/messages"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

func (p *SetPipe) ListSetsPipe(ctx context.Context, limit int, offset int, genreTag string, sort string, viewerPersonaID *uuid.UUID, cursors ...string) *shared.PipeRes[SetListResponse] {
	limit = normalizeLimit(limit)
	offset = normalizeOffset(offset)
	genreTag = strings.ToLower(strings.TrimSpace(strings.TrimPrefix(genreTag, "#")))
	if genreTag != "" && !genreTagPattern.MatchString(genreTag) {
		return shared.PipeError[SetListResponse](messages.Invalid_Set)
	}
	sort = strings.TrimSpace(sort)
	cursor := ""
	if len(cursors) > 0 {
		cursor = cursors[0]
	}
	if cursorOffset, valid := decodeSetCursor(cursor, genreTag, sort); !valid {
		return shared.PipeError[SetListResponse](messages.Invalid_Set)
	} else if cursor != "" {
		offset = cursorOffset
	}
	sets, err := p.setRepo.FindSetsWithFilters(ctx, genreTag, sort, limit+1, offset)
	if err != nil {
		return pipeInternalError[SetListResponse](err, "set.list")
	}
	if err := p.hydrateSets(ctx, sets); err != nil {
		return pipeInternalError[SetListResponse](err, "set.hydrate_list")
	}
	sets, err = p.filterViewerSets(ctx, sets, viewerPersonaID)
	if err != nil {
		return pipeInternalError[SetListResponse](err, "set.visibility")
	}
	response, err := p.listResponse(ctx, limit, offset, genreTag, sort, sets, viewerPersonaID)
	if err != nil {
		return pipeInternalError[SetListResponse](err, "set.viewer_list")
	}
	return shared.PipeSuccess(messages.Sets_Listed, response)
}

func (p *SetPipe) ListPersonaSetsPipe(ctx context.Context, personaID uuid.UUID, limit int, offset int, viewerPersonaID *uuid.UUID, cursors ...string) *shared.PipeRes[SetListResponse] {
	limit = normalizeLimit(limit)
	offset = normalizeOffset(offset)
	cursor := ""
	if len(cursors) > 0 {
		cursor = cursors[0]
	}
	if cursorOffset, valid := decodeSetCursor(cursor, "", "persona"); !valid {
		return shared.PipeError[SetListResponse](messages.Invalid_Set)
	} else if cursor != "" {
		offset = cursorOffset
	}
	sets, err := p.setRepo.FindSetsByPersonaID(ctx, personaID, limit+1, offset)
	if err != nil {
		return pipeInternalError[SetListResponse](err, "set.list_persona")
	}
	if err := p.hydrateSets(ctx, sets); err != nil {
		return pipeInternalError[SetListResponse](err, "set.hydrate_persona_list")
	}
	sets, err = p.filterViewerSets(ctx, sets, viewerPersonaID)
	if err != nil {
		return pipeInternalError[SetListResponse](err, "set.visibility")
	}
	response, err := p.listResponse(ctx, limit, offset, "", "persona", sets, viewerPersonaID)
	if err != nil {
		return pipeInternalError[SetListResponse](err, "set.viewer_persona_list")
	}
	return shared.PipeSuccess(messages.Sets_Listed, response)
}

func (p *SetPipe) listResponse(ctx context.Context, limit int, offset int, genre string, sort string, sets []*models.Set, viewerPersonaID *uuid.UUID) (*SetListResponse, error) {
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
		NextCursor: nextSetCursor(limit, offset, genre, sort, hasMore),
		Sets:       responses,
	}, nil
}

func nextSetCursor(limit int, offset int, genre string, sort string, hasMore bool) string {
	if !hasMore {
		return ""
	}
	return encodeSetCursor(offset+limit, genre, sort)
}

func (p *SetPipe) filterViewerSets(ctx context.Context, sets []*models.Set, viewerPersonaID *uuid.UUID) ([]*models.Set, error) {
	if viewerPersonaID == nil || p.preferenceRepo == nil {
		return sets, nil
	}
	viewer, err := p.personaRepo.FindPersonaByID(ctx, *viewerPersonaID)
	if err != nil {
		return nil, err
	}
	excluded, err := p.preferenceRepo.FindExcludedUserIDs(ctx, viewer.UserID)
	if err != nil {
		return nil, err
	}
	muted, err := p.preferenceRepo.FindMutedUserIDs(ctx, viewer.UserID)
	if err != nil {
		return nil, err
	}
	suppressed, err := p.preferenceRepo.FindSuppressedTargetIDs(ctx, viewer.UserID, models.SetSuppressionTargetType)
	if err != nil {
		return nil, err
	}
	filtered := make([]*models.Set, 0, len(sets))
	for _, set := range sets {
		if set == nil || set.Persona == nil || suppressed[set.ID] {
			continue
		}
		if set.Persona != nil && (excluded[set.Persona.UserID] || muted[set.Persona.UserID]) {
			continue
		}
		filtered = append(filtered, set)
	}
	return filtered, nil
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
