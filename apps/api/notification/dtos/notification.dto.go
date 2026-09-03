package dtos

import "github.com/google/uuid"

type MarkNotificationReadDTO struct {
	PersonaID uuid.UUID `json:"persona_id" validate:"required"`
}

type MarkAllNotificationsReadDTO struct {
	PersonaID uuid.UUID `json:"persona_id" validate:"required"`
}

type UpsertNotificationDeviceDTO struct {
	InstallID  string `json:"install_id" validate:"required,max=120"`
	Platform   string `json:"platform" validate:"required,oneof=ios android web"`
	PushToken  string `json:"push_token" validate:"required,max=512"`
	AppVersion string `json:"app_version" validate:"max=40"`
}

type UpdateNotificationPreferenceDTO struct {
	PersonaID        uuid.UUID `json:"persona_id" validate:"required"`
	NotificationType string    `json:"notification_type" validate:"required,max=80"`
	PushEnabled      bool      `json:"push_enabled"`
}
