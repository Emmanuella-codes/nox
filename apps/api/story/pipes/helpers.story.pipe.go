package pipes

import (
	messaging_pipes "github.com/emmanuella-codes/nox/messaging/pipes"
	"github.com/emmanuella-codes/nox/models"
)

const (
	storyItemImageDurationSeconds = 5
	maxStoryItemMediaDuration     = 120
)

func validContributionMode(mode models.StoryContributionMode) bool {
	return mode == models.PublicStoryContributionMode || mode == models.PrivateStoryContributionMode
}

func validPostingMode(mode models.PostingMode) bool {
	return mode == models.PublicPostingMode
}

func validStoryMedia(asset *models.MediaAsset) bool {
	if asset == nil {
		return false
	}
	if asset.ProcessingStatus != models.ReadyMediaStatus {
		return false
	}
	if asset.MediaKind == models.ImageMediaKind {
		return asset.DurationSeconds == storyItemImageDurationSeconds
	}
	return asset.MediaKind == models.VideoMediaKind &&
		asset.ProcessingStatus == models.ReadyMediaStatus &&
		asset.DurationSeconds > 0 &&
		asset.DurationSeconds <= maxStoryItemMediaDuration
}

func storyItemDuration(asset *models.MediaAsset) int {
	if asset == nil {
		return 0
	}
	if asset.MediaKind == models.ImageMediaKind {
		return storyItemImageDurationSeconds
	}
	return asset.DurationSeconds
}

func validStoryReactionType(reactionType models.StoryReactionType) bool {
	return reactionType == models.StoryReactionTypeLike ||
		reactionType == models.StoryReactionTypeFire ||
		reactionType == models.StoryReactionTypeHeart ||
		reactionType == models.StoryReactionTypeLaugh
}

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
