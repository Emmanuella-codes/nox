package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/models"
	notification_repo "github.com/emmanuella-codes/nox/repositories/notification"
	"github.com/google/uuid"
)

// createMessageNotifications persists recipient notifications for one new message.
func (p *MessagingPipe) createMessageNotifications(ctx context.Context, conversation *models.Conversation, message *models.Message) {
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
			ActorPersonaID:     message.SenderPersonaID,
			ConversationID:     &conversation.ID,
			MessageID:          &message.ID,
			NotificationType:   messageNotificationType(conversation.ConversationType),
		})
	}
	if len(inputs) == 0 {
		return
	}
	_, _ = p.notificationRepo.CreateNotifications(ctx, inputs)
}

// markConversationNotificationsRead marks conversation message notifications as read through one cursor.
func (p *MessagingPipe) markConversationNotificationsRead(ctx context.Context, userID uuid.UUID, personaID uuid.UUID, conversationID uuid.UUID, messageID uuid.UUID) {
	if p.notificationRepo == nil {
		return
	}
	_, _ = p.notificationRepo.MarkConversationNotificationsRead(ctx, userID, personaID, conversationID, messageID)
}

// deleteMessageNotifications removes notification rows tied to one deleted message.
func (p *MessagingPipe) deleteMessageNotifications(ctx context.Context, messageID uuid.UUID) {
	if p.notificationRepo == nil {
		return
	}
	_ = p.notificationRepo.DeleteMessageNotifications(ctx, messageID)
}

// messageNotificationType resolves the notification type for one conversation shape.
func messageNotificationType(conversationType models.ConversationType) models.NotificationType {
	if conversationType == models.GroupConversationType {
		return models.GroupMessageNotificationType
	}
	return models.DirectMessageNotificationType
}
