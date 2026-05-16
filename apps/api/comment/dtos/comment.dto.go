package dtos

import "github.com/google/uuid"

type CreateCommentDTO struct {
	PersonaID uuid.UUID  `json:"persona_id" validate:"required"`
	Body      string     `json:"body" validate:"required,max=280"`
	ParentID  *uuid.UUID `json:"parent_id"`
}
