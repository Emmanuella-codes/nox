package pipes

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/emmanuella-codes/nox/messaging/dtos"
	"github.com/emmanuella-codes/nox/messaging/messages"
	"github.com/emmanuella-codes/nox/models"
	messaging_repo "github.com/emmanuella-codes/nox/repositories/messaging"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// SendMessagePipe validates membership, payload shape, and attachments before persisting one message.
func (p *MessagingPipe) SendMessagePipe(ctx context.Context, userID uuid.UUID, conversationID uuid.UUID, dto dtos.SendMessageDTO) *shared.PipeRes[MessageResponse] {
	dto.Body = normalizeMessageBody(dto.Body)
	dto.IdempotencyKey = strings.TrimSpace(dto.IdempotencyKey)
	attachmentIDs := normalizeAttachmentIDs(dto.MediaAssetID, dto.MediaAssetIDs)
	dto.MediaAssetID = nil
	dto.MediaAssetIDs = attachmentIDs
	if dto.Body == "" && len(attachmentIDs) == 0 {
		return shared.PipeError[MessageResponse](messages.Invalid_Payload)
	}
	if len(attachmentIDs) > 5 {
		return shared.PipeError[MessageResponse](messages.Invalid_Payload)
	}
	member, message := p.requireMember(ctx, userID, conversationID, dto.SenderPersonaID)
	if message != "" {
		return shared.PipeError[MessageResponse](message)
	}
	inactiveCrew, err := p.messagingRepo.ConversationBelongsToInactiveCrew(ctx, conversationID)
	if err != nil {
		return pipeInternalError[MessageResponse](err, "messaging.check_crew_conversation")
	}
	if inactiveCrew {
		return shared.PipeError[MessageResponse](messages.Forbidden)
	}
	conversation, err := p.messagingRepo.FindConversationByID(ctx, conversationID)
	if err != nil {
		return pipeInternalError[MessageResponse](err, "messaging.find_conversation")
	}
	attachmentType := models.TextMessageType
	if len(attachmentIDs) > 0 {
		var mediaMessage shared.PipeMessage
		attachmentType, mediaMessage = p.validateMessageMedia(ctx, userID, member.PersonaID, attachmentIDs)
		if mediaMessage != "" {
			return shared.PipeError[MessageResponse](mediaMessage)
		}
	}
	dto.MessageType = normalizedMessageType(dto.MessageType, dto.Body, len(attachmentIDs) > 0, attachmentType)
	if !validMessageType(dto.MessageType) {
		return shared.PipeError[MessageResponse](messages.Invalid_Payload)
	}
	created, createdNew, err := p.messagingRepo.CreateMessage(ctx, conversationID, userID, dto)
	if err != nil {
		return pipeInternalError[MessageResponse](err, "messaging.send_message")
	}
	response := p.messageResponse(ctx, created)
	if createdNew {
		p.createMessageNotifications(ctx, conversation, created)
		p.publishMessageEvent(ctx, userID, "message.created", created)
		p.publishTypingEvent(ctx, conversationID, dto.SenderPersonaID, false)
	}
	return shared.PipeSuccess(messages.Message_Sent, &response)
}

// ListMessagesPipe lists visible messages for one conversation member.
func (p *MessagingPipe) ListMessagesPipe(ctx context.Context, userID uuid.UUID, conversationID uuid.UUID, personaID uuid.UUID, limit int, offset int) *shared.PipeRes[[]MessageResponse] {
	if _, message := p.requireMember(ctx, userID, conversationID, personaID); message != "" {
		return shared.PipeError[[]MessageResponse](message)
	}
	limit = normalizeLimit(limit, 30, 100)
	if offset < 0 {
		offset = 0
	}
	messageModels, err := p.messagingRepo.FindMessagesByConversationID(ctx, conversationID, limit, offset)
	if err != nil {
		return pipeInternalError[[]MessageResponse](err, "messaging.list_messages")
	}
	responses := p.messageResponses(ctx, messageModels)
	return shared.PipeSuccess(messages.Messages_Listed, &responses)
}

// MarkReadPipe advances the current member read cursor to a visible message.
func (p *MessagingPipe) MarkReadPipe(ctx context.Context, userID uuid.UUID, conversationID uuid.UUID, dto dtos.MarkReadDTO) *shared.PipeRes[MemberResponse] {
	memberBefore, message := p.requireMember(ctx, userID, conversationID, dto.PersonaID)
	if message != "" {
		return shared.PipeError[MemberResponse](message)
	}
	messageModel, err := p.messagingRepo.FindMessageByID(ctx, dto.MessageID)
	if err != nil {
		if errors.Is(err, messaging_repo.ErrMessageNotFound) {
			return shared.PipeError[MemberResponse](messages.Message_Not_Found)
		}
		return pipeInternalError[MemberResponse](err, "messaging.read_message")
	}
	if messageModel.ConversationID != conversationID {
		return shared.PipeError[MemberResponse](messages.Message_Not_Found)
	}
	if messageModel.DeletedAt != nil {
		return shared.PipeError[MemberResponse](messages.Message_Not_Found)
	}
	member, err := p.messagingRepo.MarkConversationRead(ctx, conversationID, dto.PersonaID, dto.MessageID)
	if err != nil {
		return pipeInternalError[MemberResponse](err, "messaging.mark_read")
	}
	personas, err := p.memberPersonas(ctx, []*models.ConversationMember{member})
	if err != nil {
		return pipeInternalError[MemberResponse](err, "messaging.mark_read_persona")
	}
	response := memberResponses([]*models.ConversationMember{member}, personas)[0]
	if member.LastReadMessageID != nil && (memberBefore.LastReadMessageID == nil || *memberBefore.LastReadMessageID != *member.LastReadMessageID) {
		p.markConversationNotificationsRead(ctx, userID, dto.PersonaID, conversationID, dto.MessageID)
		p.publishConversationEvent(ctx, conversationID, userID, "conversation.read", &messageModel.ID, response)
	}
	return shared.PipeSuccess(messages.Conversation_Read, &response)
}

