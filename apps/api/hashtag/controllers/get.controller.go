package controllers

import "github.com/gofiber/fiber/v2"

func (c *HashtagController) GetHashtag(ctx *fiber.Ctx) error {
	res := c.pipe.GetHashtagPipe(ctx.Context(), ctx.Params("tag"))
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
}
