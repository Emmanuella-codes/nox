package controllers

import (
	"github.com/emmanuella-codes/nox/auth/dtos"
	"github.com/gofiber/fiber/v2"
)

func (c *AuthController) Login(ctx *fiber.Ctx) error {
	var dto dtos.LoginDTO
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}

	res := c.pipe.LoginPipe(ctx.Context(), dto)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}

	return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
}
