package main

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/emmanuella-codes/nox/config"
	"github.com/emmanuella-codes/nox/models"
	notification_pipes "github.com/emmanuella-codes/nox/notification/pipes"
	notification_repo "github.com/emmanuella-codes/nox/repositories/notification"
	"github.com/emmanuella-codes/nox/shared/push"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type Worker struct {
	cfg      *config.Config
	repo     notification_repo.NotificationRepository
	provider push.Provider
	workerID string
}

// Builds one worker instance around the configured repository and provider.
func NewWorker(cfg *config.Config, repo notification_repo.NotificationRepository, provider push.Provider) *Worker {
	return &Worker{
		cfg:      cfg,
		repo:     repo,
		provider: provider,
		workerID: uuid.NewString(),
	}
}

// Polls the outbox until shutdown.
func (w *Worker) Run(ctx context.Context) error {
	ticker := time.NewTicker(w.cfg.PushWorkerPollInterval)
	defer ticker.Stop()

	for {
		if err := w.tick(ctx); err != nil {
			log.Error().Err(err).Msg("notification worker tick failed")
		}
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
		}
	}
}

// Claims one batch and processes each queued delivery.
func (w *Worker) tick(ctx context.Context) error {
	outboxes, err := w.repo.ClaimNotificationPushes(ctx, w.workerID, w.cfg.PushWorkerBatchSize)
	if err != nil || len(outboxes) == 0 {
		return err
	}
	for _, outbox := range outboxes {
		if err := w.processOutbox(ctx, outbox); err != nil {
			log.Error().Err(err).Str("outbox_id", outbox.ID.String()).Msg("notification outbox processing failed")
		}
	}
	return nil
}

// Skips stale deliveries and fans valid pushes out to active devices.
func (w *Worker) processOutbox(ctx context.Context, outbox *models.NotificationOutbox) error {
	canDeliver, reason, err := w.repo.CanDeliverNotificationPush(ctx, outbox.ID)
	if err != nil {
		return w.retryOrDead(ctx, outbox, err)
	}
	if !canDeliver {
		return w.repo.MarkNotificationPushSkipped(ctx, outbox.ID, reason)
	}
	var payload notification_pipes.NotificationPushPayload
	if err := json.Unmarshal(outbox.Payload, &payload); err != nil {
		return w.repo.MarkNotificationPushDead(ctx, outbox.ID, "invalid_payload")
	}
	devices, err := w.repo.FindNotificationDevices(ctx, outbox.RecipientUserID)
	if err != nil {
		return w.retryOrDead(ctx, outbox, err)
	}
	activeDevices := activeDevices(devices)
	if len(activeDevices) == 0 {
		return w.repo.MarkNotificationPushSkipped(ctx, outbox.ID, "no_active_devices")
	}
	var delivered bool
	for _, device := range activeDevices {
		err := w.provider.Send(ctx, device, push.Payload{
			Title:      payload.Title,
			Body:       payload.Body,
			TargetPath: payload.TargetPath,
			Badge:      payload.Badge,
			Raw:        outbox.Payload,
		})
		if err == nil {
			delivered = true
			continue
		}
		if errors.Is(err, push.ErrInvalidToken) {
			_ = w.repo.DisableNotificationDeviceByToken(ctx, device.PushToken)
			continue
		}
		return w.retryOrDead(ctx, outbox, err)
	}
	if !delivered {
		return w.repo.MarkNotificationPushSkipped(ctx, outbox.ID, "no_deliverable_devices")
	}
	return w.repo.MarkNotificationPushSent(ctx, outbox.ID)
}

// Reschedules transient failures and caps retry attempts.
func (w *Worker) retryOrDead(ctx context.Context, outbox *models.NotificationOutbox, err error) error {
	if outbox.AttemptCount >= 4 {
		return w.repo.MarkNotificationPushDead(ctx, outbox.ID, err.Error())
	}
	nextAttemptAt := time.Now().Add(backoffForAttempt(outbox.AttemptCount + 1))
	return w.repo.MarkNotificationPushRetry(ctx, outbox.ID, nextAttemptAt, err.Error())
}

// Drops disabled devices before delivery fanout.
func activeDevices(devices []*models.NotificationDevice) []*models.NotificationDevice {
	active := make([]*models.NotificationDevice, 0, len(devices))
	for _, device := range devices {
		if device.DisabledAt == nil {
			active = append(active, device)
		}
	}
	return active
}

// Returns the retry delay for the next failed attempt.
func backoffForAttempt(attempt int) time.Duration {
	switch attempt {
	case 1:
		return 30 * time.Second
	case 2:
		return 2 * time.Minute
	case 3:
		return 10 * time.Minute
	default:
		return 30 * time.Minute
	}
}
