package controllers

import (
	"github.com/emmanuella-codes/nox/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (c *CommentController) DeleteComment(ctx *fiber.Ctx) error {
	userID, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	commentID, err := uuid.Parse(ctx.Params("commentID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_comment_id")
	}

	res := c.pipe.DeleteCommentPipe(ctx.Context(), userID, commentID)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess[any](ctx, fiber.StatusOK, res.Message, nil)
}
