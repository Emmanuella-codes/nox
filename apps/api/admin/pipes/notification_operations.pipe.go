package pipes

import (
	"context"
	"errors"
	"strings"

	"github.com/emmanuella-codes/nox/admin/dtos"
	"github.com/emmanuella-codes/nox/admin/messages"
	"github.com/emmanuella-codes/nox/models"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (p *AdminPipe) ListNotificationDevicesPipe(ctx context.Context, actorID uuid.UUID, userID *uuid.UUID, limit, offset int) *shared.PipeRes[[]models.AdminNotificationDevice] {
	if message := p.requireActiveAdmin(ctx, actorID); message != "" {
		return shared.PipeError[[]models.AdminNotificationDevice](shared.CreatePipeMessage(message))
	}
	limit, offset = notificationPage(limit, offset)
	rows, err := p.db.Query(ctx, `SELECT id, user_id, install_id, platform, CASE WHEN length(push_token) > 8 THEN left(push_token, 4) || '...' || right(push_token, 4) ELSE '***' END, app_version, last_seen_at, disabled_at, created_at, updated_at
		FROM notification_devices WHERE ($1::uuid IS NULL OR user_id = $1) ORDER BY last_seen_at DESC, created_at DESC LIMIT $2 OFFSET $3`, userID, limit, offset)
	if err != nil {
		logInternalError(err, "notifications.devices.list")
		return shared.PipeError[[]models.AdminNotificationDevice](messages.Internal_Error)
	}
	defer rows.Close()
	items := make([]models.AdminNotificationDevice, 0)
	for rows.Next() {
		var item models.AdminNotificationDevice
		if err := rows.Scan(&item.ID, &item.UserID, &item.InstallID, &item.Platform, &item.TokenHint, &item.AppVersion, &item.LastSeenAt, &item.DisabledAt, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return shared.PipeError[[]models.AdminNotificationDevice](messages.Internal_Error)
		}
		items = append(items, item)
	}
	p.audit(ctx, actorID, "admin.notifications.devices.list", map[string]any{"user_id": nullableUUID(userID)})
	return shared.PipeSuccess(messages.Admin_Notification_Devices_Loaded, &items)
}

func (p *AdminPipe) ListNotificationOutboxPipe(ctx context.Context, actorID uuid.UUID, status string, userID *uuid.UUID, limit, offset int) *shared.PipeRes[[]models.AdminNotificationOutbox] {
	if message := p.requireActiveAdmin(ctx, actorID); message != "" {
		return shared.PipeError[[]models.AdminNotificationOutbox](shared.CreatePipeMessage(message))
	}
	limit, offset = notificationPage(limit, offset)
	rows, err := p.db.Query(ctx, `SELECT o.id, o.notification_id, o.recipient_user_id, o.recipient_persona_id, n.notification_type, o.channel, o.status, o.attempt_count, o.next_attempt_at, o.last_error, o.worker_id, o.claimed_at, o.sent_at, o.created_at, o.updated_at
		FROM notification_outbox o JOIN notifications n ON n.id = o.notification_id
		WHERE ($1 = '' OR o.status = $1) AND ($2::uuid IS NULL OR o.recipient_user_id = $2)
		ORDER BY CASE WHEN o.status IN ('failed', 'dead') THEN 0 ELSE 1 END, o.updated_at DESC LIMIT $3 OFFSET $4`, status, userID, limit, offset)
	if err != nil {
		logInternalError(err, "notifications.outbox.list")
		return shared.PipeError[[]models.AdminNotificationOutbox](messages.Internal_Error)
	}
	defer rows.Close()
	items := make([]models.AdminNotificationOutbox, 0)
	for rows.Next() {
		var item models.AdminNotificationOutbox
		if err := rows.Scan(&item.ID, &item.NotificationID, &item.RecipientUserID, &item.RecipientPersonaID, &item.NotificationType, &item.Channel, &item.Status, &item.AttemptCount, &item.NextAttemptAt, &item.LastError, &item.WorkerID, &item.ClaimedAt, &item.SentAt, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return shared.PipeError[[]models.AdminNotificationOutbox](messages.Internal_Error)
		}
		item.PayloadRedacted = true
		items = append(items, item)
	}
	p.audit(ctx, actorID, "admin.notifications.outbox.list", map[string]any{"status": status, "user_id": nullableUUID(userID)})
	return shared.PipeSuccess(messages.Admin_Notification_Outbox_Loaded, &items)
}

func (p *AdminPipe) RetryNotificationOutboxPipe(ctx context.Context, actorID, outboxID uuid.UUID, dto dtos.NotificationAdminActionDTO) *shared.PipeRes[models.AdminNotificationOutbox] {
	if message := p.requireActiveAdmin(ctx, actorID); message != "" {
		return shared.PipeError[models.AdminNotificationOutbox](shared.CreatePipeMessage(message))
	}
	dto.Reason = strings.TrimSpace(dto.Reason)
	if dto.Reason == "" {
		return shared.PipeError[models.AdminNotificationOutbox](messages.Invalid_Payload)
	}
	var updatedID uuid.UUID
	err := p.db.QueryRow(ctx, `UPDATE notification_outbox SET status = 'pending', next_attempt_at = now(), worker_id = '', claimed_at = NULL, updated_at = now() WHERE id = $1 AND status IN ('failed', 'dead') RETURNING id`, outboxID).Scan(&updatedID)
	if errors.Is(err, pgx.ErrNoRows) {
		return shared.PipeError[models.AdminNotificationOutbox](messages.Notification_Outbox_Not_Found)
	}
	if err != nil {
		logInternalError(err, "notifications.outbox.retry")
		return shared.PipeError[models.AdminNotificationOutbox](messages.Internal_Error)
	}
	item, err := p.findAdminOutbox(ctx, updatedID)
	if err != nil {
		logInternalError(err, "notifications.outbox.retry_read")
		return shared.PipeError[models.AdminNotificationOutbox](messages.Internal_Error)
	}
	p.audit(ctx, actorID, "admin.notifications.outbox.retry", map[string]any{"outbox_id": outboxID.String(), "reason": dto.Reason})
	return shared.PipeSuccess(messages.Admin_Notification_Retried, item)
}

func (p *AdminPipe) DisableNotificationDevicePipe(ctx context.Context, actorID, deviceID uuid.UUID, dto dtos.NotificationAdminActionDTO) *shared.PipeRes[any] {
	if message := p.requireActiveAdmin(ctx, actorID); message != "" {
		return shared.PipeError[any](shared.CreatePipeMessage(message))
	}
	dto.Reason = strings.TrimSpace(dto.Reason)
	if dto.Reason == "" {
		return shared.PipeError[any](messages.Invalid_Payload)
	}
	var userID uuid.UUID
	err := p.db.QueryRow(ctx, `UPDATE notification_devices SET disabled_at = COALESCE(disabled_at, now()), updated_at = now() WHERE id = $1 RETURNING user_id`, deviceID).Scan(&userID)
	if errors.Is(err, pgx.ErrNoRows) {
		return shared.PipeError[any](messages.Notification_Device_Not_Found)
	}
	if err != nil {
		logInternalError(err, "notifications.device.disable")
		return shared.PipeError[any](messages.Internal_Error)
	}
	p.audit(ctx, actorID, "admin.notifications.device.disable", map[string]any{"device_id": deviceID.String(), "user_id": userID.String(), "reason": dto.Reason})
	return shared.PipeSuccess[any](messages.Admin_Notification_Device_Disabled, nil)
}

func (p *AdminPipe) findAdminOutbox(ctx context.Context, outboxID uuid.UUID) (*models.AdminNotificationOutbox, error) {
	var item models.AdminNotificationOutbox
	err := p.db.QueryRow(ctx, `SELECT o.id, o.notification_id, o.recipient_user_id, o.recipient_persona_id, n.notification_type, o.channel, o.status, o.attempt_count, o.next_attempt_at, o.last_error, o.worker_id, o.claimed_at, o.sent_at, o.created_at, o.updated_at FROM notification_outbox o JOIN notifications n ON n.id = o.notification_id WHERE o.id = $1`, outboxID).Scan(&item.ID, &item.NotificationID, &item.RecipientUserID, &item.RecipientPersonaID, &item.NotificationType, &item.Channel, &item.Status, &item.AttemptCount, &item.NextAttemptAt, &item.LastError, &item.WorkerID, &item.ClaimedAt, &item.SentAt, &item.CreatedAt, &item.UpdatedAt)
	item.PayloadRedacted = true
	return &item, err
}

func notificationPage(limit, offset int) (int, int) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}

func nullableUUID(value *uuid.UUID) any {
	if value == nil {
		return nil
	}
	return value.String()
}
