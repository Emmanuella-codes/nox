package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/models"
	story_repo "github.com/emmanuella-codes/nox/repositories/story"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/emmanuella-codes/nox/story/dtos"
	"github.com/emmanuella-codes/nox/story/messages"
	"github.com/google/uuid"
)

// AddProfileStoryHighlightPipe adds one story to one persona's profile highlights.
func (p *StoryPipe) AddProfileStoryHighlightPipe(ctx context.Context, userID uuid.UUID, storyID uuid.UUID, dto dtos.ProfileStoryHighlightDTO) *shared.PipeRes[ProfileStoryHighlightResponse] {
	persona, message := p.ownedPersona(ctx, userID, dto.PersonaID)
	if message != "" {
		return shared.PipeError[ProfileStoryHighlightResponse](message)
	}
	story, err := p.storyRepo.FindStoryByIDAny(ctx, storyID)
	if err != nil {
		if err == story_repo.ErrStoryNotFound {
			return shared.PipeError[ProfileStoryHighlightResponse](messages.Story_Not_Found)
		}
		return pipeInternalError[ProfileStoryHighlightResponse](err, "story.profile_highlight_story")
	}
	if story.OwnerPersonaID != persona.ID {
		return shared.PipeError[ProfileStoryHighlightResponse](messages.Forbidden)
	}
	highlight, err := p.storyRepo.AddProfileStoryHighlight(ctx, persona.ID, story.ID)
	if err != nil {
		return pipeInternalError[ProfileStoryHighlightResponse](err, "story.profile_highlight_add")
	}
	response, err := p.profileHighlightResponse(ctx, highlight, nil)
	if err != nil {
		return pipeInternalError[ProfileStoryHighlightResponse](err, "story.profile_highlight_response")
	}
	return shared.PipeSuccess(messages.Profile_Highlight_Story_Added, response)
}

// ListProfileStoryHighlightsPipe lists one persona's profile highlights for the current viewer.
func (p *StoryPipe) ListProfileStoryHighlightsPipe(ctx context.Context, personaID uuid.UUID, viewerUserID *uuid.UUID, viewerPersonaID *uuid.UUID) *shared.PipeRes[[]ProfileStoryHighlightResponse] {
	viewer, message := p.viewerPersona(ctx, viewerUserID, viewerPersonaID)
	if message != "" {
		return shared.PipeError[[]ProfileStoryHighlightResponse](message)
	}
	if viewer != nil {
		viewerPersonaID = &viewer.ID
	}
	highlights, err := p.storyRepo.FindProfileStoryHighlights(ctx, personaID)
	if err != nil {
		return pipeInternalError[[]ProfileStoryHighlightResponse](err, "story.profile_highlight_list")
	}
	responses := make([]ProfileStoryHighlightResponse, 0, len(highlights))
	for _, highlight := range highlights {
		response, err := p.profileHighlightResponse(ctx, highlight, viewerPersonaID)
		if err != nil {
			if err == story_repo.ErrStoryNotFound {
				continue
			}
			return pipeInternalError[[]ProfileStoryHighlightResponse](err, "story.profile_highlight_item")
		}
		responses = append(responses, *response)
	}
	return shared.PipeSuccess(messages.Profile_Highlight_Stories_Listed, &responses)
}

// RemoveProfileStoryHighlightPipe removes one story from one persona's profile highlights.
func (p *StoryPipe) RemoveProfileStoryHighlightPipe(ctx context.Context, userID uuid.UUID, storyID uuid.UUID, personaID uuid.UUID) *shared.PipeRes[any] {
	persona, message := p.ownedPersona(ctx, userID, personaID)
	if message != "" {
		return shared.PipeError[any](message)
	}
	story, err := p.storyRepo.FindStoryByIDAny(ctx, storyID)
	if err != nil {
		if err == story_repo.ErrStoryNotFound {
			return shared.PipeError[any](messages.Story_Not_Found)
		}
		return pipeInternalError[any](err, "story.profile_highlight_remove_story")
	}
	if story.OwnerPersonaID != persona.ID {
		return shared.PipeError[any](messages.Forbidden)
	}
	if err := p.storyRepo.RemoveProfileStoryHighlight(ctx, persona.ID, story.ID); err != nil {
		if err == story_repo.ErrStoryNotFound {
			return shared.PipeError[any](messages.Story_Not_Found)
		}
		return pipeInternalError[any](err, "story.profile_highlight_remove")
	}
	return shared.PipeSuccess[any](messages.Profile_Highlight_Story_Removed, nil)
}

// profileHighlightResponse maps one profile highlight into its API response shape.
func (p *StoryPipe) profileHighlightResponse(ctx context.Context, highlight *models.ProfileStoryHighlight, viewerPersonaID *uuid.UUID) (*ProfileStoryHighlightResponse, error) {
	story, err := p.storyRepo.FindStoryByIDAny(ctx, highlight.StoryID)
	if err != nil {
		return nil, err
	}
	if viewerPersonaID != nil {
		allowed, err := p.canView(ctx, story, viewerPersonaID)
		if err != nil {
			return nil, err
		}
		if !allowed {
			return nil, story_repo.ErrStoryNotFound
		}
	}
	storyResponse, err := p.storyResponse(ctx, story, viewerPersonaID, true, true)
	if err != nil {
		return nil, err
	}
	return &ProfileStoryHighlightResponse{
		ID:             highlight.ID.String(),
		OwnerPersonaID: highlight.OwnerPersonaID.String(),
		StoryID:        highlight.StoryID.String(),
		Position:       highlight.Position,
		Story:          storyResponse,
		CreatedAt:      highlight.CreatedAt,
	}, nil
}
