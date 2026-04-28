package dtos

import (
	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
)

type CreatePostDTO struct {
	PersonaID uuid.UUID        `json:"persona_id" validate:"required"`
	EventID   uuid.UUID        `json:"event_id" validate:"required"`
	Body      string           `json:"body" validate:"required"`
	PostType  models.PostType  `json:"post_type" validate:"required"`
	MediaURL  string           `json:"media_url"`
	MediaType models.MediaType `json:"media_type"`
	Location  string           `json:"location"`
}
