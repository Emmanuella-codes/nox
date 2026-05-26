package dtos

type FollowDTO struct {
	PersonaID string `json:"persona_id" validate:"required,uuid"`
}
