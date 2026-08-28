package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/models"
	"github.com/emmanuella-codes/nox/shared/realtime"
	"github.com/google/uuid"
)

type messageEventPayload struct {
	ConversationID string          `json:"conversation_id"`
	Message        MessageResponse `json:"message"`
}

type memberEventPayload struct {
	ConversationID string         `json:"conversation_id"`
	Member         MemberResponse `json:"member"`
}

type typingEventPayload struct {
	ConversationID string `json:"conversation_id"`
	PersonaID      string `json:"persona_id"`
	Active         bool   `json:"active"`
}

// publishConversationEvent sends one realtime event to all active users in a conversation.
func (p *MessagingPipe) publishConversationEvent(ctx context.Context, conversationID uuid.UUID, eventType string, data any) {
	if p.realtimeHub == nil || p.messagingRepo == nil {
		return
	}
	userIDs, err := p.messagingRepo.FindConversationMemberUserIDs(ctx, conversationID)
	if err != nil || len(userIDs) == 0 {
		return
	}
	_ = p.realtimeHub.PublishUsers(userIDs, realtime.Event{Type: eventType, Data: data})
}

// publishMessageEvent sends one realtime message event to a conversation audience.
func (p *MessagingPipe) publishMessageEvent(ctx context.Context, eventType string, message *models.Message) {
	p.publishConversationEvent(ctx, message.ConversationID, eventType, messageEventPayload{
		ConversationID: message.ConversationID.String(),
		Message:        p.messageResponse(ctx, message),
	})
}

// publishMemberEvent sends one realtime member event to a conversation audience.
func (p *MessagingPipe) publishMemberEvent(ctx context.Context, eventType string, member *models.ConversationMember) {
	personas, err := p.memberPersonas(ctx, []*models.ConversationMember{member})
	if err != nil {
		return
	}
	p.publishConversationEvent(ctx, member.ConversationID, eventType, memberEventPayload{
		ConversationID: member.ConversationID.String(),
		Member:         memberResponses([]*models.ConversationMember{member}, personas)[0],
	})
}

// publishTypingEvent sends one realtime typing event to a conversation audience.
func (p *MessagingPipe) publishTypingEvent(ctx context.Context, conversationID uuid.UUID, personaID uuid.UUID, active bool) {
	p.publishConversationEvent(ctx, conversationID, "typing.updated", typingEventPayload{
		ConversationID: conversationID.String(),
		PersonaID:      personaID.String(),
		Active:         active,
	})
}
