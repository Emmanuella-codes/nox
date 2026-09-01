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

// scanNotification scans one notification row into the model shape.
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

// scanNotifications scans many notification rows into model values.
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
