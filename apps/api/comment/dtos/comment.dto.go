package dtos

import (
	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
)

type CreateCommentDTO struct {
	AuthorUserID uuid.UUID          `json:"-"`
	PersonaID   uuid.UUID          `json:"persona_id" validate:"required"`
	PostingMode models.PostingMode `json:"posting_mode"`
	Body        string             `json:"body" validate:"required,max=280"`
	ParentID    *uuid.UUID         `json:"parent_id"`
}
