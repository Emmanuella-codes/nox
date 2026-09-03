package pipes

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	messaging_dtos "github.com/emmanuella-codes/nox/messaging/dtos"
	"github.com/emmanuella-codes/nox/models"
	notification_dtos "github.com/emmanuella-codes/nox/notification/dtos"
	notification_repo "github.com/emmanuella-codes/nox/repositories/notification"
	"github.com/google/uuid"
)

type notificationMessagingTestRepo struct {
	messagingTestRepo
	readMember *models.ConversationMember
}

type notificationRepoStub struct {
	createdInputs     []notification_repo.CreateNotificationInput
	markedReadCalls   int
	deletedMessageIDs []uuid.UUID
}

// MarkConversationRead returns the configured advanced read member.
func (r *notificationMessagingTestRepo) MarkConversationRead(ctx context.Context, conversationID uuid.UUID, personaID uuid.UUID, messageID uuid.UUID) (*models.ConversationMember, error) {
	if r.readMember != nil {
		return r.readMember, nil
	}
	return r.messagingTestRepo.MarkConversationRead(ctx, conversationID, personaID, messageID)
}

// CreateNotifications records created notification inputs.
func (r *notificationRepoStub) CreateNotifications(ctx context.Context, inputs []notification_repo.CreateNotificationInput) ([]*models.Notification, error) {
	r.createdInputs = append(r.createdInputs, inputs...)
	return nil, nil
}

// FindPersonaNotifications is unused in these tests.
func (r *notificationRepoStub) FindPersonaNotifications(ctx context.Context, userID uuid.UUID, personaID uuid.UUID, limit int, offset int) ([]*models.Notification, error) {
	return nil, nil
}

// CountUnreadPersonaNotifications is unused in these tests.
func (r *notificationRepoStub) CountUnreadPersonaNotifications(ctx context.Context, userID uuid.UUID, personaID uuid.UUID) (int, error) {
	return 0, nil
}

// MarkNotificationRead is unused in these tests.
func (r *notificationRepoStub) MarkNotificationRead(ctx context.Context, notificationID uuid.UUID, userID uuid.UUID, personaID uuid.UUID) (*models.Notification, error) {
	return nil, nil
}

// MarkPersonaNotificationsRead is unused in these tests.
func (r *notificationRepoStub) MarkPersonaNotificationsRead(ctx context.Context, userID uuid.UUID, personaID uuid.UUID) (int64, error) {
	return 0, nil
}

// MarkConversationNotificationsRead records one conversation-read sync call.
func (r *notificationRepoStub) MarkConversationNotificationsRead(ctx context.Context, userID uuid.UUID, personaID uuid.UUID, conversationID uuid.UUID, messageID uuid.UUID) (int64, error) {
	r.markedReadCalls++
	return 1, nil
}

// DeleteMessageNotifications records one deleted-message cleanup call.
func (r *notificationRepoStub) DeleteMessageNotifications(ctx context.Context, messageID uuid.UUID) error {
	r.deletedMessageIDs = append(r.deletedMessageIDs, messageID)
	return nil
}

// UpsertNotificationDevice is unused in these tests.
func (r *notificationRepoStub) UpsertNotificationDevice(ctx context.Context, userID uuid.UUID, dto notification_dtos.UpsertNotificationDeviceDTO) (*models.NotificationDevice, error) {
	return nil, nil
}

// FindNotificationDevices is unused in these tests.
func (r *notificationRepoStub) FindNotificationDevices(ctx context.Context, userID uuid.UUID) ([]*models.NotificationDevice, error) {
	return nil, nil
}

// DisableNotificationDevice is unused in these tests.
func (r *notificationRepoStub) DisableNotificationDevice(ctx context.Context, userID uuid.UUID, deviceID uuid.UUID) error {
	return nil
}

// DisableNotificationDeviceByToken is unused in these tests.
func (r *notificationRepoStub) DisableNotificationDeviceByToken(ctx context.Context, pushToken string) error {
	return nil
}

