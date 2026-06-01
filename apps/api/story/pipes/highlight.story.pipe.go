package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/models"
	event_repo "github.com/emmanuella-codes/nox/repositories/event"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	story_repo "github.com/emmanuella-codes/nox/repositories/story"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/emmanuella-codes/nox/story/dtos"
	"github.com/emmanuella-codes/nox/story/messages"
	"github.com/google/uuid"
)

func (p *StoryPipe) AddEventHighlightStoryPipe(ctx context.Context, userID uuid.UUID, eventID uuid.UUID, dto dtos.AddEventHighlightStoryDTO) *shared.PipeRes[EventHighlightStoryResponse] {
	event, err := p.eventRepo.FindEventByID(ctx, eventID)
	if err != nil {
		if err == event_repo.ErrEventNotFound {
			return shared.PipeError[EventHighlightStoryResponse](messages.Event_Not_Found)
		}
		return pipeInternalError[EventHighlightStoryResponse](err, "story.highlight_event")
	}
	persona, err := p.personaRepo.FindPersonaByID(ctx, dto.AddedByPersonaID)
	if err != nil {
		if err == persona_repo.ErrPersonaNotFound {
			return shared.PipeError[EventHighlightStoryResponse](messages.Persona_Not_Found)
		}
		return pipeInternalError[EventHighlightStoryResponse](err, "story.highlight_persona")
	}
	if persona.UserID != userID || event.OrganizerID != persona.ID {
		return shared.PipeError[EventHighlightStoryResponse](messages.Forbidden)
	}
	story, err := p.storyRepo.FindStoryByID(ctx, dto.StoryID)
	if err != nil {
		if err == story_repo.ErrStoryNotFound {
			return shared.PipeError[EventHighlightStoryResponse](messages.Story_Not_Found)
		}
		return pipeInternalError[EventHighlightStoryResponse](err, "story.highlight_story")
	}
	if story.EventID != eventID {
		return shared.PipeError[EventHighlightStoryResponse](messages.Invalid_Story)
	}

	highlight, err := p.storyRepo.AddEventHighlightStory(ctx, eventID, dto)
	if err != nil {
		return pipeInternalError[EventHighlightStoryResponse](err, "story.add_highlight")
	}
	response, err := p.eventHighlightResponse(ctx, highlight)
	if err != nil {
		return pipeInternalError[EventHighlightStoryResponse](err, "story.highlight_response")
	}
	return shared.PipeSuccess(messages.Event_Highlight_Story_Added, response)
}

func (p *StoryPipe) ListEventHighlightStoriesPipe(ctx context.Context, eventID uuid.UUID) *shared.PipeRes[[]EventHighlightStoryResponse] {
	highlights, err := p.storyRepo.FindEventHighlightStories(ctx, eventID)
	if err != nil {
		return pipeInternalError[[]EventHighlightStoryResponse](err, "story.list_highlights")
	}
	responses := make([]EventHighlightStoryResponse, 0, len(highlights))
	for _, highlight := range highlights {
		response, err := p.eventHighlightResponse(ctx, highlight)
		if err != nil {
			return pipeInternalError[[]EventHighlightStoryResponse](err, "story.highlight_response")
		}
		responses = append(responses, *response)
	}
	return shared.PipeSuccess(messages.Event_Highlight_Stories_Listed, &responses)
}

func (p *StoryPipe) RemoveEventHighlightStoryPipe(ctx context.Context, userID uuid.UUID, eventID uuid.UUID, storyID uuid.UUID, addedByPersonaID uuid.UUID) *shared.PipeRes[any] {
	event, err := p.eventRepo.FindEventByID(ctx, eventID)
	if err != nil {
		if err == event_repo.ErrEventNotFound {
			return shared.PipeError[any](messages.Event_Not_Found)
		}
		return pipeInternalError[any](err, "story.remove_highlight_event")
	}
	persona, err := p.personaRepo.FindPersonaByID(ctx, addedByPersonaID)
	if err != nil {
		if err == persona_repo.ErrPersonaNotFound {
			return shared.PipeError[any](messages.Persona_Not_Found)
		}
		return pipeInternalError[any](err, "story.remove_highlight_persona")
	}
	if persona.UserID != userID || event.OrganizerID != persona.ID {
		return shared.PipeError[any](messages.Forbidden)
	}
	if err := p.storyRepo.RemoveEventHighlightStory(ctx, eventID, storyID); err != nil {
		if err == story_repo.ErrEventHighlightNotFound {
			return shared.PipeError[any](messages.Story_Not_Found)
		}
		return pipeInternalError[any](err, "story.remove_highlight")
	}
	return shared.PipeSuccess[any](messages.Event_Highlight_Story_Removed, nil)
}

func (p *StoryPipe) eventHighlightResponse(ctx context.Context, highlight *models.EventHighlightStory) (*EventHighlightStoryResponse, error) {
	response := &EventHighlightStoryResponse{
		ID:               uuidString(highlight.ID),
		EventID:          uuidString(highlight.EventID),
		StoryID:          uuidString(highlight.StoryID),
		AddedByPersonaID: uuidString(highlight.AddedByPersonaID),
		Position:         highlight.Position,
		CreatedAt:        highlight.CreatedAt,
	}
	story, err := p.storyRepo.FindStoryByID(ctx, highlight.StoryID)
	if err != nil {
		return nil, err
	}
	storyResponse, err := p.storyResponse(ctx, story, nil, false)
	if err != nil {
		return nil, err
	}
	response.Story = storyResponse
	return response, nil
}
