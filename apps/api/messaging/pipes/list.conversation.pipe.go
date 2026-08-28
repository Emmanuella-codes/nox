package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/messaging/messages"
	"github.com/emmanuella-codes/nox/models"
	messaging_repo "github.com/emmanuella-codes/nox/repositories/messaging"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

// ListConversationsPipe lists the conversations available to one profile.
func (p *MessagingPipe) ListConversationsPipe(ctx context.Context, userID uuid.UUID, personaID uuid.UUID, limit int, offset int) *shared.PipeRes[[]ConversationResponse] {
	if _, message := p.profilePersona(ctx, userID, personaID, true); message != "" {
		return shared.PipeError[[]ConversationResponse](message)
	}
	limit = normalizeLimit(limit, 20, 50)
	if offset < 0 {
		offset = 0
	}
	items, err := p.messagingRepo.FindPersonaConversations(ctx, userID, personaID, limit, offset)
	if err != nil {
		return pipeInternalError[[]ConversationResponse](err, "messaging.list_conversations")
	}
	responses := make([]ConversationResponse, 0, len(items))
	for _, item := range items {
		personas, err := p.memberPersonas(ctx, item.Members)
		if err != nil {
			return pipeInternalError[[]ConversationResponse](err, "messaging.list_member_personas")
		}
		responses = append(responses, p.conversationResponse(ctx, item.Conversation, item.Members, personas, item.LastMessage, item.UnreadCount))
	}
	return shared.PipeSuccess(messages.Conversations_Listed, &responses)
}

// GetConversationPipe fetches one conversation if the current user is a member.
func (p *MessagingPipe) GetConversationPipe(ctx context.Context, userID uuid.UUID, conversationID uuid.UUID) *shared.PipeRes[ConversationResponse] {
	conversation, err := p.messagingRepo.FindConversationByID(ctx, conversationID)
	if err != nil {
		if err == messaging_repo.ErrConversationNotFound {
			return shared.PipeError[ConversationResponse](messages.Conversation_Not_Found)
		}
		return pipeInternalError[ConversationResponse](err, "messaging.get_conversation")
	}
	members, err := p.messagingRepo.FindConversationMembers(ctx, conversationID)
	if err != nil {
		return pipeInternalError[ConversationResponse](err, "messaging.get_members")
	}
	if !userInConversation(userID, members) {
		return shared.PipeError[ConversationResponse](messages.Forbidden)
	}
	personas, err := p.memberPersonas(ctx, members)
	if err != nil {
		return pipeInternalError[ConversationResponse](err, "messaging.get_member_personas")
	}
	var lastMessage *models.Message
	if conversation.LastMessageID != nil {
		lastMessage, err = p.messagingRepo.FindMessageByID(ctx, *conversation.LastMessageID)
		if err != nil {
			return pipeInternalError[ConversationResponse](err, "messaging.get_last_message")
		}
		if lastMessage.DeletedAt != nil {
			lastMessage = nil
		}
	}
	response := p.conversationResponse(ctx, conversation, members, personas, lastMessage, 0)
	return shared.PipeSuccess(messages.Conversation_Fetched, &response)
}

// userInConversation checks whether the current user belongs to the conversation.
func userInConversation(userID uuid.UUID, members []*models.ConversationMember) bool {
	for _, member := range members {
		if member.UserID == userID {
			return true
		}
	}
	return false
}

// normalizeLimit bounds list sizes to the accepted window.
func normalizeLimit(limit int, fallback int, max int) int {
	if limit <= 0 {
		return fallback
	}
	if limit > max {
		return max
	}
	return limit
}
