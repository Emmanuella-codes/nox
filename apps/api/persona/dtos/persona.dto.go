package dtos

import "github.com/emmanuella-codes/nox/models"

type CreatePersonaDTO struct {
	Handle      string             `json:"handle"`
	DisplayName string             `json:"display_name" validate:"required"`
	Bio         string             `json:"bio" validate:"required"`
	AvatarURL   string             `json:"avatar_url" validate:"required"`
	CoverURL    string             `json:"cover_url" validate:"required"`
	PersonaType models.PersonaType `json:"persona_type" validate:"required"`
	GenreTags   []string           `json:"genre_tags" validate:"required"`
}

type UpdatePersonaDTO struct {
	DisplayName string   `json:"display_name" validate:"required"`
	Bio         string   `json:"bio" validate:"required"`
	AvatarURL   string   `json:"avatar_url" validate:"required"`
	CoverURL    string   `json:"cover_url" validate:"required"`
	GenreTags   []string `json:"genre_tags" validate:"required"`
}
