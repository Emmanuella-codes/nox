package controllers

import (
	"github.com/emmanuella-codes/nox/like/dtos"
	"github.com/emmanuella-codes/nox/middleware"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (c *LikeController) LikePost(ctx *fiber.Ctx) error {
	return c.likeAction(ctx, true)
}

func (c *LikeController) UnlikePost(ctx *fiber.Ctx) error {
	return c.likeAction(ctx, false)
}

func (c *LikeController) likeAction(ctx *fiber.Ctx, like bool) error {
	userID, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	postID, err := uuid.Parse(ctx.Params("postID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_post_id")
	}

	var dto dtos.LikePostDTO
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}

	var res *shared.PipeRes[any]
	if like {
		res = c.pipe.LikePostPipe(ctx.Context(), userID, postID, dto)
	} else {
		res = c.pipe.UnlikePostPipe(ctx.Context(), userID, postID, dto)
	}
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess[any](ctx, fiber.StatusOK, res.Message, nil)
}
