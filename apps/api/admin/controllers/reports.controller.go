package controllers

import (
	"strconv"

	"github.com/emmanuella-codes/nox/admin/dtos"
	"github.com/emmanuella-codes/nox/middleware"
	"github.com/emmanuella-codes/nox/models"
	adminrepo "github.com/emmanuella-codes/nox/repositories/admin"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (ac *AdminController) ListReports(ctx *fiber.Ctx) error {
	adminUserID, ok := middleware.CurrentAdminUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	params := adminrepo.ListReportsParams{}
	if value := ctx.Query("status"); value != "" {
		status := models.ReportStatus(value)
		params.Status = &status
	}
	if value := ctx.Query("target_type"); value != "" {
		target := models.ReportTargetType(value)
		params.TargetType = &target
	}
	if value := ctx.Query("assigned_admin_id"); value != "" {
		assigned, err := uuid.Parse(value)
		if err != nil {
			return pipeError(ctx, fiber.StatusBadRequest, "invalid_payload")
		}
		params.AssignedID = &assigned
	}
	params.Limit, _ = strconv.Atoi(ctx.Query("limit", "50"))
	params.Offset, _ = strconv.Atoi(ctx.Query("offset", "0"))
	res := ac.pipe.ListReportsPipe(requestContext(ctx), adminUserID, params)
	if res.Success {
		return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
	}
	return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
}

func (ac *AdminController) GetReport(ctx *fiber.Ctx) error {
	adminUserID, ok := middleware.CurrentAdminUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	reportID, err := uuid.Parse(ctx.Params("reportID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_payload")
	}
	res := ac.pipe.GetReportPipe(requestContext(ctx), adminUserID, reportID)
	if res.Success {
		return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
	}
	return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
}

func (ac *AdminController) UpdateReport(ctx *fiber.Ctx) error {
	adminUserID, ok := middleware.CurrentAdminUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	reportID, err := uuid.Parse(ctx.Params("reportID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_payload")
	}
	var dto dtos.UpdateReportDTO
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}
	res := ac.pipe.UpdateReportPipe(requestContext(ctx), adminUserID, reportID, dto)
	if res.Success {
		return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
	}
	return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
}

func (ac *AdminController) ResolveReport(ctx *fiber.Ctx) error {
	adminUserID, ok := middleware.CurrentAdminUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	reportID, err := uuid.Parse(ctx.Params("reportID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_payload")
	}
	var dto dtos.ReportActionDTO
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}
	res := ac.pipe.ResolveReportPipe(requestContext(ctx), adminUserID, reportID, dto)
	if res.Success {
		return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
	}
	return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
}
