package controllers

import (
	"github.com/emmanuella-codes/nox/middleware"
	"github.com/gofiber/fiber/v2"
)

func (ac *AdminController) Dashboard(ctx *fiber.Ctx) error {
	adminUserID, ok := middleware.CurrentAdminUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	res := ac.pipe.DashboardPipe(requestContext(ctx), adminUserID)
	if res.Success {
		return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
	}
	return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
}
