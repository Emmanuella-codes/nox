package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/models"
	"github.com/emmanuella-codes/nox/shared/realtime"
	"github.com/google/uuid"
)

// PublishCreatedNotification broadcasts one created notification and refreshed unread count.
func (p *NotificationPipe) PublishCreatedNotification(ctx context.Context, notification *models.Notification) {
	if p == nil {
		return
	}
	p.publishPersonaEvent(ctx, notification.RecipientUserID, notification.RecipientPersonaID, realtime.Event{
		Type: "notification.created",
		Data: p.notificationResponse(ctx, notification),
	})
	p.publishUnreadCount(ctx, notification.RecipientUserID, notification.RecipientPersonaID)
}

// publishPersonaEvent sends one realtime event to a user's notification subscribers.
func (p *NotificationPipe) publishPersonaEvent(ctx context.Context, userID uuid.UUID, personaID uuid.UUID, event realtime.Event) {
	if p.realtimeHub == nil {
		return
	}
	_ = p.realtimeHub.PublishUsers([]uuid.UUID{userID}, event)
}

// publishUnreadCount broadcasts the latest unread count for one persona.
func (p *NotificationPipe) publishUnreadCount(ctx context.Context, userID uuid.UUID, personaID uuid.UUID) {
	if p.realtimeHub == nil || p.notificationRepo == nil {
		return
	}
	unreadCount, err := p.notificationRepo.CountUnreadPersonaNotifications(ctx, userID, personaID)
	if err != nil {
		return
	}
	p.publishPersonaEvent(ctx, userID, personaID, realtime.Event{
		Type: "notification.unread_count",
		Data: NotificationUnreadCountResponse{UnreadCount: unreadCount},
	})
}
