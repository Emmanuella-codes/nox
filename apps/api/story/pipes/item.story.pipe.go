package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/models"
	event_repo "github.com/emmanuella-codes/nox/repositories/event"
	story_repo "github.com/emmanuella-codes/nox/repositories/story"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/emmanuella-codes/nox/story/dtos"
	"github.com/emmanuella-codes/nox/story/messages"
	"github.com/google/uuid"
)

func (p *StoryPipe) AddStoryItemPipe(ctx context.Context, userID uuid.UUID, storyID uuid.UUID, dto dtos.AddStoryItemDTO) *shared.PipeRes[StoryItemResponse] {
	if !validPostingMode(dto.PostingMode) {
		return shared.PipeError[StoryItemResponse](messages.Invalid_Story)
	}
	contributor, message := p.ownedPersona(ctx, userID, dto.ContributorPersonaID)
	if message != "" {
		return shared.PipeError[StoryItemResponse](message)
	}

	story, err := p.storyRepo.FindStoryByID(ctx, storyID)
	if err != nil {
		if err == story_repo.ErrStoryNotFound {
			return shared.PipeError[StoryItemResponse](messages.Story_Not_Found)
		}
		return pipeInternalError[StoryItemResponse](err, "story.find_for_item")
	}

	allowed, err := p.canContribute(ctx, story, contributor.ID)
	if err != nil {
		return pipeInternalError[StoryItemResponse](err, "story.contribution_permission")
	}
	if !allowed {
		return shared.PipeError[StoryItemResponse](messages.Forbidden)
	}

	asset, message := p.mediaAsset(ctx, dto.MediaAssetID)
	if message != "" {
		return shared.PipeError[StoryItemResponse](message)
	}
	if asset.OwnerUserID != userID || asset.OwnerPersonaID != contributor.ID || !validStoryVideo(asset) {
		return shared.PipeError[StoryItemResponse](messages.Invalid_Story)
	}

	label := ""
	if dto.PostingMode == models.AnonymousPostingMode {
		label = anonymousLabel(storyID, contributor.ID)
	}
	item, err := p.storyRepo.AddStoryItem(ctx, storyID, userID, asset.DurationSeconds, label, dto)
	if err != nil {
		if err == story_repo.ErrStoryDurationLimitExceeded {
			return shared.PipeError[StoryItemResponse](messages.Story_Duration_Limit_Exceeded)
		}
		if err == story_repo.ErrStoryMediaInUse {
			return shared.PipeError[StoryItemResponse](messages.Media_Asset_In_Use)
		}
		if err == story_repo.ErrStoryNotFound {
			return shared.PipeError[StoryItemResponse](messages.Story_Not_Found)
		}
		return pipeInternalError[StoryItemResponse](err, "story.add_item")
	}
	items, err := p.storyItemResponses(ctx, storyID)
	if err != nil {
		return pipeInternalError[StoryItemResponse](err, "story.item_response")
	}
	for _, response := range items {
		if response.ID == item.ID.String() {
			return shared.PipeSuccess(messages.Story_Item_Added, &response)
		}
	}
	return pipeInternalError[StoryItemResponse](story_repo.ErrStoryItemNotFound, "story.find_added_item")
}

func (p *StoryPipe) DeleteStoryItemPipe(ctx context.Context, userID uuid.UUID, storyID uuid.UUID, itemID uuid.UUID, moderatorPersonaID *uuid.UUID) *shared.PipeRes[any] {
	story, err := p.storyRepo.FindStoryByID(ctx, storyID)
	if err != nil {
		if err == story_repo.ErrStoryNotFound {
			return shared.PipeError[any](messages.Story_Not_Found)
		}
		return pipeInternalError[any](err, "story.find_for_delete_item")
	}
	items, err := p.storyRepo.FindStoryItems(ctx, storyID)
	if err != nil {
		return pipeInternalError[any](err, "story.items_for_delete")
	}
	var item *models.StoryItem
	for _, current := range items {
		if current.ID == itemID {
			item = current
			break
		}
	}
	if item == nil {
		return shared.PipeError[any](messages.Story_Item_Not_Found)
	}
	if story.OwnerUserID != userID && item.ContributorUserID != userID {
		if moderatorPersonaID == nil {
			return shared.PipeError[any](messages.Forbidden)
		}
		event, err := p.eventRepo.FindEventByID(ctx, story.EventID)
		if err != nil {
			if err == event_repo.ErrEventNotFound {
				return shared.PipeError[any](messages.Event_Not_Found)
			}
			return pipeInternalError[any](err, "story.delete_item_event")
		}
		moderator, message := p.ownedPersona(ctx, userID, *moderatorPersonaID)
		if message != "" {
			return shared.PipeError[any](message)
		}
		if moderator.ID != event.OrganizerID {
			return shared.PipeError[any](messages.Forbidden)
		}
	}
	if _, err := p.storyRepo.DeleteStoryItem(ctx, storyID, itemID); err != nil {
		if err == story_repo.ErrStoryItemNotFound {
			return shared.PipeError[any](messages.Story_Item_Not_Found)
		}
		return pipeInternalError[any](err, "story.delete_item")
	}
	return shared.PipeSuccess[any](messages.Story_Item_Deleted, nil)
}

func (p *StoryPipe) ReorderStoryItemPipe(ctx context.Context, userID uuid.UUID, storyID uuid.UUID, itemID uuid.UUID, dto dtos.ReorderStoryItemDTO) *shared.PipeRes[StoryItemResponse] {
	if dto.Position < 1 {
		return shared.PipeError[StoryItemResponse](messages.Invalid_Story)
	}
	story, err := p.storyRepo.FindStoryByID(ctx, storyID)
	if err != nil {
		if err == story_repo.ErrStoryNotFound {
			return shared.PipeError[StoryItemResponse](messages.Story_Not_Found)
		}
		return pipeInternalError[StoryItemResponse](err, "story.reorder_find")
	}
	persona, message := p.ownedPersona(ctx, userID, dto.PersonaID)
	if message != "" {
		return shared.PipeError[StoryItemResponse](message)
	}
	if story.OwnerPersonaID != persona.ID {
		return shared.PipeError[StoryItemResponse](messages.Forbidden)
	}
	item, err := p.storyRepo.ReorderStoryItem(ctx, storyID, itemID, dto.Position)
	if err != nil {
		if err == story_repo.ErrStoryItemNotFound {
			return shared.PipeError[StoryItemResponse](messages.Story_Item_Not_Found)
		}
		return pipeInternalError[StoryItemResponse](err, "story.reorder_item")
	}
	items, err := p.storyItemResponses(ctx, storyID)
	if err != nil {
		return pipeInternalError[StoryItemResponse](err, "story.reorder_response")
	}
	for _, response := range items {
		if response.ID == item.ID.String() {
			return shared.PipeSuccess(messages.Story_Item_Reordered, &response)
		}
	}
	return pipeInternalError[StoryItemResponse](story_repo.ErrStoryItemNotFound, "story.reorder_added_item")
}
