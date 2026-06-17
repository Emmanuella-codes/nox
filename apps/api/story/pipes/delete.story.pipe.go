package pipes

import (
	"context"

	story_repo "github.com/emmanuella-codes/nox/repositories/story"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/emmanuella-codes/nox/story/messages"
	"github.com/google/uuid"
)

func (p *StoryPipe) DeleteStoryPipe(ctx context.Context, userID uuid.UUID, storyID uuid.UUID) *shared.PipeRes[any] {
	story, err := p.storyRepo.FindStoryByID(ctx, storyID)
	if err != nil {
		if err == story_repo.ErrStoryNotFound {
			return shared.PipeError[any](messages.Story_Not_Found)
		}
		return pipeInternalError[any](err, "story.find_for_delete")
	}
	if story.OwnerUserID != userID {
		return shared.PipeError[any](messages.Forbidden)
	}
	if err := p.storyRepo.DeleteStory(ctx, storyID); err != nil {
		if err == story_repo.ErrStoryNotFound {
			return shared.PipeError[any](messages.Story_Not_Found)
		}
		return pipeInternalError[any](err, "story.delete")
	}
	return shared.PipeSuccess[any](messages.Story_Deleted, nil)
}
