package controllers

import (
	"github.com/emmanuella-codes/nox/admin/dtos"
	"github.com/emmanuella-codes/nox/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (ac *AdminController) ListUsers(ctx *fiber.Ctx) error {
	adminUserID, ok := middleware.CurrentAdminUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}

	res := ac.pipe.ListUsersPipe(requestContext(ctx), adminUserID)
	if res.Success {
		return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
	}

	return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
}

func (ac *AdminController) GetUser(ctx *fiber.Ctx) error {
	adminUserID, ok := middleware.CurrentAdminUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}

	targetUserID, err := uuid.Parse(ctx.Params("userID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_payload")
	}

	res := ac.pipe.GetUserPipe(requestContext(ctx), adminUserID, targetUserID)
	if res.Success {
		return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
	}

	return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
}

func (ac *AdminController) UpdateUserStatus(ctx *fiber.Ctx) error {
	adminUserID, ok := middleware.CurrentAdminUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}

	targetUserID, err := uuid.Parse(ctx.Params("userID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_payload")
	}

	dto := dtos.UpdateUserStatusDTO{}
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}

	res := ac.pipe.UpdateUserStatusPipe(requestContext(ctx), adminUserID, targetUserID, dto.Status)
	if res.Success {
		return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
	}

	return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
}

func (ac *AdminController) RevokeUserSessions(ctx *fiber.Ctx) error {
	adminUserID, ok := middleware.CurrentAdminUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}

	targetUserID, err := uuid.Parse(ctx.Params("userID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_payload")
	}

	res := ac.pipe.RevokeUserSessionsPipe(requestContext(ctx), adminUserID, targetUserID)
	if res.Success {
		return pipeSuccess[any](ctx, fiber.StatusOK, res.Message, nil)
	}

	return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
}

func (ac *AdminController) ResendUserVerification(ctx *fiber.Ctx) error {
	adminUserID, ok := middleware.CurrentAdminUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}

	targetUserID, err := uuid.Parse(ctx.Params("userID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_payload")
	}

	res := ac.pipe.ResendUserVerificationPipe(requestContext(ctx), adminUserID, targetUserID)
	if res.Success {
		return pipeSuccess[any](ctx, fiber.StatusOK, res.Message, nil)
	}

	return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
}

func (ac *AdminController) MarkUserEmailVerified(ctx *fiber.Ctx) error {
	adminUserID, ok := middleware.CurrentAdminUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}

	targetUserID, err := uuid.Parse(ctx.Params("userID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_payload")
	}

	res := ac.pipe.MarkUserEmailVerifiedPipe(requestContext(ctx), adminUserID, targetUserID)
	if res.Success {
		return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
	}

	return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
}
