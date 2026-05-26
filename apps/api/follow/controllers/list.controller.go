package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (c *FollowController) GetFollowers(ctx *fiber.Ctx) error {
	personaID, err := uuid.Parse(ctx.Params("personaID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_persona_id")
	}

	res := c.pipe.FollowersPipe(ctx.Context(), personaID, queryLimit(ctx, 20))
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}

	return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
}

func (c *FollowController) GetFollowing(ctx *fiber.Ctx) error {
	personaID, err := uuid.Parse(ctx.Params("personaID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_persona_id")
	}

	res := c.pipe.FollowingPipe(ctx.Context(), personaID, queryLimit(ctx, 20))
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}

	return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
}
