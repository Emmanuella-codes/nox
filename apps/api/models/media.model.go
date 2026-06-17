package models

import (
	"time"

	"github.com/google/uuid"
)

type MediaKind string
type MediaProcessingStatus string

const (
	ImageMediaKind MediaKind = "image"
	VideoMediaKind MediaKind = "video"
)

const (
	PendingMediaStatus MediaProcessingStatus = "pending"
	ReadyMediaStatus   MediaProcessingStatus = "ready"
	FailedMediaStatus  MediaProcessingStatus = "failed"
)

type MediaAsset struct {
	ID               uuid.UUID             `json:"id"`
	OwnerUserID      uuid.UUID             `json:"-"`
	OwnerPersonaID   uuid.UUID             `json:"owner_persona_id"`
	MediaKind        MediaKind             `json:"media_kind"`
	StorageKey       string                `json:"storage_key"`
	PlaybackURL      string                `json:"playback_url"`
	ThumbnailURL     string                `json:"thumbnail_url"`
	MimeType         string                `json:"mime_type"`
	DurationSeconds  int                   `json:"duration_seconds"`
	SizeBytes        int64                 `json:"size_bytes"`
	ProcessingStatus MediaProcessingStatus `json:"processing_status"`
	CreatedAt        time.Time             `json:"created_at"`
	UpdatedAt        time.Time             `json:"updated_at"`
}
