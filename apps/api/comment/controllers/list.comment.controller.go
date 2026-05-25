package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (c *CommentController) ListComments(ctx *fiber.Ctx) error {
	postID, err := uuid.Parse(ctx.Params("postID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_post_id")
	}

	res := c.pipe.ListCommentsPipe(ctx.Context(), postID, queryLimit(ctx, 50))
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
}
