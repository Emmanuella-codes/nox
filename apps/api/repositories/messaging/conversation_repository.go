package messaging

import (
	"context"
	"errors"

	"github.com/emmanuella-codes/nox/messaging/dtos"
	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
)

// CreateDirectConversation creates one direct conversation between two profiles.
func (r *pgRepository) CreateDirectConversation(ctx context.Context, creator *models.Persona, recipient *models.Persona) (*models.Conversation, error) {
	personaAID, personaBID := orderedPersonaIDs(creator.ID, recipient.ID)
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	conversation, err := insertConversation(ctx, tx, models.DirectConversationType, "", creator.ID)
	if err != nil {
		return nil, err
	}
	if _, err := tx.Exec(ctx, `
		INSERT INTO direct_conversations (conversation_id, persona_a_id, persona_b_id)
		VALUES ($1, $2, $3)
	`, conversation.ID, personaAID, personaBID); err != nil {
		return nil, err
	}
	if _, err := insertMember(ctx, tx, conversation.ID, creator.UserID, creator.ID, models.ConversationMemberRoleMember); err != nil {
		return nil, err
	}
	if recipient.ID != creator.ID {
		if _, err := insertMember(ctx, tx, conversation.ID, recipient.UserID, recipient.ID, models.ConversationMemberRoleMember); err != nil {
			return nil, err
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return conversation, nil
}

// FindDirectConversationBetweenPersonas finds the direct conversation for two profiles.
func (r *pgRepository) FindDirectConversationBetweenPersonas(ctx context.Context, personaAID uuid.UUID, personaBID uuid.UUID) (*models.Conversation, error) {
	personaAID, personaBID = orderedPersonaIDs(personaAID, personaBID)
	row := r.db.QueryRow(ctx, `
		SELECT c.id, c.conversation_type, c.title, c.created_by, c.last_message_id, c.created_at, c.updated_at
		FROM direct_conversations dc
		JOIN conversations c ON c.id = dc.conversation_id
		WHERE dc.persona_a_id = $1 AND dc.persona_b_id = $2
	`, personaAID, personaBID)
	return scanConversation(row)
}

// CreateGroupConversation creates one group conversation for the creator and invited members.
func (r *pgRepository) CreateGroupConversation(ctx context.Context, creator *models.Persona, members []*models.Persona, dto dtos.CreateGroupConversationDTO) (*models.Conversation, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	conversation, err := insertConversation(ctx, tx, models.GroupConversationType, dto.Title, creator.ID)
	if err != nil {
		return nil, err
	}
	if _, err := insertMember(ctx, tx, conversation.ID, creator.UserID, creator.ID, models.ConversationMemberRoleAdmin); err != nil {
		return nil, err
	}
	for _, member := range members {
		if member.ID == creator.ID {
			continue
		}
		if _, err := insertMember(ctx, tx, conversation.ID, member.UserID, member.ID, models.ConversationMemberRoleMember); err != nil {
			return nil, err
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return conversation, nil
}

// FindConversationByID fetches one conversation by id.
func (r *pgRepository) FindConversationByID(ctx context.Context, conversationID uuid.UUID) (*models.Conversation, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, conversation_type, title, created_by, last_message_id, created_at, updated_at
		FROM conversations
		WHERE id = $1
	`, conversationID)
	return scanConversation(row)
}

// DeleteConversation removes one conversation and all dependent rows.
func (r *pgRepository) DeleteConversation(ctx context.Context, conversationID uuid.UUID) error {
	commandTag, err := r.db.Exec(ctx, `
		DELETE FROM conversations
		WHERE id = $1
	`, conversationID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return ErrConversationNotFound
	}
	return nil
}

// ConversationBelongsToInactiveCrew reports whether a crew-linked conversation is inactive.
func (r *pgRepository) ConversationBelongsToInactiveCrew(ctx context.Context, conversationID uuid.UUID) (bool, error) {
	var inactive bool
	err := r.db.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM event_crews
			WHERE conversation_id = $1
			  AND (status = 'ended' OR expires_at <= now())
		)
	`, conversationID).Scan(&inactive)
	return inactive, err
}

// FindPersonaConversations lists the conversations for one member profile.
func (r *pgRepository) FindPersonaConversations(ctx context.Context, userID uuid.UUID, personaID uuid.UUID, limit int, offset int) ([]*ConversationListItem, error) {
	rows, err := r.db.Query(ctx, `
		SELECT c.id, c.conversation_type, c.title, c.created_by, c.last_message_id, c.created_at, c.updated_at,
		       (
		         SELECT COUNT(*)
		         FROM messages unread
		         WHERE unread.conversation_id = c.id
		           AND unread.deleted_at IS NULL
		           AND unread.sender_user_id <> $1
		           AND (
		             cm.last_read_message_id IS NULL
		             OR unread.created_at > COALESCE((SELECT created_at FROM messages WHERE id = cm.last_read_message_id), cm.joined_at)
		           )
		       )::INT AS unread_count
		FROM conversation_members cm
		JOIN conversations c ON c.id = cm.conversation_id
		WHERE cm.user_id = $1 AND cm.persona_id = $2 AND cm.left_at IS NULL
		ORDER BY c.updated_at DESC
		LIMIT $3 OFFSET $4
	`, userID, personaID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*ConversationListItem, 0)
	for rows.Next() {
		conversation, unreadCount, err := scanConversationListItem(rows)
		if err != nil {
			return nil, err
		}
		item := &ConversationListItem{Conversation: conversation, UnreadCount: unreadCount}
		if conversation.LastMessageID != nil {
			message, err := r.FindMessageByID(ctx, *conversation.LastMessageID)
			if err != nil && !errors.Is(err, ErrMessageNotFound) {
				return nil, err
			}
			if err == nil && message.DeletedAt == nil {
				item.LastMessage = message
			}
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return items, nil
	}

	ids := make([]uuid.UUID, 0, len(items))
	itemByID := make(map[uuid.UUID]*ConversationListItem, len(items))
	for _, item := range items {
		ids = append(ids, item.Conversation.ID)
		itemByID[item.Conversation.ID] = item
	}
	members, err := r.findMembersByConversationIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	for _, member := range members {
		if item := itemByID[member.ConversationID]; item != nil {
			item.Members = append(item.Members, member)
		}
	}
	return items, nil
}
