package pipes

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type NotificationPushPayload struct {
	NotificationID string            `json:"notification_id"`
	PersonaID      string            `json:"persona_id"`
	Type           string            `json:"type"`
	Title          string            `json:"title"`
	Body           string            `json:"body"`
	TargetPath     string            `json:"target_path"`
	Badge          int               `json:"badge"`
	Data           map[string]string `json:"data"`
}

func (p *NotificationPipe) enqueuePushDelivery(ctx context.Context, notification *models.Notification) {
	if p == nil || p.notificationRepo == nil {
		return
	}
	enabled, err := p.notificationRepo.PushEnabledForPersona(ctx, notification.RecipientPersonaID, notification.NotificationType)
	if err != nil || !enabled {
		return
	}
	payload, err := p.notificationPushPayload(ctx, notification)
	if err != nil {
		log.Error().Err(err).Str("notification_id", notification.ID.String()).Msg("failed to build notification push payload")
		return
	}
	bytes, err := json.Marshal(payload)
	if err != nil {
		return
	}
	_, _ = p.notificationRepo.EnqueueNotificationPush(ctx, notification, bytes)
}

func (p *NotificationPipe) notificationPushPayload(ctx context.Context, notification *models.Notification) (*NotificationPushPayload, error) {
	actorName := "Someone"
	if notification.ActorPostingMode == models.AnonymousPostingMode && notification.ActorAnonymousHandle != "" {
		actorName = notification.ActorAnonymousHandle
	}
	if notification.ActorPostingMode == models.PublicPostingMode && notification.ActorPersonaID != nil {
		persona, err := p.personaRepo.FindPersonaByID(ctx, *notification.ActorPersonaID)
		if err == nil && persona.DisplayName != "" {
			actorName = persona.DisplayName
		}
	}
	badge, err := p.notificationRepo.CountUnreadPersonaNotifications(ctx, notification.RecipientUserID, notification.RecipientPersonaID)
	if err != nil {
		return nil, err
	}
	title, body, targetPath := notificationCopy(notification, actorName)
	return &NotificationPushPayload{
		NotificationID: notification.ID.String(),
		PersonaID:      notification.RecipientPersonaID.String(),
		Type:           string(notification.NotificationType),
		Title:          title,
		Body:           body,
		TargetPath:     targetPath,
		Badge:          badge,
		Data: map[string]string{
			"notification_id": notification.ID.String(),
			"persona_id":      notification.RecipientPersonaID.String(),
			"type":            string(notification.NotificationType),
			"target_path":     targetPath,
		},
	}, nil
}

func notificationCopy(notification *models.Notification, actorName string) (string, string, string) {
	switch notification.NotificationType {
	case models.FollowNotificationType:
		return "New follower", fmt.Sprintf("%s followed you", actorName), "/notifications"
	case models.LikeNotificationType:
		return "New like", fmt.Sprintf("%s liked your post", actorName), "/notifications"
	case models.CommentNotificationType:
		return "New comment", fmt.Sprintf("%s commented on your post", actorName), "/notifications"
	case models.DirectMessageNotificationType:
		return "New message", fmt.Sprintf("%s sent you a message", actorName), pathOrFallback(notification.ConversationID, "/messages")
	case models.GroupMessageNotificationType:
		return "Group message", fmt.Sprintf("%s sent a group message", actorName), pathOrFallback(notification.ConversationID, "/messages")
	case models.StoryContributionRequestNotificationType:
		return "Story contribution", fmt.Sprintf("%s wants to add to your story", actorName), "/notifications"
	case models.StoryContributionAcceptedNotificationType:
		return "Contribution accepted", "Your story contribution was accepted", "/notifications"
	case models.StoryContributionRejectedNotificationType:
		return "Contribution rejected", "Your story contribution was rejected", "/notifications"
	case models.StoryHighlightAddedNotificationType:
		return "Story highlighted", fmt.Sprintf("%s added your story to event highlights", actorName), "/notifications"
	case models.StoryHighlightRemovedNotificationType:
		return "Story removed", fmt.Sprintf("%s removed your story from event highlights", actorName), "/notifications"
	case models.StoryReactionNotificationType:
		return "Story reaction", fmt.Sprintf("%s reacted to your story", actorName), "/notifications"
	default:
		return "Notification", "You have a new notification", "/notifications"
	}
}

func pathOrFallback(conversationID *uuid.UUID, fallback string) string {
	if conversationID == nil {
		return fallback
	}
	return "/messages/" + conversationID.String()
}
