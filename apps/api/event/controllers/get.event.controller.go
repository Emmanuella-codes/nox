package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (c *EventController) GetEvent(ctx *fiber.Ctx) error {
	eventID, err := uuid.Parse(ctx.Params("eventID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_event_id")
	}

	res := c.pipe.GetEventPipe(ctx.Context(), eventID)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
}
