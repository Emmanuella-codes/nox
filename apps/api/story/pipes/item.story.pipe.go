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

	anonymousLabel := ""
	if dto.PostingMode == models.AnonymousPostingMode {
		anonymousLabel = "anonymous"
	}
	item, err := p.storyRepo.AddStoryItem(ctx, storyID, userID, asset.DurationSeconds, anonymousLabel, dto)
	if err != nil {
		if err == story_repo.ErrStoryDurationLimitExceeded {
			return shared.PipeError[StoryItemResponse](messages.Story_Duration_Limit_Exceeded)
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

func (p *StoryPipe) DeleteStoryItemPipe(ctx context.Context, userID uuid.UUID, storyID uuid.UUID, itemID uuid.UUID) *shared.PipeRes[any] {
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
		return shared.PipeError[any](messages.Forbidden)
	}
	if _, err := p.storyRepo.DeleteStoryItem(ctx, storyID, itemID); err != nil {
		if err == story_repo.ErrStoryItemNotFound {
			return shared.PipeError[any](messages.Story_Item_Not_Found)
		}
		return pipeInternalError[any](err, "story.delete_item")
	}
	return shared.PipeSuccess[any](messages.Story_Item_Deleted, nil)
}
