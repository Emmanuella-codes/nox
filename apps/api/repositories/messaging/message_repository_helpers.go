package messaging

import (
	"context"

	"github.com/emmanuella-codes/nox/messaging/dtos"
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

// returns the first attachment id for legacy message storage.
func firstAttachmentID(attachmentIDs []uuid.UUID) any {
	if len(attachmentIDs) == 0 {
		return nil
	}
	return attachmentIDs[0]
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
	if err := db.QueryRow(ctx, `
		SELECT id
		FROM messages
		WHERE conversation_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT 1
	`, conversationID).Scan(&lastMessageID); err != nil && err != pgx.ErrNoRows {
		return err
	}
	err := db.QueryRow(ctx, `
		UPDATE conversations
		SET last_message_id = $2,
		    updated_at = now()
		WHERE id = $1
		RETURNING id
	`, conversationID, lastMessageID).Scan(&conversationID)
	return err
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
