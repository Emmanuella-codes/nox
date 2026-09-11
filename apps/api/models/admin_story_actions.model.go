package models

import (
	"time"

	"github.com/google/uuid"
)

type AdminFeaturedSet struct {
	SetID            uuid.UUID  `json:"set_id"`
	Position         int        `json:"position"`
	FeaturedByUserID uuid.UUID  `json:"featured_by_user_id"`
	CreatedAt        time.Time  `json:"created_at"`
	ExpiresAt        *time.Time `json:"expires_at,omitempty"`
}

type AdminStoryContribution struct {
	ID                   uuid.UUID                      `json:"id"`
	StoryID              uuid.UUID                      `json:"story_id"`
	MediaAssetID         uuid.UUID                      `json:"media_asset_id"`
	ContributorUserID    uuid.UUID                      `json:"contributor_user_id"`
	ContributorPersonaID uuid.UUID                      `json:"contributor_persona_id"`
	Status               StoryContributionRequestStatus `json:"status"`
	ReviewedByPersonaID  *uuid.UUID                     `json:"reviewed_by_persona_id,omitempty"`
	StoryItemID          *uuid.UUID                     `json:"story_item_id,omitempty"`
	CreatedAt            time.Time                      `json:"created_at"`
	ReviewedAt           *time.Time                     `json:"reviewed_at,omitempty"`
}

type AdminHighlight struct {
	ID               uuid.UUID `json:"id"`
	HighlightType    string    `json:"highlight_type"`
	EventID          *uuid.UUID `json:"event_id,omitempty"`
	OwnerPersonaID   *uuid.UUID `json:"owner_persona_id,omitempty"`
	StoryID          uuid.UUID `json:"story_id"`
	AddedByPersonaID *uuid.UUID `json:"added_by_persona_id,omitempty"`
	Position         int       `json:"position"`
	CreatedAt        time.Time `json:"created_at"`
}
