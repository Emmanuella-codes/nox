package dtos

import "github.com/google/uuid"

type CreateMediaAssetDTO struct {
	OwnerPersonaID  uuid.UUID `json:"owner_persona_id" validate:"required"`
	StorageKey      string    `json:"storage_key" validate:"required"`
	PlaybackURL     string    `json:"playback_url" validate:"required"`
	ThumbnailURL    string    `json:"thumbnail_url"`
	MimeType        string    `json:"mime_type" validate:"required"`
	DurationSeconds int       `json:"duration_seconds" validate:"required,min=1,max=900"`
	SizeBytes       int64     `json:"size_bytes" validate:"required,min=1"`
}
