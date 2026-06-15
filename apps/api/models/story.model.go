package models

import (
	"time"

	"github.com/google/uuid"
)

type StoryContributionMode string

const (
	PublicStoryContributionMode    StoryContributionMode = "public"
	FollowersStoryContributionMode StoryContributionMode = "followers"
)

type Story struct {
	ID                   uuid.UUID             `json:"id"`
	EventID              uuid.UUID             `json:"event_id"`
	OwnerUserID          uuid.UUID             `json:"-"`
	OwnerPersonaID       uuid.UUID             `json:"owner_persona_id"`
	Title                string                `json:"title"`
	ContributionMode     StoryContributionMode `json:"contribution_mode"`
	TotalDurationSeconds int                   `json:"total_duration_seconds"`
	ExpiresAt            time.Time             `json:"expires_at"`
	CreatedAt            time.Time             `json:"created_at"`
	UpdatedAt            time.Time             `json:"updated_at"`
}

type StoryItem struct {
	ID                   uuid.UUID   `json:"id"`
	StoryID              uuid.UUID   `json:"story_id"`
	MediaAssetID         uuid.UUID   `json:"media_asset_id"`
	ContributorUserID    uuid.UUID   `json:"-"`
	ContributorPersonaID uuid.UUID   `json:"contributor_persona_id"`
	PostingMode          PostingMode `json:"posting_mode"`
	AnonymousLabel       string      `json:"anonymous_label,omitempty"`
	DurationSeconds      int         `json:"duration_seconds"`
	Position             int         `json:"position"`
	CreatedAt            time.Time   `json:"created_at"`
}

type EventHighlightStory struct {
	ID               uuid.UUID `json:"id"`
	EventID          uuid.UUID `json:"event_id"`
	StoryID          uuid.UUID `json:"story_id"`
	AddedByPersonaID uuid.UUID `json:"added_by_persona_id"`
	Position         int       `json:"position"`
	CreatedAt        time.Time `json:"created_at"`
}
