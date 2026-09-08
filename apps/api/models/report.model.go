package models

import (
	"time"

	"github.com/google/uuid"
)

type ReportStatus string
type ReportTargetType string

const (
	ReportStatusOpen      ReportStatus = "open"
	ReportStatusReviewing ReportStatus = "reviewing"
	ReportStatusResolved  ReportStatus = "resolved"
	ReportStatusDismissed ReportStatus = "dismissed"
)

const (
	ReportTargetPersona   ReportTargetType = "persona"
	ReportTargetPost      ReportTargetType = "post"
	ReportTargetComment   ReportTargetType = "comment"
	ReportTargetStory     ReportTargetType = "story"
	ReportTargetStoryItem ReportTargetType = "story_item"
	ReportTargetSet       ReportTargetType = "set"
	ReportTargetEvent     ReportTargetType = "event"
	ReportTargetMessage   ReportTargetType = "message"
)

type Report struct {
	ID                uuid.UUID        `json:"id"`
	ReporterUserID    uuid.UUID        `json:"reporter_user_id"`
	ReporterPersonaID uuid.UUID        `json:"reporter_persona_id"`
	TargetType        ReportTargetType `json:"target_type"`
	TargetID          uuid.UUID        `json:"target_id"`
	Reason            string           `json:"reason"`
	Description       string           `json:"description"`
	Status            ReportStatus     `json:"status"`
	AssignedAdminID   *uuid.UUID       `json:"assigned_admin_id,omitempty"`
	ResolutionNote    string           `json:"resolution_note"`
	CreatedAt         time.Time        `json:"created_at"`
	UpdatedAt         time.Time        `json:"updated_at"`
	ResolvedAt        *time.Time       `json:"resolved_at,omitempty"`
}
