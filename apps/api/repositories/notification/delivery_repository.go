package notification

import (
	"context"
	"errors"

	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// Checks whether a queued push is still valid to deliver.
func (r *pgRepository) CanDeliverNotificationPush(ctx context.Context, outboxID uuid.UUID) (bool, string, error) {
	type deliveryState struct {
		notificationType   models.NotificationType
		conversationID     *uuid.UUID
		messageID          *uuid.UUID
		recipientPersonaID uuid.UUID
	}

	var state deliveryState
	err := r.db.QueryRow(ctx, `
		SELECT n.notification_type, n.conversation_id, n.message_id, n.recipient_persona_id
		FROM notification_outbox o
		JOIN notifications n ON n.id = o.notification_id
		WHERE o.id = $1
	`, outboxID).Scan(&state.notificationType, &state.conversationID, &state.messageID, &state.recipientPersonaID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, "missing_notification", nil
		}
		return false, "", err
	}
	if state.notificationType != models.DirectMessageNotificationType && state.notificationType != models.GroupMessageNotificationType {
		return true, "", nil
	}
	if state.conversationID == nil || state.messageID == nil {
		return false, "invalid_message_notification", nil
	}

	var deleted bool
	err = r.db.QueryRow(ctx, `
		SELECT deleted_at IS NOT NULL
		FROM messages
		WHERE id = $1
	`, *state.messageID).Scan(&deleted)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, "message_missing", nil
		}
		return false, "", err
	}
	if deleted {
		return false, "message_deleted", nil
	}

	var membershipExists bool
	err = r.db.QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1
			FROM conversation_members
			WHERE conversation_id = $1 AND persona_id = $2 AND left_at IS NULL
		)
	`, *state.conversationID, state.recipientPersonaID).Scan(&membershipExists)
	if err != nil {
		return false, "", err
	}
	if !membershipExists {
		return false, "recipient_left", nil
	}

	var alreadyRead bool
	err = r.db.QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1
			FROM conversation_members cm
			JOIN messages current_message ON current_message.id = cm.last_read_message_id
			JOIN messages target_message ON target_message.id = $3
			WHERE cm.conversation_id = $1
			  AND cm.persona_id = $2
			  AND cm.left_at IS NULL
			  AND (
				current_message.created_at > target_message.created_at
				OR (current_message.created_at = target_message.created_at AND current_message.id::text >= target_message.id::text)
			  )
		)
	`, *state.conversationID, state.recipientPersonaID, *state.messageID).Scan(&alreadyRead)
	if err != nil {
		return false, "", err
	}
	if alreadyRead {
		return false, "already_read", nil
	}
	return true, "", nil
}
