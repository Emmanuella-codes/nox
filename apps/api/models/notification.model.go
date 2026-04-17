package models

import (
	"time"

	"github.com/google/uuid"
)

type NotificationType string

const (
	FollowNotificationType  NotificationType = "follow"
	LikeNotificationType    NotificationType = "like"
	CommentNotificationType NotificationType = "comment"
	RepostNotificationType  NotificationType = "repost"
	MentionNotificationType NotificationType = "mention"
)

type Notification struct {
	ID               uuid.UUID        `json:"id"`
	PersonaID        uuid.UUID        `json:"persona_id"`
	ActorID          uuid.UUID        `json:"actor_id"`
	PostID           uuid.UUID        `json:"post_id"`
	CommentID        uuid.UUID        `json:"comment_id"`
	IsRead           bool             `json:"is_read"`
	NotificationType NotificationType `json:"notification_type"`
	CreatedAt        time.Time        `json:"created_at"`
}
