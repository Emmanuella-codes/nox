package notification

import (
	"context"

	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
)

// CreateNotifications persists a batch of notification records.
func (r *pgRepository) CreateNotifications(ctx context.Context, inputs []CreateNotificationInput) ([]*models.Notification, error) {
	notifications := make([]*models.Notification, 0, len(inputs))
	for _, input := range inputs {
		row := r.db.QueryRow(ctx, `
			INSERT INTO notifications (
				recipient_user_id, recipient_persona_id, actor_persona_id, actor_posting_mode,
				actor_anonymous_handle, actor_anonymous_avatar_key, conversation_id,
				message_id, post_id, comment_id, notification_type
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
			ON CONFLICT DO NOTHING
			RETURNING id, recipient_user_id, recipient_persona_id, actor_persona_id, actor_posting_mode,
			          actor_anonymous_handle, actor_anonymous_avatar_key, conversation_id,
			          message_id, post_id, comment_id, is_read, read_at, notification_type, created_at
		`, input.RecipientUserID, input.RecipientPersonaID, input.ActorPersonaID, input.ActorPostingMode, input.ActorAnonymousHandle, input.ActorAnonymousAvatarKey, input.ConversationID, input.MessageID, input.PostID, input.CommentID, input.NotificationType)
		notification, err := scanNotification(row)
		if err != nil {
			if err == ErrNotificationNotFound {
				continue
			}
			return nil, err
		}
		notifications = append(notifications, notification)
	}
	return notifications, nil
}

// FindPersonaNotifications lists one persona's notifications.
func (r *pgRepository) FindPersonaNotifications(ctx context.Context, userID uuid.UUID, personaID uuid.UUID, limit int, offset int) ([]*models.Notification, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, recipient_user_id, recipient_persona_id, actor_persona_id, actor_posting_mode,
		       actor_anonymous_handle, actor_anonymous_avatar_key, conversation_id,
		       message_id, post_id, comment_id, is_read, read_at, notification_type, created_at
		FROM notifications
		WHERE recipient_user_id = $1 AND recipient_persona_id = $2
		ORDER BY created_at DESC, id DESC
		LIMIT $3 OFFSET $4
	`, userID, personaID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanNotifications(rows)
}

// CountUnreadPersonaNotifications returns the unread notification count for one persona.
func (r *pgRepository) CountUnreadPersonaNotifications(ctx context.Context, userID uuid.UUID, personaID uuid.UUID) (int, error) {
	var count int
	err := r.db.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM notifications
		WHERE recipient_user_id = $1 AND recipient_persona_id = $2 AND is_read = FALSE
	`, userID, personaID).Scan(&count)
	return count, err
}

// MarkNotificationRead marks one notification as read.
func (r *pgRepository) MarkNotificationRead(ctx context.Context, notificationID uuid.UUID, userID uuid.UUID, personaID uuid.UUID) (*models.Notification, error) {
	row := r.db.QueryRow(ctx, `
		UPDATE notifications
		SET is_read = TRUE,
		    read_at = COALESCE(read_at, now())
		WHERE id = $1 AND recipient_user_id = $2 AND recipient_persona_id = $3
		RETURNING id, recipient_user_id, recipient_persona_id, actor_persona_id, actor_posting_mode,
		          actor_anonymous_handle, actor_anonymous_avatar_key, conversation_id,
		          message_id, post_id, comment_id, is_read, read_at, notification_type, created_at
	`, notificationID, userID, personaID)
	return scanNotification(row)
}

// MarkPersonaNotificationsRead marks all persona notifications as read.
func (r *pgRepository) MarkPersonaNotificationsRead(ctx context.Context, userID uuid.UUID, personaID uuid.UUID) (int64, error) {
	tag, err := r.db.Exec(ctx, `
		UPDATE notifications
		SET is_read = TRUE,
		    read_at = COALESCE(read_at, now())
		WHERE recipient_user_id = $1 AND recipient_persona_id = $2 AND is_read = FALSE
	`, userID, personaID)
	return tag.RowsAffected(), err
}

// MarkConversationNotificationsRead marks message notifications read through one message cursor.
func (r *pgRepository) MarkConversationNotificationsRead(ctx context.Context, userID uuid.UUID, personaID uuid.UUID, conversationID uuid.UUID, messageID uuid.UUID) (int64, error) {
	tag, err := r.db.Exec(ctx, `
		UPDATE notifications n
		SET is_read = TRUE,
		    read_at = COALESCE(read_at, now())
		FROM messages target, messages current
		WHERE target.id = $4
		  AND current.id = n.message_id
		  AND n.recipient_user_id = $1
		  AND n.recipient_persona_id = $2
		  AND n.conversation_id = $3
		  AND n.message_id IS NOT NULL
		  AND n.is_read = FALSE
		  AND (
		    current.created_at < target.created_at
		    OR (current.created_at = target.created_at AND current.id::text <= target.id::text)
		  )
	`, userID, personaID, conversationID, messageID)
	return tag.RowsAffected(), err
}

// DeleteMessageNotifications removes notifications tied to one deleted message.
func (r *pgRepository) DeleteMessageNotifications(ctx context.Context, messageID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `
		DELETE FROM notifications
		WHERE message_id = $1
	`, messageID)
	return err
}
