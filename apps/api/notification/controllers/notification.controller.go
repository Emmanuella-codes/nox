package controllers

import (
	"github.com/emmanuella-codes/nox/middleware"
	"github.com/emmanuella-codes/nox/notification/dtos"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// ListNotifications lists notifications for one owned persona.
func (c *NotificationController) ListNotifications(ctx *fiber.Ctx) error {
	userID, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	personaID, err := uuid.Parse(ctx.Query("persona_id"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_persona_id")
	}
	res := c.pipe.ListNotificationsPipe(ctx.Context(), userID, personaID, queryLimit(ctx, 20), queryOffset(ctx))
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
}

// GetUnreadCount returns the unread notification count for one owned persona.
func (c *NotificationController) GetUnreadCount(ctx *fiber.Ctx) error {
	userID, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	personaID, err := uuid.Parse(ctx.Query("persona_id"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_persona_id")
	}
	res := c.pipe.GetUnreadCountPipe(ctx.Context(), userID, personaID)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
}

// MarkNotificationRead marks one notification as read.
func (c *NotificationController) MarkNotificationRead(ctx *fiber.Ctx) error {
	userID, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	notificationID, err := uuid.Parse(ctx.Params("notificationID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_notification_id")
	}
	var dto dtos.MarkNotificationReadDTO
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}
	res := c.pipe.MarkNotificationReadPipe(ctx.Context(), userID, notificationID, dto.PersonaID)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
}

// MarkAllNotificationsRead marks all notifications as read for one persona.
func (c *NotificationController) MarkAllNotificationsRead(ctx *fiber.Ctx) error {
	userID, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	var dto dtos.MarkAllNotificationsReadDTO
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}
	res := c.pipe.MarkAllNotificationsReadPipe(ctx.Context(), userID, dto.PersonaID)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess[any](ctx, fiber.StatusOK, res.Message, nil)
}