// EditMessagePipe updates the body of a sender-owned message within the mutation window.
func (p *MessagingPipe) EditMessagePipe(ctx context.Context, userID uuid.UUID, messageID uuid.UUID, dto dtos.EditMessageDTO) *shared.PipeRes[MessageResponse] {
	dto.Body = normalizeMessageBody(dto.Body)
	if dto.Body == "" {
		return shared.PipeError[MessageResponse](messages.Invalid_Payload)
	}
	messageModel, err := p.messagingRepo.FindMessageByID(ctx, messageID)
	if err != nil {
		if errors.Is(err, messaging_repo.ErrMessageNotFound) {
			return shared.PipeError[MessageResponse](messages.Message_Not_Found)
		}
		return pipeInternalError[MessageResponse](err, "messaging.find_edit_message")
	}
	if messageModel.SenderUserID != userID || messageModel.DeletedAt != nil || !canMutateMessage(messageModel, time.Now()) {
		return shared.PipeError[MessageResponse](messages.Forbidden)
	}
	updated, err := p.messagingRepo.UpdateMessageBody(ctx, messageID, dto.Body)
	if err != nil {
		if errors.Is(err, messaging_repo.ErrMessageNotFound) {
			return shared.PipeError[MessageResponse](messages.Message_Not_Found)
		}
		return pipeInternalError[MessageResponse](err, "messaging.edit_message")
	}
	response := p.messageResponse(ctx, updated)
	p.publishMessageEvent(ctx, userID, "message.updated", updated)
	return shared.PipeSuccess(messages.Message_Updated, &response)
}

// DeleteMessagePipe hides a sender-owned message for all members within the mutation window.
func (p *MessagingPipe) DeleteMessagePipe(ctx context.Context, userID uuid.UUID, messageID uuid.UUID) *shared.PipeRes[MessageResponse] {
	messageModel, err := p.messagingRepo.FindMessageByID(ctx, messageID)
	if err != nil {
		if errors.Is(err, messaging_repo.ErrMessageNotFound) {
			return shared.PipeError[MessageResponse](messages.Message_Not_Found)
		}
		return pipeInternalError[MessageResponse](err, "messaging.find_delete_message")
	}
	if messageModel.SenderUserID != userID || messageModel.DeletedAt != nil || !canMutateMessage(messageModel, time.Now()) {
		return shared.PipeError[MessageResponse](messages.Forbidden)
	}
	deleted, err := p.messagingRepo.SoftDeleteMessage(ctx, messageID)
	if err != nil {
		return pipeInternalError[MessageResponse](err, "messaging.delete_message")
	}
	response := p.messageResponse(ctx, deleted)
	p.deleteMessageNotifications(ctx, messageID)
	p.publishMessageEvent(ctx, userID, "message.deleted", deleted)
	return shared.PipeSuccess(messages.Message_Deleted, &response)
}

// validateMessageMedia verifies that all supplied attachments belong to the sender and are chat-safe.
func (p *MessagingPipe) validateMessageMedia(ctx context.Context, userID uuid.UUID, personaID uuid.UUID, attachmentIDs []uuid.UUID) (models.MessageType, shared.PipeMessage) {
	if p.mediaRepo == nil || len(attachmentIDs) == 0 {
		return "", messages.Invalid_Payload
	}
	messageType := models.TextMessageType
	for _, attachmentID := range attachmentIDs {
		asset, err := p.mediaRepo.FindMediaAssetByID(ctx, attachmentID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return "", messages.Invalid_Payload
			}
			return "", messages.Internal_Error
		}
		if asset.OwnerUserID != userID || asset.OwnerPersonaID != personaID || asset.ProcessingStatus != models.ReadyMediaStatus {
			return "", messages.Forbidden
		}
		switch asset.MediaKind {
		case models.ImageMediaKind:
			if messageType == models.TextMessageType {
				messageType = models.ImageMessageType
			}
		case models.VideoMediaKind:
			if messageType == models.TextMessageType {
				messageType = models.VideoMessageType
			}
		case models.AudioMediaKind:
			if messageType == models.TextMessageType {
				messageType = models.AudioMessageType
			}
		default:
			return "", messages.Invalid_Payload
		}
	}
	return messageType, ""
}

// normalizedMessageType derives a stable message type from the payload and attachments.
func normalizedMessageType(messageType models.MessageType, body string, hasAttachments bool, attachmentType models.MessageType) models.MessageType {
	if messageType != "" {
		return messageType
	}
	if body != "" || !hasAttachments {
		return models.TextMessageType
	}
	return attachmentType
}

// UpdateTypingPipe broadcasts typing state for one active conversation member.
func (p *MessagingPipe) UpdateTypingPipe(ctx context.Context, userID uuid.UUID, conversationID uuid.UUID, dto dtos.TypingDTO) *shared.PipeRes[any] {
	if _, message := p.requireMember(ctx, userID, conversationID, dto.PersonaID); message != "" {
		return shared.PipeError[any](message)
	}
	p.publishTypingEvent(ctx, conversationID, dto.PersonaID, dto.Active)
	return shared.PipeSuccess[any](messages.Typing_Updated, nil)
}
