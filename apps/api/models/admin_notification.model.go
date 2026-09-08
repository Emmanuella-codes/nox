package models

import (
	"time"

	"github.com/google/uuid"
)

type AdminNotificationDevice struct {
	ID         uuid.UUID                  `json:"id"`
	UserID     uuid.UUID                  `json:"user_id"`
	InstallID  string                     `json:"install_id"`
	Platform   NotificationDevicePlatform `json:"platform"`
	TokenHint  string                     `json:"token_hint"`
	AppVersion string                     `json:"app_version"`
	LastSeenAt time.Time                  `json:"last_seen_at"`
	DisabledAt *time.Time                 `json:"disabled_at,omitempty"`
	CreatedAt  time.Time                  `json:"created_at"`
	UpdatedAt  time.Time                  `json:"updated_at"`
}

type AdminNotificationOutbox struct {
	ID                 uuid.UUID                `json:"id"`
	NotificationID     uuid.UUID                `json:"notification_id"`
	RecipientUserID    uuid.UUID                `json:"recipient_user_id"`
	RecipientPersonaID uuid.UUID                `json:"recipient_persona_id"`
	NotificationType   NotificationType         `json:"notification_type"`
	Channel            string                   `json:"channel"`
	Status             NotificationOutboxStatus `json:"status"`
	PayloadRedacted    bool                     `json:"payload_redacted"`
	AttemptCount       int                      `json:"attempt_count"`
	NextAttemptAt      time.Time                `json:"next_attempt_at"`
	LastError          string                   `json:"last_error"`
	WorkerID           string                   `json:"worker_id"`
	ClaimedAt          *time.Time               `json:"claimed_at,omitempty"`
	SentAt             *time.Time               `json:"sent_at,omitempty"`
	CreatedAt          time.Time                `json:"created_at"`
	UpdatedAt          time.Time                `json:"updated_at"`
}
