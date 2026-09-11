package models

import (
	"time"

	"github.com/google/uuid"
)

type AdminContentRecord struct {
	EntityType       ModerationEntityType `json:"entity_type"`
	EntityID         uuid.UUID            `json:"entity_id"`
	Data             map[string]any       `json:"data"`
	ModerationStatus ModerationStatus     `json:"moderation_status"`
	ModerationReason string               `json:"moderation_reason"`
	ModeratedAt      *time.Time           `json:"moderated_at,omitempty"`
	ModeratedBy      *uuid.UUID           `json:"moderated_by,omitempty"`
}

type AdminContentPage struct {
	EntityType string               `json:"entity_type"`
	Items      []AdminContentRecord `json:"items"`
	Limit      int                  `json:"limit"`
	Offset     int                  `json:"offset"`
	HasMore    bool                 `json:"has_more"`
}
