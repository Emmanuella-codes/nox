package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/notification/messages"
	notification_repo "github.com/emmanuella-codes/nox/repositories/notification"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

// ListNotificationsPipe lists notifications for one owned persona.
func (p *NotificationPipe) ListNotificationsPipe(ctx context.Context, userID uuid.UUID, personaID uuid.UUID, limit int, offset int) *shared.PipeRes[[]NotificationResponse] {
	if _, message := p.profilePersona(ctx, userID, personaID); message != "" {
		return shared.PipeError[[]NotificationResponse](message)
	}
	if limit <= 0 {
		limit = 20
	}
	if limit > 50 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	notifications, err := p.notificationRepo.FindPersonaNotifications(ctx, userID, personaID, limit, offset)
	if err != nil {
		return pipeInternalError[[]NotificationResponse](err, "notification.list")
	}
	responses := p.notificationResponses(ctx, notifications)
	return shared.PipeSuccess(messages.Notifications_Listed, &responses)
}

// MarkNotificationReadPipe marks one notification as read for the current persona.
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
	return shared.PipeSuccess(messages.Notification_Read, &response)
}

// MarkAllNotificationsReadPipe marks all notifications as read for one owned persona.
func (p *NotificationPipe) MarkAllNotificationsReadPipe(ctx context.Context, userID uuid.UUID, personaID uuid.UUID) *shared.PipeRes[any] {
	if _, message := p.profilePersona(ctx, userID, personaID); message != "" {
		return shared.PipeError[any](message)
	}
	if _, err := p.notificationRepo.MarkPersonaNotificationsRead(ctx, userID, personaID); err != nil {
		return pipeInternalError[any](err, "notification.mark_all_read")
	}
	return shared.PipeSuccess[any](messages.Notifications_Read, nil)
}
