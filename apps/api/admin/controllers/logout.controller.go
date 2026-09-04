package controllers

import (
	"github.com/emmanuella-codes/nox/admin/dtos"
	"github.com/gofiber/fiber/v2"
)

func (ac *AdminController) Logout(ctx *fiber.Ctx) error {
	dto := dtos.RefreshDTO{}
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}

	res := ac.pipe.LogoutPipe(requestContext(ctx), dto.RefreshToken)
	if res.Success {
		return pipeSuccess[any](ctx, 200, res.Message, nil)
	}

	return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
}
