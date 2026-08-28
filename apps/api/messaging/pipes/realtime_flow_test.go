package pipes

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	messaging_dtos "github.com/emmanuella-codes/nox/messaging/dtos"
	"github.com/emmanuella-codes/nox/models"
	messaging_repo "github.com/emmanuella-codes/nox/repositories/messaging"
	"github.com/emmanuella-codes/nox/shared/realtime"
	"github.com/google/uuid"
)

type realtimeMessagingTestRepo struct {
	messagingTestRepo
	appendedEvents []*models.ConversationEvent
}

// AppendConversationEvent records one replayable event for assertions.
func (r *messagingTestRepo) AppendConversationEvent(ctx context.Context, conversationID uuid.UUID, actorUserID uuid.UUID, eventType string, messageID *uuid.UUID, payload json.RawMessage) (*models.ConversationEvent, error) {
	return &models.ConversationEvent{ID: 1, ConversationID: conversationID, ActorUserID: actorUserID, EventType: eventType, MessageID: messageID, Payload: payload, CreatedAt: time.Now()}, nil
}

// FindConversationEventsAfter returns no replay events in pipe tests.
func (r *messagingTestRepo) FindConversationEventsAfter(ctx context.Context, userID uuid.UUID, afterID int64, limit int) ([]*models.ConversationEvent, error) {
	return nil, nil
}

// FindConversationMemberUserIDs returns distinct active member user ids for one conversation.
func (r *realtimeMessagingTestRepo) FindConversationMemberUserIDs(ctx context.Context, conversationID uuid.UUID) ([]uuid.UUID, error) {
	members := r.membersByConversation[conversationID]
	seen := map[uuid.UUID]bool{}
	userIDs := make([]uuid.UUID, 0, len(members))
	for _, member := range members {
		if seen[member.UserID] {
			continue
		}
		seen[member.UserID] = true
		userIDs = append(userIDs, member.UserID)
	}
	return userIDs, nil
}

// AppendConversationEvent records one replayable event for assertions.
func (r *realtimeMessagingTestRepo) AppendConversationEvent(ctx context.Context, conversationID uuid.UUID, actorUserID uuid.UUID, eventType string, messageID *uuid.UUID, payload json.RawMessage) (*models.ConversationEvent, error) {
	event := &models.ConversationEvent{
		ID:             int64(len(r.appendedEvents) + 1),
		ConversationID: conversationID,
		ActorUserID:    actorUserID,
		EventType:      eventType,
		MessageID:      messageID,
		Payload:        payload,
		CreatedAt:      time.Now(),
	}
	r.appendedEvents = append(r.appendedEvents, event)
	return event, nil
}

