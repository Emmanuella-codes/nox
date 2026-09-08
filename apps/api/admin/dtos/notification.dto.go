package dtos

type NotificationAdminActionDTO struct {
	Reason string `json:"reason" validate:"required,max=500"`
}
