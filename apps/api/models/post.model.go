package models

import (
	"time"

	"github.com/google/uuid"
)

type PostType string

const (
	TextPostType     PostType = "text"
	ImagePostType    PostType = "image"
	SetPostType      PostType = "set"
	EventTagPostType PostType = "event_tag"
)

type MediaType string

const (
	ImageMediaType MediaType = "image"
	VideoMediaType MediaType = "video"
)

type Post struct {
	ID           uuid.UUID `json:"id"`
	PersonaID    uuid.UUID `json:"persona_id"`
	EventID      uuid.UUID `json:"event_id"`
	Body         string    `json:"body"` // 350 characters
	PostType     PostType  `json:"post_type"`
	MediaURL     []string  `json:"media_url"`
	MediaType    MediaType `json:"media_type"`
	Location     string    `json:"location"`
	LikeCount    int       `json:"like_count"`
	CommentCount int       `json:"comment_count"`
	RepostCount  int       `json:"repost_count"`
	Is_Repost    bool      `json:"is_repost"`
	RepostOf     uuid.UUID `json:"repost_of"`
	CreatedAt    time.Time `json:"created_at"`
}
