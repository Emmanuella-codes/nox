package models

import (
	"time"

	"github.com/google/uuid"
)

type AdminMediaReference struct {
	EntityType string    `json:"entity_type"`
	EntityID   uuid.UUID `json:"entity_id"`
}

type AdminMediaAsset struct {
	ID                      uuid.UUID             `json:"id"`
	OwnerUserID             uuid.UUID             `json:"owner_user_id"`
	OwnerPersonaID          uuid.UUID             `json:"owner_persona_id"`
	MediaKind               MediaKind             `json:"media_kind"`
	StorageKey              string                `json:"storage_key"`
	PlaybackURL             string                `json:"playback_url"`
	ThumbnailURL            string                `json:"thumbnail_url"`
	MimeType                string                `json:"mime_type"`
	DurationSeconds         int                   `json:"duration_seconds"`
	SizeBytes               int64                 `json:"size_bytes"`
	ProcessingStatus        MediaProcessingStatus `json:"processing_status"`
	ProcessingAttempts      int                   `json:"processing_attempts"`
	LastProcessingError     string                `json:"last_processing_error,omitempty"`
	LastProcessingStartedAt *time.Time            `json:"last_processing_started_at,omitempty"`
	CreatedAt               time.Time             `json:"created_at"`
	UpdatedAt               time.Time             `json:"updated_at"`
	ModerationStatus        ModerationStatus      `json:"moderation_status"`
	ModerationReason        string                `json:"moderation_reason"`
	ModeratedByUserID       *uuid.UUID            `json:"moderated_by_user_id,omitempty"`
	ModeratedAt             *time.Time            `json:"moderated_at,omitempty"`
	IsOrphaned              bool                  `json:"is_orphaned"`
	ReferenceCount          int                   `json:"reference_count"`
	References              []AdminMediaReference `json:"references,omitempty"`
}

type AdminMediaPage struct {
	Items   []AdminMediaAsset `json:"items"`
	Limit   int               `json:"limit"`
	Offset  int               `json:"offset"`
	HasMore bool              `json:"has_more"`
}

type MediaCleanupResult struct {
	DeletedCount int64 `json:"deleted_count"`
}
