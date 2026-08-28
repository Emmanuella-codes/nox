package pipes

import (
	"context"
	"encoding/json"
	"strconv"

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

// publishConversationEvent stores and broadcasts one replayable conversation event.
func (p *MessagingPipe) publishConversationEvent(ctx context.Context, conversationID uuid.UUID, actorUserID uuid.UUID, eventType string, messageID *uuid.UUID, data any) {
	if p.messagingRepo == nil {
		return
	}
	payload, err := json.Marshal(data)
	if err != nil {
		return
	}
	stored, err := p.messagingRepo.AppendConversationEvent(ctx, conversationID, actorUserID, eventType, messageID, payload)
	if err != nil {
		return
	}
	if p.realtimeHub == nil {
		return
	}
	userIDs, err := p.messagingRepo.FindConversationMemberUserIDs(ctx, conversationID)
	if err != nil || len(userIDs) == 0 {
		return
	}
	_ = p.realtimeHub.PublishUsers(userIDs, realtime.Event{
		ID:        strconv.FormatInt(stored.ID, 10),
		Type:      eventType,
		Data:      data,
		CreatedAt: stored.CreatedAt.Format(timeFormat),
	})
}

// publishMessageEvent stores and broadcasts one replayable message event.
func (p *MessagingPipe) publishMessageEvent(ctx context.Context, actorUserID uuid.UUID, eventType string, message *models.Message) {
	p.publishConversationEvent(ctx, message.ConversationID, actorUserID, eventType, &message.ID, messageEventPayload{
		ConversationID: message.ConversationID.String(),
		Message:        p.messageResponse(ctx, message),
	})
}

// publishMemberEvent stores and broadcasts one replayable member event.
func (p *MessagingPipe) publishMemberEvent(ctx context.Context, actorUserID uuid.UUID, eventType string, member *models.ConversationMember) {
	personas, err := p.memberPersonas(ctx, []*models.ConversationMember{member})
	if err != nil {
		return
	}
	p.publishConversationEvent(ctx, member.ConversationID, actorUserID, eventType, nil, memberEventPayload{
		ConversationID: member.ConversationID.String(),
		Member:         memberResponses([]*models.ConversationMember{member}, personas)[0],
	})
}

// publishTypingEvent broadcasts one ephemeral typing event.
func (p *MessagingPipe) publishTypingEvent(ctx context.Context, conversationID uuid.UUID, personaID uuid.UUID, active bool) {
	if p.realtimeHub == nil || p.messagingRepo == nil {
		return
	}
	userIDs, err := p.messagingRepo.FindConversationMemberUserIDs(ctx, conversationID)
	if err != nil || len(userIDs) == 0 {
		return
	}
	_ = p.realtimeHub.PublishUsers(userIDs, realtime.Event{
		Type: "typing.updated",
		Data: typingEventPayload{
			ConversationID: conversationID.String(),
			PersonaID:      personaID.String(),
			Active:         active,
		},
	})
}
