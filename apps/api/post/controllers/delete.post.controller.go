package controllers

import (
	"github.com/emmanuella-codes/nox/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (c *PostController) DeletePost(ctx *fiber.Ctx) error {
	userID, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}

	postID, err := uuid.Parse(ctx.Params("postID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_post_id")
	}

	res := c.pipe.DeletePostPipe(ctx.Context(), userID, postID)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}

	return pipeSuccess[any](ctx, fiber.StatusOK, res.Message, nil)
}
