package controllers

import (
	"github.com/emmanuella-codes/nox/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (c *SetController) DeleteSet(ctx *fiber.Ctx) error {
	userID, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}

	setID, err := uuid.Parse(ctx.Params("setID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_set_id")
	}

	res := c.pipe.DeleteSetPipe(ctx.Context(), userID, setID)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
}
