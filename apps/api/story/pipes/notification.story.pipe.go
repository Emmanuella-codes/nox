package pipes

import (
	"context"
	"encoding/json"
	"strconv"

	messaging_dtos "github.com/emmanuella-codes/nox/messaging/dtos"
	"github.com/emmanuella-codes/nox/models"
	messaging_repo "github.com/emmanuella-codes/nox/repositories/messaging"
	notification_repo "github.com/emmanuella-codes/nox/repositories/notification"
	"github.com/emmanuella-codes/nox/shared/realtime"
	"github.com/google/uuid"
)

// createNotifications persists story notifications and broadcasts any created rows.
func (p *StoryPipe) createNotifications(ctx context.Context, inputs []notification_repo.CreateNotificationInput) error {
	if p.notificationRepo == nil || len(inputs) == 0 {
		return nil
	}
	notifications, err := p.notificationRepo.CreateNotifications(ctx, inputs)
	if err != nil {
		return err
	}
	if p.notificationPublisher == nil {
		return nil
	}
	for _, notification := range notifications {
		p.notificationPublisher.PublishCreatedNotification(ctx, notification)
	}
	return nil
}

// notifyStoryContributionRequest publishes one pending contribution request notification to the story owner.
func (p *StoryPipe) notifyStoryContributionRequest(ctx context.Context, story *models.Story, contributor *models.Persona, request *models.StoryContributionRequest) error {
	if story.OwnerUserID == contributor.UserID {
		return nil
	}
	return p.createNotifications(ctx, []notification_repo.CreateNotificationInput{{
		RecipientUserID:            story.OwnerUserID,
		RecipientPersonaID:         story.OwnerPersonaID,
		ActorPersonaID:             &contributor.ID,
		ActorPostingMode:           models.PublicPostingMode,
		StoryID:                    &story.ID,
		StoryContributionRequestID: &request.ID,
		NotificationType:           models.StoryContributionRequestNotificationType,
	}})
}

// notifyStoryContributionDecision publishes one accepted or rejected contribution decision back to the contributor.
func (p *StoryPipe) notifyStoryContributionDecision(ctx context.Context, story *models.Story, contributor *models.Persona, notificationType models.NotificationType, request *models.StoryContributionRequest) error {
	if story.OwnerUserID == contributor.UserID {
		return nil
	}
	return p.createNotifications(ctx, []notification_repo.CreateNotificationInput{{
		RecipientUserID:            contributor.UserID,
		RecipientPersonaID:         contributor.ID,
		ActorPersonaID:             &story.OwnerPersonaID,
		ActorPostingMode:           models.PublicPostingMode,
		StoryID:                    &story.ID,
		StoryContributionRequestID: &request.ID,
		NotificationType:           notificationType,
	}})
}

// notifyStoryHighlightChange publishes one highlight add or remove notification to the story owner.
func (p *StoryPipe) notifyStoryHighlightChange(ctx context.Context, story *models.Story, actor *models.Persona, notificationType models.NotificationType) error {
	if story.OwnerUserID == actor.UserID {
		return nil
	}
	return p.createNotifications(ctx, []notification_repo.CreateNotificationInput{{
		RecipientUserID:    story.OwnerUserID,
		RecipientPersonaID: story.OwnerPersonaID,
		ActorPersonaID:     &actor.ID,
		ActorPostingMode:   models.PublicPostingMode,
		EventID:            &story.EventID,
		StoryID:            &story.ID,
		NotificationType:   notificationType,
	}})
}

// notifyStoryReaction publishes one story reaction notification to the story owner.
func (p *StoryPipe) notifyStoryReaction(ctx context.Context, story *models.Story, item *models.StoryItem, actor *models.Persona) error {
	if story.OwnerUserID == actor.UserID {
		return nil
	}
	return p.createNotifications(ctx, []notification_repo.CreateNotificationInput{{
		RecipientUserID:    story.OwnerUserID,
		RecipientPersonaID: story.OwnerPersonaID,
		ActorPersonaID:     &actor.ID,
		ActorPostingMode:   models.PublicPostingMode,
		StoryID:            &story.ID,
		StoryItemID:        &item.ID,
		NotificationType:   models.StoryReactionNotificationType,
	}})
}

