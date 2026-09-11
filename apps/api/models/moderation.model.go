package models

import (
	"time"

	"github.com/google/uuid"
)

type ModerationStatus string
type ModerationEntityType string

const (
	ModerationStatusActive  ModerationStatus = "active"
	ModerationStatusHidden  ModerationStatus = "hidden"
	ModerationStatusRemoved ModerationStatus = "removed"
)

const (
	ModerationEntityPersona   ModerationEntityType = "persona"
	ModerationEntityPost      ModerationEntityType = "post"
	ModerationEntityComment   ModerationEntityType = "comment"
	ModerationEntityStory     ModerationEntityType = "story"
	ModerationEntityStoryItem ModerationEntityType = "story_item"
	ModerationEntitySet       ModerationEntityType = "set"
	ModerationEntityEvent     ModerationEntityType = "event"
)

type ModerationAction struct {
	ID             uuid.UUID            `json:"id"`
	EntityType     ModerationEntityType `json:"entity_type"`
	EntityID       uuid.UUID            `json:"entity_id"`
	PreviousStatus *ModerationStatus    `json:"previous_status,omitempty"`
	Status         ModerationStatus     `json:"status"`
	Reason         string               `json:"reason"`
	AdminUserID    uuid.UUID            `json:"admin_user_id"`
	CreatedAt      time.Time            `json:"created_at"`
}

type ModerationState struct {
	EntityType  ModerationEntityType `json:"entity_type"`
	EntityID    uuid.UUID            `json:"entity_id"`
	Status      ModerationStatus     `json:"status"`
	Reason      string               `json:"reason"`
	ModeratedAt *time.Time           `json:"moderated_at,omitempty"`
	ModeratedBy *uuid.UUID           `json:"moderated_by,omitempty"`
}

func PublicModerationStatuses() []ModerationStatus {
	return []ModerationStatus{ModerationStatusActive}
}