// FindNotificationPreferences is unused in these tests.
func (r *notificationRepoStub) FindNotificationPreferences(ctx context.Context, personaID uuid.UUID) ([]*models.NotificationPreference, error) {
	return nil, nil
}

// UpsertNotificationPreference is unused in these tests.
func (r *notificationRepoStub) UpsertNotificationPreference(ctx context.Context, personaID uuid.UUID, notificationType models.NotificationType, pushEnabled bool) (*models.NotificationPreference, error) {
	return nil, nil
}

// PushEnabledForPersona is unused in these tests.
func (r *notificationRepoStub) PushEnabledForPersona(ctx context.Context, personaID uuid.UUID, notificationType models.NotificationType) (bool, error) {
	return true, nil
}

// EnqueueNotificationPush is unused in these tests.
func (r *notificationRepoStub) EnqueueNotificationPush(ctx context.Context, notification *models.Notification, payload json.RawMessage) (*models.NotificationOutbox, error) {
	return nil, nil
}

// ClaimNotificationPushes is unused in these tests.
func (r *notificationRepoStub) ClaimNotificationPushes(ctx context.Context, workerID string, limit int) ([]*models.NotificationOutbox, error) {
	return nil, nil
}

// MarkNotificationPushSent is unused in these tests.
func (r *notificationRepoStub) MarkNotificationPushSent(ctx context.Context, outboxID uuid.UUID) error {
	return nil
}

// MarkNotificationPushRetry is unused in these tests.
func (r *notificationRepoStub) MarkNotificationPushRetry(ctx context.Context, outboxID uuid.UUID, nextAttemptAt time.Time, lastError string) error {
	return nil
}

// MarkNotificationPushDead is unused in these tests.
func (r *notificationRepoStub) MarkNotificationPushDead(ctx context.Context, outboxID uuid.UUID, lastError string) error {
	return nil
}

// MarkNotificationPushSkipped is unused in these tests.
func (r *notificationRepoStub) MarkNotificationPushSkipped(ctx context.Context, outboxID uuid.UUID, reason string) error {
	return nil
}

