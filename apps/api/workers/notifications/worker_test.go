package main

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/emmanuella-codes/nox/config"
	"github.com/emmanuella-codes/nox/models"
	notification_dtos "github.com/emmanuella-codes/nox/notification/dtos"
	notification_pipes "github.com/emmanuella-codes/nox/notification/pipes"
	notification_repo "github.com/emmanuella-codes/nox/repositories/notification"
	"github.com/emmanuella-codes/nox/shared/push"
	"github.com/google/uuid"
)

type workerRepoStub struct {
	canDeliver     bool
	skipReason     string
	canDeliverErr  error
	devices        []*models.NotificationDevice
	findDevicesErr error
	skippedReasons []string
	sentIDs        []uuid.UUID
	disabledTokens []string
	retryCalls     int
	deadCalls      int
}

func (r *workerRepoStub) CreateNotifications(ctx context.Context, inputs []notification_repo.CreateNotificationInput) ([]*models.Notification, error) {
	return nil, nil
}

func (r *workerRepoStub) FindPersonaNotifications(ctx context.Context, userID uuid.UUID, personaID uuid.UUID, limit int, offset int) ([]*models.Notification, error) {
	return nil, nil
}

func (r *workerRepoStub) CountUnreadPersonaNotifications(ctx context.Context, userID uuid.UUID, personaID uuid.UUID) (int, error) {
	return 0, nil
}

func (r *workerRepoStub) MarkNotificationRead(ctx context.Context, notificationID uuid.UUID, userID uuid.UUID, personaID uuid.UUID) (*models.Notification, error) {
	return nil, nil
}

func (r *workerRepoStub) MarkPersonaNotificationsRead(ctx context.Context, userID uuid.UUID, personaID uuid.UUID) (int64, error) {
	return 0, nil
}

func (r *workerRepoStub) MarkConversationNotificationsRead(ctx context.Context, userID uuid.UUID, personaID uuid.UUID, conversationID uuid.UUID, messageID uuid.UUID) (int64, error) {
	return 0, nil
}

func (r *workerRepoStub) DeleteMessageNotifications(ctx context.Context, messageID uuid.UUID) error {
	return nil
}

func (r *workerRepoStub) UpsertNotificationDevice(ctx context.Context, userID uuid.UUID, dto notification_dtos.UpsertNotificationDeviceDTO) (*models.NotificationDevice, error) {
	return nil, nil
}

func (r *workerRepoStub) FindNotificationDevices(ctx context.Context, userID uuid.UUID) ([]*models.NotificationDevice, error) {
	return r.devices, r.findDevicesErr
}

func (r *workerRepoStub) DisableNotificationDevice(ctx context.Context, userID uuid.UUID, deviceID uuid.UUID) error {
	return nil
}

func (r *workerRepoStub) DisableNotificationDeviceByToken(ctx context.Context, pushToken string) error {
	r.disabledTokens = append(r.disabledTokens, pushToken)
	return nil
}

func (r *workerRepoStub) FindNotificationPreferences(ctx context.Context, personaID uuid.UUID) ([]*models.NotificationPreference, error) {
	return nil, nil
}

func (r *workerRepoStub) UpsertNotificationPreference(ctx context.Context, personaID uuid.UUID, notificationType models.NotificationType, pushEnabled bool) (*models.NotificationPreference, error) {
	return nil, nil
}

func (r *workerRepoStub) PushEnabledForPersona(ctx context.Context, personaID uuid.UUID, notificationType models.NotificationType) (bool, error) {
	return true, nil
}

func (r *workerRepoStub) EnqueueNotificationPush(ctx context.Context, notification *models.Notification, payload json.RawMessage) (*models.NotificationOutbox, error) {
	return nil, nil
}

func (r *workerRepoStub) ClaimNotificationPushes(ctx context.Context, workerID string, limit int) ([]*models.NotificationOutbox, error) {
	return nil, nil
}

func (r *workerRepoStub) CanDeliverNotificationPush(ctx context.Context, outboxID uuid.UUID) (bool, string, error) {
	return r.canDeliver, r.skipReason, r.canDeliverErr
}

func (r *workerRepoStub) MarkNotificationPushSent(ctx context.Context, outboxID uuid.UUID) error {
	r.sentIDs = append(r.sentIDs, outboxID)
	return nil
}

func (r *workerRepoStub) MarkNotificationPushRetry(ctx context.Context, outboxID uuid.UUID, nextAttemptAt time.Time, lastError string) error {
	r.retryCalls++
	return nil
}

