package messaging

import (
	"context"
	"errors"

	"github.com/emmanuella-codes/nox/messaging/dtos"
	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrConversationNotFound = errors.New("conversation not found")
	ErrMessageNotFound      = errors.New("message not found")
	ErrMembershipNotFound   = errors.New("conversation membership not found")
)

type ConversationListItem struct {
	Conversation *models.Conversation
	Members      []*models.ConversationMember
	LastMessage  *models.Message
	UnreadCount  int
}

type MessagingRepository interface {
	// CreateDirectConversation creates one direct conversation between two profiles.
	CreateDirectConversation(ctx context.Context, creator *models.Persona, recipient *models.Persona) (*models.Conversation, error)
	// FindDirectConversationBetweenPersonas finds the direct conversation for two profiles.
	FindDirectConversationBetweenPersonas(ctx context.Context, personaAID uuid.UUID, personaBID uuid.UUID) (*models.Conversation, error)
	// CreateGroupConversation creates one group conversation for the creator and invited members.
	CreateGroupConversation(ctx context.Context, creator *models.Persona, members []*models.Persona, dto dtos.CreateGroupConversationDTO) (*models.Conversation, error)
	// FindConversationByID fetches one conversation by id.
	FindConversationByID(ctx context.Context, conversationID uuid.UUID) (*models.Conversation, error)
	// DeleteConversation removes one conversation and all dependent rows.
	DeleteConversation(ctx context.Context, conversationID uuid.UUID) error
	// ConversationBelongsToInactiveCrew reports whether a crew-linked conversation is inactive.
	ConversationBelongsToInactiveCrew(ctx context.Context, conversationID uuid.UUID) (bool, error)
	// FindPersonaConversations lists the conversations for one member profile.
	FindPersonaConversations(ctx context.Context, userID uuid.UUID, personaID uuid.UUID, limit int, offset int) ([]*ConversationListItem, error)
	// FindConversationMembers fetches active members for one conversation.
	FindConversationMembers(ctx context.Context, conversationID uuid.UUID) ([]*models.ConversationMember, error)
	// FindMember fetches one active conversation member by profile id.
	FindMember(ctx context.Context, conversationID uuid.UUID, personaID uuid.UUID) (*models.ConversationMember, error)
	// AddConversationMembers inserts or reactivates members in a conversation.
	AddConversationMembers(ctx context.Context, conversationID uuid.UUID, members []*models.Persona) ([]*models.ConversationMember, error)
	// UpdateConversationMemberRole updates one member role in a group conversation.
	UpdateConversationMemberRole(ctx context.Context, conversationID uuid.UUID, personaID uuid.UUID, role models.ConversationMemberRole) (*models.ConversationMember, error)
	// RemoveConversationMember marks one member as left and dissolves empty groups.
	RemoveConversationMember(ctx context.Context, conversationID uuid.UUID, personaID uuid.UUID) error
	// CreateMessage persists one message and its attachments.
	CreateMessage(ctx context.Context, conversationID uuid.UUID, senderUserID uuid.UUID, dto dtos.SendMessageDTO) (*models.Message, error)
	// UpdateMessageBody edits one message body and marks it as edited.
	UpdateMessageBody(ctx context.Context, messageID uuid.UUID, body string) (*models.Message, error)
	// FindMessagesByConversationID lists visible messages in one conversation.
	FindMessagesByConversationID(ctx context.Context, conversationID uuid.UUID, limit int, offset int) ([]*models.Message, error)
	// FindMessageByID fetches one message by id, including deleted messages.
	FindMessageByID(ctx context.Context, messageID uuid.UUID) (*models.Message, error)
	// FindMessageAttachmentsByMessageIDs fetches attachments for a set of message ids.
	FindMessageAttachmentsByMessageIDs(ctx context.Context, messageIDs []uuid.UUID) (map[uuid.UUID][]*models.MediaAsset, error)
	// MarkConversationRead advances the member read cursor for one conversation.
	MarkConversationRead(ctx context.Context, conversationID uuid.UUID, personaID uuid.UUID, messageID uuid.UUID) (*models.ConversationMember, error)
	// SoftDeleteMessage hides one message from all members and refreshes conversation state.
	SoftDeleteMessage(ctx context.Context, messageID uuid.UUID) (*models.Message, error)
}

// NewMessagingRepository builds the messaging repository from a database pool.
func NewMessagingRepository(db *pgxpool.Pool) MessagingRepository {
	return newPgRepository(db)
}
