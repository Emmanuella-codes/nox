package dtos

import "github.com/google/uuid"

type LikePostDTO struct {
	PersonaID uuid.UUID `json:"persona_id" validate:"required"`
}
