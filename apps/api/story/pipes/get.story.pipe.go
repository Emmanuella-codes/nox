package pipes

import (
	"context"

	story_repo "github.com/emmanuella-codes/nox/repositories/story"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/emmanuella-codes/nox/story/messages"
	"github.com/google/uuid"
)

func (p *StoryPipe) GetStoryPipe(ctx context.Context, storyID uuid.UUID, viewerPersonaID *uuid.UUID) *shared.PipeRes[StoryResponse] {
	story, err := p.storyRepo.FindStoryByID(ctx, storyID)
	if err != nil {
		if err == story_repo.ErrStoryNotFound {
			return shared.PipeError[StoryResponse](messages.Story_Not_Found)
		}
		return pipeInternalError[StoryResponse](err, "story.get")
	}
	response, err := p.storyResponse(ctx, story, viewerPersonaID, true)
	if err != nil {
		return pipeInternalError[StoryResponse](err, "story.response")
	}
	return shared.PipeSuccess(messages.Story_Fetched, response)
}

func (p *StoryPipe) ListStoryItemsPipe(ctx context.Context, storyID uuid.UUID) *shared.PipeRes[[]StoryItemResponse] {
	if _, err := p.storyRepo.FindStoryByID(ctx, storyID); err != nil {
		if err == story_repo.ErrStoryNotFound {
			return shared.PipeError[[]StoryItemResponse](messages.Story_Not_Found)
		}
		return pipeInternalError[[]StoryItemResponse](err, "story.find_for_items")
	}
	items, err := p.storyItemResponses(ctx, storyID)
	if err != nil {
		return pipeInternalError[[]StoryItemResponse](err, "story.items")
	}
	return shared.PipeSuccess(messages.Story_Items_Listed, &items)
}
