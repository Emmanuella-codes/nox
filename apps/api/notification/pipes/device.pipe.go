package pipes

import (
	"context"

	notification_dtos "github.com/emmanuella-codes/nox/notification/dtos"
	"github.com/emmanuella-codes/nox/notification/messages"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

func (p *NotificationPipe) UpsertNotificationDevicePipe(ctx context.Context, userID uuid.UUID, dto notification_dtos.UpsertNotificationDeviceDTO) *shared.PipeRes[NotificationDeviceResponse] {
	device, err := p.notificationRepo.UpsertNotificationDevice(ctx, userID, dto)
	if err != nil {
		return pipeInternalError[NotificationDeviceResponse](err, "notification.upsert_device")
	}
	response := notificationDeviceResponse(device)
	return shared.PipeSuccess(messages.Notification_Device_Upserted, &response)
}

func (p *NotificationPipe) RemoveNotificationDevicePipe(ctx context.Context, userID uuid.UUID, deviceID uuid.UUID) *shared.PipeRes[any] {
	if err := p.notificationRepo.DisableNotificationDevice(ctx, userID, deviceID); err != nil {
		return pipeInternalError[any](err, "notification.remove_device")
	}
	return shared.PipeSuccess[any](messages.Notification_Device_Removed, nil)
}
