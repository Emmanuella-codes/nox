package pipes

import (
	"time"

	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
)

type StoryResponse struct {
	ID                   string                       `json:"id"`
	EventID              string                       `json:"event_id"`
	Owner                PersonaResponse              `json:"owner"`
	Title                string                       `json:"title"`
	ContributionMode     models.StoryContributionMode `json:"contribution_mode"`
	TotalDurationSeconds int                          `json:"total_duration_seconds"`
	CanContribute        bool                         `json:"can_contribute"`
	Items                []StoryItemResponse          `json:"items"`
	CreatedAt            time.Time                    `json:"created_at"`
	UpdatedAt            time.Time                    `json:"updated_at"`
}

type PersonaResponse struct {
	ID          string                 `json:"id"`
	Handle      string                 `json:"handle"`
	DisplayName string                 `json:"display_name"`
	AvatarURL   string                 `json:"avatar_url"`
	Category    models.PersonaCategory `json:"category"`
}

type StoryItemResponse struct {
	ID              string             `json:"id"`
	StoryID         string             `json:"story_id"`
	MediaAsset      *models.MediaAsset `json:"media_asset"`
	Contributor     *PersonaResponse   `json:"contributor,omitempty"`
	PostingMode     models.PostingMode `json:"posting_mode"`
	AnonymousLabel  string             `json:"anonymous_label,omitempty"`
	DurationSeconds int                `json:"duration_seconds"`
	Position        int                `json:"position"`
	CreatedAt       time.Time          `json:"created_at"`
}

type StoryListResponse struct {
	Limit      int             `json:"limit"`
	Offset     int             `json:"offset"`
	HasMore    bool            `json:"has_more"`
	NextOffset *int            `json:"next_offset,omitempty"`
	Stories    []StoryResponse `json:"stories"`
}

type EventHighlightStoryResponse struct {
	ID               string         `json:"id"`
	EventID          string         `json:"event_id"`
	StoryID          string         `json:"story_id"`
	AddedByPersonaID string         `json:"added_by_persona_id"`
	Position         int            `json:"position"`
	Story            *StoryResponse `json:"story,omitempty"`
	CreatedAt        time.Time      `json:"created_at"`
}

func personaResponse(persona *models.Persona) PersonaResponse {
	return PersonaResponse{
		ID:          persona.ID.String(),
		Handle:      persona.Handle,
		DisplayName: persona.DisplayName,
		AvatarURL:   persona.AvatarURL,
		Category:    persona.Category,
	}
}

func nextOffset(limit int, offset int, hasMore bool) *int {
	if !hasMore {
		return nil
	}
	next := offset + limit
	return &next
}

func uuidString(id uuid.UUID) string {
	return id.String()
}
