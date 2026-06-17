package dtos

import (
	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
)

type CreateMediaAssetDTO struct {
	OwnerPersonaID  uuid.UUID        `json:"owner_persona_id" validate:"required"`
	MediaKind       models.MediaKind `json:"media_kind"`
	StorageKey      string           `json:"storage_key" validate:"required"`
	PlaybackURL     string           `json:"playback_url" validate:"required"`
	ThumbnailURL    string           `json:"thumbnail_url"`
	MimeType        string           `json:"mime_type" validate:"required"`
	DurationSeconds int              `json:"duration_seconds" validate:"required,min=1,max=900"`
	SizeBytes       int64            `json:"size_bytes" validate:"required,min=1"`
}

type InitiatePostMediaUploadDTO struct {
	OwnerPersonaID uuid.UUID        `json:"owner_persona_id" validate:"required"`
	MediaKind      models.MediaKind `json:"media_kind" validate:"required"`
	MimeType       string           `json:"mime_type" validate:"required"`
	SizeBytes      int64            `json:"size_bytes" validate:"required,min=1"`
}

type ConfirmPostMediaUploadDTO struct {
	OwnerPersonaID  uuid.UUID        `json:"owner_persona_id" validate:"required"`
	MediaKind       models.MediaKind `json:"media_kind" validate:"required"`
	PublicID        string           `json:"public_id" validate:"required"`
	SecureURL       string           `json:"secure_url" validate:"required"`
	ThumbnailURL    string           `json:"thumbnail_url"`
	MimeType        string           `json:"mime_type" validate:"required"`
	DurationSeconds int              `json:"duration_seconds"`
	SizeBytes       int64            `json:"size_bytes" validate:"required,min=1"`
}

type InitiateSetVideoUploadDTO struct {
	OwnerPersonaID uuid.UUID `json:"owner_persona_id" validate:"required"`
	MimeType       string    `json:"mime_type" validate:"required"`
	SizeBytes      int64     `json:"size_bytes" validate:"required,min=1"`
}

type InitiateStoryVideoUploadDTO struct {
	OwnerPersonaID uuid.UUID `json:"owner_persona_id" validate:"required"`
	MimeType       string    `json:"mime_type" validate:"required"`
	SizeBytes      int64     `json:"size_bytes" validate:"required,min=1"`
}

type CompleteMediaProcessingDTO struct {
	PlaybackURL     string `json:"playback_url" validate:"required"`
	ThumbnailURL    string `json:"thumbnail_url"`
	MimeType        string `json:"mime_type" validate:"required"`
	DurationSeconds int    `json:"duration_seconds" validate:"required,min=1,max=900"`
	SizeBytes       int64  `json:"size_bytes" validate:"required,min=1"`
}

type FailMediaProcessingDTO struct {
	Reason string `json:"reason"`
}
