package controllers

import (
	"github.com/emmanuella-codes/nox/auth/dtos"
	"github.com/gofiber/fiber/v2"
)

func (c *AuthController) ResendVerification(ctx *fiber.Ctx) error {
	var dto dtos.ResendVerificationDTO
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}

	res := c.pipe.ResendVerificationPipe(ctx.Context(), dto)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}

	return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
}
