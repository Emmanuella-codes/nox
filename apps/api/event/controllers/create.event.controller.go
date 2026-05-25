package controllers

import (
	"github.com/emmanuella-codes/nox/event/dtos"
	"github.com/emmanuella-codes/nox/middleware"
	"github.com/gofiber/fiber/v2"
)

func (c *EventController) CreateEvent(ctx *fiber.Ctx) error {
	userID, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}

	var dto dtos.CreateEventDTO
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}

	res := c.pipe.CreateEventPipe(ctx.Context(), userID, dto)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusCreated, res.Message, res.Data)
}
