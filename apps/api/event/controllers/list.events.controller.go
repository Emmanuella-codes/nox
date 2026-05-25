package controllers

import "github.com/gofiber/fiber/v2"

func (c *EventController) ListEvents(ctx *fiber.Ctx) error {
	res := c.pipe.ListEventsPipe(ctx.Context(), queryLimit(ctx, 30))
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
}
