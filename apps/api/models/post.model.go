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
	ImageMediaType      MediaType = "image"
	YouTubeMediaType    MediaType = "youtube"
	SoundCloudMediaType MediaType = "soundcloud"
)

type PostingMode string

const (
	PublicPostingMode    PostingMode = "public"
	AnonymousPostingMode PostingMode = "anonymous"
)

type Post struct {
	ID           uuid.UUID   `json:"id"`
	AuthorUserID uuid.UUID   `json:"-"`
	PersonaID    *uuid.UUID  `json:"persona_id,omitempty"`
	PostingMode  PostingMode `json:"posting_mode"`
	EventID      *uuid.UUID  `json:"event_id,omitempty"`
	Body         string      `json:"body"`
	PostType     PostType    `json:"post_type"`
	MediaURL     string      `json:"media_url,omitempty"`
	MediaType    MediaType   `json:"media_type,omitempty"`
	Location     string      `json:"location,omitempty"`
	LikeCount    int         `json:"like_count"`
	CommentCount int         `json:"comment_count"`
	RepostCount  int         `json:"repost_count"`
	IsRepost     bool        `json:"is_repost"`
	RepostOf     *uuid.UUID  `json:"repost_of,omitempty"`
	CreatedAt    time.Time   `json:"created_at"`
}
