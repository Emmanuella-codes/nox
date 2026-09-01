package pipes

import (
	messaging_pipes "github.com/emmanuella-codes/nox/messaging/pipes"
	"github.com/emmanuella-codes/nox/models"
)

// validContributionMode validates the supported story contribution modes.
func validContributionMode(mode models.StoryContributionMode) bool {
	return mode == models.PublicStoryContributionMode || mode == models.PrivateStoryContributionMode
}

// validPostingMode validates the supported direct story item posting modes.
func validPostingMode(mode models.PostingMode) bool {
	return mode == models.PublicPostingMode
}

// validStoryVideo verifies that one media asset is usable as a story item.
func validStoryVideo(asset *models.MediaAsset) bool {
	return asset.MediaKind == models.VideoMediaKind &&
		asset.ProcessingStatus == models.ReadyMediaStatus &&
		asset.DurationSeconds > 0 &&
		asset.DurationSeconds <= 300
}

// validStoryReactionType validates the supported story reaction types.
func validStoryReactionType(reactionType models.StoryReactionType) bool {
	return reactionType == models.StoryReactionTypeLike ||
		reactionType == models.StoryReactionTypeFire ||
		reactionType == models.StoryReactionTypeHeart ||
		reactionType == models.StoryReactionTypeLaugh
}

// storyMessageResponse maps one story-linked direct message into the shared messaging response shape.
func storyMessageResponse(message *models.Message) messaging_pipes.MessageResponse {
	var storyID *string
	if message.StoryID != nil {
		value := message.StoryID.String()
		storyID = &value
	}
	var storyItemID *string
	if message.StoryItemID != nil {
		value := message.StoryItemID.String()
		storyItemID = &value
	}
	return messaging_pipes.MessageResponse{
		ID:              message.ID.String(),
		ConversationID:  message.ConversationID.String(),
		SenderPersonaID: message.SenderPersonaID.String(),
		Body:            message.Body,
		MessageType:     message.MessageType,
		Attachments:     []*models.MediaAsset{},
		StoryID:         storyID,
		StoryItemID:     storyItemID,
		Deleted:         message.DeletedAt != nil,
		Edited:          message.EditedAt != nil,
		CreatedAt:       message.CreatedAt.Format(messageTimeFormat),
	}
}

const messageTimeFormat = "2006-01-02T15:04:05.999999999Z07:00"
