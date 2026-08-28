package messaging

import (
	"context"
	"strings"
	"time"

	"github.com/emmanuella-codes/nox/messaging/dtos"
	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// merges legacy and multi-attachment inputs into one ordered list.
func normalizedAttachmentIDs(dto dtos.SendMessageDTO) []uuid.UUID {
	ids := make([]uuid.UUID, 0, len(dto.MediaAssetIDs)+1)
	seen := map[uuid.UUID]bool{}
	if dto.MediaAssetID != nil && *dto.MediaAssetID != uuid.Nil {
		seen[*dto.MediaAssetID] = true
		ids = append(ids, *dto.MediaAssetID)
	}
	for _, mediaAssetID := range dto.MediaAssetIDs {
		if mediaAssetID == uuid.Nil || seen[mediaAssetID] {
			continue
		}
		seen[mediaAssetID] = true
		ids = append(ids, mediaAssetID)
	}
	return ids
}

// normalizedIdempotencyKey trims and bounds one message request key.
func normalizedIdempotencyKey(value string) string {
	return strings.TrimSpace(value)
}

// returns the first attachment id for legacy message storage.
func firstAttachmentID(attachmentIDs []uuid.UUID) any {
	if len(attachmentIDs) == 0 {
		return nil
	}
	return attachmentIDs[0]
}

// findMessageByRequestKey resolves one existing message for a repeated send request.
func findMessageByRequestKey(ctx context.Context, db execQuerier, conversationID uuid.UUID, senderUserID uuid.UUID, key string) (*models.Message, error) {
	row := db.QueryRow(ctx, `
		SELECT m.id, m.conversation_id, m.sender_user_id, m.sender_persona_id, m.body, m.message_type,
		       m.media_asset_id, m.created_at, m.edited_at, m.deleted_at
		FROM message_request_keys mrk
		JOIN messages m ON m.id = mrk.message_id
		WHERE mrk.conversation_id = $1
		  AND mrk.sender_user_id = $2
		  AND mrk.idempotency_key = $3
	`, conversationID, senderUserID, key)
	return scanMessage(row)
}

// storeMessageRequestKey persists one idempotent send mapping for later retries.
func storeMessageRequestKey(ctx context.Context, db execQuerier, conversationID uuid.UUID, senderUserID uuid.UUID, key string, messageID uuid.UUID) error {
	return db.QueryRow(ctx, `
		INSERT INTO message_request_keys (conversation_id, sender_user_id, idempotency_key, message_id)
		VALUES ($1, $2, $3, $4)
		RETURNING message_id
	`, conversationID, senderUserID, key, messageID).Scan(&messageID)
}

// persists attachment rows for one message.
func insertMessageAttachments(ctx context.Context, db execQuerier, messageID uuid.UUID, attachmentIDs []uuid.UUID) error {
	for position, mediaAssetID := range attachmentIDs {
		if err := db.QueryRow(ctx, `
			INSERT INTO message_attachments (message_id, media_asset_id, position)
			VALUES ($1, $2, $3)
			RETURNING message_id
		`, messageID, mediaAssetID, position).Scan(&messageID); err != nil {
			return err
		}
	}
	return nil
}

// updates the conversation pointer to the latest visible message.
func refreshConversationLastMessage(ctx context.Context, db execQuerier, conversationID uuid.UUID) error {
	var lastMessageID *uuid.UUID
	var lastMessageCreatedAt *time.Time
	if err := db.QueryRow(ctx, `
		SELECT id, created_at
		FROM messages
		WHERE conversation_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC, id DESC
		LIMIT 1
	`, conversationID).Scan(&lastMessageID, &lastMessageCreatedAt); err != nil && err != pgx.ErrNoRows {
		return err
	}
	err := db.QueryRow(ctx, `
		UPDATE conversations
		SET last_message_id = $2,
		    updated_at = COALESCE($3::timestamptz, created_at)
		WHERE id = $1
		RETURNING id
	`, conversationID, lastMessageID, lastMessageCreatedAt).Scan(&conversationID)
	return err
}

// messageCursorAfter reports whether one message is newer than another message cursor.
func messageCursorAfter(candidate *models.Message, current *models.Message) bool {
	if candidate.CreatedAt.After(current.CreatedAt) {
		return true
	}
	if candidate.CreatedAt.Before(current.CreatedAt) {
		return false
	}
	return candidate.ID.String() > current.ID.String()
}

// deletes a conversation when no active members remain.
func dissolveConversationIfEmpty(ctx context.Context, db execQuerier, conversationID uuid.UUID) error {
	var activeCount int
	if err := db.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM conversation_members
		WHERE conversation_id = $1 AND left_at IS NULL
	`, conversationID).Scan(&activeCount); err != nil {
		return err
	}
	if activeCount > 0 {
		return nil
	}
	err := db.QueryRow(ctx, `
		DELETE FROM conversations
		WHERE id = $1
		RETURNING id
	`, conversationID).Scan(&conversationID)
	if err == pgx.ErrNoRows {
		return nil
	}
	return err
}
