package messaging

import (
	"context"

	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
)

// fetches active members for one conversation.
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

// FindConversationMemberUserIDs fetches active member user ids for one conversation.
func (r *pgRepository) FindConversationMemberUserIDs(ctx context.Context, conversationID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := r.db.Query(ctx, `
		SELECT DISTINCT user_id
		FROM conversation_members
		WHERE conversation_id = $1 AND left_at IS NULL
	`, conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanUUIDs(rows)
}

// FindRelatedConversationUserIDs fetches users who share active conversations with one user.
func (r *pgRepository) FindRelatedConversationUserIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := r.db.Query(ctx, `
		SELECT DISTINCT peers.user_id
		FROM conversation_members mine
		JOIN conversation_members peers ON peers.conversation_id = mine.conversation_id
		WHERE mine.user_id = $1
		  AND mine.left_at IS NULL
		  AND peers.left_at IS NULL
		  AND peers.user_id <> $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanUUIDs(rows)
}

// fetches one active conversation member by profile id.
func (r *pgRepository) FindMember(ctx context.Context, conversationID uuid.UUID, personaID uuid.UUID) (*models.ConversationMember, error) {
	row := r.db.QueryRow(ctx, `
		SELECT conversation_id, user_id, persona_id, role, last_read_message_id, joined_at, left_at
		FROM conversation_members
		WHERE conversation_id = $1 AND persona_id = $2 AND left_at IS NULL
	`, conversationID, personaID)
	return scanMember(row)
}

// inserts or reactivates members in a conversation.
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

// updates one member role in a group conversation.
func (r *pgRepository) UpdateConversationMemberRole(ctx context.Context, conversationID uuid.UUID, personaID uuid.UUID, role models.ConversationMemberRole) (*models.ConversationMember, error) {
	row := r.db.QueryRow(ctx, `
		UPDATE conversation_members
		SET role = $3
		WHERE conversation_id = $1 AND persona_id = $2 AND left_at IS NULL
		RETURNING conversation_id, user_id, persona_id, role, last_read_message_id, joined_at, left_at
	`, conversationID, personaID, role)
	return scanMember(row)
}

// marks one member as left and dissolves empty groups.
func (r *pgRepository) RemoveConversationMember(ctx context.Context, conversationID uuid.UUID, personaID uuid.UUID) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	commandTag, err := tx.Exec(ctx, `
		UPDATE conversation_members
		SET left_at = now()
		WHERE conversation_id = $1 AND persona_id = $2 AND left_at IS NULL
	`, conversationID, personaID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return ErrMembershipNotFound
	}
	if err := dissolveConversationIfEmpty(ctx, tx, conversationID); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// fetches members for a set of conversations.
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
