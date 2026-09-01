package pipes

import (
	"context"
	"strings"

	messaging_pipes "github.com/emmanuella-codes/nox/messaging/pipes"
	"github.com/emmanuella-codes/nox/models"
	story_repo "github.com/emmanuella-codes/nox/repositories/story"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/emmanuella-codes/nox/story/dtos"
	"github.com/emmanuella-codes/nox/story/messages"
	"github.com/google/uuid"
)

// ViewStoryItemPipe records one story item view for the current viewer.
func (p *StoryPipe) ViewStoryItemPipe(ctx context.Context, userID uuid.UUID, storyID uuid.UUID, itemID uuid.UUID, dto dtos.StoryItemViewDTO) *shared.PipeRes[StoryItemResponse] {
	viewer, message := p.ownedPersona(ctx, userID, dto.ViewerPersonaID)
	if message != "" {
		return shared.PipeError[StoryItemResponse](message)
	}
	story, item, failure := p.loadViewableStoryItem(ctx, storyID, itemID, &viewer.ID)
	if failure != nil {
		return failure
	}
	if _, _, err := p.storyRepo.UpsertStoryItemView(ctx, story.ID, item.ID, userID, viewer.ID); err != nil {
		return pipeInternalError[StoryItemResponse](err, "story.view_item")
	}
	response, err := p.storyItemResponse(ctx, item, &viewer.ID)
	if err != nil {
		return pipeInternalError[StoryItemResponse](err, "story.view_item_response")
	}
	return shared.PipeSuccess(messages.Story_Item_Viewed, response)
}

// ListStoryItemViewersPipe lists viewer identities for one owner-visible story item.
func (p *StoryPipe) ListStoryItemViewersPipe(ctx context.Context, userID uuid.UUID, storyID uuid.UUID, itemID uuid.UUID) *shared.PipeRes[StoryViewersResponse] {
	story, err := p.storyRepo.FindStoryByIDAny(ctx, storyID)
	if err != nil {
		if err == story_repo.ErrStoryNotFound {
			return shared.PipeError[StoryViewersResponse](messages.Story_Not_Found)
		}
		return pipeInternalError[StoryViewersResponse](err, "story.find_viewers_story")
	}
	if story.OwnerUserID != userID {
		return shared.PipeError[StoryViewersResponse](messages.Forbidden)
	}
	item, err := p.storyRepo.FindStoryItemByID(ctx, storyID, itemID)
	if err != nil {
		if err == story_repo.ErrStoryItemNotFound {
			return shared.PipeError[StoryViewersResponse](messages.Story_Item_Not_Found)
		}
		return pipeInternalError[StoryViewersResponse](err, "story.find_viewers_item")
	}
	viewerIDs, err := p.storyRepo.FindStoryItemViewerPersonaIDs(ctx, storyID, itemID)
	if err != nil {
		return pipeInternalError[StoryViewersResponse](err, "story.find_viewers")
	}
	viewers := make([]PersonaResponse, 0, len(viewerIDs))
	for _, viewerID := range viewerIDs {
		persona, err := p.personaRepo.FindPersonaByID(ctx, viewerID)
		if err != nil {
			return pipeInternalError[StoryViewersResponse](err, "story.viewer_persona")
		}
		viewers = append(viewers, personaResponse(persona))
	}
	response := &StoryViewersResponse{
		StoryID:     story.ID.String(),
		StoryItemID: item.ID.String(),
		ViewCount:   len(viewers),
		Viewers:     viewers,
	}
	return shared.PipeSuccess(messages.Story_Item_Viewers_Listed, response)
}

// ReactToStoryItemPipe upserts one viewer reaction on one story item.
func (p *StoryPipe) ReactToStoryItemPipe(ctx context.Context, userID uuid.UUID, storyID uuid.UUID, itemID uuid.UUID, dto dtos.StoryItemReactionDTO) *shared.PipeRes[StoryItemResponse] {
	if !validStoryReactionType(dto.ReactionType) {
		return shared.PipeError[StoryItemResponse](messages.Invalid_Payload)
	}
	reactor, message := p.ownedPersona(ctx, userID, dto.PersonaID)
	if message != "" {
		return shared.PipeError[StoryItemResponse](message)
	}
	story, item, failure := p.loadViewableStoryItem(ctx, storyID, itemID, &reactor.ID)
	if failure != nil {
		return failure
	}
	if _, err := p.storyRepo.UpsertStoryItemReaction(ctx, story.ID, item.ID, userID, reactor.ID, dto.ReactionType); err != nil {
		return pipeInternalError[StoryItemResponse](err, "story.react_item")
	}
	if err := p.notifyStoryReaction(ctx, story, item, reactor); err != nil {
		return pipeInternalError[StoryItemResponse](err, "story.react_notify")
	}
	response, err := p.storyItemResponse(ctx, item, &reactor.ID)
	if err != nil {
		return pipeInternalError[StoryItemResponse](err, "story.react_response")
	}
	return shared.PipeSuccess(messages.Story_Item_Reaction_Updated, response)
}

