package controllers

import (
	"github.com/emmanuella-codes/nox/auth/dtos"
	"github.com/gofiber/fiber/v2"
)

func (c *AuthController) Refresh(ctx *fiber.Ctx) error {
	var dto dtos.RefreshDTO
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}

	res := c.pipe.RefreshPipe(ctx.Context(), dto.RefreshToken)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}

	return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
}
