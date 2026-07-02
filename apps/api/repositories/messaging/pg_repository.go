package messaging

import (
	"context"
	"errors"

	"github.com/emmanuella-codes/nox/messaging/dtos"
	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgRepository struct {
	db *pgxpool.Pool
}

func newPgRepository(db *pgxpool.Pool) *pgRepository {
	return &pgRepository{db: db}
}

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
	if _, err := insertMember(ctx, tx, conversation.ID, recipient.UserID, recipient.ID, models.ConversationMemberRoleMember); err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return conversation, nil
}

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

func (r *pgRepository) FindConversationByID(ctx context.Context, conversationID uuid.UUID) (*models.Conversation, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, conversation_type, title, created_by, last_message_id, created_at, updated_at
		FROM conversations
		WHERE id = $1
	`, conversationID)
	return scanConversation(row)
}

func (r *pgRepository) FindUserConversations(ctx context.Context, userID uuid.UUID, limit int, offset int) ([]*ConversationListItem, error) {
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
		WHERE cm.user_id = $1 AND cm.left_at IS NULL
		ORDER BY c.updated_at DESC
		LIMIT $2 OFFSET $3
	`, userID, limit, offset)
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
			if err != nil {
				return nil, err
			}
			item.LastMessage = message
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

func (r *pgRepository) FindConversationMembers(ctx context.Context, conversationID uuid.UUID) ([]*models.ConversationMember, error) {
	rows, err := r.db.Query(ctx, `
		SELECT conversation_id, user_id, persona_id, role, last_read_message_id, joined_at, left_at
		FROM conversation_members
		WHERE conversation_id = $1 AND left_at IS NULL
		ORDER BY joined_at ASC
	`, conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanMembers(rows)
}

func (r *pgRepository) FindMember(ctx context.Context, conversationID uuid.UUID, personaID uuid.UUID) (*models.ConversationMember, error) {
	row := r.db.QueryRow(ctx, `
		SELECT conversation_id, user_id, persona_id, role, last_read_message_id, joined_at, left_at
		FROM conversation_members
		WHERE conversation_id = $1 AND persona_id = $2 AND left_at IS NULL
	`, conversationID, personaID)
	return scanMember(row)
}

func (r *pgRepository) AddConversationMembers(ctx context.Context, conversationID uuid.UUID, members []*models.Persona) ([]*models.ConversationMember, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	inserted := make([]*models.ConversationMember, 0, len(members))
	for _, member := range members {
		conversationMember, err := insertMember(ctx, tx, conversationID, member.UserID, member.ID, models.ConversationMemberRoleMember)
		if err != nil {
			return nil, err
		}
		inserted = append(inserted, conversationMember)
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return inserted, nil
}

func (r *pgRepository) RemoveConversationMember(ctx context.Context, conversationID uuid.UUID, personaID uuid.UUID) error {
	tag, err := r.db.Exec(ctx, `
		UPDATE conversation_members
		SET left_at = now()
		WHERE conversation_id = $1 AND persona_id = $2 AND left_at IS NULL
	`, conversationID, personaID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrMembershipNotFound
	}
	return nil
}

func (r *pgRepository) CreateMessage(ctx context.Context, conversationID uuid.UUID, senderUserID uuid.UUID, dto dtos.SendMessageDTO) (*models.Message, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	row := tx.QueryRow(ctx, `
		INSERT INTO messages (conversation_id, sender_user_id, sender_persona_id, body, message_type, media_asset_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, conversation_id, sender_user_id, sender_persona_id, body, message_type,
		          media_asset_id, created_at, edited_at, deleted_at
	`, conversationID, senderUserID, dto.SenderPersonaID, dto.Body, dto.MessageType, dto.MediaAssetID)
	message, err := scanMessage(row)
	if err != nil {
		return nil, err
	}
	if _, err := tx.Exec(ctx, `
		UPDATE conversations
		SET last_message_id = $2, updated_at = now()
		WHERE id = $1
	`, conversationID, message.ID); err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return message, nil
}

func (r *pgRepository) FindMessagesByConversationID(ctx context.Context, conversationID uuid.UUID, limit int, offset int) ([]*models.Message, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, conversation_id, sender_user_id, sender_persona_id, body, message_type,
		       media_asset_id, created_at, edited_at, deleted_at
		FROM messages
		WHERE conversation_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`, conversationID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanMessages(rows)
}

func (r *pgRepository) FindMessageByID(ctx context.Context, messageID uuid.UUID) (*models.Message, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, conversation_id, sender_user_id, sender_persona_id, body, message_type,
		       media_asset_id, created_at, edited_at, deleted_at
		FROM messages
		WHERE id = $1
	`, messageID)
	return scanMessage(row)
}

func (r *pgRepository) MarkConversationRead(ctx context.Context, conversationID uuid.UUID, personaID uuid.UUID, messageID uuid.UUID) (*models.ConversationMember, error) {
	row := r.db.QueryRow(ctx, `
		UPDATE conversation_members
		SET last_read_message_id = $3
		WHERE conversation_id = $1 AND persona_id = $2 AND left_at IS NULL
		RETURNING conversation_id, user_id, persona_id, role, last_read_message_id, joined_at, left_at
	`, conversationID, personaID, messageID)
	return scanMember(row)
}

func (r *pgRepository) SoftDeleteMessage(ctx context.Context, messageID uuid.UUID) (*models.Message, error) {
	row := r.db.QueryRow(ctx, `
		UPDATE messages
		SET body = '', media_asset_id = NULL, deleted_at = COALESCE(deleted_at, now())
		WHERE id = $1
		RETURNING id, conversation_id, sender_user_id, sender_persona_id, body, message_type,
		          media_asset_id, created_at, edited_at, deleted_at
	`, messageID)
	return scanMessage(row)
}

func (r *pgRepository) findMembersByConversationIDs(ctx context.Context, conversationIDs []uuid.UUID) ([]*models.ConversationMember, error) {
	rows, err := r.db.Query(ctx, `
		SELECT conversation_id, user_id, persona_id, role, last_read_message_id, joined_at, left_at
		FROM conversation_members
		WHERE conversation_id = ANY($1) AND left_at IS NULL
		ORDER BY joined_at ASC
	`, conversationIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanMembers(rows)
}

func orderedPersonaIDs(a uuid.UUID, b uuid.UUID) (uuid.UUID, uuid.UUID) {
	if a.String() > b.String() {
		return b, a
	}
	return a, b
}

func mapNotFound(err error, fallback error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return fallback
	}
	return err
}
