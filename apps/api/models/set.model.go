package models

import (
	"time"

	"github.com/google/uuid"
)

type EmbedType string

const (
	YouTubeEmbedType    EmbedType = "youtube"
	SoundCloudEmbedType EmbedType = "soundcloud"
)

type Set struct {
	ID          uuid.UUID `json:"id"`
	PersonaID   uuid.UUID `json:"persona_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	EmbedURL    string    `json:"embed_url"`
	EmbedType   EmbedType `json:"embed_type"`
	Duration    int       `json:"duration"`
	GenreTags   []string  `json:"genre_tags"`
	PlayCount   int       `json:"play_count"`
	LikeCount   int       `json:"like_count"`
	RecordedAt  time.Time `json:"recorded_at"`
	CreatedAt   time.Time `json:"created_at"`
}
