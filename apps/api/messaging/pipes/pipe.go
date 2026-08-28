package pipes

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/emmanuella-codes/nox/messaging/messages"
	"github.com/emmanuella-codes/nox/models"
	follow_repo "github.com/emmanuella-codes/nox/repositories/follow"
	media_repo "github.com/emmanuella-codes/nox/repositories/media"
	messaging_repo "github.com/emmanuella-codes/nox/repositories/messaging"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/emmanuella-codes/nox/shared/realtime"
	"github.com/google/uuid"
)

type MessagingPipe struct {
	messagingRepo messaging_repo.MessagingRepository
	personaRepo   persona_repo.PersonaRepository
	mediaRepo     media_repo.MediaRepository
	followRepo    follow_repo.FollowRepository
	realtimeHub   *realtime.Hub
}

// NewMessagingPipe builds the messaging orchestration layer from repositories.
func NewMessagingPipe(messagingRepo messaging_repo.MessagingRepository, personaRepo persona_repo.PersonaRepository, mediaRepo media_repo.MediaRepository, followRepo follow_repo.FollowRepository, deps ...any) *MessagingPipe {
	pipe := &MessagingPipe{messagingRepo: messagingRepo, personaRepo: personaRepo, mediaRepo: mediaRepo, followRepo: followRepo}
	for _, dep := range deps {
		if hub, ok := dep.(*realtime.Hub); ok {
			pipe.realtimeHub = hub
		}
	}
	return pipe
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
	ID              string               `json:"id"`
	ConversationID  string               `json:"conversation_id"`
	SenderPersonaID string               `json:"sender_persona_id"`
	Body            string               `json:"body"`
	MessageType     models.MessageType   `json:"message_type"`
	Attachments     []*models.MediaAsset `json:"attachments"`
	MediaAssetID    *string              `json:"media_asset_id,omitempty"`
	Media           *models.MediaAsset   `json:"media,omitempty"`
	Deleted         bool                 `json:"deleted"`
	Edited          bool                 `json:"edited"`
	CreatedAt       string               `json:"created_at"`
	EditedAt        *string              `json:"edited_at,omitempty"`
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
	attachments := p.messageAttachments(ctx, []uuid.UUID{message.ID})
	return messageResponseWithAttachments(message, attachments[message.ID])
}

// messageResponses maps a slice of messages into API response shape.
func (p *MessagingPipe) messageResponses(ctx context.Context, messageModels []*models.Message) []MessageResponse {
	messageIDs := make([]uuid.UUID, 0, len(messageModels))
	for _, message := range messageModels {
		messageIDs = append(messageIDs, message.ID)
	}
	attachments := p.messageAttachments(ctx, messageIDs)
	responses := make([]MessageResponse, 0, len(messageModels))
	for _, message := range messageModels {
		responses = append(responses, messageResponseWithAttachments(message, attachments[message.ID]))
	}
	return responses
}

// messageAttachments loads attachments for the supplied message ids.
func (p *MessagingPipe) messageAttachments(ctx context.Context, messageIDs []uuid.UUID) map[uuid.UUID][]*models.MediaAsset {
	attachments := make(map[uuid.UUID][]*models.MediaAsset, len(messageIDs))
	if p.messagingRepo == nil || len(messageIDs) == 0 {
		return attachments
	}
	loaded, err := p.messagingRepo.FindMessageAttachmentsByMessageIDs(ctx, messageIDs)
	if err != nil {
		return attachments
	}
	return loaded
}

// messageResponseWithAttachments maps one message and its attachments into the API response shape.
func messageResponseWithAttachments(message *models.Message, attachments []*models.MediaAsset) MessageResponse {
	var mediaAssetID *string
	var media *models.MediaAsset
	if len(attachments) > 0 {
		value := attachments[0].ID.String()
		mediaAssetID = &value
		media = attachments[0]
	}
	var editedAt *string
	if message.EditedAt != nil {
		value := message.EditedAt.Format(timeFormat)
		editedAt = &value
	}
	return MessageResponse{
		ID:              message.ID.String(),
		ConversationID:  message.ConversationID.String(),
		SenderPersonaID: message.SenderPersonaID.String(),
		Body:            message.Body,
		MessageType:     message.MessageType,
		Attachments:     attachmentsOrEmpty(attachments),
		MediaAssetID:    mediaAssetID,
		Media:           media,
		Deleted:         message.DeletedAt != nil,
		Edited:          message.EditedAt != nil,
		CreatedAt:       message.CreatedAt.Format(timeFormat),
		EditedAt:        editedAt,
	}
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
	return messageType == models.TextMessageType || messageType == models.ImageMessageType || messageType == models.VideoMessageType || messageType == models.AudioMessageType
}

// normalizeMessageBody trims whitespace from message bodies.
func normalizeMessageBody(body string) string {
	return strings.TrimSpace(body)
}

// normalizeAttachmentIDs merges legacy and multi-attachment payload inputs.
func normalizeAttachmentIDs(legacy *uuid.UUID, attachmentIDs []uuid.UUID) []uuid.UUID {
	normalized := make([]uuid.UUID, 0, len(attachmentIDs)+1)
	seen := map[uuid.UUID]bool{}
	if legacy != nil && *legacy != uuid.Nil {
		seen[*legacy] = true
		normalized = append(normalized, *legacy)
	}
	for _, attachmentID := range attachmentIDs {
		if attachmentID == uuid.Nil || seen[attachmentID] {
			continue
		}
		seen[attachmentID] = true
		normalized = append(normalized, attachmentID)
	}
	return normalized
}

// attachmentsOrEmpty normalizes nil attachment slices into empty response arrays.
func attachmentsOrEmpty(attachments []*models.MediaAsset) []*models.MediaAsset {
	if attachments == nil {
		return []*models.MediaAsset{}
	}
	return attachments
}

// canMutateMessage reports whether a message is still inside the one-hour mutation window.
func canMutateMessage(message *models.Message, now time.Time) bool {
	return now.Before(message.CreatedAt.Add(messageMutationWindow))
}

const timeFormat = "2006-01-02T15:04:05.999999999Z07:00"
const messageMutationWindow = time.Hour
