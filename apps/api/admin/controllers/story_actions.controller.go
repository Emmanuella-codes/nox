package controllers

import (
	"strconv"
	"strings"

	"github.com/emmanuella-codes/nox/admin/dtos"
	"github.com/emmanuella-codes/nox/middleware"
	"github.com/emmanuella-codes/nox/models"
	adminrepo "github.com/emmanuella-codes/nox/repositories/admin"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (ac *AdminController) RestoreSet(ctx *fiber.Ctx) error {
	return ac.fixedModeration(ctx, models.ModerationEntitySet, models.ModerationStatusActive, false)
}

func (ac *AdminController) RemoveSet(ctx *fiber.Ctx) error {
	return ac.fixedModeration(ctx, models.ModerationEntitySet, models.ModerationStatusRemoved, true)
}

func (ac *AdminController) RestoreStory(ctx *fiber.Ctx) error {
	return ac.fixedModeration(ctx, models.ModerationEntityStory, models.ModerationStatusActive, false)
}

func (ac *AdminController) RemoveStory(ctx *fiber.Ctx) error {
	return ac.fixedModeration(ctx, models.ModerationEntityStory, models.ModerationStatusRemoved, true)
}

func (ac *AdminController) RestoreStoryItem(ctx *fiber.Ctx) error {
	return ac.fixedModeration(ctx, models.ModerationEntityStoryItem, models.ModerationStatusActive, false)
}

func (ac *AdminController) RemoveStoryItem(ctx *fiber.Ctx) error {
	return ac.fixedModeration(ctx, models.ModerationEntityStoryItem, models.ModerationStatusRemoved, true)
}

func (ac *AdminController) fixedModeration(ctx *fiber.Ctx, entityType models.ModerationEntityType, status models.ModerationStatus, requiresReason bool) error {
	adminUserID, ok := middleware.CurrentAdminUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	entityID, err := uuid.Parse(ctx.Params("entityID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_payload")
	}
	dto := dtos.ModerateContentDTO{Status: status}
	if ctx.Request().Header.ContentLength() > 0 {
		if err := ctx.BodyParser(&dto); err != nil {
			return validationError(ctx, err)
		}
	}
	dto.Status = status
	dto.Reason = strings.TrimSpace(dto.Reason)
	if requiresReason && dto.Reason == "" {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_payload")
	}
	res := ac.pipe.ModerateContentPipe(requestContext(ctx), adminUserID, entityType, entityID, dto)
	if res.Success {
		return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
	}
	return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
}

func (ac *AdminController) ListFeaturedSets(ctx *fiber.Ctx) error {
	adminUserID, ok := middleware.CurrentAdminUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	params := adminrepo.ListFeaturedSetsParams{}
	params.Limit, _ = strconv.Atoi(ctx.Query("limit", "50"))
	params.Offset, _ = strconv.Atoi(ctx.Query("offset", "0"))
	res := ac.pipe.ListFeaturedSetsPipe(requestContext(ctx), adminUserID, params)
	if res.Success {
		return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
	}
	return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
}

func (ac *AdminController) FeatureSet(ctx *fiber.Ctx) error {
	adminUserID, ok := middleware.CurrentAdminUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	setID, err := uuid.Parse(ctx.Params("setID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_payload")
	}
	var dto dtos.FeatureSetDTO
	if ctx.Request().Header.ContentLength() > 0 {
		if err := parseAndValidate(ctx, &dto); err != nil {
			return validationError(ctx, err)
		}
	}
	res := ac.pipe.FeatureSetPipe(requestContext(ctx), adminUserID, setID, dto)
	if res.Success {
		return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
	}
	return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
}

func (ac *AdminController) UnfeatureSet(ctx *fiber.Ctx) error {
	adminUserID, ok := middleware.CurrentAdminUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	setID, err := uuid.Parse(ctx.Params("setID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_payload")
	}
	res := ac.pipe.UnfeatureSetPipe(requestContext(ctx), adminUserID, setID)
	if res.Success {
		return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
	}
	return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
}

func (ac *AdminController) ListStoryContributions(ctx *fiber.Ctx) error {
	adminUserID, ok := middleware.CurrentAdminUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	params := adminrepo.ListStoryContributionsParams{}
	if value := ctx.Query("status"); value != "" {
		status := models.StoryContributionRequestStatus(value)
		params.Status = &status
	}
	params.Limit, _ = strconv.Atoi(ctx.Query("limit", "50"))
	params.Offset, _ = strconv.Atoi(ctx.Query("offset", "0"))
	res := ac.pipe.ListStoryContributionsPipe(requestContext(ctx), adminUserID, params)
	if res.Success {
		return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
	}
	return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
}

func (ac *AdminController) ReviewStoryContribution(ctx *fiber.Ctx) error {
	adminUserID, ok := middleware.CurrentAdminUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	requestID, err := uuid.Parse(ctx.Params("requestID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_payload")
	}
	var dto dtos.ReviewStoryContributionDTO
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}
	res := ac.pipe.ReviewStoryContributionPipe(requestContext(ctx), adminUserID, requestID, dto)
	if res.Success {
		return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
	}
	return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
}

func (ac *AdminController) ListAdminHighlights(ctx *fiber.Ctx) error {
	adminUserID, ok := middleware.CurrentAdminUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	params := adminrepo.ListAdminHighlightsParams{HighlightType: ctx.Query("type")}
	params.Limit, _ = strconv.Atoi(ctx.Query("limit", "50"))
	params.Offset, _ = strconv.Atoi(ctx.Query("offset", "0"))
	res := ac.pipe.ListAdminHighlightsPipe(requestContext(ctx), adminUserID, params)
	if res.Success {
		return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
	}
	return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
}

func (ac *AdminController) RemoveAdminHighlight(ctx *fiber.Ctx) error {
	adminUserID, ok := middleware.CurrentAdminUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	highlightID, err := uuid.Parse(ctx.Params("highlightID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_payload")
	}
	var body struct {
		Reason string `json:"reason"`
	}
	if err := ctx.BodyParser(&body); err != nil {
		return validationError(ctx, err)
	}
	res := ac.pipe.RemoveAdminHighlightPipe(requestContext(ctx), adminUserID, highlightID, ctx.Params("highlightType"), body.Reason)
	if res.Success {
		return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
	}
	return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
}

func (ac *AdminController) GetReportScopedPrivateContent(ctx *fiber.Ctx) error {
	adminUserID, ok := middleware.CurrentAdminUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	reportID, err := uuid.Parse(ctx.Params("reportID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_payload")
	}
	res := ac.pipe.GetReportScopedPrivateContentPipe(requestContext(ctx), adminUserID, reportID)
	if res.Success {
		return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
	}
	return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
}
