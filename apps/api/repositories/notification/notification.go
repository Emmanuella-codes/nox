package notification

import (
	"context"
	"encoding/json"
	"time"

	"github.com/emmanuella-codes/nox/models"
	notification_dtos "github.com/emmanuella-codes/nox/notification/dtos"
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
	StoryItemID                *uuid.UUID
	StoryContributionRequestID *uuid.UUID
	NotificationType           models.NotificationType
}

type NotificationRepository interface {
	CreateNotifications(ctx context.Context, inputs []CreateNotificationInput) ([]*models.Notification, error)
	FindPersonaNotifications(ctx context.Context, userID uuid.UUID, personaID uuid.UUID, limit int, offset int) ([]*models.Notification, error)
	CountUnreadPersonaNotifications(ctx context.Context, userID uuid.UUID, personaID uuid.UUID) (int, error)
	MarkNotificationRead(ctx context.Context, notificationID uuid.UUID, userID uuid.UUID, personaID uuid.UUID) (*models.Notification, error)
	MarkPersonaNotificationsRead(ctx context.Context, userID uuid.UUID, personaID uuid.UUID) (int64, error)
	MarkConversationNotificationsRead(ctx context.Context, userID uuid.UUID, personaID uuid.UUID, conversationID uuid.UUID, messageID uuid.UUID) (int64, error)
	DeleteMessageNotifications(ctx context.Context, messageID uuid.UUID) error
	UpsertNotificationDevice(ctx context.Context, userID uuid.UUID, dto notification_dtos.UpsertNotificationDeviceDTO) (*models.NotificationDevice, error)
	FindNotificationDevices(ctx context.Context, userID uuid.UUID) ([]*models.NotificationDevice, error)
	DisableNotificationDevice(ctx context.Context, userID uuid.UUID, deviceID uuid.UUID) error
	DisableNotificationDeviceByToken(ctx context.Context, pushToken string) error
	FindNotificationPreferences(ctx context.Context, personaID uuid.UUID) ([]*models.NotificationPreference, error)
	UpsertNotificationPreference(ctx context.Context, personaID uuid.UUID, notificationType models.NotificationType, pushEnabled bool) (*models.NotificationPreference, error)
	PushEnabledForPersona(ctx context.Context, personaID uuid.UUID, notificationType models.NotificationType) (bool, error)
	EnqueueNotificationPush(ctx context.Context, notification *models.Notification, payload json.RawMessage) (*models.NotificationOutbox, error)
	ClaimNotificationPushes(ctx context.Context, workerID string, limit int) ([]*models.NotificationOutbox, error)
	CanDeliverNotificationPush(ctx context.Context, outboxID uuid.UUID) (bool, string, error)
	MarkNotificationPushSent(ctx context.Context, outboxID uuid.UUID) error
	MarkNotificationPushRetry(ctx context.Context, outboxID uuid.UUID, nextAttemptAt time.Time, lastError string) error
	MarkNotificationPushDead(ctx context.Context, outboxID uuid.UUID, lastError string) error
	MarkNotificationPushSkipped(ctx context.Context, outboxID uuid.UUID, reason string) error
}

func NewNotificationRepository(db *pgxpool.Pool) NotificationRepository {
	return newPgRepository(db)
}
