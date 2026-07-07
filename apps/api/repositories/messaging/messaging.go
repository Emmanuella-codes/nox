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
	CreateDirectConversation(ctx context.Context, creator *models.Persona, recipient *models.Persona) (*models.Conversation, error)
	FindDirectConversationBetweenPersonas(ctx context.Context, personaAID uuid.UUID, personaBID uuid.UUID) (*models.Conversation, error)
	CreateGroupConversation(ctx context.Context, creator *models.Persona, members []*models.Persona, dto dtos.CreateGroupConversationDTO) (*models.Conversation, error)
	FindConversationByID(ctx context.Context, conversationID uuid.UUID) (*models.Conversation, error)
	ConversationBelongsToInactiveCrew(ctx context.Context, conversationID uuid.UUID) (bool, error)
	FindPersonaConversations(ctx context.Context, userID uuid.UUID, personaID uuid.UUID, limit int, offset int) ([]*ConversationListItem, error)
	FindConversationMembers(ctx context.Context, conversationID uuid.UUID) ([]*models.ConversationMember, error)
	FindMember(ctx context.Context, conversationID uuid.UUID, personaID uuid.UUID) (*models.ConversationMember, error)
	AddConversationMembers(ctx context.Context, conversationID uuid.UUID, members []*models.Persona) ([]*models.ConversationMember, error)
	RemoveConversationMember(ctx context.Context, conversationID uuid.UUID, personaID uuid.UUID) error
	CreateMessage(ctx context.Context, conversationID uuid.UUID, senderUserID uuid.UUID, dto dtos.SendMessageDTO) (*models.Message, error)
	FindMessagesByConversationID(ctx context.Context, conversationID uuid.UUID, limit int, offset int) ([]*models.Message, error)
	FindMessageByID(ctx context.Context, messageID uuid.UUID) (*models.Message, error)
	MarkConversationRead(ctx context.Context, conversationID uuid.UUID, personaID uuid.UUID, messageID uuid.UUID) (*models.ConversationMember, error)
	SoftDeleteMessage(ctx context.Context, messageID uuid.UUID) (*models.Message, error)
}

func NewMessagingRepository(db *pgxpool.Pool) MessagingRepository {
	return newPgRepository(db)
}
