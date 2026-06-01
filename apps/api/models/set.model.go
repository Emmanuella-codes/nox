package models

import (
	"time"

	"github.com/google/uuid"
)

type Set struct {
	ID              uuid.UUID   `json:"id"`
	AuthorUserID    uuid.UUID   `json:"-"`
	PersonaID       uuid.UUID   `json:"persona_id"`
	MediaAssetID    uuid.UUID   `json:"media_asset_id"`
	Title           string      `json:"title"`
	Description     string      `json:"description"`
	GenreTags       []string    `json:"genre_tags"`
	DurationSeconds int         `json:"duration_seconds"`
	LikeCount       int         `json:"like_count"`
	CommentCount    int         `json:"comment_count"`
	PlayCount       int         `json:"play_count"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
	Persona         *Persona    `json:"persona,omitempty"`
	MediaAsset      *MediaAsset `json:"media_asset,omitempty"`
}
