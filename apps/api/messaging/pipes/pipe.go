package pipes

import (
	"context"
	"errors"
	"strings"

	"github.com/emmanuella-codes/nox/messaging/messages"
	"github.com/emmanuella-codes/nox/models"
	media_repo "github.com/emmanuella-codes/nox/repositories/media"
	messaging_repo "github.com/emmanuella-codes/nox/repositories/messaging"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

type MessagingPipe struct {
	messagingRepo messaging_repo.MessagingRepository
	personaRepo   persona_repo.PersonaRepository
	mediaRepo     media_repo.MediaRepository
}

// NewMessagingPipe builds the messaging orchestration layer from repositories.
func NewMessagingPipe(messagingRepo messaging_repo.MessagingRepository, personaRepo persona_repo.PersonaRepository, mediaRepo media_repo.MediaRepository) *MessagingPipe {
	return &MessagingPipe{messagingRepo: messagingRepo, personaRepo: personaRepo, mediaRepo: mediaRepo}
}

// pipeInternalError maps internal messaging errors to pipe responses.
func pipeInternalError[T any](err error, operation string) *shared.PipeRes[T] {
	return shared.PipeInternalError[T](err, "messaging", operation, messages.Internal_Error)
}

type ConversationResponse struct {
	ID               string           `json:"id"`
	ConversationType string           `json:"conversation_type"`
	Title            string           `json:"title"`
	CreatedBy        string           `json:"created_by"`
	LastMessageID    *string          `json:"last_message_id,omitempty"`
	Members          []MemberResponse `json:"members"`
	LastMessage      *MessageResponse `json:"last_message,omitempty"`
	UnreadCount      int              `json:"unread_count"`
	CreatedAt        string           `json:"created_at"`
	UpdatedAt        string           `json:"updated_at"`
}

type MemberResponse struct {
	PersonaID         string                 `json:"persona_id"`
	Persona           *MemberPersonaResponse `json:"persona,omitempty"`
	Role              string                 `json:"role"`
	LastReadMessageID *string                `json:"last_read_message_id,omitempty"`
	JoinedAt          string                 `json:"joined_at"`
}

type MemberPersonaResponse struct {
	ID          string `json:"id"`
	Handle      string `json:"handle"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
}

type MessageResponse struct {
	ID              string             `json:"id"`
	ConversationID  string             `json:"conversation_id"`
	SenderPersonaID string             `json:"sender_persona_id"`
	Body            string             `json:"body"`
	MessageType     models.MessageType `json:"message_type"`
	MediaAssetID    *string            `json:"media_asset_id,omitempty"`
	Media           *models.MediaAsset `json:"media,omitempty"`
	Deleted         bool               `json:"deleted"`
	CreatedAt       string             `json:"created_at"`
	EditedAt        *string            `json:"edited_at,omitempty"`
}

// conversationResponse maps one conversation and its related state into the API response shape.
func (p *MessagingPipe) conversationResponse(ctx context.Context, conversation *models.Conversation, members []*models.ConversationMember, personas map[uuid.UUID]*models.Persona, lastMessage *models.Message, unreadCount int) ConversationResponse {
	var lastMessageID *string
	if conversation.LastMessageID != nil {
		value := conversation.LastMessageID.String()
		lastMessageID = &value
	}
	response := ConversationResponse{
		ID:               conversation.ID.String(),
		ConversationType: string(conversation.ConversationType),
		Title:            conversation.Title,
		CreatedBy:        conversation.CreatedBy.String(),
		LastMessageID:    lastMessageID,
		Members:          memberResponses(members, personas),
		UnreadCount:      unreadCount,
		CreatedAt:        conversation.CreatedAt.Format(timeFormat),
		UpdatedAt:        conversation.UpdatedAt.Format(timeFormat),
	}
	if lastMessage != nil {
		message := p.messageResponse(ctx, lastMessage)
		response.LastMessage = &message
	}
	return response
}

// memberResponses maps conversation members into the API response shape.
func memberResponses(members []*models.ConversationMember, personas map[uuid.UUID]*models.Persona) []MemberResponse {
	responses := make([]MemberResponse, 0, len(members))
	for _, member := range members {
		var lastReadMessageID *string
		if member.LastReadMessageID != nil {
			value := member.LastReadMessageID.String()
			lastReadMessageID = &value
		}
		var personaResponse *MemberPersonaResponse
		if persona := personas[member.PersonaID]; persona != nil {
			personaResponse = &MemberPersonaResponse{
				ID:          persona.ID.String(),
				Handle:      persona.Handle,
				DisplayName: persona.DisplayName,
				AvatarURL:   persona.AvatarURL,
			}
		}
		responses = append(responses, MemberResponse{
			PersonaID:         member.PersonaID.String(),
			Persona:           personaResponse,
			Role:              string(member.Role),
			LastReadMessageID: lastReadMessageID,
			JoinedAt:          member.JoinedAt.Format(timeFormat),
		})
	}
	return responses
}

// memberPersonas hydrates the public profiles referenced by conversation members.
func (p *MessagingPipe) memberPersonas(ctx context.Context, members []*models.ConversationMember) (map[uuid.UUID]*models.Persona, error) {
	personas := make(map[uuid.UUID]*models.Persona, len(members))
	for _, member := range members {
		if _, ok := personas[member.PersonaID]; ok {
			continue
		}
		persona, err := p.personaRepo.FindPersonaByID(ctx, member.PersonaID)
		if err != nil {
			return nil, err
		}
		personas[member.PersonaID] = persona
	}
	return personas, nil
}

// messageResponse maps one message into the API response shape.
func (p *MessagingPipe) messageResponse(ctx context.Context, message *models.Message) MessageResponse {
	var mediaAssetID *string
	var media *models.MediaAsset
	if message.MediaAssetID != nil {
		value := message.MediaAssetID.String()
		mediaAssetID = &value
		if p.mediaRepo != nil && message.DeletedAt == nil {
			asset, err := p.mediaRepo.FindMediaAssetByID(ctx, *message.MediaAssetID)
			if err == nil {
				media = asset
			}
		}
	}
	var editedAt *string
	if message.EditedAt != nil {
		value := message.EditedAt.Format(timeFormat)
		editedAt = &value
	}
	body := message.Body
	if message.DeletedAt != nil {
		body = ""
	}
	return MessageResponse{
		ID:              message.ID.String(),
		ConversationID:  message.ConversationID.String(),
		SenderPersonaID: message.SenderPersonaID.String(),
		Body:            body,
		MessageType:     message.MessageType,
		MediaAssetID:    mediaAssetID,
		Media:           media,
		Deleted:         message.DeletedAt != nil,
		CreatedAt:       message.CreatedAt.Format(timeFormat),
		EditedAt:        editedAt,
	}
}

// messageResponses maps a slice of messages into API response shape.
func (p *MessagingPipe) messageResponses(ctx context.Context, messageModels []*models.Message) []MessageResponse {
	responses := make([]MessageResponse, 0, len(messageModels))
	for _, message := range messageModels {
		responses = append(responses, p.messageResponse(ctx, message))
	}
	return responses
}

// profilePersona validates that a messaging participant is a real public profile.
func (p *MessagingPipe) profilePersona(ctx context.Context, userID uuid.UUID, personaID uuid.UUID, requireOwner bool) (*models.Persona, shared.PipeMessage) {
	persona, err := p.personaRepo.FindPersonaByID(ctx, personaID)
	if err != nil {
		if errors.Is(err, persona_repo.ErrPersonaNotFound) {
			return nil, messages.Persona_Not_Found
		}
		return nil, messages.Internal_Error
	}
	if requireOwner && persona.UserID != userID {
		return nil, messages.Forbidden
	}
	return persona, ""
}

// requireMember validates that the current user owns the supplied profile and belongs to the conversation.
func (p *MessagingPipe) requireMember(ctx context.Context, userID uuid.UUID, conversationID uuid.UUID, personaID uuid.UUID) (*models.ConversationMember, shared.PipeMessage) {
	persona, message := p.profilePersona(ctx, userID, personaID, true)
	if message != "" {
		return nil, message
	}
	member, err := p.messagingRepo.FindMember(ctx, conversationID, persona.ID)
	if err != nil {
		if errors.Is(err, messaging_repo.ErrMembershipNotFound) {
			return nil, messages.Forbidden
		}
		return nil, messages.Internal_Error
	}
	if member.UserID != userID {
		return nil, messages.Forbidden
	}
	return member, ""
}

// validMessageType validates the supported message payload types.
func validMessageType(messageType models.MessageType) bool {
	return messageType == models.TextMessageType || messageType == models.ImageMessageType || messageType == models.VideoMessageType
}

// normalizeMessageBody trims whitespace from message bodies.
func normalizeMessageBody(body string) string {
	return strings.TrimSpace(body)
}

const timeFormat = "2006-01-02T15:04:05.999999999Z07:00"
