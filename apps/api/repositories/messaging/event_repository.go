package messaging

import (
	"context"
	"encoding/json"

	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
)

// stores one replayable conversation event.
func (r *pgRepository) AppendConversationEvent(ctx context.Context, conversationID uuid.UUID, actorUserID uuid.UUID, eventType string, messageID *uuid.UUID, payload json.RawMessage) (*models.ConversationEvent, error) {
	row := r.db.QueryRow(ctx, `
		INSERT INTO conversation_events (conversation_id, actor_user_id, event_type, message_id, payload)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, conversation_id, actor_user_id, event_type, message_id, payload, created_at
	`, conversationID, actorUserID, eventType, messageID, payload)
	return scanConversationEvent(row)
}

// lists replayable conversation events for one active user after a cursor.
func (r *pgRepository) FindConversationEventsAfter(ctx context.Context, userID uuid.UUID, afterID int64, limit int) ([]*models.ConversationEvent, error) {
	rows, err := r.db.Query(ctx, `
		SELECT ce.id, ce.conversation_id, ce.actor_user_id, ce.event_type, ce.message_id, ce.payload, ce.created_at
		FROM conversation_events ce
		WHERE ce.id > $2
		  AND EXISTS (
		    SELECT 1
		    FROM conversation_members cm
		    WHERE cm.conversation_id = ce.conversation_id
		      AND cm.user_id = $1
		      AND cm.left_at IS NULL
		  )
		ORDER BY ce.id ASC
		LIMIT $3
	`, userID, afterID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanConversationEvents(rows)
}
