package dtos

import (
	"time"

	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
)

type CreateStoryDTO struct {
	EventID          uuid.UUID                    `json:"event_id" validate:"required"`
	OwnerPersonaID   uuid.UUID                    `json:"owner_persona_id" validate:"required"`
	Title            string                       `json:"title" validate:"required"`
	ContributionMode models.StoryContributionMode `json:"contribution_mode" validate:"required"`
	ExpiresAt        time.Time                    `json:"expires_at"`
}

type AddStoryItemDTO struct {
	ContributorPersonaID uuid.UUID          `json:"contributor_persona_id" validate:"required"`
	MediaAssetID         uuid.UUID          `json:"media_asset_id" validate:"required"`
	PostingMode          models.PostingMode `json:"posting_mode" validate:"required"`
}

type CreateStoryContributionRequestDTO struct {
	ContributorPersonaID uuid.UUID `json:"contributor_persona_id" validate:"required"`
	MediaAssetID         uuid.UUID `json:"media_asset_id" validate:"required"`
}

type AddEventHighlightStoryDTO struct {
	StoryID          uuid.UUID `json:"story_id" validate:"required"`
	AddedByPersonaID uuid.UUID `json:"added_by_persona_id" validate:"required"`
}

type ReorderStoryItemDTO struct {
	PersonaID uuid.UUID `json:"persona_id" validate:"required"`
	Position  int       `json:"position" validate:"required,min=1"`
}

type ReorderEventHighlightStoryDTO struct {
	PersonaID uuid.UUID `json:"persona_id" validate:"required"`
	Position  int       `json:"position" validate:"required,min=1"`
}
