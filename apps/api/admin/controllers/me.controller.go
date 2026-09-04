package controllers

import (
	"github.com/emmanuella-codes/nox/middleware"
	"github.com/gofiber/fiber/v2"
)

func (ac *AdminController) Me(ctx *fiber.Ctx) error {
	userID, ok := middleware.CurrentAdminUserID(ctx)
	if !ok {
		return pipeError(ctx, 401, "invalid_token")
	}

	res := ac.pipe.MePipe(requestContext(ctx), userID)
	if res.Success {
		return pipeSuccess(ctx, 200, res.Message, res.Data)
	}

	return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
}
