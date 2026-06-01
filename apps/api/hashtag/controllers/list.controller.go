package controllers

import "github.com/gofiber/fiber/v2"

func (c *HashtagController) Trending(ctx *fiber.Ctx) error {
	res := c.pipe.TrendingPipe(ctx.Context(), queryLimit(ctx, 20))
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
}
