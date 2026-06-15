package pipes

import (
	"context"
	"strings"

	event_repo "github.com/emmanuella-codes/nox/repositories/event"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/emmanuella-codes/nox/story/dtos"
	"github.com/emmanuella-codes/nox/story/messages"
	"github.com/google/uuid"
)

func (p *StoryPipe) CreateStoryPipe(ctx context.Context, userID uuid.UUID, dto dtos.CreateStoryDTO) *shared.PipeRes[StoryResponse] {
	dto.Title = strings.TrimSpace(dto.Title)
	if dto.Title == "" || !validContributionMode(dto.ContributionMode) {
		return shared.PipeError[StoryResponse](messages.Invalid_Story)
	}

	if _, err := p.eventRepo.FindEventByID(ctx, dto.EventID); err != nil {
		if err == event_repo.ErrEventNotFound {
			return shared.PipeError[StoryResponse](messages.Event_Not_Found)
		}
		return pipeInternalError[StoryResponse](err, "story.find_event")
	}

	if _, message := p.ownedPersona(ctx, userID, dto.OwnerPersonaID); message != "" {
		return shared.PipeError[StoryResponse](message)
	}
	dto.ExpiresAt = defaultStoryExpiry(dto.ExpiresAt)

	story, err := p.storyRepo.CreateStory(ctx, userID, dto)
	if err != nil {
		return pipeInternalError[StoryResponse](err, "story.create")
	}
	response, err := p.storyResponse(ctx, story, &dto.OwnerPersonaID, false)
	if err != nil {
		return pipeInternalError[StoryResponse](err, "story.response")
	}
	return shared.PipeSuccess(messages.Story_Created, response)
}
