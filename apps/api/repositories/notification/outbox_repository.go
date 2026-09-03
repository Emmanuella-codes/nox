package notification

import (
	"context"
	"encoding/json"
	"time"

	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
)

func (r *pgRepository) EnqueueNotificationPush(ctx context.Context, notification *models.Notification, payload json.RawMessage) (*models.NotificationOutbox, error) {
	row := r.db.QueryRow(ctx, `
		INSERT INTO notification_outbox (
			notification_id, recipient_user_id, recipient_persona_id, channel, status, payload
		)
		VALUES ($1, $2, $3, 'push', 'pending', $4)
		ON CONFLICT (notification_id, channel) DO NOTHING
		RETURNING id, notification_id, recipient_user_id, recipient_persona_id, channel, status,
		          payload, attempt_count, next_attempt_at, last_error, worker_id, claimed_at,
		          sent_at, created_at, updated_at
	`, notification.ID, notification.RecipientUserID, notification.RecipientPersonaID, payload)
	return scanNotificationOutbox(row)
}

func (r *pgRepository) ClaimNotificationPushes(ctx context.Context, workerID string, limit int) ([]*models.NotificationOutbox, error) {
	rows, err := r.db.Query(ctx, `
		WITH claimed AS (
			SELECT id
			FROM notification_outbox
			WHERE status IN ('pending', 'failed')
			  AND next_attempt_at <= now()
			ORDER BY next_attempt_at ASC, created_at ASC
			FOR UPDATE SKIP LOCKED
			LIMIT $2
		)
		UPDATE notification_outbox o
		SET status = 'processing',
		    worker_id = $1,
		    claimed_at = now(),
		    updated_at = now()
		FROM claimed
		WHERE o.id = claimed.id
		RETURNING o.id, o.notification_id, o.recipient_user_id, o.recipient_persona_id, o.channel,
		          o.status, o.payload, o.attempt_count, o.next_attempt_at, o.last_error,
		          o.worker_id, o.claimed_at, o.sent_at, o.created_at, o.updated_at
	`, workerID, normalizeClaimLimit(limit))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanNotificationOutboxes(rows)
}

func (r *pgRepository) MarkNotificationPushSent(ctx context.Context, outboxID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `
		UPDATE notification_outbox
		SET status = 'sent',
		    sent_at = now(),
		    updated_at = now()
		WHERE id = $1
	`, outboxID)
	return err
}

func (r *pgRepository) MarkNotificationPushRetry(ctx context.Context, outboxID uuid.UUID, nextAttemptAt time.Time, lastError string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE notification_outbox
		SET status = 'failed',
		    attempt_count = attempt_count + 1,
		    next_attempt_at = $2,
		    last_error = $3,
		    updated_at = now()
		WHERE id = $1
	`, outboxID, nextAttemptAt, lastError)
	return err
}

func (r *pgRepository) MarkNotificationPushDead(ctx context.Context, outboxID uuid.UUID, lastError string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE notification_outbox
		SET status = 'dead',
		    attempt_count = attempt_count + 1,
		    last_error = $2,
		    updated_at = now()
		WHERE id = $1
	`, outboxID, lastError)
	return err
}

func (r *pgRepository) MarkNotificationPushSkipped(ctx context.Context, outboxID uuid.UUID, reason string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE notification_outbox
		SET status = 'skipped',
		    last_error = $2,
		    updated_at = now()
		WHERE id = $1
	`, outboxID, reason)
	return err
}

func normalizeClaimLimit(limit int) int {
	if limit <= 0 {
		return 25
	}
	if limit > 100 {
		return 100
	}
	return limit
}
