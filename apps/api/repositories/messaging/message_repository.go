package messaging

import (
	"context"
	"errors"

	"github.com/emmanuella-codes/nox/messaging/dtos"
	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
)

// CreateMessage persists one message and its attachments.
func (r *pgRepository) CreateMessage(ctx context.Context, conversationID uuid.UUID, senderUserID uuid.UUID, dto dtos.SendMessageDTO) (*models.Message, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	attachments := normalizedAttachmentIDs(dto)
	row := tx.QueryRow(ctx, `
		INSERT INTO messages (conversation_id, sender_user_id, sender_persona_id, body, message_type, media_asset_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, conversation_id, sender_user_id, sender_persona_id, body, message_type,
		          media_asset_id, created_at, edited_at, deleted_at
	`, conversationID, senderUserID, dto.SenderPersonaID, dto.Body, dto.MessageType, firstAttachmentID(attachments))
	message, err := scanMessage(row)
	if err != nil {
		return nil, err
	}
	if err := insertMessageAttachments(ctx, tx, message.ID, attachments); err != nil {
		return nil, err
	}
	if err := refreshConversationLastMessage(ctx, tx, conversationID); err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return message, nil
}

// UpdateMessageBody edits one message body and marks it as edited.
func (r *pgRepository) UpdateMessageBody(ctx context.Context, messageID uuid.UUID, body string) (*models.Message, error) {
	row := r.db.QueryRow(ctx, `
		UPDATE messages
		SET body = $2,
		    edited_at = now()
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING id, conversation_id, sender_user_id, sender_persona_id, body, message_type,
		          media_asset_id, created_at, edited_at, deleted_at
	`, messageID, body)
	return scanMessage(row)
}

// FindMessagesByConversationID lists visible messages in one conversation.
func (r *pgRepository) FindMessagesByConversationID(ctx context.Context, conversationID uuid.UUID, limit int, offset int) ([]*models.Message, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, conversation_id, sender_user_id, sender_persona_id, body, message_type,
		       media_asset_id, created_at, edited_at, deleted_at
		FROM messages
		WHERE conversation_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`, conversationID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanMessages(rows)
}

// FindMessageByID fetches one message by id, including deleted messages.
func (r *pgRepository) FindMessageByID(ctx context.Context, messageID uuid.UUID) (*models.Message, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, conversation_id, sender_user_id, sender_persona_id, body, message_type,
		       media_asset_id, created_at, edited_at, deleted_at
		FROM messages
		WHERE id = $1
	`, messageID)
	return scanMessage(row)
}

// FindMessageAttachmentsByMessageIDs fetches attachments for a set of message ids.
func (r *pgRepository) FindMessageAttachmentsByMessageIDs(ctx context.Context, messageIDs []uuid.UUID) (map[uuid.UUID][]*models.MediaAsset, error) {
	attachmentsByMessage := make(map[uuid.UUID][]*models.MediaAsset, len(messageIDs))
	if len(messageIDs) == 0 {
		return attachmentsByMessage, nil
	}
	rows, err := r.db.Query(ctx, `
		SELECT ma.message_id, a.id, a.owner_user_id, a.owner_persona_id, a.media_kind, a.storage_key, a.playback_url,
		       COALESCE(a.thumbnail_url, ''), a.mime_type, a.duration_seconds, a.size_bytes, a.processing_status,
		       a.created_at, a.updated_at
		FROM message_attachments ma
		JOIN media_assets a ON a.id = ma.media_asset_id
		WHERE ma.message_id = ANY($1)
		ORDER BY ma.message_id, ma.position ASC
	`, messageIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var messageID uuid.UUID
		asset, err := scanMessageAttachmentAsset(rows, &messageID)
		if err != nil {
			return nil, err
		}
		attachmentsByMessage[messageID] = append(attachmentsByMessage[messageID], asset)
	}
	return attachmentsByMessage, rows.Err()
}

// MarkConversationRead advances the member read cursor for one conversation.
func (r *pgRepository) MarkConversationRead(ctx context.Context, conversationID uuid.UUID, personaID uuid.UUID, messageID uuid.UUID) (*models.ConversationMember, error) {
	row := r.db.QueryRow(ctx, `
		UPDATE conversation_members
		SET last_read_message_id = $3
		WHERE conversation_id = $1 AND persona_id = $2 AND left_at IS NULL
		  AND (
		    last_read_message_id IS NULL
		    OR (SELECT created_at FROM messages WHERE id = $3 AND conversation_id = $1) >
		       COALESCE((SELECT created_at FROM messages WHERE id = last_read_message_id), joined_at)
		  )
		RETURNING conversation_id, user_id, persona_id, role, last_read_message_id, joined_at, left_at
	`, conversationID, personaID, messageID)
	member, err := scanMember(row)
	if err == nil {
		return member, nil
	}
	if errors.Is(err, ErrMembershipNotFound) {
		return r.FindMember(ctx, conversationID, personaID)
	}
	return nil, err
}

// SoftDeleteMessage hides one message from all members and refreshes conversation state.
func (r *pgRepository) SoftDeleteMessage(ctx context.Context, messageID uuid.UUID) (*models.Message, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	row := tx.QueryRow(ctx, `
		UPDATE messages
		SET body = '',
		    media_asset_id = NULL,
		    deleted_at = COALESCE(deleted_at, now())
		WHERE id = $1
		RETURNING id, conversation_id, sender_user_id, sender_persona_id, body, message_type,
		          media_asset_id, created_at, edited_at, deleted_at
	`, messageID)
	message, err := scanMessage(row)
	if err != nil {
		return nil, err
	}
	if _, err := tx.Exec(ctx, `
		DELETE FROM message_attachments
		WHERE message_id = $1
	`, messageID); err != nil {
		return nil, err
	}
	if err := refreshConversationLastMessage(ctx, tx, message.ConversationID); err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return message, nil
}
