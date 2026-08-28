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

// AddEventHighlightStoryPipe adds one story to an event's highlight list.
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
	if err := p.notifyStoryHighlightChange(ctx, story, persona, models.StoryHighlightAddedNotificationType); err != nil {
		return pipeInternalError[EventHighlightStoryResponse](err, "story.add_highlight_notify")
	}
	response, err := p.eventHighlightResponse(ctx, highlight)
	if err != nil {
		return pipeInternalError[EventHighlightStoryResponse](err, "story.highlight_response")
	}
	return shared.PipeSuccess(messages.Event_Highlight_Story_Added, response)
}

// ListEventHighlightStoriesPipe lists highlight stories visible to the current viewer.
func (p *StoryPipe) ListEventHighlightStoriesPipe(ctx context.Context, eventID uuid.UUID, limit int, offset int, viewerUserID *uuid.UUID, viewerPersonaID *uuid.UUID) *shared.PipeRes[[]EventHighlightStoryResponse] {
	viewer, message := p.viewerPersona(ctx, viewerUserID, viewerPersonaID)
	if message != "" {
		return shared.PipeError[[]EventHighlightStoryResponse](message)
	}
	if viewer != nil {
		viewerPersonaID = &viewer.ID
	}
	highlights, err := p.storyRepo.FindEventHighlightStories(ctx, eventID)
	if err != nil {
		return pipeInternalError[[]EventHighlightStoryResponse](err, "story.list_highlights")
	}
	responses := make([]EventHighlightStoryResponse, 0, len(highlights))
	for _, highlight := range highlights {
		story, err := p.storyRepo.FindStoryByID(ctx, highlight.StoryID)
		if err != nil {
			return pipeInternalError[[]EventHighlightStoryResponse](err, "story.highlight_story")
		}
		allowed, err := p.canView(ctx, story, viewerPersonaID)
		if err != nil {
			return pipeInternalError[[]EventHighlightStoryResponse](err, "story.highlight_view")
		}
		if !allowed {
			continue
		}
		response, err := p.eventHighlightResponse(ctx, highlight)
		if err != nil {
			return pipeInternalError[[]EventHighlightStoryResponse](err, "story.highlight_response")
		}
		responses = append(responses, *response)
	}
	limit = normalizeLimit(limit)
	offset = normalizeOffset(offset)
	if offset >= len(responses) {
		responses = []EventHighlightStoryResponse{}
	} else {
		end := offset + limit
		if end > len(responses) {
			end = len(responses)
		}
		responses = responses[offset:end]
	}
	return shared.PipeSuccess(messages.Event_Highlight_Stories_Listed, &responses)
}

// RemoveEventHighlightStoryPipe removes one story from an event's highlight list.
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
	story, err := p.storyRepo.FindStoryByID(ctx, storyID)
	if err != nil {
		if err == story_repo.ErrStoryNotFound {
			return shared.PipeError[any](messages.Story_Not_Found)
		}
		return pipeInternalError[any](err, "story.remove_highlight_story")
	}
	if err := p.storyRepo.RemoveEventHighlightStory(ctx, eventID, storyID); err != nil {
		if err == story_repo.ErrEventHighlightNotFound {
			return shared.PipeError[any](messages.Story_Not_Found)
		}
		return pipeInternalError[any](err, "story.remove_highlight")
	}
	if err := p.notifyStoryHighlightChange(ctx, story, persona, models.StoryHighlightRemovedNotificationType); err != nil {
		return pipeInternalError[any](err, "story.remove_highlight_notify")
	}
	return shared.PipeSuccess[any](messages.Event_Highlight_Story_Removed, nil)
}

// ReorderEventHighlightStoryPipe reorders one highlighted story within an event.
func (p *StoryPipe) ReorderEventHighlightStoryPipe(ctx context.Context, userID uuid.UUID, eventID uuid.UUID, storyID uuid.UUID, dto dtos.ReorderEventHighlightStoryDTO) *shared.PipeRes[EventHighlightStoryResponse] {
	if dto.Position < 1 {
		return shared.PipeError[EventHighlightStoryResponse](messages.Invalid_Story)
	}
	event, err := p.eventRepo.FindEventByID(ctx, eventID)
	if err != nil {
		if err == event_repo.ErrEventNotFound {
			return shared.PipeError[EventHighlightStoryResponse](messages.Event_Not_Found)
		}
		return pipeInternalError[EventHighlightStoryResponse](err, "story.reorder_highlight_event")
	}
	persona, err := p.personaRepo.FindPersonaByID(ctx, dto.PersonaID)
	if err != nil {
		if err == persona_repo.ErrPersonaNotFound {
			return shared.PipeError[EventHighlightStoryResponse](messages.Persona_Not_Found)
		}
		return pipeInternalError[EventHighlightStoryResponse](err, "story.reorder_highlight_persona")
	}
	if persona.UserID != userID || event.OrganizerID != persona.ID {
		return shared.PipeError[EventHighlightStoryResponse](messages.Forbidden)
	}
	highlight, err := p.storyRepo.ReorderEventHighlightStory(ctx, eventID, storyID, dto.Position)
	if err != nil {
		if err == story_repo.ErrEventHighlightNotFound {
			return shared.PipeError[EventHighlightStoryResponse](messages.Story_Not_Found)
		}
		return pipeInternalError[EventHighlightStoryResponse](err, "story.reorder_highlight")
	}
	response, err := p.eventHighlightResponse(ctx, highlight)
	if err != nil {
		return pipeInternalError[EventHighlightStoryResponse](err, "story.reorder_highlight_response")
	}
	return shared.PipeSuccess(messages.Event_Highlight_Story_Reordered, response)
}

// eventHighlightResponse maps one highlight row into the API response shape.
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
