package dtos

import (
	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
)

type CreatePostDTO struct {
	PersonaID     *uuid.UUID         `json:"persona_id"`
	PostingMode   models.PostingMode `json:"posting_mode" validate:"required"`
	EventID       uuid.UUID          `json:"event_id"`
	Body          string             `json:"body" validate:"required"`
	PostType      models.PostType    `json:"post_type" validate:"required"`
	MediaURL      string             `json:"media_url"`
	MediaType     models.MediaType   `json:"media_type"`
	MediaAssetIDs []uuid.UUID        `json:"media_asset_ids"`
	Location      string             `json:"location"`
}
