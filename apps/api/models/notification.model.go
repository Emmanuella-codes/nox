package models

import (
	"time"

	"github.com/google/uuid"
)

type NotificationType string

const (
	FollowNotificationType        NotificationType = "follow"
	LikeNotificationType          NotificationType = "like"
	CommentNotificationType       NotificationType = "comment"
	RepostNotificationType        NotificationType = "repost"
	MentionNotificationType       NotificationType = "mention"
	DirectMessageNotificationType NotificationType = "message_direct"
	GroupMessageNotificationType  NotificationType = "message_group"
)

type Notification struct {
	ID                      uuid.UUID        `json:"id"`
	RecipientUserID         uuid.UUID        `json:"-"`
	RecipientPersonaID      uuid.UUID        `json:"persona_id"`
	ActorPersonaID          *uuid.UUID       `json:"actor_persona_id,omitempty"`
	ActorPostingMode        PostingMode      `json:"actor_posting_mode"`
	ActorAnonymousHandle    string           `json:"actor_anonymous_handle,omitempty"`
	ActorAnonymousAvatarKey string           `json:"actor_anonymous_avatar_key,omitempty"`
	ConversationID          *uuid.UUID       `json:"conversation_id,omitempty"`
	MessageID               *uuid.UUID       `json:"message_id,omitempty"`
	PostID                  *uuid.UUID       `json:"post_id,omitempty"`
	CommentID               *uuid.UUID       `json:"comment_id,omitempty"`
	IsRead                  bool             `json:"is_read"`
	ReadAt                  *time.Time       `json:"read_at,omitempty"`
	NotificationType        NotificationType `json:"notification_type"`
	CreatedAt               time.Time        `json:"created_at"`
}
