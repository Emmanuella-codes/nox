package controllers

import (
	"github.com/emmanuella-codes/nox/middleware"
	notification_dtos "github.com/emmanuella-codes/nox/notification/dtos"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (c *NotificationController) UpsertNotificationDevice(ctx *fiber.Ctx) error {
	userID, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	var dto notification_dtos.UpsertNotificationDeviceDTO
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}
	res := c.pipe.UpsertNotificationDevicePipe(ctx.Context(), userID, dto)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusCreated, res.Message, res.Data)
}

func (c *NotificationController) ListNotificationDevices(ctx *fiber.Ctx) error {
	userID, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	res := c.pipe.ListNotificationDevicesPipe(ctx.Context(), userID)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
}

func (c *NotificationController) RemoveNotificationDevice(ctx *fiber.Ctx) error {
	userID, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	deviceID, err := uuid.Parse(ctx.Params("deviceID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_device_id")
	}
	res := c.pipe.RemoveNotificationDevicePipe(ctx.Context(), userID, deviceID)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess[any](ctx, fiber.StatusOK, res.Message, nil)
}
