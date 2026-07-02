package dtos

import (
	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
)

type CreateDirectConversationDTO struct {
	SenderPersonaID    uuid.UUID `json:"sender_persona_id" validate:"required"`
	RecipientPersonaID uuid.UUID `json:"recipient_persona_id" validate:"required"`
}

type CreateGroupConversationDTO struct {
	CreatorPersonaID uuid.UUID   `json:"creator_persona_id" validate:"required"`
	Title            string      `json:"title" validate:"required,max=80"`
	MemberPersonaIDs []uuid.UUID `json:"member_persona_ids" validate:"required,min=1,max=49"`
}

type SendMessageDTO struct {
	SenderPersonaID uuid.UUID          `json:"sender_persona_id" validate:"required"`
	Body            string             `json:"body" validate:"max=2000"`
	MessageType     models.MessageType `json:"message_type"`
	MediaAssetID    *uuid.UUID         `json:"media_asset_id"`
}

type AddMembersDTO struct {
	AdminPersonaID   uuid.UUID   `json:"admin_persona_id" validate:"required"`
	MemberPersonaIDs []uuid.UUID `json:"member_persona_ids" validate:"required,min=1,max=50"`
}

type RemoveMemberDTO struct {
	AdminPersonaID uuid.UUID `json:"admin_persona_id" validate:"required"`
}

type MarkReadDTO struct {
	PersonaID uuid.UUID `json:"persona_id" validate:"required"`
	MessageID uuid.UUID `json:"message_id" validate:"required"`
}
