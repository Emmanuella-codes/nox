package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/models"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/emmanuella-codes/nox/story/messages"
	"github.com/google/uuid"
)

func (p *StoryPipe) ListEventStoriesPipe(ctx context.Context, eventID uuid.UUID, limit int, offset int, viewerUserID *uuid.UUID, viewerPersonaID *uuid.UUID) *shared.PipeRes[StoryListResponse] {
	viewer, message := p.viewerPersona(ctx, viewerUserID, viewerPersonaID)
	if message != "" {
		return shared.PipeError[StoryListResponse](message)
	}
	if viewer != nil {
		viewerPersonaID = &viewer.ID
	}
	limit = normalizeLimit(limit)
	offset = normalizeOffset(offset)
	stories, err := p.storyRepo.FindStoriesByEventID(ctx, eventID, limit+1, offset)
	if err != nil {
		return pipeInternalError[StoryListResponse](err, "story.list_event")
	}
	return p.listResponse(ctx, stories, limit, offset, viewerPersonaID)
}

func (p *StoryPipe) ListPersonaStoriesPipe(ctx context.Context, personaID uuid.UUID, limit int, offset int, viewerUserID *uuid.UUID, viewerPersonaID *uuid.UUID) *shared.PipeRes[StoryListResponse] {
	viewer, message := p.viewerPersona(ctx, viewerUserID, viewerPersonaID)
	if message != "" {
		return shared.PipeError[StoryListResponse](message)
	}
	if viewer != nil {
		viewerPersonaID = &viewer.ID
	}
	limit = normalizeLimit(limit)
	offset = normalizeOffset(offset)
	stories, err := p.storyRepo.FindStoriesByOwnerPersonaID(ctx, personaID, limit+1, offset)
	if err != nil {
		return pipeInternalError[StoryListResponse](err, "story.list_persona")
	}
	return p.listResponse(ctx, stories, limit, offset, viewerPersonaID)
}

func (p *StoryPipe) listResponse(ctx context.Context, stories []*models.Story, limit int, offset int, viewerPersonaID *uuid.UUID) *shared.PipeRes[StoryListResponse] {
	hasMore := len(stories) > limit
	if hasMore {
		stories = stories[:limit]
	}
	responses := make([]StoryResponse, 0, len(stories))
	for _, story := range stories {
		allowed, err := p.canView(ctx, story, viewerPersonaID)
		if err != nil {
			return pipeInternalError[StoryListResponse](err, "story.list_view")
		}
		if !allowed {
			continue
		}
		response, err := p.storyResponse(ctx, story, viewerPersonaID, true, false)
		if err != nil {
			return pipeInternalError[StoryListResponse](err, "story.list_response")
		}
		responses = append(responses, *response)
	}
	return shared.PipeSuccess(messages.Stories_Listed, &StoryListResponse{
		Limit:      limit,
		Offset:     offset,
		HasMore:    hasMore,
		NextOffset: nextOffset(limit, offset, hasMore),
		Stories:    responses,
	})
}
