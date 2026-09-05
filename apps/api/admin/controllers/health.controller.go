package controllers

import "github.com/gofiber/fiber/v2"

func (ac *AdminController) Health(ctx *fiber.Ctx) error {
	res := ac.pipe.HealthPipe(requestContext(ctx))
	if res.Success {
		return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
	}

	return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
}
