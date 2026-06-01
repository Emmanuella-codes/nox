package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (c *SetController) GetSet(ctx *fiber.Ctx) error {
	setID, err := uuid.Parse(ctx.Params("setID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_set_id")
	}

	res := c.pipe.GetSetPipe(ctx.Context(), setID)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
}
