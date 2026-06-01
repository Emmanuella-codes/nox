package dtos

import "github.com/google/uuid"

type CreateSetDTO struct {
	PersonaID    uuid.UUID `json:"persona_id" validate:"required"`
	MediaAssetID uuid.UUID `json:"media_asset_id" validate:"required"`
	Title        string    `json:"title" validate:"required"`
	Description  string    `json:"description"`
	GenreTags    []string  `json:"genre_tags" validate:"required"`
}

type SetPersonaActionDTO struct {
	PersonaID uuid.UUID `json:"persona_id" validate:"required"`
}

type CreateSetCommentDTO struct {
	PersonaID uuid.UUID `json:"persona_id" validate:"required"`
	Body      string    `json:"body" validate:"required"`
	ParentID  uuid.UUID `json:"parent_id"`
}