// RemoveStoryItemReactionPipe removes one viewer reaction from one story item.
func (p *StoryPipe) RemoveStoryItemReactionPipe(ctx context.Context, userID uuid.UUID, storyID uuid.UUID, itemID uuid.UUID, personaID uuid.UUID) *shared.PipeRes[StoryItemResponse] {
	reactor, message := p.ownedPersona(ctx, userID, personaID)
	if message != "" {
		return shared.PipeError[StoryItemResponse](message)
	}
	story, item, failure := p.loadViewableStoryItem(ctx, storyID, itemID, &reactor.ID)
	if failure != nil {
		return failure
	}
	if err := p.storyRepo.DeleteStoryItemReaction(ctx, story.ID, item.ID, reactor.ID); err != nil {
		return pipeInternalError[StoryItemResponse](err, "story.remove_reaction")
	}
	p.deleteStoryReactionNotification(ctx, story, item, reactor)
	response, err := p.storyItemResponse(ctx, item, &reactor.ID)
	if err != nil {
		return pipeInternalError[StoryItemResponse](err, "story.remove_reaction_response")
	}
	return shared.PipeSuccess(messages.Story_Item_Reaction_Removed, response)
}

// ReplyToStoryItemPipe sends one story reply as a direct message to the story owner.
func (p *StoryPipe) ReplyToStoryItemPipe(ctx context.Context, userID uuid.UUID, storyID uuid.UUID, itemID uuid.UUID, dto dtos.StoryItemReplyDTO) *shared.PipeRes[messaging_pipes.MessageResponse] {
	replier, message := p.ownedPersona(ctx, userID, dto.PersonaID)
	if message != "" {
		return shared.PipeError[messaging_pipes.MessageResponse](message)
	}
	story, item, failure := p.loadViewableStoryItem(ctx, storyID, itemID, &replier.ID)
	if failure != nil {
		return failureMessageResponse(failure)
	}
	body := strings.TrimSpace(dto.Body)
	if body == "" {
		return shared.PipeError[messaging_pipes.MessageResponse](messages.Invalid_Payload)
	}
	owner, err := p.personaRepo.FindPersonaByID(ctx, story.OwnerPersonaID)
	if err != nil {
		return pipeInternalError[messaging_pipes.MessageResponse](err, "story.reply_owner")
	}
	_, messageModel, err := p.createStoryReplyMessage(ctx, replier, owner, story.ID, item.ID, body, strings.TrimSpace(dto.IdempotencyKey))
	if err != nil {
		return pipeInternalError[messaging_pipes.MessageResponse](err, "story.reply_message")
	}
	response := storyMessageResponse(messageModel)
	return shared.PipeSuccess(messages.Story_Item_Replied, &response)
}

// loadViewableStoryItem loads one story and item pair and enforces story visibility.
func (p *StoryPipe) loadViewableStoryItem(ctx context.Context, storyID uuid.UUID, itemID uuid.UUID, viewerPersonaID *uuid.UUID) (*models.Story, *models.StoryItem, *shared.PipeRes[StoryItemResponse]) {
	story, err := p.storyRepo.FindStoryByID(ctx, storyID)
	if err != nil {
		if err == story_repo.ErrStoryNotFound {
			return nil, nil, shared.PipeError[StoryItemResponse](messages.Story_Not_Found)
		}
		return nil, nil, pipeInternalError[StoryItemResponse](err, "story.find_item_story")
	}
	allowed, err := p.canView(ctx, story, viewerPersonaID)
	if err != nil {
		return nil, nil, pipeInternalError[StoryItemResponse](err, "story.item_view")
	}
	if !allowed {
		return nil, nil, shared.PipeError[StoryItemResponse](messages.Story_Not_Found)
	}
	item, err := p.storyRepo.FindStoryItemByID(ctx, storyID, itemID)
	if err != nil {
		if err == story_repo.ErrStoryItemNotFound {
			return nil, nil, shared.PipeError[StoryItemResponse](messages.Story_Item_Not_Found)
		}
		return nil, nil, pipeInternalError[StoryItemResponse](err, "story.find_item")
	}
	return story, item, nil
}

// storyItemResponse maps one story item into a response with viewer state hydrated.
func (p *StoryPipe) storyItemResponse(ctx context.Context, item *models.StoryItem, viewerPersonaID *uuid.UUID) (*StoryItemResponse, error) {
	responses, err := p.storyItemResponses(ctx, item.StoryID, viewerPersonaID)
	if err != nil {
		return nil, err
	}
	for _, response := range responses {
		if response.ID == item.ID.String() {
			return &response, nil
		}
	}
	return nil, story_repo.ErrStoryItemNotFound
}

// failureMessageResponse converts one story item pipe failure into a message-shaped failure.
func failureMessageResponse(res *shared.PipeRes[StoryItemResponse]) *shared.PipeRes[messaging_pipes.MessageResponse] {
	if res.Success {
		return shared.PipeSuccess[messaging_pipes.MessageResponse](res.Message, nil)
	}
	return shared.PipeError[messaging_pipes.MessageResponse](res.Message)
}
