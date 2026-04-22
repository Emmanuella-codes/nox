package controllers

import (
	"github.com/emmanuella-codes/nox/auth/dtos"
	"github.com/gofiber/fiber/v2"
)

func (c *AuthController) VerifyEmail(ctx *fiber.Ctx) error {
	var dto dtos.VerifyEmailDTO
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}

	res := c.pipe.VerifyEmailPipe(ctx.Context(), dto)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}

	return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
}
