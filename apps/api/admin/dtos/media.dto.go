package dtos

import (
	"time"

	"github.com/emmanuella-codes/nox/models"
)

type CorrectMediaDTO struct {
	Status       models.MediaProcessingStatus `json:"status" validate:"required"`
	PlaybackURL  string                       `json:"playback_url"`
	ThumbnailURL string                       `json:"thumbnail_url"`
	MimeType     string                       `json:"mime_type"`
	Duration     int                          `json:"duration_seconds"`
	SizeBytes    int64                        `json:"size_bytes"`
	Reason       string                       `json:"reason" validate:"max=500"`
}

type ModerateMediaDTO struct {
	Status models.ModerationStatus `json:"status" validate:"required"`
	Reason string                  `json:"reason" validate:"max=500"`
}

type MediaListQuery struct {
	Status     string     `json:"-"`
	Kind       string     `json:"-"`
	Moderation string     `json:"-"`
	OwnerID    string     `json:"-"`
	Orphaned   string     `json:"-"`
	OlderThan  *time.Time `json:"-"`
	Limit      int        `json:"-"`
	Offset     int        `json:"-"`
}
