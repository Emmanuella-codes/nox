package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type NotificationType string
type NotificationDevicePlatform string
type NotificationOutboxStatus string

const (
	FollowNotificationType                    NotificationType = "follow"
	LikeNotificationType                      NotificationType = "like"
	CommentNotificationType                   NotificationType = "comment"
	RepostNotificationType                    NotificationType = "repost"
	MentionNotificationType                   NotificationType = "mention"
	DirectMessageNotificationType             NotificationType = "message_direct"
	GroupMessageNotificationType              NotificationType = "message_group"
	StoryContributionRequestNotificationType  NotificationType = "story_contribution_request"
	StoryContributionAcceptedNotificationType NotificationType = "story_contribution_accepted"
	StoryContributionRejectedNotificationType NotificationType = "story_contribution_rejected"
	StoryHighlightAddedNotificationType       NotificationType = "story_highlight_added"
	StoryHighlightRemovedNotificationType     NotificationType = "story_highlight_removed"
	StoryReactionNotificationType             NotificationType = "story_reaction"
)

const (
	NotificationDevicePlatformIOS     NotificationDevicePlatform = "ios"
	NotificationDevicePlatformAndroid NotificationDevicePlatform = "android"
	NotificationDevicePlatformWeb     NotificationDevicePlatform = "web"
)

const (
	NotificationOutboxStatusPending    NotificationOutboxStatus = "pending"
	NotificationOutboxStatusProcessing NotificationOutboxStatus = "processing"
	NotificationOutboxStatusSent       NotificationOutboxStatus = "sent"
	NotificationOutboxStatusFailed     NotificationOutboxStatus = "failed"
	NotificationOutboxStatusDead       NotificationOutboxStatus = "dead"
	NotificationOutboxStatusSkipped    NotificationOutboxStatus = "skipped"
)

type Notification struct {
	ID                         uuid.UUID        `json:"id"`
	RecipientUserID            uuid.UUID        `json:"-"`
	RecipientPersonaID         uuid.UUID        `json:"persona_id"`
	ActorPersonaID             *uuid.UUID       `json:"actor_persona_id,omitempty"`
	ActorPostingMode           PostingMode      `json:"actor_posting_mode"`
	ActorAnonymousHandle       string           `json:"actor_anonymous_handle,omitempty"`
	ActorAnonymousAvatarKey    string           `json:"actor_anonymous_avatar_key,omitempty"`
	ConversationID             *uuid.UUID       `json:"conversation_id,omitempty"`
	MessageID                  *uuid.UUID       `json:"message_id,omitempty"`
	PostID                     *uuid.UUID       `json:"post_id,omitempty"`
	CommentID                  *uuid.UUID       `json:"comment_id,omitempty"`
	EventID                    *uuid.UUID       `json:"event_id,omitempty"`
	StoryID                    *uuid.UUID       `json:"story_id,omitempty"`
	StoryItemID                *uuid.UUID       `json:"story_item_id,omitempty"`
	StoryContributionRequestID *uuid.UUID       `json:"story_contribution_request_id,omitempty"`
	IsRead                     bool             `json:"is_read"`
	ReadAt                     *time.Time       `json:"read_at,omitempty"`
	NotificationType           NotificationType `json:"notification_type"`
	CreatedAt                  time.Time        `json:"created_at"`
}

type NotificationDevice struct {
	ID         uuid.UUID                  `json:"id"`
	UserID     uuid.UUID                  `json:"-"`
	InstallID  string                     `json:"install_id"`
	Platform   NotificationDevicePlatform `json:"platform"`
	PushToken  string                     `json:"-"`
	AppVersion string                     `json:"app_version"`
	LastSeenAt time.Time                  `json:"last_seen_at"`
	DisabledAt *time.Time                 `json:"disabled_at,omitempty"`
	CreatedAt  time.Time                  `json:"created_at"`
	UpdatedAt  time.Time                  `json:"updated_at"`
}

type NotificationPreference struct {
	PersonaID        uuid.UUID        `json:"persona_id"`
	NotificationType NotificationType `json:"notification_type"`
	PushEnabled      bool             `json:"push_enabled"`
	UpdatedAt        time.Time        `json:"updated_at"`
}

type NotificationOutbox struct {
	ID                 uuid.UUID                `json:"id"`
	NotificationID     uuid.UUID                `json:"notification_id"`
	RecipientUserID    uuid.UUID                `json:"recipient_user_id"`
	RecipientPersonaID uuid.UUID                `json:"recipient_persona_id"`
	Channel            string                   `json:"channel"`
	Status             NotificationOutboxStatus `json:"status"`
	Payload            json.RawMessage          `json:"payload"`
	AttemptCount       int                      `json:"attempt_count"`
	NextAttemptAt      time.Time                `json:"next_attempt_at"`
	LastError          string                   `json:"last_error"`
	WorkerID           string                   `json:"worker_id"`
	ClaimedAt          *time.Time               `json:"claimed_at,omitempty"`
	SentAt             *time.Time               `json:"sent_at,omitempty"`
	CreatedAt          time.Time                `json:"created_at"`
	UpdatedAt          time.Time                `json:"updated_at"`
}
