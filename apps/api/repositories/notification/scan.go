package notification

import (
	"errors"

	"github.com/emmanuella-codes/nox/models"
	"github.com/jackc/pgx/v5"
)

var ErrNotificationNotFound = errors.New("notification not found")

type scanner interface {
	Scan(dest ...any) error
}

func scanNotification(row scanner) (*models.Notification, error) {
	var notification models.Notification
	err := row.Scan(
		&notification.ID,
		&notification.RecipientUserID,
		&notification.RecipientPersonaID,
		&notification.ActorPersonaID,
		&notification.ActorPostingMode,
		&notification.ActorAnonymousHandle,
		&notification.ActorAnonymousAvatarKey,
		&notification.ConversationID,
		&notification.MessageID,
		&notification.PostID,
		&notification.CommentID,
		&notification.EventID,
		&notification.StoryID,
		&notification.StoryItemID,
		&notification.StoryContributionRequestID,
		&notification.IsRead,
		&notification.ReadAt,
		&notification.NotificationType,
		&notification.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotificationNotFound
		}
		return nil, err
	}
	return &notification, nil
}

func scanNotifications(rows pgx.Rows) ([]*models.Notification, error) {
	notifications := make([]*models.Notification, 0)
	for rows.Next() {
		notification, err := scanNotification(rows)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}
	return notifications, rows.Err()
}

func scanNotificationDevice(row scanner) (*models.NotificationDevice, error) {
	var device models.NotificationDevice
	if err := row.Scan(&device.ID, &device.UserID, &device.InstallID, &device.Platform, &device.PushToken, &device.AppVersion, &device.LastSeenAt, &device.DisabledAt, &device.CreatedAt, &device.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotificationNotFound
		}
		return nil, err
	}
	return &device, nil
}

func scanNotificationDevices(rows pgx.Rows) ([]*models.NotificationDevice, error) {
	devices := make([]*models.NotificationDevice, 0)
	for rows.Next() {
		device, err := scanNotificationDevice(rows)
		if err != nil {
			return nil, err
		}
		devices = append(devices, device)
	}
	return devices, rows.Err()
}

func scanNotificationPreference(row scanner) (*models.NotificationPreference, error) {
	var preference models.NotificationPreference
	if err := row.Scan(&preference.PersonaID, &preference.NotificationType, &preference.PushEnabled, &preference.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotificationNotFound
		}
		return nil, err
	}
	return &preference, nil
}

func scanNotificationPreferences(rows pgx.Rows) ([]*models.NotificationPreference, error) {
	preferences := make([]*models.NotificationPreference, 0)
	for rows.Next() {
		preference, err := scanNotificationPreference(rows)
		if err != nil {
			return nil, err
		}
		preferences = append(preferences, preference)
	}
	return preferences, rows.Err()
}

func scanNotificationOutbox(row scanner) (*models.NotificationOutbox, error) {
	var outbox models.NotificationOutbox
	if err := row.Scan(&outbox.ID, &outbox.NotificationID, &outbox.RecipientUserID, &outbox.RecipientPersonaID, &outbox.Channel, &outbox.Status, &outbox.Payload, &outbox.AttemptCount, &outbox.NextAttemptAt, &outbox.LastError, &outbox.WorkerID, &outbox.ClaimedAt, &outbox.SentAt, &outbox.CreatedAt, &outbox.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotificationNotFound
		}
		return nil, err
	}
	return &outbox, nil
}

func scanNotificationOutboxes(rows pgx.Rows) ([]*models.NotificationOutbox, error) {
	outboxes := make([]*models.NotificationOutbox, 0)
	for rows.Next() {
		outbox, err := scanNotificationOutbox(rows)
		if err != nil {
			return nil, err
		}
		outboxes = append(outboxes, outbox)
	}
	return outboxes, rows.Err()
}
