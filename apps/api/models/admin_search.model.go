package models

import (
	"time"

	"github.com/google/uuid"
)

type AdminHashtag struct {
	ID                 uuid.UUID        `json:"id"`
	Tag                string           `json:"tag"`
	PostCount          int              `json:"post_count"`
	CreatedAt          time.Time        `json:"created_at"`
	ModerationStatus   ModerationStatus `json:"moderation_status"`
	ModerationReason   string           `json:"moderation_reason"`
	ModeratedByUserID  *uuid.UUID       `json:"moderated_by_user_id,omitempty"`
	ModeratedAt        *time.Time       `json:"moderated_at,omitempty"`
	IsSuppressed       bool             `json:"is_suppressed"`
	SuppressionReason  string           `json:"suppression_reason,omitempty"`
	SuppressionExpires *time.Time       `json:"suppression_expires_at,omitempty"`
}

type AdminHashtagPage struct {
	Items   []AdminHashtag `json:"items"`
	Limit   int            `json:"limit"`
	Offset  int            `json:"offset"`
	HasMore bool           `json:"has_more"`
}

type SearchSuppression struct {
	ID              uuid.UUID  `json:"id"`
	NormalizedQuery string     `json:"normalized_query"`
	Reason          string     `json:"reason"`
	CreatedBy       uuid.UUID  `json:"created_by_user_id"`
	CreatedAt       time.Time  `json:"created_at"`
	ExpiresAt       *time.Time `json:"expires_at,omitempty"`
}

type SearchSuppressionPage struct {
	Items   []SearchSuppression `json:"items"`
	Limit   int                 `json:"limit"`
	Offset  int                 `json:"offset"`
	HasMore bool                `json:"has_more"`
}

type SearchReindexResult struct {
	CacheKeysDeleted int `json:"cache_keys_deleted"`
}