// TestSendMessagePipeCreatesMessageNotifications verifies recipient notifications are persisted for new messages.
func TestSendMessagePipeCreatesMessageNotifications(t *testing.T) {
	userID, recipientUserID := uuid.New(), uuid.New()
	personaID, recipientPersonaID := uuid.New(), uuid.New()
	conversationID, messageID := uuid.New(), uuid.New()
	repo := &notificationMessagingTestRepo{
		messagingTestRepo: messagingTestRepo{
			conversations: map[uuid.UUID]*models.Conversation{
				conversationID: {ID: conversationID, ConversationType: models.GroupConversationType, CreatedBy: personaID, CreatedAt: time.Now(), UpdatedAt: time.Now()},
			},
			memberByConversationPersona: map[string]*models.ConversationMember{
				memberKey(conversationID, personaID): testMember(conversationID, userID, personaID, models.ConversationMemberRoleAdmin),
			},
			membersByConversation: map[uuid.UUID][]*models.ConversationMember{
				conversationID: {
					testMember(conversationID, userID, personaID, models.ConversationMemberRoleAdmin),
					testMember(conversationID, recipientUserID, recipientPersonaID, models.ConversationMemberRoleMember),
				},
			},
			createMessageCreated: true,
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
	notifications := &notificationRepoStub{}
	pipe := NewMessagingPipe(
		repo,
		&messagingTestPersonaRepo{personas: map[uuid.UUID]*models.Persona{personaID: testPersona(userID, personaID, "sender")}},
		&messagingTestMediaRepo{},
		&messagingTestFollowRepo{},
		notifications,
	)

	res := pipe.SendMessagePipe(context.Background(), userID, conversationID, messaging_dtos.SendMessageDTO{
		SenderPersonaID: personaID,
		Body:            "hello",
	})
	if !res.Success {
		t.Fatalf("expected success, got %q", res.Message)
	}
	if len(notifications.createdInputs) != 1 {
		t.Fatalf("expected 1 notification input, got %d", len(notifications.createdInputs))
	}
	if notifications.createdInputs[0].NotificationType != models.GroupMessageNotificationType {
		t.Fatalf("expected group message notification, got %q", notifications.createdInputs[0].NotificationType)
	}
	if notifications.createdInputs[0].RecipientUserID != recipientUserID || notifications.createdInputs[0].RecipientPersonaID != recipientPersonaID {
		t.Fatal("expected notification to target the non-sender member")
	}
}

// TestMarkReadPipeMarksConversationNotifications verifies conversation reads also clear message notifications.
func TestMarkReadPipeMarksConversationNotifications(t *testing.T) {
	userID, personaID, conversationID, messageID := uuid.New(), uuid.New(), uuid.New(), uuid.New()
	before := testMember(conversationID, userID, personaID, models.ConversationMemberRoleMember)
	after := testMember(conversationID, userID, personaID, models.ConversationMemberRoleMember)
	after.LastReadMessageID = &messageID
	repo := &notificationMessagingTestRepo{
		messagingTestRepo: messagingTestRepo{
			memberByConversationPersona: map[string]*models.ConversationMember{
				memberKey(conversationID, personaID): before,
			},
			messages: map[uuid.UUID]*models.Message{
				messageID: {ID: messageID, ConversationID: conversationID, CreatedAt: time.Now()},
			},
		},
		readMember: after,
	}
	notifications := &notificationRepoStub{}
	pipe := NewMessagingPipe(
		repo,
		&messagingTestPersonaRepo{personas: map[uuid.UUID]*models.Persona{personaID: testPersona(userID, personaID, "reader")}},
		&messagingTestMediaRepo{},
		&messagingTestFollowRepo{},
		notifications,
	)

	res := pipe.MarkReadPipe(context.Background(), userID, conversationID, messaging_dtos.MarkReadDTO{PersonaID: personaID, MessageID: messageID})
	if !res.Success {
		t.Fatalf("expected success, got %q", res.Message)
	}
	if notifications.markedReadCalls != 1 {
		t.Fatalf("expected one notification read sync, got %d", notifications.markedReadCalls)
	}
}

// TestDeleteMessagePipeDeletesMessageNotifications verifies deleted messages clear tied notifications.
func TestDeleteMessagePipeDeletesMessageNotifications(t *testing.T) {
	userID, personaID, conversationID, messageID := uuid.New(), uuid.New(), uuid.New(), uuid.New()
	now := time.Now()
	repo := &notificationMessagingTestRepo{
		messagingTestRepo: messagingTestRepo{
			messages: map[uuid.UUID]*models.Message{
				messageID: {ID: messageID, ConversationID: conversationID, SenderUserID: userID, SenderPersonaID: personaID, CreatedAt: now},
			},
			deletedMessageResult: &models.Message{
				ID:              messageID,
				ConversationID:  conversationID,
				SenderUserID:    userID,
				SenderPersonaID: personaID,
				CreatedAt:       now,
				DeletedAt:       timePtr(now),
			},
		},
	}
	notifications := &notificationRepoStub{}
	pipe := NewMessagingPipe(repo, &messagingTestPersonaRepo{}, &messagingTestMediaRepo{}, &messagingTestFollowRepo{}, notifications)

	res := pipe.DeleteMessagePipe(context.Background(), userID, messageID)
	if !res.Success {
		t.Fatalf("expected success, got %q", res.Message)
	}
	if len(notifications.deletedMessageIDs) != 1 || notifications.deletedMessageIDs[0] != messageID {
		t.Fatal("expected deleted message notifications to be cleared")
	}
}

var _ notification_repo.NotificationRepository = (*notificationRepoStub)(nil)
