package controllers

import (
	"github.com/emmanuella-codes/nox/admin/dtos"
	"github.com/gofiber/fiber/v2"
)

func (ac *AdminController) Login(ctx *fiber.Ctx) error {
	dto := dtos.LoginDTO{}
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}

	res := ac.pipe.LoginPipe(requestContext(ctx), dto)
	if res.Success {
		return pipeSuccess(ctx, 200, res.Message, res.Data)
	}

	return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
}
