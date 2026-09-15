package controllers

import (
	"github.com/emmanuella-codes/nox/middleware"
	"github.com/emmanuella-codes/nox/set/dtos"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (c *SetController) LikeSet(ctx *fiber.Ctx) error {
	userID, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	setID, err := uuid.Parse(ctx.Params("setID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_set_id")
	}
	var dto dtos.SetPersonaActionDTO
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}
	res := c.pipe.LikeSetPipe(ctx.Context(), userID, setID, dto)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess[any](ctx, fiber.StatusOK, res.Message, nil)
}

func (c *SetController) UnlikeSet(ctx *fiber.Ctx) error {
	userID, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	setID, err := uuid.Parse(ctx.Params("setID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_set_id")
	}
	var dto dtos.SetPersonaActionDTO
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}
	res := c.pipe.UnlikeSetPipe(ctx.Context(), userID, setID, dto)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess[any](ctx, fiber.StatusOK, res.Message, nil)
}

func (c *SetController) RecordSetPlay(ctx *fiber.Ctx) error {
	userID, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	setID, err := uuid.Parse(ctx.Params("setID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_set_id")
	}
	res := c.pipe.RecordSetPlayPipe(ctx.Context(), userID, setID)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess[any](ctx, fiber.StatusOK, res.Message, nil)
}

func (c *SetController) CreateSetComment(ctx *fiber.Ctx) error {
	userID, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	setID, err := uuid.Parse(ctx.Params("setID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_set_id")
	}
	var dto dtos.CreateSetCommentDTO
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}
	res := c.pipe.CreateSetCommentPipe(ctx.Context(), userID, setID, dto)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusCreated, res.Message, res.Data)
}

func (c *SetController) ListSetComments(ctx *fiber.Ctx) error {
	setID, err := uuid.Parse(ctx.Params("setID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_set_id")
	}
	res := c.pipe.ListSetCommentsPipe(ctx.Context(), setID, queryLimit(ctx, 20), queryOffset(ctx))
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
}

func (c *SetController) UpdateSetComment(ctx *fiber.Ctx) error {
	userID, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	commentID, err := uuid.Parse(ctx.Params("commentID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_comment_id")
	}
	var dto dtos.UpdateSetCommentDTO
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}
	res := c.pipe.UpdateSetCommentPipe(ctx.Context(), userID, commentID, dto)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
}

func (c *SetController) DeleteSetComment(ctx *fiber.Ctx) error {
	userID, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	commentID, err := uuid.Parse(ctx.Params("commentID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_comment_id")
	}
	personaID, err := uuid.Parse(ctx.Query("persona_id"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_persona_id")
	}
	res := c.pipe.DeleteSetCommentPipe(ctx.Context(), userID, commentID, personaID)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess[any](ctx, fiber.StatusOK, res.Message, nil)
}

func (c *SetController) LikeSetComment(ctx *fiber.Ctx) error {
	return c.setCommentLike(ctx, false)
}

func (c *SetController) UnlikeSetComment(ctx *fiber.Ctx) error {
	return c.setCommentLike(ctx, true)
}

func (c *SetController) setCommentLike(ctx *fiber.Ctx, unlike bool) error {
	userID, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	commentID, err := uuid.Parse(ctx.Params("commentID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_comment_id")
	}
	personaID, err := uuid.Parse(ctx.Query("persona_id"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_persona_id")
	}
	res := c.pipe.SetCommentLikePipe(ctx.Context(), userID, commentID, personaID, unlike)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess[any](ctx, fiber.StatusOK, res.Message, nil)
}
