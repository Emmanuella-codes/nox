package controllers

import (
	"github.com/emmanuella-codes/nox/admin/dtos"
	"github.com/emmanuella-codes/nox/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (ac *AdminController) ListAdminUsers(ctx *fiber.Ctx) error {
	adminUserID, ok := middleware.CurrentAdminUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}

	res := ac.pipe.ListAdminUsersPipe(requestContext(ctx), adminUserID)
	if res.Success {
		return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
	}

	return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
}

func (ac *AdminController) GetAdminUser(ctx *fiber.Ctx) error {
	adminUserID, ok := middleware.CurrentAdminUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}

	targetUserID, err := uuid.Parse(ctx.Params("userID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_payload")
	}

	res := ac.pipe.GetAdminUserPipe(requestContext(ctx), adminUserID, targetUserID)
	if res.Success {
		return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
	}

	return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
}

func (ac *AdminController) CreateAdminUser(ctx *fiber.Ctx) error {
	adminUserID, ok := middleware.CurrentAdminUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}

	dto := dtos.CreateAdminUserDTO{}
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}

	res := ac.pipe.CreateAdminUserPipe(requestContext(ctx), adminUserID, dto)
	if res.Success {
		return pipeSuccess(ctx, fiber.StatusCreated, res.Message, res.Data)
	}

	return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
}

func (ac *AdminController) UpdateAdminUserRole(ctx *fiber.Ctx) error {
	adminUserID, ok := middleware.CurrentAdminUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}

	targetUserID, err := uuid.Parse(ctx.Params("userID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_payload")
	}

	dto := dtos.UpdateAdminRoleDTO{}
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}

	res := ac.pipe.UpdateAdminUserRolePipe(requestContext(ctx), adminUserID, targetUserID, dto)
	if res.Success {
		return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
	}

	return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
}

func (ac *AdminController) UpdateAdminUserStatus(ctx *fiber.Ctx) error {
	adminUserID, ok := middleware.CurrentAdminUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}

	targetUserID, err := uuid.Parse(ctx.Params("userID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_payload")
	}

	dto := dtos.UpdateAdminStatusDTO{}
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}

	res := ac.pipe.UpdateAdminUserStatusPipe(requestContext(ctx), adminUserID, targetUserID, dto)
	if res.Success {
		return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
	}

	return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
}