// deleteStoryReactionNotification removes one story reaction notification when the reaction is undone.
func (p *StoryPipe) deleteStoryReactionNotification(ctx context.Context, story *models.Story, item *models.StoryItem, actor *models.Persona) {
	if p.notificationRepo == nil || story.OwnerUserID == actor.UserID {
		return
	}
	repo, ok := p.notificationRepo.(interface {
		DeleteStoryReactionNotification(context.Context, uuid.UUID, uuid.UUID, uuid.UUID) error
	})
	if !ok {
		return
	}
	_ = repo.DeleteStoryReactionNotification(ctx, story.OwnerPersonaID, actor.ID, item.ID)
}

// createStoryReplyMessage persists one story reply as a direct message and broadcasts the standard messaging side effects.
func (p *StoryPipe) createStoryReplyMessage(ctx context.Context, replier *models.Persona, owner *models.Persona, storyID uuid.UUID, itemID uuid.UUID, body string, idempotencyKey string) (*models.Conversation, *models.Message, error) {
	if p.messagingRepo == nil {
		return nil, nil, messaging_repo.ErrConversationNotFound
	}
	conversation, err := p.messagingRepo.FindDirectConversationBetweenPersonas(ctx, replier.ID, owner.ID)
	if err == messaging_repo.ErrConversationNotFound {
		conversation, err = p.messagingRepo.CreateDirectConversation(ctx, replier, owner)
	}
	if err != nil {
		return nil, nil, err
	}
	message, createdNew, err := p.messagingRepo.CreateMessage(ctx, conversation.ID, replier.UserID, messaging_dtos.SendMessageDTO{
		SenderPersonaID: replier.ID,
		Body:            body,
		MessageType:     models.TextMessageType,
		StoryID:         &storyID,
		StoryItemID:     &itemID,
		IdempotencyKey:  idempotencyKey,
	})
	if err != nil {
		return nil, nil, err
	}
	if createdNew {
		p.createDirectMessageNotification(ctx, conversation, message)
		p.publishStoryReplyEvent(ctx, replier.UserID, message)
	}
	return conversation, message, nil
}

// createDirectMessageNotification persists the standard direct-message notification for one story reply.
func (p *StoryPipe) createDirectMessageNotification(ctx context.Context, conversation *models.Conversation, message *models.Message) {
	if p.notificationRepo == nil || p.messagingRepo == nil {
		return
	}
	members, err := p.messagingRepo.FindConversationMembers(ctx, conversation.ID)
	if err != nil {
		return
	}
	inputs := make([]notification_repo.CreateNotificationInput, 0, len(members))
	for _, member := range members {
		if member.UserID == message.SenderUserID {
			continue
		}
		inputs = append(inputs, notification_repo.CreateNotificationInput{
			RecipientUserID:    member.UserID,
			RecipientPersonaID: member.PersonaID,
			ActorPersonaID:     &message.SenderPersonaID,
			ActorPostingMode:   models.PublicPostingMode,
			ConversationID:     &conversation.ID,
			MessageID:          &message.ID,
			NotificationType:   models.DirectMessageNotificationType,
		})
	}
	_ = p.createNotifications(ctx, inputs)
}

// publishStoryReplyEvent appends and broadcasts one durable direct-message event for a story reply.
func (p *StoryPipe) publishStoryReplyEvent(ctx context.Context, actorUserID uuid.UUID, message *models.Message) {
	if p.messagingRepo == nil {
		return
	}
	payload := map[string]any{
		"conversation_id": message.ConversationID.String(),
		"message":         storyMessageResponse(message),
	}
	bytes, err := json.Marshal(payload)
	if err != nil {
		return
	}
	stored, err := p.messagingRepo.AppendConversationEvent(ctx, message.ConversationID, actorUserID, "message.created", &message.ID, bytes)
	if err != nil || p.realtimeHub == nil {
		return
	}
	userIDs, err := p.messagingRepo.FindConversationMemberUserIDs(ctx, message.ConversationID)
	if err != nil || len(userIDs) == 0 {
		return
	}
	_ = p.realtimeHub.PublishUsers(userIDs, realtime.Event{
		ID:        strconv.FormatInt(stored.ID, 10),
		Type:      "message.created",
		Data:      payload,
		CreatedAt: stored.CreatedAt.Format(messageTimeFormat),
	})
}