// TestSendMessagePipeTrimsIdempotencyKey verifies send requests normalize idempotency keys.
func TestSendMessagePipeTrimsIdempotencyKey(t *testing.T) {
	userID, personaID, conversationID := uuid.New(), uuid.New(), uuid.New()
	repo := &messagingTestRepo{
		conversations: map[uuid.UUID]*models.Conversation{
			conversationID: {ID: conversationID, ConversationType: models.DirectConversationType, CreatedBy: personaID, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		},
		memberByConversationPersona: map[string]*models.ConversationMember{
			memberKey(conversationID, personaID): testMember(conversationID, userID, personaID, models.ConversationMemberRoleMember),
		},
	}
	pipe := NewMessagingPipe(
		repo,
		&messagingTestPersonaRepo{personas: map[uuid.UUID]*models.Persona{personaID: testPersona(userID, personaID, "sender")}},
		&messagingTestMediaRepo{},
		&messagingTestFollowRepo{},
	)

	res := pipe.SendMessagePipe(context.Background(), userID, conversationID, messaging_dtos.SendMessageDTO{
		SenderPersonaID: personaID,
		Body:            "hello",
		IdempotencyKey:  "  retry-key  ",
	})
	if !res.Success {
		t.Fatalf("expected success, got %q", res.Message)
	}
	if repo.createdMessageDTO.IdempotencyKey != "retry-key" {
		t.Fatalf("expected trimmed idempotency key, got %q", repo.createdMessageDTO.IdempotencyKey)
	}
}

// TestSendMessagePipeSkipsReplayableEventOnIdempotentRetry verifies repeated sends do not emit new create events.
func TestSendMessagePipeSkipsReplayableEventOnIdempotentRetry(t *testing.T) {
	userID, personaID, conversationID, messageID := uuid.New(), uuid.New(), uuid.New(), uuid.New()
	repo := &realtimeMessagingTestRepo{
		messagingTestRepo: messagingTestRepo{
			memberByConversationPersona: map[string]*models.ConversationMember{
				memberKey(conversationID, personaID): testMember(conversationID, userID, personaID, models.ConversationMemberRoleMember),
			},
			conversations: map[uuid.UUID]*models.Conversation{
				conversationID: {ID: conversationID, ConversationType: models.DirectConversationType, CreatedBy: personaID, CreatedAt: time.Now(), UpdatedAt: time.Now()},
			},
			membersByConversation: map[uuid.UUID][]*models.ConversationMember{
				conversationID: {testMember(conversationID, userID, personaID, models.ConversationMemberRoleMember)},
			},
			createMessageCreated: false,
			createMessageResult: &models.Message{
				ID:              messageID,
				ConversationID:  conversationID,
				SenderUserID:    userID,
				SenderPersonaID: personaID,
				Body:            "hello",
				MessageType:     models.TextMessageType,
				CreatedAt:       time.Now(),
			},
		},
	}
	pipe := NewMessagingPipe(
		repo,
		&messagingTestPersonaRepo{personas: map[uuid.UUID]*models.Persona{personaID: testPersona(userID, personaID, "sender")}},
		&messagingTestMediaRepo{},
		&messagingTestFollowRepo{},
		realtime.NewHub(),
	)

	res := pipe.SendMessagePipe(context.Background(), userID, conversationID, messaging_dtos.SendMessageDTO{
		SenderPersonaID: personaID,
		Body:            "hello",
		IdempotencyKey:  "retry-key",
	})
	if !res.Success {
		t.Fatalf("expected success, got %q", res.Message)
	}
	if len(repo.appendedEvents) != 0 {
		t.Fatalf("expected no new replayable event, got %d", len(repo.appendedEvents))
	}
}

// TestMarkReadPipeSkipsReplayableEventWhenCursorDoesNotAdvance verifies stale read calls do not emit events.
func TestMarkReadPipeSkipsReplayableEventWhenCursorDoesNotAdvance(t *testing.T) {
	userID, personaID, conversationID, messageID := uuid.New(), uuid.New(), uuid.New(), uuid.New()
	member := testMember(conversationID, userID, personaID, models.ConversationMemberRoleMember)
	member.LastReadMessageID = &messageID
	repo := &realtimeMessagingTestRepo{
		messagingTestRepo: messagingTestRepo{
			memberByConversationPersona: map[string]*models.ConversationMember{
				memberKey(conversationID, personaID): member,
			},
			messages: map[uuid.UUID]*models.Message{
				messageID: {ID: messageID, ConversationID: conversationID, CreatedAt: time.Now()},
			},
		},
	}
	pipe := NewMessagingPipe(
		repo,
		&messagingTestPersonaRepo{personas: map[uuid.UUID]*models.Persona{personaID: testPersona(userID, personaID, "reader")}},
		&messagingTestMediaRepo{},
		&messagingTestFollowRepo{},
		realtime.NewHub(),
	)

	res := pipe.MarkReadPipe(context.Background(), userID, conversationID, messaging_dtos.MarkReadDTO{PersonaID: personaID, MessageID: messageID})
	if !res.Success {
		t.Fatalf("expected success, got %q", res.Message)
	}
	if len(repo.appendedEvents) != 0 {
		t.Fatalf("expected no new replayable event, got %d", len(repo.appendedEvents))
	}
}

var _ messaging_repo.MessagingRepository = (*realtimeMessagingTestRepo)(nil)
