package messaging

import (
	"context"
	"encoding/json"

	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type conversationScanner interface {
	Scan(dest ...any) error
}

type execQuerier interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// inserts one conversation and returns the created row.
func insertConversation(ctx context.Context, db execQuerier, conversationType models.ConversationType, title string, createdBy uuid.UUID) (*models.Conversation, error) {
	row := db.QueryRow(ctx, `
		INSERT INTO conversations (conversation_type, title, created_by)
		VALUES ($1, $2, $3)
		RETURNING id, conversation_type, title, created_by, last_message_id, created_at, updated_at
	`, conversationType, title, createdBy)
	return scanConversation(row)
}

// scans one conversation row into the model shape.
func scanConversation(scanner conversationScanner) (*models.Conversation, error) {
	var conversation models.Conversation
	err := scanner.Scan(
		&conversation.ID,
		&conversation.ConversationType,
		&conversation.Title,
		&conversation.CreatedBy,
		&conversation.LastMessageID,
		&conversation.CreatedAt,
		&conversation.UpdatedAt,
	)
	if err != nil {
		return nil, mapNotFound(err, ErrConversationNotFound)
	}
	return &conversation, nil
}

// scans one conversation row with unread metadata.
func scanConversationListItem(scanner conversationScanner) (*models.Conversation, int, error) {
	var conversation models.Conversation
	var unreadCount int
	err := scanner.Scan(
		&conversation.ID,
		&conversation.ConversationType,
		&conversation.Title,
		&conversation.CreatedBy,
		&conversation.LastMessageID,
		&conversation.CreatedAt,
		&conversation.UpdatedAt,
		&unreadCount,
	)
	if err != nil {
		return nil, 0, err
	}
	return &conversation, unreadCount, nil
}

// inserts or reactivates one conversation member row.
func insertMember(ctx context.Context, db execQuerier, conversationID uuid.UUID, userID uuid.UUID, personaID uuid.UUID, role models.ConversationMemberRole) (*models.ConversationMember, error) {
	row := db.QueryRow(ctx, `
		INSERT INTO conversation_members (conversation_id, user_id, persona_id, role)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (conversation_id, persona_id)
		DO UPDATE SET left_at = NULL, joined_at = now(), last_read_message_id = NULL, role = EXCLUDED.role
		RETURNING conversation_id, user_id, persona_id, role, last_read_message_id, joined_at, left_at
	`, conversationID, userID, personaID, role)
	return scanMember(row)
}

// scans one conversation member row into the model shape.
func scanMember(scanner conversationScanner) (*models.ConversationMember, error) {
	var member models.ConversationMember
	err := scanner.Scan(
		&member.ConversationID,
		&member.UserID,
		&member.PersonaID,
		&member.Role,
		&member.LastReadMessageID,
		&member.JoinedAt,
		&member.LeftAt,
	)
	if err != nil {
		return nil, mapNotFound(err, ErrMembershipNotFound)
	}
	return &member, nil
}

// scans many conversation member rows into model values.
func scanMembers(rows pgx.Rows) ([]*models.ConversationMember, error) {
	members := make([]*models.ConversationMember, 0)
	for rows.Next() {
		member, err := scanMember(rows)
		if err != nil {
			return nil, err
		}
		members = append(members, member)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return members, nil
}

// scans one message row into the model shape.
func scanMessage(scanner conversationScanner) (*models.Message, error) {
	var message models.Message
	err := scanner.Scan(
		&message.ID,
		&message.ConversationID,
		&message.SenderUserID,
		&message.SenderPersonaID,
		&message.Body,
		&message.MessageType,
		&message.MediaAssetID,
		&message.StoryID,
		&message.StoryItemID,
		&message.CreatedAt,
		&message.EditedAt,
		&message.DeletedAt,
	)
	if err != nil {
		return nil, mapNotFound(err, ErrMessageNotFound)
	}
	return &message, nil
}

// scans many message rows into model values.
func scanMessages(rows pgx.Rows) ([]*models.Message, error) {
	messages := make([]*models.Message, 0)
	for rows.Next() {
		message, err := scanMessage(rows)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return messages, nil
}

// scans one attachment join row into a media asset.
func scanMessageAttachmentAsset(scanner conversationScanner, messageID *uuid.UUID) (*models.MediaAsset, error) {
	var asset models.MediaAsset
	err := scanner.Scan(
		messageID,
		&asset.ID,
		&asset.OwnerUserID,
		&asset.OwnerPersonaID,
		&asset.MediaKind,
		&asset.StorageKey,
		&asset.PlaybackURL,
		&asset.ThumbnailURL,
		&asset.MimeType,
		&asset.DurationSeconds,
		&asset.SizeBytes,
		&asset.ProcessingStatus,
		&asset.CreatedAt,
		&asset.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

// scans one-column uuid row sets into a slice.
func scanUUIDs(rows pgx.Rows) ([]uuid.UUID, error) {
	ids := make([]uuid.UUID, 0)
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

// scans one event-log row into the model shape.
func scanConversationEvent(scanner conversationScanner) (*models.ConversationEvent, error) {
	var event models.ConversationEvent
	var payload []byte
	err := scanner.Scan(
		&event.ID,
		&event.ConversationID,
		&event.ActorUserID,
		&event.EventType,
		&event.MessageID,
		&payload,
		&event.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	event.Payload = json.RawMessage(payload)
	return &event, nil
}

// scans many event-log rows into model values.
func scanConversationEvents(rows pgx.Rows) ([]*models.ConversationEvent, error) {
	events := make([]*models.ConversationEvent, 0)
	for rows.Next() {
		event, err := scanConversationEvent(rows)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, rows.Err()
}
