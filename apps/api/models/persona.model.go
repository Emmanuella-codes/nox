package models

import (
	"time"

	"github.com/google/uuid"
)

type PersonaType string
type PersonaCategory string

const (
	GhostPersonaType   PersonaType = "ghost"
	VisiblePersonaType PersonaType = "visible"
)

const (
	PatronPersonaCategory    PersonaCategory = "patron"
	DJPersonaCategory        PersonaCategory = "dj"
	OrganizerPersonaCategory PersonaCategory = "organizer"

	FanPersonaCategory     PersonaCategory = PatronPersonaCategory
	CreatorPersonaCategory PersonaCategory = PatronPersonaCategory
)

type Persona struct {
	ID             uuid.UUID       `json:"id"`
	UserID         uuid.UUID       `json:"-"`
	Handle         string          `json:"handle"`
	DisplayName    string          `json:"display_name"`
	Bio            string          `json:"bio"`
	AvatarURL      string          `json:"avatar_url"`
	CoverURL       string          `json:"cover_url"`
	PersonaType    PersonaType     `json:"persona_type"`
	Category       PersonaCategory `json:"category"`
	GenreTags      []string        `json:"genre_tags"`
	FollowerCount  int             `json:"follower_count"`
	FollowingCount int             `json:"following_count"`
	PostCount      int             `json:"post_count"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

// IsOwnedBy reports whether this public profile belongs to the given user.
func (p *Persona) IsOwnedBy(userID uuid.UUID) bool {
	return p != nil && p.UserID == userID
}
