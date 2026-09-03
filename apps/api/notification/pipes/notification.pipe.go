package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/models"
	"github.com/emmanuella-codes/nox/notification/messages"
	notification_repo "github.com/emmanuella-codes/nox/repositories/notification"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/emmanuella-codes/nox/shared/realtime"
	"github.com/google/uuid"
)

func (p *NotificationPipe) ListNotificationsPipe(ctx context.Context, userID uuid.UUID, personaID uuid.UUID, limit int, offset int) *shared.PipeRes[NotificationListResponse] {
	if _, message := p.profilePersona(ctx, userID, personaID); message != "" {
		return shared.PipeError[NotificationListResponse](message)
	}
	limit = normalizeLimit(limit, 20, 50)
	if offset < 0 {
		offset = 0
	}
	notifications, err := p.notificationRepo.FindPersonaNotifications(ctx, userID, personaID, limit+1, offset)
	if err != nil {
		return pipeInternalError[NotificationListResponse](err, "notification.list")
	}
	unreadCount, err := p.notificationRepo.CountUnreadPersonaNotifications(ctx, userID, personaID)
	if err != nil {
		return pipeInternalError[NotificationListResponse](err, "notification.count_unread")
	}
	hasMore := len(notifications) > limit
	if hasMore {
		notifications = notifications[:limit]
	}
	response := NotificationListResponse{
		Limit:         limit,
		Offset:        offset,
		HasMore:       hasMore,
		NextOffset:    nextOffset(limit, offset, hasMore),
		UnreadCount:   unreadCount,
		Notifications: p.notificationResponses(ctx, notifications),
	}
	return shared.PipeSuccess(messages.Notifications_Listed, &response)
}

func (p *NotificationPipe) GetUnreadCountPipe(ctx context.Context, userID uuid.UUID, personaID uuid.UUID) *shared.PipeRes[NotificationUnreadCountResponse] {
	if _, message := p.profilePersona(ctx, userID, personaID); message != "" {
		return shared.PipeError[NotificationUnreadCountResponse](message)
	}
	unreadCount, err := p.notificationRepo.CountUnreadPersonaNotifications(ctx, userID, personaID)
	if err != nil {
		return pipeInternalError[NotificationUnreadCountResponse](err, "notification.get_unread_count")
	}
	response := NotificationUnreadCountResponse{UnreadCount: unreadCount}
	return shared.PipeSuccess(messages.Notifications_Listed, &response)
}

func (p *NotificationPipe) MarkNotificationReadPipe(ctx context.Context, userID uuid.UUID, notificationID uuid.UUID, personaID uuid.UUID) *shared.PipeRes[NotificationResponse] {
	if _, message := p.profilePersona(ctx, userID, personaID); message != "" {
		return shared.PipeError[NotificationResponse](message)
	}
	notification, err := p.notificationRepo.MarkNotificationRead(ctx, notificationID, userID, personaID)
	if err != nil {
		if err == notification_repo.ErrNotificationNotFound {
			return shared.PipeError[NotificationResponse](messages.Notification_Not_Found)
		}
		return pipeInternalError[NotificationResponse](err, "notification.mark_read")
	}
	response := p.notificationResponse(ctx, notification)
	p.publishPersonaEvent(ctx, userID, personaID, realtime.Event{Type: "notification.updated", Data: response})
	p.publishUnreadCount(ctx, userID, personaID)
	return shared.PipeSuccess(messages.Notification_Read, &response)
}

func (p *NotificationPipe) MarkAllNotificationsReadPipe(ctx context.Context, userID uuid.UUID, personaID uuid.UUID) *shared.PipeRes[any] {
	if _, message := p.profilePersona(ctx, userID, personaID); message != "" {
		return shared.PipeError[any](message)
	}
	if _, err := p.notificationRepo.MarkPersonaNotificationsRead(ctx, userID, personaID); err != nil {
		return pipeInternalError[any](err, "notification.mark_all_read")
	}
	p.publishPersonaEvent(ctx, userID, personaID, realtime.Event{
		Type: "notification.read_all",
		Data: map[string]string{"persona_id": personaID.String()},
	})
	p.publishUnreadCount(ctx, userID, personaID)
	return shared.PipeSuccess[any](messages.Notifications_Read, nil)
}

func normalizeLimit(limit int, fallback int, max int) int {
	if limit <= 0 {
		return fallback
	}
	if limit > max {
		return max
	}
	return limit
}

func nextOffset(limit int, offset int, hasMore bool) *int {
	if !hasMore {
		return nil
	}
	next := offset + limit
	return &next
}

func buildNotificationInput(recipient *models.Persona, actor *models.Persona, notificationType models.NotificationType) notification_repo.CreateNotificationInput {
	input := notification_repo.CreateNotificationInput{
		RecipientUserID:    recipient.UserID,
		RecipientPersonaID: recipient.ID,
		ActorPostingMode:   models.PublicPostingMode,
		NotificationType:   notificationType,
	}
	if actor != nil {
		input.ActorPersonaID = &actor.ID
	}
	return input
}
