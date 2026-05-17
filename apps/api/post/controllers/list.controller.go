package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (c *PostController) GetPost(ctx *fiber.Ctx) error {
	postID, err := uuid.Parse(ctx.Params("postID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_post_id")
	}

	res := c.pipe.GetPostPipe(ctx.Context(), postID)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}

	return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
}

func (c *PostController) GetPersonaPosts(ctx *fiber.Ctx) error {
	personaID, err := uuid.Parse(ctx.Params("personaID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_persona_id")
	}

	res := c.pipe.GetPersonaPostsPipe(ctx.Context(), personaID, queryLimit(ctx, 20))
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}

	return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
}

func (c *PostController) GetFeed(ctx *fiber.Ctx) error {
	personaID, err := uuid.Parse(ctx.Params("personaID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_persona_id")
	}

	res := c.pipe.GetFeedPipe(ctx.Context(), personaID, queryLimit(ctx, 20))
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}

	return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
}