func (r *workerRepoStub) MarkNotificationPushDead(ctx context.Context, outboxID uuid.UUID, lastError string) error {
	r.deadCalls++
	return nil
}

func (r *workerRepoStub) MarkNotificationPushSkipped(ctx context.Context, outboxID uuid.UUID, reason string) error {
	r.skippedReasons = append(r.skippedReasons, reason)
	return nil
}

type workerProviderStub struct {
	err   error
	calls int
}

func (p *workerProviderStub) Send(ctx context.Context, device *models.NotificationDevice, payload push.Payload) error {
	p.calls++
	return p.err
}

func TestProcessOutboxSkipsAlreadyReadMessagePush(t *testing.T) {
	repo := &workerRepoStub{canDeliver: false, skipReason: "already_read"}
	worker := NewWorker(&config.Config{PushWorkerPollInterval: time.Second, PushWorkerBatchSize: 1}, repo, &workerProviderStub{})
	outbox := &models.NotificationOutbox{ID: uuid.New(), RecipientUserID: uuid.New(), Payload: mustPayload(t)}

	if err := worker.processOutbox(context.Background(), outbox); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(repo.skippedReasons) != 1 || repo.skippedReasons[0] != "already_read" {
		t.Fatalf("expected already_read skip, got %#v", repo.skippedReasons)
	}
}

func TestProcessOutboxSkipsDeletedMessagePush(t *testing.T) {
	repo := &workerRepoStub{canDeliver: false, skipReason: "message_deleted"}
	worker := NewWorker(&config.Config{PushWorkerPollInterval: time.Second, PushWorkerBatchSize: 1}, repo, &workerProviderStub{})
	outbox := &models.NotificationOutbox{ID: uuid.New(), RecipientUserID: uuid.New(), Payload: mustPayload(t)}

	if err := worker.processOutbox(context.Background(), outbox); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(repo.skippedReasons) != 1 || repo.skippedReasons[0] != "message_deleted" {
		t.Fatalf("expected message_deleted skip, got %#v", repo.skippedReasons)
	}
}

func TestProcessOutboxRetriesGuardFailure(t *testing.T) {
	repo := &workerRepoStub{canDeliverErr: errors.New("db unavailable")}
	worker := NewWorker(&config.Config{PushWorkerPollInterval: time.Second, PushWorkerBatchSize: 1}, repo, &workerProviderStub{})
	outbox := &models.NotificationOutbox{ID: uuid.New(), RecipientUserID: uuid.New(), AttemptCount: 1, Payload: mustPayload(t)}

	if err := worker.processOutbox(context.Background(), outbox); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if repo.retryCalls != 1 {
		t.Fatalf("expected one retry call, got %d", repo.retryCalls)
	}
}

func TestProcessOutboxSendsDeliverableMessagePush(t *testing.T) {
	device := &models.NotificationDevice{
		ID:         uuid.New(),
		UserID:     uuid.New(),
		InstallID:  "device-1",
		Platform:   models.NotificationDevicePlatformWeb,
		PushToken:  "push-token",
		LastSeenAt: time.Now(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	repo := &workerRepoStub{
		canDeliver: true,
		devices:    []*models.NotificationDevice{device},
	}
	provider := &workerProviderStub{}
	worker := NewWorker(&config.Config{PushWorkerPollInterval: time.Second, PushWorkerBatchSize: 1}, repo, provider)
	outbox := &models.NotificationOutbox{ID: uuid.New(), RecipientUserID: device.UserID, Payload: mustPayload(t)}

	if err := worker.processOutbox(context.Background(), outbox); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if provider.calls != 1 {
		t.Fatalf("expected one provider call, got %d", provider.calls)
	}
	if len(repo.sentIDs) != 1 || repo.sentIDs[0] != outbox.ID {
		t.Fatalf("expected sent outbox id, got %#v", repo.sentIDs)
	}
}

func mustPayload(t *testing.T) json.RawMessage {
	t.Helper()
	payload, err := json.Marshal(notification_pipes.NotificationPushPayload{
		NotificationID: uuid.NewString(),
		PersonaID:      uuid.NewString(),
		Type:           string(models.DirectMessageNotificationType),
		Title:          "New message",
		Body:           "hello",
		TargetPath:     "/messages/" + uuid.NewString(),
		Badge:          1,
	})
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}
	return payload
}
