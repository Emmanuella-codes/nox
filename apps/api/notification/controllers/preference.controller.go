package controllers

import (
	"github.com/emmanuella-codes/nox/middleware"
	"github.com/emmanuella-codes/nox/models"
	notification_dtos "github.com/emmanuella-codes/nox/notification/dtos"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (c *NotificationController) ListNotificationPreferences(ctx *fiber.Ctx) error {
	userID, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	personaID, err := uuid.Parse(ctx.Query("persona_id"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_persona_id")
	}
	res := c.pipe.ListNotificationPreferencesPipe(ctx.Context(), userID, personaID)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
}

func (c *NotificationController) UpdateNotificationPreference(ctx *fiber.Ctx) error {
	userID, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	var dto notification_dtos.UpdateNotificationPreferenceDTO
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}
	res := c.pipe.UpdateNotificationPreferencePipe(ctx.Context(), userID, dto.PersonaID, models.NotificationType(dto.NotificationType), dto.PushEnabled)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
}
