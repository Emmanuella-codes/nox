package dtos

import "github.com/google/uuid"

type FollowDTO struct {
	PersonaID uuid.UUID `json:"persona_id" validate:"required"`
}
