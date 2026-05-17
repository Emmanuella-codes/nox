package dtos

type RefreshDTO struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
