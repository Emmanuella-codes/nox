package models

import "github.com/google/uuid"

type Hashtag struct {
	ID        uuid.UUID `json:"id"`
	Tag       string    `json:"tag"`
	PostCount int       `json:"post_count"`
}

type PostHashtag struct {
	PostID    uuid.UUID `json:"post_id"`
	HashtagID uuid.UUID `json:"hashtag_id"`
}
