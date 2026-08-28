package notification

import (
	"context"

	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateNotificationInput struct {
	RecipientUserID            uuid.UUID
	RecipientPersonaID         uuid.UUID
	ActorPersonaID             *uuid.UUID
	ActorPostingMode           models.PostingMode
	ActorAnonymousHandle       string
	ActorAnonymousAvatarKey    string
	ConversationID             *uuid.UUID
	MessageID                  *uuid.UUID
	PostID                     *uuid.UUID
	CommentID                  *uuid.UUID
	EventID                    *uuid.UUID
	StoryID                    *uuid.UUID
	StoryContributionRequestID *uuid.UUID
	NotificationType           models.NotificationType
}

type NotificationRepository interface {
	// CreateNotifications persists a batch of notification records.
	CreateNotifications(ctx context.Context, inputs []CreateNotificationInput) ([]*models.Notification, error)
	// FindPersonaNotifications lists one persona's notifications.
	FindPersonaNotifications(ctx context.Context, userID uuid.UUID, personaID uuid.UUID, limit int, offset int) ([]*models.Notification, error)
	// CountUnreadPersonaNotifications returns the unread notification count for one persona.
	CountUnreadPersonaNotifications(ctx context.Context, userID uuid.UUID, personaID uuid.UUID) (int, error)
	// MarkNotificationRead marks one notification as read.
	MarkNotificationRead(ctx context.Context, notificationID uuid.UUID, userID uuid.UUID, personaID uuid.UUID) (*models.Notification, error)
	// MarkPersonaNotificationsRead marks all persona notifications as read.
	MarkPersonaNotificationsRead(ctx context.Context, userID uuid.UUID, personaID uuid.UUID) (int64, error)
	// MarkConversationNotificationsRead marks message notifications read through one message cursor.
	MarkConversationNotificationsRead(ctx context.Context, userID uuid.UUID, personaID uuid.UUID, conversationID uuid.UUID, messageID uuid.UUID) (int64, error)
	// DeleteMessageNotifications removes notifications tied to one deleted message.
	DeleteMessageNotifications(ctx context.Context, messageID uuid.UUID) error
}

// NewNotificationRepository builds the notification repository from a database pool.
func NewNotificationRepository(db *pgxpool.Pool) NotificationRepository {
	return newPgRepository(db)
}
