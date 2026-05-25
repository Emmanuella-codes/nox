package controllers

import "github.com/gofiber/fiber/v2"

func (c *SearchController) Search(ctx *fiber.Ctx) error {
	res := c.pipe.Search(ctx.Context(), ctx.Query("q"), queryLimit(ctx, 10))
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
}
