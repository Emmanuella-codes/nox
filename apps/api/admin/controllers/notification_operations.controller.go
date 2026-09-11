package controllers

import (
	"strconv"

	"github.com/emmanuella-codes/nox/admin/dtos"
	"github.com/emmanuella-codes/nox/admin/messages"
	"github.com/emmanuella-codes/nox/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (ac *AdminController) ListNotificationDevices(ctx *fiber.Ctx) error {
	adminID, ok := middleware.CurrentAdminUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	userID, err := optionalUUID(ctx.Query("user_id"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, messages.Invalid_Payload)
	}
	limit, _ := strconv.Atoi(ctx.Query("limit", "50"))
	offset, _ := strconv.Atoi(ctx.Query("offset", "0"))
	res := ac.pipe.ListNotificationDevicesPipe(requestContext(ctx), adminID, userID, limit, offset)
	if res.Success {
		return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
	}
	return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
}

func (ac *AdminController) ListNotificationOutbox(ctx *fiber.Ctx) error {
	adminID, ok := middleware.CurrentAdminUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	userID, err := optionalUUID(ctx.Query("user_id"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, messages.Invalid_Payload)
	}
	limit, _ := strconv.Atoi(ctx.Query("limit", "50"))
	offset, _ := strconv.Atoi(ctx.Query("offset", "0"))
	res := ac.pipe.ListNotificationOutboxPipe(requestContext(ctx), adminID, ctx.Query("status"), userID, limit, offset)
	if res.Success {
		return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
	}
	return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
}

func (ac *AdminController) RetryNotificationOutbox(ctx *fiber.Ctx) error {
	adminID, ok := middleware.CurrentAdminUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	outboxID, err := uuid.Parse(ctx.Params("outboxID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, messages.Invalid_Payload)
	}
	var dto dtos.NotificationAdminActionDTO
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}
	res := ac.pipe.RetryNotificationOutboxPipe(requestContext(ctx), adminID, outboxID, dto)
	if res.Success {
		return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
	}
	return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
}

func (ac *AdminController) DisableNotificationDevice(ctx *fiber.Ctx) error {
	adminID, ok := middleware.CurrentAdminUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	deviceID, err := uuid.Parse(ctx.Params("deviceID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, messages.Invalid_Payload)
	}
	var dto dtos.NotificationAdminActionDTO
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}
	res := ac.pipe.DisableNotificationDevicePipe(requestContext(ctx), adminID, deviceID, dto)
	if res.Success {
		return pipeSuccess[any](ctx, fiber.StatusOK, res.Message, nil)
	}
	return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
}
