package models

import (
	"time"

	"github.com/google/uuid"
)

type PersonaType string

const (
	GhostPersonaType   PersonaType = "ghost"
	VisiblePersonaType PersonaType = "visible"
)

type Persona struct {
	ID             uuid.UUID   `json:"id"`
	UserID         uuid.UUID   `json:"user_id"`
	Handle         string      `json:"handle"`
	DisplayName    string      `json:"display_name"`
	Bio            string      `json:"bio"`
	AvatarURL      string      `json:"avatar_url"`
	CoverURL       string      `json:"cover_url"`
	PersonaType    PersonaType `json:"persona_type"`
	GenreTags      []string    `json:"genre_tags"`
	FollowerCount  int         `json:"follower_count"`
	FollowingCount int         `json:"following_count"`
	PostCount      int         `json:"post_count"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
}
