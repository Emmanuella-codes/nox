package messaging

import (
	"context"
	"errors"

	"github.com/emmanuella-codes/nox/messaging/dtos"
	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

// persists one message and its attachments or reuses an idempotent send.
func (r *pgRepository) CreateMessage(ctx context.Context, conversationID uuid.UUID, senderUserID uuid.UUID, dto dtos.SendMessageDTO) (*models.Message, bool, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, false, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	idempotencyKey := normalizedIdempotencyKey(dto.IdempotencyKey)
	if idempotencyKey != "" {
		existing, err := findMessageByRequestKey(ctx, tx, conversationID, senderUserID, idempotencyKey)
		if err == nil {
			return existing, false, nil
		}
		if !errors.Is(err, ErrMessageNotFound) {
			return nil, false, err
		}
	}
	attachments := normalizedAttachmentIDs(dto)
	row := tx.QueryRow(ctx, `
		INSERT INTO messages (conversation_id, sender_user_id, sender_persona_id, body, message_type, media_asset_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, conversation_id, sender_user_id, sender_persona_id, body, message_type,
		          media_asset_id, created_at, edited_at, deleted_at
	`, conversationID, senderUserID, dto.SenderPersonaID, dto.Body, dto.MessageType, firstAttachmentID(attachments))
	message, err := scanMessage(row)
	if err != nil {
		return nil, false, err
	}
	if err := insertMessageAttachments(ctx, tx, message.ID, attachments); err != nil {
		return nil, false, err
	}
	if idempotencyKey != "" {
		if err := storeMessageRequestKey(ctx, tx, conversationID, senderUserID, idempotencyKey, message.ID); err != nil {
			var pgErr *pgconn.PgError
			if !errors.As(err, &pgErr) || pgErr.Code != "23505" {
				return nil, false, err
			}
			existing, existingErr := findMessageByRequestKey(ctx, r.db, conversationID, senderUserID, idempotencyKey)
			if existingErr != nil {
				return nil, false, existingErr
			}
			return existing, false, nil
		}
	}
	if err := refreshConversationLastMessage(ctx, tx, conversationID); err != nil {
		return nil, false, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, false, err
	}
	return message, true, nil
}

// edits one message body and marks it as edited.
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

// lists visible messages in one conversation.
func (r *pgRepository) FindMessagesByConversationID(ctx context.Context, conversationID uuid.UUID, limit int, offset int) ([]*models.Message, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, conversation_id, sender_user_id, sender_persona_id, body, message_type,
		       media_asset_id, created_at, edited_at, deleted_at
		FROM messages
		WHERE conversation_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC, id DESC
		LIMIT $2 OFFSET $3
	`, conversationID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanMessages(rows)
}

// fetches one message by id, including deleted messages.
func (r *pgRepository) FindMessageByID(ctx context.Context, messageID uuid.UUID) (*models.Message, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, conversation_id, sender_user_id, sender_persona_id, body, message_type,
		       media_asset_id, created_at, edited_at, deleted_at
		FROM messages
		WHERE id = $1
	`, messageID)
	return scanMessage(row)
}

// fetches attachments for a set of message ids.
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

// advances the member read cursor for one conversation.
func (r *pgRepository) MarkConversationRead(ctx context.Context, conversationID uuid.UUID, personaID uuid.UUID, messageID uuid.UUID) (*models.ConversationMember, error) {
	member, err := r.FindMember(ctx, conversationID, personaID)
	if err != nil {
		return nil, err
	}
	candidate, err := r.FindMessageByID(ctx, messageID)
	if err != nil {
		return nil, err
	}
	if candidate.ConversationID != conversationID || candidate.DeletedAt != nil {
		return member, nil
	}
	if member.LastReadMessageID != nil {
		current, err := r.FindMessageByID(ctx, *member.LastReadMessageID)
		if err != nil && !errors.Is(err, ErrMessageNotFound) {
			return nil, err
		}
		if err == nil && !messageCursorAfter(candidate, current) {
			return member, nil
		}
	}
	row := r.db.QueryRow(ctx, `
		UPDATE conversation_members
		SET last_read_message_id = $3
		WHERE conversation_id = $1 AND persona_id = $2 AND left_at IS NULL
		RETURNING conversation_id, user_id, persona_id, role, last_read_message_id, joined_at, left_at
	`, conversationID, personaID, messageID)
	return scanMember(row)
}

// hides one message from all members and refreshes conversation state.
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
