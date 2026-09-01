package models

import (
	"time"

	"github.com/google/uuid"
)

type StoryContributionMode string
type StoryReactionType string

const (
	PublicStoryContributionMode  StoryContributionMode = "public"
	PrivateStoryContributionMode StoryContributionMode = "private"

	FollowersStoryContributionMode StoryContributionMode = PrivateStoryContributionMode
)

const (
	StoryReactionTypeLike  StoryReactionType = "like"
	StoryReactionTypeFire  StoryReactionType = "fire"
	StoryReactionTypeHeart StoryReactionType = "heart"
	StoryReactionTypeLaugh StoryReactionType = "laugh"
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

type StoryContributionRequestStatus string

const (
	PendingStoryContributionRequestStatus  StoryContributionRequestStatus = "pending"
	AcceptedStoryContributionRequestStatus StoryContributionRequestStatus = "accepted"
	RejectedStoryContributionRequestStatus StoryContributionRequestStatus = "rejected"
)

type StoryContributionRequest struct {
	ID                   uuid.UUID                      `json:"id"`
	StoryID              uuid.UUID                      `json:"story_id"`
	MediaAssetID         uuid.UUID                      `json:"media_asset_id"`
	ContributorUserID    uuid.UUID                      `json:"-"`
	ContributorPersonaID uuid.UUID                      `json:"contributor_persona_id"`
	Status               StoryContributionRequestStatus `json:"status"`
	ReviewedByPersonaID  *uuid.UUID                     `json:"reviewed_by_persona_id,omitempty"`
	StoryItemID          *uuid.UUID                     `json:"story_item_id,omitempty"`
	CreatedAt            time.Time                      `json:"created_at"`
	ReviewedAt           *time.Time                     `json:"reviewed_at,omitempty"`
}

type EventHighlightStory struct {
	ID               uuid.UUID `json:"id"`
	EventID          uuid.UUID `json:"event_id"`
	StoryID          uuid.UUID `json:"story_id"`
	AddedByPersonaID uuid.UUID `json:"added_by_persona_id"`
	Position         int       `json:"position"`
	CreatedAt        time.Time `json:"created_at"`
}

type StoryItemView struct {
	StoryID         uuid.UUID `json:"story_id"`
	StoryItemID     uuid.UUID `json:"story_item_id"`
	ViewerUserID    uuid.UUID `json:"-"`
	ViewerPersonaID uuid.UUID `json:"viewer_persona_id"`
	CreatedAt       time.Time `json:"created_at"`
}

type StoryItemReaction struct {
	StoryID          uuid.UUID         `json:"story_id"`
	StoryItemID      uuid.UUID         `json:"story_item_id"`
	ReactorUserID    uuid.UUID         `json:"-"`
	ReactorPersonaID uuid.UUID         `json:"reactor_persona_id"`
	ReactionType     StoryReactionType `json:"reaction_type"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
}

type ProfileStoryHighlight struct {
	ID             uuid.UUID `json:"id"`
	OwnerPersonaID uuid.UUID `json:"owner_persona_id"`
	StoryID        uuid.UUID `json:"story_id"`
	Position       int       `json:"position"`
	CreatedAt      time.Time `json:"created_at"`
}
