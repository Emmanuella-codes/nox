package pipes

import (
	"context"

	story_repo "github.com/emmanuella-codes/nox/repositories/story"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/emmanuella-codes/nox/story/messages"
	"github.com/google/uuid"
)

func (p *StoryPipe) GetStoryPipe(ctx context.Context, storyID uuid.UUID, viewerUserID *uuid.UUID, viewerPersonaID *uuid.UUID) *shared.PipeRes[StoryResponse] {
	viewer, message := p.viewerPersona(ctx, viewerUserID, viewerPersonaID)
	if message != "" {
		return shared.PipeError[StoryResponse](message)
	}
	story, err := p.storyRepo.FindStoryByID(ctx, storyID)
	if err != nil {
		if err == story_repo.ErrStoryNotFound {
			return shared.PipeError[StoryResponse](messages.Story_Not_Found)
		}
		return pipeInternalError[StoryResponse](err, "story.get")
	}
	allowed, err := p.canView(ctx, story, viewerPersonaID)
	if err != nil {
		return pipeInternalError[StoryResponse](err, "story.view")
	}
	if !allowed {
		return shared.PipeError[StoryResponse](messages.Story_Not_Found)
	}
	if viewer != nil {
		viewerPersonaID = &viewer.ID
	}
	response, err := p.storyResponse(ctx, story, viewerPersonaID, true, false)
	if err != nil {
		return pipeInternalError[StoryResponse](err, "story.response")
	}
	return shared.PipeSuccess(messages.Story_Fetched, response)
}

func (p *StoryPipe) ListStoryItemsPipe(ctx context.Context, storyID uuid.UUID, limit int, offset int, viewerUserID *uuid.UUID, viewerPersonaID *uuid.UUID) *shared.PipeRes[[]StoryItemResponse] {
	viewer, message := p.viewerPersona(ctx, viewerUserID, viewerPersonaID)
	if message != "" {
		return shared.PipeError[[]StoryItemResponse](message)
	}
	story, err := p.storyRepo.FindStoryByID(ctx, storyID)
	if err != nil {
		if err == story_repo.ErrStoryNotFound {
			return shared.PipeError[[]StoryItemResponse](messages.Story_Not_Found)
		}
		return pipeInternalError[[]StoryItemResponse](err, "story.find_for_items")
	}
	if viewer != nil {
		viewerPersonaID = &viewer.ID
	}
	allowed, err := p.canView(ctx, story, viewerPersonaID)
	if err != nil {
		return pipeInternalError[[]StoryItemResponse](err, "story.items_view")
	}
	if !allowed {
		return shared.PipeError[[]StoryItemResponse](messages.Story_Not_Found)
	}
	items, err := p.storyItemResponses(ctx, storyID, viewerPersonaID, false)
	if err != nil {
		return pipeInternalError[[]StoryItemResponse](err, "story.items")
	}
	limit = normalizeLimit(limit)
	offset = normalizeOffset(offset)
	if offset >= len(items) {
		items = []StoryItemResponse{}
	} else {
		end := offset + limit
		if end > len(items) {
			end = len(items)
		}
		items = items[offset:end]
	}
	return shared.PipeSuccess(messages.Story_Items_Listed, &items)
}
