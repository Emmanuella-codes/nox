package pipes

import (
	"context"
	"errors"

	"github.com/emmanuella-codes/nox/messaging/dtos"
	"github.com/emmanuella-codes/nox/messaging/messages"
	"github.com/emmanuella-codes/nox/models"
	messaging_repo "github.com/emmanuella-codes/nox/repositories/messaging"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (p *MessagingPipe) SendMessagePipe(ctx context.Context, userID uuid.UUID, conversationID uuid.UUID, dto dtos.SendMessageDTO) *shared.PipeRes[MessageResponse] {
	dto.Body = normalizeMessageBody(dto.Body)
	if dto.MessageType == "" {
		dto.MessageType = models.TextMessageType
	}
	if !validMessageType(dto.MessageType) || (dto.Body == "" && dto.MediaAssetID == nil) {
		return shared.PipeError[MessageResponse](messages.Invalid_Payload)
	}
	if dto.MediaAssetID == nil && dto.MessageType != models.TextMessageType {
		return shared.PipeError[MessageResponse](messages.Invalid_Payload)
	}
	if dto.MediaAssetID != nil && dto.MessageType == models.TextMessageType {
		return shared.PipeError[MessageResponse](messages.Invalid_Payload)
	}
	member, message := p.requireMember(ctx, userID, conversationID, dto.SenderPersonaID)
	if message != "" {
		return shared.PipeError[MessageResponse](message)
	}
	if dto.MediaAssetID != nil {
		if message := p.validateMessageMedia(ctx, userID, member.PersonaID, dto); message != "" {
			return shared.PipeError[MessageResponse](message)
		}
	}
	created, err := p.messagingRepo.CreateMessage(ctx, conversationID, userID, dto)
	if err != nil {
		return pipeInternalError[MessageResponse](err, "messaging.send_message")
	}
	response := p.messageResponse(ctx, created)
	return shared.PipeSuccess(messages.Message_Sent, &response)
}

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

func (p *MessagingPipe) MarkReadPipe(ctx context.Context, userID uuid.UUID, conversationID uuid.UUID, dto dtos.MarkReadDTO) *shared.PipeRes[MemberResponse] {
	if _, message := p.requireMember(ctx, userID, conversationID, dto.PersonaID); message != "" {
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
	member, err := p.messagingRepo.MarkConversationRead(ctx, conversationID, dto.PersonaID, dto.MessageID)
	if err != nil {
		return pipeInternalError[MemberResponse](err, "messaging.mark_read")
	}
	personas, err := p.memberPersonas(ctx, []*models.ConversationMember{member})
	if err != nil {
		return pipeInternalError[MemberResponse](err, "messaging.mark_read_persona")
	}
	response := memberResponses([]*models.ConversationMember{member}, personas)[0]
	return shared.PipeSuccess(messages.Conversation_Read, &response)
}

func (p *MessagingPipe) DeleteMessagePipe(ctx context.Context, userID uuid.UUID, messageID uuid.UUID) *shared.PipeRes[MessageResponse] {
	messageModel, err := p.messagingRepo.FindMessageByID(ctx, messageID)
	if err != nil {
		if errors.Is(err, messaging_repo.ErrMessageNotFound) {
			return shared.PipeError[MessageResponse](messages.Message_Not_Found)
		}
		return pipeInternalError[MessageResponse](err, "messaging.find_delete_message")
	}
	if messageModel.SenderUserID != userID {
		return shared.PipeError[MessageResponse](messages.Forbidden)
	}
	deleted, err := p.messagingRepo.SoftDeleteMessage(ctx, messageID)
	if err != nil {
		return pipeInternalError[MessageResponse](err, "messaging.delete_message")
	}
	response := p.messageResponse(ctx, deleted)
	return shared.PipeSuccess(messages.Message_Deleted, &response)
}

func (p *MessagingPipe) validateMessageMedia(ctx context.Context, userID uuid.UUID, personaID uuid.UUID, dto dtos.SendMessageDTO) shared.PipeMessage {
	if p.mediaRepo == nil || dto.MediaAssetID == nil {
		return messages.Invalid_Payload
	}
	asset, err := p.mediaRepo.FindMediaAssetByID(ctx, *dto.MediaAssetID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return messages.Invalid_Payload
		}
		return messages.Internal_Error
	}
	if asset.OwnerUserID != userID || asset.OwnerPersonaID != personaID || asset.ProcessingStatus != models.ReadyMediaStatus {
		return messages.Forbidden
	}
	if dto.MessageType == models.ImageMessageType && asset.MediaKind != models.ImageMediaKind {
		return messages.Invalid_Payload
	}
	if dto.MessageType == models.VideoMessageType && asset.MediaKind != models.VideoMediaKind {
		return messages.Invalid_Payload
	}
	return ""
}
