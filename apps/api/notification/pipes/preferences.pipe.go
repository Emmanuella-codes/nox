package pipes

import (
	"context"
	"slices"

	"github.com/emmanuella-codes/nox/models"
	"github.com/emmanuella-codes/nox/notification/messages"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

type NotificationDeviceResponse struct {
	ID         string  `json:"id"`
	InstallID  string  `json:"install_id"`
	Platform   string  `json:"platform"`
	AppVersion string  `json:"app_version"`
	LastSeenAt string  `json:"last_seen_at"`
	DisabledAt *string `json:"disabled_at,omitempty"`
}

type NotificationPreferenceResponse struct {
	PersonaID        string                  `json:"persona_id"`
	NotificationType models.NotificationType `json:"notification_type"`
	PushEnabled      bool                    `json:"push_enabled"`
	UpdatedAt        string                  `json:"updated_at"`
}

func (p *NotificationPipe) ListNotificationDevicesPipe(ctx context.Context, userID uuid.UUID) *shared.PipeRes[[]NotificationDeviceResponse] {
	devices, err := p.notificationRepo.FindNotificationDevices(ctx, userID)
	if err != nil {
		return pipeInternalError[[]NotificationDeviceResponse](err, "notification.list_devices")
	}
	responses := make([]NotificationDeviceResponse, 0, len(devices))
	for _, device := range devices {
		responses = append(responses, notificationDeviceResponse(device))
	}
	return shared.PipeSuccess(messages.Notification_Devices_Listed, &responses)
}

func (p *NotificationPipe) ListNotificationPreferencesPipe(ctx context.Context, userID uuid.UUID, personaID uuid.UUID) *shared.PipeRes[[]NotificationPreferenceResponse] {
	if _, message := p.profilePersona(ctx, userID, personaID); message != "" {
		return shared.PipeError[[]NotificationPreferenceResponse](message)
	}
	existing, err := p.notificationRepo.FindNotificationPreferences(ctx, personaID)
	if err != nil {
		return pipeInternalError[[]NotificationPreferenceResponse](err, "notification.list_preferences")
	}
	byType := make(map[models.NotificationType]*models.NotificationPreference, len(existing))
	for _, preference := range existing {
		byType[preference.NotificationType] = preference
	}
	responses := make([]NotificationPreferenceResponse, 0, len(notificationPreferenceTypes()))
	for _, notificationType := range notificationPreferenceTypes() {
		preference := byType[notificationType]
		if preference == nil {
			responses = append(responses, NotificationPreferenceResponse{
				PersonaID:        personaID.String(),
				NotificationType: notificationType,
				PushEnabled:      true,
			})
			continue
		}
		responses = append(responses, notificationPreferenceResponse(preference))
	}
	return shared.PipeSuccess(messages.Notification_Preferences_Listed, &responses)
}

func (p *NotificationPipe) UpdateNotificationPreferencePipe(ctx context.Context, userID uuid.UUID, personaID uuid.UUID, notificationType models.NotificationType, pushEnabled bool) *shared.PipeRes[NotificationPreferenceResponse] {
	if _, message := p.profilePersona(ctx, userID, personaID); message != "" {
		return shared.PipeError[NotificationPreferenceResponse](message)
	}
	if !slices.Contains(notificationPreferenceTypes(), notificationType) {
		return shared.PipeError[NotificationPreferenceResponse](messages.Invalid_Payload)
	}
	preference, err := p.notificationRepo.UpsertNotificationPreference(ctx, personaID, notificationType, pushEnabled)
	if err != nil {
		return pipeInternalError[NotificationPreferenceResponse](err, "notification.update_preference")
	}
	response := notificationPreferenceResponse(preference)
	return shared.PipeSuccess(messages.Notification_Preference_Updated, &response)
}

func notificationDeviceResponse(device *models.NotificationDevice) NotificationDeviceResponse {
	var disabledAt *string
	if device.DisabledAt != nil {
		value := device.DisabledAt.Format(timeFormat)
		disabledAt = &value
	}
	return NotificationDeviceResponse{
		ID:         device.ID.String(),
		InstallID:  device.InstallID,
		Platform:   string(device.Platform),
		AppVersion: device.AppVersion,
		LastSeenAt: device.LastSeenAt.Format(timeFormat),
		DisabledAt: disabledAt,
	}
}

func notificationPreferenceResponse(preference *models.NotificationPreference) NotificationPreferenceResponse {
	return NotificationPreferenceResponse{
		PersonaID:        preference.PersonaID.String(),
		NotificationType: preference.NotificationType,
		PushEnabled:      preference.PushEnabled,
		UpdatedAt:        preference.UpdatedAt.Format(timeFormat),
	}
}

func notificationPreferenceTypes() []models.NotificationType {
	return []models.NotificationType{
		models.FollowNotificationType,
		models.LikeNotificationType,
		models.CommentNotificationType,
		models.DirectMessageNotificationType,
		models.GroupMessageNotificationType,
		models.StoryContributionRequestNotificationType,
		models.StoryContributionAcceptedNotificationType,
		models.StoryContributionRejectedNotificationType,
		models.StoryHighlightAddedNotificationType,
		models.StoryHighlightRemovedNotificationType,
		models.StoryReactionNotificationType,
	}
}
